# Conventions

<!--
Human-authored coding standards live here. /agent-context reads this file on the
next run and distills it into docs/agents/conventions.md. Keep directives terse
and imperative ("MUST NOT commit X", "Always use Y"). Delete sections you do not
need; add more as conventions emerge.
-->

## Safety

- MUST NOT commit secrets, `.env` files, or credentials.

Bad:

```go
const token = "sk_live_secret"
```

Good:

```go
token := os.Getenv("API_TOKEN")
```

- MUST NOT disable, skip, or bypass tests to make code pass.

Bad:

```go
func TestSave(t *testing.T) {
	t.Skip("flaky")
}
```

Good:

```go
func TestSave(t *testing.T) {
	if err := save(settings); err != nil {
		t.Fatal(err)
	}
}
```

- MUST run `gofmt` on Go code before completion.

Bad:

```go
func save( path string)error{return nil}
```

Good:

```go
func save(path string) error {
	return nil
}
```

## Naming

- Name functions by abstraction level: use broader names for high-level functions and more specific names for lower-level helpers.

Bad:

```go
func DetermineFileExtensionAndParseConfigurationFile(path string) (Config, error) {
	return parseJSONConfigurationFile(path)
}
```

Good:

```go
func Parse(path string) (Config, error) {
	return parseJSON(path)
}
```

- Name variables by scope: use descriptive names for wider scopes and short names only in narrow local scopes.

Bad:

```go
func save(browsers []Browser) {
	for _, browserWithConfiguredLaunchArguments := range browsers {
		open(browserWithConfiguredLaunchArguments)
	}
}
```

Good:

```go
func save(browsers []Browser) {
	for _, browser := range browsers {
		open(browser)
	}
}
```

- MUST NOT name variables after their types when the domain meaning is clearer.

Bad:

```go
var stringValue string
```

Good:

```go
var browserName string
```

## Patterns

- Follow KISS: prefer the simplest design that solves the current task.

Bad:

```go
type LauncherFactory struct{}

func (LauncherFactory) Build() Launcher {
	return Launcher{}
}
```

Good:

```go
launcher := Launcher{}
```

- Follow DRY: remove meaningful duplication, but MUST NOT create premature abstractions.

Bad:

```go
chromePath := strings.TrimSpace(chrome.Path)
edgePath := strings.TrimSpace(edge.Path)
```

Good:

```go
func cleanPath(path string) string {
	return strings.TrimSpace(path)
}
```

- Keep functions small and focused on one responsibility.

Bad:

```go
func SaveAndLaunch(settings Settings, url string) error {
	if err := save(settings); err != nil {
		return err
	}
	return launch(settings.DefaultBrowser, url)
}
```

Good:

```go
func Save(settings Settings) error {
	return save(settings)
}

func Launch(browser Browser, url string) error {
	return launch(browser, url)
}
```

- Prefer early returns over deep nesting.

Bad:

```go
func open(browser Browser, url string) error {
	if browser.Enabled {
		if url != "" {
			return launch(browser, url)
		}
	}
	return ErrInvalidLaunch
}
```

Good:

```go
func open(browser Browser, url string) error {
	if !browser.Enabled {
		return ErrInvalidLaunch
	}
	if url == "" {
		return ErrInvalidLaunch
	}
	return launch(browser, url)
}
```

- Declare variables close to their first use.

Bad:

```go
func open(browser Browser) error {
	path := browser.Path
	validate(browser)
	return exec.Command(path).Start()
}
```

Good:

```go
func open(browser Browser) error {
	validate(browser)
	path := browser.Path
	return exec.Command(path).Start()
}
```

- Keep mutable scope small; avoid global mutable state.

Bad:

```go
var currentSettings Settings

func SetTheme(theme string) {
	currentSettings.Theme = theme
}
```

Good:

```go
func SetTheme(settings Settings, theme string) Settings {
	settings.Theme = theme
	return settings
}
```

- Keep function signatures small; use option structs when parameter lists become unclear.

Bad:

```go
func Register(name string, path string, args string, hidden bool, uwp bool) Browser
```

Good:

```go
type BrowserOptions struct {
	Name   string
	Path   string
	Args   string
	Hidden bool
	UWP    bool
}

func Register(opts BrowserOptions) Browser
```

- Prefer concrete return types; accept small interfaces where useful.

Bad:

```go
func NewBrowser() interface{} {
	return Browser{}
}
```

Good:

```go
func NewBrowser() Browser {
	return Browser{}
}
```

- Avoid exposing `interface{}` or `any` unless required at system boundaries.

Bad:

```go
func Save(value any) error {
	browser := value.(Browser)
	return save(browser)
}
```

Good:

```go
func Save(browser Browser) error {
	return save(browser)
}
```

- Use pointers intentionally and avoid expanding mutable ownership unnecessarily.

Bad:

```go
func (store *Store) Get(id string) *Browser {
	return store.browsers[id]
}
```

Good:

```go
func (store *Store) Get(id string) Browser {
	return *store.browsers[id]
}
```

## Workflow

- MUST NOT add explanatory comments inside implementation code.

Bad:

```go
// Loop through browsers and launch the selected one.
for _, browser := range browsers {
	launch(browser)
}
```

Good:

```go
for _, browser := range browsers {
	launch(browser)
}
```

- Code should explain itself through names, structure, and small functions.

Bad:

```go
func p(b Browser) bool {
	return b.Path != ""
}
```

Good:

```go
func hasLaunchPath(browser Browser) bool {
	return browser.Path != ""
}
```

- Only add comments for exported public APIs, interfaces, or genuinely non-obvious why-level constraints.

Bad:

```go
// Browser stores browser data.
type Browser struct {
	Name string
}
```

Good:

```go
// Launcher opens URLs with a configured browser.
type Launcher interface {
	Open(url string) error
}
```
