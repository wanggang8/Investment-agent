import type {
  LocalKnowledgeImportConfirmRequest,
  LocalKnowledgeImportConfirmResponse,
  LocalKnowledgeImportValidationRequest,
  LocalKnowledgeImportValidationResponse,
} from '../types/localKnowledge'
import { apiRequest } from './client'

export function validateLocalKnowledgeImport(body: LocalKnowledgeImportValidationRequest) {
  return apiRequest<LocalKnowledgeImportValidationResponse>('/api/v1/local-knowledge/imports/validate', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export function confirmLocalKnowledgeImport(body: LocalKnowledgeImportConfirmRequest) {
  return apiRequest<LocalKnowledgeImportConfirmResponse>('/api/v1/local-knowledge/imports/confirm', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}
