# Switchy

A lightning-fast Windows browser picker that intercepts links and routes them to your choice of browser with zero friction using Wails, Vue 3, and Go.

A Go codebase built with Wails, Tailwind CSS, and Vue.

<!-- agent-context v0.0.8 | graph e30e2f0 | 2026-05-16 | regenerate: /agent-context --force -->

## 1. Commands

**One way to run things. Don't invent alternatives.**

```
wails dev
wails build
go fmt ./...
go vet ./...
cd frontend && npm run lint:fix
```

The test: a fresh clone should run green after pasting the install and dev commands.

## 2. Boundaries

**Three tiers. No exceptions, no shortcuts.**

Always:
- run `go fmt` before committing.
- use the documented test command above.
- ask before changing files outside the task scope.

Ask first:
- schema migrations or changes to shared config.

Never:
- commit secrets, `.env` files, or credentials.
- edit or delete applied migrations.
- run destructive commands without explicit approval.
- push `--force` to `main`.

## 3. Module Map

**Layers are disjoint. Don't blur them.**

- Business Logic & Rules (6) — Core application logic: URL routing rules, browser selection, configuration management.
- Windows Integration (6) — Windows-specific functionality: registry operations, API calls, window management.
- Frontend (Vue 3) (8) — Vue 3 user interface components, browser selection UI, settings UI.
- Data Access & Configuration (4) — Configuration persistence, browser inventory, settings management.
- Process Management (4) — Browser launching, process spawning, execution control.
- UI Framework & Wails Bridge (2) — Wails framework integration, backend bindings to Vue frontend.
- Entry Point (6) — Process bootstrap and command-line argument routing.
- Utilities & Cross-Cutting Concerns (2) — Logging, browser utilities, common helpers.

The test: every file under `internal/` maps to exactly one layer above.

## 4. Architectural Altitude

**Business Logic & Rules is the main stage. Windows Integration is the backstage.**

- To understand URL interception & routing, start at `main.go`.
- To understand browser selection & launching, start at `app.go`.
- To understand URL routing rules & automation, start at `internal/rules/matcher.go`.
- To understand settings & configuration management, start at `internal/config/service.go`.
- To understand browser detection & inventory, start at `internal/browser/launch.go`.

The test: open AGENTS.md cold, name the top two entry points without scrolling.

## 5. Non-Obvious Conventions

**Match existing shape. Don't normalise the outliers.**

- Business Logic & Rules rarely imports Utilities & Cross-Cutting Concerns directly — the one exception is `app.go` importing `internal/logger/logger.go`. Don't remove it.
- Most Frontend (Vue 3) files end in `.vue`. Don't rename the exceptions — `frontend/src/utils/browser.ts`, `frontend/src/composables/useTheme.ts` — they are intentional.
- `internal/browser/icon_windows.go` lives under Process Management by path but is classified as Windows Integration. Don't reorganise it.
- Most functions are async — don't mix sync I/O in async call chains.
- Functions in Business Logic & Rules start with uppercase (OpenWithBrowser, GetBrowsers) — match it when adding new ones.

The test: grep for the convention in two more places before assuming it holds.

## 6. Absolute rules

**Read and follow. No exceptions, no workarounds.**

### Safety

- MUST NOT commit secrets, `.env` files, or credentials.
- MUST NOT edit migrations after they have been applied.
- MUST NOT disable tests to make them pass.
- MUST NOT run destructive commands without explicit human approval.
- When a hook blocks a command, stop and ask — never work around it.

### While coding

- MUST NOT add abstractions beyond what is planned.
- MUST NOT improve or refactor adjacent unrelated code.
- MUST state assumptions explicitly; if uncertain, ask before proceeding.

### Project-specific

- MUST run `gofmt` on Go code before completion.
- MUST NOT disable, skip, or bypass tests to make code pass.
- Keep functions small and focused on one responsibility.

## 7. Deeper Context

**AGENTS.md is the kernel. Below it, read on demand.**

- @docs/agents/architecture.md — project overview, stack, quick start, layer map.
- @docs/agents/flow.md — entry points, business flows, execution paths.
- @docs/agents/patterns.md — recurring patterns with file:line exemplars.
- @docs/agents/conventions.md — AI-targeted coding directives.
- @docs/agents/testing.md — runner, layout, mock stance.
- @docs/agents/tech-debt.md — known gotchas.

The test: if the answer is in AGENTS.md, don't open `docs/agents/`.

---

Working if: agents stop asking "where does X live?", hook denials are respected, and PRs match the conventions above without being told.
