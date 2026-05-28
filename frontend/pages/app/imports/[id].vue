<script setup lang="ts">
definePageMeta({
  layout: 'app',
  middleware: ['require-auth'],
})

const route = useRoute()
const id = computed(() => route.params.id as string)

useHead({ title: () => `Импорт ${id.value} — Veltiq` })

const { data, isLoading, isError, error } = useImportPolling(id)

const isDone = computed(
  () =>
    data.value?.status === 'done' || data.value?.status === 'partial_failed',
)
const isFailed = computed(() => data.value?.status === 'failed')
const isTerminal = computed(() => isDone.value || isFailed.value)
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
      <h1 class="mt-2 text-2xl font-semibold text-surface-900">
        Импорт <span class="text-surface-400">#{{ id }}</span>
      </h1>
    </div>

    <VCard>
      <template v-if="isLoading">
        <div class="flex items-center gap-3">
          <Icon name="lucide:loader-2" class="size-5 animate-spin text-brand-500" />
          <span class="text-sm text-surface-600">Получаем статус…</span>
        </div>
      </template>

      <template v-else-if="isError">
        <div class="text-sm text-red-600">
          Не удалось получить статус: {{ (error as Error)?.message ?? 'ошибка' }}
        </div>
      </template>

      <template v-else-if="data">
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div class="flex items-center gap-3">
            <span class="text-sm text-surface-500">Статус:</span>
            <ImportStatusBadge :status="data.status" />
          </div>
          <VButton
            v-if="isDone"
            :to="`/app/reports/${id}`"
            icon-right="lucide:arrow-right"
          >
            Открыть отчёт
          </VButton>
        </div>

        <p v-if="isFailed && data.errorCode" class="mt-4 text-sm text-red-600">
          Код ошибки: {{ data.errorCode }}
        </p>
        <p v-else-if="!isTerminal" class="mt-4 text-sm text-surface-500">
          Veltiq обрабатывает файл. Эта страница обновится автоматически.
        </p>
      </template>
    </VCard>
  </div>
</template>
