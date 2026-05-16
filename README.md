<p align="center">
  <img src="build/appicon.svg" width="130" alt="Switchy" />
</p>

<h1 align="center">Switchy</h1>

<p align="center">
  <strong>Every link. Your choice. Zero friction.</strong><br/>
  The lightning-fast browser picker for Windows that routes URLs exactly where <em>you</em> want them — every single time.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/platform-Windows%2010%2F11-0078D4?style=flat-square&logo=windows11&logoColor=white" alt="Platform" />
  <img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/Vue-3-4FC08D?style=flat-square&logo=vuedotjs&logoColor=white" alt="Vue" />
  <img src="https://img.shields.io/badge/Wails-v2-red?style=flat-square" alt="Wails" />
  <img src="https://img.shields.io/badge/license-MIT-22C55E?style=flat-square" alt="License" />
</p>

---

## Why Switchy?

You open Chrome for work, Firefox for personal browsing, Edge for that one internal app, and Brave for everything else.
But Windows lets you pick only **one** default browser — and every link goes there, no matter what.

**Switchy fixes that.**

Click any link → Switchy intercepts it → a sleek, acrylic-glass picker floats up in under a second → you tap a key and
it's open in exactly the right browser. You can even teach Switchy rules so work URLs go straight to Chrome, your team's
Notion links open in Firefox, and everything else still asks you. Pure control, zero compromise.

## Showcases

<p align="center">
  <img width="1200" alt="demo" src="https://github.com/user-attachments/assets/41200cfe-0dec-4701-82e4-3fe04bdd8bb8" />
</p>

## Features

- **Instant browser picker** — frameless, always-on-top acrylic window appears the moment a link is clicked; dismiss
  with `Esc` or pick with `1`–`9`
- **Two layout modes** — vertical list or horizontal dock; switch in settings
- **Timed auto-selection** — optional countdown that opens your default browser automatically if no key is pressed;
  timer and duration are configurable
- **Smart URL rules** — domain (`d$`), regex (`r$`), and substring (`s$`) patterns auto-open URLs without asking; force
  the picker with `_Switchy`
- **Alternate browser profiles** — launch a browser's incognito or work-profile variant directly from the picker
- **Single-instance IPC** — a second click forwards its URL to the running picker over a named pipe; no double windows
- **Native Windows feel** — Acrylic backdrop on the picker, Mica on settings, frameless with custom drag region; looks
  right at home on Windows 11
- **Auto-detect browsers** — scans `StartMenuInternet` registry to find every installed browser in seconds
- **Icon extraction** — rasterizes browser icons directly from executable resources at configurable sizes (16–256 px)
- **Settings persist** — all configuration saved to `%APPDATA%\Switchy\settings.json`

---

## Project Structure

```
switchy/
├── main.go               # Entry point & CLI modes (--register, --unregister, --settings)
├── app.go                # Wails-bound App struct — the frontend ↔ backend bridge
├── go.mod / go.sum
├── wails.json
├── .golangci.yml         # golangci-lint configuration
│
├── internal/             # Private Go packages
│   ├── browser/
│   │   ├── launch.go        # Launch browsers; DetectInstalled via registry; FilterLegacyIE
│   │   └── icon_windows.go  # Extract & rasterize HICON → base64 PNG (shell32/gdi32)
│   ├── config/
│   │   ├── service.go    # Load / Get / Save / Invalidate settings; schema migration
│   │   └── types.go      # UserSettings, Browser, AlternateLaunch, Ruleset, Rule, AppSettings
│   ├── ipc/
│   │   └── pipe.go       # Named-pipe single-instance IPC (go-winio, \\.\pipe\SwitchyPipe)
│   ├── registry/
│   │   └── registry.go   # Register / Unregister Switchy as Windows default browser
│   └── rules/
│       └── matcher.go    # URL-to-browser rule engine (domain / regex / substring)
│
├── frontend/             # Vue 3 + TypeScript + Tailwind CSS
│   ├── eslint.config.js  # ESLint 9 flat config (TS + Vue + Prettier-compat)
│   ├── .prettierrc.json  # Prettier config (no-semi, single-quote, Tailwind class sort)
│   ├── .prettierignore
│   └── src/
│       ├── App.vue           # Root router — renders SelectorView or SettingsView
│       ├── types.ts          # Frontend type definitions (mirrors Go structs)
│       ├── views/
│       │   ├── SelectorView.vue  # Browser picker (vertical list or horizontal dock)
│       │   └── SettingsView.vue  # Configuration hub (browsers, rulesets, app settings)
│       ├── components/
│       │   ├── BrowserList.vue       # Vertical stacked browser list
│       │   ├── BrowserButton.vue     # Single row item with icon, name, shortcut
│       │   ├── BrowserDock.vue       # Horizontal dock layout
│       │   ├── BrowserDockTile.vue   # Single dock tile with large icon and keycap badge
│       │   ├── UrlCard.vue           # URL display with protocol icon and copy button
│       │   ├── BrowserEditor.vue     # Browser config row in settings
│       │   └── RulesetEditor.vue     # Ruleset config row in settings
│       ├── composables/
│       │   └── useBrowserIcon.ts     # Async icon loader (calls GetBrowserIconAt backend)
│       └── utils/
│           ├── browser.ts    # Shortcut key helpers, animation delay calculations
│           └── config.ts     # Layout / overflow value normalizers
│
├── build/                # Wails build assets
│   ├── appicon.svg       # Source logo — regenerate appicon.png / icon.ico from this
│   ├── appicon.png       # 1024×1024 app icon (used by Wails build)
│   └── windows/
│       ├── icon.ico
│       ├── info.json         # Product metadata (version, company, copyright)
│       ├── wails.exe.manifest
│       └── installer/
│           ├── installer.iss   # Inno Setup 6 script → Switchy_Installer.exe
│           └── project.nsi     # NSIS script (wails build --nsis)
│
└── scripts/
    └── build.ps1         # Full release build: Wails → switchy.exe → Inno Setup installer
```

