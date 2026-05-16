# Known Gotchas

- **Windows-only:** Switchy is Windows-specific. Registry operations, Windows API calls, and path handling assume Windows. No cross-platform support.
- **Single instance via named pipes:** Single-instance enforcement relies on named pipes. On network shares or certain environments, pipe communication may fail.
- **Chromium profile detection:** Profile detection for Chrome, Edge, Brave depends on filesystem paths. Custom profile locations or portable installs may not be detected.
- **Icon extraction:** Browser icon extraction via Windows API may fail for non-standard installations or UWP apps without fallback icon handling.
- **Registry permissions:** Registration/unregistration operations require admin privileges. Graceful error handling for permission denied is important.
