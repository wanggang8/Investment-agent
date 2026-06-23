import type { PageResult } from '../types/api'
import type { AuditEvent } from '../types/audit'
import { apiRequest } from './client'

export function listAuditEvents() {
  return apiRequest<PageResult<AuditEvent>>('/api/v1/audit-events')
}
