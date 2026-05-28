export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  devtools: { enabled: true },

  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt',
    '@vueuse/nuxt',
    '@nuxt/icon',
    '@nuxt/eslint',
  ],

  css: ['~/assets/css/main.css'],

  imports: {
    dirs: ['stores', 'shared/utils', 'shared/types'],
  },

  components: [
    { path: '~/components/ui', prefix: 'V' },
    { path: '~/components', pathPrefix: false },
  ],

  routeRules: {
    '/': { prerender: true },
    '/login': { ssr: true },
    '/register': { ssr: true },
    '/app/**': { ssr: false },
  },

  runtimeConfig: {
    public: {
      apiBase: '',
    },
  },

  nitro: {
    devProxy: {
      '/api/v1': {
        target: 'http://localhost:8080/api/v1',
        changeOrigin: true,
      },
    },
  },

  icon: {
    // Явно бандлим коллекцию lucide в server runtime, чтобы /api/_nuxt_icon
    // работал и без выхода в интернет.
    serverBundle: {
      collections: ['lucide'],
    },
  },

  app: {
    head: {
      title: 'Veltiq — аналитика чековых книг',
      htmlAttrs: { lang: 'ru' },
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        {
          name: 'description',
          content:
            'SaaS для малых розничных сетей: сокращайте неликвид через аналитику чеков.',
        },
      ],
      link: [{ rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' }],
    },
  },

  typescript: {
    strict: true,
    typeCheck: false,
  },
})
