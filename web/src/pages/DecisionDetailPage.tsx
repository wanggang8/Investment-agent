import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { DecisionTrace } from '../components/decision/DecisionTrace'
import { APIClientError } from '../services/client'
import { getDecision, createConfirmation, consultDecision } from '../services/decision'
import type { ConfirmationRequest, DecisionDetailResponse } from '../types/decision'

export function DecisionDetailPage() {
  const { decisionId } = useParams()
  const isConsultation = !decisionId
  const [decision, setDecision] = useState<DecisionDetailResponse>()
  const [message, setMessage] = useState(decisionId ? '选择一条建议后展示完整裁决链路。' : '')
  const [question, setQuestion] = useState('')
  const [symbol, setSymbol] = useState('')
  const [scenario, setScenario] = useState('hold_review')
  const [previousBaseMidpointPercent, setPreviousBaseMidpointPercent] = useState('')
  const [targetReturnPercent, setTargetReturnPercent] = useState('')

  useEffect(() => {
    if (!decisionId) return
    getDecision(decisionId)
      .then((res) => {
        setDecision(res.data)
        setMessage('')
      })
      .catch((error: unknown) => {
        setMessage(error instanceof APIClientError ? error.message : '暂时无法读取决策详情。')
      })
  }, [decisionId])

  function handleConsult() {
    if (!question.trim() || !symbol.trim()) {
      setMessage('请填写咨询问题和标的代码。')
      return
    }
    consultDecision({
      question: question.trim(),
      symbol: symbol.trim(),
      scenario,
      expected_return_previous_base_midpoint: percentInputToRate(previousBaseMidpointPercent),
      expected_return_target_return_rate: percentInputToRate(targetReturnPercent),
    })
      .then((res) => {
        setDecision(res.data)
        setMessage('')
      })
      .catch((error: unknown) => {
        setMessage(error instanceof APIClientError ? error.message : '暂时无法提交咨询。')
      })
  }

  function handleConfirm(id: string, payload: ConfirmationRequest) {
    createConfirmation(id, payload)
      .then(() => {
        setMessage('确认已记录。')
        return getDecision(id)
      })
      .then((res) => setDecision(res.data))
      .catch((error: unknown) => {
        setMessage(error instanceof APIClientError ? error.message : '暂时无法提交确认。')
      })
  }

  return (
    <div>
      <h1 className="page-title">{isConsultation ? '主动咨询' : '决策详情'}</h1>
      {isConsultation && (
        <article className="cockpit-card">
          <div className="state-label">咨询输入</div>
          <h2>输入假设</h2>
          <p>填写标的、场景和问题后，系统只生成本地分析材料和规则裁决记录；不会自动交易、自动确认或自动应用规则。</p>
          <div className="form-grid">
            <label>标的代码<input value={symbol} onChange={(event) => setSymbol(event.target.value)} /></label>
            <label>咨询场景
              <select value={scenario} onChange={(event) => setScenario(event.target.value)}>
                <option value="hold_review">持有评估</option>
                <option value="buy_review">买入评估</option>
                <option value="sell_review">卖出评估</option>
                <option value="rebalance_review">再平衡评估</option>
              </select>
            </label>
            <label>上一轮基准情景中枢（%）<input inputMode="decimal" value={previousBaseMidpointPercent} onChange={(event) => setPreviousBaseMidpointPercent(event.target.value)} placeholder="可选" /></label>
            <label>目标收益率（%）<input inputMode="decimal" value={targetReturnPercent} onChange={(event) => setTargetReturnPercent(event.target.value)} placeholder="可选" /></label>
          </div>
          <label>咨询问题<textarea value={question} onChange={(event) => setQuestion(event.target.value)} /></label>
          <div className="action-row">
            <button type="button" onClick={handleConsult}>提交咨询</button>
          </div>
        </article>
      )}
      {message && <div className="page-placeholder">{message}</div>}
      {isConsultation && decision && (
        <article className="cockpit-card">
          <div className="state-label">解释路径</div>
          <h2>生成结果可追踪</h2>
          <p>已生成本地决策材料。请先阅读裁决、安全边界、证据和闭环状态，再决定是否线下记录人工动作。</p>
          <div className="link-row">
            <Link to={`/decisions/${decision.decision_id}`}>打开生成的决策详情</Link>
            <Link to="/evidence">查看证据</Link>
            <Link to="/decision-loop">查看决策闭环</Link>
            <Link to="/audit">查看审计</Link>
          </div>
        </article>
      )}
      {decision && <DecisionTrace decision={decision} onConfirm={handleConfirm} />}
    </div>
  )
}

function percentInputToRate(value: string) {
  const trimmed = value.trim()
  if (!trimmed) return undefined
  const parsed = Number(trimmed)
  return Number.isFinite(parsed) ? parsed / 100 : undefined
}
