export interface ApiError {
  status: number
  code?: string
  message: string
  details?: Record<string, unknown>
}

export interface Paginated<T> {
  items: T[]
  total: number
  cursor?: string | null
}

export type Result<T, E = ApiError> =
  | { ok: true; value: T }
  | { ok: false; error: E }
