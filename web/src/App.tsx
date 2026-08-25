import { SceneRenderer, type FixtureSet, type Palette } from './rendering/scene'

const previewFixture: FixtureSet = {
  schemaVersion: '1',
  id: 'renderer-preview',
  scenes: [{
    id: 'preview-scene', family: 'terminal-agent',
    regions: [{
      id: 'preview-frame', kind: 'frame', state: 'default', background: 'background',
      foreground: 'foreground', border: 'surface',
      blocks: [{ id: 'preview-copy', kind: 'text', state: 'info', text: 'RiceBench', foreground: 'foreground' }],
    }],
  }],
}

const previewPalette: Palette = {
  schemaVersion: '1',
  id: 'preview',
  semanticCore: [
    { id: 'background', value: { srgb: 'Canvas' } },
    { id: 'foreground', value: { srgb: 'CanvasText' } },
    { id: 'surface', value: { srgb: 'ButtonFace' } },
  ],
  terminal: { ansi: [], aliases: [] },
}

export function App() {
  return <SceneRenderer fixtureSet={previewFixture} palette={previewPalette} />
}
