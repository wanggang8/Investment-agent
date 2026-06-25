import { expect, test } from '@playwright/test'
import type { Locator, Page } from '@playwright/test'
import { mkdirSync, writeFileSync } from 'node:fs'
import path from 'node:path'

const artifactDir =
  process.env.P119_ARTIFACT_DIR ||
  path.resolve(process.cwd(), '../docs/release/ui-audit-assets/2026-06-25-p119-full-ui-control-and-affordance-acceptance/browser')

type ControlRecord = {
  route_id: string
  route: string
  tag: string
  name: string
  href: string
  disabled: boolean
  category: string
}

const routes = [
  { id: 'R01', path: '/', screenshot: 'r01-dashboard.png' },
  { id: 'R02', path: '/workbench', screenshot: 'r02-workbench.png' },
  { id: 'R03', path: '/decision-loop', screenshot: 'r03-decision-loop.png' },
  { id: 'R04', path: '/data-quality', screenshot: 'r04-data-quality.png' },
  { id: 'R05', path: '/positions', screenshot: 'r05-positions.png' },
  { id: 'R06', path: '/consultation', screenshot: 'r06-consultation.png' },
  { id: 'R07', path: '/decisions', screenshot: 'r07-decisions.png' },
  { id: 'R08', path: '/decisions/decision_p119_confirm', screenshot: 'r08-decision-detail.png' },
  { id: 'R09', path: '/evidence', screenshot: 'r09-evidence.png' },
  { id: 'R10', path: '/rules', screenshot: 'r10-rules.png' },
  { id: 'R11', path: '/audit', screenshot: 'r11-audit.png' },
  { id: 'R12', path: '/notifications', screenshot: 'r12-notifications.png' },
  { id: 'R13', path: '/risk-alerts', screenshot: 'r13-risk-alerts.png' },
  { id: 'R14', path: '/risk-alerts/risk_p119_active', screenshot: 'r14-risk-alert-detail.png' },
  { id: 'R15', path: '/daily-auto-run', screenshot: 'r15-daily-auto-run.png' },
  { id: 'R16', path: '/daily-discipline/reports', screenshot: 'r16-daily-reports.png' },
  { id: 'R17', path: '/daily-discipline/reports/p119_report_01', screenshot: 'r17-daily-report-detail.png' },
  { id: 'R18', path: '/review', screenshot: 'r18-review.png' },
  { id: 'R19', path: '/local-install', screenshot: 'r19-local-install.png' },
  { id: 'R20', path: '/local-knowledge', screenshot: 'r20-local-knowledge.png' },
  { id: 'R21', path: '/settings', screenshot: 'r21-settings.png' },
  { id: 'R22', path: '/api-diagnostics', screenshot: 'r22-api-diagnostics.png' },
]

const mobileRoutes = [
  { id: 'M01', path: '/', screenshot: 'm01-dashboard.png' },
  { id: 'M02', path: '/workbench', screenshot: 'm02-workbench.png' },
  { id: 'M03', path: '/positions', screenshot: 'm03-positions.png' },
  { id: 'M04', path: '/decisions/decision_p119_confirm', screenshot: 'm04-decision-detail.png' },
  { id: 'M05', path: '/data-quality', screenshot: 'm05-data-quality.png' },
  { id: 'M06', path: '/risk-alerts', screenshot: 'm06-risk-alerts.png' },
  { id: 'M07', path: '/local-knowledge', screenshot: 'm07-local-knowledge.png' },
  { id: 'M08', path: '/settings', screenshot: 'm08-settings.png' },
]

