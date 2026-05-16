---
paths:
  - "app.go"
  - "internal/rules/**/*"
  - "internal/config/**/*"
---

# Business Logic & Rules Layer

Core application logic: URL routing rules, browser selection, configuration management.

## Key Responsibilities

- URL-to-browser rule evaluation (domain, regex, string matching)
- Browser selection and launch coordination
- Configuration persistence and caching
- Settings management interface

## Conventions for This Layer

- All rule matching logic centralizes in `internal/rules/matcher.go`
- Configuration service (Load, Save, Get) in `internal/config/service.go`
- Avoid duplicating rule evaluation — use the matcher engine, don't re-implement
- Configuration defaults applied when file missing — don't fail on missing config
- Functions in Business Logic & Rules start with uppercase (OpenWithBrowser, GetBrowsers) — match it when adding new ones
