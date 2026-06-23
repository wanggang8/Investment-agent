import type { MarketSnapshot } from './market'

export interface SystemSettings {
  notification_enabled: boolean
  page_preference?: string
  data_sources: string[]
}

export interface CapabilitySettings {
  capability_id?: string
  asset_types?: string[]
  symbols?: string[]
  excluded_symbols?: string[]
  strategy_scope?: string[]
  updated_at?: string
  allowed_asset_types?: string[]
  allowed_symbols?: string[]
  notes?: string
}

export interface SystemStatus {
  sqlite_status: string
  sqlite_path?: string
  veclite_status: string
  veclite_path?: string
  deepseek_status: string
  data_sources: string[]
  log_level: string
}

export interface SettingsPageData {
  system?: SystemStatus
  capability?: CapabilitySettings
  settings?: SystemSettings
  marketSnapshot?: MarketSnapshot
}
