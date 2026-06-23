export type OpsTone = 'success' | 'warning' | 'danger' | 'unknown'

export interface OpsMetric {
  label: string
  value: string
  detail?: string
  tone?: OpsTone
}

export interface OpsAction {
  label: string
  detail: string
  href: string
}

export function countBy<T extends string>(items: T[]) {
  return items.reduce<Record<string, number>>((acc, item) => {
    acc[item] = (acc[item] ?? 0) + 1
    return acc
  }, {})
}

export function hasForbiddenCopy(text: string) {
  return /自动交易|一键交易|代下单|外部推送|短信|邮件|第三方推送|自动确认|自动规则应用|自动修复|覆盖真实库|收益承诺/.test(text)
}

