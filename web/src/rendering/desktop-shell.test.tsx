import { readFile } from 'node:fs/promises'
import { fileURLToPath } from 'node:url'
import { renderToStaticMarkup } from 'react-dom/server'
import { beforeAll, describe, expect, it } from 'vitest'

import { SceneRenderer, type FixtureSet, type Palette } from './scene'

const goldenPath = fileURLToPath(
  new URL('../../../internal/fixture/testdata/desktop-shell.golden.json', import.meta.url),
)

const semanticRoleIDs = [
  'background', 'surface', 'surface-elevated', 'foreground', 'text-secondary', 'text-muted',
  'accent', 'accent-foreground', 'focus', 'cursor', 'selection-background', 'selection-foreground',
  'error', 'warning', 'success', 'info', 'diff-add', 'diff-remove', 'diff-modify',
]

// Every role resolves to the same gray, so any remaining distinction between states is
// structural rather than color.
const grayscalePalette: Palette = {
  schemaVersion: '1',
  id: 'grayscale',
  semanticCore: semanticRoleIDs.map((id) => ({ id, value: { srgb: '#808080' } })),
  terminal: { ansi: [], aliases: [] },
}

let family: FixtureSet

beforeAll(async () => {
  family = JSON.parse(await readFile(goldenPath, 'utf8')) as FixtureSet
})

describe('desktop-shell family', () => {
  it('keeps focused and unfocused windows distinguishable under a grayscale role map', () => {
    const markup = render(family)

    expect(markup).toContain('data-content-kind="focus-ring"')
    expect(markup.match(/data-content-kind="focus-ring"/g)).toHaveLength(1)
    expect(markup).toContain('data-semantic-state="focused"')
    expect(markup).toContain('data-semantic-state="inactive"')
  })

  it('keeps active and inactive workspaces distinguishable under a grayscale role map', () => {
    const markup = render(family)

    expect(markup).toContain('data-semantic-state="active"')
    expect(markup).toContain('>●</span>')
    expect(markup).toContain('data-semantic-state="inactive"')
    expect(markup).toContain('>○</span>')
  })

  it('keeps selected and neighboring launcher results distinguishable under a grayscale role map', () => {
    const markup = render(family)

    expect(markup).toContain('data-semantic-state="selected"')
    expect(markup).toContain('>◆</span>')
    expect(markup).toContain('data-semantic-state="default"')
    expect(markup).toContain('>·</span>')
  })

  it('snapshots the rendered structure with role references stripped', () => {
    expect(withoutPaletteValues(render(family))).toMatchSnapshot()
  })
})

function render(fixtureSet: FixtureSet): string {
  return renderToStaticMarkup(<SceneRenderer fixtureSet={fixtureSet} palette={grayscalePalette} />)
}

function withoutPaletteValues(markup: string): string {
  return markup.replace(/--palette-role-[a-z0-9-]+:[^;"]+/g, '')
}
