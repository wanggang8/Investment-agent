import { useMemo, useState, type ChangeEvent } from 'react'
import { Button, Field, SummaryCard, type UITone } from '../components/ui'
import { buildLocalOpsModel } from '../features/governance'
import { redactSensitiveText } from '../shared/utils'

type StepSummary = {
  name: string
  status: string
  exit_code: number | null
  command: string
  artifact?: string | null
}

type InstallSummary = {
  generated_at?: string
  generated_dir?: string
  steps?: StepSummary[]
}

type ConfigDraftForm = {
  serverHost: string
  serverPort: string
  sqlitePath: string
  veclitePath: string
  deepseekBaseUrl: string
  deepseekModel: string
}

const defaultDraft: ConfigDraftForm = {
  serverHost: '127.0.0.1',
  serverPort: '8080',
  sqlitePath: './data/investment-agent.db',
  veclitePath: './data/veclite',
  deepseekBaseUrl: 'https://api.deepseek.com',
  deepseekModel: 'gpt-5.4-mini',
}

function buildConfigDraft(form: ConfigDraftForm) {
  return `server:
  host: "${form.serverHost}"
  port: ${form.serverPort}

sqlite:
  path: "<local-sqlite-path>"

veclite:
  path: "<local-veclite-path>"

deepseek:
  api_key: ""
  base_url: "${form.deepseekBaseUrl}"
  model: "${form.deepseekModel}"
  timeout_seconds: 15
`
}

function redactInstallText(value: string | null | undefined) {
  return redactSensitiveText(value, {
    key: '[REDACTED_KEY]',
    sql: 'SQL redacted',
    prompt: 'prompt redacted',
    raw: 'raw diagnostic redacted',
    stack: 'redacted stack summary',
    path: '<local-path>',
  })
}

function redactInstallSummary(summary: InstallSummary): InstallSummary {
  return {
    generated_at: redactInstallText(summary.generated_at),
    generated_dir: summary.generated_dir ? '<local-path>' : '',
    steps: (summary.steps ?? []).map((step) => ({
      name: redactInstallText(step.name),
      status: redactInstallText(step.status),
      exit_code: step.exit_code,
      command: '',
      artifact: null,
    })),
  }
}

