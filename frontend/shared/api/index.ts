import type { ApiClient } from './client'
import { createAuthApi, type AuthApi } from './auth'
import { createImportsApi, type ImportsApi } from './imports'
import { createProfileApi, type ProfileApi } from './profile'

export * from './client'
export * from './errors'
export type { AuthApi } from './auth'
export type { ImportsApi } from './imports'
export type { ProfileApi } from './profile'

export interface VeltiqApi {
  auth: AuthApi
  imports: ImportsApi
  profile: ProfileApi
}

export function createVeltiqApi(client: ApiClient, baseURL: string): VeltiqApi {
  return {
    auth: createAuthApi(client),
    imports: createImportsApi(client, baseURL),
    profile: createProfileApi(client),
  }
}
