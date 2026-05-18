<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, watch, nextTick } from 'vue'
import { EventsOn, ScreenGetAll } from '@wailsjs/runtime'
import {
  GetURL,
  GetBrowsers,
  GetSettings,
  OpenWithBrowser,
  RememberBrowserForDomain,
  CopyURL,
  DismissWindow,
  OpenSettings,
  ResizeSelector,
} from '@wailsjs/go/main/App'
import type { Browser } from '@/types'
import BrowserList from '@/components/BrowserList.vue'
import BrowserDock from '@/components/BrowserDock.vue'
import UrlCard from '@/components/UrlCard.vue'
import brandMarkUrl from '@/assets/brand-mark.png'
import { normLayout } from '@/utils/config'
import { applyTheme, applyAppearance } from '@/composables/useTheme'

const DOCK_TILE_WIDTH_PX = 110
const DOCK_GAP_PX = 48
const DOCK_SIDE_PADDING_PX = 22
const HORIZONTAL_DOCK_HEIGHT_PX = 220
const FOOTER_HEIGHT_PX = 8

const VERTICAL_WIDTH_PX = 420
const VERTICAL_ROW_HEIGHT_PX = 65

const SCREEN_CAP_RATIO = 0.80
const COPY_FEEDBACK_DURATION_MS = 1500
const TIMER_TICK_INTERVAL_MS = 100
const TIMER_TICK_STEP = 0.1

const NOTICE_DURATION_MS = 3000

const currentUrl = ref('')
const browsers = ref<Browser[]>([])
const loading = ref(true)
const copied = ref(false)
const notice = ref('')
const ctrlHeld = ref(false)
let noticeTimer: ReturnType<typeof setTimeout> | null = null
const timerActive = ref(false)
const timerSecs = ref(0)
const timerRemain = ref(0)
let timerInterval: ReturnType<typeof setInterval> | null = null

const selectorLayout = ref<'vertical' | 'horizontal'>('vertical')

const dragStripEl = ref<HTMLElement | null>(null)
const browserZoneEl = ref<HTMLElement | null>(null)
const urlStripEl = ref<InstanceType<typeof UrlCard> | null>(null)

const timerDisplay = computed(() => timerRemain.value.toFixed(1) + 's')

const displayHostname = computed(() => {
  try {
    return new URL(currentUrl.value).hostname || currentUrl.value
  } catch {
    return currentUrl.value
  }
})
const timerBarScale = computed(() =>
  Math.max(0, Math.min(1, timerRemain.value / Math.max(timerSecs.value, 0.0001)))
)

const windowClass = computed(() => ({
  'selector-window': true,
  'selector-window--horizontal': selectorLayout.value === 'horizontal',
}))

const browserZoneClass = computed(() => ({
  'browser-zone': true,
  'browser-zone--horizontal': selectorLayout.value === 'horizontal',
}))

async function applyWindowResize() {
  if (loading.value) return
  await nextTick()
  await new Promise<void>((resolve) => requestAnimationFrame(() => resolve()))
  await nextTick()

  const strip = dragStripEl.value
  const urlStrip = urlStripEl.value?.root
  const zone = browserZoneEl.value
  if (!urlStrip || !zone) return

  const stripH = strip?.offsetHeight ?? 0

  let count = browsers.value.length
  let delta_vertical_row_height_px: number
  if (count < 3) {
    delta_vertical_row_height_px = 65
  } else if (count < 6) {
    delta_vertical_row_height_px = 40
  } else if (count < 10) {
    delta_vertical_row_height_px = 15
  } else {
    delta_vertical_row_height_px = 0
  }

  let cur: { width: number; height: number } | undefined
  try {
    const screens = await ScreenGetAll()
    cur =
      screens.find((s: { isCurrent: boolean }) => s.isCurrent) ||
      screens.find((s: { isPrimary: boolean }) => s.isPrimary) ||
      screens[0]
  } catch {}

  let w: number
  let h: number

  if (selectorLayout.value === 'horizontal') {
    w = DOCK_SIDE_PADDING_PX * 2 + count * DOCK_TILE_WIDTH_PX + Math.max(0, count - 1) * DOCK_GAP_PX
    if (cur) w = Math.min(w, Math.floor(cur.width * SCREEN_CAP_RATIO))
    h = Math.ceil(stripH + HORIZONTAL_DOCK_HEIGHT_PX + urlStrip.offsetHeight + FOOTER_HEIGHT_PX)
  } else {
    w = cur ? Math.floor(cur.width / 5) : VERTICAL_WIDTH_PX
    h = Math.ceil(
      stripH + urlStrip.offsetHeight + count * VERTICAL_ROW_HEIGHT_PX + delta_vertical_row_height_px
    )
  }

  h = Math.min(h, Math.floor(window.screen.height * 0.85))

  try {
    await ResizeSelector(w, h)
  } catch {}
}

