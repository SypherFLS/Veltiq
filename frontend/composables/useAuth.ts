import { normalizeError } from '~/shared/api/errors'
import type { ApiError } from '~/shared/types/api'

export function useAuth() {
  const store = useAuthStore()
  const toast = useToast()

  async function login(email: string, password: string): Promise<ApiError | null> {
    try {
      await store.login(email, password)
      return null
    } catch (err) {
      const apiErr = normalizeError(err)
      toast.error('Не удалось войти', { description: apiErr.message })
      return apiErr
    }
  }

  async function register(email: string, password: string): Promise<ApiError | null> {
    try {
      await store.register(email, password)
      return null
    } catch (err) {
      const apiErr = normalizeError(err)
      toast.error('Не удалось зарегистрироваться', {
        description: apiErr.message,
      })
      return apiErr
    }
  }

  async function logout() {
    await store.logout()
    const { $queryClient } = useNuxtApp()
    $queryClient?.clear()
    await navigateTo('/login')
  }

  return {
    user: computed(() => store.user),
    isAuthenticated: computed(() => store.isAuthenticated),
    tenantId: computed(() => store.tenantId),
    login,
    register,
    logout,
  }
}
