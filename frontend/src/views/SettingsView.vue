<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  ApplyAppearance,
  GetSettings,
  SaveSettings,
  DetectBrowsers,
  RegisterAsDefault,
  UnregisterAsDefault,
  CloseSettings,
} from '@wailsjs/go/main/App'
import { WindowMinimise, WindowToggleMaximise } from '@wailsjs/runtime/runtime'
import type { config } from '@wailsjs/go/models'
import type { UserSettings } from '@/types'
import { normLayout } from '@/utils/config'
import { applyTheme, applyAppearance } from '@/composables/useTheme'
import BrowserEditor from '@/components/BrowserEditor.vue'
import RulesetEditor from '@/components/RulesetEditor.vue'
import brandMarkUrl from '@/assets/brand-mark.png'

type Section = 'browsers' | 'rulesets' | 'app'

const SAVE_MESSAGE_DURATION_MS = 2000
const ACTION_MESSAGE_DURATION_MS = 3000
const DEFAULT_BROWSER_NAME = 'New Browser'
const DEFAULT_RULESET_NAME = 'New Ruleset'

const activeSection = ref<Section>('browsers')
const settings = ref<UserSettings | null>(null)
const saving = ref(false)
const message = ref('')
const registering = ref(false)
const unsavedCount = ref(0)

const pageInfo = computed(() => {
  if (activeSection.value === 'browsers')
    return {
      eyebrow: 'Configure / 01',
      title: 'Browsers',
      description:
        'Every browser Switchy can route to. Auto-detect scans the Windows registry; manual entries override paths and add alternate launch profiles.',
    }
  if (activeSection.value === 'rulesets')
    return {
      eyebrow: 'Configure / 02',
      title: 'Rulesets',
      description:
        'Pattern-based rules that route matching URLs automatically. Prefix d$ for domain, r$ for regex, s$ for contains; omit a prefix for a simple contains match.',
    }
  return {
    eyebrow: 'Configure / 03',
    title: 'Application',
    description:
      'Global preferences, timed auto-selection, and Windows default browser registration.',
  }
})

const verticalLayout = computed(
  () => normLayout(settings.value?.AppSettings?.SelectorLayout as string) !== 'horizontal'
)

const browserCount = computed(() => String(settings.value?.Browsers.length ?? 0).padStart(2, '0'))

const rulesetCount = computed(() => String(settings.value?.Rulesets.length ?? 0).padStart(2, '0'))

const browserNames = computed(() => settings.value?.Browsers.map((b) => b.Name) ?? [])

function markUnsaved() {
  unsavedCount.value++
}

async function load() {
  settings.value = await GetSettings()
  settings.value.Browsers ??= []
  settings.value.Rulesets ??= []
  unsavedCount.value = 0
  if (settings.value.Browsers.length > 0) activeSection.value = 'browsers'
  await applyTheme(settings.value.AppSettings.Theme ?? 'default')
  await applyAppearance(settings.value.AppSettings.Appearance ?? 'acrylic')
}

function setAppearance(appearance: 'acrylic' | 'mica') {
  if (!settings.value) return
  settings.value.AppSettings.Appearance = appearance
  applyAppearance(appearance)
  ApplyAppearance(appearance)
  markUnsaved()
}

function setTheme(theme: string) {
  if (!settings.value) return
  settings.value.AppSettings.Theme = theme
  applyTheme(theme)
  markUnsaved()
}

async function save() {
  try {
    saving.value = true
    await SaveSettings(settings.value! as unknown as config.UserSettings)
    unsavedCount.value = 0
    message.value = 'Saved.'
    setTimeout(() => (message.value = ''), SAVE_MESSAGE_DURATION_MS)
  } catch (e: unknown) {
    message.value = 'Error: ' + (e instanceof Error ? e.message : String(e))
  } finally {
    saving.value = false
  }
}

async function detect() {
  const found = await DetectBrowsers()
  if (!found?.length) {
    message.value = 'No new browsers found.'
    setTimeout(() => (message.value = ''), ACTION_MESSAGE_DURATION_MS)
    return
  }
  const existingPaths = new Set(settings.value?.Browsers.map((b) => b.ExePath) ?? [])
  const newBrowsers = found.filter((b) => !existingPaths.has(b.ExePath))
  settings.value!.Browsers = [...settings.value!.Browsers, ...newBrowsers]
  markUnsaved()
  message.value = `Added ${newBrowsers.length} browser(s).`
  setTimeout(() => (message.value = ''), ACTION_MESSAGE_DURATION_MS)
}

