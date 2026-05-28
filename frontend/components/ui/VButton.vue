<script setup lang="ts">
import { cn } from '~/shared/utils/classes'

type Variant = 'primary' | 'secondary' | 'ghost' | 'danger' | 'outline'
type Size = 'sm' | 'md' | 'lg'

interface Props {
  variant?: Variant
  size?: Size
  type?: 'button' | 'submit' | 'reset'
  loading?: boolean
  disabled?: boolean
  block?: boolean
  icon?: string
  iconRight?: string
  to?: string
  href?: string
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  type: 'button',
})

const variantClass: Record<Variant, string> = {
  primary:
    'bg-brand-600 text-white hover:bg-brand-700 active:bg-brand-800 disabled:bg-brand-300',
  secondary:
    'bg-surface-100 text-surface-800 hover:bg-surface-200 active:bg-surface-300 disabled:opacity-60',
  outline:
    'bg-white text-surface-800 border border-surface-200 hover:bg-surface-50 active:bg-surface-100 disabled:opacity-60',
  ghost:
    'bg-transparent text-surface-700 hover:bg-surface-100 active:bg-surface-200 disabled:opacity-60',
  danger:
    'bg-red-600 text-white hover:bg-red-700 active:bg-red-800 disabled:bg-red-300',
}

const sizeClass: Record<Size, string> = {
  sm: 'h-8 px-3 text-sm gap-1.5',
  md: 'h-10 px-4 text-sm gap-2',
  lg: 'h-12 px-5 text-base gap-2',
}

const classes = computed(() =>
  cn(
    'inline-flex items-center justify-center rounded-xl font-medium transition focus-ring select-none whitespace-nowrap',
    variantClass[props.variant],
    sizeClass[props.size],
    props.block && 'w-full',
    (props.disabled || props.loading) && 'pointer-events-none opacity-70',
  ),
)

const tag = computed(() => {
  if (props.to) return resolveComponent('NuxtLink')
  if (props.href) return 'a'
  return 'button'
})
</script>

<template>
  <component
    :is="tag"
    :to="to"
    :href="href"
    :type="!to && !href ? type : undefined"
    :disabled="disabled || loading"
    :class="classes"
  >
    <Icon
      v-if="loading"
      name="lucide:loader-2"
      class="size-4 animate-spin"
    />
    <Icon v-else-if="icon" :name="icon" class="size-4" />
    <span v-if="$slots.default" class="leading-none"><slot /></span>
    <Icon v-if="iconRight && !loading" :name="iconRight" class="size-4" />
  </component>
</template>
