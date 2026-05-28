<script setup lang="ts">
import type { ImportStatus } from '~/shared/types/import'

interface Props {
  status: ImportStatus
}

defineProps<Props>()

const labelByStatus: Record<ImportStatus, string> = {
  pending: 'В очереди',
  processing: 'Обработка',
  partial_failed: 'Частично',
  done: 'Готово',
  failed: 'Ошибка',
}

const toneByStatus: Record<
  ImportStatus,
  'neutral' | 'info' | 'success' | 'danger' | 'warning'
> = {
  pending: 'neutral',
  processing: 'info',
  partial_failed: 'warning',
  done: 'success',
  failed: 'danger',
}

const iconByStatus: Record<ImportStatus, string> = {
  pending: 'lucide:clock',
  processing: 'lucide:loader-2',
  partial_failed: 'lucide:alert-triangle',
  done: 'lucide:check-circle-2',
  failed: 'lucide:alert-circle',
}
</script>

<template>
  <VBadge :tone="toneByStatus[status]" :icon="iconByStatus[status]">
    <span :class="status === 'processing' && 'animate-pulse'">
      {{ labelByStatus[status] }}
    </span>
  </VBadge>
</template>
