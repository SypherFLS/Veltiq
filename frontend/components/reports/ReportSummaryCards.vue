<script setup lang="ts">
import { formatMoney, formatNumber } from '~/shared/utils/format'
import type { ReportSummary } from '~/shared/types/report'

interface Props {
  summary: ReportSummary | null | undefined
  loading?: boolean
}

const props = defineProps<Props>()

const cards = computed(() => [
  {
    label: 'Чеков',
    value: props.summary ? formatNumber(props.summary.receiptsCount) : '—',
    icon: 'lucide:receipt',
  },
  {
    label: 'Оборот',
    value: props.summary ? formatMoney(props.summary.totalSum) : '—',
    icon: 'lucide:wallet',
  },
  {
    label: 'Наличные',
    value: props.summary ? formatMoney(props.summary.cashSum) : '—',
    icon: 'lucide:banknote',
  },
  {
    label: 'Безнал',
    value: props.summary ? formatMoney(props.summary.cardSum) : '—',
    icon: 'lucide:credit-card',
  },
])
</script>

<template>
  <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
    <VCard v-for="c in cards" :key="c.label" padding="sm">
      <div class="flex items-start justify-between">
        <div>
          <div class="text-xs uppercase tracking-wide text-surface-500">
            {{ c.label }}
          </div>
          <div class="mt-2 text-2xl font-semibold text-surface-900">
            <VSkeleton v-if="loading" height="2rem" width="60%" />
            <template v-else>{{ c.value }}</template>
          </div>
        </div>
        <div
          class="flex size-9 items-center justify-center rounded-xl bg-brand-50 text-brand-600"
        >
          <Icon :name="c.icon" class="size-5" />
        </div>
      </div>
    </VCard>
  </div>
</template>
