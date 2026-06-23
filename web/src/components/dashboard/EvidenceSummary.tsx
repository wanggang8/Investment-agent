import type { EvidenceSummary as EvidenceSummaryData } from '../../types/evidence'
import { textOrRaw, verificationStatusText } from '../../shared/mappers'

interface Props {
  summary?: EvidenceSummaryData
  dashboardState: string
}

export function EvidenceSummary({ summary, dashboardState }: Props) {
  const isInsufficient = dashboardState === 'insufficient_data'

  return (
    <article className="cockpit-card">
      <div className="state-label">证据摘要</div>
      {summary ? (
        <dl className="compact-list">
          <div>
            <dt>独立信源</dt>
            <dd>{summary.source_count} 个</dd>
          </div>
          <div>
            <dt>最高等级</dt>
            <dd>{summary.highest_source_level}</dd>
          </div>
          <div>
            <dt>核验状态</dt>
            <dd>{textOrRaw(verificationStatusText, summary.verification_status)}</dd>
          </div>
        </dl>
      ) : (
        <p>{isInsufficient ? '缺少有效证据或索引不可用，暂停交易类建议。' : '暂无证据摘要。'}</p>
      )}
    </article>
  )
}
