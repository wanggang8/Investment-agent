import { APIClientError } from '../../services/client'
import type { DisplayState } from '../../types/api'

export interface PageErrorState {
  state: DisplayState
  message: string
  code?: string
}

export function toPageErrorState(error: unknown, fallbackMessage = '暂时无法连接本地 API。'): PageErrorState {
  if (error instanceof APIClientError) {
    return {
      state: error.state.displayState,
      message: error.message,
      code: error.code,
    }
  }
  return { state: 'generic_failure', message: fallbackMessage }
}

export function formatCurrency(value?: number) {
  return new Intl.NumberFormat('zh-CN', {
    style: 'currency',
    currency: 'CNY',
    minimumFractionDigits: 2,
  }).format(value ?? 0)
}

export function formatPercent(value?: number) {
  return `${((value ?? 0) * 100).toFixed(2)}%`
}

export function formatDateTime(value?: string) {
  if (!value) return '暂无'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString('zh-CN', { hour12: false })
}

export { redactSensitiveText, type RedactionLabels } from './redaction'
