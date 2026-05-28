export default defineNuxtRouteMiddleware(() => {
  if (import.meta.server) return
  const auth = useAuthStore()
  if (auth.isAuthenticated) {
    return navigateTo('/app')
  }
})
