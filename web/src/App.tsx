import { useMemo } from 'react'
import { ControlCard } from './components/ControlCard'
import { DeploymentCard } from './components/DeploymentCard'
import { NodeCard } from './components/NodeCard'
import { PodCard } from './components/PodCard'
import { StatusPill } from './components/StatusPill'
import { useClusterMetrics } from './hooks/useClusterMetrics'
import { useTargets } from './hooks/useTargets'
import { useWebSocketMetrics } from './hooks/useWebSocketMetrics'

export default function App() {
  const { nodes, status } = useWebSocketMetrics()
  const { targets, maxReplicas } = useTargets()
  const { podMetrics, deploymentMetrics } = useClusterMetrics()

  const nodeList = useMemo(() => Object.values(nodes), [nodes])

  return (
    <div className="app">
      <header className="header">
        <div className="title">Resource Checker Control Deck</div>
        <div className="subtitle">Live node metrics + deployment scaling</div>
        <StatusPill status={status} />
      </header>

      <section className="section">
        <h2>Control Plane</h2>
        <div className="grid">
          {targets.map((target) => (
            <ControlCard key={target.name} target={target} maxReplicas={maxReplicas} />
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
            <NodeCard key={node.node} metric={node} />
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

      <section className="section">
        <h2>Deployment Metrics</h2>
        <div className="grid">
          {deploymentMetrics.map((metric) => (
            <DeploymentCard key={metric.name} metric={metric} />
          ))}
          {deploymentMetrics.length === 0 && (
            <div className="card">
              <h3>No deployment metrics</h3>
              <div className="metric">
                <span>Waiting for metrics-server</span>
                <strong>...</strong>
              </div>
            </div>
          )}
        </div>
      </section>

      <section className="section">
        <h2>Pod Metrics</h2>
        <div className="grid">
          {podMetrics.map((metric) => (
            <PodCard key={metric.name} metric={metric} />
          ))}
          {podMetrics.length === 0 && (
            <div className="card">
              <h3>No pod metrics</h3>
              <div className="metric">
                <span>Waiting for metrics-server</span>
                <strong>...</strong>
              </div>
            </div>
          )}
        </div>
      </section>
    </div>
  )
}
