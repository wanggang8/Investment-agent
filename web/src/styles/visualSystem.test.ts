/// <reference types="node" />

import { readFileSync } from 'node:fs'
import { describe, expect, it } from 'vitest'

const css = readFileSync(`${process.cwd()}/src/styles/global.css`, 'utf8')

describe('P110 visual system contract', () => {
  it('defines command-center tokens that apply beyond a single page', () => {
    expect(css).toContain('--color-command-bg')
    expect(css).toContain('--color-command-surface')
    expect(css).toContain('--color-command-accent')
    expect(css).toContain('--shadow-command-panel')
  })

  it('provides shared layout language for command, ledger, and manual-action surfaces', () => {
    expect(css).toContain('.command-center-shell')
    expect(css).toContain('.ledger-surface')
    expect(css).toContain('.manual-action-panel')
    expect(css).toContain('.trust-signal-strip')
  })

  it('keeps responsive guardrails for the redesigned system', () => {
    expect(css).toContain('@media (max-width: 760px)')
    expect(css).toContain('.command-center-shell')
    expect(css).toContain('.ledger-surface')
  })
})