async function registerDefault() {
  try {
    registering.value = true
    await RegisterAsDefault()
    message.value = 'Registered as default browser.'
    setTimeout(() => (message.value = ''), ACTION_MESSAGE_DURATION_MS)
  } catch (e: unknown) {
    message.value = 'Error: ' + (e instanceof Error ? e.message : String(e))
  } finally {
    registering.value = false
  }
}

async function unregisterDefault() {
  try {
    await UnregisterAsDefault()
    message.value = 'Unregistered.'
    setTimeout(() => (message.value = ''), ACTION_MESSAGE_DURATION_MS)
  } catch (e: unknown) {
    message.value = 'Error: ' + (e instanceof Error ? e.message : String(e))
  }
}

function addBrowser() {
  settings.value!.Browsers.push({
    Name: DEFAULT_BROWSER_NAME,
    ExePath: '',
    IconSpec: '',
    LaunchArgs: '',
    IsUwp: false,
    Hidden: false,
    AlternateLaunches: [],
  })
  markUnsaved()
}

function removeBrowser(idx: number) {
  settings.value!.Browsers.splice(idx, 1)
  markUnsaved()
}

function moveBrowserUp(idx: number) {
  if (idx === 0) return
  const arr = settings.value!.Browsers
  ;[arr[idx - 1], arr[idx]] = [arr[idx], arr[idx - 1]]
  markUnsaved()
}

function moveBrowserDown(idx: number) {
  const arr = settings.value!.Browsers
  if (idx >= arr.length - 1) return
  ;[arr[idx], arr[idx + 1]] = [arr[idx + 1], arr[idx]]
  markUnsaved()
}

function addRuleset() {
  settings.value!.Rulesets.push({ Name: DEFAULT_RULESET_NAME, Enabled: true, Rules: [] })
  markUnsaved()
}

function removeRuleset(idx: number) {
  settings.value!.Rulesets.splice(idx, 1)
  markUnsaved()
}

function toggleUseRules() {
  if (!settings.value) return
  settings.value.AppSettings.UseRules = !settings.value.AppSettings.UseRules
  markUnsaved()
}

function toggleTimedSelection() {
  if (!settings.value) return
  settings.value.AppSettings.TimedSelection = !settings.value.AppSettings.TimedSelection
  markUnsaved()
}

function setLayoutVertical() {
  if (!settings.value) return
  settings.value.AppSettings.SelectorLayout = 'vertical'
  markUnsaved()
}

function setLayoutHorizontal() {
  if (!settings.value) return
  settings.value.AppSettings.SelectorLayout = 'horizontal'
  markUnsaved()
}

onMounted(load)
</script>

