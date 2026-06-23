import type { DisplayState } from '../../types/api'

export interface StatusCopy {
  title: string
  message: string
}

const stateCopy: Record<DisplayState, StatusCopy> = {
  first_use: { title: '等待初始化', message: '请先录入账户、持仓和基础规则。' },
  normal: { title: '状态正常', message: '当前数据可用于展示纪律建议。' },
  insufficient_data: { title: '信息不足', message: '当前缺少证据、市场或账户数据，暂停交易类建议。' },
  frozen_watch: { title: '冻结观察', message: '证据核验或规则条件未满足，建议继续观察。' },
  high_risk: { title: '高危状态', message: '风险指标处于高位，禁止新增买入。' },
  data_source_unavailable: { title: '数据源不可用', message: '行情、情报或索引暂不可用，请刷新或稍后重试。' },
  generic_failure: { title: '读取失败', message: '页面暂时无法读取本地 API 响应。' },
}

const codeCopy: Record<string, StatusCopy> = {
  DATA_STALE: { title: '数据已过期', message: '本地数据已过期，请刷新后再查看。' },
  ANALYST_UNAVAILABLE: { title: '分析服务降级', message: '分析服务暂不可用，页面仅展示规则与已有数据。' },
  VECTOR_INDEX_UNAVAILABLE: { title: '索引不可用', message: '索引暂不可用，请稍后重试或重建索引。' },
}

export function getStatusNoticeCopy(state: DisplayState, code?: string): StatusCopy {
  return (code && codeCopy[code]) || stateCopy[state]
}
