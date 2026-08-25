// These vocabularies mirror fixture schema v1. Keeping the renderer side explicit makes an
// unsupported fixture value fail before it can degrade into a generic visual treatment.
export const regionKinds = [
  'frame',
  'surface',
  'row',
  'column',
  'tabs',
  'list',
  'table',
  'status',
  'overlay',
] as const

export const contentKinds = [
  'text',
  'code',
  'tab',
  'list-item',
  'table-cell',
  'status-item',
  'selection',
  'focus-ring',
  'cursor',
  'badge',
] as const

export const semanticStates = [
  'default',
  'active',
  'inactive',
  'focused',
  'selected',
  'muted',
  'success',
  'warning',
  'error',
  'info',
  'urgent',
  'added',
  'removed',
  'modified',
  'search-match',
] as const

export type RegionKind = (typeof regionKinds)[number]
export type ContentKind = (typeof contentKinds)[number]
export type SemanticState = (typeof semanticStates)[number]

export function formatCandidateLabel(candidateNumber: number): string {
  if (!Number.isInteger(candidateNumber) || candidateNumber < 1 || candidateNumber > 8) {
    throw new Error(`candidate number must be an integer from 1 through 8, got ${candidateNumber}`)
  }

  return `P${candidateNumber}`
}
