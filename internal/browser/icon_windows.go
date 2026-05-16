//go:build windows

package browser

import (
	"debug/pe"
	"encoding/base64"
	"image/png"
	"image"
	"image/color"
	"encoding/binary"
	"bytes"
	"os"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modShell32 = windows.NewLazySystemDLL("shell32.dll")
	modUser32  = windows.NewLazySystemDLL("user32.dll")
	modGdi32   = windows.NewLazySystemDLL("gdi32.dll")

	procExtractIconExW       = modShell32.NewProc("ExtractIconExW")
	procPrivateExtractIconsW = modShell32.NewProc("PrivateExtractIconsW")
	procDestroyIcon          = modUser32.NewProc("DestroyIcon")
	procGetDC                = modUser32.NewProc("GetDC")
	procReleaseDC            = modUser32.NewProc("ReleaseDC")
	procDrawIconEx           = modUser32.NewProc("DrawIconEx")
	procCreateCompatibleDC   = modGdi32.NewProc("CreateCompatibleDC")
	procCreateDIBSection     = modGdi32.NewProc("CreateDIBSection")
	procSelectObject         = modGdi32.NewProc("SelectObject")
	procDeleteDC             = modGdi32.NewProc("DeleteDC")
	procDeleteObject         = modGdi32.NewProc("DeleteObject")

	pngHeaderBytes = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
)

type resEntry struct {
	nameOrID uint32
	offset   uint32
}

func bmpToImage(data []byte) *image.RGBA {
	if len(data) < 40 || binary.LittleEndian.Uint32(data[:4]) < 40 {
		return nil
	}
	w := int(binary.LittleEndian.Uint32(data[4:8]))
	h := binary.LittleEndian.Uint32(data[8:12])
	if binary.LittleEndian.Uint16(data[14:16]) != 32 {
		return nil
	}
	topDown := h < 0
	if h < 0 {
		h = -h
	}
	img := image.NewRGBA(image.Rect(0, 0, w, int(h)))
	pixels := data[40:]
	iw := uint32(w)
	for y := uint32(0); y < h; y++ {
		row := y
		if !topDown {
			row = h - 1 - y
		}
		for x := uint32(0); x < iw; x++ {
			s := (row*iw + x) * 4
			if s+4 > uint32(len(pixels)) {
				break
			}
			dst := int(row)*img.Stride + int(x)*4
			img.Pix[dst+0] = pixels[s+2]
			img.Pix[dst+1] = pixels[s+1]
			img.Pix[dst+2] = pixels[s+0]
			img.Pix[dst+3] = pixels[s+3]
		}
	}
	return img
}

