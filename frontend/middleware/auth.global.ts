export default defineNuxtRouteMiddleware(async (to) => {
  if (import.meta.server) return

  const auth = useAuthStore()
  if (!auth.hydrated) {
    await auth.hydrate()
  }

  const isAppRoute = to.path.startsWith('/app')
  const isGuestRoute = to.path === '/login' || to.path === '/register'

  if (isAppRoute && !auth.isAuthenticated) {
    return navigateTo({
      path: '/login',
      query: to.fullPath !== '/login' ? { redirect: to.fullPath } : undefined,
    })
  }

  if (isGuestRoute && auth.isAuthenticated) {
    return navigateTo('/app')
  }
})
