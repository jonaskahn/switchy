# Business Flows

Eight core user and system flows orchestrate Switchy's behavior.

## URL Interception → Manual Browser Selection

User clicks a link, Switchy intercepts it. No matching rule found or rules disabled. User sees browser picker and selects one.

Entry: User clicks http(s) link in browser, triggering registered Switchy handler

- URL handler invocation — Windows registry directs link to Switchy with URL argument
- IPC forward attempt — new process forwards URL to existing instance if running
- Automatic rule matching — system evaluates enabled rules (domain, regex, containment match)
- Selector UI display — Wails launches with Vue selector component, loads settings and browser list
- Browser selection UI interaction — user selects browser via click or keyboard (1-9)
- Browser process launch — selected browser spawned with URL and launch arguments
- Switchy exit — selector window closes, Switchy terminates

Exit: Selected browser opens with URL, Switchy exits

## URL Interception → Automatic Rule Routing

User clicks link, Switchy intercepts it. Matching rule found and rules enabled. Browser launches automatically without UI.

Entry: User clicks http(s) link in browser, triggering registered Switchy handler

- URL handler invocation
- IPC forward attempt
- Automatic rule matching → match found
- Browser process launch
- Switchy exit

Exit: Matched browser opens with URL, Switchy exits silently

## Browser Detection & Inventory

System scans for installed browsers, extracts metadata, builds inventory with icons and launch configurations. Triggered on startup or manual refresh.

Entry: App startup or user triggers 'Detect Browsers' action

- Windows registry scan — searches HKLM and HKCU for StartMenuInternet registry entries
- Browser metadata extraction — reads display name, executable path, icon spec from each entry
- Chromium profile detection — detects individual user profiles for Chrome, Edge, Brave
- Browser icon extraction — extracts icons via Windows API, returns base64 PNG
- Browser deduplication — removes duplicate entries (same name+exe from multiple registry locations)
- Inventory caching — saves browser list to config file and memory

Exit: Browser list saved to configuration, available for selection

## Settings Configuration & Persistence

User opens Settings UI, modifies rules, browser visibility, appearance, and saves. Configuration persisted to disk with validation.

Entry: User opens Settings window (--settings flag or from selector)

- Settings UI load — Wails launches in settings mode, loads config from disk
- Configuration load — current settings read from JSON file with defaults applied
- Settings display — Vue form populated with rules, browser toggles, appearance options
- User edit rules — user adds, removes, or edits pattern-to-browser mappings
- User edit appearance — user changes theme, backdrop effect, layout preference
- Configuration save — modified settings serialized to JSON and written to disk
- Switchy exit

Exit: Configuration persisted, changes take effect on next URL

## Registration as Default URL Handler

Admin operation to register Switchy as the system's http/https protocol handler. Modifies Windows registry to route clicked links through Switchy.

Entry: User runs: switchy.exe --register

- Write registry client entry — creates HKLM SOFTWARE\Clients\StartMenuInternet\Switchy
- Write registry ProgID entry — creates HKLM SOFTWARE\Classes\SwitchyHTML
- Write registry capabilities entry — adds capabilities key for Default Apps visibility
- Write registry URL associations — registers http and https to route through SwitchyHTML
- Write registry shell command entry — registers shell\open\command with Switchy.exe and URL placeholder
- Write registry application entry — adds Switchy to RegisteredApplications
- Open default apps settings — launches Windows Settings for user confirmation

Exit: Windows Settings opens for user confirmation

## Unregistration as Default URL Handler

Admin operation to remove Switchy from system URL protocol handling. Cleans up Windows registry entries.

Entry: User runs: switchy.exe --unregister

- Registry cleanup: URL associations — deletes URLAssociations registry key
- Registry cleanup: capabilities — deletes Capabilities registry key
- Registry cleanup: shell commands — deletes shell\open\command entries
- Registry cleanup: client entry — deletes Switchy entry from StartMenuInternet
- Registry cleanup: ProgID entry — deletes SwitchyHTML ProgID registry key
- Registry cleanup: registered applications — removes Switchy from RegisteredApplications

Exit: Registry cleaned, Switchy no longer handles links

## Single Instance Enforcement

When user clicks a link and another Switchy instance is already running, the new invocation forwards the URL to existing instance and exits.

Entry: User clicks link while Switchy selector already open

- Second process spawned — Windows registry spawns new Switchy.exe with URL argument
- IPC connection attempt — new process attempts named pipe connection with 2-second timeout
- URL message encoding — URL wrapped in JSON Message struct
- IPC URL forward — encoded message transmitted to listening process
- First process receives message — existing instance receives URL from named pipe
- First process updates URL — running instance updates internal state, triggers UI event
- Second process exits — new process exits after successful forward

Exit: URL added to existing selector, second process exits

## Alternate Browser Launch Profile

User selects an alternate launch configuration for a browser (e.g., Private/Incognito mode). Browser launches with profile-specific arguments.

Entry: User selects alternate profile option in browser selector UI

- Alternate profile selection — user clicks context menu or secondary button in browser button
- Alternate config lookup — backend searches browser's AlternateLaunches array for matching config
- Alternate arguments build — launch arguments combined with URL to create process command line
- Browser launch with alternate profile — browser launched with alternate arguments (--private-window, etc)
- Switchy exit

Exit: Browser opens with alternate launch arguments (private mode, etc)
