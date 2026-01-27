import type { ScaleTarget } from '../types'
import { scaleDeployment } from '../api/scale'

export function ControlCard({
  target,
  maxReplicas
}: {
  target: ScaleTarget
  maxReplicas: number
}) {
  return (
    <div className="card">
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
              onClick={() =>
                scaleDeployment({
                  name: target.name,
                  replicas: Math.max(0, target.replicas - 1)
                })
              }
              disabled={target.replicas <= 0}
            >
              -1
            </button>
            <button
              className="secondary"
              onClick={() =>
                scaleDeployment({
                  name: target.name,
                  replicas: Math.min(maxReplicas, target.replicas + 1)
                })
              }
              disabled={target.replicas >= maxReplicas}
            >
              +1
            </button>
          </div>
          <button
            onClick={() => scaleDeployment({ name: target.name, replicas: 0 })}
            disabled={target.replicas === 0}
          >
            Stop
          </button>
        </div>
      </div>
    </div>
  )
}
