//go:build windows

package browser

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"switchy/internal/config"
	"switchy/internal/logger"

	"golang.org/x/sys/windows/registry"
)

const ownName = "Switchy"

// OpenURL opens a URL with the given browser.
func OpenURL(b config.Browser, rawURL string) error {
	args := buildArgs(b.LaunchArgs, rawURL)
	if b.IsUwp {
		return exec.Command("explorer.exe", args...).Start()
	}
	logger.Debug("Launching browser", "exe", b.ExePath, "args", args)
	return exec.Command(b.ExePath, args...).Start()
}

// OpenURLWithAlternate opens a URL using another browser's configuration.
func OpenURLWithAlternate(b config.Browser, altName string, rawURL string) error {
	for _, alt := range b.AlternateLaunches {
		if alt.Name == altName {
			logger.Debug("Launching alternate", "exe", b.ExePath, "alt", altName, "args", alt.LaunchArgs)
			return exec.Command(b.ExePath, buildArgs(alt.LaunchArgs, rawURL)...).Start()
		}
	}
	return fmt.Errorf("alternate launch %q not found on browser %q", altName, b.Name)
}

func buildArgs(launchArgs, rawURL string) []string {
	if launchArgs == "" {
		return []string{rawURL}
	}

	parsed := splitArgs(launchArgs)

	hasPlaceholder := false
	for i, arg := range parsed {
		if arg == "%URL%" {
			parsed[i] = rawURL
			hasPlaceholder = true
		}
	}

	if !hasPlaceholder {
		parsed = append(parsed, rawURL)
	}

	return parsed
}

