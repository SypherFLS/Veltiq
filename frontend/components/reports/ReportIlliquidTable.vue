<script setup lang="ts">
import type { IlliquidItem } from '~/shared/types/report'
import { formatDate, formatNumber } from '~/shared/utils/format'

interface Props {
  items: IlliquidItem[]
  loading?: boolean
}

const props = defineProps<Props>()

type SortKey = 'daysWithoutSale' | 'salesQuantity' | 'name'

type SortDir = 'asc' | 'desc'

const sortKey = ref<SortKey>('daysWithoutSale')
const sortDir = ref<SortDir>('desc')
const query = ref('')
const recommendationFilter = ref<string>('all')

const recommendationOptions = [
  { value: 'all', label: 'Все рекомендации' },
  { value: 'discount', label: 'Скидка' },
  { value: 'bundle', label: 'Бандл' },
  { value: 'writeoff', label: 'Списать' },
  { value: 'monitor', label: 'Наблюдать' },
]

const recommendationLabel: Record<string, string> = {
  discount: 'Скидка',
  bundle: 'Бандл',
  writeoff: 'Списать',
  monitor: 'Наблюдать',
}

function recommendationTone(rec?: string): 'neutral' | 'warning' | 'danger' | 'info' {
  switch (rec) {
    case 'writeoff':
      return 'danger'
    case 'discount':
      return 'warning'
    case 'bundle':
      return 'info'
    default:
      return 'neutral'
  }
}

const filtered = computed(() => {
  const q = query.value.trim().toLowerCase()
  return props.items.filter((item) => {
    if (recommendationFilter.value !== 'all') {
      if ((item.recommendation ?? 'monitor') !== recommendationFilter.value)
        return false
    }
    if (!q) return true
    return (
      item.name.toLowerCase().includes(q) ||
      item.sku.toLowerCase().includes(q) ||
      (item.category?.toLowerCase().includes(q) ?? false)
    )
  })
})

const sorted = computed(() => {
  const dir = sortDir.value === 'asc' ? 1 : -1
  return [...filtered.value].sort((a, b) => {
    const k = sortKey.value
    const av = a[k] as number | string | undefined
    const bv = b[k] as number | string | undefined
    if (typeof av === 'number' && typeof bv === 'number') return (av - bv) * dir
    return String(av ?? '').localeCompare(String(bv ?? '')) * dir
  })
})

function toggleSort(k: SortKey) {
  if (sortKey.value === k) {
    sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortKey.value = k
    sortDir.value = 'desc'
  }
}

function sortIcon(k: SortKey): string {
  if (sortKey.value !== k) return 'lucide:chevrons-up-down'
  return sortDir.value === 'asc' ? 'lucide:chevron-up' : 'lucide:chevron-down'
}
</script>

<template>
  <div class="space-y-4">
    <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <VInput
        v-model="query"
        placeholder="Поиск по товару или SKU"
        icon-left="lucide:search"
        class="sm:w-72"
      />
      <select
        v-model="recommendationFilter"
        class="h-10 rounded-xl border border-surface-200 bg-white px-3 text-sm text-surface-800 focus-ring"
      >
        <option v-for="opt in recommendationOptions" :key="opt.value" :value="opt.value">
          {{ opt.label }}
        </option>
      </select>
    </div>

    <div v-if="loading" class="space-y-2">
      <VSkeleton v-for="i in 5" :key="i" height="2.5rem" />
    </div>

    <VEmptyState
      v-else-if="!sorted.length"
      icon="lucide:package-search"
      title="Ничего не найдено"
      description="Попробуйте изменить поисковый запрос или фильтр рекомендации."
    />

    <div v-else class="overflow-hidden rounded-xl border border-surface-100">
      <table class="min-w-full divide-y divide-surface-100 text-sm">
        <thead class="bg-surface-50 text-left text-xs uppercase tracking-wide text-surface-500">
          <tr>
            <th class="px-4 py-3">
              <button class="inline-flex items-center gap-1 hover:text-surface-700" @click="toggleSort('name')">
                Товар <Icon :name="sortIcon('name')" class="size-3.5" />
              </button>
            </th>
            <th class="px-4 py-3 hidden sm:table-cell">Категория</th>
            <th class="px-4 py-3 text-right">
              <button class="inline-flex items-center gap-1 hover:text-surface-700" @click="toggleSort('salesQuantity')">
                Продано за период <Icon :name="sortIcon('salesQuantity')" class="size-3.5" />
              </button>
            </th>
            <th class="px-4 py-3 text-right">
              <button class="inline-flex items-center gap-1 hover:text-surface-700" @click="toggleSort('daysWithoutSale')">
                Дней без продаж <Icon :name="sortIcon('daysWithoutSale')" class="size-3.5" />
              </button>
            </th>
            <th class="px-4 py-3 hidden md:table-cell">Последняя продажа</th>
            <th class="px-4 py-3 text-right">Рекомендация</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-surface-100 bg-white">
          <tr v-for="item in sorted" :key="item.sku" class="hover:bg-surface-50">
            <td class="px-4 py-3">
              <div class="font-medium text-surface-900">{{ item.name }}</div>
              <div class="text-xs text-surface-500">SKU {{ item.sku }}</div>
            </td>
            <td class="px-4 py-3 hidden text-surface-600 sm:table-cell">
              {{ item.category ?? '—' }}
            </td>
            <td class="px-4 py-3 text-right tabular-nums text-surface-800">
              {{ formatNumber(item.salesQuantity) }}
            </td>
            <td class="px-4 py-3 text-right tabular-nums text-surface-800">
              {{ item.daysWithoutSale }}
            </td>
            <td class="px-4 py-3 hidden text-surface-600 md:table-cell">
              {{ formatDate(item.lastSaleAt) }}
            </td>
            <td class="px-4 py-3 text-right">
              <VBadge :tone="recommendationTone(item.recommendation)">
                {{ item.recommendationNote ?? recommendationLabel[item.recommendation ?? 'monitor'] ?? '—' }}
              </VBadge>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
