# Architecture

Switchy is a lightweight Windows utility that intercepts clicked links and routes them to your selected browser using Go's Wails framework and a Vue 3 frontend.

## Quick Start

```
wails dev       # Development server with hot reload
wails build     # Production binary
go fmt ./...    # Format code
go vet ./...    # Run linter
```

## Layer Map

### Entry Point
Process bootstrap and command-line argument routing. Detects flags (--register, --unregister, --settings) and dispatches to appropriate handler or browser selector UI.

Entry points:
- `main.go` — program entry, argument parsing, mode routing

### Business Logic & Rules
Core application logic including URL-to-browser routing rules, browser selection decisions, and configuration access. Implements rule matching (domain, regex, string containment), browser launching logic, and settings management interface.

Key files:
- `app.go` — Wails backend controller with 27 public methods
- `internal/rules/matcher.go` — rule evaluation engine
- `internal/config/service.go` — config I/O and caching

### Windows Integration
Windows-specific operations: registry manipulation for protocol handler registration, window management, DWM backdrop effects, IPC via named pipes. Handles browser detection via registry scanning.

Key files:
- `internal/registry/registry.go` — registry read/write
- `internal/browser/launch.go` — browser detection, profile scanning
- `internal/ipc/pipe.go` — single-instance enforcement

### Data Access & Configuration
Configuration persistence (JSON on disk), browser inventory management, user settings. Provides defaults, caching, and structured loading/saving.

Key files:
- `internal/config/service.go` — Load/Save interface
- `internal/config/types.go` — data structures

### Process Management
Browser launching with argument building, process spawning via os/exec, handling UWP app launches via explorer.exe, profile-specific launch configurations.

Key files:
- `internal/browser/launch.go` — OpenURL, DetectInstalled
- `internal/browser/icon_windows.go` — icon extraction

### UI Framework & Wails Bridge
Wails framework integration binding Go backend to Vue 3 frontend. Exposes backend methods via Wails runtime for frontend consumption.

Key files:
- `app.go` — App struct with exported methods
- `main.go` — Wails options setup

### Frontend (Vue 3)
Vue 3 single-page application providing browser selector UI and settings configuration. Communicates with Go backend via Wails IPC. Uses Tailwind CSS for styling.

Key files:
- `frontend/src/App.vue` — root component
- `frontend/src/views/SelectorView.vue` — browser picker UI
- `frontend/src/views/SettingsView.vue` — settings/rules editor
- `frontend/src/components/BrowserList.vue` — list rendering
- `frontend/src/components/BrowserButton.vue` — individual browser entry

### Utilities & Cross-Cutting Concerns
Logging service, browser utility functions, common helpers used across layers.

Key files:
- `internal/logger/logger.go` — structured logging (zap)
- `frontend/src/utils/browser.ts` — frontend browser helpers

## Guided Tour

### Understanding URL Interception
When a user clicks a link, Windows invokes Switchy via the registered protocol handler. The URL flows through argument parsing in `main.go`, IPC forwarding (if another instance exists), rule evaluation, and either automatic launch or UI display.

Start at `main.go:firstURL` to see how the URL arrives, then `main.go:autoRule` for rule matching.

### Understanding Browser Selection
The browser selector UI displays available browsers and responds to user selection. The backend (app.go) maintains the browser inventory, loaded from Windows registry on startup or manual refresh.

Start at `frontend/src/views/SelectorView.vue` for the UI, then `app.go:GetBrowsers` for backend browser list, then `internal/browser/launch.go:DetectInstalled` for registry scanning.

### Understanding Configuration Persistence
Settings, rules, and browser visibility are persisted to a JSON file in the user's AppData directory. The service layer (internal/config/service.go) handles loading, caching, and saving with validation.

Start at `internal/config/service.go:Save` for the persistence point, then trace backward to caller in app.go.

## Cross-Layer Dependencies

| From | To | Count | Pattern |
|------|-----|-------|---------|
| Business Logic → Windows Integration | 4 | Backend needs registry, browser launching, IPC |
| Business Logic → Data Access | 3 | Config reads for rules and settings |
| Windows Integration → Data Access | 2 | Browser detection saves to config |
| Frontend → Business Logic | 5 | Vue calls exported backend methods |
| Entry Point → Business Logic | 2 | Mode dispatch to core logic |
