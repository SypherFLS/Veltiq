import { ofetch, type FetchOptions } from 'ofetch'
import {
  createApiClient,
  createVeltiqApi,
  type RawFetch,
  type VeltiqApi,
} from '~/shared/api'

declare module '#app' {
  interface NuxtApp {
    $api: VeltiqApi
  }
}

declare module 'vue' {
  interface ComponentCustomProperties {
    $api: VeltiqApi
  }
}

const NO_REFRESH_PATHS = [
  '/api/v1/auth/refresh',
  '/api/v1/auth/logout',
  '/api/v1/auth/session',
  '/api/v1/login',
  '/api/v1/register',
]

function shouldSkipRefresh(url: string): boolean {
  return NO_REFRESH_PATHS.some((p) => url.endsWith(p))
}

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase

  const base = ofetch.create({
    baseURL,
    credentials: 'include',
    retry: 0,
  })

  let refreshInFlight: Promise<void> | null = null

  async function ensureRefresh(): Promise<boolean> {
    if (!refreshInFlight) {
      refreshInFlight = base('/api/v1/auth/refresh', { method: 'POST' })
        .then(() => undefined)
        .finally(() => {
          refreshInFlight = null
        })
    }
    try {
      await refreshInFlight
      return true
    } catch {
      return false
    }
  }

  const withAuthRetry: RawFetch = async <T = unknown>(
    url: string,
    options: FetchOptions<'json'> = {},
  ): Promise<T> => {
    try {
      return (await base(url, options)) as T
    } catch (err) {
      const status = (err as { response?: { status?: number } })?.response?.status

      if (status !== 401 || shouldSkipRefresh(url)) {
        throw err
      }

      const refreshed = await ensureRefresh()
      if (!refreshed) {
        if (import.meta.client) {
          const auth = useAuthStore()
          auth.clear()
          const route = useRoute()
          if (route.path.startsWith('/app')) {
            await navigateTo('/login')
          }
        }
        throw err
      }

      return (await base(url, options)) as T
    }
  }

  const client = createApiClient(withAuthRetry)
  const api = createVeltiqApi(client, baseURL)

  return {
    provide: {
      api,
    },
  }
})
