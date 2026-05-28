import type { ApiClient } from './client'
import type {
  AuthStatusResponse,
  LoginPayload,
  RegisterPayload,
  SessionResponse,
} from '~/shared/types/auth'

export function createAuthApi(client: ApiClient) {
  return {
    register(payload: RegisterPayload) {
      return client.post<AuthStatusResponse>('/api/v1/register', payload)
    },

    login(payload: LoginPayload) {
      return client.post<AuthStatusResponse>('/api/v1/login', payload)
    },

    logout() {
      return client.post<AuthStatusResponse>('/api/v1/auth/logout')
    },

    refresh() {
      return client.post<AuthStatusResponse>('/api/v1/auth/refresh')
    },

    session() {
      return client.get<SessionResponse>('/api/v1/auth/session')
    },
  }
}

export type AuthApi = ReturnType<typeof createAuthApi>
