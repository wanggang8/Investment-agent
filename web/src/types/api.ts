// APIErrorBody 对应后端统一响应信封中的 error 字段。
export interface APIErrorBody {
  code: string
  message: string
  detail?: string
}

// APIMeta 承载规则版本、生成时间等非主体数据。
export interface APIMeta {
  generated_at?: string
  rule_version?: string
}

// APIEnvelope 是业务接口统一响应结构。
// P0 健康检查是例外，前端 API client 会转换成该结构。
export interface APIEnvelope<T = unknown> {
  request_id: string
  data?: T
  meta?: APIMeta
  error?: APIErrorBody
}

// PageResult 对应后端列表接口的本地单页返回结构。
export interface PageResult<T> {
  items: T[]
  total: number
}

export type DisplayState =
  | 'first_use'
  | 'normal'
  | 'insufficient_data'
  | 'frozen_watch'
  | 'high_risk'
  | 'data_source_unavailable'
  | 'generic_failure'

export interface APIErrorState {
  requestId: string
  code: string
  message: string
  detail?: string
  displayState: DisplayState
  httpStatus?: number
}

// HealthData 是健康检查返回数据。
export interface HealthData {
  status: string
}