test('P119 inventories all routes, controls, layout states, and key UI-backed writes', async ({ page, request }) => {
  test.setTimeout(300_000)
  mkdirSync(artifactDir, { recursive: true })

  const consoleErrors: string[] = []
  const pageErrors: string[] = []
  const failedApiResponses: Array<{ url: string; status: number; method: string }> = []
  page.on('console', (message) => {
    if (message.type() === 'error' && !message.text().includes('Failed to load resource')) consoleErrors.push(message.text())
  })
  page.on('pageerror', (error) => pageErrors.push(error.message))
  page.on('response', (response) => {
    if (response.url().includes('/api/v1/') && response.status() >= 500) {
      failedApiResponses.push({ url: response.url(), status: response.status(), method: response.request().method() })
    }
  })

  await expect(await request.get('/api/v1/health')).toBeOK()

  const routeEvidence: any[] = []
  const controls: ControlRecord[] = []
  const layoutIssues: any[] = []
  const productCopyIssues: any[] = []

  await page.setViewportSize({ width: 1440, height: 960 })
  for (const route of routes) {
    const evidence = await visitAndInventory(page, route.path, route.id, route.screenshot)
    routeEvidence.push(evidence)
    controls.push(...evidence.controls)
    layoutIssues.push(...evidence.layout_issues)
    productCopyIssues.push(...evidence.product_copy_issues)
  }

  const toggleEvidence = await exerciseUpstreamToggleInteractions(page)
  const actionEvidence = await exerciseKeyActions(page)

  await page.setViewportSize({ width: 390, height: 844 })
  const mobileEvidence: any[] = []
  for (const route of mobileRoutes) {
    const evidence = await visitAndInventory(page, route.path, route.id, route.screenshot, true)
    mobileEvidence.push(evidence)
    layoutIssues.push(...evidence.layout_issues)
    productCopyIssues.push(...evidence.product_copy_issues)
  }

  const unnamedControls = controls.filter((control) => !control.name)
  const unclassifiedControls = controls.filter((control) => control.category === 'unclassified')
  const categoryCounts = controls.reduce<Record<string, number>>((acc, control) => {
    acc[control.category] = (acc[control.category] || 0) + 1
    return acc
  }, {})

  expect(unnamedControls).toEqual([])
  expect(unclassifiedControls).toEqual([])
  expect(layoutIssues).toEqual([])
  expect(productCopyIssues).toEqual([])
  expect(consoleErrors).toEqual([])
  expect(pageErrors).toEqual([])
  expect(failedApiResponses).toEqual([])

  const payload = {
    status: 'passed',
    generated_at: new Date().toISOString(),
    route_count: routeEvidence.length,
    mobile_route_count: mobileEvidence.length,
    control_count: controls.length,
    control_category_counts: categoryCounts,
    unnamed_controls: unnamedControls.length,
    unclassified_controls: unclassifiedControls.length,
    layout_issue_count: layoutIssues.length,
    product_copy_issue_count: productCopyIssues.length,
    console_error_count: consoleErrors.length,
    page_error_count: pageErrors.length,
    api_5xx_count: failedApiResponses.length,
    toggle_interaction_count: toggleEvidence.length,
    toggle_issue_count: 0,
    route_evidence: routeEvidence,
    mobile_evidence: mobileEvidence,
    toggle_evidence: toggleEvidence,
    action_evidence: actionEvidence,
    console_errors: consoleErrors,
    page_errors: pageErrors,
    failed_api_responses: failedApiResponses,
    safety_evidence: {
      forbidden_visible_affordances: 0,
      broker_order_buttons: 0,
      return_guarantee_claims: 0,
    },
  }
  writeFileSync(path.join(artifactDir, 'p119-browser-results.json'), JSON.stringify(payload, null, 2))
})

async function visitAndInventory(page: Page, route: string, routeId: string, screenshot: string, mobile = false) {
  await page.goto(route)
  await expect(page.locator('h1.page-title').first()).toBeVisible()
  await expect(page.locator('body')).not.toHaveText('')
  await page.screenshot({ path: path.join(artifactDir, screenshot), fullPage: true })
  const bodyText = await page.locator('body').innerText()
  const controls = await collectControls(page, route, routeId)
  const layoutIssues = await collectLayoutIssues(page, route, routeId)
  const productCopyIssues = productCopyScan(bodyText, route, routeId)
  await assertNoForbiddenVisibleAffordance(bodyText)
  return {
    route_id: routeId,
    route,
    viewport: page.viewportSize(),
    screenshot_path: path.join(artifactDir, screenshot),
    heading: await page.locator('h1.page-title').first().innerText(),
    body_text_length: bodyText.length,
    controls,
    control_count: controls.length,
    layout_issues: layoutIssues,
    product_copy_issues: productCopyIssues,
    mobile,
  }
}

