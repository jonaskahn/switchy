package config

const (
	layoutVertical    = "vertical"
	appearanceAcrylic = "acrylic"
	appearanceMica    = "mica"
	themeDefault      = "default"
)

// AlternateLaunch configures browser-specific launch overrides.
type AlternateLaunch struct {
	Name       string `json:"Name"`
	LaunchArgs string `json:"LaunchArgs"`
}

// Browser represents a detected browser with its launch configuration and
// display metadata.
type Browser struct {
	Name              string            `json:"Name"`
	ExePath           string            `json:"ExePath"`
	IconSpec          string            `json:"IconSpec"`
	LaunchArgs        string            `json:"LaunchArgs"`
	IsUwp             bool              `json:"IsUwp"`
	Hidden            bool              `json:"Hidden"`
	AlternateLaunches []AlternateLaunch `json:"AlternateLaunches"`
	ProfileIcon       string            `json:"ProfileIcon"`
}

// Rule represents a single URL routing rule.
type Rule struct {
	Pattern string `json:"Pattern"`
	Browser string `json:"Browser"`
}

// Ruleset is a named collection of rules.
type Ruleset struct {
	Name    string `json:"Name"`
	Enabled bool   `json:"Enabled"`
	Rules   []Rule `json:"Rules"`
}

// AppSettings holds application-level UI and behaviour settings.
type AppSettings struct {
	UseRules       bool   `json:"UseRules"`
	TimedSelection bool   `json:"TimedSelection"`
	TimedSeconds   int    `json:"TimedSeconds"`
	DefaultBrowser string `json:"DefaultBrowser"`
	Appearance     string `json:"Appearance"`
	Theme          string `json:"Theme"`
	SelectorLayout string `json:"SelectorLayout"`
}

// UserSettings is the top-level settings structure persisted to disk.
type UserSettings struct {
	Browsers    []Browser   `json:"Browsers"`
	Rulesets    []Ruleset   `json:"Rulesets"`
	AppSettings AppSettings `json:"AppSettings"`
}

// DefaultSettings returns a UserSettings populated with sensible defaults.
func DefaultSettings() *UserSettings {
	return &UserSettings{
		Browsers: []Browser{},
		Rulesets: []Ruleset{},
		AppSettings: AppSettings{
			TimedSeconds:   5,
			SelectorLayout: layoutVertical,
			Appearance:     appearanceAcrylic,
			Theme:          themeDefault,
		},
	}
}
