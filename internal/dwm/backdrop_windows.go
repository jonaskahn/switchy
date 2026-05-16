//go:build windows

package dwm

import (
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	dwmwaSystemBackdropType = uintptr(38)
	dwmsbtMainWindow        = uint32(2)
	dwmsbtTransientWindow   = uint32(3)
)

var (
	dwmapi           = windows.NewLazySystemDLL("dwmapi.dll")
	dwmSetWindowAttr = dwmapi.NewProc("DwmSetWindowAttribute")

	currentBackdropType uint32 = dwmsbtTransientWindow
	enumCb              uintptr
)

func init() {
	pid := uint32(os.Getpid())
	enumCb = syscall.NewCallback(func(hwnd, _ uintptr) uintptr {
		var winPID uint32
		windows.GetWindowThreadProcessId(windows.HWND(hwnd), &winPID)
		if winPID == pid {
			dwmSetWindowAttr.Call(hwnd, dwmwaSystemBackdropType,
				uintptr(unsafe.Pointer(&currentBackdropType)), 4)
		}
		return 1
	})
}

// SetAppearance applies the window backdrop style. useMica=true enables Mica,
// false enables acrylic.
func SetAppearance(useMica bool) {
	if useMica {
		currentBackdropType = dwmsbtMainWindow
	} else {
		currentBackdropType = dwmsbtTransientWindow
	}
	windows.EnumWindows(enumCb, unsafe.Pointer(nil))
}