export function LocalInstallPage() {
  const [form, setForm] = useState(defaultDraft)
  const [summary, setSummary] = useState<InstallSummary>()
  const [summaryError, setSummaryError] = useState('')

  const draftText = useMemo(() => buildConfigDraft(form), [form])
  const failedStepCount = summary?.steps?.filter((item) => item.status === 'failed').length ?? 0
  const localModel = buildLocalOpsModel({})

  function updateField(field: keyof ConfigDraftForm, value: string) {
    const nextValue = field === 'sqlitePath' || field === 'veclitePath' ? redactInstallText(value) : value
    setForm((next) => ({ ...next, [field]: nextValue }))
  }

  function handleSummaryUpload(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0]
    if (!file) {
      return
    }

    setSummaryError('')
    file.text().then((content) => {
      try {
        const parsed = JSON.parse(content) as InstallSummary
        setSummary(redactInstallSummary({
          generated_at: parsed.generated_at,
          generated_dir: parsed.generated_dir,
          steps: Array.isArray(parsed.steps) ? parsed.steps : [],
        }))
      } catch (error: unknown) {
        setSummary(undefined)
        setSummaryError(error instanceof Error ? error.message : '摘要文件不是有效 JSON。')
      }
    })
  }

  function clearSummary() {
    setSummary(undefined)
    setSummaryError('')
  }

  return (
    <div className="reference-tight-page">
      <h1 className="page-title">本地安装与诊断</h1>

      <section className={`daily-hero daily-tone-${localModel.overallTone}`} aria-label="本地配置与诊断总览">
        <div className="daily-hero-main">
          <div className="state-label">本地配置与诊断状态</div>
          <h2>{localModel.overallLabel}</h2>
          <p>{localModel.safetyNotes[0]}</p>
          <div className="daily-signal-grid quality-signal-grid">
            {localModel.metrics.map((metric) => (
              <SummaryCard key={metric.label} title={metric.label} value={metric.value} detail={metric.detail} tone={(metric.tone ?? 'unknown') as UITone} />
            ))}
          </div>
        </div>
        <aside className="daily-hero-side" aria-label="本地配置下一步">
          <strong>下一步本地复验</strong>
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
        该页用于本地安装引导、配置草稿与诊断打包的只读查看。这里只提供命令与验证路径，不提供交易触发、规则自动生效入口，也不外连交易通道。
      </p>

      <section className="cockpit-card">
        <div className="state-label">配置向导</div>
        <h2>启动草稿</h2>
        <Field id="local-install-server-host" label="server host" hint="本地监听地址，默认只绑定本机。">
          <input value={form.serverHost} onChange={(event) => updateField('serverHost', event.target.value)} />
        </Field>
        <Field id="local-install-server-port" label="server port">
          <input value={form.serverPort} onChange={(event) => updateField('serverPort', event.target.value)} />
        </Field>
        <Field id="local-install-sqlite-path" label="sqlite 路径" hint="页面展示会脱敏本地路径。">
          <input value={form.sqlitePath} onChange={(event) => updateField('sqlitePath', event.target.value)} />
        </Field>
        <Field id="local-install-veclite-path" label="veclite 路径" hint="页面展示会脱敏本地路径。">
          <input value={form.veclitePath} onChange={(event) => updateField('veclitePath', event.target.value)} />
        </Field>
        <Field id="local-install-deepseek-base-url" label="deepseek base URL">
          <input value={form.deepseekBaseUrl} onChange={(event) => updateField('deepseekBaseUrl', event.target.value)} />
        </Field>
        <Field id="local-install-deepseek-model" label="deepseek model">
          <input value={form.deepseekModel} onChange={(event) => updateField('deepseekModel', event.target.value)} />
        </Field>
        <pre aria-label="启动配置草稿">{draftText}</pre>
      </section>

      <section className="cockpit-card">
        <div className="state-label">关键命令</div>
        <h2>建议命令</h2>
        <ul>
          <li><code>cd /ABSOLUTE/PATH/TO/Investment-agent</code></li>
          <li><code>go run ./cmd/agent --preflight --diagnostics ./tmp/preflight.json</code></li>
          <li><code>go run ./cmd/agent --task market-refresh</code></li>
          <li><code>bash scripts/local-install-diagnostics.sh --skip-e2e</code></li>
          <li><code>bash scripts/e2e-smoke.sh</code></li>
          <li><code>bash scripts/recovery-smoke.sh</code></li>
        </ul>
      </section>

      <section className="cockpit-card">
        <div className="state-label">诊断摘要导入</div>
        <h2>上传 install-summary.json</h2>
        <Field id="local-install-summary-file" label="选择脚本导出的摘要文件" hint="只读取脱敏后的 install-summary.json。">
          <input type="file" accept=".json,application/json" onChange={handleSummaryUpload} />
        </Field>
        {summaryError ? <p role="alert">{summaryError}</p> : null}
        {summary ? (
          <>
            <Button variant="secondary" onClick={clearSummary}>清除展示</Button>
            <p>生成时间：{summary.generated_at || '未知'}</p>
            <p>失败步骤：{failedStepCount} 个</p>
            <p>摘要路径：{summary.generated_dir || '未知'}</p>
            <ul>
              {(summary.steps ?? []).map((step, index) => (
                <li key={`${step.name}-${index}`}>
                  <strong>{step.name}</strong> · {step.status} · code: {step.exit_code ?? 'n/a'}
                </li>
              ))}
            </ul>
          </>
        ) : null}
      </section>

      <section className="cockpit-card">
        <div className="state-label">边界说明</div>
        <h2>安全边界</h2>
        <ul>
          <li>本页仅展示本地诊断产物，不读取数据库路径、完整 key、SQL 或原始 HTTP 响应。</li>
          <li>脚本导出的摘要仅用于复现和汇报，不带入下单动作或投资收益口径判断。</li>
          <li>如发现 failed，可按日志提示修复后重跑脚本，不建议跳过本地预检。</li>
        </ul>
      </section>
    </div>
  )
}
