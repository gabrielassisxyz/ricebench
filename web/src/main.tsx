import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'

import { App } from './App'

const container = document.getElementById('root')
if (!container) {
  throw new Error('missing #root: index.html and this entry point disagree')
}

createRoot(container).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
