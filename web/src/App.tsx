import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { AppLayout } from './app/AppLayout'
import { AuditPage } from './pages/AuditPage'
import { DataQualityPage } from './pages/DataQualityPage'
import { DailyAutoRunPage } from './pages/DailyAutoRunPage'
import { DailyDisciplineReportDetailPage } from './pages/DailyDisciplineReportDetailPage'
import { DailyDisciplineReportsPage } from './pages/DailyDisciplineReportsPage'
import { DashboardPage } from './pages/DashboardPage'
import { DecisionLoopPage } from './pages/DecisionLoopPage'
import { DecisionDetailPage } from './pages/DecisionDetailPage'
import { EvidencePage } from './pages/EvidencePage'
import { NotificationPage } from './pages/NotificationPage'
import { PortfolioPage } from './pages/PortfolioPage'
import { ReviewSummaryPage } from './pages/ReviewSummaryPage'
import { RiskAlertPage } from './pages/RiskAlertPage'
import { RulesPage } from './pages/RulesPage'
import { SettingsPage } from './pages/SettingsPage'
import { WorkbenchPage } from './pages/WorkbenchPage'
import { LocalInstallPage } from './pages/LocalInstallPage'
import { LocalKnowledgePage } from './pages/LocalKnowledgePage'

// App 定义产品路由入口，所有可见路由均指向已接入业务数据的页面。
function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<AppLayout />}>
          <Route index element={<DashboardPage />} />
          <Route path="workbench" element={<WorkbenchPage />} />
          <Route path="decision-loop" element={<DecisionLoopPage />} />
          <Route path="data-quality" element={<DataQualityPage />} />
          <Route path="positions" element={<PortfolioPage />} />
          <Route path="consultation" element={<DecisionDetailPage />} />
          <Route path="decisions/:decisionId" element={<DecisionDetailPage />} />
          <Route path="evidence" element={<EvidencePage />} />
          <Route path="rules" element={<RulesPage />} />
          <Route path="audit" element={<AuditPage />} />
          <Route path="notifications" element={<NotificationPage />} />
          <Route path="risk-alerts" element={<RiskAlertPage />} />
          <Route path="risk-alerts/:alertId" element={<RiskAlertPage />} />
          <Route path="daily-auto-run" element={<DailyAutoRunPage />} />
          <Route path="daily-discipline/reports" element={<DailyDisciplineReportsPage />} />
          <Route path="daily-discipline/reports/:reportId" element={<DailyDisciplineReportDetailPage />} />
          <Route path="review" element={<ReviewSummaryPage />} />
          <Route path="local-install" element={<LocalInstallPage />} />
          <Route path="local-knowledge" element={<LocalKnowledgePage />} />
          <Route path="settings" element={<SettingsPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