func splitArgs(s string) []string {
	var args []string
	var current strings.Builder
	inQuote := false
	quoteChar := byte(0)

	for i := 0; i < len(s); i++ {
		c := s[i]

		if inQuote {
			if c == quoteChar {
				inQuote = false
				quoteChar = 0
			} else {
				current.WriteByte(c)
			}
		} else {
			switch c {
			case '"', '\'':
				inQuote = true
				quoteChar = c
			case ' ', '\t':
				if current.Len() > 0 {
					args = append(args, current.String())
					current.Reset()
				}
			default:
				current.WriteByte(c)
			}
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

func isLegacyInternetExplorer(displayName, exePath string) bool {
	cleanPath := filepath.Clean(strings.Trim(exePath, `"`))
	lowerPath := strings.ToLower(cleanPath)
	if strings.Contains(lowerPath, `\internet explorer`) {
		return true
	}
	if strings.EqualFold(filepath.Base(cleanPath), "iexplore.exe") {
		return true
	}
	return strings.TrimSpace(strings.ToLower(displayName)) == "internet explorer"
}

func isLegacyIERegistryKey(regKey string) bool {
	k := strings.TrimSpace(strings.ToLower(regKey))
	return k == "iexplore.exe" || k == "internet explorer"
}

// FilterLegacyInternetExplorer removes Internet Explorer entries that ship
// with Windows (IE mode of Edge and the legacy IE application).
func FilterLegacyInternetExplorer(entries []config.Browser) []config.Browser {
	if len(entries) == 0 {
		return entries
	}
	out := make([]config.Browser, 0, len(entries))
	for _, b := range entries {
		if !isLegacyInternetExplorer(b.Name, b.ExePath) {
			out = append(out, b)
		}
	}
	return out
}

type chromeProfile struct {
	Name string `json:"name"`
}

type chromeLocalState struct {
	Profile struct {
		InfoCache map[string]chromeProfile `json:"info_cache"`
	} `json:"profile"`
}

func isChromiumBrowser(b config.Browser) bool {
	name := strings.ToLower(b.Name)
	exe := strings.ToLower(b.ExePath)

	return strings.Contains(name, "chrome") ||
		strings.Contains(name, "chromium") ||
		strings.Contains(name, "edge") ||
		strings.Contains(exe, "chrome.exe") ||
		strings.Contains(exe, "chromium.exe") ||
		strings.Contains(exe, "msedge.exe") ||
		strings.Contains(exe, "brave.exe")
}

func getUserDataDir(b config.Browser) string {
	if testDir := os.Getenv("CHROME_TEST_USER_DATA_DIR"); testDir != "" {
		return testDir
	}

	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		userProfile := os.Getenv("USERPROFILE")
		if userProfile == "" {
			return ""
		}
		localAppData = filepath.Join(userProfile, "AppData", "Local")
	}

	exe := strings.ToLower(b.ExePath)
	if strings.Contains(exe, "msedge.exe") {
		return filepath.Join(localAppData, "Microsoft", "Edge", "User Data")
	}
	if strings.Contains(exe, "brave.exe") {
		return filepath.Join(localAppData, "BraveSoftware", "Brave-Browser", "User Data")
	}
	if strings.Contains(exe, "chrome.exe") || strings.Contains(exe, "chromium.exe") {
		return filepath.Join(localAppData, "Google", "Chrome", "User Data")
	}

	// Chromium-based fallback: derive from exe path pattern
	// e.g. ...\SomeBrowser\...\somebrowser.exe -> guess User Data
	return ""
}

func getProfileIcon(browserName, userDataDir, dirName string) string {
	profileDir := filepath.Join(userDataDir, dirName)

	var candidates []string
	if strings.Contains(strings.ToLower(browserName), "edge") {
		candidates = []string{
			profileDir + "/Edge Profile Picture.png",
		}
	} else {
		candidates = []string{
			profileDir + "/Google Profile Picture",
			profileDir + "/Google Profile Picture.png",
			profileDir + "/Google Profile Picture.jpg",
			profileDir + "/Google Profile Picture.jpeg",
		}
	}
	for _, candidate := range candidates {
		data, err := os.ReadFile(candidate)
		if err != nil {
			continue
		}
		ext := strings.ToLower(filepath.Ext(candidate))
		mimeType := getMimeType(ext)
		return "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(data)
	}
	return ""
}

func getMimeType(ext string) string {
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".ico":
		return "image/x-icon"
	default:
		return "image/png"
	}
}

func makeProfileBrowser(base config.Browser, displayName, dirName, userDataDir string) config.Browser {
	return config.Browser{
		Name:              fmt.Sprintf("%s (%s)", base.Name, displayName),
		ExePath:           base.ExePath,
		IconSpec:          base.IconSpec,
		LaunchArgs:        fmt.Sprintf(`--profile-directory="%s"`, dirName),
		IsUwp:             base.IsUwp,
		Hidden:            base.Hidden,
		AlternateLaunches: base.AlternateLaunches,
		ProfileIcon:       getProfileIcon(base.Name, userDataDir, dirName),
	}
}

func detectChromiumProfiles(baseBrowser config.Browser) []config.Browser {
	var profiles []config.Browser

	userDataDir := getUserDataDir(baseBrowser)
	if userDataDir == "" {
		logger.Warn("Chrome user data directory not found")
		return profiles
	}

	localStatePath := filepath.Join(userDataDir, "Local State")
	data, err := os.ReadFile(localStatePath)
	if err != nil {
		logger.Debug("Cannot read Chrome Local State, falling back to directory scan", "error", err)
		return detectProfilesByDirectory(userDataDir, baseBrowser)
	}

	var localState chromeLocalState
	if err := json.Unmarshal(data, &localState); err != nil {
		logger.Debug("Cannot parse Chrome Local State, falling back to directory scan", "error", err)
		return detectProfilesByDirectory(userDataDir, baseBrowser)
	}

	for dirName, profile := range localState.Profile.InfoCache {
		if dirName == "Guest Profile" || dirName == "System Profile" {
			continue
		}

		displayName := profile.Name
		if displayName == "" {
			displayName = dirName
		}

		if dirName == "Default" && displayName == dirName {
			displayName = "Personal"
		}

		profiles = append(profiles, makeProfileBrowser(baseBrowser, displayName, dirName, userDataDir))
		logger.Debug("Detected Chromium profile", "name", baseBrowser.Name, "dir", dirName)
	}

	return profiles
}

func detectProfilesByDirectory(userDataDir string, baseBrowser config.Browser) []config.Browser {
	var profiles []config.Browser

	dirs, err := os.ReadDir(userDataDir)
	if err != nil {
		logger.Debug("Cannot read Chrome user data directory", "error", err)
		return profiles
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		dirName := dir.Name()

		if dirName == "Guest Profile" || dirName == "System Profile" ||
			(!strings.HasPrefix(dirName, "Profile ") && dirName != "Default") {
			continue
		}

		displayName := dirName
		if dirName == "Default" {
			displayName = "Personal"
		} else if strings.HasPrefix(dirName, "Profile ") {
			displayName = strings.TrimPrefix(dirName, "Profile ")
			if displayName == "" {
				displayName = "Profile"
			}
		}

		profiles = append(profiles, makeProfileBrowser(baseBrowser, displayName, dirName, userDataDir))
		logger.Debug("Detected Chromium profile (directory scan)", "name", baseBrowser.Name, "dir", dirName)
	}

	return profiles
}

// DetectInstalled discovers browsers installed on the local system.
func DetectInstalled() []config.Browser {
	var found []config.Browser
	roots := []registry.Key{registry.LOCAL_MACHINE, registry.CURRENT_USER}
	subpath := `SOFTWARE\Clients\StartMenuInternet`

	for _, root := range roots {
		k, err := registry.OpenKey(root, subpath, registry.READ)
		if err != nil {
			continue
		}
		names, _ := k.ReadSubKeyNames(-1)
		k.Close()

		for _, name := range names {
			if strings.EqualFold(name, ownName) || isLegacyIERegistryKey(name) {
				continue
			}
			b, ok := readBrowserEntry(root, subpath, name)
			if !ok {
				continue
			}

			if isChromiumBrowser(b) {
				profiles := detectChromiumProfiles(b)
				if len(profiles) > 0 {
					found = append(found, profiles...)
				} else {
					found = append(found, b)
				}
			} else {
				found = append(found, b)
			}
		}
	}
	// Deduplicate browsers with the same name+exe but from different registry roots
	// (e.g. Firefox registered under both HKLM and HKCU).
	// Chrome profiles share the same ExePath but have distinct names like "Chrome (Personal)",
	// so they are preserved by the composite key.
	return FilterLegacyInternetExplorer(deduplicate(found))
}

func deduplicate(browsers []config.Browser) []config.Browser {
	type browserKey struct {
		name   string
		exeKey string // normalized to lowercase for dedup
	}
	seen := make(map[browserKey]bool, len(browsers))
	var result []config.Browser
	for _, b := range browsers {
		key := browserKey{name: b.Name, exeKey: strings.ToLower(b.ExePath)}
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, b)
	}
	return result
}

