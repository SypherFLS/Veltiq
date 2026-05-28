<script setup lang="ts">
type Tone = 'neutral' | 'brand' | 'success' | 'warning' | 'danger' | 'info'

interface Props {
  tone?: Tone
  icon?: string
}

const props = withDefaults(defineProps<Props>(), { tone: 'neutral' })

const toneClass: Record<Tone, string> = {
  neutral: 'bg-surface-100 text-surface-700',
  brand: 'bg-brand-50 text-brand-700',
  success: 'bg-emerald-50 text-emerald-700',
  warning: 'bg-amber-50 text-amber-700',
  danger: 'bg-red-50 text-red-700',
  info: 'bg-sky-50 text-sky-700',
}
</script>

<template>
  <span
    :class="[
      'inline-flex items-center gap-1 rounded-full px-2.5 py-1 text-xs font-medium',
      toneClass[props.tone],
    ]"
  >
    <Icon v-if="icon" :name="icon" class="size-3.5" />
    <slot />
  </span>
</template>
