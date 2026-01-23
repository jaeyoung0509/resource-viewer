import { useEffect, useState } from 'react'
import type { ScaleTarget, TargetsResponse } from '../types'

const DEFAULT_MAX = 5

export function useTargets() {
  const [targets, setTargets] = useState<ScaleTarget[]>([])
  const [maxReplicas, setMaxReplicas] = useState(DEFAULT_MAX)

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

  return { targets, maxReplicas }
}
