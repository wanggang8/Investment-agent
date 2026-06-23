export type RedactionLabels = {
  key?: string
  sql?: string
  prompt?: string
  raw?: string
  stack?: string
  path?: string
}

const defaultLabels: Required<RedactionLabels> = {
  key: '已脱敏密钥',
  sql: 'SQL 已脱敏',
  prompt: 'prompt 已脱敏',
  raw: 'raw 诊断已脱敏',
  stack: '已脱敏错误摘要',
  path: '已脱敏路径',
}

export function redactSensitiveText(value: string | null | undefined, labels: RedactionLabels = {}) {
  if (!value) {
    return value ?? ''
  }
  const merged = { ...defaultLabels, ...labels }
  return value
    .replace(/sk-[A-Za-z0-9_-]+/g, merged.key)
    .replace(/\bSELECT\s+\*\s+FROM\s+[A-Za-z0-9_.-]+/gi, merged.sql)
    .replace(/\b(INSERT\s+INTO|UPDATE|DELETE\s+FROM|DROP\s+TABLE|ALTER\s+TABLE|CREATE\s+TABLE)\s+[A-Za-z0-9_.-]+/gi, merged.sql)
    .replace(/\bprompt\s*[:=]\s*[^\s，。；;]*/gi, merged.prompt)
    .replace(/raw\s+(vendor|provider|http|stack|payload)/gi, merged.raw)
    .replace(/stack trace/gi, merged.stack)
    .replace(/(?:[A-Za-z]:\\Users\\|[A-Za-z]:\\|\/Users\/|\/tmp\/|\/opt\/private\/|\/home\/)[^\s"']+/g, merged.path)
}