async function scheduleContentFit() {
  await applyWindowResize()
  await nextTick()
  void applyWindowResize()
}

async function load() {
  try {
    loading.value = true
    const [url, brs, cfg] = await Promise.all([GetURL(), GetBrowsers(), GetSettings()])
    currentUrl.value = url ?? ''
    browsers.value = (brs ?? []).filter((b: Browser) => !b.Hidden)

    if (browsers.value.length === 0) {
      openSettings()
      return
    }

    if (cfg?.AppSettings) {
      selectorLayout.value = normLayout(cfg.AppSettings.SelectorLayout as string)
      await applyTheme(cfg.AppSettings.Theme ?? 'default')
      await applyAppearance(cfg.AppSettings.Appearance ?? 'acrylic')
    }

    if (cfg?.AppSettings?.TimedSelection && cfg.AppSettings.TimedSeconds > 0) {
      timerSecs.value = cfg.AppSettings.TimedSeconds
      timerRemain.value = cfg.AppSettings.TimedSeconds
      timerActive.value = true
      startTimer(cfg.AppSettings.DefaultBrowser)
    }
  } finally {
    loading.value = false
    void scheduleContentFit()
  }
}

watch(
  () => [loading.value, browsers.value.length, selectorLayout.value] as const,
  () => void scheduleContentFit()
)

function startTimer(defaultBrowser: string) {
  timerInterval = setInterval(() => {
    timerRemain.value = Math.max(0, timerRemain.value - TIMER_TICK_STEP)
    if (timerRemain.value <= 0) {
      stopTimer()
      if (defaultBrowser) openBrowser(defaultBrowser)
      else DismissWindow()
    }
  }, TIMER_TICK_INTERVAL_MS)
}

function stopTimer() {
  if (timerInterval) {
    clearInterval(timerInterval)
    timerInterval = null
  }
  timerActive.value = false
}

function showNotice(msg: string) {
  if (noticeTimer) clearTimeout(noticeTimer)
  notice.value = msg
  noticeTimer = setTimeout(() => {
    notice.value = ''
    noticeTimer = null
  }, NOTICE_DURATION_MS)
}

async function openBrowser(name: string, remember = false) {
  stopTimer()
  if (remember) {
    try {
      const hostname = await RememberBrowserForDomain(name)
      showNotice(`Remembered ${name} for ${hostname}`)
    } catch {
      showNotice('Cannot remember this URL type')
    }
  }
  await OpenWithBrowser(name)
}

async function copyUrl() {
  await CopyURL()
  copied.value = true
  setTimeout(() => (copied.value = false), COPY_FEEDBACK_DURATION_MS)
}

function dismiss() {
  stopTimer()
  DismissWindow()
}
function openSettings() {
  stopTimer()
  DismissWindow()
  OpenSettings()
}

function onKeyDown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    dismiss()
    return
  }
  if (e.key.toLowerCase() === 'c' && !e.ctrlKey) {
    copyUrl()
    return
  }
  if (e.key === ',') {
    openSettings()
    return
  }
  const idx = parseInt(e.key, 10) - 1
  if (!isNaN(idx) && idx >= 0 && idx < browsers.value.length)
    openBrowser(browsers.value[idx].Name, e.ctrlKey)
}

