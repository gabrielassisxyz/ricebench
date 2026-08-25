import { mkdir, mkdtemp, readFile, rm, writeFile } from 'node:fs/promises'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { fileURLToPath } from 'node:url'
import { afterEach, describe, expect, it } from 'vitest'

import {
  contentKinds,
  formatCandidateLabel,
  regionKinds,
  semanticStates,
} from './contracts'
import { scanRendererSources, scanSource } from './source-scan'

const temporaryDirectories: string[] = []

describe('fixture rendering contract', () => {
  it('covers every schema-v1 renderer vocabulary', async () => {
    const schema = fileURLToPath(new URL('../../../internal/fixture/schema.go', import.meta.url))

    expect(regionKinds).toEqual(await schemaValues(schema, 'RegionKind'))
    expect(contentKinds).toEqual(await schemaValues(schema, 'ContentKind'))
    expect(semanticStates).toEqual(await schemaValues(schema, 'SemanticState'))
    expect({ contentKinds, regionKinds, semanticStates }).toMatchSnapshot()
  })

  it('exposes only neutral participant-facing candidate labels', () => {
    expect(formatCandidateLabel(1)).toBe('P1')
    expect(formatCandidateLabel(8)).toBe('P8')
    expect(() => formatCandidateLabel(0)).toThrow('candidate number must be an integer')
    expect(() => formatCandidateLabel(1.5)).toThrow('candidate number must be an integer')
  })
})

describe('renderer source conformance', () => {
  afterEach(async () => {
    await Promise.all(temporaryDirectories.splice(0).map((directory) => rm(directory, { recursive: true })))
  })

  it('rejects raw colors', () => {
    expect(scanSource('const foreground = "#123456"', 'example.ts')).toEqual([
      { file: 'example.ts', rule: 'raw-color', match: '#123456' },
    ])
  })

  it('rejects internal candidate identifiers', () => {
    expect(scanSource('return candidate.theme', 'example.ts')).toEqual([
      { file: 'example.ts', rule: 'candidate-identifier', match: 'candidate.theme' },
    ])
  })

  it('reports raw colors in recursively scanned source files', async () => {
    const root = await rendererSourceDirectory()
    const file = join(root, 'nested', 'example.ts')
    await writeFile(file, 'const foreground = "#123456"')

    await expect(scanRendererSources(root)).resolves.toEqual([
      { file, rule: 'raw-color', match: '#123456' },
    ])
  })

  it('reports internal candidate identifiers in scanned source files', async () => {
    const root = await rendererSourceDirectory()
    const file = join(root, 'example.tsx')
    await writeFile(file, 'export const label = candidate.theme')

    await expect(scanRendererSources(root)).resolves.toEqual([
      { file, rule: 'candidate-identifier', match: 'candidate.theme' },
    ])
  })

  it('keeps all frontend sources free of forbidden syntax', async () => {
    const root = fileURLToPath(new URL('../', import.meta.url))
    await expect(scanRendererSources(root)).resolves.toEqual([])
  })
})

async function schemaValues(schemaPath: string, type: string): Promise<string[]> {
  const schema = await readFile(schemaPath, 'utf8')
  const declaration = new RegExp(`\\w+\\s+${type}\\s+=\\s+"([^"]+)"`, 'g')

  return Array.from(schema.matchAll(declaration), (match) => match[1])
}

async function rendererSourceDirectory(): Promise<string> {
  const root = await mkdtemp(join(tmpdir(), 'ricebench-renderer-'))
  const nested = join(root, 'nested')
  await mkdir(nested)
  temporaryDirectories.push(root)
  return root
}
