<script setup lang="ts">
interface Props {
  width?: string
  height?: string
  rounded?: 'sm' | 'md' | 'lg' | 'xl' | 'full'
}

const props = withDefaults(defineProps<Props>(), {
  height: '1rem',
  rounded: 'md',
})

const roundedClass = {
  sm: 'rounded',
  md: 'rounded-md',
  lg: 'rounded-lg',
  xl: 'rounded-xl',
  full: 'rounded-full',
}[props.rounded]
</script>

<template>
  <div
    :class="['animate-pulse bg-surface-200/70', roundedClass]"
    :style="{ width: width ?? '100%', height }"
  />
</template>
