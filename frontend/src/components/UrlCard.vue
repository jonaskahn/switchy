<script setup lang="ts">
import { ref, computed, nextTick } from 'vue'
import { SetURL } from '@wailsjs/go/main/App'

const URL_QUERY_DISPLAY_LENGTH = 20
const COPY_TITLE_DEFAULT = 'Copy URL (C)'
const COPY_TITLE_DONE = 'Copied!'

const props = defineProps<{
  url: string
  copied: boolean
  compact?: boolean
}>()

const emit = defineEmits<{ copy: [] }>()

const root = ref<HTMLElement | null>(null)
const editInput = ref<HTMLInputElement | null>(null)
defineExpose({ root })

const editing = ref(false)
const editValue = ref('')

const urlKind = computed(() => {
  try {
    const { protocol } = new URL(props.url)
    if (protocol === 'mailto:') return 'email'
    if (protocol === 'https:') return 'https'
    if (protocol === 'http:') return 'http'
  } catch {}
  return 'http'
})

const displayUrl = computed(() => {
  try {
    const parsed = new URL(props.url)
    const path = parsed.pathname.length > 1 ? parsed.pathname : ''
    const query = truncateQuery(parsed.search)
    return parsed.hostname + path + query
  } catch {
    return props.url
  }
})

const copyTitle = computed(() => (props.copied ? COPY_TITLE_DONE : COPY_TITLE_DEFAULT))

function truncateQuery(search: string): string {
  if (!search) return ''
  const truncated = search.slice(0, URL_QUERY_DISPLAY_LENGTH)
  return search.length > URL_QUERY_DISPLAY_LENGTH ? truncated + '…' : truncated
}

function startEdit() {
  editValue.value = props.url
  editing.value = true
  nextTick(() => editInput.value?.focus())
}

function cancelEdit() {
  editing.value = false
  editValue.value = ''
}

async function saveEdit() {
  const newUrl = editValue.value.trim()
  if (newUrl && newUrl !== props.url) {
    try {
      new URL(newUrl)
      await SetURL(newUrl)
    } catch {}
  }
  editing.value = false
}

function onEditKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter') {
    saveEdit()
  } else if (e.key === 'Escape') {
    cancelEdit()
  }
}
</script>

<template>
  <div ref="root" :class="['url-card', { 'url-strip': compact }]">
    <div class="url-icon" :class="urlKind">
      <svg
        v-if="urlKind === 'email'"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="1.6"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <rect x="2" y="4" width="20" height="16" rx="2" />
        <path d="M22 4L12 13 2 4" />
      </svg>
      <svg v-else-if="urlKind === 'https'" viewBox="0 0 64 64" xmlns="http://www.w3.org/2000/svg">
        <path
          d="M32 5 L54 14 V28 C54 43 45 54 32 60 C19 54 10 43 10 28 V14 Z"
          fill="var(--accent)"
          stroke="var(--accent)"
          stroke-width="3"
          stroke-linejoin="round"
          opacity="0.25"
        />
        <path
          d="M32 5 L54 14 V28 C54 43 45 54 32 60 C19 54 10 43 10 28 V14 Z"
          fill="none"
          stroke="var(--accent)"
          stroke-width="3"
          stroke-linejoin="round"
        />
      </svg>
      <svg v-else-if="urlKind === 'http'" viewBox="0 0 64 64" xmlns="http://www.w3.org/2000/svg">
        <defs>
          <linearGradient id="redGradient" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stop-color="#f87171" />
            <stop offset="100%" stop-color="#b91c1c" />
          </linearGradient>
        </defs>
        <path
          d="M32 5 L54 14 V28 C54 43 45 54 32 60 C19 54 10 43 10 28 V14 Z"
          fill="url(#redGradient)"
          stroke="#7f1d1d"
          stroke-width="3"
          stroke-linejoin="round"
        />
        <path d="M32 10 L48 17 V28 C48 40 42 48 32 54 Z" fill="#ffffff" opacity="0.16" />
      </svg>
      <svg
        v-else
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="1.6"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71" />
        <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71" />
      </svg>
    </div>
    <input
      v-if="editing"
      ref="editInput"
      v-model="editValue"
      class="url-edit"
      @blur="saveEdit"
      @keydown="onEditKeydown"
    />
    <span v-else class="url-value" :class="{ empty: !displayUrl }" :title="url">
      {{ displayUrl || '(no URL)' }}
    </span>
    <button v-if="!editing" class="icon-btn no-drag" title="Edit URL" @click="startEdit">
      <svg
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
        fill="var(--accent)"
        stroke="var(--accent)"
        stroke-width="1"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <g transform="scale(0.9) translate(2 2)">
          <path d="M16.5 3.5a2.1 2.1 0 0 1 3 3L7 19l-4 1 1-4Z" />
        </g>
      </svg>
    </button>
    <template v-else>
      <button class="icon-btn no-drag" title="Save (Enter)" @click="saveEdit">
        <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" fill="none">
          <circle cx="12" cy="12" r="11" fill="var(--accent)" />
          <path
            d="M8 12.5l3 3 5-6"
            stroke="white"
            stroke-width="3"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
      </button>
      <button class="icon-btn no-drag" title="Cancel" @click="cancelEdit">
        <svg
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="3.5"
          stroke-linecap="round"
        >
          <path d="M6 6 L18 18 M18 6 L6 18" />
        </svg>
      </button>
    </template>
    <button
      v-if="!editing"
      class="copy-btn no-drag"
      :title="copyTitle"
      :disabled="copied"
      :aria-disabled="copied"
      @click="emit('copy')"
    >
      <svg v-if="copied" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" fill="none">
        <circle cx="12" cy="12" r="11" fill="var(--accent)" />
        <path
          d="M8 12.5l3 3 5-6"
          stroke="white"
          stroke-width="3"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
      </svg>
      <svg
        v-else
        viewBox="0 0 24 24"
        fill="var(--accent)"
        stroke="var(--accent)"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <rect x="9" y="9" width="11" height="11" rx="2" />
        <rect x="4" y="4" width="11" height="11" rx="2" />
      </svg>
    </button>
  </div>
