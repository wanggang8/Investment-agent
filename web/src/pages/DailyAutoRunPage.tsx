import { useEffect, useState } from 'react'
import { getDailyAutoRunStatus } from '../services/dailyAutoRun'
import type { DailyAutoRunStatus } from '../types/dailyAutoRun'
import { toPageErrorState } from '../shared/utils'
import { buildDailyOpsModel } from '../features/governance'

const statusLabels: Record<DailyAutoRunStatus['status'], string> = {
  disabled: '关闭',
  scheduled: '已计划',
  running: '运行中',
  success: '成功',
  degraded: '部分成功',
  failed: '失败',
}

const failureReasonLabels: Record<string, string> = {
  missing_prerequisites: '缺少本地账户或持仓。',
  timeout: '每日自动运行超时，请查看审计记录。',
  market_refresh_failed: '每日自动运行暂时无法完成，请查看审计记录。',
  evidence_refresh_failed: '每日自动运行暂时无法完成，请查看审计记录。',
  daily_discipline_failed: '每日自动运行暂时无法完成，请查看审计记录。',
}

function safeFailureReason(code?: string) {
  if (!code) {
    return '每日自动运行暂时无法完成，请查看审计记录。'
  }
  return failureReasonLabels[code] ?? '每日自动运行暂时无法完成，请查看审计记录。'
}

export function DailyAutoRunPage() {
  const [status, setStatus] = useState<DailyAutoRunStatus>()
  const [error, setError] = useState('')

  useEffect(() => {
    getDailyAutoRunStatus()
      .then((res) => {
        setStatus(res.data)
        setError('')
      })
      .catch((err: unknown) => {
        const nextError = toPageErrorState(err, '每日自动运行状态加载失败')
        setStatus(undefined)
        setError(nextError.message)
      })
  }, [])

  const dailyModel = buildDailyOpsModel({ autoRun: status })

  return (
    <section className="reference-tight-page">
      <h1 className="page-title">每日自动运行</h1>
      <section className={`daily-hero daily-tone-${dailyModel.overallTone}`} aria-label="每日自动运行总览">
        <div className="daily-hero-main">
          <div className="state-label">每日自动运行健康</div>
          <h2>{dailyModel.overallLabel}</h2>
          <p>{dailyModel.safetyNotes[0]}</p>
          <div className="daily-signal-grid quality-signal-grid">
            {dailyModel.metrics.map((metric) => (
              <article key={metric.label} className={`daily-signal daily-tone-${metric.tone ?? 'unknown'}`}>
                <h2>{metric.label}</h2>
                <strong>{metric.value}</strong>
                {metric.detail ? <p>{metric.detail}</p> : null}
              </article>
            ))}
          </div>
        </div>
        <aside className="daily-hero-side" aria-label="每日自动运行下一步">
          <strong>下一步本地复验</strong>
          <ul>
            {dailyModel.nextActions.map((action) => (
              <li key={action.label}>
                <a href={action.href} aria-label={`${action.label}入口`}>{action.label}</a>
                <span>{action.detail}</span>
              </li>
            ))}
          </ul>
        </aside>
      </section>
      <p className="reference-page-note">展示本地每日刷新、纪律评估、通知与审计状态。当前状态：{status ? statusLabels[status.status] : '加载中'}。</p>
      {error ? <p role="alert">{error}</p> : null}
      {status ? (
        <div className="panel-list">
          <article className="panel-card">
            <div className="row-between">
              <strong>运行状态</strong>
              <span>{status.enabled ? '已启用' : '未启用'}</span>
            </div>
            <p>{status.enabled ? `计划时间：${status.run_time ?? '未配置'} ${status.timezone ?? ''}` : '每日自动运行未启用。'}</p>
            <p>范围：{status.scope === 'holdings' ? '本地当前持仓' : status.scope ?? '未配置'}</p>
            {status.last_run_at ? <p>上次运行：{status.last_run_at}</p> : null}
            {status.next_run_at ? <p>下次运行：{status.next_run_at}</p> : null}
          </article>

          {(status.failure_reason || status.missing_action) ? (
            <article className="panel-card">
              <strong>缺失项与诊断</strong>
              {status.failure_reason ? <p>{safeFailureReason(status.failure_code)}</p> : null}
              {status.failure_code ? <small>错误分类：{status.failure_code}</small> : null}
              {status.missing_action ? <p>{status.missing_action}</p> : null}
            </article>
          ) : null}

          <article className="panel-card">
            <strong>追踪入口</strong>
            <div className="action-row">
              {status.latest_decision_link ? <a href={status.latest_decision_link}>查看最新每日决策</a> : null}
              {status.latest_notification_link ? <a href={status.latest_notification_link}>查看通知</a> : null}
              {status.latest_audit_link ? <a href={status.latest_audit_link}>查看审计详情</a> : null}
            </div>
          </article>

          <article className="panel-card">
            <strong>安全边界</strong>
            <p>{status.safety_note}</p>
          </article>
        </div>
      ) : null}
    </section>
  )
}
