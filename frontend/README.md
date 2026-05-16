# Switchy — Frontend

Vue 3 + TypeScript + Tailwind CSS frontend for the Switchy browser picker, bundled by Vite and embedded into the Wails binary at build time.

## Stack

| Tool | Version | Role |
|---|---|---|
| Vue 3 | 3.5+ | UI framework (`<script setup>` SFCs) |
| TypeScript | 5.7+ | Type checking |
| Vite | 8+ | Dev server & bundler |
| Tailwind CSS | 3 | Utility-first styling |
| Lucide Vue | 1+ | Icon set |
| ESLint 9 | flat config | Linting (TS + Vue rules) |
| Prettier 3 | — | Formatting (with Tailwind class sorting) |

## Project Layout

```
frontend/
├── eslint.config.js        # ESLint flat config
├── .prettierrc.json        # Prettier config
├── .prettierignore
├── tailwind.config.js
├── vite.config.ts
├── tsconfig.json
├── index.html
└── src/
    ├── App.vue             # Root router — switches between SelectorView / SettingsView
    ├── types.ts            # Shared TypeScript types (mirror of Go structs)
    ├── main.ts             # Vue app bootstrap
    ├── style.css           # Global styles + Tailwind directives
    ├── views/
    │   ├── SelectorView.vue    # Browser picker window (vertical list or horizontal dock)
    │   └── SettingsView.vue    # Configuration hub (browsers, rulesets, app settings)
    ├── components/
    │   ├── BrowserList.vue       # Vertical stacked browser list
    │   ├── BrowserButton.vue     # Single list row (icon, name, shortcut key)
    │   ├── BrowserDock.vue       # Horizontal dock container
    │   ├── BrowserDockTile.vue   # Single dock tile (large icon, keycap badge)
    │   ├── UrlCard.vue           # URL display with protocol icon and copy button
    │   ├── TimerBar.vue          # Countdown bar for timed auto-selection
    │   ├── BrowserEditor.vue     # Browser config row (settings)
    │   └── RulesetEditor.vue     # Ruleset config row (settings)
    ├── composables/
    │   └── useBrowserIcon.ts     # Reactive icon loader — calls GetBrowserIconAt() backend
    └── utils/
        ├── browser.ts    # Shortcut key labels, stagger animation delay helpers
        └── config.ts     # Layout / overflow value normalizers
```

## Development

Run from the **repo root** — Wails manages the Vite dev server automatically:

```powershell
wails dev
```

To work on the frontend in isolation (without the Go backend):

```powershell
cd frontend
npm run dev
# → http://localhost:5173
```

> Backend calls (`window.go.*`) will fail without Wails, but layout and styling can be iterated freely.

## Code Quality

```powershell
npm run lint          # ESLint — report issues
npm run lint:fix      # ESLint — auto-fix
npm run format        # Prettier — reformat all src files
npm run format:check  # Prettier — CI check (no writes)
```

## Wails Integration

Wails generates Go ↔ JS bindings into `wailsjs/` on every `wails build` or `wails dev`. This directory is gitignored — never edit it by hand.

The frontend calls backend methods via the generated runtime:

```ts
import { GetBrowsers, OpenWithBrowser } from '../wailsjs/go/main/App'

const browsers = await GetBrowsers()
await OpenWithBrowser('Firefox')
```

## Recommended IDE Setup

**VS Code** with the [Vue - Official](https://marketplace.visualstudio.com/items?itemName=Vue.volar) extension (Volar). Make sure the built-in TypeScript extension is enabled alongside Volar — Take Over mode is no longer needed in modern Volar.
