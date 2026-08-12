import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// The Go binary embeds internal/web/dist, so the build has to land there rather than in the
// default web/dist. emptyOutDir stays off because that directory carries a tracked .gitkeep
// which go:embed depends on; bin/build-web clears the previous build instead, so hashed
// assets cannot accumulate inside the binary.
export default defineConfig({
  plugins: [react()],
  build: {
    outDir: '../internal/web/dist',
    emptyOutDir: false,
  },
})
