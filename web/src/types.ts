export type NodeMetric = {
  node: string
  cpu: number
  memory: number
  disk: number
  timestamp: number
}

export type ScaleTarget = {
  name: string
  replicas: number
  availableReplicas: number
  error?: string
}

export type TargetsResponse = {
  maxReplicas: number
  targets: ScaleTarget[]
}

export type WsStatus = 'connecting' | 'online' | 'offline'

export type PodMetric = {
  name: string
  namespace: string
  cpuMilli: number
  memoryBytes: number
  timestamp: number
}

export type DeploymentMetric = {
  name: string
  replicas: number
  available: number
  cpuMilli: number
  memoryBytes: number
  podCount: number
  missingMetrics: number
}

export type ScaleRequest = {
  name: string
  replicas: number
}
