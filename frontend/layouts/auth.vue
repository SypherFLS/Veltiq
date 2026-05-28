<template>
  <div class="flex min-h-full items-center justify-center bg-surface-50 px-4 py-12">
    <div class="w-full max-w-md">
      <NuxtLink to="/" class="mb-8 flex items-center justify-center gap-2">
        <span
          class="flex size-10 items-center justify-center rounded-xl bg-brand-600 text-white"
        >
          <Icon name="lucide:bar-chart-3" class="size-5" />
        </span>
        <span class="text-xl font-semibold text-surface-900">Veltiq</span>
      </NuxtLink>

      <div class="card-base p-6 sm:p-8">
        <slot />
      </div>
    </div>
  </div>
</template>
