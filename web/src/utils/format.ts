export function formatPercent(value: number) {
  if (Number.isNaN(value)) return '-'
  return `${value.toFixed(1)}%`
}

export function formatCPU(milli: number) {
  if (Number.isNaN(milli)) return '-'
  const cores = milli / 1000
  return `${cores.toFixed(2)} vCPU`
}

export function formatMemory(bytes: number) {
  if (Number.isNaN(bytes)) return '-'
  const mib = bytes / 1024 / 1024
  return `${mib.toFixed(1)} MiB`
}
