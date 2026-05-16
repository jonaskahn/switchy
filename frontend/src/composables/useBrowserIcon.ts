import { ref, type Ref, watch } from 'vue'
import { GetBrowserIconAt } from '@wailsjs/go/main/App'

const defaultPixelSize = 256

/** Loads a PNG data URL for a browser tile from backend icon extraction. */
export function useBrowserIcon(
  browserSource: Ref<{ IconSpec?: string; ExePath?: string } | undefined>,
  opts?: { pixelSize?: number }
) {
  const iconDataUrl = ref('')
  const pixelSize = Math.min(256, Math.max(64, opts?.pixelSize ?? defaultPixelSize))

  async function loadIcon() {
    const b = browserSource.value
    if (!b) {
      iconDataUrl.value = ''
      return
    }
    const spec = (b.IconSpec || b.ExePath || '').trim()
    if (!spec) {
      iconDataUrl.value = ''
      return
    }
    try {
      const data = await GetBrowserIconAt(spec, pixelSize)
      iconDataUrl.value = data || ''
    } catch {}
  }

  watch(
    () => [browserSource.value?.ExePath, browserSource.value?.IconSpec] as const,
    () => {
      void loadIcon()
    }
  )

  return { iconDataUrl, loadIcon, pixelSize }
}