async function collectControls(page: Page, route: string, routeId: string): Promise<ControlRecord[]> {
  const raw = await page.evaluate(() => {
    function nameFor(el: Element) {
      const element = el as HTMLElement
      const direct = element.getAttribute('aria-label') || element.getAttribute('title') || element.innerText || element.textContent || ''
      if (direct.trim()) return direct.replace(/\s+/g, ' ').trim()
      const id = element.getAttribute('id')
      if (id) {
        const label = document.querySelector(`label[for="${CSS.escape(id)}"]`)
        if (label?.textContent?.trim()) return label.textContent.replace(/\s+/g, ' ').trim()
      }
      const wrapped = element.closest('label')
      if (wrapped?.textContent?.trim()) return wrapped.textContent.replace(/\s+/g, ' ').trim()
      return element.getAttribute('placeholder') || element.getAttribute('name') || id || ''
    }

    return Array.from(document.querySelectorAll('button,a[href],input,select,textarea,summary'))
      .map((el) => {
        const element = el as HTMLElement
        const rect = element.getBoundingClientRect()
        const style = window.getComputedStyle(element)
        const closedDetails = element.tagName === 'SUMMARY' ? null : element.closest('details:not([open])')
        return {
          tag: element.tagName.toLowerCase(),
          name: nameFor(element),
          href: element instanceof HTMLAnchorElement ? element.getAttribute('href') || '' : '',
          disabled: element instanceof HTMLButtonElement || element instanceof HTMLInputElement || element instanceof HTMLSelectElement || element instanceof HTMLTextAreaElement ? element.disabled : false,
          visible: !closedDetails && rect.width > 0 && rect.height > 0 && style.display !== 'none' && style.visibility !== 'hidden',
        }
      })
      .filter((item) => item.visible)
  })
  return raw.map((item) => ({
    route_id: routeId,
    route,
    tag: item.tag,
    name: item.name,
    href: item.href,
    disabled: item.disabled,
    category: classifyControl(item.tag, item.name, item.href, item.disabled),
  }))
}

function classifyControl(tag: string, name: string, href: string, disabled: boolean) {
  const text = name.trim()
  if (disabled) return 'disabled_expected'
  if (tag === 'a' && href) return 'navigation'
  if (tag === 'summary') return 'light_interaction'
  if (tag === 'input' || tag === 'select' || tag === 'textarea') return 'form_input'
  if (['导航', '展开证据链', '收起证据链', '清除展示'].includes(text) || text.startsWith('展开') || text.startsWith('收起') || text.startsWith('查看')) return 'light_interaction'
  if (/(刷新摘要|刷新当前页面摘要|切换标的|检查门禁处置|刷新情报|重建索引|校验预览|刷新市场数据|运行季度再平衡复核|校验批量导入)/.test(text)) return 'read_action'
  if (/(保存|提交咨询|记录计划|标记待观察|已手动执行|标记错误|记录线下交易|移除当前持仓|提交确认|全部标记已读|标记已读|记录继续观察|记录升级复核|记录本地解除预警|记录处置|退役处置|写入本地事实|记录修正审计)/.test(text)) return 'write_local_fact'
  if (/(生成 SOP 补充提案|确认送审|拒绝提案|确认应用到正式规则|拒绝应用)/.test(text)) return 'governance_confirm'
  return 'unclassified'
}

