import { useMemo, useState } from 'react'
import { StatusNotice } from '../components/status/StatusNotice'
import { Button, Field, SummaryCard, type UITone } from '../components/ui'
import { buildLocalOpsModel, localOpsMetricTitle } from '../features/governance'
import {
  confirmLocalKnowledgeImport,
  validateLocalKnowledgeImport,
} from '../services/localKnowledge'
import type { PageErrorState } from '../shared/utils'
import { toPageErrorState } from '../shared/utils'
import type {
  LocalKnowledgeImportConfirmResponse,
  LocalKnowledgeImportRow,
  LocalKnowledgeImportValidationResponse,
} from '../types/localKnowledge'

const sampleRowsText = JSON.stringify([
  {
    title: '510300 估值观察',
    text: '本地研究记录：指数估值处于偏高区间，后续仅作为背景材料参与检索。',
    symbol: '510300',
    tags: ['估值', '本地研究'],
  },
], null, 2)

type FormState = {
  sourceLabel: string
  defaultSymbol: string
  rowsText: string
  confirmReason: string
}

const initialForm: FormState = {
  sourceLabel: 'local_research_notes',
  defaultSymbol: '510300',
  rowsText: sampleRowsText,
  confirmReason: '人工确认导入为本地背景材料',
}

export function LocalKnowledgePage() {
  const [form, setForm] = useState<FormState>(initialForm)
  const [parseError, setParseError] = useState('')
  const [apiError, setApiError] = useState<PageErrorState>()
  const [validation, setValidation] = useState<LocalKnowledgeImportValidationResponse>()
  const [confirmResult, setConfirmResult] = useState<LocalKnowledgeImportConfirmResponse>()
  const [isValidating, setIsValidating] = useState(false)
  const [isConfirming, setIsConfirming] = useState(false)

  const parsedRows = useMemo(() => parseRows(form.rowsText), [form.rowsText])
  const canConfirm = Boolean(validation && validation.summary.blocking_count === 0 && !parsedRows.error)
  const localModel = buildLocalOpsModel({ validation })

  function updateField(field: keyof FormState, value: string) {
    setForm((next) => ({ ...next, [field]: value }))
  }

  async function handleValidate() {
    const nextRows = parseRows(form.rowsText)
    setParseError(nextRows.error ?? '')
    setApiError(undefined)
    setConfirmResult(undefined)
    if (nextRows.error) {
      setValidation(undefined)
      return
    }

    setIsValidating(true)
    try {
      const response = await validateLocalKnowledgeImport({
        source_label: form.sourceLabel,
        default_symbol: form.defaultSymbol,
        rows: nextRows.rows,
      })
      setValidation(response.data)
    } catch (error: unknown) {
      setValidation(undefined)
      setApiError(toPageErrorState(error))
    } finally {
      setIsValidating(false)
    }
  }

  async function handleConfirm() {
    const nextRows = parseRows(form.rowsText)
    setParseError(nextRows.error ?? '')
    setApiError(undefined)
    setConfirmResult(undefined)
    if (!validation || nextRows.error) {
      return
    }

    setIsConfirming(true)
    try {
      const response = await confirmLocalKnowledgeImport({
        import_batch_id: validation.import_batch_id,
        confirm_reason: form.confirmReason,
        source_label: form.sourceLabel,
        default_symbol: form.defaultSymbol,
        rows: nextRows.rows,
      })
      setConfirmResult(response.data)
    } catch (error: unknown) {
      setApiError(toPageErrorState(error))
    } finally {
      setIsConfirming(false)
    }
  }

  return (
    <div className="reference-tight-page">
      <h1 className="page-title">本地知识导入</h1>

      <section className={`daily-hero daily-tone-${localModel.overallTone}`} aria-label="本地知识导入总览">
        <div className="daily-hero-main">
          <div className="state-label">本地配置与诊断状态</div>
          <h2>{localModel.overallLabel}</h2>
          <p>{localModel.safetyNotes[0]}</p>
          <div className="daily-signal-grid quality-signal-grid">
            {localModel.metrics.map((metric) => (
              <SummaryCard key={metric.label} title={localOpsMetricTitle(metric.label)} value={metric.value} detail={metric.detail} tone={(metric.tone ?? 'unknown') as UITone} />
            ))}
          </div>
        </div>
        <aside className="daily-hero-side" aria-label="本地知识导入下一步">
          <strong>下一步人工复核</strong>
          <ul>
            {localModel.nextActions.map((action) => (
              <li key={action.label}>
                <a href={action.href} aria-label={`${action.label}入口`}>{action.label}</a>
                <span>{action.detail}</span>
              </li>
            ))}
          </ul>
        </aside>
      </section>
      <p className="reference-page-note">
        将本地研究记录先校验为脱敏预览，再由人工写入本地背景事实。页面只展示预览、批次和索引计划，不触发交易、不改变规则、不发送到外部通道。
      </p>

      <section className="cockpit-grid" aria-label="本地知识导入区域">
        <article className="cockpit-card form-card">
          <div className="state-label">导入草稿</div>
          <h2>来源与内容</h2>
          <p>先用来源、标的和记录数量建立导入摘要；结构化记录可在高级详情中编辑。</p>
          <div className="product-summary-grid" aria-label="本地知识草稿摘要">
            <div className="product-summary-card"><span>来源标签</span><strong>{sourceLabelText(form.sourceLabel)}</strong></div>
            <div className="product-summary-card"><span>默认标的</span><strong>{form.defaultSymbol || '待填写'}</strong></div>
            <div className="product-summary-card"><span>记录数量</span><strong>{parsedRows.error ? '需修正' : `${parsedRows.rows.length} 条`}</strong></div>
          </div>
          <details className="product-detail">
            <summary>编辑结构化记录</summary>
            <div className="product-detail-body">
              <div className="form-grid form-grid-compact">
                <Field id="local-knowledge-source-label" label="来源标签" hint="只作为本地背景材料来源标记。">
                  <input value={form.sourceLabel} onChange={(event) => updateField('sourceLabel', event.target.value)} />
                </Field>
                <Field id="local-knowledge-default-symbol" label="默认标的">
                  <input value={form.defaultSymbol} onChange={(event) => updateField('defaultSymbol', event.target.value)} />
                </Field>
              </div>
              <Field id="local-knowledge-rows-json" label="记录 JSON" hint="字段支持 title、text、symbol、as_of_date、tags。提交前会先做本地预览校验。" error={parseError || undefined}>
                <textarea
                  value={form.rowsText}
                  rows={12}
                  onChange={(event) => updateField('rowsText', event.target.value)}
                />
              </Field>
            </div>
          </details>
          <div className="action-row">
            <Button onClick={handleValidate} isWorking={isValidating} workingLabel="校验中">校验预览</Button>
          </div>
        </article>

        <article className="cockpit-card form-card">
          <div className="state-label">脱敏预览</div>
          <h2>校验结果</h2>
          {apiError ? <StatusNotice state={apiError.state} safeMessage={apiError.message} code={apiError.code} /> : null}
          {validation ? (
            <>
              <div className="metric-grid">
                <div><span>批次</span><strong>{validation.import_batch_id}</strong></div>
                <div><span>总数</span><strong>{validation.summary.total_count}</strong></div>
                <div><span>可写入</span><strong>{validation.summary.valid_count}</strong></div>
                <div><span>需关注</span><strong>{validation.summary.warning_count}</strong></div>
                <div><span>阻断</span><strong>{validation.summary.blocking_count}</strong></div>
              </div>
              <div className="table-wrap">
                <table>
                  <thead>
                    <tr>
                      <th>行</th>
                      <th>状态</th>
                      <th>标的</th>
                      <th>标题预览</th>
                      <th>内容预览</th>
                      <th>风险</th>
                    </tr>
                  </thead>
                  <tbody>
                    {validation.rows.map((row) => (
                      <tr key={`${row.row_number}-${row.content_hash}`}>
                        <td>{row.row_number}</td>
                        <td>{statusText(row.status)}</td>
                        <td>{row.symbol || '暂无'}</td>
                        <td>{row.title_preview || '暂无'}</td>
                        <td>{row.text_preview}</td>
                        <td>{row.risks.length ? row.risks.map((risk) => risk.message).join('；') : '无'}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </>
          ) : (
            <p>尚未校验。先提交草稿，确认预览无阻断后再写入。</p>
          )}
        </article>

        <article className="cockpit-card form-card">
          <div className="state-label">索引计划</div>
          <h2>背景事实写入</h2>
          {validation ? (
            <>
              <p>预计片段：{validation.index_plan.rag_chunk_count}</p>
              <p>索引状态：{indexStatusText(validation.index_plan.index_status)}</p>
              <p>{safeSafetyNote()}</p>
              <Field id="local-knowledge-confirm-reason" label="确认理由" hint="必须由人工确认写入本地事实。">
                <input value={form.confirmReason} onChange={(event) => updateField('confirmReason', event.target.value)} />
              </Field>
              <div className="action-row">
                <Button onClick={handleConfirm} disabled={!canConfirm} isWorking={isConfirming} workingLabel="写入中">写入本地事实</Button>
              </div>
            </>
          ) : (
            <p>校验后展示索引计划。</p>
          )}
          {confirmResult ? (
            <div className="state-card state-normal">
              <p>导入批次：{confirmResult.import_batch_id}</p>
              <p>情报：{confirmResult.intelligence_item_count} 条；摘要：{confirmResult.summary_count} 条；片段：{confirmResult.rag_chunk_count} 条。</p>
              <p>核验记录：{confirmResult.verification_count} 条；审计事件：{confirmResult.audit_event_ids.length} 条。</p>
              <p>索引状态：{indexStatusText(confirmResult.index_status)}</p>
            </div>
          ) : null}
        </article>
      </section>
    </div>
  )
}

function parseRows(input: string): { rows: LocalKnowledgeImportRow[]; error?: string } {
  try {
    const parsed = JSON.parse(input) as unknown
    if (!Array.isArray(parsed)) {
      return { rows: [], error: '记录 JSON 必须是数组。' }
    }
    const rows: LocalKnowledgeImportRow[] = parsed.map((item) => {
      const row = item as Record<string, unknown>
      return {
        title: safeString(row.title),
        text: safeString(row.text),
        symbol: safeString(row.symbol),
        as_of_date: safeString(row.as_of_date),
        tags: Array.isArray(row.tags) ? row.tags.map((tag) => safeString(tag)).filter(Boolean) : undefined,
      }
    })
    if (rows.length === 0) {
      return { rows, error: '至少需要一条记录。' }
    }
    if (rows.some((row) => !row.text)) {
      return { rows, error: '每条记录都需要 text。' }
    }
    return { rows }
  } catch {
    return { rows: [], error: '记录 JSON 格式无效。' }
  }
}

function safeString(value: unknown) {
  return typeof value === 'string' ? value.trim() : ''
}

function sourceLabelText(value: string) {
  return value === 'local_research_notes' ? '本地研究记录' : value || '待填写'
}

function statusText(status: string) {
  if (status === 'valid') return '可写入'
  if (status === 'warning') return '需关注'
  if (status === 'blocking') return '阻断'
  return '已记录'
}

function indexStatusText(status: string) {
  if (status === 'pending') return '待重建'
  if (status === 'success') return '完成'
  return status || '暂无'
}

function safeSafetyNote() {
  return '仅写入本地背景材料，后续需由用户在相关页面人工复核。'
}