<template>
  <div class="settings-window">
    <div class="titlebar drag-region">
      <div class="brand no-drag">
        <img
          class="brand-mark"
          :src="brandMarkUrl"
          width="18"
          height="18"
          alt=""
          draggable="false"
          decoding="sync"
        />
        <span class="brand-name">Switchy</span>
        <span class="brand-divider"></span>
        <span class="brand-sub">Settings</span>
      </div>
      <div class="titlebar-spacer"></div>
      <div class="titlebar-actions no-drag">
        <button class="icon-btn" aria-label="Minimize" @click="WindowMinimise()">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M5 12h14" />
          </svg>
        </button>
        <button class="icon-btn" aria-label="Maximize" @click="WindowToggleMaximise()">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6">
            <rect x="5" y="5" width="14" height="14" rx="1" />
          </svg>
        </button>
        <button class="icon-btn danger" aria-label="Close" @click="CloseSettings()">
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.8"
            stroke-linecap="round"
          >
            <path d="M6 6 L18 18 M18 6 L6 18" />
          </svg>
        </button>
      </div>
    </div>

    <div v-if="settings" class="body">
      <nav class="sidebar no-drag">
        <div class="nav-group">
          <div class="nav-label">Configure</div>

          <div
            class="nav-item"
            :class="{ 'is-active': activeSection === 'browsers' }"
            @click="activeSection = 'browsers'"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.6"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <circle cx="12" cy="12" r="10" />
              <path d="M2 12h20M12 2a15 15 0 0 1 0 20M12 2a15 15 0 0 0 0 20" />
            </svg>
            Browsers
            <span class="nav-count">{{ browserCount }}</span>
          </div>

          <div
            class="nav-item"
            :class="{ 'is-active': activeSection === 'rulesets' }"
            @click="activeSection = 'rulesets'"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.6"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M4 6h16M4 12h16M4 18h10" />
            </svg>
            Rulesets
            <span class="nav-count">{{ rulesetCount }}</span>
          </div>

          <div
            class="nav-item"
            :class="{ 'is-active': activeSection === 'app' }"
            @click="activeSection = 'app'"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.6"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <circle cx="12" cy="12" r="3" />
              <path
                d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"
              />
            </svg>
            Application
          </div>
        </div>

        <div class="sidebar-spacer"></div>

        <div class="sidebar-meta">
          <div class="meta-row"><span>VERSION</span><span>1.0.0</span></div>
          <div class="meta-row">
            <span>DEFAULT</span>
            <span style="color: var(--accent)">● ACTIVE</span>
          </div>
          <div class="meta-row"><span>BUILD</span><span>amd64</span></div>
        </div>
      </nav>

      <main class="content no-drag">
        <header class="page-header">
          <div class="page-title">
            <div class="page-eyebrow">{{ pageInfo.eyebrow }}</div>
            <h1>{{ pageInfo.title }}</h1>
            <p>{{ pageInfo.description }}</p>
          </div>
          <div class="page-actions">
            <template v-if="activeSection === 'browsers'">
              <button class="btn" @click="detect">Auto-detect</button>
              <button class="btn primary" @click="addBrowser">Add browser</button>
            </template>
            <template v-else-if="activeSection === 'rulesets'">
              <button class="btn primary" @click="addRuleset">Add ruleset</button>
            </template>
          </div>
        </header>

        <div class="content-body">
          <div v-if="activeSection === 'browsers'">
            <div v-if="settings.Browsers.length === 0" class="empty-state">
              <svg
                width="48"
                height="48"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <circle cx="12" cy="12" r="10" />
                <path d="M2 12h20M12 2a15 15 0 0 1 0 20M12 2a15 15 0 0 0 0 20" />
              </svg>
              <span>No browsers detected yet.</span>
              <div class="empty-state-actions">
                <button class="btn" @click="detect">Auto-detect</button>
                <button class="btn primary" @click="addBrowser">Add browser</button>
              </div>
            </div>
            <div v-else class="section">
              <div class="section-head">
                <h2>Founded · {{ settings.Browsers.length }}</h2>
              </div>
              <div class="table">
                <BrowserEditor
                  v-for="(browser, idx) in settings.Browsers"
                  :key="idx"
                  v-model="settings.Browsers[idx]"
                  :is-default="idx === 0"
                  @remove="removeBrowser(idx)"
                  @move-up="moveBrowserUp(idx)"
                  @move-down="moveBrowserDown(idx)"
                  @change="markUnsaved"
                />
              </div>
            </div>
          </div>

          <div v-else-if="activeSection === 'rulesets'" class="rulesets-panel">
            <div class="rule-master-toggle">
              <div class="field-text">
                Enable rule matching
                <small>When off, every link shows the picker regardless of rules.</small>
              </div>
              <button
                type="button"
                class="switch"
                :class="{ on: settings.AppSettings.UseRules }"
                :aria-pressed="settings.AppSettings.UseRules"
                @click="toggleUseRules"
              ></button>
            </div>

            <div v-if="settings.Rulesets.length === 0" class="empty-state empty-state-muted">
              <svg
                width="36"
                height="36"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.15"
                stroke-linecap="round"
                stroke-linejoin="round"
                aria-hidden="true"
              >
                <path d="M4 6h16M4 12h16M4 18h10" />
              </svg>
              <span>No rulesets yet. Click "Add ruleset" to route URLs automatically.</span>
            </div>

            <div v-else class="section rulesets-section">
              <div class="section-head">
                <h2>Rulesets · {{ settings.Rulesets.length }}</h2>
              </div>
              <div class="table">
                <RulesetEditor
                  v-for="(ruleset, idx) in settings.Rulesets"
                  :key="idx"
                  v-model="settings.Rulesets[idx]"
                  :browsers="browserNames"
                  @remove="removeRuleset(idx)"
                  @change="markUnsaved"
                />
              </div>
            </div>
          </div>

          <div v-else-if="activeSection === 'app'" class="space-y-4">
            <div class="settings-card">
              <div class="card-title">Appearance</div>
              <div class="field-row-setting field-row-setting--top">
                <div class="field-text">
                  Window material
                  <small
                    >Backdrop effect for all Switchy windows. Reopen to apply OS-level
                    change.</small
                  >
                </div>
                <div class="lf-picker">
                  <button
                    type="button"
                    class="lf-card"
                    :class="{
                      'is-active': (settings.AppSettings.Appearance || 'acrylic') === 'acrylic',
                    }"
                    @click="setAppearance('acrylic')"
                  >
                    <div class="lf-swatch lf-swatch--acrylic">
                      <div class="lf-swatch-bar" />
                      <div class="lf-swatch-row" />
                      <div class="lf-swatch-row lf-swatch-row--short" />
                    </div>
                    <span class="lf-label">Acrylic</span>
                    <span class="lf-desc">Bold · frosted glass</span>
                  </button>
                  <button
                    type="button"
                    class="lf-card"
                    :class="{ 'is-active': settings.AppSettings.Appearance === 'mica' }"
                    @click="setAppearance('mica')"
                  >
                    <div class="lf-swatch lf-swatch--mica">
                      <div class="lf-swatch-bar" />
                      <div class="lf-swatch-row" />
                      <div class="lf-swatch-row lf-swatch-row--short" />
                    </div>
                    <span class="lf-label">Mica</span>
                    <span class="lf-desc">Subtle · desktop-tinted</span>
                  </button>
                </div>
              </div>
            </div>

            <div class="settings-card">
              <div class="card-title">Theme</div>
              <div class="field-row-setting field-row-setting--top">
                <div class="field-text">
                  Color palette
                  <small>Accent colors and ink scale applied across all windows.</small>
                </div>
                <div class="lf-picker">
                  <button
                    type="button"
                    class="lf-card"
                    :class="{
                      'is-active':
                        !settings.AppSettings.Theme || settings.AppSettings.Theme === 'default',
                    }"
                    @click="setTheme('default')"
                  >
                    <div class="lf-palette">
                      <span class="lf-dot" style="background: #db1a1a" />
                      <span class="lf-dot" style="background: #c9b59c" />
                      <span class="lf-dot" style="background: #8b7058" />
                      <span
                        class="lf-dot"
                        style="
                          background: #f9f8f6;
                          box-shadow: inset 0 0 0 1px rgba(249, 248, 246, 0.25);
                        "
                      />
                    </div>
                    <span class="lf-label">Default</span>
                    <span class="lf-desc">Sand · warm tan</span>
                  </button>
                </div>
              </div>
            </div>

            <div class="settings-card">
              <div class="card-title">Selector behaviour</div>
              <div class="field-row-setting">
                <div class="field-text">
                  Timed auto-selection
                  <small>Automatically open the default browser after N seconds.</small>
                </div>
                <button
                  class="switch"
                  :class="{ on: settings.AppSettings.TimedSelection }"
                  @click="toggleTimedSelection"
                ></button>
              </div>
              <div v-if="settings.AppSettings.TimedSelection" class="field-row-setting">
                <div class="field-text">
                  Seconds before auto-open
                  <small>Between 1 and 30.</small>
                </div>
                <input
                  v-model.number="settings.AppSettings.TimedSeconds"
                  type="number"
                  min="1"
                  max="30"
                  class="s-input s-input--number"
                  @input="markUnsaved"
                />
              </div>
              <div class="field-row-setting">
                <div class="field-text">
                  Default browser on timeout
                  <small>Opened automatically when the timer expires.</small>
                </div>
                <select
                  v-model="settings.AppSettings.DefaultBrowser"
                  class="s-input s-input--select"
                  @change="markUnsaved"
                >
                  <option value="">— none —</option>
                  <option v-for="b in settings.Browsers" :key="b.Name" :value="b.Name">
                    {{ b.Name }}
                  </option>
                </select>
              </div>
              <p v-if="settings.Browsers.length > 0" class="help-text" style="margin-top: 2px">
                "Default" here means the browser Switchy launches on its own when the countdown
                reaches zero — not your Windows system default browser.
              </p>
            </div>

            <div class="settings-card">
              <div class="card-title">Selector layout</div>
              <div class="field-row-setting">
                <div class="field-text">
                  Layout
                  <small
                    >Horizontal matches the dock-row picker; vertical is the stacked list (reopen
                    the picker after saving).</small
                  >
                </div>
                <div class="layout-toggle">
                  <button
                    type="button"
                    class="btn"
                    :class="{ primary: verticalLayout }"
                    @click="setLayoutVertical"
                  >
                    Vertical
                  </button>
                  <button
                    type="button"
                    class="btn"
                    :class="{ primary: !verticalLayout }"
                    @click="setLayoutHorizontal"
                  >
                    Horizontal
                  </button>
                </div>
              </div>
            </div>

            <div class="settings-card">
              <div class="card-title">Windows Integration</div>
              <div class="flex flex-wrap gap-2" style="margin-bottom: 12px">
                <button class="btn primary" :disabled="registering" @click="registerDefault">
                  {{ registering ? 'Registering…' : 'Register & Open Default Apps' }}
                </button>
                <button class="btn btn--text-danger" @click="unregisterDefault">Unregister</button>
              </div>
              <p class="help-text">
                Windows requires you to confirm in <strong>Default Apps</strong> settings after
                clicking Register. Admin rights may be required.
              </p>
            </div>
          </div>
        </div>

        <div class="statusbar">
          <span>
            <span class="status-dot"></span>
            Switchy is the active browser handler.
          </span>
          <span class="status-spacer"></span>
          <span v-if="message" class="status-message">
            {{ message }}
          </span>
          <span v-else-if="unsavedCount > 0" class="status-unsaved">
            UNSAVED CHANGES · {{ unsavedCount }}
          </span>
          <button class="btn primary btn-save" :disabled="saving" @click="save">
            {{ saving ? 'Saving…' : 'Save' }}
          </button>
        </div>
      </main>
    </div>

    <div v-else class="loading-state">Loading…</div>
  </div>