async function collectLayoutIssues(page: Page, route: string, routeId: string) {
  return page.evaluate(({ route, routeId }) => {
    const issues: any[] = []
    const interactive = Array.from(document.querySelectorAll('button,a[href],input,select,textarea,summary')) as HTMLElement[]
    const viewportWidth = window.innerWidth
    for (const el of interactive) {
      const rect = el.getBoundingClientRect()
      const style = window.getComputedStyle(el)
      if (el.tagName !== 'SUMMARY' && el.closest('details:not([open])')) continue
      if (rect.width <= 0 || rect.height <= 0 || style.display === 'none' || style.visibility === 'hidden') continue
      if (rect.left < -2 || rect.right > viewportWidth + 2) {
        issues.push({ route_id: routeId, route, type: 'interactive_overflow', name: el.innerText || el.getAttribute('aria-label') || el.getAttribute('id') || el.tagName, left: rect.left, right: rect.right, viewportWidth })
      }
      if ((el.tagName === 'BUTTON' || el.tagName === 'A') && rect.height > 96) {
        issues.push({ route_id: routeId, route, type: 'control_height_suspicious', name: el.innerText || el.getAttribute('aria-label') || el.tagName, height: rect.height })
      }
    }
    const elements = interactive.filter((el) => {
      const rect = el.getBoundingClientRect()
      const style = window.getComputedStyle(el)
      return !el.closest('details:not([open])') && ['BUTTON', 'INPUT', 'SELECT', 'TEXTAREA'].includes(el.tagName) && rect.width > 0 && rect.height > 0 && style.display !== 'none' && style.visibility !== 'hidden'
    })
    for (let i = 0; i < elements.length; i += 1) {
      const a = elements[i].getBoundingClientRect()
      for (let j = i + 1; j < elements.length; j += 1) {
        const b = elements[j].getBoundingClientRect()
        const x = Math.max(0, Math.min(a.right, b.right) - Math.max(a.left, b.left))
        const y = Math.max(0, Math.min(a.bottom, b.bottom) - Math.max(a.top, b.top))
        if (x * y > 64) {
          const sameElementNesting = elements[i].contains(elements[j]) || elements[j].contains(elements[i])
          if (!sameElementNesting) issues.push({ route_id: routeId, route, type: 'interactive_overlap', first: elements[i].innerText || elements[i].tagName, second: elements[j].innerText || elements[j].tagName })
        }
      }
    }
    return issues
  }, { route, routeId })
}

function productCopyScan(bodyText: string, route: string, routeId: string) {
  const badPatterns = [/undefined/i, /\bNaN\b/i, /stack trace/i, /panic:/i, /sk-[A-Za-z0-9]{8,}/, /TODO/i]
  const issues: any[] = []
  for (const pattern of badPatterns) {
    if (pattern.test(bodyText)) issues.push({ route_id: routeId, route, pattern: String(pattern) })
  }
  return issues
}

async function assertNoForbiddenVisibleAffordance(bodyText: string) {
  const forbidden = ['一键交易', '代下单', '券商下单', '下单按钮', '收益保证', '保证收益']
  for (const term of forbidden) {
    expect(hasUnsafeAffordance(bodyText, term)).toBe(false)
  }
}

function hasUnsafeAffordance(bodyText: string, term: string) {
  let offset = bodyText.indexOf(term)
  while (offset >= 0) {
    const before = bodyText.slice(Math.max(0, offset - 8), offset)
    const after = bodyText.slice(offset + term.length, offset + term.length + 8)
    const boundaryContext = /(不|无|非|禁止|不会|不能|不得|未|只记录|不提供|不构成|不支持)/.test(before) || /(不可用|不存在|能力)/.test(after)
    if (!boundaryContext) return true
    offset = bodyText.indexOf(term, offset + term.length)
  }
  return false
}

