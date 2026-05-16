<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { Browser } from '@/types'
import { useBrowserIcon } from '@/composables/useBrowserIcon'

const DEFAULT_ALTERNATE_PROFILE_NAME = 'Profile'

const props = defineProps<{
  modelValue: Browser
  isDefault?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [Browser]
  remove: []
  'move-up': []
  'move-down': []
  change: []
}>()

const expanded = ref(false)
const browserRef = computed(() => props.modelValue)
const { iconDataUrl, loadIcon } = useBrowserIcon(browserRef)

function update(partial: Partial<Browser>) {
  emit('update:modelValue', { ...props.modelValue, ...partial })
  emit('change')
}

function setField(key: keyof Browser, e: Event) {
  const el = e.target as HTMLInputElement
  update({ [key]: el.type === 'checkbox' ? el.checked : el.value } as Partial<Browser>)
}

function addAlternate() {
  update({
    AlternateLaunches: [
      ...(props.modelValue.AlternateLaunches || []),
      { Name: DEFAULT_ALTERNATE_PROFILE_NAME, LaunchArgs: '' },
    ],
  })
}
function removeAlternate(idx: number) {
  const alts = [...(props.modelValue.AlternateLaunches || [])]
  alts.splice(idx, 1)
  update({ AlternateLaunches: alts })
}
function setAltName(idx: number, e: Event) {
  const alts = [...(props.modelValue.AlternateLaunches || [])]
  alts[idx] = { ...alts[idx], Name: (e.target as HTMLInputElement).value }
  update({ AlternateLaunches: alts })
}
function setAltArgs(idx: number, e: Event) {
  const alts = [...(props.modelValue.AlternateLaunches || [])]
  alts[idx] = { ...alts[idx], LaunchArgs: (e.target as HTMLInputElement).value }
  update({ AlternateLaunches: alts })
}

const altCount = computed(() => props.modelValue.AlternateLaunches?.length ?? 0)

onMounted(loadIcon)
</script>

