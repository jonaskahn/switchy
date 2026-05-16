# Patterns

Recurring code patterns and hotspots across Switchy.

## Complexity Hotspots

| File | Complexity | Purpose |
|------|-----------|---------|
| app.go | complex | Core Wails controller with 27 public methods for browser launching, settings access, window management |
| main.go | moderate | CLI dispatcher routing arguments to selector, settings, register, or unregister modes |

## Function Exemplars by Layer

### Business Logic & Rules

- `app.go:OpenWithBrowser` — exemplar for browser launch coordination (app controller pattern)
- `internal/rules/matcher.go:Match` — exemplar for rule evaluation (pattern matching engine)
- `internal/config/service.go:Save` — exemplar for persistence with error handling

### Windows Integration

- `internal/registry/registry.go:writeClient` — exemplar for registry manipulation
- `internal/browser/launch.go:DetectInstalled` — exemplar for Windows API enumeration
- `internal/ipc/pipe.go:TryForward` — exemplar for named pipe IPC with timeout

### Frontend

- `frontend/src/views/SelectorView.vue` — exemplar for browser list UI interaction
- `frontend/src/composables/useTheme.ts` — exemplar for reactive state composition
- `frontend/src/components/BrowserButton.vue` — exemplar for event emission and profile selection

## Recurring Imports

| File | Import Count | Role |
|------|--------------|------|
| app.go | 8+ | Core backend controller, imported by main for UI mode dispatch |
| internal/config/service.go | 6+ | Config persistence, imported across layers for settings access |
| internal/browser/launch.go | 5+ | Browser detection and launching, imported by app and rules |
| internal/logger/logger.go | 4+ | Structured logging, cross-cutting import across backend |

## Naming Patterns

### Go Backend
- **Method prefix:** Public methods on `App` struct start uppercase (`GetBrowsers`, `OpenWithBrowser`, `RegisterAsDefault`)
- **Unexported helpers:** Internal functions lowercase (`autoRule`, `firstURL`, `readBrowserEntry`)
- **Error handling:** Errors returned as second return value (Go idiom), no exceptions thrown

### Vue Frontend
- **Component files:** PascalCase (SelectorView.vue, BrowserButton.vue, BrowserList.vue)
- **Composable files:** camelCase with use prefix (useTheme.ts, useBrowserIcon.ts)
- **Utility files:** camelCase (browser.ts)

## Architecture Patterns

### Wails Bridge Pattern
Backend (`app.go`) exposes public methods that Wails automatically binds to frontend. Frontend calls backend via `window.runtime` or Wails bindings. Single-threaded sequential communication, no concurrent requests.

Pattern location: `app.go` (27 exported methods), `frontend/src/App.vue` (Wails setup)

### Registry Abstraction
Windows registry operations wrapped in `internal/registry/registry.go` to isolate platform-specific code. All registry reads/writes centralized for maintainability.

Pattern location: `internal/registry/registry.go`, called from `app.go` for registration operations

### Rule Matching Engine
Pattern matching (domain, regex, substring) implemented in `internal/rules/matcher.go`, decoupled from UI and browser launching. Supports extensible rule types.

Pattern location: `internal/rules/matcher.go`, called by `main.go:autoRule`, `app.go:OpenWithBrowser`

### Browser Profile Detection
For Chromium browsers (Chrome, Edge, Brave), profiles detected by scanning filesystem (User Data\Default, etc). Each profile creates a separate browser entry in UI.

Pattern location: `internal/browser/launch.go:detectChromiumProfiles`

### Structured Logging
All logging uses `go.uber.org/zap` structured logger with leveled output (debug, info, warn, error). File-based with rotation via lumberjack.

Pattern location: `internal/logger/logger.go`
