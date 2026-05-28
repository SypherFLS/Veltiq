<script setup lang="ts">
import { formatMoney, formatNumber } from '~/shared/utils/format'

definePageMeta({
  layout: 'app',
  middleware: ['require-auth'],
})
useHead({ title: 'Дашборд — Veltiq' })

const { user } = useAuth()
const { data: importsData, isLoading: importsLoading, isError: importsError } = useImportsList(5)

const items = computed(() => importsData.value?.items ?? [])

const aggregate = computed(() => {
  const done = items.value.filter((i) => i.status === 'done')
  const inFlight = items.value.filter(
    (i) => i.status === 'pending' || i.status === 'processing',
  )
  const failed = items.value.filter((i) => i.status === 'failed')
  return {
    total: importsData.value?.total ?? items.value.length,
    done: done.length,
    inFlight: inFlight.length,
    failed: failed.length,
  }
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-wrap items-end justify-between gap-3">
      <div>
        <h1 class="text-2xl font-semibold text-surface-900">Дашборд</h1>
        <p class="text-sm text-surface-500">
          Привет{{ user?.email ? `, ${user.email}` : '' }}. Загрузите чековую книгу,
          чтобы получить аналитику.
        </p>
      </div>
      <VButton to="/app/imports/new" icon="lucide:upload-cloud">
        Новый импорт
      </VButton>
    </div>

    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <KpiCard
        label="Всего импортов"
        :value="formatNumber(aggregate.total)"
        icon="lucide:database"
        :loading="importsLoading"
      />
      <KpiCard
        label="Обработано"
        :value="formatNumber(aggregate.done)"
        icon="lucide:check-circle-2"
        hint="готовых отчётов"
        :loading="importsLoading"
      />
      <KpiCard
        label="В работе"
        :value="formatNumber(aggregate.inFlight)"
        icon="lucide:loader"
        hint="в очереди или обрабатываются"
        :loading="importsLoading"
      />
      <KpiCard
        label="С ошибками"
        :value="formatNumber(aggregate.failed)"
        icon="lucide:alert-triangle"
        :loading="importsLoading"
      />
    </div>

    <VCard
      title="Последние импорты"
      description="Файлы чековых книг, загруженные за последнее время"
    >
      <template #actions>
        <VButton variant="ghost" size="sm" to="/app/imports" icon-right="lucide:arrow-right">
          Все импорты
        </VButton>
      </template>

      <div
        v-if="importsError"
        class="rounded-xl border border-amber-200 bg-amber-50 p-3 text-sm text-amber-700"
      >
        Список импортов недоступен — возможно, эндпоинт ещё не реализован на бэкенде.
      </div>
      <ImportList v-else :items="items" :loading="importsLoading" />
    </VCard>
  </div>
</template>
