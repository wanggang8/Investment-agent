import { useState } from 'react'
import type { ConfirmationRequest, ConfirmationType } from '../../types/decision'

const actionText: Record<ConfirmationType, string> = {
  planned: '记录计划',
  executed_manually: '已手动执行',
  watch: '标记待观察',
  marked_error: '标记错误',
}

const statusText: Record<string, string> = {
  pending: '待确认',
  planned: '已记录计划',
  executed_manually: '已手动执行',
  watch: '待观察',
  marked_error: '已标记错误',
  not_required: '无需确认',
}

const rootCauseOptions = [
  { value: 'evidence_missed', label: '证据遗漏' },
  { value: 'rule_threshold_issue', label: '规则阈值问题' },
  { value: 'analyst_error', label: '分析材料错误' },
  { value: 'user_context_missing', label: '用户上下文缺失' },
  { value: 'market_exception', label: '市场异常' },
]

interface Props {
  availableActions?: string[]
  confirmationStatus: string
  decisionId?: string
  onSubmit?: (decisionId: string, payload: ConfirmationRequest) => void
}

export function UserConfirmationPanel({ availableActions = [], confirmationStatus, decisionId = '', onSubmit }: Props) {
  const [selected, setSelected] = useState<ConfirmationType | null>(null)
  const [operationType, setOperationType] = useState<'buy' | 'sell' | 'reduce'>('buy')
  const [symbol, setSymbol] = useState('')
  const [quantity, setQuantity] = useState('')
  const [price, setPrice] = useState('')
  const [fees, setFees] = useState('')
  const [executedAt, setExecutedAt] = useState('')
  const [note, setNote] = useState('')
  const [actualOutcome, setActualOutcome] = useState('')
  const [rootCauseTag, setRootCauseTag] = useState('')
  const [lessonLearned, setLessonLearned] = useState('')
  const [error, setError] = useState('')

  const safeActions = availableActions.filter((action): action is ConfirmationType =>
    ['planned', 'executed_manually', 'watch', 'marked_error'].includes(action),
  )

  function submit() {
    if (!selected || !onSubmit || !decisionId) return
    const payload: ConfirmationRequest = { confirmation_type: selected }
    if (selected === 'executed_manually') {
      const quantityValue = Number(quantity)
      const priceValue = Number(price)
      const feesValue = fees ? Number(fees) : undefined
      if (!symbol || !operationType || !Number.isFinite(quantityValue) || quantityValue <= 0 || !Number.isFinite(priceValue) || priceValue <= 0 || !executedAt) {
        setError('请完整填写标的、动作、数量、价格和执行时间。')
        return
      }
      if (new Date(executedAt).getTime() > Date.now()) {
        setError('执行时间不能晚于当前时间。')
        return
      }
      if (feesValue !== undefined && (!Number.isFinite(feesValue) || feesValue < 0)) {
        setError('费用必须是非负数字。')
        return
      }
      payload.symbol = symbol.trim()
      payload.operation_type = operationType
      payload.quantity = quantityValue
      payload.price = priceValue
      if (feesValue !== undefined) payload.fees = feesValue
      payload.executed_at = toRFC3339(executedAt)
    }
    if (selected === 'marked_error') {
      if (!actualOutcome || !rootCauseTag || !lessonLearned) {
        setError('请完整填写实际结果、原因标签和复盘记录。')
        return
      }
      payload.actual_outcome = actualOutcome
      payload.root_cause_tag = rootCauseTag
      payload.lesson_learned = lessonLearned
    }
    if (note) payload.note = note
    setError('')
    onSubmit(decisionId, payload)
  }

  return (
    <article className="cockpit-card confirmation-card">
      <div className="state-label">用户确认</div>
      <p>当前确认状态：{statusText[confirmationStatus] ?? '未知状态'}</p>
      <div className="action-row">
        {safeActions.length === 0 ? (
          <span>当前无需确认。</span>
        ) : (
          safeActions.map((action) => (
            <button key={action} type="button" onClick={() => setSelected(action)}>
              {actionText[action]}
            </button>
          ))
        )}
      </div>

      {selected && (
        <div className="confirmation-form" aria-label="确认表单">
          <strong>{actionText[selected]}</strong>
          {selected === 'executed_manually' && (
            <>
              <label>标的代码<input value={symbol} onChange={(event) => setSymbol(event.target.value)} /></label>
              <label>线下动作
                <select value={operationType} onChange={(event) => setOperationType(event.target.value as 'buy' | 'sell' | 'reduce')}>
                  <option value="buy">买入</option>
                  <option value="sell">卖出</option>
                  <option value="reduce">减仓</option>
                </select>
              </label>
              <label>数量<input type="number" value={quantity} onChange={(event) => setQuantity(event.target.value)} /></label>
              <label>价格<input type="number" value={price} onChange={(event) => setPrice(event.target.value)} /></label>
              <label>费用<input type="number" min="0" value={fees} onChange={(event) => setFees(event.target.value)} /></label>
              <label>执行时间<input type="datetime-local" value={executedAt} onChange={(event) => setExecutedAt(event.target.value)} /></label>
            </>
          )}
          {selected === 'marked_error' && (
            <>
              <label>实际结果<input value={actualOutcome} onChange={(event) => setActualOutcome(event.target.value)} /></label>
              <label>原因标签
                <select value={rootCauseTag} onChange={(event) => setRootCauseTag(event.target.value)}>
                  <option value="">请选择原因标签</option>
                  {rootCauseOptions.map((option) => <option key={option.value} value={option.value}>{option.label}</option>)}
                </select>
              </label>
              <label>复盘记录<textarea value={lessonLearned} onChange={(event) => setLessonLearned(event.target.value)} /></label>
            </>
          )}
          <label>备注<textarea value={note} onChange={(event) => setNote(event.target.value)} /></label>
          {error && <p className="error-text">{error}</p>}
          <button type="button" onClick={submit}>提交确认</button>
        </div>
      )}

      <small>系统只记录你的线下动作，不会替你买入或卖出。</small>
    </article>
  )
}

function toRFC3339(value: string) {
  if (!value) return value
  return value.endsWith('Z') ? value : `${value}:00Z`
}