</template>

<style scoped>
/* ── Window shell ─────────────────────────────────────────────────────────── */
.settings-window {
  width: 100%;
  height: 100%;
  display: grid;
  grid-template-rows: 36px 1fr;
  background: linear-gradient(180deg, var(--substrate-top), var(--substrate));
  backdrop-filter: var(--backdrop-blur);
  -webkit-backdrop-filter: var(--backdrop-blur);
  overflow: hidden;
}

/* ── Titlebar ─────────────────────────────────────────────────────────────── */
.titlebar {
  display: flex;
  align-items: center;
  padding: 0 12px;
  border-bottom: 1px solid var(--hairline);
}
.brand {
  display: flex;
  align-items: center;
  gap: 8px;
}
.brand-mark {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
  object-fit: contain;
  display: block;
}
.brand-name {
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--ink-80);
}
.brand-divider {
  width: 1px;
  height: 12px;
  background: var(--ink-20);
  margin: 0 12px;
}
.brand-sub {
  font-size: 11px;
  font-weight: 500;
  color: var(--ink-60);
  letter-spacing: 0.04em;
}
.titlebar-spacer {
  flex: 1;
}
.titlebar-actions {
  display: flex;
  gap: 2px;
}

/* ── Body ─────────────────────────────────────────────────────────────────── */
.body {
  display: grid;
  grid-template-columns: 220px 1fr;
  min-height: 0;
}