func readExePath(root registry.Key, subpath, name string) (string, bool) {
	k, err := registry.OpenKey(root, subpath+`\`+name+`\shell\open\command`, registry.READ)
	if err != nil {
		return "", false
	}
	defer k.Close()
	v, _, err := k.GetStringValue("")
	if err != nil {
		return "", false
	}
	return strings.Trim(v, `"`), true
}

func readDisplayName(root registry.Key, subpath, name, fallback string) string {
	k, err := registry.OpenKey(root, subpath+`\`+name, registry.READ)
	if err != nil {
		return fallback
	}
	defer k.Close()
	v, _, err := k.GetStringValue("")
	if err != nil || v == "" {
		return fallback
	}
	return v
}

func readIconSpec(root registry.Key, subpath, name string) string {
	k, err := registry.OpenKey(root, subpath+`\`+name+`\DefaultIcon`, registry.READ)
	if err != nil {
		return ""
	}
	defer k.Close()
	v, _, _ := k.GetStringValue("")
	return v
}

func readBrowserEntry(root registry.Key, subpath, name string) (config.Browser, bool) {
	exePath, ok := readExePath(root, subpath, name)
	if !ok {
		return config.Browser{}, false
	}
	return config.Browser{
		Name:              readDisplayName(root, subpath, name, name),
		ExePath:           exePath,
		IconSpec:          readIconSpec(root, subpath, name),
		AlternateLaunches: []config.AlternateLaunch{},
	}, true
}