function onCtrlDown(e: KeyboardEvent) {
  if (e.key === 'Control') ctrlHeld.value = true
}
function onCtrlUp(e: KeyboardEvent) {
  if (e.key === 'Control') ctrlHeld.value = false
}
function onWindowBlur() {
  ctrlHeld.value = false
}
function onMouseMove(e: MouseEvent) {
  ctrlHeld.value = e.ctrlKey
}

onMounted(() => {
  load()
  window.addEventListener('keydown', onKeyDown)
  window.addEventListener('keydown', onCtrlDown)
  window.addEventListener('keyup', onCtrlUp)
  window.addEventListener('blur', onWindowBlur)
  window.addEventListener('mousemove', onMouseMove)
  EventsOn('url:new', (url: string) => {
    currentUrl.value = url
  })
})
onUnmounted(() => {
  window.removeEventListener('keydown', onKeyDown)
  window.removeEventListener('keydown', onCtrlDown)
  window.removeEventListener('keyup', onCtrlUp)
  window.removeEventListener('blur', onWindowBlur)
  window.removeEventListener('mousemove', onMouseMove)
  stopTimer()
})
</script>

<template>
  <div :class="windowClass">
    <div ref="dragStripEl" class="window-drag-strip drag-region" aria-hidden="true" />
    <template v-if="selectorLayout === 'vertical'">
      <UrlCard ref="urlStripEl" :url="currentUrl" :copied="copied" @copy="copyUrl" />
      <div ref="browserZoneEl" :class="browserZoneClass">
        <BrowserList
          :browsers="browsers"
          :loading="loading"
          @open="(name: string, remember: boolean) => openBrowser(name, remember)"
        />
      </div>
    </template>

    <template v-else>
      <div ref="browserZoneEl" :class="browserZoneClass">
        <BrowserDock
          :browsers="browsers"
          :loading="loading"
          @open="(name: string, remember: boolean) => openBrowser(name, remember)"
        />
      </div>
      <UrlCard ref="urlStripEl" :url="currentUrl" :copied="copied" :compact="true" @copy="copyUrl" />
    </template>

    <div class="footer" :class="{ 'footer--horizontal': selectorLayout === 'horizontal' }">
      <div class="footer-row">
        <div class="footer-brand no-drag">
          <img
            class="footer-brand-mark"
            :src="brandMarkUrl"
            width="16"
            height="16"
            alt=""
            draggable="false"
            decoding="sync"
          />
          <span class="footer-brand-text">Switchy</span>
        </div>

        <div class="footer-hints no-drag">
          <template v-if="notice">
            <span class="hint hint--notice">{{ notice }}</span>
          </template>
          <template v-else-if="ctrlHeld">
            <span class="hint hint--ctrl"><kbd>Ctrl</kbd> held — click to remember for {{ displayHostname }}</span>
          </template>
          <template v-else>
            <span class="hint"><kbd>Ctrl</kbd> Remember</span>
            <span class="hint"><kbd>C</kbd> Copy</span>
            <span class="hint"><kbd>Esc</kbd> Close</span>
          </template>
        </div>

        <div class="footer-trailing no-drag">
          <div
            v-if="timerActive"
            class="footer-timer"
            :class="{ 'footer-timer--mockup': selectorLayout === 'horizontal' }"
          >
            <span class="timer-dot"></span>{{ timerDisplay }}
          </div>
          <button
            class="icon-btn icon-btn--footer"
            title="Settings (,)"
            aria-label="Settings"
            @click="openSettings"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.8"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <circle cx="12" cy="12" r="3" />
              <path
                d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"
              />
            </svg>
          </button>
          <button
            class="icon-btn danger icon-btn--footer"
            title="Close (Esc)"
            aria-label="Close"
            @click="dismiss"
          >
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

      <div
        v-if="timerActive && selectorLayout === 'horizontal'"
        class="footer-timer-bar"
        :style="{ transform: `scaleX(${timerBarScale})` }"
      />
    </div>
  </div>
</template>