</template>

<style scoped>
.url-card {
  margin: 4px 10px 10px;
  padding: 8px 14px;
  background: var(--ink-10);
  border: 1px solid var(--hairline);
  border-radius: var(--r-row);
  display: grid;
  grid-template-columns: auto 1fr auto auto;
  gap: 6px;
  align-items: center;
  flex-shrink: 0;
}
.url-strip {
  margin: 16px 28px 52px;
}
.url-icon {
  width: 28px;
  height: 28px;
  display: grid;
  place-items: center;
  flex-shrink: 0;
}
.url-icon svg {
  width: 20px;
  height: 20px;
}
.url-icon.link svg {
  color: var(--ink-60);
}
.url-icon.http svg {
  color: #ef4444;
}
.url-icon.email svg {
  color: var(--accent);
}
.url-icon.link svg {
  color: var(--ink-60);
}
.url-value {
  font-family: var(--font-mono);
  font-size: 12.5px;
  color: var(--ink-100);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  letter-spacing: -0.005em;
  text-align: center;
}
.url-strip .url-value {
  font-size: 12px;
}
.url-value.empty {
  color: var(--ink-40);
  font-style: italic;
  font-family: var(--font-ui);
  font-size: 12px;
}
.url-edit {
  font-family: var(--font-mono);
  font-size: 12.5px;
  color: var(--ink-100);
  background: var(--substrate);
  border: 1px solid var(--accent);
  border-radius: 4px;
  padding: 2px 4px;
  outline: none;
  width: 100%;
}
.icon-btn {
  width: 28px;
  height: 28px;
  border: 0;
  background: transparent;
  border-radius: 8px;
  color: var(--ink-60);
  display: grid;
  place-items: center;
  cursor: pointer;
  transition:
    background 120ms,
    color 120ms;
}
.icon-btn:hover {
  background: var(--ink-10);
  color: var(--ink-100);
}
.icon-btn svg {
  width: 20px;
  height: 20px;
}
.copy-btn {
  width: 28px;
  height: 28px;
  border: 0;
  background: transparent;
  border-radius: 8px;
  color: var(--ink-60);
  display: grid;
  place-items: center;
  cursor: pointer;
  transition:
    background 120ms,
    color 120ms;
}
.copy-btn:hover {
  background: var(--ink-10);
  color: var(--ink-100);
}
.copy-btn svg {
  width: 20px;
  height: 20px;
}
</style>
