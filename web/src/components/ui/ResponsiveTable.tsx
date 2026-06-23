import type { ReactNode } from 'react'

export type ResponsiveTableColumn<T> = {
  key: string
  header: string
  render: (row: T) => ReactNode
}

type ResponsiveTableProps<T> = {
  caption: string
  columns: Array<ResponsiveTableColumn<T>>
  rows: T[]
  getRowKey: (row: T) => string
  emptyText?: string
}

export function ResponsiveTable<T>({ caption, columns, rows, getRowKey, emptyText = '暂无记录。' }: ResponsiveTableProps<T>) {
  if (!rows.length) {
    return <p className="muted-text">{emptyText}</p>
  }

  return (
    <div className="table-wrap ui-responsive-table">
      <table className="responsive-table">
        <caption>{caption}</caption>
        <thead>
          <tr>
            {columns.map((column) => (
              <th key={column.key} scope="col">{column.header}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row) => (
            <tr key={getRowKey(row)}>
              {columns.map((column) => (
                <td key={column.key} data-label={column.header}>{column.render(row)}</td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