<style scoped>
/* Window shell */
.selector-window {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: linear-gradient(180deg, var(--substrate-top), var(--substrate));
  backdrop-filter: var(--backdrop-blur);
  -webkit-backdrop-filter: var(--backdrop-blur);
  overflow: hidden;
  animation: rise 240ms cubic-bezier(0.2, 0.8, 0.2, 1) both;
}

@keyframes rise {
  from {
    opacity: 0;
    transform: translateY(8px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.selector-window--horizontal {
  backdrop-filter: var(--backdrop-blur-h);
  -webkit-backdrop-filter: var(--backdrop-blur-h);
  animation: riseH 240ms cubic-bezier(0.2, 0.8, 0.2, 1) both;
}
@keyframes riseH {
  from {
    opacity: 0;
    transform: translateY(6px) scale(0.985);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.window-drag-strip {
  height: 10px;
  flex-shrink: 0;
  cursor: move;
  -webkit-app-region: drag;
}

/* Browser zone */
.browser-zone {
  min-height: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.browser-zone--horizontal {
  flex: 0 0 auto;
}

/* Footer — fixed size, overlays scroll area. */
.footer {
  position: absolute;
  left: 0;
  bottom: 0;
  width: 100%;
  height: 48px;
  z-index: 10;
  border-top: 1px solid var(--hairline);
  background: rgba(0, 0, 0, 0.2);
  padding: 0 14px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  justify-content: center;
}
.footer-row {
  display: flex;
  flex-wrap: nowrap;
  align-items: center;
  gap: 8px 10px;
  width: 100%;
  min-width: 0;
  min-height: 26px;
}
.footer-brand {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}
.footer-brand-mark {
  width: 16px;
  height: 16px;
  object-fit: contain;
  display: block;
  flex-shrink: 0;
}
.footer-brand-text {
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--ink-80);
}
.footer-trailing {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
  margin-left: auto;
}
.icon-btn--footer {
  width: 26px;
  height: 26px;
  border-radius: 6px;
}
.icon-btn--footer svg {
  width: 13px;
  height: 13px;
}
.footer-hints {
  flex: 1 1 0;
  display: flex;
  flex-wrap: nowrap;
  align-items: center;
  justify-content: center;
  gap: 8px 10px;
  min-width: 0;
  overflow: hidden;
}
.footer-timer {
  font-family: var(--font-mono);
  font-size: 10.5px;
  color: var(--ink-40);
  letter-spacing: 0.02em;
  display: flex;
  align-items: center;
  flex-shrink: 0;
}
.footer-timer--mockup {
  font-size: 10px;
}
.hint {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 10px;
  font-weight: 500;
  color: var(--ink-60);
  letter-spacing: 0.02em;
}
.footer--horizontal .hint {
  font-size: 10px;
}
.hint kbd {
  font-family: var(--font-mono);
  font-size: 9.5px;
  font-weight: 500;
  background: var(--ink-10);
  border: 1px solid var(--hairline);
  border-bottom-width: 2px;
  padding: 2px 5px;
  border-radius: 4px;
  color: var(--ink-80);
  line-height: 1;
}
.footer--horizontal .hint kbd {
  font-size: 9.5px;
}
.hint--notice {
  color: var(--accent);
  font-weight: 600;
  letter-spacing: 0.01em;
}
.hint--ctrl {
  color: var(--ink-80);
  font-style: italic;
}
.timer-dot {
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--accent);
  margin-right: 6px;
  box-shadow: 0 0 8px var(--accent);
  animation: pulse 1.6s ease-in-out infinite;
}
.footer-timer--mockup .timer-dot {
  width: 5px;
  height: 5px;
  box-shadow: 0 0 6px var(--accent);
}

.footer-timer-bar {
  position: absolute;
  left: 0;
  bottom: 0;
  height: 1.5px;
  width: 100%;
  background: var(--accent);
  transform-origin: left center;
  transition: transform 100ms linear;
  box-shadow: 0 0 6px var(--accent);
  pointer-events: none;
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.35;
  }
}
</style>
