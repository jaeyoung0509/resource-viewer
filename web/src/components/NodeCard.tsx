import type { NodeMetric } from '../types'
import { formatPercent } from '../utils/format'

export function NodeCard({ metric }: { metric: NodeMetric }) {
  return (
    <div className="card">
      <h3>{metric.node}</h3>
      <div className="metric">
        <span>CPU</span>
        <strong>{formatPercent(metric.cpu)}</strong>
      </div>
      <div className="metric">
        <span>Memory</span>
        <strong>{formatPercent(metric.memory)}</strong>
      </div>
      <div className="metric">
        <span>Disk</span>
        <strong>{formatPercent(metric.disk)}</strong>
      </div>
      <div className="metric">
        <span>Last Update</span>
        <strong>
          {metric.timestamp ? new Date(metric.timestamp * 1000).toLocaleTimeString() : '-'}
        </strong>
      </div>
    </div>
  )
}
