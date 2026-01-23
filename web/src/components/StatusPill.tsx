import type { WsStatus } from '../types'

export function StatusPill({ status }: { status: WsStatus }) {
  return (
    <div className="status-pill">
      <span className={`status-dot ${status === 'online' ? '' : 'offline'}`} />
      <span>WebSocket: {status}</span>
    </div>
  )
}
