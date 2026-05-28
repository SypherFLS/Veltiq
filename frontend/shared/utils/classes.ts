export function cn(
  ...values: Array<string | false | null | undefined | Record<string, boolean>>
): string {
  const out: string[] = []
  for (const v of values) {
    if (!v) continue
    if (typeof v === 'string') {
      out.push(v)
    } else if (typeof v === 'object') {
      for (const [k, val] of Object.entries(v)) if (val) out.push(k)
    }
  }
  return out.join(' ')
}