async function exerciseKeyActions(page: Page) {
  const evidence: any[] = []

  await page.setViewportSize({ width: 1440, height: 960 })

  await page.goto('/positions')
  await page.getByLabel('标的代码').fill('512000')
  await page.getByLabel('标的名称').fill('券商ETF')
  await page.getByLabel('数量').fill('10')
  await page.getByLabel('成本价').fill('1.20')
  await page.getByLabel('现价').fill('1.25')
  await page.getByLabel('买入日期').fill('2026-06-20')
  await page.getByLabel('买入理由').fill('P119 UI 控件真实点击补记')
  await page.getByLabel('资产标签').fill('satellite')
  await page.getByRole('button', { name: '记录线下交易' }).click()
  await expect(page.getByText('线下交易已记录为本地事实。')).toBeVisible()
  evidence.push({ route: '/positions', action: '记录线下交易', assertion: 'success message visible' })

  await page.goto('/decisions/decision_p119_confirm')
  await page.getByRole('button', { name: '已手动执行' }).click()
  const confirmForm = page.locator('[aria-label="确认表单"]')
  await confirmForm.locator('label').filter({ hasText: '标的代码' }).locator('input').fill('510300')
  await confirmForm.locator('label').filter({ hasText: '数量' }).locator('input').fill('5')
  await confirmForm.locator('label').filter({ hasText: '价格' }).locator('input').fill('4.10')
  await confirmForm.locator('label').filter({ hasText: '费用' }).locator('input').fill('0')
  await confirmForm.locator('label').filter({ hasText: '执行时间' }).locator('input').fill('2026-06-24T09:00')
  await confirmForm.getByRole('button', { name: '提交确认' }).click()
  await expect(page.getByText('确认已记录。')).toBeVisible()
  evidence.push({ route: '/decisions/decision_p119_confirm', action: '已手动执行/提交确认', assertion: 'confirmation recorded' })

  await page.goto('/decisions/decision_p119_error')
  await page.getByRole('button', { name: '标记错误' }).click()
  const errorForm = page.locator('[aria-label="确认表单"]')
  await errorForm.locator('label').filter({ hasText: '实际结果' }).locator('input').fill('P119 实际结果与建议不一致')
  await errorForm.locator('label').filter({ hasText: '原因标签' }).locator('select').selectOption('user_context_missing')
  await errorForm.locator('label').filter({ hasText: '复盘记录' }).locator('textarea').fill('P119 错误标注用于验证复盘链路')
  await errorForm.getByRole('button', { name: '提交确认' }).click()
  await expect(page.getByText('确认已记录。')).toBeVisible()
  evidence.push({ route: '/decisions/decision_p119_error', action: '标记错误/提交确认', assertion: 'error case recorded' })

  await page.goto('/risk-alerts/risk_p119_active')
  await page.getByRole('button', { name: '记录本地解除预警' }).click()
  await expect(page.getByText('已解除')).toBeVisible()
  evidence.push({ route: '/risk-alerts/risk_p119_active', action: '记录本地解除预警', assertion: 'risk lifecycle resolved' })

  await page.goto('/notifications')
  await page.getByRole('button', { name: '全部标记已读' }).click()
  await expect(page.getByRole('button', { name: '全部标记已读' })).toBeDisabled()
  evidence.push({ route: '/notifications', action: '全部标记已读', assertion: 'mark-all disabled after read' })

  await page.goto('/data-quality?symbol=000300')
  const recordResolution = page.getByRole('button', { name: '记录处置' })
  await recordResolution.scrollIntoViewIfNeeded()
  await expect(recordResolution).toBeEnabled()
  await recordResolution.click()
  await expect(page.getByText('范围排除').first()).toBeVisible()
  evidence.push({ route: '/data-quality?symbol=000300', action: '记录处置', assertion: 'data quality resolution visible' })

  await page.goto('/rules')
  await page.getByRole('button', { name: '生成 SOP 补充提案' }).click()
  await expect(page.getByText('SOP 补充提案已生成，等待人工确认。')).toBeVisible()
  const confirmProposal = page.getByRole('button', { name: '确认送审' }).first()
  if (await confirmProposal.isVisible()) {
    await confirmProposal.click()
    await expect(page.getByText('已确认送审。')).toBeVisible()
  }
  evidence.push({ route: '/rules', action: '生成 SOP 补充提案/确认送审', assertion: 'rule governance action visible' })

  await page.goto('/local-knowledge')
  await page.getByRole('button', { name: '校验预览' }).click()
  await expect(page.getByRole('button', { name: '写入本地事实' })).toBeEnabled()
  await page.getByRole('button', { name: '写入本地事实' }).click()
  await expect(page.getByText('导入批次：')).toBeVisible()
  evidence.push({ route: '/local-knowledge', action: '校验预览/写入本地事实', assertion: 'local knowledge imported' })

  await page.goto('/evidence')
  await page.getByRole('button', { name: '重建索引' }).click()
  await expect(page.getByText(/索引重建完成|索引不可用/)).toBeVisible()
  evidence.push({ route: '/evidence', action: '重建索引', assertion: 'index rebuild success or productized unavailable state visible' })

  await page.goto('/settings')
  const marketRefresh = page.getByRole('button', { name: '刷新市场数据' })
  if (await marketRefresh.isVisible()) {
    await marketRefresh.click()
    await expect(page.getByText('市场刷新完成；只更新本地行情事实和审计记录，不会执行交易。')).toBeVisible()
  }
  evidence.push({ route: '/settings', action: '刷新市场数据', assertion: 'market refresh visible if supported' })

  await page.goto('/consultation')
  await page.getByRole('button', { name: '提交咨询' }).click()
  await expect(page.getByText('请填写咨询问题和标的代码。')).toBeVisible()
  evidence.push({ route: '/consultation', action: '提交咨询空表单', assertion: 'validation message visible and no write' })

  return evidence
}

