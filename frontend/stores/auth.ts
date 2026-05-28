import { defineStore } from 'pinia'
import type { User } from '~/shared/types/auth'

interface AuthState {
  user: User | null
  hydrated: boolean
  hydrating: boolean
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    hydrated: false,
    hydrating: false,
  }),

  getters: {
    isAuthenticated: (state) => state.user !== null,
    tenantId: (state) => state.user?.tenantId ?? null,
  },

  actions: {
    setUser(user: User | null) {
      this.user = user
    },

    clear() {
      this.user = null
    },

    async hydrate() {
      if (this.hydrated || this.hydrating) return
      this.hydrating = true
      const { $api } = useNuxtApp()
      try {
        const res = await $api.auth.session()
        if (res.valid && res.user_id) {
          if (res.access_expired) {
            try {
              await $api.auth.refresh()
            } catch {
              this.user = null
              return
            }
          }
          this.user = {
            id: res.user_id,
            tenantId: res.tenant_id ?? '',
          }
        } else {
          this.user = null
        }
      } catch {
        this.user = null
      } finally {
        this.hydrating = false
        this.hydrated = true
      }
    },

    async login(email: string, password: string) {
      const { $api } = useNuxtApp()
      await $api.auth.login({ email, password })
      this.hydrated = false
      await this.hydrate()
    },

    async register(email: string, password: string) {
      const { $api } = useNuxtApp()
      await $api.auth.register({ email, password })
    },

    async logout() {
      const { $api } = useNuxtApp()
      try {
        await $api.auth.logout()
      } finally {
        this.clear()
      }
    },
  },
})
