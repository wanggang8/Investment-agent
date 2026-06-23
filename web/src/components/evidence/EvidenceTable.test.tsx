import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { EvidenceTable } from './EvidenceTable'

const items = [
  { evidence_id: 'e1', source_name: '公告', source_level: 'A', evidence_role: 'formal', verification_status: 'satisfied', summary: '公告摘要', published_at: '2026-05-31' },
  { evidence_id: 'e2', source_name: '论坛', source_level: 'C', evidence_role: 'background', verification_status: 'background_only', summary: '论坛摘要', original_url: 'https://example.com/e2', content_hash: 'hash-e2', time_weight: 0.4, relevance_score: 0.6, high_grade_independent_source_count: 2 },
]

describe('EvidenceTable', () => {
  it('filters and expands evidence rows', () => {
    render(<EvidenceTable items={items} />)

    fireEvent.change(screen.getByLabelText('筛选证据角色'), { target: { value: 'background' } })
    expect(screen.queryByText('公告摘要')).not.toBeInTheDocument()
    expect(screen.getByText('论坛摘要')).toBeInTheDocument()

    fireEvent.click(screen.getByRole('button', { name: '展开摘要' }))
    expect(screen.getByText('URL')).toBeInTheDocument()
    expect(screen.getByText('hash-e2')).toBeInTheDocument()
    expect(screen.getByText('0.4')).toBeInTheDocument()
    expect(screen.getByText('0.6')).toBeInTheDocument()
    expect(screen.getByText('高等级独立信源数')).toBeInTheDocument()
    expect(screen.getByText('2')).toBeInTheDocument()
  })
})
