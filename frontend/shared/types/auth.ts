export interface User {
  id: string
  email?: string
  tenantId: string
}

export interface SessionResponse {
  valid: boolean
  user_id?: string
  tenant_id?: string
  access_expired?: boolean
  error?: string
}

export interface AuthStatusResponse {
  status: string
}

export interface LoginPayload {
  email: string
  password: string
}

export interface RegisterPayload {
  email: string
  password: string
}
