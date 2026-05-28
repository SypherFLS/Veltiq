<script setup lang="ts">
definePageMeta({
  layout: 'app',
  middleware: ['require-auth'],
})
useHead({ title: 'Неликвид — Veltiq' })
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-semibold text-surface-900">Неликвид</h1>
      <p class="text-sm text-surface-500">
        Кандидаты на распродажу, скидку или списание по всем отчётам.
      </p>
    </div>
    <VEmptyState
      icon="lucide:package-search"
      title="Данных пока нет"
      description="Загрузите чековую книгу и дождитесь анализа, чтобы увидеть рекомендации."
    />
  </div>
</template>