func encodeIconData(bmpData []byte) string {
	img := bmpToImage(bmpData)
	if img == nil {
		return ""
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return ""
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}

func extractICOFile(path string) string {
	raw, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	if len(raw) < 6 || binary.LittleEndian.Uint16(raw[2:4]) != 1 {
		return ""
	}
	count := int(binary.LittleEndian.Uint16(raw[4:6]))
	dirEnd := 6 + count*16
	if count == 0 || len(raw) < dirEnd {
		return ""
	}

	bestIdx := -1
	bestSize := 0
	for i := 0; i < count; i++ {
		off := 6 + i*16
		w := int(raw[off])
		if w == 0 {
			w = 256
		}
		if w > bestSize {
			bestSize = w
			bestIdx = i
		}
	}
	if bestIdx < 0 {
		return ""
	}
	eOff := 6 + bestIdx*16
	imgOff := int(binary.LittleEndian.Uint32(raw[eOff+12 : eOff+16]))
	imgLen := int(binary.LittleEndian.Uint32(raw[eOff+8 : eOff+12]))
	if imgOff+imgLen > len(raw) {
		return ""
	}
	imgData := raw[imgOff : imgOff+imgLen]
	if bytes.HasPrefix(imgData, pngHeaderBytes) {
		return "data:image/png;base64," + base64.StdEncoding.EncodeToString(imgData)
	}
	return encodeIconData(imgData)
}

func extractResourceIcons(path string) string {
	f, err := pe.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	var rva, _ uint32
	switch oh := f.OptionalHeader.(type) {
	case *pe.OptionalHeader64:
		d := oh.DataDirectory[pe.IMAGE_DIRECTORY_ENTRY_RESOURCE]
		rva, _ = d.VirtualAddress, d.Size
	case *pe.OptionalHeader32:
		d := oh.DataDirectory[pe.IMAGE_DIRECTORY_ENTRY_RESOURCE]
		rva, _ = d.VirtualAddress, d.Size
	default:
		return ""
	}
	if rva == 0 {
		return ""
	}

	fileRaw, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	r := &peReader{raw: fileRaw, sections: f.Sections, baseRVA: rva}

	rootEntries := r.readDir(0)
	l1Off := uint32(0)
	for _, e := range rootEntries {
		if e.nameOrID&0x80000000 == 0 && e.nameOrID == 3 {
			l1Off = e.offset & 0x7FFFFFFF
			break
		}
	}
	if l1Off == 0 {
		return ""
	}

	var bestW uint32
	var bestData []byte

	l2Entries := r.readDir(l1Off)
	for _, l2 := range l2Entries {
		// L2: ID or named entry → L3: language subdirectory
		if l2.offset&0x80000000 == 0 {
			continue
		}
		l3Off := l2.offset & 0x7FFFFFFF
		l3Entries := r.readDir(l3Off)
		for _, l3 := range l3Entries {
			// L3: language entry → L4: data entry
			if l3.offset&0x80000000 == 0 {
				// Direct data entry
				dRVA, dSz := r.readData(l3.offset)
				if dSz < 40 || dSz > 500000 {
					continue
				}
				bmpData := r.readSlice(dRVA, dSz)
				if bmpData != nil && len(bmpData) >= 40 &&
					binary.LittleEndian.Uint32(bmpData[:4]) >= 40 &&
					binary.LittleEndian.Uint16(bmpData[14:16]) == 32 {
					w := binary.LittleEndian.Uint32(bmpData[4:8])
					if w > bestW {
						bestW = w
						bestData = bmpData
					}
				}
			} else {
				// Points to another subdirectory (L4)
				l4Off := l3.offset & 0x7FFFFFFF
				l4Entries := r.readDir(l4Off)
				for _, l4 := range l4Entries {
					if l4.offset&0x80000000 != 0 {
						continue
					}
					dRVA, dSz := r.readData(l4.offset)
					if dSz < 40 || dSz > 500000 {
						continue
					}
					bmpData := r.readSlice(dRVA, dSz)
					if bmpData != nil && len(bmpData) >= 40 &&
						binary.LittleEndian.Uint32(bmpData[:4]) >= 40 &&
						binary.LittleEndian.Uint16(bmpData[14:16]) == 32 {
						w := binary.LittleEndian.Uint32(bmpData[4:8])
						if w > bestW {
							bestW = w
							bestData = bmpData
						}
					}
				}
			}
		}
	}

	if bestData != nil {
		return encodeIconData(bestData)
	}
	return ""
}

type peReader struct {
	raw      []byte
	sections []*pe.Section
	baseRVA  uint32
}

func (r *peReader) rvaToOffset(rva uint32) int {
	for _, s := range r.sections {
		if rva >= s.VirtualAddress && rva < s.VirtualAddress+s.VirtualSize {
			return int(rva-s.VirtualAddress) + int(s.Offset)
		}
		if rva >= s.VirtualAddress && rva < s.VirtualAddress+s.Size {
			return int(rva-s.VirtualAddress) + int(s.Offset)
		}
	}
	return -1
}

func (r *peReader) readDir(off uint32) []resEntry {
	if off+16 > uint32(len(r.raw)) {
		return nil
	}
	named := binary.LittleEndian.Uint16(r.raw[off+12 : off+14])
	n := binary.LittleEndian.Uint16(r.raw[off+14 : off+16])
	total := int(named) + int(n)
	entries := make([]resEntry, total)
	for i := 0; i < total; i++ {
		e := off + 16 + uint32(i)*8
		entries[i] = resEntry{
			nameOrID: binary.LittleEndian.Uint32(r.raw[e:e+4]),
			offset:   binary.LittleEndian.Uint32(r.raw[e+4:e+8]),
		}
	}
	return entries
}

func (r *peReader) readData(off uint32) (uint32, uint32) {
	if off+8 > uint32(len(r.raw)) {
		return 0, 0
	}
	rva := binary.LittleEndian.Uint32(r.raw[off : off+4])
	sz := binary.LittleEndian.Uint32(r.raw[off+4 : off+8])
	return rva, sz
}

func (r *peReader) readSlice(rva uint32, sz uint32) []byte {
	off := r.rvaToOffset(rva)
	if off < 0 || off+int(sz) > len(r.raw) {
		return nil
	}
	return r.raw[off : off+int(sz)]
}

func parseIconSpec(spec string) (string, int) {
	spec = strings.TrimSpace(spec)
	lastComma := strings.LastIndex(spec, ",")
	if lastComma < 0 {
		return spec, 0
	}
	path := strings.TrimSpace(spec[:lastComma])
	idx, _ := strconv.Atoi(strings.TrimSpace(spec[lastComma+1:]))
	return path, idx
}

func ObtainIcon(path string, idx, size int) (uintptr, bool) {
	pathPtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, false
	}
	if h := privateExtractBest(pathPtr, idx, size); h != 0 {
		return h, true
	}
	var largeIcon, smallIcon uintptr
	ret, _, _ := procExtractIconExW.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(uint32(idx)),
		uintptr(unsafe.Pointer(&largeIcon)),
		uintptr(unsafe.Pointer(&smallIcon)),
		1,
	)
	if ret == 0 {
		if largeIcon != 0 { procDestroyIcon.Call(largeIcon) }
		if smallIcon != 0 { procDestroyIcon.Call(smallIcon) }
		return 0, false
	}
	icon := largeIcon
	if icon == 0 { icon = smallIcon }
	if icon == 0 { return 0, false }
	if largeIcon != 0 && largeIcon != icon { procDestroyIcon.Call(largeIcon) }
	if smallIcon != 0 && smallIcon != icon { procDestroyIcon.Call(smallIcon) }
	return icon, true
}

