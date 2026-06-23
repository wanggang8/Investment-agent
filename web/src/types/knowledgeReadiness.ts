export interface KnowledgeReadiness {
  symbol: string
  status: string
  symbol_profile: KnowledgeSymbolProfile
  knowledge_references: KnowledgeEntry[]
  data_dependencies: KnowledgeDataDependency[]
  feature_impacts: KnowledgeFeatureImpact[]
  llm_context_summary: string
  safety_notes: string[] | null
}

export interface KnowledgeSymbolProfile {
  symbol: string
  name?: string
  asset_type?: string
  tracked_index_symbol?: string
  tracked_index_name?: string
  known: boolean
}

export interface KnowledgeEntry {
  knowledge_id: string
  title: string
  category: string
  summary: string
  applies_to: string[] | null
  rule_mapping: string[] | null
  llm_context_allowed: boolean
  formal_evidence_allowed: boolean
  safety_boundary: string
}

export interface KnowledgeDataDependency {
  category: string
  status: string
  required: boolean
  source_name?: string
  source_level?: string
  source_type?: string
  freshness?: string
  data_date?: string
  request_id?: string
  affected_symbols?: string[] | null
  affected_features?: string[] | null
  safe_degradation?: string
}

export interface KnowledgeFeatureImpact {
  feature: string
  category: string
  impact: string
  claims?: string[] | null
}
