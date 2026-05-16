<script setup lang="ts">
import { ref, onMounted } from 'vue'
import SelectorView from '@/views/SelectorView.vue'
import SettingsView from '@/views/SettingsView.vue'
import { applyTheme, applyAppearance } from '@/composables/useTheme'

const mode = ref<'selector' | 'settings' | null>(null)

onMounted(async () => {
  try {
    const { GetMode, GetSettings } = await import('@wailsjs/go/main/App')
    const [m, cfg] = await Promise.all([GetMode(), GetSettings()])
    // Wait for both stylesheets to finish loading before showing the view,
    // so no frame renders with undefined CSS variables.
    await Promise.all([
      applyTheme(cfg?.AppSettings?.Theme ?? 'default'),
      applyAppearance(cfg?.AppSettings?.Appearance ?? 'acrylic'),
    ])
    mode.value = m as 'selector' | 'settings'
  } catch {
    mode.value = 'selector'
  }
})
</script>

<template>
  <div class="app-root">
    <SelectorView v-if="mode === 'selector'" />
    <SettingsView v-else-if="mode === 'settings'" />
    <div v-else class="loading-panel">
      <span class="loading-text">Loading…</span>
    </div>
  </div>
</template>

<style scoped>
.app-root {
  height: 100%;
  width: 100%;
}
.loading-panel {
  height: 100%;
  width: 100%;
  display: grid;
  place-items: center;
}
.loading-text {
  font-family: var(--font-mono);
  font-size: 12.5px;
  color: var(--ink-40);
  letter-spacing: 0.04em;
}
</style>
