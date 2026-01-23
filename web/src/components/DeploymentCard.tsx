import type { DeploymentMetric } from '../types'
import { formatCPU, formatMemory } from '../utils/format'

export function DeploymentCard({ metric }: { metric: DeploymentMetric }) {
  return (
    <div className="card">
      <h3>{metric.name}</h3>
      <div className="metric">
        <span>Replicas</span>
        <strong>
          {metric.available} / {metric.replicas}
        </strong>
      </div>
      <div className="metric">
        <span>CPU</span>
        <strong>{formatCPU(metric.cpuMilli)}</strong>
      </div>
      <div className="metric">
        <span>Memory</span>
        <strong>{formatMemory(metric.memoryBytes)}</strong>
      </div>
      <div className="metric">
        <span>Pods</span>
        <strong>{metric.podCount}</strong>
      </div>
      {metric.missingMetrics > 0 && (
        <div className="metric">
          <span>Missing metrics</span>
          <strong>{metric.missingMetrics}</strong>
        </div>
      )}
    </div>
  )
}
