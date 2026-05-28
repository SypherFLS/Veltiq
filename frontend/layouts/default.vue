<template>
  <div class="min-h-full bg-surface-50">
    <header class="border-b border-surface-100 bg-white/80 backdrop-blur">
      <div class="mx-auto flex h-16 max-w-6xl items-center justify-between gap-4 px-4">
        <NuxtLink to="/" class="flex items-center gap-2">
          <span
            class="flex size-9 items-center justify-center rounded-xl bg-brand-600 text-white"
          >
            <Icon name="lucide:bar-chart-3" class="size-5" />
          </span>
          <span class="text-lg font-semibold text-surface-900">Veltiq</span>
        </NuxtLink>
        <nav class="flex items-center gap-2">
          <VButton to="/login" variant="ghost" size="sm">Войти</VButton>
          <VButton to="/register" size="sm">Регистрация</VButton>
        </nav>
      </div>
    </header>

    <main>
      <slot />
    </main>

    <footer class="mt-16 border-t border-surface-100 bg-white">
      <div
        class="mx-auto flex max-w-6xl flex-col items-start justify-between gap-2 px-4 py-6 text-sm text-surface-500 sm:flex-row sm:items-center"
      >
        <span>© {{ new Date().getFullYear() }} Veltiq</span>
        <span>SaaS-аналитика чеков для розничных сетей</span>
      </div>
    </footer>
  </div>
</template>