async function exerciseUpstreamToggleInteractions(page: Page) {
  const evidence: any[] = []

  await page.setViewportSize({ width: 1440, height: 960 })

  await page.goto('/workbench')
  await page.getByRole('button', { name: '刷新摘要' }).click()
  await expect(page.locator('h1.page-title')).toHaveText('用户决策工作台')
  evidence.push({ route: '/workbench', control: '刷新摘要', assertion: 'reload keeps route and title stable' })

  await page.setViewportSize({ width: 390, height: 844 })
  await page.goto('/')
  const navToggle = page.getByRole('button', { name: '导航' })
  await expect(navToggle).toHaveAttribute('aria-expanded', 'false')
  await navToggle.click()
  await expect(navToggle).toHaveAttribute('aria-expanded', 'true')
  await page.getByRole('link', { name: '数据质量' }).click()
  await expect(page.locator('h1.page-title')).toHaveText('数据质量可观测')
  await expect(navToggle).toHaveAttribute('aria-expanded', 'false')
  evidence.push({ route: '/', control: '导航', assertion: 'mobile nav opens, navigates, and closes' })

  await page.setViewportSize({ width: 1440, height: 960 })

  await page.goto('/data-quality?symbol=000300')
  const symbolInput = page.getByRole('textbox', { name: '数据质量标的' })
  await symbolInput.fill('159915')
  await page.getByRole('button', { name: '切换标的' }).click()
  await expect(page).toHaveURL(/symbol=159915/)
  await expect(page.getByText('当前查看：159915')).toBeVisible()
  evidence.push({ route: '/data-quality', control: '切换标的', assertion: 'query symbol and visible state update' })
  await openDetails(page, '查看其余知识引用')
  evidence.push({ route: '/data-quality', control: '查看其余知识引用', assertion: 'details opens' })
  await openDetails(page, '查看审计线索')
  evidence.push({ route: '/data-quality', control: '查看审计线索', assertion: 'details opens' })

  await page.goto('/decisions/decision_p119_confirm')
  const evidenceButton = page.getByRole('button', { name: '收起证据链' })
  await evidenceButton.click()
  await expect(page.getByRole('button', { name: '展开证据链' })).toBeVisible()
  await page.getByRole('button', { name: '展开证据链' }).click()
  await expect(page.getByRole('button', { name: '收起证据链' })).toBeVisible()
  evidence.push({ route: '/decisions/decision_p119_confirm', control: '展开/收起证据链', assertion: 'button label and panel state toggle both ways' })
  await page.getByRole('button', { name: /展开 \d+ 份分析材料/ }).click()
  await expect(page.getByText('P119LocalAnalyst')).toBeVisible()
  await page.getByRole('button', { name: '收起分析材料' }).click()
  await expect(page.getByText('P119LocalAnalyst')).toBeHidden()
  evidence.push({ route: '/decisions/decision_p119_confirm', control: '展开/收起分析材料', assertion: 'analysis panel opens and closes' })

  await page.goto('/evidence')
  await page.getByLabel('筛选证据角色').selectOption('background')
  await expect(page.getByLabel('筛选证据角色')).toHaveValue('background')
  await page.getByLabel('筛选证据角色').selectOption('formal')
  await expect(page.getByLabel('筛选证据角色')).toHaveValue('formal')
  await expect(page.getByRole('button', { name: '展开摘要' }).first()).toBeVisible()
  await page.getByRole('button', { name: '展开摘要' }).first().click()
  await expect(page.getByText('内容哈希')).toBeVisible()
  await page.getByRole('button', { name: '收起摘要' }).first().click()
  await expect(page.getByText('内容哈希')).toBeHidden()
  evidence.push({ route: '/evidence', control: '筛选证据角色/展开摘要', assertion: 'filter and row expansion toggle correctly' })

  await page.goto('/audit')
  await page.getByLabel('筛选审计状态').selectOption('success')
  const auditExpand = page.getByRole('button', { name: '展开引用' }).first()
  await auditExpand.click()
  await expect(page.getByRole('region', { name: '审计引用详情' }).first()).toBeVisible()
  await expect(page.getByRole('button', { name: '收起引用' }).first()).toHaveAttribute('aria-expanded', 'true')
  await page.getByRole('button', { name: '收起引用' }).first().click()
  await expect(page.getByRole('region', { name: '审计引用详情' })).toHaveCount(0)
  evidence.push({ route: '/audit', control: '筛选审计状态/展开引用', assertion: 'filter and aria-expanded detail toggle work' })

  await page.goto('/rules')
  await openDetails(page, '查看规则阈值')
  evidence.push({ route: '/rules', control: '查看规则阈值', assertion: 'details opens' })
  await openDetails(page, '查看影响范围')
  evidence.push({ route: '/rules', control: '查看影响范围', assertion: 'details opens' })

  await page.goto('/local-knowledge')
  await openDetails(page, '编辑结构化记录')
  await expect(page.getByLabel('记录 JSON')).toBeVisible()
  evidence.push({ route: '/local-knowledge', control: '编辑结构化记录', assertion: 'details exposes editable local JSON textarea' })

  await page.goto('/local-install')
  await openDetails(page, '编辑配置并查看配置文本')
  await expect(page.getByLabel('server host')).toBeVisible()
  await openDetails(page, '查看本地复验命令')
  await expect(page.getByText('bash scripts/local-install-diagnostics.sh --skip-e2e')).toBeVisible()
  evidence.push({ route: '/local-install', control: '编辑配置/查看复验命令', assertion: 'details expose config and command blocks' })

  await page.goto('/positions')
  await page.getByLabel('纪律状态').selectOption('frozen_watch')
  await expect(page.getByLabel('纪律状态')).toHaveValue('frozen_watch')
  await page.getByLabel('线下交易类型').selectOption('sell')
  await expect(page.getByLabel('线下交易类型')).toHaveValue('sell')
  evidence.push({ route: '/positions', control: '纪律状态/线下交易类型', assertion: 'select controls update without submitting' })

  await page.goto('/consultation')
  await page.getByLabel('咨询场景').selectOption('rebalance_review')
  await expect(page.getByLabel('咨询场景')).toHaveValue('rebalance_review')
  evidence.push({ route: '/consultation', control: '咨询场景', assertion: 'scenario select updates without submitting' })

  evidence.push(...await exerciseAllDetailsSummaryInstances(page))
  return evidence
}

