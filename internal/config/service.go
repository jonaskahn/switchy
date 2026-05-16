package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const settingsDir = "Switchy"
const settingsFile = "settings.json"

var (
	mu             sync.RWMutex
	cachedSettings *UserSettings
)

func settingsPath() (string, error) {
	roaming, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(roaming, settingsDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, settingsFile), nil
}

var cachedPathOnce sync.Once

var cachedPath struct {
	path string
}

func cachedSettingsPath() (string, error) {
	cachedPathOnce.Do(func() {
		roaming, err := os.UserConfigDir()
		if err != nil {
			// Leave path empty to signal error on next call
			return
		}
		dir := filepath.Join(roaming, settingsDir)
		if err := os.MkdirAll(dir, 0o755); err == nil {
			cachedPath.path = filepath.Join(dir, settingsFile)
		}
	})
	if cachedPath.path == "" {
		return settingsPath()
	}
	return cachedPath.path, nil
}

// Load reads settings from disk, applies migrations, and caches the result.
func Load() (*UserSettings, error) {
	mu.Lock()
	defer mu.Unlock()

	path, err := cachedSettingsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		defaults := DefaultSettings()
		if saveErr := saveUnlocked(defaults); saveErr != nil {
			return nil, saveErr
		}
		cachedSettings = defaults
		return cachedSettings, nil
	}
	if err != nil {
		return nil, err
	}

	var s UserSettings
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	migrateUserSettings(&s)
	cachedSettings = &s
	return cachedSettings, nil
}

// Get returns the cached settings, loading from disk if not yet cached.
func Get() (*UserSettings, error) {
	mu.RLock()
	if cachedSettings != nil {
		defer mu.RUnlock()
		return cachedSettings, nil
	}
	mu.RUnlock()
	return Load()
}

// Save writes settings to disk and updates the cache.
func Save(s *UserSettings) error {
	mu.Lock()
	defer mu.Unlock()
	return saveUnlocked(s)
}

func saveUnlocked(s *UserSettings) error {
	path, err := cachedSettingsPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return err
	}
	cachedSettings = s
	return nil
}

// Invalidate clears the cached settings, forcing a fresh load on next Get.
func Invalidate() {
	mu.Lock()
	cachedSettings = nil
	mu.Unlock()
}

func migrateUserSettings(s *UserSettings) {
	initEmptySlices(s)
	migrateLayoutDefaults(&s.AppSettings)
	migrateAppearance(&s.AppSettings)
}

func migrateAppearance(a *AppSettings) {
	if a.Appearance != "" {
		return
	}
	switch a.Theme {
	case appearanceMica:
		a.Appearance = appearanceMica
	default:
		a.Appearance = appearanceAcrylic
	}
	a.Theme = themeDefault
}

func initEmptySlices(s *UserSettings) {
	if s.Browsers == nil {
		s.Browsers = []Browser{}
	}
	if s.Rulesets == nil {
		s.Rulesets = []Ruleset{}
	}
}

func migrateLayoutDefaults(a *AppSettings) {
	if strings.TrimSpace(a.SelectorLayout) == "" {
		a.SelectorLayout = layoutVertical
	}
}
