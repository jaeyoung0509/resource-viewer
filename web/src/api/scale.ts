import type { ScaleRequest } from '../types'

export async function scaleDeployment(request: ScaleRequest) {
  try {
    await fetch('/api/scale', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request)
    })
  } catch (err) {
    console.error('scale failed', err)
  }
}
