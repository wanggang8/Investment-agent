export type UITone = 'success' | 'warning' | 'danger' | 'degraded' | 'unknown' | 'readonly' | 'blocked'

export function joinClassNames(...names: Array<string | undefined | false>) {
  return names.filter(Boolean).join(' ')
}