async function exerciseAllDetailsSummaryInstances(page: Page) {
  const evidence: any[] = []
  await page.setViewportSize({ width: 1440, height: 960 })
  for (const route of routes) {
    await page.goto(route.path)
    await expect(page.locator('h1.page-title').first()).toBeVisible()
    const summaries = page.locator('summary')
    const summaryCount = await summaries.count()
    for (let index = 0; index < summaryCount; index += 1) {
      const summary = summaries.nth(index)
      if (!(await summary.isVisible())) continue
      const label = (await summary.innerText()).replace(/\s+/g, ' ').trim()
      await summary.scrollIntoViewIfNeeded()
      const initiallyOpen = await isSummaryParentOpen(summary)
      await summary.click()
      await expect.poll(async () => isSummaryParentOpen(summary)).toBe(!initiallyOpen)
      await summary.click()
      await expect.poll(async () => isSummaryParentOpen(summary)).toBe(initiallyOpen)
      evidence.push({
        route: route.path,
        control: label,
        assertion: 'visible details summary toggles open and closed',
        sweep: 'all_visible_details_summaries',
      })
    }
  }
  return evidence
}

async function openDetails(page: Page, summaryText: string) {
  const summary = page.locator('summary', { hasText: summaryText }).first()
  await expect(summary).toBeVisible()
  await summary.click()
  await expect.poll(async () => isSummaryParentOpen(summary)).toBe(true)
}

async function isSummaryParentOpen(summary: Locator) {
  return summary.evaluate((element) => Boolean(element.parentElement && 'open' in element.parentElement && (element.parentElement as HTMLDetailsElement).open))
}
