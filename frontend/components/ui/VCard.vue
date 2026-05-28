<script setup lang="ts">
interface Props {
  title?: string
  description?: string
  padding?: 'sm' | 'md' | 'lg' | 'none'
}

const props = withDefaults(defineProps<Props>(), {
  padding: 'md',
})

const paddingClass = {
  none: '',
  sm: 'p-4',
  md: 'p-5 sm:p-6',
  lg: 'p-6 sm:p-8',
}[props.padding]
</script>

<template>
  <section :class="['card-base', paddingClass]">
    <header v-if="title || $slots.header" class="mb-4 flex items-start justify-between gap-3">
      <div>
        <h3 v-if="title" class="text-base font-semibold text-surface-900">
          {{ title }}
        </h3>
        <p v-if="description" class="mt-1 text-sm text-surface-500">
          {{ description }}
        </p>
      </div>
      <slot name="actions" />
    </header>
    <slot />
  </section>
</template>
