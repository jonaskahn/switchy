import defaultThemeUrl from '@/themes/theme-default.css?url'
import midnightThemeUrl from '@/themes/theme-midnight.css?url'
import forestThemeUrl from '@/themes/theme-forest.css?url'
import crimsonThemeUrl from '@/themes/theme-crimson.css?url'
import auroraThemeUrl from '@/themes/theme-aurora.css?url'
import glacierThemeUrl from '@/themes/theme-glacier.css?url'
import emberThemeUrl from '@/themes/theme-ember.css?url'
import sakuraThemeUrl from '@/themes/theme-sakura.css?url'
import monochromeThemeUrl from '@/themes/theme-monochrome.css?url'
import copperThemeUrl from '@/themes/theme-copper.css?url'
import voidThemeUrl from '@/themes/theme-void.css?url'
import oceanThemeUrl from '@/themes/theme-ocean.css?url'
import acrylicAppUrl from '@/themes/appearance-acrylic.css?url'
import micaAppUrl from '@/themes/appearance-mica.css?url'

export type Appearance = 'acrylic' | 'mica'
export type Theme =
  | 'default'
  | 'midnight'
  | 'forest'
  | 'crimson'
  | 'aurora'
  | 'glacier'
  | 'ember'
  | 'sakura'
  | 'monochrome'
  | 'copper'
  | 'void'
  | 'ocean'

const THEME_URLS: Record<string, string> = {
  default: defaultThemeUrl,
  midnight: midnightThemeUrl,
  forest: forestThemeUrl,
  crimson: crimsonThemeUrl,
  aurora: auroraThemeUrl,
  glacier: glacierThemeUrl,
  ember: emberThemeUrl,
  sakura: sakuraThemeUrl,
  monochrome: monochromeThemeUrl,
  copper: copperThemeUrl,
  void: voidThemeUrl,
  ocean: oceanThemeUrl,
}

const APPEARANCE_URLS: Record<string, string> = {
  acrylic: acrylicAppUrl,
  mica: micaAppUrl,
}

function swapLink(id: string, href: string): Promise<void> {
  return new Promise((resolve) => {
    let link = document.getElementById(id) as HTMLLinkElement | null
    if (!link) {
      link = document.createElement('link')
      link.id = id
      link.rel = 'stylesheet'
      document.head.appendChild(link)
    }
    // Already loaded with this exact href — resolve immediately
    if (link.getAttribute('data-href') === href && link.sheet) {
      resolve()
      return
    }
    link.setAttribute('data-href', href)
    link.addEventListener('load', () => resolve(), { once: true })
    link.addEventListener('error', () => resolve(), { once: true })
    link.href = href
  })
}

/** Loads the color palette theme stylesheet. */
export function applyTheme(theme: string): Promise<void> {
  return swapLink('switchy-theme', THEME_URLS[theme] ?? THEME_URLS.default)
}

/** Loads the window material appearance stylesheet.
 *  Always appended after the theme link so appearance can
 *  fine-tune low-opacity surface vars (--ink-10, --ink-20, etc.). */
export function applyAppearance(appearance: string): Promise<void> {
  return swapLink('switchy-appearance', APPEARANCE_URLS[appearance] ?? APPEARANCE_URLS.acrylic)
}
