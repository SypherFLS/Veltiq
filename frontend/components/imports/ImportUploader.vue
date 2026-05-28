<script setup lang="ts">
import { cn } from '~/shared/utils/classes'
import {
  ACCEPTED_IMPORT_MIME,
  MAX_IMPORT_FILE_BYTES,
  importFileSchema,
} from '~/shared/schemas/import'

interface Props {
  uploading?: boolean
  progress?: number
  error?: string | null
}

const props = withDefaults(defineProps<Props>(), {
  uploading: false,
  progress: 0,
  error: null,
})

const emit = defineEmits<{
  upload: [file: File]
  cancel: []
}>()

const file = ref<File | null>(null)
const localError = ref<string | null>(null)
const isDragging = ref(false)
const inputRef = ref<HTMLInputElement | null>(null)

const errorMessage = computed(() => props.error ?? localError.value)

const accept = computed(() =>
  [...ACCEPTED_IMPORT_MIME, '.csv', '.xlsx', '.xls', '.json', '.xml'].join(','),
)

function bytesToReadable(bytes: number): string {
  if (bytes < 1024) return `${bytes} Б`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} КБ`
  return `${(bytes / 1024 / 1024).toFixed(1)} МБ`
}

function validateAndSet(f: File | null) {
  localError.value = null
  if (!f) {
    file.value = null
    return
  }
  const result = importFileSchema.safeParse(f)
  if (!result.success) {
    localError.value = result.error.errors[0]?.message ?? 'Файл не подходит'
    file.value = null
    return
  }
  file.value = f
}

function onPick(event: Event) {
  const target = event.target as HTMLInputElement
  validateAndSet(target.files?.[0] ?? null)
}

function onDrop(event: DragEvent) {
  isDragging.value = false
  const dropped = event.dataTransfer?.files?.[0]
  if (dropped) validateAndSet(dropped)
}

function onDragOver(event: DragEvent) {
  event.preventDefault()
  if (!props.uploading) isDragging.value = true
}

function onDragLeave() {
  isDragging.value = false
}

function clear() {
  file.value = null
  localError.value = null
  if (inputRef.value) inputRef.value.value = ''
}

function submit() {
  if (!file.value || props.uploading) return
  emit('upload', file.value)
}
</script>

<template>
  <div class="space-y-4">
    <label
      :class="cn(
        'flex cursor-pointer flex-col items-center justify-center gap-3 rounded-2xl border-2 border-dashed bg-white px-6 py-10 text-center transition',
        isDragging
          ? 'border-brand-500 bg-brand-50/40'
          : 'border-surface-200 hover:border-brand-300 hover:bg-surface-50',
        (uploading || file) && 'pointer-events-none opacity-60',
      )"
      @dragover.prevent="onDragOver"
      @dragenter.prevent="onDragOver"
      @dragleave="onDragLeave"
      @drop.prevent="onDrop"
    >
      <input
        ref="inputRef"
        type="file"
        :accept="accept"
        class="sr-only"
        :disabled="uploading"
        @change="onPick"
      />
      <div class="flex size-12 items-center justify-center rounded-xl bg-brand-50 text-brand-600">
        <Icon name="lucide:upload-cloud" class="size-6" />
      </div>
      <div>
        <p class="text-sm font-medium text-surface-900">
          Перетащите файл сюда или нажмите, чтобы выбрать
        </p>
        <p class="mt-1 text-xs text-surface-500">
          CSV, XLS/XLSX, JSON или XML · до
          {{ Math.round(MAX_IMPORT_FILE_BYTES / 1024 / 1024) }} МБ
        </p>
      </div>
    </label>

    <div
      v-if="file"
      class="flex items-center justify-between gap-3 rounded-xl border border-surface-200 bg-white p-3"
    >
      <div class="flex min-w-0 items-center gap-3">
        <div class="flex size-9 items-center justify-center rounded-lg bg-surface-100 text-surface-600">
          <Icon name="lucide:file-spreadsheet" class="size-5" />
        </div>
        <div class="min-w-0">
          <p class="truncate text-sm font-medium text-surface-900">{{ file.name }}</p>
          <p class="text-xs text-surface-500">{{ bytesToReadable(file.size) }}</p>
        </div>
      </div>
      <button
        v-if="!uploading"
        type="button"
        class="rounded-lg p-1.5 text-surface-400 hover:bg-surface-100 hover:text-surface-700 focus-ring"
        aria-label="Убрать файл"
        @click="clear"
      >
        <Icon name="lucide:x" class="size-4" />
      </button>
    </div>

    <div v-if="uploading" class="space-y-1.5">
      <div class="h-2 overflow-hidden rounded-full bg-surface-100">
        <div
          class="h-full bg-brand-500 transition-all"
          :style="{ width: `${Math.min(100, Math.max(0, progress))}%` }"
        />
      </div>
      <div class="flex items-center justify-between text-xs text-surface-500">
        <span>Загружаем… {{ progress }}%</span>
        <button
          type="button"
          class="font-medium text-brand-600 hover:underline focus-ring"
          @click="emit('cancel')"
        >
          Отмена
        </button>
      </div>
    </div>

    <p v-if="errorMessage" class="field-error">{{ errorMessage }}</p>

    <div class="flex justify-end gap-2">
      <VButton
        variant="outline"
        :disabled="!file || uploading"
        @click="clear"
      >
        Очистить
      </VButton>
      <VButton
        icon="lucide:upload"
        :loading="uploading"
        :disabled="!file"
        @click="submit"
      >
        Загрузить
      </VButton>
    </div>
  </div>
</template>
