import { DashboardFeature } from '../features/dashboard'

// DashboardPage 仅负责路由组合，业务 UI 位于 features/dashboard。
export function DashboardPage() {
  return <DashboardFeature />
}
