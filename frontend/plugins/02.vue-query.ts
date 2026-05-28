import {
  QueryClient,
  VueQueryPlugin,
  type VueQueryPluginOptions,
} from '@tanstack/vue-query'

export default defineNuxtPlugin((nuxtApp) => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 30_000,
        retry: 1,
        refetchOnWindowFocus: false,
      },
    },
  })

  const options: VueQueryPluginOptions = { queryClient }
  nuxtApp.vueApp.use(VueQueryPlugin, options)

  return {
    provide: {
      queryClient,
    },
  }
})
