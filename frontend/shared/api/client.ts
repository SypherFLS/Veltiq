import type { FetchOptions } from 'ofetch'

export type RawFetch = <T = unknown>(
  url: string,
  options?: FetchOptions<'json'>,
) => Promise<T>

export interface ApiClient {
  raw: RawFetch
  get<T = unknown>(url: string, opts?: FetchOptions<'json'>): Promise<T>
  post<T = unknown>(url: string, body?: unknown, opts?: FetchOptions<'json'>): Promise<T>
  put<T = unknown>(url: string, body?: unknown, opts?: FetchOptions<'json'>): Promise<T>
  patch<T = unknown>(url: string, body?: unknown, opts?: FetchOptions<'json'>): Promise<T>
  delete<T = unknown>(url: string, opts?: FetchOptions<'json'>): Promise<T>
}

export function createApiClient(raw: RawFetch): ApiClient {
  return {
    raw,
    get: (url, opts) => raw(url, { method: 'GET', ...opts }),
    post: (url, body, opts) => raw(url, { method: 'POST', body, ...opts }),
    put: (url, body, opts) => raw(url, { method: 'PUT', body, ...opts }),
    patch: (url, body, opts) => raw(url, { method: 'PATCH', body, ...opts }),
    delete: (url, opts) => raw(url, { method: 'DELETE', ...opts }),
  }
}
