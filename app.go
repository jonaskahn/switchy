package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"sync"
	"unsafe"

	"switchy/internal/browser"
	"switchy/internal/config"
	"switchy/internal/dwm"
	"switchy/internal/logger"
	"switchy/internal/registry"
	"switchy/internal/rules"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/sys/windows"
)

type App struct {
	ctx    context.Context
	mu     sync.Mutex
	mode   string
	rawURL string
}

func NewApp(mode, rawURL string) *App {
	return &App{mode: mode, rawURL: rawURL}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	logger.Debug("App startup in mode:", a.mode)
}

func (a *App) currentURL() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.rawURL
}

func (a *App) setCurrentURL(url string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.rawURL = url
}

func (a *App) GetMode() string { return a.mode }

func (a *App) GetSettings() (*config.UserSettings, error) {
	s, err := config.Load()
	if err != nil {
		logger.Error("Failed to load settings:", err)
		return nil, err
	}
	logger.Debug("Settings loaded successfully")
	out := *s
	out.Browsers = browser.FilterLegacyInternetExplorer(append([]config.Browser(nil), s.Browsers...))
	return &out, nil
}

func (a *App) GetURL() string {
	return a.currentURL()
}

func (a *App) SetURL(url string) {
	a.setCurrentURL(url)
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "url:new", url)
	}
}

func (a *App) PushURL(url string) {
	a.setCurrentURL(url)
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "url:new", url)
		runtime.WindowShow(a.ctx)
	}
}

func (a *App) GetBrowsers() ([]config.Browser, error) {
	s, err := config.Get()
	if err != nil {
		return nil, err
	}
	filtered := browser.FilterLegacyInternetExplorer(s.Browsers)
	visible := make([]config.Browser, 0, len(filtered))
	for _, b := range filtered {
		if !b.Hidden {
			visible = append(visible, b)
		}
	}
	return visible, nil
}

func (a *App) OpenWithBrowser(browserName string) error {
	url := a.currentURL()
	logger.Info("Opening URL with browser:", "browser", browserName, "url", url)
	s, err := config.Get()
	if err != nil {
		logger.Error("Failed to get config for browser open:", err)
		return err
	}
	b, ok := findBrowserByName(s.Browsers, browserName)
	if !ok {
		logger.Warn("Browser not found:", "name", browserName)
		return nil
	}
	if err := browser.OpenURL(b, url); err != nil {
		logger.Error("Failed to open URL with browser:", "browser", browserName, "error", err)
		return err
	}
	logger.Info("Browser opened successfully", "browser", browserName)
	runtime.Quit(a.ctx)
	return nil
}

func (a *App) OpenWithAlternate(browserName, altName string) error {
	url := a.currentURL()
	s, err := config.Get()
	if err != nil {
		return err
	}
	b, ok := findBrowserByName(s.Browsers, browserName)
	if !ok {
		return nil
	}
	if err := browser.OpenURLWithAlternate(b, altName, url); err != nil {
		return err
	}
	runtime.Quit(a.ctx)
	return nil
}

func (a *App) CopyURL() error {
	return runtime.ClipboardSetText(a.ctx, a.currentURL())
}

func (a *App) DismissWindow() {
	runtime.Quit(a.ctx)
}

func (a *App) ResizeSelector(width, height int) {
	if width < 200 {
		width = 200
	}
	if height < 160 {
		height = 160
	}
	if a.ctx == nil {
		return
	}
	runtime.WindowSetSize(a.ctx, width, height)

	// Center on the primary work area after every resize.
	var rect windows.Rect
	_, _, _ = windows.NewLazySystemDLL("user32.dll").NewProc("SystemParametersInfoW").Call(
		0x0030, 0, uintptr(unsafe.Pointer(&rect)), 0,
	)
	workW := int(rect.Right - rect.Left)
	workH := int(rect.Bottom - rect.Top)
	workX := int(rect.Left)
	workY := int(rect.Top)
	x := workX + (workW-width)/2
	y := workY + (workH-height)/2
	if x < workX {
		x = workX
	}
	if y < workY {
		y = workY
	}
	runtime.WindowSetPosition(a.ctx, x, y)
}

