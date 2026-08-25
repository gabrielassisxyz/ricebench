import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import {
  contentKinds,
  regionKinds,
  semanticStates,
  type ContentKind,
  type RegionKind,
  type SemanticState,
} from './contracts'
import { measurementStyle, paletteStyle, SceneRenderer, type FixtureSet, type Palette } from './scene'

const fixture: FixtureSet = {
  schemaVersion: '1', id: 'fixture', scenes: [{
    id: 'scene', family: 'terminal-agent',
    regions: regionKinds.map((kind, regionIndex) => ({
      id: `region-${kind}`, kind, state: semanticStates[regionIndex], background: 'background',
      foreground: 'foreground', border: 'surface',
      blocks: contentKinds.map((contentKind, blockIndex) => ({
        id: `block-${kind}-${contentKind}`, kind: contentKind, state: semanticStates[blockIndex],
        text: `${contentKind} copy`, background: 'surface', foreground: 'foreground', border: 'accent',
        runs: contentKind === 'text'
          ? semanticStates.map((state) => ({ text: state, state, foreground: 'foreground' }))
          : undefined,
      })),
    })),
  }],
}

const semanticRoleIDs = [
  'background', 'surface', 'surface-elevated', 'foreground', 'text-secondary', 'text-muted',
  'accent', 'accent-foreground', 'focus', 'cursor', 'selection-background', 'selection-foreground',
  'error', 'warning', 'success', 'info', 'diff-add', 'diff-remove', 'diff-modify',
]

const firstPalette = completePalette('first')

const snapshotFixture: FixtureSet = {
  schemaVersion: '1', id: 'snapshot', scenes: [{
    id: 'scene', family: 'terminal-agent',
    regions: [{
      id: 'frame', kind: 'frame', state: 'default', background: 'background', foreground: 'foreground', border: 'surface',
      blocks: [
        {
          id: 'copy', kind: 'text', state: 'info', background: 'surface', foreground: 'foreground',
          runs: [
            { text: 'prompt', state: 'default', foreground: 'foreground' },
            { text: 'match', state: 'search-match', foreground: 'accent' },
          ],
        },
        { id: 'badge', kind: 'badge', state: 'success', text: 'ok', background: 'success', foreground: 'accent-foreground' },
      ],
    }],
  }],
}

const secondPalette: Palette = {
  ...firstPalette,
  id: 'second',
  semanticCore: firstPalette.semanticCore.map((role) => ({
    ...role,
    value: role.value && { srgb: role.value.srgb.replace('first', 'second') },
  })),
}

describe('SceneRenderer', () => {
  describe('region primitives', () => {
    it.each(regionKinds)('renders region kind %s through the shared primitive', (kind) => {
      const markup = renderToStaticMarkup(<SceneRenderer fixtureSet={regionFixture(kind)} palette={firstPalette} />)

      expect(markup).toContain(`data-region-kind="${kind}"`)
      expect(markup).toContain(`class="region region-${kind}"`)
    })
  })

  describe('content primitives', () => {
    it.each(contentKinds)('renders content kind %s through the shared primitive', (kind) => {
      const markup = renderToStaticMarkup(<SceneRenderer fixtureSet={contentFixture(kind)} palette={firstPalette} />)

      expect(markup).toContain(`data-content-kind="${kind}"`)
      expect(markup).toContain(`class="content content-${kind}"`)
    })
  })

  it('renders every semantic state through the shared run primitive', () => {
    const markup = renderToStaticMarkup(<SceneRenderer fixtureSet={fixture} palette={firstPalette} />)

    for (const state of semanticStates) expect(markup).toContain(`data-semantic-state="${state}"`)
  })

  it('keeps structure and role-reference identifiers identical across palette swaps', () => {
    const first = renderToStaticMarkup(<SceneRenderer fixtureSet={snapshotFixture} palette={firstPalette} />)
    const second = renderToStaticMarkup(<SceneRenderer fixtureSet={snapshotFixture} palette={secondPalette} />)
    const structure = withoutPaletteValues(first)

    expect(structure).toMatchSnapshot()
    expect(withoutPaletteValues(second)).toBe(structure)
    expect(first).toContain('--palette-role-background:first-background')
    expect(second).toContain('--palette-role-background:second-background')
  })

  it('keeps all non-color measurement tokens fixed', () => {
    expect(measurementStyle()).toEqual(measurementStyle())
    expect(paletteStyle({ ...firstPalette, terminal: { ...firstPalette.terminal, aliases: [{ id: 'terminal-background', alias: { target: 'background' } }] } }))
      .toMatchObject({ '--palette-role-terminal-background': 'first-background' })
  })

  it('makes equal-copy states visibly distinct without palette changes', () => {
    const defaultRun = renderToStaticMarkup(<SceneRenderer fixtureSet={runFixture('default')} palette={firstPalette} />)
    const urgentRun = renderToStaticMarkup(<SceneRenderer fixtureSet={runFixture('urgent')} palette={firstPalette} />)

    expect(defaultRun).toContain('>·</span>same copy')
    expect(urgentRun).toContain('>‼</span>same copy')
  })
})

function regionFixture(kind: RegionKind): FixtureSet {
  return {
    schemaVersion: '1', id: `region-${kind}`, scenes: [{
      id: 'scene', family: 'terminal-agent',
      regions: [{ id: 'region', kind, state: 'default', background: 'background', foreground: 'foreground', blocks: [] }],
    }],
  }
}

function contentFixture(kind: ContentKind): FixtureSet {
  return {
    schemaVersion: '1', id: `content-${kind}`, scenes: [{
      id: 'scene', family: 'terminal-agent',
      regions: [{
        id: 'region', kind: 'frame', state: 'default', background: 'background', foreground: 'foreground',
        blocks: [{ id: 'block', kind, state: 'default', text: 'copy', background: 'surface', foreground: 'foreground' }],
      }],
    }],
  }
}

function runFixture(state: SemanticState): FixtureSet {
  return {
    schemaVersion: '1', id: `run-${state}`, scenes: [{
      id: 'scene', family: 'terminal-agent',
      regions: [{
        id: 'region', kind: 'frame', state: 'default', background: 'background', foreground: 'foreground',
        blocks: [{ id: 'block', kind: 'text', state: 'default', runs: [{ text: 'same copy', state, foreground: 'foreground' }] }],
      }],
    }],
  }
}

function completePalette(prefix: string): Palette {
  return {
    schemaVersion: '1',
    id: prefix,
    semanticCore: semanticRoleIDs.map((id) => ({ id, value: { srgb: `${prefix}-${id}` } })),
    terminal: {
      ansi: Array.from({ length: 16 }, (_, index) => ({ id: `ansi-${index}`, value: { srgb: `${prefix}-ansi-${index}` } })),
      aliases: [
        { id: 'terminal-background', alias: { target: 'background' } },
        { id: 'terminal-foreground', alias: { target: 'foreground' } },
        { id: 'terminal-cursor', alias: { target: 'cursor' } },
        { id: 'terminal-selection-background', alias: { target: 'selection-background' } },
        { id: 'terminal-selection-foreground', alias: { target: 'selection-foreground' } },
      ],
    },
  }
}

function withoutPaletteValues(markup: string): string {
  return markup.replace(/--palette-role-[a-z0-9-]+:[^;"]+/g, '')
}
