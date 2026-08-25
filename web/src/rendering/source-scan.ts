import { readdir, readFile } from 'node:fs/promises'
import { join } from 'node:path'

const sourceExtensions = new Set(['.ts', '.tsx'])
const rawColorPattern = /#[0-9a-f]{3,8}\b|\b(?:rgb|rgba|hsl|hsla|oklch|oklab)\s*\(/i
const candidateIdentifierPattern = /\bcandidate\.(?:name|title|label|theme|slug)\b/

// Color arithmetic must live in Go, never in the renderer. The named operations and the
// WCAG luminance weights are the signatures of conversion, contrast and gamut mapping, so
// a renderer that computes any of them is caught here rather than drifting from the
// authoritative implementation.
const colorMathPattern = /\b(?:relativeLuminance|luminance|contrastRatio|contrast|gamut|deltaE|deltaEOK|srgbToLinear|linearToSrgb|toOklch|toOklab|gammaEncode|gammaDecode|linearize|wcag)\b|0\.(?:2126|7152|0722)/i

export type SourceViolation = {
  file: string
  rule: 'raw-color' | 'candidate-identifier' | 'color-math'
  match: string
}

export function scanSource(source: string, file: string): SourceViolation[] {
  return [
    violationFor(source, file, 'raw-color', rawColorPattern),
    violationFor(source, file, 'candidate-identifier', candidateIdentifierPattern),
    violationFor(source, file, 'color-math', colorMathPattern),
  ].filter((violation): violation is SourceViolation => violation !== undefined)
}

function violationFor(
  source: string,
  file: string,
  rule: SourceViolation['rule'],
  pattern: RegExp,
): SourceViolation | undefined {
  const match = source.match(pattern)
  if (!match) {
    return undefined
  }

  return { file, rule, match: match[0] }
}

export async function scanRendererSources(root: string): Promise<SourceViolation[]> {
  const files = await sourceFiles(root)
  const violations = await Promise.all(
    files.map(async (file) => scanSource(await readFile(file, 'utf8'), file)),
  )

  return violations.flat()
}

async function sourceFiles(directory: string): Promise<string[]> {
  const entries = await readdir(directory, { withFileTypes: true })
  const nested = await Promise.all(
    entries.map(async (entry) => {
      const path = join(directory, entry.name)
      if (entry.isDirectory()) {
        return sourceFiles(path)
      }
      if (!sourceExtensions.has(extension(entry.name)) || entry.name.includes('.test.')) {
        return []
      }
      // This module defines the patterns which detect forbidden syntax, so scanning it would
      // only make the scanner reject its own implementation rather than renderer code.
      if (entry.name === 'source-scan.ts') {
        return []
      }
      return [path]
    }),
  )

  return nested.flat()
}

function extension(file: string): string {
  const index = file.lastIndexOf('.')
  return index === -1 ? '' : file.slice(index)
}
