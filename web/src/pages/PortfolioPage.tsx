import { useEffect, useState } from 'react'
import { MetricBarChart } from '../components/charts/MetricBarChart'
import { PortfolioTable } from '../components/portfolio/PortfolioTable'
import { StatusNotice } from '../components/status/StatusNotice'
import { Button, Field } from '../components/ui'
import { buildPortfolioExperienceModel } from '../features/portfolio/portfolioExperienceModel'
import { APIClientError } from '../services/client'
import {
  adjustPortfolio,
  confirmPortfolioImport,
  correctPortfolioFact,
  editHolding,
  getPortfolioCurrent,
  recordOfflineTransaction,
  removeHolding,
  reviewQuarterlyRebalance,
  validatePortfolioImport,
} from '../services/portfolio'
import { buildPortfolioAllocationData } from '../shared/mappers/charts'
import type { PageErrorState } from '../shared/utils'
import { formatCurrency, formatPercent, toPageErrorState } from '../shared/utils'
import type { BatchImportRow, PortfolioCurrentResponse, Position, RebalanceReviewResponse } from '../types/portfolio'

const confirmation = '我确认这只是本地事实记录，不连接券商、不自动交易。'

type PositionForm = {
  cash: string
  totalAssets: string
  symbol: string
  name: string
  quantity: string
  costPrice: string
  currentPrice: string
  buyDate: string
  positionState: 'normal' | 'sell_only' | 'frozen_watch'
  buyReason: string
  assetTag: string
  riskPreference: string
  operationType: 'buy' | 'sell' | 'reduce'
}

const emptyForm: PositionForm = { cash: '', totalAssets: '', symbol: '', name: '', quantity: '', costPrice: '', currentPrice: '', buyDate: '', positionState: 'normal', buyReason: '', assetTag: '', riskPreference: '', operationType: 'buy' }

