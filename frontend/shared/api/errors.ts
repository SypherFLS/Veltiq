import type { ApiError } from '~/shared/types/api'

interface FetchLikeError {
  response?: { status?: number; _data?: unknown }
  data?: unknown
  statusCode?: number
  message?: string
}

export function normalizeError(err: unknown): ApiError {
  const e = err as FetchLikeError
  const status = e?.response?.status ?? e?.statusCode ?? 0
  const payload = (e?.response?._data ?? e?.data) as
    | Record<string, unknown>
    | undefined

  const message =
    (payload?.message as string | undefined) ??
    (payload?.error as string | undefined) ??
    e?.message ??
    'Неизвестная ошибка'

  const code = payload?.code as string | undefined

  const details = payload && typeof payload === 'object'
    ? (payload as Record<string, unknown>)
    : undefined

  return { status, code, message, details }
}

export function isApiError(value: unknown): value is ApiError {
  return (
    typeof value === 'object' &&
    value !== null &&
    'status' in value &&
    'message' in value
  )
}
