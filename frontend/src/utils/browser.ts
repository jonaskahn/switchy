export const MAX_KEYBOARD_SHORTCUT_INDEX = 9
export const ANIMATION_BASE_DELAY_MS = 40
export const ANIMATION_STEP_DELAY_MS = 40

export function browserShortcut(idx: number): string {
  return idx < MAX_KEYBOARD_SHORTCUT_INDEX ? String(idx + 1) : ''
}

export function browserAnimationDelay(idx: number): number {
  return ANIMATION_BASE_DELAY_MS + idx * ANIMATION_STEP_DELAY_MS
}

export function browserAnimationStyle(idx: number): Record<string, string> {
  return { animationDelay: `${browserAnimationDelay(idx)}ms` }
}
