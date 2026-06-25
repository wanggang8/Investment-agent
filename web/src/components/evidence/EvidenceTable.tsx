import { useMemo, useState } from 'react'
import type { EvidenceItem } from '../../types/evidence'
import { evidenceRoleText, textOrRaw, verificationStatusText } from '../../shared/mappers'

interface Props {
  items: EvidenceItem[]
}

export function EvidenceTable({ items }: Props) {
  const [roleFilter, setRoleFilter] = useState('all')
  const [expandedId, setExpandedId] = useState<string>()
  const filteredItems = useMemo(() => items.filter((item) => roleFilter === 'all' || (item.evidence_role ?? 'formal') === roleFilter), [items, roleFilter])

  return (
    <article className="cockpit-card">
      <div className="state-label">证据列表</div>
      <label className="filter-label">
        筛选证据角色
        <select value={roleFilter} onChange={(event) => setRoleFilter(event.target.value)}>
          <option value="all">全部</option>
          <option value="formal">正式证据</option>
          <option value="background">背景材料</option>
        </select>
      </label>
      {filteredItems.length === 0 ? (
        <p className="muted-text">暂无匹配证据。检索索引不可用或证据不足时，请查看刷新状态。</p>
      ) : (
        <div className="table-wrap">
          <table className="responsive-table">
            <thead>
              <tr>
                <th>信源</th>
                <th>信源等级</th>
                <th>证据角色</th>
                <th>核验状态</th>
                <th>摘要</th>
              </tr>
            </thead>
            <tbody>
              {filteredItems.map((item) => (
                <tr key={item.evidence_id}>
                  <td data-label="信源">{item.source_name}</td>
                  <td data-label="信源等级">{item.source_level} 级</td>
                  <td data-label="证据角色">{textOrRaw(evidenceRoleText, item.evidence_role ?? 'formal')}</td>
                  <td data-label="核验状态">{textOrRaw(verificationStatusText, item.verification_status)}</td>
                  <td data-label="摘要">
                    <button className="link-button" type="button" onClick={() => setExpandedId(expandedId === item.evidence_id ? undefined : item.evidence_id)}>
                      {expandedId === item.evidence_id ? '收起' : '展开'}摘要
                    </button>
                    <p>{item.summary}</p>
                    {expandedId === item.evidence_id && (
                      <dl className="compact-list">
                        <div><dt>发布时间</dt><dd>{item.published_at || '暂无'}</dd></div>
                        <div><dt>捕获时间</dt><dd>{item.captured_at || '暂无'}</dd></div>
                        <div><dt>URL</dt><dd>{item.original_url || '暂无'}</dd></div>
                        <div><dt>内容哈希</dt><dd>{item.content_hash || '暂无'}</dd></div>
                        <div><dt>时间权重</dt><dd>{item.time_weight ?? '暂无'}</dd></div>
                        <div><dt>相关性</dt><dd>{item.relevance_score ?? '暂无'}</dd></div>
                        <div><dt>高等级独立信源数</dt><dd>{item.high_grade_independent_source_count ?? '暂无'}</dd></div>
                      </dl>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </article>
  )
}
