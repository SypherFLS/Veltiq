<script setup lang="ts">
interface Props {
  label: string
  value: string | number
  icon?: string
  hint?: string
  trend?: number
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  icon: 'lucide:activity',
})

const trendTone = computed(() => {
  if (props.trend == null) return ''
  return props.trend >= 0 ? 'text-emerald-600' : 'text-red-600'
})
</script>

<template>
  <VCard padding="sm">
    <div class="flex items-start justify-between gap-3">
      <div class="min-w-0">
        <div class="text-xs uppercase tracking-wide text-surface-500">
          {{ label }}
        </div>
        <div class="mt-2 text-2xl font-semibold text-surface-900">
          <VSkeleton v-if="loading" height="2rem" width="60%" />
          <template v-else>{{ value }}</template>
        </div>
        <div v-if="hint || trend != null" class="mt-1 flex items-center gap-2 text-xs">
          <span v-if="trend != null" :class="['inline-flex items-center gap-0.5 font-medium', trendTone]">
            <Icon
              :name="trend >= 0 ? 'lucide:trending-up' : 'lucide:trending-down'"
              class="size-3.5"
            />
            {{ trend >= 0 ? '+' : '' }}{{ trend }}%
          </span>
          <span v-if="hint" class="text-surface-500">{{ hint }}</span>
        </div>
      </div>
      <div
        class="flex size-9 shrink-0 items-center justify-center rounded-xl bg-brand-50 text-brand-600"
      >
        <Icon :name="icon" class="size-5" />
      </div>
    </div>
  </VCard>
</template>
