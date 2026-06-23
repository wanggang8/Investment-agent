import type { CapabilitySettings, SystemSettings, SystemStatus } from '../types/settings'
import { apiRequest } from './client'

export function getSystemSettings() {
  return apiRequest<SystemStatus>('/api/v1/settings/system')
}

export function updateSystemSettings(body: SystemSettings) {
  return apiRequest<SystemSettings>('/api/v1/settings', {
    method: 'PUT',
    body: JSON.stringify(body),
  })
}

export function getCapabilitySettings() {
  return apiRequest<CapabilitySettings>('/api/v1/settings/capability')
}

export function updateCapabilitySettings(body: CapabilitySettings) {
  return apiRequest<CapabilitySettings>('/api/v1/settings/capability', {
    method: 'PUT',
    body: JSON.stringify(body),
  })
}
