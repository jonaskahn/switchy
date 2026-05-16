<script setup lang="ts">
import { onMounted, computed, toRef } from 'vue'
import type { Browser } from '@/types'
import { useBrowserIcon } from '@/composables/useBrowserIcon'

const props = defineProps<{
  browser: Browser
  shortcut: string
  focused: boolean
  animationDelayMs: number
}>()

defineEmits<{ (e: 'open'): void }>()

const browserRef = toRef(props, 'browser')
const { iconDataUrl, loadIcon } = useBrowserIcon(browserRef, { pixelSize: 256 })

const tileClass = computed(() => ({
  tile: true,
  'is-focused': props.focused,
}))

const animationStyle = computed(() => ({ animationDelay: `${props.animationDelayMs}ms` }))

onMounted(() => {
  void loadIcon()
})
</script>

<template>
  <button
    type="button"
    :class="tileClass"
    role="option"
    :aria-selected="focused ? 'true' : 'false'"
    :style="animationStyle"
    @click="$emit('open')"
  >
    <span v-if="shortcut" class="keycap">{{ shortcut }}</span>
    <div class="tile-icon">
      <img
        v-if="iconDataUrl"
        class="tile-img"
        :src="iconDataUrl"
        :alt="browser.Name"
        decoding="async"
        draggable="false"
      />
      <svg
        v-else
        viewBox="0 0 24 24"
        width="40"
        height="40"
        fill="none"
        stroke="currentColor"
        stroke-width="1.35"
        stroke-linecap="round"
        stroke-linejoin="round"
        style="color: var(--ink-60)"
      >
        <circle cx="12" cy="12" r="10" />
        <path d="M12 2a14.5 14.5 0 0 0 0 20 14.5 14.5 0 0 0 0-20" />
        <path d="M2 12h20" />
      </svg>
      <img
        v-if="browser.ProfileIcon"
        class="profile-avatar tile-profile"
        :src="browser.ProfileIcon"
        decoding="async"
        alt="Profile"
      />
    </div>
    <div class="tile-name">{{ browser.Name }}</div>
  </button>
</template>

<style scoped>
.tile {
  position: relative;
  flex: 0 0 110px;
  padding: 14px 10px 12px;
  border-radius: 14px;
  cursor: pointer;
  border: 1px solid transparent;
  background: transparent;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  transition:
    background 140ms ease,
    border-color 140ms ease,
    transform 180ms cubic-bezier(0.2, 0.8, 0.2, 1);
  animation: tileIn 320ms cubic-bezier(0.2, 0.8, 0.2, 1) both;
  color: inherit;
  font: inherit;
}

@keyframes tileIn {
  from {
    opacity: 0;
    transform: translateY(6px) scale(0.96);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.tile:hover {
  background: var(--ink-10);
  border-color: var(--hairline-hi);
  transform: translateY(-2px);
}

.tile.is-focused {
  background: var(--accent-tint);
  border-color: rgb(from var(--accent) r g b / 0.45);
  transform: translateY(-2px);
}

.tile .keycap {
  position: absolute;
  top: -8px;
  right: 8px;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  border-radius: 5px;
  background: var(--ink-10);
  border: 1px solid var(--hairline);
  border-bottom-width: 2px;
  color: var(--ink-80);
  font-family: var(--font-mono);
  font-size: 10px;
  font-weight: 500;
  display: grid;
  place-items: center;
  line-height: 1;
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
}

.tile.is-focused .keycap {
  background: var(--accent);
  border-color: var(--accent-dim);
  color: var(--accent-text);
}

.tile-icon {
  position: relative;
  width: 56px;
  height: 56px;
  border-radius: 0;
  display: grid;
  place-items: center;
  background: transparent;
  overflow: visible;
}

.tile-profile {
  width: 18px;
  height: 18px;
  bottom: -4px;
  right: -4px;
}

.tile-img {
  width: 56px;
  height: 56px;
  object-fit: contain;
  display: block;
  pointer-events: none;
  image-rendering: -webkit-optimize-contrast;
}

.tile-name {
  font-size: 11.5px;
  font-weight: 500;
  color: var(--ink-100);
  letter-spacing: -0.005em;
  line-height: 1.2;
  text-align: center;
  white-space: nowrap;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
