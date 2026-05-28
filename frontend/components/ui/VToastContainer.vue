<script setup lang="ts">
import { TransitionGroup } from 'vue'
import type { ToastKind } from '~/stores/toast'

const toast = useToastStore()

const iconByKind: Record<ToastKind, string> = {
  info: 'lucide:info',
  success: 'lucide:check-circle-2',
  error: 'lucide:alert-circle',
  warning: 'lucide:alert-triangle',
}

const colorByKind: Record<ToastKind, string> = {
  info: 'border-l-sky-500 text-sky-600',
  success: 'border-l-emerald-500 text-emerald-600',
  error: 'border-l-red-500 text-red-600',
  warning: 'border-l-amber-500 text-amber-600',
}
</script>

<template>
  <ClientOnly>
    <div
      class="pointer-events-none fixed inset-x-0 top-4 z-[60] flex flex-col items-center gap-2 px-4 sm:items-end sm:right-4 sm:left-auto"
    >
      <TransitionGroup
        enter-active-class="transition duration-200 ease-out"
        enter-from-class="opacity-0 -translate-y-2"
        enter-to-class="opacity-100 translate-y-0"
        leave-active-class="transition duration-150 ease-in"
        leave-from-class="opacity-100 translate-y-0"
        leave-to-class="opacity-0 translate-y-1"
      >
        <div
          v-for="t in toast.items"
          :key="t.id"
          :class="[
            'pointer-events-auto flex w-full max-w-sm items-start gap-3 rounded-xl border-l-4 bg-white p-4 shadow-elevated',
            colorByKind[t.kind],
          ]"
        >
          <Icon :name="iconByKind[t.kind]" class="mt-0.5 size-5 shrink-0" />
          <div class="flex-1">
            <p class="text-sm font-semibold text-surface-900">{{ t.title }}</p>
            <p v-if="t.description" class="mt-0.5 text-sm text-surface-600">
              {{ t.description }}
            </p>
          </div>
          <button
            type="button"
            class="rounded-md p-1 text-surface-400 hover:text-surface-700 focus-ring"
            aria-label="Закрыть"
            @click="toast.dismiss(t.id)"
          >
            <Icon name="lucide:x" class="size-4" />
          </button>
        </div>
      </TransitionGroup>
    </div>
  </ClientOnly>
</template>
