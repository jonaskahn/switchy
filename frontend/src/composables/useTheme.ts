import defaultThemeUrl from '@/themes/theme-default.css?url'
import acrylicAppUrl from '@/themes/appearance-acrylic.css?url'
import micaAppUrl from '@/themes/appearance-mica.css?url'

export type Appearance = 'acrylic' | 'mica'
export type Theme = 'default'

const THEME_URLS: Record<string, string> = {
  default: defaultThemeUrl,
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
