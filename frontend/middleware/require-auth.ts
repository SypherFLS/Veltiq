export default defineNuxtRouteMiddleware((to) => {
  if (import.meta.server) return
  const auth = useAuthStore()
  if (!auth.isAuthenticated) {
    return navigateTo({
      path: '/login',
      query: { redirect: to.fullPath },
    })
  }
})
