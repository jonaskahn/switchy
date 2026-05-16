//go:build !windows

package dwm

// SetAppearance applies the window backdrop style. useMica=true enables Mica,
// false enables acrylic (no-op on non-Windows).
func SetAppearance(useMica bool) {}