/* ── Sidebar ──────────────────────────────────────────────────────────────── */
.sidebar {
  background: var(--sidebar);
  border-right: 1px solid var(--hairline);
  padding: 18px 12px;
  display: flex;
  flex-direction: column;
  min-height: 0;
  overflow-y: auto;
}
.nav-group {
  margin-bottom: 8px;
}
.nav-label {
  font-size: 9.5px;
  font-weight: 600;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: var(--ink-40);
  padding: 8px 10px 6px;
}
.sidebar-spacer {
  flex: 1;
}
.sidebar-meta {
  padding: 10px 10px 4px;
  border-top: 1px solid var(--hairline);
  margin-top: 8px;
}
.meta-row {
  display: flex;
  justify-content: space-between;
  font-size: 10.5px;
  color: var(--ink-40);
  font-family: var(--font-mono);
  padding: 3px 0;
}
.meta-row span:last-child {
  color: var(--ink-60);
}

/* ── Content area ─────────────────────────────────────────────────────────── */
.content {
  display: grid;
  grid-template-rows: auto 1fr auto;
  min-height: 0;
}

/* Page header */
.page-header {
  padding: 24px 32px 18px;
  display: flex;
  align-items: flex-end;
  gap: 16px;
  border-bottom: 1px solid var(--hairline);
  flex-shrink: 0;
}
.page-title {
  flex: 1;
  min-width: 0;
}
.page-eyebrow {
  font-size: 9.5px;
  font-weight: 600;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: var(--ink-40);
  margin-bottom: 6px;
}
.page-title h1 {
  font-size: 22px;
  font-weight: 600;
  letter-spacing: -0.018em;
  color: var(--ink-100);
  line-height: 1.1;
}
.page-title p {
  margin-top: 6px;
  font-size: 12.5px;
  color: var(--ink-60);
  max-width: 56ch;
  line-height: 1.45;
}
.page-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