---

## Prerequisites

| Tool          | Version | Install                                                                 |
|---------------|---------|-------------------------------------------------------------------------|
| Go            | 1.25+   | [go.dev/dl](https://go.dev/dl/)                                         |
| Wails CLI     | v2      | `go install github.com/wailsapp/wails/v2/cmd/wails@latest`              |
| Node.js       | 18+     | [nodejs.org](https://nodejs.org/)                                       |
| golangci-lint | latest  | `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest` |
| Inno Setup    | 6       | [jrsoftware.org](https://jrsoftware.org/isinfo.php) *(installer only)*  |

---

## Development

```powershell
# From the switchy/ directory — hot-reload on both Go and Vue changes
wails dev
```

Go methods are also exposed as a dev server at `http://localhost:34115` so you can call them from the browser DevTools.

---

## Code Quality

### Frontend

```powershell
cd frontend

npm run lint          # report ESLint issues
npm run lint:fix      # auto-fix what's fixable
npm run format        # format all src files with Prettier
npm run format:check  # CI-safe Prettier check (no writes)
```

### Go

```powershell
# Format
gofmt -w .

# Vet
go vet ./...

# Full lint (requires golangci-lint)
golangci-lint run ./...
```

---

## Building

### Executable only

```powershell
wails build -platform windows/amd64
# → build/bin/switchy.exe
```

### Full release (exe + Inno Setup installer)

```powershell
.\scripts\build.ps1
# → build/bin/Switchy_Installer.exe
```

---

## Register as Default Browser

After building or installing, run once:

```powershell
.\build\bin\switchy.exe --register
```

Windows will open **Default apps** settings — scroll to *Web browser* and select **Switchy** to complete the setup. (
Windows 10/11 requires this manual confirmation step; no app can bypass it.)

To remove Switchy from the registry:

```powershell
.\build\bin\switchy.exe --unregister
```

---

## URL Rule Patterns

Rules are evaluated in order across enabled rulesets. The first match wins.

| Prefix   | Type      | Example             | Matches                           |
|----------|-----------|---------------------|-----------------------------------|
| `d$`     | Domain    | `d$example.com`     | `example.com` and `*.example.com` |
| `r$`     | Regex     | `r$^https://work\.` | Full URL matched against regex    |
| `s$`     | Substring | `s$/confluence/`    | URL contains the string           |
| *(none)* | Substring | `github.com`        | Same as `s$`                      |

Set a rule's browser target to `_Switchy` to force the picker even when other rules would match.

---

## CLI Flags

| Flag                           | Effect                                                        |
|--------------------------------|---------------------------------------------------------------|
| *(no flag, with URL argument)* | Open the browser picker for the given URL                     |
| `--register`                   | Write Windows registry entries and open Default Apps settings |
| `--unregister`                 | Remove Switchy registry entries                               |
| `--settings`                   | Open the settings window                                      |

---

## Inspiration

> Switchy is lovingly inspired by **[Hurl](https://github.com/U-C-S/Hurl)** — a fantastic open-source browser picker for
> Windows that proved the concept beautifully.
> Switchy builds on that idea with a Wails-powered UI, native-pipe IPC, URL rule matching, and alternate browser
> profiles.
> Huge respect to the Hurl team for laying the foundation.

---

## Agent Context

This repository uses [**agent-context**](https://github.com/jonaskahn/agent-context) to generate and maintain its
AI-agent context files (e.g., `AGENTS.md`). If you're working on Switchy with an AI assistant, those files provide the
grounding architecture, conventions, and module map needed to navigate the codebase effectively.

---

## License

[MIT](LICENSE)
