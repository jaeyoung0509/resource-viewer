import { useEffect, useState } from 'react'
import type { DeploymentMetric, PodMetric } from '../types'

export function useClusterMetrics() {
  const [podMetrics, setPodMetrics] = useState<PodMetric[]>([])
  const [deploymentMetrics, setDeploymentMetrics] = useState<DeploymentMetric[]>([])

  useEffect(() => {
    let isMounted = true

    async function fetchMetrics() {
      try {
        const [podsRes, deployRes] = await Promise.all([
          fetch('/api/metrics/pods'),
          fetch('/api/metrics/deployments')
        ])
        if (podsRes.ok) {
          const pods = (await podsRes.json()) as PodMetric[]
          if (isMounted) setPodMetrics(pods)
        }
        if (deployRes.ok) {
          const deployments = (await deployRes.json()) as DeploymentMetric[]
          if (isMounted) setDeploymentMetrics(deployments)
        }
      } catch (err) {
        console.error('failed to fetch metrics', err)
      }
    }

    fetchMetrics()
    const timer = setInterval(fetchMetrics, 4000)
    return () => {
      isMounted = false
      clearInterval(timer)
    }
  }, [])

  return { podMetrics, deploymentMetrics }
}