/* Content body */
.content-body {
  overflow-y: auto;
  padding: 22px 32px 28px;
  min-height: 0;
}

/* Section */
.section + .section {
  margin-top: 24px;
}
.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 4px 10px;
}
.section-head h2 {
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: var(--ink-40);
}
.table {
  border: 1px solid var(--hairline);
  border-radius: var(--r-row);
  overflow: hidden;
  background: var(--fill-layer-alt);
}

/* Status bar */
.statusbar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 8px 16px;
  border-top: 1px solid var(--hairline);
  background: rgba(0, 0, 0, 0.2);
  font-size: 10.5px;
  color: var(--ink-40);
  font-family: var(--font-mono);
  flex-shrink: 0;
}
.status-dot {
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--accent);
  box-shadow: 0 0 6px var(--accent);
  margin-right: 6px;
  vertical-align: middle;
}
.status-spacer {
  flex: 1;
}
.status-message {
  color: var(--accent);
  font-family: var(--font-mono);
  font-size: 10.5px;
}
.status-unsaved {
  color: var(--ink-60);
  font-family: var(--font-mono);
  font-size: 10.5px;
}
.btn-save {
  height: 26px;
  padding: 0 14px;
}

/* App settings */
.settings-card {
  padding: 18px 20px;
  border-radius: var(--r-row);
  background: rgb(from var(--color-sand) r g b / 0.03);
  border: 1px solid var(--hairline);
  margin-bottom: 16px;
}
.card-title {
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: var(--ink-40);
  margin-bottom: 14px;
}
.field-row-setting {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 0;
}
.field-text {
  font-size: 12.5px;
  color: var(--ink-100);
  flex: 1;
}
.field-text small {
  display: block;
  margin-top: 2px;
  font-size: 11px;
  color: var(--ink-60);
  font-weight: 400;
}
.help-text {
  font-size: 11px;
  color: var(--ink-40);
  line-height: 1.5;
}
.help-text strong {
  color: var(--ink-60);
  font-weight: 600;
}

.s-input--number {
  width: 64px;
  text-align: center;
}
.s-input--select {
  width: auto;
  padding: 0 8px;
}
.btn--text-danger {
  color: var(--color-danger);
}

