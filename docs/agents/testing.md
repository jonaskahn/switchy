# Testing

## Test Runner

**Go:** `go test ./...`  
**Frontend:** `npm run lint` (ESLint, no unit tests configured)

## Test Layout

Test files are **co-located** with source files using the `_test.go` suffix (Go convention). No separate `tests/` directory.

Example:
```
internal/rules/matcher.go
internal/rules/matcher_test.go
```

## Test Command

```bash
go test ./...
```

Runs all Go tests with coverage. Frontend uses lint-only (no unit test runner configured).

## Mock Stance

[to fill]
