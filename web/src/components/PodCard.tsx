import type { PodMetric } from '../types'
import { formatCPU, formatMemory } from '../utils/format'

export function PodCard({ metric }: { metric: PodMetric }) {
  return (
    <div className="card">
      <h3>{metric.name}</h3>
      <div className="metric">
        <span>CPU</span>
        <strong>{formatCPU(metric.cpuMilli)}</strong>
      </div>
      <div className="metric">
        <span>Memory</span>
        <strong>{formatMemory(metric.memoryBytes)}</strong>
      </div>
      <div className="metric">
        <span>Namespace</span>
        <strong>{metric.namespace}</strong>
      </div>
    </div>
  )
}