.rulesets-panel {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.rule-master-toggle {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 14px 18px;
  border-radius: var(--r-row);
  background: rgb(from var(--color-sand) r g b / 0.03);
  border: 1px solid var(--hairline);
}

.rulesets-section {
  margin-bottom: 0;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 52px 24px;
  text-align: center;
  font-size: 14px;
  color: var(--ink-60);
}
.empty-state svg {
  color: var(--ink-40);
}
.empty-state-actions {
  display: flex;
  gap: 10px;
  margin-top: 8px;
}

.empty-state-muted {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
  padding: 36px 24px;
  border: 1px dashed var(--hairline);
  border-radius: var(--r-row);
  color: var(--ink-60);
}

.empty-state-muted svg {
  color: var(--ink-40);
  flex-shrink: 0;
}

.empty-state-muted span {
  max-width: 42ch;
  line-height: 1.45;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  font-size: 13px;
  color: var(--ink-60);
}

.space-y-3 > * + * {
  margin-top: 12px;
}
.space-y-4 > * + * {
  margin-top: 16px;
}
.flex {
  display: flex;
}
.items-center {
  align-items: center;
}
.gap-2 {
  gap: 8px;
}
.flex-wrap {
  flex-wrap: wrap;
}
.text-sm {
  font-size: 13px;
}

.layout-toggle {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

/* ── Look & Feel pickers (Appearance + Theme) ─────────────────────────────── */
.field-row-setting--top {
  align-items: flex-start;
  padding-top: 0;
}

.lf-picker {
  display: flex;
  gap: 10px;
  flex-shrink: 0;
  flex-wrap: wrap;
}

.lf-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 10px;
  width: 108px;
  border-radius: var(--r-row);
  border: 1px solid var(--hairline-hi);
  background: var(--ink-10);
  color: var(--ink-60);
  cursor: pointer;
  font-family: var(--font-ui);
  transition:
    background 140ms,
    border-color 140ms,
    color 140ms;
}
.lf-card:hover {
  background: var(--ink-20);
  color: var(--ink-100);
}
.lf-card.is-active {
  border-color: var(--accent);
  background: var(--accent-tint);
  color: var(--ink-100);
}

/* Appearance swatches */
.lf-swatch {
  width: 88px;
  height: 56px;
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 8px 8px 0;
}
.lf-swatch--acrylic {
  background: linear-gradient(
    160deg,
    rgb(from var(--color-surface-hi) r g b / 0.88),
    rgb(from var(--color-surface) r g b / 0.92)
  );
  box-shadow: inset 0 0 0 1px rgb(from var(--color-sand) r g b / 0.1);
}
.lf-swatch--mica {
  background: linear-gradient(
    160deg,
    rgb(from var(--color-surface-hi) r g b / 0.42),
    rgb(from var(--color-surface) r g b / 0.46)
  );
  box-shadow: inset 0 0 0 1px rgb(from var(--color-sand) r g b / 0.06);
}
.lf-swatch-bar {
  height: 6px;
  border-radius: 3px;
  background: rgb(from var(--color-sand) r g b / 0.16);
  flex-shrink: 0;
}
.lf-swatch--mica .lf-swatch-bar {
  background: rgb(from var(--color-sand) r g b / 0.09);
}
.lf-swatch-row {
  height: 5px;
  border-radius: 2.5px;
  background: rgb(from var(--color-sand) r g b / 0.1);
  flex-shrink: 0;
}
.lf-swatch--mica .lf-swatch-row {
  background: rgb(from var(--color-sand) r g b / 0.06);
}
.lf-swatch-row--short {
  width: 60%;
}

.lf-card.is-active .lf-swatch-bar,
.lf-card.is-active .lf-swatch-row {
  background: rgb(from var(--accent) r g b / 0.3);
}
.lf-card.is-active .lf-swatch--mica .lf-swatch-bar,
.lf-card.is-active .lf-swatch--mica .lf-swatch-row {
  background: rgb(from var(--accent) r g b / 0.18);
}

/* Theme (palette) swatches */
.lf-palette {
  width: 88px;
  height: 56px;
  border-radius: 8px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr;
  gap: 4px;
  padding: 8px;
  background: rgb(from var(--color-sand) r g b / 0.04);
  box-shadow: inset 0 0 0 1px rgb(from var(--color-sand) r g b / 0.08);
}
.lf-dot {
  border-radius: 50%;
  display: block;
  width: 20px;
  height: 20px;
  justify-self: center;
  align-self: center;
}

/* Shared labels */
.lf-label {
  font-size: 12px;
  font-weight: 600;
  letter-spacing: -0.005em;
}
.lf-desc {
  font-size: 10px;
  font-weight: 400;
  color: var(--ink-40);
  text-align: center;
}
.lf-card.is-active .lf-desc {
  color: var(--ink-60);
}
</style>
