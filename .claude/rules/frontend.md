---
paths:
  - "frontend/src/**/*"
---

# Frontend (Vue 3) Layer

Vue 3 user interface components, browser selection UI, settings UI, real-time theme/appearance switching.

## Key Responsibilities

- Browser selector UI (list, selection, keyboard shortcuts)
- Settings configuration interface (rules editor, appearance toggles)
- Real-time theme switching and appearance management
- Wails bridge communication with Go backend
- Tailwind CSS styling and responsive layout

## Conventions for This Layer

- Most Frontend (Vue 3) files end in `.vue`. Don't rename the exceptions — `frontend/src/utils/browser.ts`, `frontend/src/composables/useTheme.ts` — they are intentional.
- Vue components use PascalCase (SelectorView.vue, BrowserButton.vue)
- Composables use camelCase with `use` prefix (useTheme.ts)
- Utility files camelCase (browser.ts)
- All backend calls go through Wails bindings — don't make direct API calls
- Keep component state minimal; prefer computed properties and reactive composables
