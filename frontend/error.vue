<script setup lang="ts">
interface ErrorProps {
  error: {
    statusCode: number
    statusMessage?: string
    message?: string
  }
}

const props = defineProps<ErrorProps>()

function reload() {
  clearError({ redirect: '/' })
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-surface-50 px-4">
    <div class="w-full max-w-md text-center">
      <div
        class="mx-auto flex size-14 items-center justify-center rounded-2xl bg-red-50 text-red-600"
      >
        <Icon name="lucide:alert-triangle" class="size-7" />
      </div>
      <h1 class="mt-6 text-3xl font-semibold text-surface-900">
        {{ props.error.statusCode }}
      </h1>
      <p class="mt-2 text-surface-600">
        {{ props.error.statusMessage || props.error.message || 'Что-то пошло не так' }}
      </p>
      <button
        type="button"
        class="mt-6 inline-flex h-10 items-center justify-center rounded-xl bg-brand-600 px-4 text-sm font-medium text-white hover:bg-brand-700"
        @click="reload"
      >
        На главную
      </button>
    </div>
  </div>
</template>
