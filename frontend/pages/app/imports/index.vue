<script setup lang="ts">
definePageMeta({
  layout: 'app',
  middleware: ['require-auth'],
})
useHead({ title: 'Импорты — Veltiq' })

const { data, isLoading, isError } = useImportsList(50)

const items = computed(() => data.value?.items ?? [])
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-wrap items-end justify-between gap-3">
      <div>
        <h1 class="text-2xl font-semibold text-surface-900">Импорты</h1>
        <p class="text-sm text-surface-500">
          Загруженные чековые книги и статус их обработки.
        </p>
      </div>
      <VButton to="/app/imports/new" icon="lucide:upload-cloud">
        Новый импорт
      </VButton>
    </div>

    <VCard>
      <div
        v-if="isError"
        class="rounded-xl border border-amber-200 bg-amber-50 p-3 text-sm text-amber-700"
      >
        Список импортов недоступен — возможно, эндпоинт ещё не реализован на бэкенде.
      </div>
      <ImportList v-else :items="items" :loading="isLoading" />
    </VCard>
  </div>
</template>
