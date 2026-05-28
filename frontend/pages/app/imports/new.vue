<script setup lang="ts">
import { normalizeError } from '~/shared/api/errors'

definePageMeta({
  layout: 'app',
  middleware: ['require-auth'],
})
useHead({ title: 'Новый импорт — Veltiq' })

const { $api } = useNuxtApp()
const toast = useToast()

const uploading = ref(false)
const progress = ref(0)
const uploadError = ref<string | null>(null)
const abortController = ref<AbortController | null>(null)

async function onUpload(file: File) {
  uploading.value = true
  progress.value = 0
  uploadError.value = null
  abortController.value = new AbortController()

  try {
    const res = await $api.imports.upload({
      file,
      onProgress: (p) => {
        progress.value = p
      },
      signal: abortController.value.signal,
    })
    toast.success('Файл загружен', {
      description: 'Сейчас Veltiq разбирает чековую книгу.',
    })
    await navigateTo(`/app/imports/${res.id}`)
  } catch (err) {
    if (abortController.value?.signal.aborted) {
      uploadError.value = 'Загрузка отменена'
    } else {
      const apiErr = normalizeError(err)
      uploadError.value = apiErr.message
      toast.error('Не удалось загрузить', { description: apiErr.message })
    }
  } finally {
    uploading.value = false
    abortController.value = null
  }
}

function onCancel() {
  abortController.value?.abort()
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <NuxtLink
        to="/app/imports"
        class="inline-flex items-center gap-1 text-sm text-surface-500 hover:text-surface-700"
      >
        <Icon name="lucide:arrow-left" class="size-4" />
        К списку импортов
      </NuxtLink>
      <h1 class="mt-2 text-2xl font-semibold text-surface-900">Новый импорт</h1>
      <p class="text-sm text-surface-500">
        Загрузите файл чековой книги — Veltiq распарсит его и подготовит отчёт.
      </p>
    </div>

    <VCard title="Файл чековой книги" description="Принимаются CSV, XLS/XLSX, JSON, XML">
      <ImportUploader
        :uploading="uploading"
        :progress="progress"
        :error="uploadError"
        @upload="onUpload"
        @cancel="onCancel"
      />
    </VCard>

    <VCard title="Что произойдёт дальше" padding="sm">
      <ol class="space-y-3 text-sm text-surface-600">
        <li class="flex gap-3">
          <span class="flex size-6 shrink-0 items-center justify-center rounded-full bg-brand-50 text-xs font-semibold text-brand-700">1</span>
          Файл попадает в очередь и парсится бэкендом.
        </li>
        <li class="flex gap-3">
          <span class="flex size-6 shrink-0 items-center justify-center rounded-full bg-brand-50 text-xs font-semibold text-brand-700">2</span>
          Страница импорта обновляется автоматически каждые несколько секунд.
        </li>
        <li class="flex gap-3">
          <span class="flex size-6 shrink-0 items-center justify-center rounded-full bg-brand-50 text-xs font-semibold text-brand-700">3</span>
          Когда статус станет «Готово», откроется отчёт с графиком и кандидатами на распродажу.
        </li>
      </ol>
    </VCard>
  </div>
</template>