<template>
  <div class="row" :class="{ 'is-open': expanded }">
    <div class="row-head" @click="expanded = !expanded">
      <svg
        class="chev"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2.2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path d="M9 6l6 6-6 6" />
      </svg>

      <div class="row-icon">
        <img v-if="iconDataUrl" :src="iconDataUrl" decoding="async" :alt="modelValue.Name" />
        <svg
          v-else
          viewBox="0 0 24 24"
          width="18"
          height="18"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          stroke-linecap="round"
          stroke-linejoin="round"
          style="color: var(--ink-60)"
        >
          <circle cx="12" cy="12" r="10" />
          <path d="M12 2a14.5 14.5 0 0 0 0 20 14.5 14.5 0 0 0 0-20" />
          <path d="M2 12h20" />
        </svg>
        <img
          v-if="modelValue.ProfileIcon"
          class="profile-avatar"
          :src="modelValue.ProfileIcon"
          decoding="async"
          alt="Profile"
        />
      </div>

      <div class="row-meta" @click.stop>
        <div class="row-name">{{ modelValue.Name || 'Unnamed' }}</div>
        <div class="row-path">{{ modelValue.ExePath || '—' }}</div>
      </div>

      <span v-if="isDefault" class="tag default" @click.stop>DEFAULT</span>
      <span v-else-if="altCount > 0" class="tag" @click.stop
        >{{ altCount }} PROFILE{{ altCount > 1 ? 'S' : '' }}</span
      >
      <span v-else-if="modelValue.Hidden" class="tag" @click.stop>HIDDEN</span>
      <span v-else-if="modelValue.IsUwp" class="tag" @click.stop>UWP</span>
      <span v-else class="tag" style="opacity: 0" aria-hidden="true" @click.stop>&nbsp;</span>

      <div class="row-actions" @click.stop>
        <button class="icon-btn" aria-label="Move up" title="Move up" @click="emit('move-up')">
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.8"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M18 15l-6-6-6 6" />
          </svg>
        </button>
        <button
          class="icon-btn"
          aria-label="Move down"
          title="Move down"
          @click="emit('move-down')"
        >
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.8"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M6 9l6 6 6-6" />
          </svg>
        </button>
        <button class="icon-btn danger" aria-label="Delete" title="Delete" @click="emit('remove')">
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.6"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path
              d="M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"
            />
          </svg>
        </button>
      </div>
    </div>

    <div v-show="expanded" class="row-body">
      <div class="field-grid">
        <div class="field-label">Display name</div>
        <input
          class="s-input"
          :value="modelValue.Name"
          placeholder="My Browser"
          @input="setField('Name', $event)"
        />

        <div class="field-label">Executable</div>
        <div class="input-with-btn">
          <input
            :value="modelValue.ExePath"
            placeholder="C:\Program Files\..."
            @input="setField('ExePath', $event)"
          />
          <button type="button">BROWSE</button>
        </div>

        <div class="field-label">Launch arguments</div>
        <input
          class="s-input"
          :value="modelValue.LaunchArgs"
          placeholder='e.g. --new-window "%1"'
          @input="setField('LaunchArgs', $event)"
        />

        <div class="field-label">Profile variant</div>
        <input
          class="s-input"
          :value="modelValue.IconSpec"
          placeholder="Optional · e.g. Work, Private"
          @input="setField('IconSpec', $event)"
        />

        <div class="field-label">Behaviour</div>
        <div>
          <div class="field-row-b">
            <div class="field-text-b">
              Show in selector
              <small>If disabled, this browser is reachable only via rulesets.</small>
            </div>
            <button
              class="switch"
              :class="{ on: !modelValue.Hidden }"
              @click="update({ Hidden: !modelValue.Hidden })"
            ></button>
          </div>
          <div class="field-row-b">
            <div class="field-text-b">
              UWP / Store app
              <small>Enable if the browser is installed from the Microsoft Store.</small>
            </div>
            <button
              class="switch"
              :class="{ on: modelValue.IsUwp }"
              @click="update({ IsUwp: !modelValue.IsUwp })"
            ></button>
          </div>
        </div>

        <div class="field-label" style="padding-top: 16px">Profiles</div>
        <div>
          <div v-if="!modelValue.AlternateLaunches?.length" class="alt-empty">
            None — <button class="link-btn" @click.prevent="addAlternate">add one</button>
          </div>
          <div v-for="(alt, idx) in modelValue.AlternateLaunches || []" :key="idx" class="alt-row">
            <input
              class="s-input"
              style="width: 130px; flex-shrink: 0"
              :value="alt.Name"
              placeholder="Profile name"
              @input="setAltName(idx, $event)"
            />
            <input
              class="s-input"
              style="flex: 1; font-family: var(--font-mono)"
              :value="alt.LaunchArgs"
              placeholder="--profile-directory=Default %URL%"
              @input="setAltArgs(idx, $event)"
            />
            <button class="icon-btn danger" @click="removeAlternate(idx)">
              <svg
                viewBox="0 0 24 24"
                width="12"
                height="12"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
              >
                <path d="M18 6 6 18M6 6l12 12" />
              </svg>
            </button>
          </div>
          <button
            v-if="modelValue.AlternateLaunches?.length"
            class="link-btn"
            @click="addAlternate"
          >
            + Add profile
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* ── Row shell ────────────────────────────────────────────────────────────── */
.row {
  border-bottom: 1px solid var(--hairline);
  transition: background 100ms;
}
.row:last-child {
  border-bottom: 0;
}

/* ── Row head ─────────────────────────────────────────────────────────────── */
.row-head {
  display: grid;
  grid-template-columns: 32px 36px 1fr auto auto;
  align-items: center;
  gap: 14px;
  padding: 12px 14px;
  cursor: pointer;
  transition: background 100ms;
}
.row:not(.is-open) .row-head:hover {
  background: var(--ink-10);
}

.chev {
  width: 14px;
  height: 14px;
  color: var(--ink-40);
  transition: transform 180ms ease;
  flex-shrink: 0;
}
.row.is-open .chev {
  transform: rotate(90deg);
  color: var(--accent);
}

.row-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: transparent;
  display: grid;
  place-items: center;
  flex-shrink: 0;
  overflow: visible;
  position: relative;
}

.row-meta {
  min-width: 0;
}
.row-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--ink-100);
  letter-spacing: -0.005em;
}
.row-path {
  margin-top: 2px;
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--ink-40);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.row-actions {
  display: flex;
  gap: 2px;
}

/* ── Row body ─────────────────────────────────────────────────────────────── */
.row-body {
  padding: 0 14px 16px 96px;
  border-top: 1px solid var(--hairline);
  background: rgba(0, 0, 0, 0.18);
}

.field-grid {
  display: grid;
  grid-template-columns: 140px 1fr;
  gap: 14px 18px;
  padding: 16px 0 4px;
}

.field-label {
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--ink-60);
  padding-top: 8px;
  align-self: start;
}

/* Alternate launches */
.alt-empty {
  font-size: 11px;
  color: var(--ink-40);
  padding: 6px 0;
}
.alt-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}
</style>
