<script setup lang="ts">
interface Props {
  icon?: string
  title: string
  description?: string
}

withDefaults(defineProps<Props>(), {
  icon: 'lucide:inbox',
})
</script>

<template>
  <div class="flex flex-col items-center justify-center rounded-2xl border border-dashed border-surface-200 bg-white px-6 py-12 text-center">
    <div class="flex size-12 items-center justify-center rounded-xl bg-surface-100 text-surface-500">
      <Icon :name="icon" class="size-6" />
    </div>
    <h3 class="mt-4 text-base font-semibold text-surface-900">{{ title }}</h3>
    <p v-if="description" class="mt-1 max-w-sm text-sm text-surface-500">
      {{ description }}
    </p>
    <div v-if="$slots.default" class="mt-5">
      <slot />
    </div>
  </div>
</template>
