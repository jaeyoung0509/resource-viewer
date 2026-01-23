import { useEffect, useMemo, useState } from 'react'

const DEFAULT_MAX = 5

type NodeMetric = {
  node: string
  cpu: number
  memory: number
  disk: number
  timestamp: number
}

type ScaleTarget = {
  name: string
  replicas: number
  availableReplicas: number
}

type TargetsResponse = {
  maxReplicas: number
  targets: ScaleTarget[]
}

type WsStatus = 'connecting' | 'online' | 'offline'

function formatPercent(value: number) {
  if (Number.isNaN(value)) return '-'
  return `${value.toFixed(1)}%`
}

export default function App() {
  const [nodes, setNodes] = useState<Record<string, NodeMetric>>({})
  const [targets, setTargets] = useState<ScaleTarget[]>([])
  const [maxReplicas, setMaxReplicas] = useState(DEFAULT_MAX)
  const [wsStatus, setWsStatus] = useState<WsStatus>('connecting')

  const nodeList = useMemo(() => Object.values(nodes), [nodes])

  useEffect(() => {
    const scheme = window.location.protocol === 'https:' ? 'wss' : 'ws'
    const ws = new WebSocket(`${scheme}://${window.location.host}/ws`)

    ws.onopen = () => setWsStatus('online')
    ws.onclose = () => setWsStatus('offline')
    ws.onerror = () => setWsStatus('offline')

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

  useEffect(() => {
    let isMounted = true
    async function fetchTargets() {
      try {
        const res = await fetch('/api/targets')
        if (!res.ok) return
        const data = (await res.json()) as TargetsResponse
        if (!isMounted) return
        setTargets(data.targets || [])
        setMaxReplicas(data.maxReplicas || DEFAULT_MAX)
      } catch (err) {
        console.error('failed to fetch targets', err)
      }
    }

    fetchTargets()
    const timer = setInterval(fetchTargets, 2000)
    return () => {
      isMounted = false
      clearInterval(timer)
    }
  }, [])

  async function scale(name: string, replicas: number) {
    try {
      await fetch('/api/scale', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, replicas })
      })
    } catch (err) {
      console.error('scale failed', err)
    }
  }

  return (
    <div className="app">
      <header className="header">
        <div className="title">Resource Checker Control Deck</div>
        <div className="subtitle">Live node metrics + deployment scaling</div>
        <div className="status-pill">
          <span className={`status-dot ${wsStatus === 'online' ? '' : 'offline'}`} />
          <span>WebSocket: {wsStatus}</span>
        </div>
      </header>

      <section className="section">
        <h2>Control Plane</h2>
        <div className="grid">
          {targets.map((target) => (
            <div key={target.name} className="card">
              <h3>{target.name}</h3>
              <div className="metric">
                <span>Replicas</span>
                <strong>
                  {target.replicas} / {maxReplicas}
                </strong>
              </div>
              <div className="metric">
                <span>Available</span>
                <strong>{target.availableReplicas}</strong>
              </div>
              <div className="controls">
                <div className="control-row">
                  <div className="control-buttons">
                    <button
                      onClick={() => scale(target.name, Math.max(0, target.replicas - 1))}
                      disabled={target.replicas <= 0}
                    >
                      -1
                    </button>
                    <button
                      className="secondary"
                      onClick={() => scale(target.name, Math.min(maxReplicas, target.replicas + 1))}
                      disabled={target.replicas >= maxReplicas}
                    >
                      +1
                    </button>
                  </div>
                  <button
                    onClick={() => scale(target.name, 0)}
                    disabled={target.replicas === 0}
                  >
                    Stop
                  </button>
                </div>
              </div>
            </div>
          ))}
          {targets.length === 0 && (
            <div className="card">
              <h3>No targets</h3>
              <div className="metric">
                <span>Waiting for /api/targets</span>
                <strong>...</strong>
              </div>
            </div>
          )}
        </div>
      </section>

      <section className="section">
        <h2>Live Metrics</h2>
        <div className="grid">
          {nodeList.map((node) => (
            <div key={node.node} className="card">
              <h3>{node.node}</h3>
              <div className="metric">
                <span>CPU</span>
                <strong>{formatPercent(node.cpu)}</strong>
              </div>
              <div className="metric">
                <span>Memory</span>
                <strong>{formatPercent(node.memory)}</strong>
              </div>
              <div className="metric">
                <span>Disk</span>
                <strong>{formatPercent(node.disk)}</strong>
              </div>
              <div className="metric">
                <span>Last Update</span>
                <strong>
                  {node.timestamp ? new Date(node.timestamp * 1000).toLocaleTimeString() : '-'}
                </strong>
              </div>
            </div>
          ))}
          {nodeList.length === 0 && (
            <div className="card">
              <h3>No nodes yet</h3>
              <div className="metric">
                <span>Waiting for metrics</span>
                <strong>...</strong>
              </div>
            </div>
          )}
        </div>
      </section>
    </div>
  )
}
