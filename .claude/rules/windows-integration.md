---
paths:
  - "internal/registry/**/*"
  - "internal/browser/**/*"
  - "internal/ipc/**/*"
  - "internal/dwm/**/*"
---

# Windows Integration Layer

Windows-specific functionality: registry operations, API calls, window management, browser detection, IPC.

## Key Responsibilities

- Registry manipulation for protocol handler registration
- Browser detection via Windows registry scanning
- Chromium profile detection and enumeration
- Single-instance enforcement via named pipes
- DWM backdrop effects and window management
- Browser icon extraction via Windows API

## Conventions for This Layer

- Windows Integration rarely imports Utilities & Cross-Cutting Concerns directly — the one exception is importing logger.go. Don't remove it.
- Registry operations isolated in `internal/registry/registry.go` — don't scatter registry code across other modules
- Browser launching and detection in `internal/browser/launch.go` — centralize browser-specific logic here
- Named pipe IPC in `internal/ipc/pipe.go` — don't re-implement pipe handling
- All Windows API calls assume Windows platform — document any assumptions about registry structure or API availability
