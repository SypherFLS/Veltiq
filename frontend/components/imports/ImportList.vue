<script setup lang="ts">
import { formatDateTime } from '~/shared/utils/format'
import type { ImportRecord } from '~/shared/types/import'

interface Props {
  items: ImportRecord[]
  loading?: boolean
}

defineProps<Props>()
</script>

<template>
  <div v-if="loading" class="space-y-2">
    <VSkeleton v-for="i in 4" :key="i" height="3.5rem" />
  </div>

  <VEmptyState
    v-else-if="!items.length"
    icon="lucide:upload-cloud"
    title="Импортов пока нет"
    description="Загрузите чековую книгу, и здесь появится список последних импортов."
  >
    <VButton to="/app/imports/new" icon="lucide:plus">
      Загрузить файл
    </VButton>
  </VEmptyState>

  <ul v-else class="divide-y divide-surface-100">
    <li
      v-for="item in items"
      :key="item.id"
      class="flex items-center justify-between gap-3 py-3"
    >
      <div class="min-w-0">
        <NuxtLink
          :to="`/app/imports/${item.id}`"
          class="block truncate text-sm font-medium text-surface-900 hover:text-brand-600"
        >
          {{ item.fileName ?? `Импорт #${item.id}` }}
        </NuxtLink>
        <div class="text-xs text-surface-500">
          {{ formatDateTime(item.createdAt) }}
        </div>
      </div>
      <div class="flex items-center gap-3">
        <ImportStatusBadge :status="item.status" />
        <NuxtLink
          v-if="item.status === 'done'"
          :to="`/app/reports/${item.id}`"
          class="hidden text-sm font-medium text-brand-600 hover:underline sm:inline"
        >
          Отчёт
        </NuxtLink>
      </div>
    </li>
  </ul>
</template>
