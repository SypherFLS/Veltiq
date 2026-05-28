<script setup lang="ts">
import {
  Dialog,
  DialogPanel,
  DialogTitle,
  TransitionChild,
  TransitionRoot,
} from '@headlessui/vue'

interface Props {
  modelValue: boolean
  title?: string
  description?: string
  size?: 'sm' | 'md' | 'lg'
}

const props = withDefaults(defineProps<Props>(), { size: 'md' })

const emit = defineEmits<{
  'update:modelValue': [v: boolean]
}>()

const sizeClass = {
  sm: 'max-w-sm',
  md: 'max-w-md',
  lg: 'max-w-2xl',
}[props.size]

function close() {
  emit('update:modelValue', false)
}
</script>

<template>
  <TransitionRoot appear :show="modelValue" as="template">
    <Dialog as="div" class="relative z-50" @close="close">
      <TransitionChild
        as="template"
        enter="ease-out duration-200"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="ease-in duration-150"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <div class="fixed inset-0 bg-surface-900/40 backdrop-blur-sm" />
      </TransitionChild>

      <div class="fixed inset-0 overflow-y-auto">
        <div class="flex min-h-full items-center justify-center p-4">
          <TransitionChild
            as="template"
            enter="ease-out duration-200"
            enter-from="opacity-0 scale-95"
            enter-to="opacity-100 scale-100"
            leave="ease-in duration-150"
            leave-from="opacity-100 scale-100"
            leave-to="opacity-0 scale-95"
          >
            <DialogPanel
              :class="[
                'w-full transform overflow-hidden rounded-2xl bg-white p-6 text-left shadow-elevated transition-all',
                sizeClass,
              ]"
            >
              <DialogTitle
                v-if="title"
                class="text-lg font-semibold text-surface-900"
              >
                {{ title }}
              </DialogTitle>
              <p v-if="description" class="mt-1 text-sm text-surface-500">
                {{ description }}
              </p>
              <div class="mt-4">
                <slot :close="close" />
              </div>
              <div
                v-if="$slots.footer"
                class="mt-6 flex justify-end gap-2 border-t border-surface-100 pt-4"
              >
                <slot name="footer" :close="close" />
              </div>
            </DialogPanel>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>
