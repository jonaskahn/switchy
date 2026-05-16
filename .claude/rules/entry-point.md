---
paths:
  - "main.go"
---

# Entry Point Layer

Process bootstrap and command-line argument routing.

## Key Responsibilities

- Program entry point (main function)
- CLI argument parsing (--register, --unregister, --settings)
- Mode dispatch to selector, settings, registration, or unregistration
- Wails application initialization and runtime setup
- Rule evaluation and automatic browser routing

## Conventions for This Layer

- All mode routing centralizes in main.go — don't scatter mode logic elsewhere
- Early argument parsing and dispatch — don't pass raw arguments deep into the call stack
- Rule matching happens here before selector UI — use `autoRule` for early return
- Entry point starts the Wails app — keep app bootstrap here
