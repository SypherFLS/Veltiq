<script setup lang="ts">
import { cn } from '~/shared/utils/classes'

interface Props {
  modelValue?: string | number
  type?: string
  label?: string
  placeholder?: string
  name?: string
  autocomplete?: string
  disabled?: boolean
  error?: string
  hint?: string
  iconLeft?: string
  iconRight?: string
  required?: boolean
}

withDefaults(defineProps<Props>(), {
  type: 'text',
})

defineEmits<{
  'update:modelValue': [value: string]
  blur: [event: FocusEvent]
}>()

const uid = useId()
</script>

<template>
  <div>
    <label v-if="label" :for="uid" class="field-label">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>
    <div class="relative">
      <Icon
        v-if="iconLeft"
        :name="iconLeft"
        class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 size-4 text-surface-400"
      />
      <input
        :id="uid"
        :type="type"
        :name="name"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :autocomplete="autocomplete"
        :class="cn(
          'block w-full rounded-xl border bg-white text-sm text-surface-900 placeholder:text-surface-400',
          'transition focus-ring h-10',
          iconLeft ? 'pl-9' : 'pl-3.5',
          iconRight ? 'pr-9' : 'pr-3.5',
          error
            ? 'border-red-400 focus:border-red-500'
            : 'border-surface-200 focus:border-brand-500',
          disabled && 'bg-surface-50 cursor-not-allowed',
        )"
        @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
        @blur="$emit('blur', $event)"
      />
      <Icon
        v-if="iconRight"
        :name="iconRight"
        class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 size-4 text-surface-400"
      />
    </div>
    <p v-if="error" class="field-error">{{ error }}</p>
    <p v-else-if="hint" class="mt-1 text-sm text-surface-500">{{ hint }}</p>
  </div>
</template>
