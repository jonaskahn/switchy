export function normLayout(raw?: string): 'vertical' | 'horizontal' {
  return String(raw || '')
    .trim()
    .toLowerCase() === 'horizontal'
    ? 'horizontal'
    : 'vertical'
}
