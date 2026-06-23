import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import { LocalInstallPage } from './LocalInstallPage'

describe('LocalInstallPage', () => {
  afterEach(() => {
    cleanup()
  })

  it('renders install guidance and forbid automation affordances', () => {
    render(<LocalInstallPage />)

    expect(screen.getByRole('heading', { name: '本地安装与诊断' })).toBeInTheDocument()
    expect(screen.getByText('本地配置与诊断状态')).toBeInTheDocument()
    expect(screen.getByText('复验本地安装')).toBeInTheDocument()
    expect(screen.getByText('查看设置')).toBeInTheDocument()
    expect(screen.getByText(/启动草稿/)).toBeInTheDocument()
    expect(screen.getByText(/关键命令/)).toBeInTheDocument()
    expect(screen.getByText('go run ./cmd/agent --preflight --diagnostics ./tmp/preflight.json')).toBeInTheDocument()
    expect(screen.getByText(/该页用于本地安装引导、配置草稿与诊断打包的只读查看/)).toBeInTheDocument()
    expect(screen.getByText(/本页仅展示本地诊断产物，不读取数据库路径、完整 key、SQL 或原始 HTTP 响应。/)).toBeInTheDocument()
  })

  it('builds and updates startup config draft and renders uploaded summary read-only', async () => {
    render(<LocalInstallPage />)

    const hostInput = screen.getByLabelText('server host')
    const portInput = screen.getByLabelText('server port')
    const sqliteInput = screen.getByLabelText('sqlite 路径')
    const vecliteInput = screen.getByLabelText('veclite 路径')
    fireEvent.change(hostInput, { target: { value: '0.0.0.0' } })
    fireEvent.change(portInput, { target: { value: '9090' } })
    fireEvent.change(sqliteInput, { target: { value: '/opt/private/investment-agent.db' } })
    fireEvent.change(vecliteInput, { target: { value: 'C:\\Users\\vick\\veclite' } })

    const draftBlock = screen.getByLabelText('启动配置草稿')
    expect(draftBlock).toHaveTextContent('host: "0.0.0.0"')
    expect(draftBlock).toHaveTextContent('port: 9090')
    expect(draftBlock).toHaveTextContent('path: "<local-sqlite-path>"')
    expect(draftBlock).toHaveTextContent('path: "<local-veclite-path>"')

    const summary = {
      generated_at: '2026-06-16T12:00:00Z',
      generated_dir: '/Users/private/p44',
      steps: [
        { name: 'preflight', status: 'pass', exit_code: 0, command: 'go run ./cmd/agent --preflight --diagnostics ./tmp/preflight.json', artifact: './tmp/preflight.json' },
        { name: 'recovery_smoke', status: 'failed', exit_code: 1, command: 'bash scripts/recovery-smoke.sh DELETE FROM accounts prompt = sk-secret raw stack /opt/private/file', artifact: 'C:\\Users\\vick\\tmp\\recovery.log' },
        { name: 'DROP TABLE raw_stack', status: 'failed', exit_code: 2, command: 'DROP TABLE audit_events', artifact: '/opt/private/raw.log' },
        { name: 'e2e_smoke', status: 'skipped', exit_code: null, command: '-', artifact: null },
      ],
    }

    const file = new File([JSON.stringify(summary)], 'install-summary.json', {
      type: 'application/json',
    })
    const fileInput = screen.getByLabelText('选择脚本导出的摘要文件')
    fireEvent.change(fileInput, { target: { files: [file] } })

    await waitFor(() => expect(screen.getByText('生成时间：2026-06-16T12:00:00Z')).toBeInTheDocument())
    expect(screen.getByText('失败步骤：2 个')).toBeInTheDocument()
    expect(screen.getByText('preflight', { selector: 'strong' })).toBeInTheDocument()
    expect(screen.getByText('recovery_smoke', { selector: 'strong' })).toBeInTheDocument()
    expect(screen.getByText('e2e_smoke', { selector: 'strong' })).toBeInTheDocument()
    expect(screen.getByText('摘要路径：<local-path>')).toBeInTheDocument()
    expect(document.body.textContent).not.toMatch(/\/Users\/private|\/opt\/private|C:\\Users|SELECT \* FROM|DELETE FROM|DROP TABLE|prompt[:=]|sk-secret|raw stack|raw_stack/)
    expect(screen.queryByLabelText('步骤原始JSON')).not.toBeInTheDocument()
    expect(screen.getByRole('button', { name: '清除展示' })).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: '清除展示' }))
    expect(screen.queryByText('生成时间：2026-06-16T12:00:00Z')).not.toBeInTheDocument()
  })
})