func privateExtractBest(pathPtr *uint16, idx, size int) uintptr {
	if err := procPrivateExtractIconsW.Find(); err != nil {
		return 0
	}
	var hIcon uintptr
	var iconID uint32
	ret, _, _ := procPrivateExtractIconsW.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(idx),
		uintptr(size), uintptr(size),
		uintptr(unsafe.Pointer(&hIcon)),
		uintptr(unsafe.Pointer(&iconID)),
		1, 0,
	)
	if ret == 0 || uint32(ret) == 0xffffffff || hIcon == 0 { return 0 }
	return hIcon
}

func RasterizeIconToPNG(icon uintptr, size int) string {
	type bi struct {
		sz uint32; w, h int32; p uint16; bc uint16; c uint32
		si uint32; xp, yp int32; cu, ci uint32
	}
	b := bi{sz: uint32(unsafe.Sizeof(bi{})), w: int32(size), h: -int32(size), p: 1, bc: 32, c: 0}
	dc, _, _ := procGetDC.Call(0)
	if dc == 0 { return "" }
	defer procReleaseDC.Call(0, dc)
	mdc, _, _ := procCreateCompatibleDC.Call(dc)
	if mdc == 0 { return "" }
	defer procDeleteDC.Call(mdc)
	var pp unsafe.Pointer
	bm, _, _ := procCreateDIBSection.Call(mdc, uintptr(unsafe.Pointer(&b)), 0, uintptr(unsafe.Pointer(&pp)), 0, 0)
	if bm == 0 || pp == nil { return "" }
	defer procDeleteObject.Call(bm)
	procSelectObject.Call(mdc, bm)
	ps := unsafe.Slice((*uint32)(pp), size*size)
	for i := range ps { ps[i] = 0 }
	procDrawIconEx.Call(mdc, 0, 0, icon, uintptr(size), uintptr(size), 0, 0, 3)
	ha := false
	for _, p := range ps { if uint8(p>>24) != 0 { ha = true; break } }
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for i, p := range ps { img.SetRGBA(i%size, i/size, color.RGBA{R: uint8(p>>16), G: uint8(p>>8), B: uint8(p), A: func() uint8 { if ha { return uint8(p>>24) }; return 255 }() }) }
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil { return "" }
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}

func ExtractIconAsBase64PNG(iconSpec string) string {
	path, idx := parseIconSpec(iconSpec)
	if path == "" { return "" }
	// 1. Standalone .ico
	if r := extractICOFile(path); r != "" { return r }
	// 2. PE resources
	if r := extractResourceIcons(path); r != "" { return r }
	// 3. GDI fallback
	hIcon, ok := ObtainIcon(path, idx, 256)
	if !ok || hIcon == 0 { return "" }
	defer procDestroyIcon.Call(hIcon)
	return RasterizeIconToPNG(hIcon, 256)
}

func ExtractIconAsBase64PNGSize(iconSpec string, _ int) string {
	return ExtractIconAsBase64PNG(iconSpec)
}