export function PortfolioPage() {
  const [portfolio, setPortfolio] = useState<PortfolioCurrentResponse>()
  const [calibrationMessage, setCalibrationMessage] = useState('')
  const [actionMessage, setActionMessage] = useState('')
  const [form, setForm] = useState<PositionForm>(emptyForm)
  const [positionStateTouched, setPositionStateTouched] = useState(false)
  const [importBatchID, setImportBatchID] = useState('')
  const [importReady, setImportReady] = useState(false)
  const [rebalanceReview, setRebalanceReview] = useState<RebalanceReviewResponse>()
  const [errorState, setErrorState] = useState<PageErrorState>()

  useEffect(() => {
    loadPortfolio()
  }, [])

  function loadPortfolio() {
    getPortfolioCurrent()
      .then((res) => {
        setPortfolio(res.data)
        setErrorState(undefined)
      })
      .catch((error: unknown) => {
        setPortfolio(undefined)
        if (error instanceof APIClientError && error.code === 'NOT_FOUND') {
          setErrorState(undefined)
          return
        }
        setErrorState(toPageErrorState(error))
      })
  }

  const allocation = buildPortfolioAllocationData(portfolio?.positions ?? [])
  const firstPosition = portfolio?.positions[0]
  const effectivePositionState = positionStateTouched ? form.positionState : ((firstPosition?.position_state as PositionForm['positionState'] | undefined) ?? form.positionState)
  const experience = buildPortfolioExperienceModel({ portfolio, importReady, importBatchID, error: errorState })

  function updateForm(field: keyof PositionForm, value: string) {
    if (field === 'positionState') {
      setPositionStateTouched(true)
    }
    setForm((current) => ({ ...current, [field]: value }))
  }

  function formPositionInput(seed?: Position) {
    const symbol = (form.symbol || seed?.symbol || '').trim()
    const name = (form.name || seed?.name || '').trim()
    const buyReason = (form.buyReason || seed?.buy_reason || '').trim()
    return {
      symbol,
      name,
      quantity: Number(form.quantity || seed?.quantity || 0),
      cost_price: Number(form.costPrice || seed?.cost_price || 0),
      current_price: Number(form.currentPrice || seed?.current_price || 0),
      buy_date: (form.buyDate || seed?.buy_date || '').trim(),
      position_state: positionStateTouched || !seed ? form.positionState : seed.position_state || 'normal',
      buy_reason: buyReason,
      asset_tag: (form.assetTag || seed?.asset_tag || '').trim(),
    }
  }

  function requirePositionInput(seed?: Position) {
    const position = formPositionInput(seed)
    if (!position.symbol || !position.name || position.quantity <= 0 || position.cost_price <= 0 || position.current_price <= 0 || !position.buy_reason) {
      setActionMessage('请先填写账户和持仓必填信息。')
      setCalibrationMessage('')
      return undefined
    }
    return position
  }

  function requireOfflineTransactionInput(seed?: Position) {
    const position = formPositionInput(seed)
    if (!position.symbol || !position.name || position.quantity <= 0 || position.current_price <= 0 || (form.operationType === 'buy' && !position.buy_reason)) {
      setActionMessage('请先填写账户和持仓必填信息。')
      setCalibrationMessage('')
      return undefined
    }
    return position
  }

  function setActionSuccess(message: string) {
    setActionMessage(message)
    setCalibrationMessage('')
    setErrorState(undefined)
    setImportReady(false)
    loadPortfolio()
  }

  function handleError(error: unknown) {
    setActionMessage('')
    setCalibrationMessage('')
    setErrorState(toPageErrorState(error))
  }

  function handleCalibration() {
    if (!firstPosition && !form.cash && !form.totalAssets && !form.symbol && !form.name && !form.quantity && !form.costPrice && !form.currentPrice && !form.buyReason) {
      setActionMessage('请先填写账户和持仓必填信息。')
      setCalibrationMessage('')
      return
    }
    const position = form.symbol || firstPosition ? requirePositionInput(firstPosition) : undefined
    if ((form.symbol || firstPosition) && !position) return
    adjustPortfolio({
      cash: Number(form.cash || portfolio?.snapshot?.cash || 0),
      total_assets: Number(form.totalAssets || portfolio?.snapshot?.total_assets || 0),
      adjust_reason: '用户本地账户校准',
      positions: position ? [position] : [],
    })
      .then(() => {
        setCalibrationMessage('账户校准已保存为本地事实；不会连接交易接口。')
        setActionMessage('')
        setErrorState(undefined)
        loadPortfolio()
      })
      .catch(handleError)
  }

  function handleEditHolding() {
    const position = requirePositionInput(firstPosition)
    if (!position) return
    editHolding({ position_id: firstPosition?.position_id, reason: '用户本地持仓编辑', confirmation, position })
      .then(() => setActionSuccess('持仓编辑已保存为本地事实。'))
      .catch(handleError)
  }

  function handleRemoveHolding() {
    if (!firstPosition?.position_id) return
    removeHolding({ position_id: firstPosition.position_id, reason: '用户本地移除持仓', confirmation })
      .then(() => setActionSuccess('持仓移除已保存为本地事实。'))
      .catch(handleError)
  }

  function handleOfflineTransaction() {
    const position = requireOfflineTransactionInput(firstPosition)
    if (!position) return
    recordOfflineTransaction({
      operation_type: form.operationType,
      symbol: position.symbol,
      name: position.name,
      quantity: position.quantity,
      price: position.current_price,
      fees: 0,
      executed_at: new Date().toISOString(),
      note: '用户补记线下交易',
      buy_reason: position.buy_reason,
      asset_tag: position.asset_tag,
    })
      .then(() => setActionSuccess('线下交易已记录为本地事实。'))
      .catch(handleError)
  }

  function importRows(): BatchImportRow[] | undefined {
    const row = requirePositionInput(firstPosition)
    return row ? [{ row_number: 1, row_type: 'holding', ...row }] : undefined
  }

  function handleValidateImport() {
    const rows = importRows()
    if (!rows) return
    validatePortfolioImport({ rows })
      .then((res) => {
        if (!res.data) {
          throw new Error('导入校验未返回结果')
        }
        setImportBatchID(res.data.import_batch_id)
        setImportReady(res.data.summary.invalid_count === 0)
        setActionMessage(`导入校验完成：有效 ${res.data.summary.valid_count} 行，无效 ${res.data.summary.invalid_count} 行。`)
        setCalibrationMessage('')
        setErrorState(undefined)
      })
      .catch(handleError)
  }

  function handleConfirmImport() {
    if (!importReady || !importBatchID) return
    const rows = importRows()
    if (!rows) return
    confirmPortfolioImport({
      import_batch_id: importBatchID,
      confirm_reason: '确认导入本地账户事实',
      rows,
    })
      .then(() => setActionSuccess('批量导入已确认并保存。'))
      .catch(handleError)
  }

  function handleCorrection() {
    const after = requirePositionInput(firstPosition)
    if (!after) return
    correctPortfolioFact({
      target_type: 'position',
      target_id: firstPosition?.position_id || 'manual_position',
      before_json: JSON.stringify(firstPosition ? { symbol: firstPosition.symbol, quantity: firstPosition.quantity } : {}),
      after_json: JSON.stringify(after),
      correction_reason: '用户更正本地事实',
    })
      .then(() => setActionSuccess('错误修正已保存为本地事实。'))
      .catch(handleError)
  }

  function handleRebalanceReview() {
    reviewQuarterlyRebalance({
      target_core_ratio: 0.5,
      target_satellite_ratio: 0.2,
      target_cash_ratio: 0.3,
      drift_threshold: 0.15,
      review_date: new Date().toISOString().slice(0, 10),
    })
      .then((res) => {
        setRebalanceReview(res.data)
        setActionMessage('季度再平衡复核已生成，仅作为人工计划。')
        setCalibrationMessage('')
        setErrorState(undefined)
      })
      .catch(handleError)
  }

  function bucketText(bucket: string) {
    return ({ core: '核心', satellite: '卫星', cash: '现金' } as Record<string, string>)[bucket] || bucket
  }

  function recommendationText(value: string) {
    return ({
      hold: '保持',
      manual_sell_or_reduce: '人工计划卖出/减仓',
      manual_buy_or_add: '人工计划买入/加仓',
      manual_raise_cash: '人工计划提高现金',
    } as Record<string, string>)[value] || value
  }

  return (
    <div>
      <h1 className="page-title">组合与持仓维护</h1>
      {errorState && <StatusNotice state={errorState.state} safeMessage={errorState.message} code={errorState.code} />}
      <section className={`daily-hero daily-tone-${experience.statusTone}`} aria-label="组合维护状态">
        <div className="daily-hero-main">
          <div className="state-label">组合维护状态</div>
          <h2>{experience.statusLabel}</h2>
          <p>{experience.stageDescription}</p>
          <dl className="daily-hero-meta">
            {experience.summaryMetrics.map((metric) => (
              <div key={metric.label}>
                <dt>{metric.label}</dt>
                <dd>{metric.value}</dd>
              </div>
            ))}
          </dl>
        </div>
        <aside className="daily-hero-side" aria-label="组合下一步">
          <div>
            <strong>{experience.stageLabel}</strong>
            <p className="muted-text">{experience.safetyNotes[0]}</p>
          </div>
          <ul>
            {experience.nextActions.map((action) => (
              <li key={action.label}>
                {action.href ? <a href={action.href}>{action.label}</a> : <strong>{action.label}</strong>}
                <span>{action.detail}</span>
              </li>
            ))}
          </ul>
        </aside>
      </section>
      <section className="daily-signal-grid" aria-label="组合维护模式">
        {experience.maintenanceModes.map((mode) => (
          <article key={mode.id} className="daily-signal daily-tone-success">
            <h2>{mode.id === 'offline_transaction' ? '线下交易补记' : mode.label}</h2>
            <p>{mode.description}</p>
          </article>
        ))}
      </section>
      {portfolio?.snapshot && (
        <div className="cockpit-card">
          <div className="state-label">账户快照</div>
          <p>快照：{portfolio.snapshot.snapshot_id}</p>
          <p>总资产：{formatCurrency(portfolio.snapshot.total_assets)}</p>
          <p>现金：{formatCurrency(portfolio.snapshot.cash)}</p>
          <p>现金占比：{formatPercent(portfolio.snapshot.cash_ratio)}</p>
        </div>
      )}
      <article className="cockpit-card form-card">
        <div className="state-label">首次使用引导</div>
        <h2>初始化/校准</h2>
        <p>这里只记录账户与持仓事实，不提供自动下单能力，也不连接交易接口。</p>
        <div className="form-grid form-grid-wide">
          <Field id="portfolio-cash" label="现金" hint="本地账户现金事实，不连接交易接口。">
            <input inputMode="decimal" value={form.cash} onChange={(event) => updateForm('cash', event.target.value)} />
          </Field>
          <Field id="portfolio-total-assets" label="总资产" hint="用于本地纪律评估的账户总额。">
            <input inputMode="decimal" value={form.totalAssets} onChange={(event) => updateForm('totalAssets', event.target.value)} />
          </Field>
          <Field id="portfolio-symbol" label="标的代码" hint="只记录用户本地持仓事实。">
            <input value={form.symbol} onChange={(event) => updateForm('symbol', event.target.value)} />
          </Field>
          <Field id="portfolio-name" label="标的名称">
            <input value={form.name} onChange={(event) => updateForm('name', event.target.value)} />
          </Field>
          <Field id="portfolio-quantity" label="数量">
            <input inputMode="decimal" value={form.quantity} onChange={(event) => updateForm('quantity', event.target.value)} />
          </Field>
          <Field id="portfolio-cost-price" label="成本价">
            <input inputMode="decimal" value={form.costPrice} onChange={(event) => updateForm('costPrice', event.target.value)} />
          </Field>
          <Field id="portfolio-current-price" label="现价">
            <input inputMode="decimal" value={form.currentPrice} onChange={(event) => updateForm('currentPrice', event.target.value)} />
          </Field>
          <Field id="portfolio-buy-date" label="买入日期" hint="用于持仓回溯与纪律复盘。">
            <input type="date" value={form.buyDate} onChange={(event) => updateForm('buyDate', event.target.value)} />
          </Field>
          <Field id="portfolio-position-state" label="纪律状态" hint="用于标记正常、只卖不买或冻结观察。">
            <select value={effectivePositionState} onChange={(event) => updateForm('positionState', event.target.value as PositionForm['positionState'])}>
              <option value="normal">正常</option>
              <option value="sell_only">只卖不买</option>
              <option value="frozen_watch">冻结观察</option>
            </select>
          </Field>
          <Field id="portfolio-buy-reason" label="买入理由" hint="用于复盘解释，不生成收益承诺。">
            <input value={form.buyReason} onChange={(event) => updateForm('buyReason', event.target.value)} />
          </Field>
          <Field id="portfolio-asset-tag" label="资产标签">
            <input value={form.assetTag} onChange={(event) => updateForm('assetTag', event.target.value)} />
          </Field>
          <Field id="portfolio-risk-preference" label="风险偏好">
            <input value={form.riskPreference} onChange={(event) => updateForm('riskPreference', event.target.value)} />
          </Field>
        </div>
        <div className="action-row">
          <Button onClick={handleCalibration}>保存本地校准</Button>
        </div>
        {calibrationMessage && <p>{calibrationMessage}</p>}
      </article>
      <article className="cockpit-card form-card">
        <div className="state-label">持仓维护</div>
        <p>编辑、移除、修正都只产生本地审计事实。</p>
        <div className="action-row">
          <Button onClick={handleEditHolding}>保存持仓编辑</Button>
          <Button variant="danger" onClick={handleRemoveHolding} disabled={!firstPosition}>移除当前持仓</Button>
        </div>
      </article>
      <article className="cockpit-card form-card">
        <div className="state-label">线下交易记录</div>
        <p>用于补记已由用户自行完成的线下买入、卖出或减仓。</p>
        <div className="form-grid">
          <Field id="portfolio-offline-transaction-type" label="线下交易类型">
            <select aria-label="线下交易类型" value={form.operationType} onChange={(event) => updateForm('operationType', event.target.value as PositionForm['operationType'])}>
              <option value="buy">买入</option>
              <option value="sell">卖出</option>
              <option value="reduce">减仓</option>
            </select>
          </Field>
        </div>
        <div className="action-row">
          <Button onClick={handleOfflineTransaction}>记录线下交易</Button>
        </div>
      </article>
      <article className="cockpit-card form-card">
        <div className="state-label">季度再平衡复核</div>
        <p>按核心 50%、卫星 20%、现金 30% 和 +/-15% 偏离阈值生成离线人工计划。</p>
        <div className="action-row">
          <Button onClick={handleRebalanceReview}>运行季度再平衡复核</Button>
        </div>
        {rebalanceReview && (
          <section className="proposal-item">
            <h3>复核结果</h3>
            <p>复核日期：{rebalanceReview.review_date}</p>
            <p>总资产：{formatCurrency(rebalanceReview.total_assets)}</p>
            <p>偏离阈值：{formatPercent(rebalanceReview.drift_threshold)}</p>
            <ul>
              {rebalanceReview.items.map((item) => (
                <li key={item.bucket}>
                  {bucketText(item.bucket)}：目标 {formatPercent(item.target_ratio)}，实际 {formatPercent(item.actual_ratio)}，偏离 {formatPercent(item.drift_ratio)}，{recommendationText(item.recommendation)} {formatCurrency(item.manual_amount)}
                </li>
              ))}
            </ul>
            <p>{rebalanceReview.safety_statement}</p>
          </section>
        )}
      </article>
      <article className="cockpit-card form-card">
        <div className="state-label">批量导入与错误修正</div>
        <div className="action-row">
          <Button onClick={handleValidateImport}>校验批量导入</Button>
          <Button onClick={handleConfirmImport} disabled={!importReady || !importBatchID}>确认批量导入</Button>
          <Button variant="secondary" onClick={handleCorrection}>记录修正审计</Button>
        </div>
        <p className="muted">修正只记录审计；若需改变当前持仓或现金，请使用持仓编辑或线下交易记录。</p>
        {actionMessage && <p>{actionMessage}</p>}
      </article>
      <MetricBarChart title="持仓结构" data={allocation} emptyText="暂无持仓结构数据。" />
      <PortfolioTable positions={portfolio?.positions ?? []} />
    </div>
  )
}
