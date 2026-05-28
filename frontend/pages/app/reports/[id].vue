<script setup lang="ts">
definePageMeta({
  layout: 'app',
  middleware: ['require-auth'],
})

const route = useRoute()
const id = computed(() => route.params.id as string)

useHead({ title: () => `Отчёт ${id.value} — Veltiq` })

const { data: report, isLoading: reportLoading, isError: reportError, refetch } = useReport(id)
const { data: insights, isLoading: insightsLoading } = useReportInsights(id)

const summary = computed(() => report.value?.data ?? null)
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-wrap items-end justify-between gap-3">
      <div>
        <NuxtLink
          to="/app/reports"
          class="inline-flex items-center gap-1 text-sm text-surface-500 hover:text-surface-700"
        >
          <Icon name="lucide:arrow-left" class="size-4" />
          К списку отчётов
        </NuxtLink>
        <h1 class="mt-2 text-2xl font-semibold text-surface-900">
          Отчёт <span class="text-surface-400">#{{ id }}</span>
        </h1>
        <p v-if="summary?.note" class="mt-1 text-sm text-amber-700">
          {{ summary.note }}
        </p>
      </div>
      <VButton
        variant="outline"
        icon="lucide:refresh-cw"
        :loading="reportLoading"
        @click="refetch()"
      >
        Обновить
      </VButton>
    </div>

    <div
      v-if="reportError"
      class="rounded-2xl border border-red-200 bg-red-50 p-4 text-sm text-red-700"
    >
      Не удалось загрузить отчёт. Попробуйте обновить страницу.
    </div>

    <template v-else>
      <ReportSummaryCards :summary="summary" :loading="reportLoading" />

      <VCard title="Структура продаж" description="Распределение по способам оплаты">
        <ReportSalesChart
          :cash-sum="summary?.cashSum ?? 0"
          :card-sum="summary?.cardSum ?? 0"
          :total-sum="summary?.totalSum ?? 0"
          :loading="reportLoading"
        />
      </VCard>

      <VCard
        title="Кандидаты на распродажу/списание"
        description="Товары, которые давно не продавались — фильтруйте и сортируйте"
      >
        <ReportIlliquidTable
          :items="insights?.items ?? []"
          :loading="insightsLoading"
        />
      </VCard>
    </template>
  </div>
</template>
