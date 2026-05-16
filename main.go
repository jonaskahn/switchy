package main

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"switchy/internal/config"
	"switchy/internal/ipc"
	"switchy/internal/logger"
	"switchy/internal/rules"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

const (
	selectorFallbackWidth  = 580
	selectorFallbackHeight = 420
)

func main() {
	logger.Init()
	defer logger.Sync()

	logger.Info("Switchy starting")
	args := os.Args[1:]

	switch {
	case hasFlag(args, "--register"):
		logger.Info("Registering as default browser")
		handleRegister()
	case hasFlag(args, "--unregister"):
		logger.Info("Unregistering as default browser")
		handleUnregister()
	case hasFlag(args, "--settings"):
		logger.Info("Starting settings mode")
		runSettings()
	default:
		logger.Info("Starting selector mode", "url", firstURL(args))
		runSelector(firstURL(args))
	}
}

func handleRegister() {
	app := NewApp("settings", "")
	if err := app.RegisterAsDefault(); err != nil {
		fmt.Fprintln(os.Stderr, "register:", err)
		os.Exit(1)
	}
	fmt.Println("Switchy registered. Opening Windows Default Apps settings to complete setup.")
}

func handleUnregister() {
	app := NewApp("settings", "")
	if err := app.UnregisterAsDefault(); err != nil {
		fmt.Fprintln(os.Stderr, "unregister:", err)
		os.Exit(1)
	}
	fmt.Println("Switchy unregistered.")
}

func runSelector(rawURL string) {
	if rawURL != "" {
		if ipc.TryForward(rawURL) {
			return
		}
		if matched := autoRule(rawURL); matched != "" && matched != rules.ForcePickerBrowser {
			app := NewApp("selector", rawURL)
			_ = app.OpenWithBrowser(matched)
			return
		}
	}

	app := NewApp("selector", rawURL)

	go func() {
		_ = ipc.ListenAndServe(func(url string) {
			if matched := autoRule(url); matched != "" && matched != rules.ForcePickerBrowser {
				_ = app.OpenWithBrowser(matched)
				return
			}
			app.PushURL(url)
		})
	}()

	winW, winH := selectorFallbackWidth, selectorFallbackHeight
	appOptions := baseWailsOptions(app)
	appOptions.Title = "Switchy"
	appOptions.Width = winW
	appOptions.Height = winH
	appOptions.DisableResize = true
	appOptions.Frameless = true
	appOptions.AlwaysOnTop = true

	err := wails.Run(&appOptions)
	if err != nil {
		fmt.Fprintln(os.Stderr, "selector:", err)
	}
}

func runSettings() {
	app := NewApp("settings", "")
	appOptions := baseWailsOptions(app)
	appOptions.Title = "Switchy Settings"
	appOptions.Width = 960
	appOptions.Height = 680
	appOptions.MinWidth = 800
	appOptions.MinHeight = 600
	appOptions.Frameless = true

	err := wails.Run(&appOptions)
	if err != nil {
		fmt.Fprintln(os.Stderr, "settings:", err)
	}
}

func baseWailsOptions(app *App) options.App {
	return options.App{
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		AssetServer:      &assetserver.Options{Assets: assets},
		OnStartup:        app.startup,
		Bind:             []any{app},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			BackdropType:         resolvedBackdropType(),
		},
	}
}

func resolvedBackdropType() windows.BackdropType {
	s, err := config.Get()
	if err == nil && s.AppSettings.Appearance == "mica" {
		return windows.Mica
	}
	return windows.Acrylic
}

func autoRule(rawURL string) string {
	s, err := config.Get()
	if err != nil || !s.AppSettings.UseRules {
		return ""
	}
	return rules.Match(rawURL, s.Rulesets)
}

func hasFlag(args []string, flag string) bool {
	for _, a := range args {
		if strings.EqualFold(a, flag) {
			return true
		}
	}
	return false
}

func firstURL(args []string) string {
	for _, a := range args {
		if !strings.HasPrefix(a, "--") {
			return a
		}
	}
	return ""
}
