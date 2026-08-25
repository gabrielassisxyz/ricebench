import type { CSSProperties, ReactNode } from 'react'

import type { ContentKind, RegionKind, SemanticState } from './contracts'

export interface FixtureSet {
  schemaVersion: string
  id: string
  scenes: Scene[]
}

export interface Scene {
  id: string
  family: string
  regions: Region[]
}

export interface Region {
  id: string
  kind: RegionKind
  state: SemanticState
  background?: string
  foreground?: string
  border?: string
  blocks: ContentBlock[]
}

export interface ContentBlock {
  id: string
  kind: ContentKind
  state: SemanticState
  text?: string
  runs?: ContentRun[]
  background?: string
  foreground?: string
  border?: string
}

export interface ContentRun {
  text: string
  state: SemanticState
  background?: string
  foreground?: string
}

export interface PaletteRole {
  id: string
  value?: { srgb: string }
  alias?: { target: string }
}

export interface Palette {
  schemaVersion: string
  id: string
  semanticCore: PaletteRole[]
  terminal: {
    ansi: PaletteRole[]
    aliases: PaletteRole[]
  }
}

export const measurementTokens = {
  fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace',
  fontSize: '14px',
  lineHeight: '20px',
  space: '8px',
  radius: '4px',
  borderWidth: '1px',
  iconSize: '14px',
  sceneWidth: '960px',
  sceneHeight: '540px',
  zoom: '1',
  neutralBackdrop: 'Canvas',
} as const

export function measurementStyle(): CSSProperties {
  return {
    '--measurement-font-family': measurementTokens.fontFamily,
    '--measurement-font-size': measurementTokens.fontSize,
    '--measurement-line-height': measurementTokens.lineHeight,
    '--measurement-space': measurementTokens.space,
    '--measurement-radius': measurementTokens.radius,
    '--measurement-border-width': measurementTokens.borderWidth,
    '--measurement-icon-size': measurementTokens.iconSize,
    '--measurement-scene-width': measurementTokens.sceneWidth,
    '--measurement-scene-height': measurementTokens.sceneHeight,
    '--measurement-zoom': measurementTokens.zoom,
    '--measurement-neutral-backdrop': measurementTokens.neutralBackdrop,
  } as CSSProperties
}

export function paletteStyle(palette: Palette): CSSProperties {
  const roles = new Map(allPaletteRoles(palette).map((role) => [role.id, role]))
  const style: Record<string, string> = {}

  for (const role of roles.values()) {
    style[`--palette-role-${role.id}`] = resolveRole(role, roles, [])
  }

  return style as CSSProperties
}

export function SceneRenderer({ fixtureSet, palette }: { fixtureSet: FixtureSet; palette: Palette }) {
  return (
    <div
      className="fixture-set"
      data-fixture-set-id={fixtureSet.id}
      data-fixture-schema-version={fixtureSet.schemaVersion}
      style={{
        ...measurementStyle(), ...paletteStyle(palette), fontFamily: 'var(--measurement-font-family)',
        fontSize: 'var(--measurement-font-size)', lineHeight: 'var(--measurement-line-height)',
        backgroundColor: 'var(--measurement-neutral-backdrop)', transform: 'scale(var(--measurement-zoom))',
      }}
    >
      {fixtureSet.scenes.map((scene) => <ScenePrimitive key={scene.id} scene={scene} />)}
    </div>
  )
}

function ScenePrimitive({ scene }: { scene: Scene }) {
  return (
    <section
      className="scene"
      data-scene-family={scene.family}
      data-scene-id={scene.id}
      style={{ width: 'var(--measurement-scene-width)', minHeight: 'var(--measurement-scene-height)' }}
    >
      {scene.regions.map((region) => <RegionPrimitive key={region.id} region={region} />)}
    </section>
  )
}

function RegionPrimitive({ region }: { region: Region }) {
  return (
    <div
      className={`region region-${region.kind}`}
      data-region-id={region.id}
      data-region-kind={region.kind}
      data-semantic-state={region.state}
      {...roleAttributes(region)}
      style={primitiveStyle(region)}
    >
      <StateSymbol state={region.state} />
      {region.blocks.map((block) => <ContentPrimitive key={block.id} block={block} />)}
    </div>
  )
}

function ContentPrimitive({ block }: { block: ContentBlock }) {
  return (
    <div
      className={`content content-${block.kind}`}
      data-content-id={block.id}
      data-content-kind={block.kind}
      data-semantic-state={block.state}
      {...roleAttributes(block)}
      style={primitiveStyle(block)}
    >
      <StateSymbol state={block.state} />
      {block.runs?.length ? block.runs.map((run, index) => <RunPrimitive key={index} run={run} />) : block.text}
    </div>
  )
}

function RunPrimitive({ run }: { run: ContentRun }) {
  return (
    <span className="content-run" data-semantic-state={run.state} {...roleAttributes(run)} style={primitiveStyle(run)}>
      <StateSymbol state={run.state} />
      {run.text}
    </span>
  )
}

const stateSymbolStyle: CSSProperties = {
  display: 'inline-block',
  width: 'var(--measurement-icon-size)',
  fontSize: 'var(--measurement-icon-size)',
  lineHeight: 'var(--measurement-icon-size)',
  textAlign: 'center',
}

function StateSymbol({ state }: { state: SemanticState }) {
  return (
    <span aria-hidden="true" className="state-symbol" style={stateSymbolStyle}>
      {stateSymbol(state)}
    </span>
  )
}

function allPaletteRoles(palette: Palette): PaletteRole[] {
  return [...palette.semanticCore, ...palette.terminal.ansi, ...palette.terminal.aliases]
}

function resolveRole(role: PaletteRole, roles: Map<string, PaletteRole>, trail: string[]): string {
  if (role.value) return role.value.srgb
  if (!role.alias) throw new Error(`palette role ${role.id} has neither value nor alias`)
  if (trail.includes(role.id)) throw new Error(`palette role alias cycle at ${role.id}`)

  const target = roles.get(role.alias.target)
  if (!target) throw new Error(`palette role ${role.id} aliases missing role ${role.alias.target}`)
  return resolveRole(target, roles, [...trail, role.id])
}

function roleAttributes(value: { background?: string; foreground?: string; border?: string }) {
  return {
    'data-role-background': value.background,
    'data-role-foreground': value.foreground,
    'data-role-border': value.border,
  }
}

function primitiveStyle(value: { background?: string; foreground?: string; border?: string }): CSSProperties {
  return {
    backgroundColor: paletteVariable(value.background),
    color: paletteVariable(value.foreground),
    border: value.border ? 'var(--measurement-border-width) solid var(--renderer-border)' : undefined,
    borderRadius: 'var(--measurement-radius)',
    gap: 'var(--measurement-space)',
    padding: 'var(--measurement-space)',
    '--renderer-border': paletteVariable(value.border),
  } as CSSProperties
}

function paletteVariable(role: string | undefined): string | undefined {
  return role ? `var(--palette-role-${role})` : undefined
}

function stateSymbol(state: SemanticState): ReactNode {
  const symbols: Record<SemanticState, string> = {
    default: '·', active: '●', inactive: '○', focused: '◉', selected: '◆', muted: '◌',
    success: '✓', warning: '!', error: '×', info: 'i', urgent: '‼', added: '+',
    removed: '−', modified: '~', 'search-match': '⌕',
  }
  return symbols[state]
}
