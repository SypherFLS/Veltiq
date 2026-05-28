import type { ApiClient } from './client'

export interface ProfileResponse {
  id: string
  email?: string
  tenantId?: string
}

export function createProfileApi(client: ApiClient) {
  return {
    me() {
      return client.get<ProfileResponse>('/api/v1/profile')
    },
  }
}

export type ProfileApi = ReturnType<typeof createProfileApi>
