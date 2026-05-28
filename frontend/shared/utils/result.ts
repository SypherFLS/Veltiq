import type { ApiError, Result } from '~/shared/types/api'
import { normalizeError } from '~/shared/api/errors'

export async function safeCall<T>(
  fn: () => Promise<T>,
): Promise<Result<T, ApiError>> {
  try {
    const value = await fn()
    return { ok: true, value }
  } catch (err) {
    return { ok: false, error: normalizeError(err) }
  }
}
