import { useEffect, useState } from 'react'
import type { NodeMetric, WsStatus } from '../types'

export function useWebSocketMetrics() {
  const [nodes, setNodes] = useState<Record<string, NodeMetric>>({})
  const [status, setStatus] = useState<WsStatus>('connecting')

  useEffect(() => {
    const scheme = window.location.protocol === 'https:' ? 'wss' : 'ws'
    const ws = new WebSocket(`${scheme}://${window.location.host}/ws`)

    ws.onopen = () => setStatus('online')
    ws.onclose = () => setStatus('offline')
    ws.onerror = () => setStatus('offline')

    ws.onmessage = (event) => {
      try {
        const data: NodeMetric = JSON.parse(event.data)
        if (!data || !data.node) return
        setNodes((prev) => ({ ...prev, [data.node]: data }))
      } catch (err) {
        console.error('invalid payload', err)
      }
    }

    return () => ws.close()
  }, [])

  return { nodes, status }
}
