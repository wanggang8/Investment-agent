import type { PageResult } from '../types/api'
import type {
  EvidenceItem,
  EvidenceRefreshRequest,
  EvidenceRefreshResponse,
  RebuildIndexResponse,
  SourceVerification,
} from '../types/evidence'
import { apiRequest } from './client'

export function refreshEvidence(body: EvidenceRefreshRequest) {
  return apiRequest<EvidenceRefreshResponse>('/api/v1/evidence/refresh', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export function listEvidence() {
  return apiRequest<PageResult<EvidenceItem>>('/api/v1/evidence')
}

export function getEvidenceVerification() {
  return apiRequest<SourceVerification>('/api/v1/evidence/verification')
}

export function rebuildEvidenceIndex() {
  return apiRequest<RebuildIndexResponse>('/api/v1/evidence/rebuild-index', {
    method: 'POST',
  })
}
