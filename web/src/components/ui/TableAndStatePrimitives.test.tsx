import { render, screen, within } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { EmptyState } from './EmptyState'
import { ErrorState } from './ErrorState'
import { ResponsiveTable } from './ResponsiveTable'

type Row = {
  id: string
  source: string
  status: string
}

describe('table and state UI primitives', () => {
  it('renders responsive table captions and mobile data labels', () => {
    render(
      <ResponsiveTable<Row>
        caption="数据源健康明细"
        columns={[
          { key: 'source', header: '来源', render: (row) => row.source },
          { key: 'status', header: '状态', render: (row) => row.status },
        ]}
        rows={[{ id: '1', source: 'csindex', status: '过期' }]}
        getRowKey={(row) => row.id}
      />,
    )

    const table = screen.getByRole('table', { name: '数据源健康明细' })
    expect(table).toHaveClass('responsive-table')
    const row = within(table).getByRole('row', { name: /csindex/ })
    expect(within(row).getByText('csindex')).toHaveAttribute('data-label', '来源')
    expect(within(row).getByText('过期')).toHaveAttribute('data-label', '状态')
  })

  it('renders safe empty state with a local manual action', () => {
    render(<EmptyState title="暂无风险预警" description="当前没有需要人工处置的本地风险。" action={{ label: '返回工作台', href: '/workbench' }} />)

    expect(screen.getByRole('region', { name: '暂无风险预警' })).toHaveTextContent('当前没有需要人工处置的本地风险。')
    expect(screen.getByRole('link', { name: '返回工作台' })).toHaveAttribute('href', '/workbench')
  })

  it('renders safe error state without raw secret-shaped diagnostics', () => {
    render(<ErrorState title="加载失败" message="请求失败 sk-1234567890abcdef /Users/private/db.sqlite stack trace SELECT * FROM secrets prompt: raw payload raw vendor payload /tmp/raw.log /opt/private/audit.log C:\\Users\\vick\\db.sqlite" retryLabel="重新加载" onRetry={() => undefined} />)

    const region = screen.getByRole('alert')
    expect(region).toHaveTextContent('加载失败')
    expect(region).not.toHaveTextContent('sk-1234567890abcdef')
    expect(region).not.toHaveTextContent('/Users/private/db.sqlite')
    expect(region).not.toHaveTextContent('stack trace')
    expect(region).not.toHaveTextContent('SELECT * FROM')
    expect(region).not.toHaveTextContent('prompt:')
    expect(region).not.toHaveTextContent('raw vendor payload')
    expect(region).not.toHaveTextContent('/tmp/raw.log')
    expect(region).not.toHaveTextContent('/opt/private/audit.log')
    expect(region).not.toHaveTextContent('C:\\Users\\vick\\db.sqlite')
    expect(screen.getByRole('button', { name: '重新加载' })).toBeVisible()
  })
})
