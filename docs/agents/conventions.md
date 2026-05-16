# Coding Conventions

This document distills the project's human-authored standards in `docs/CONVENTIONS.md` into AI-targeted directives.

## Safety

- MUST NOT commit secrets, `.env` files, or credentials. Load secrets from environment variables: `os.Getenv("API_TOKEN")`.
- MUST NOT disable, skip, or bypass tests to make code pass. If a test is flaky, fix the underlying issue, don't skip.
- MUST run `gofmt` on Go code before completion. Go code must pass `gofmt` formatting.

## Naming

- Name functions by abstraction level: high-level functions broad names, helpers specific. E.g., `Parse()` over `DetermineFileExtensionAndParseConfigurationFile()`.
- Name variables by scope: wider scope → descriptive, narrow scope → short. E.g., `browser` in a loop, not `browserWithConfiguredLaunchArguments`.
- MUST NOT name variables after their types when domain meaning is clearer. E.g., `browserName` not `stringValue`.

## Patterns

- Follow KISS: prefer the simplest design solving the current task. Don't add factories, builders, or abstractions until needed.
- Follow DRY: remove meaningful duplication, but don't create premature abstractions. Three lines is OK if each serves distinct purpose.
- Keep functions small and focused on one responsibility. Bad: `SaveAndLaunch()`. Good: separate `Save()` and `Launch()`.
- Prefer early returns over deep nesting. Guard clauses at function start reduce indentation and improve readability.
- Declare variables close to their first use, not at function start.
- Keep mutable scope small; avoid global mutable state. Pass settings as parameters, don't use package-level vars.
- Keep function signatures small. Use option structs when parameter lists exceed 3–4 items.
- Prefer concrete return types; accept small interfaces where useful. Return `Browser`, not `interface{}`.
- Avoid exposing `interface{}` or `any` unless required at system boundaries (JSON decoding, reflection).
- Use pointers intentionally. Avoid expanding mutable ownership unnecessarily. Return copies of values when possible.

## Workflow

- MUST NOT add explanatory comments inside implementation code. Self-explain through clear naming and small functions.
- Code should explain itself through names, structure, and small functions. E.g., `hasLaunchPath()` not `p()`.
- Only add comments for exported public APIs, interfaces, or genuinely non-obvious why-level constraints. Skip "this stores browser data" comments.

## Related Files

Existing `docs/CONVENTIONS.md` for full examples and detailed rationale.
