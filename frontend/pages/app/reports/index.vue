<script setup lang="ts">
import { formatDateTime } from '~/shared/utils/format'

definePageMeta({
  layout: 'app',
  middleware: ['require-auth'],
})
useHead({ title: 'Отчёты — Veltiq' })

const { data, isLoading, isError } = useImportsList(100)

const reports = computed(() =>
  (data.value?.items ?? []).filter(
    (item) => item.status === 'done' || item.status === 'partial_failed',
  ),
)
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-wrap items-end justify-between gap-3">
      <div>
        <h1 class="text-2xl font-semibold text-surface-900">Отчёты</h1>
        <p class="text-sm text-surface-500">
          Аналитика по импортированным чековым книгам. Каждый завершённый импорт
          даёт один отчёт.
        </p>
      </div>
      <VButton to="/app/imports/new" icon="lucide:upload-cloud" variant="outline">
        Новый импорт
      </VButton>
    </div>

    <VCard>
      <div
        v-if="isError"
        class="rounded-xl border border-amber-200 bg-amber-50 p-3 text-sm text-amber-700"
      >
        Не удалось загрузить список отчётов.
      </div>

      <div v-else-if="isLoading" class="space-y-2">
        <VSkeleton v-for="i in 4" :key="i" height="3.5rem" />
      </div>

      <VEmptyState
        v-else-if="!reports.length"
        icon="lucide:line-chart"
        title="Отчётов пока нет"
        description="Когда хотя бы один импорт перейдёт в статус «Готово», его отчёт появится здесь."
      >
        <VButton to="/app/imports/new" icon="lucide:plus">
          Загрузить файл
        </VButton>
      </VEmptyState>

      <ul v-else class="divide-y divide-surface-100">
        <li
          v-for="item in reports"
          :key="item.id"
          class="flex items-center justify-between gap-3 py-3"
        >
          <div class="min-w-0">
            <NuxtLink
              :to="`/app/reports/${item.id}`"
              class="block truncate text-sm font-medium text-surface-900 hover:text-brand-600"
            >
              Отчёт #{{ item.id }}
            </NuxtLink>
            <div class="text-xs text-surface-500">
              {{ formatDateTime(item.createdAt) }}
            </div>
          </div>
          <div class="flex items-center gap-3">
            <ImportStatusBadge :status="item.status" />
            <NuxtLink
              :to="`/app/reports/${item.id}`"
              class="hidden text-sm font-medium text-brand-600 hover:underline sm:inline"
            >
              Открыть
            </NuxtLink>
          </div>
        </li>
      </ul>
    </VCard>
  </div>
</template>