type WindowPos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type WindowSize struct {
	W int `json:"w"`
	H int `json:"h"`
}

type WorkArea struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

func (a *App) GetWorkArea() WorkArea {
	var rect windows.Rect
	_, _, _ = windows.NewLazySystemDLL("user32.dll").NewProc("SystemParametersInfoW").Call(
		0x0030,
		0,
		uintptr(unsafe.Pointer(&rect)),
		0,
	)
	return WorkArea{
		X: int(rect.Left),
		Y: int(rect.Top),
		W: int(rect.Right - rect.Left),
		H: int(rect.Bottom - rect.Top),
	}
}

func (a *App) GetWindowPosition() WindowPos {
	if a.ctx == nil {
		return WindowPos{}
	}
	x, y := runtime.WindowGetPosition(a.ctx)
	return WindowPos{X: x, Y: y}
}

func (a *App) GetWindowSize() WindowSize {
	if a.ctx == nil {
		return WindowSize{}
	}
	w, h := runtime.WindowGetSize(a.ctx)
	return WindowSize{W: w, H: h}
}

func (a *App) OpenSettings() {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	_ = exec.Command(exe, "--settings").Start()
}

func (a *App) ApplyAutoRules() (string, error) {
	s, err := config.Get()
	if err != nil {
		return "", err
	}
	return rules.Match(a.currentURL(), s.Rulesets), nil
}

func (a *App) RememberBrowserForDomain(browserName string) (string, error) {
	parsed, err := url.Parse(a.currentURL())
	if err != nil || parsed.Hostname() == "" {
		return "", fmt.Errorf("no hostname in URL %q", a.currentURL())
	}
	hostname := parsed.Hostname()

	s, err := config.Load()
	if err != nil {
		return "", err
	}

	const quickRules = "Quick Rules"
	idx := -1
	for i, rs := range s.Rulesets {
		if rs.Name == quickRules {
			idx = i
			break
		}
	}
	if idx == -1 {
		s.Rulesets = append(s.Rulesets, config.Ruleset{
			Name:    quickRules,
			Enabled: true,
			Rules:   []config.Rule{},
		})
		idx = len(s.Rulesets) - 1
	}

	pattern := "d$" + hostname
	updated := false
	for i, r := range s.Rulesets[idx].Rules {
		if r.Pattern == pattern {
			s.Rulesets[idx].Rules[i].Browser = browserName
			updated = true
			break
		}
	}
	if !updated {
		s.Rulesets[idx].Rules = append(s.Rulesets[idx].Rules, config.Rule{
			Pattern: pattern,
			Browser: browserName,
		})
	}

	s.AppSettings.UseRules = true

	if err := config.Save(s); err != nil {
		return "", err
	}
	return hostname, nil
}

func (a *App) SaveSettings(s *config.UserSettings) error {
	return config.Save(s)
}

func (a *App) DetectBrowsers() []config.Browser {
	return browser.DetectInstalled()
}

func (a *App) GetBrowserIcon(iconSpec string) string {
	return browser.ExtractIconAsBase64PNG(iconSpec)
}

func (a *App) GetBrowserIconAt(iconSpec string, sizePx int) string {
	return browser.ExtractIconAsBase64PNGSize(iconSpec, sizePx)
}

func (a *App) RegisterAsDefault() error {
	if err := registry.Register(""); err != nil {
		return err
	}
	return registry.OpenDefaultAppsSettings()
}

func (a *App) UnregisterAsDefault() error {
	return registry.Unregister()
}

func (a *App) ApplyAppearance(appearance string) {
	dwm.SetAppearance(appearance == "mica")
}

func (a *App) CloseSettings() {
	runtime.Quit(a.ctx)
}

func findBrowserByName(browsers []config.Browser, name string) (config.Browser, bool) {
	for _, b := range browsers {
		if b.Name == name {
			return b, true
		}
	}
	return config.Browser{}, false
}
