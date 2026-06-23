import type { ChartDatum } from '../../shared/mappers/charts'

interface Props {
  title: string
  data: ChartDatum[]
  emptyText?: string
}

export function MetricBarChart({ title, data, emptyText = '暂无可展示数据。' }: Props) {
  const max = Math.max(...data.map((item) => item.value), 100)

  return (
    <article className="cockpit-card chart-card">
      <div className="state-label">{title}</div>
      {data.length === 0 ? (
        <p className="muted-text">{emptyText}</p>
      ) : (
        <div className="chart-bars" role="list" aria-label={title}>
          {data.map((item) => {
            const width = Math.max(4, Math.round((item.value / max) * 100))
            return (
              <div key={item.label} className="chart-row" role="listitem">
                <div className="chart-row-head">
                  <span>{item.label}</span>
                  <strong>{item.value}</strong>
                </div>
                <div className="chart-track" aria-hidden="true">
                  <span className={`chart-fill chart-fill-${item.tone}`} style={{ width: `${width}%` }} />
                </div>
              </div>
            )
          })}
        </div>
      )}
    </article>
  )
}
