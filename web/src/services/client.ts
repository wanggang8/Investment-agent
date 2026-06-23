import type { APIEnvelope, APIErrorState, DisplayState } from '../types/api'

const defaultBase = ''

// 这些错误码是前端状态的稳定依据，页面不得解析后端底层错误文本。
const errorDisplayState: Record<string, DisplayState> = {
  DATA_REQUIRED: 'first_use',
  DATA_STALE: 'insufficient_data',
  EVIDENCE_NOT_FOUND: 'insufficient_data',
  VECTOR_INDEX_UNAVAILABLE: 'insufficient_data',
  ANALYST_UNAVAILABLE: 'insufficient_data',
  DECISION_RECORD_FAILED: 'insufficient_data',
  SOURCE_VERIFICATION_FAILED: 'frozen_watch',
  DATA_SOURCE_UNAVAILABLE: 'data_source_unavailable',
  MARKET_SNAPSHOT_WRITE_FAILED: 'generic_failure',
  RULE_VERSION_MISSING: 'high_risk',
  CONFLICT: 'generic_failure',
  INVALID_STATE: 'frozen_watch',
  INTERNAL_ERROR: 'generic_failure',
}

const safeErrorMessages: Record<string, string> = {
  DATA_REQUIRED: '需要先录入账户与持仓数据。',
  DATA_STALE: '本地数据已过期，请刷新后再查看。',
  EVIDENCE_NOT_FOUND: '证据不足，当前暂停交易类建议。',
  VECTOR_INDEX_UNAVAILABLE: '索引暂不可用，请稍后重试或重建索引。',
  ANALYST_UNAVAILABLE: '分析服务暂不可用，页面仅展示规则与已有数据。',
  DECISION_RECORD_FAILED: '决策记录暂不可用，请稍后重试。',
  SOURCE_VERIFICATION_FAILED: '证据核验未完成，当前进入观察状态。',
  DATA_SOURCE_UNAVAILABLE: '数据源暂不可用，请检查数据源状态。',
  MARKET_SNAPSHOT_WRITE_FAILED: '市场快照更新失败，请稍后重试。',
  RULE_VERSION_MISSING: '当前缺少可用规则版本。',
  CONFLICT: '当前操作存在状态冲突，请刷新后重试。',
  INVALID_STATE: '当前状态不允许执行该操作。',
  INTERNAL_ERROR: '系统暂时无法处理请求，请稍后重试。',
  HTTP_ERROR: '本地 API 返回了无法识别的响应。',
  UNKNOWN: '系统暂时无法处理请求，请稍后重试。',
}

export class APIClientError extends Error {
  readonly state: APIErrorState

  constructor(state: APIErrorState) {
    super(state.message)
    this.name = 'APIClientError'
    this.state = state
  }

  get requestId() {
    return this.state.requestId
  }

  get code() {
    return this.state.code
  }

  get detail() {
    return this.state.detail
  }
}

function generateRequestId(): string {
  const ts = new Date().toISOString().replace(/[-:TZ.]/g, '').slice(0, 14)
  const rand = Math.random().toString(36).slice(2, 8)
  return `req_${ts}_${rand}`
}

export function mapErrorCodeToDisplayState(code: string, status?: number): DisplayState {
  if (errorDisplayState[code]) {
    return errorDisplayState[code]
  }
  if (status === 503) {
    return 'data_source_unavailable'
  }
  if (status === 409) {
    return 'insufficient_data'
  }
  return 'generic_failure'
}

export function getSafeErrorMessage(code: string, status?: number): string {
  if (safeErrorMessages[code]) {
    return safeErrorMessages[code]
  }
  if (status === 503) {
    return safeErrorMessages.DATA_SOURCE_UNAVAILABLE
  }
  if (status === 409) {
    return '当前信息不足或状态冲突，请刷新后重试。'
  }
  return safeErrorMessages.UNKNOWN
}

function makeClientError(
  requestId: string,
  code: string,
  _message: string,
  detail?: string,
  httpStatus?: number,
) {
  const safeMessage = getSafeErrorMessage(code, httpStatus)
  return new APIClientError({
    requestId,
    code,
    message: safeMessage,
    // detail 仅保留在错误对象里方便开发排查，页面层不得展示该字段。
    detail,
    httpStatus,
    displayState: mapErrorCodeToDisplayState(code, httpStatus),
  })
}

export async function apiRequest<T>(
  path: string,
  init?: RequestInit,
): Promise<APIEnvelope<T>> {
  const requestId = generateRequestId()
  const headers = new Headers(init?.headers)
  headers.set('Accept', 'application/json')
  headers.set('X-Request-Id', requestId)

  if (init?.body && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json')
  }

  const res = await fetch(`${defaultBase}${path}`, { ...init, headers })
  const contentType = res.headers.get('content-type') ?? ''

  if (!contentType.includes('application/json')) {
    throw makeClientError(requestId, 'HTTP_ERROR', `HTTP ${res.status}`, undefined, res.status)
  }

  const body = (await res.json()) as APIEnvelope<T> & { status?: string }

  // P0 健康检查返回普通 JSON，这里转换为统一信封，避免页面写特例。
  if (path.includes('/health') && body && 'status' in body && !('request_id' in body)) {
    const healthBody = body as { status: string }
    return {
      request_id: requestId,
      data: { status: healthBody.status } as T,
    }
  }

  const envelope = body as APIEnvelope<T>
  const rid = envelope.request_id || requestId

  if (!res.ok || envelope.error) {
    throw makeClientError(
      rid,
      envelope.error?.code ?? 'UNKNOWN',
      envelope.error?.message ?? `HTTP ${res.status}`,
      envelope.error?.detail,
      res.status,
    )
  }

  return { ...envelope, request_id: rid }
}

export function getHealth() {
  return apiRequest<{ status: string }>('/api/v1/health')
}
