import { describe, expect, it } from 'vitest'
import { redactSensitiveText } from './redaction'

describe('redactSensitiveText', () => {
  it('redacts key-shaped tokens, SQL, prompts, raw diagnostics, stack traces, and local paths', () => {
    const input = '请求失败 sk-1234567890abcdef /Users/private/db.sqlite stack trace SELECT * FROM secrets prompt: raw payload raw vendor payload /tmp/raw.log C:\\Users\\vick\\db.sqlite'

    const output = redactSensitiveText(input)

    expect(output).toContain('已脱敏密钥')
    expect(output).toContain('SQL 已脱敏')
    expect(output).toContain('prompt 已脱敏')
    expect(output).toContain('raw 诊断已脱敏')
    expect(output).toContain('已脱敏错误摘要')
    expect(output).toContain('已脱敏路径')
    expect(output).not.toMatch(/sk-1234567890abcdef|\/Users\/private|\/tmp\/raw\.log|C:\\Users|SELECT \* FROM|prompt:|raw vendor payload|stack trace/)
  })

  it('allows caller-specific replacement labels', () => {
    const input = 'DELETE FROM accounts prompt = sk-secret raw stack /opt/private/file'

    const output = redactSensitiveText(input, {
      key: '[KEY]',
      sql: '[SQL]',
      prompt: '[PROMPT]',
      raw: '[RAW]',
      path: '<path>',
    })

    expect(output).toContain('[SQL]')
    expect(output).toContain('<path>')
    expect(output).not.toMatch(/DELETE FROM|sk-secret|prompt =|raw stack|\/opt\/private/)
  })
})
