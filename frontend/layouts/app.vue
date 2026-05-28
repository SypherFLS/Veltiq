<script setup lang="ts">
import {
  Dialog,
  DialogPanel,
  TransitionChild,
  TransitionRoot,
} from '@headlessui/vue'

const ui = useUIStore()
const route = useRoute()

watch(() => route.fullPath, () => ui.closeMobileSidebar())
</script>

<template>
  <div class="flex h-full bg-surface-50">
    <AppSidebar />

    <TransitionRoot :show="ui.mobileSidebarOpen" as="template">
      <Dialog as="div" class="relative z-40 lg:hidden" @close="ui.closeMobileSidebar">
        <TransitionChild
          as="template"
          enter="transition-opacity ease-linear duration-200"
          enter-from="opacity-0"
          enter-to="opacity-100"
          leave="transition-opacity ease-linear duration-150"
          leave-from="opacity-100"
          leave-to="opacity-0"
        >
          <div class="fixed inset-0 bg-surface-900/40" />
        </TransitionChild>
        <div class="fixed inset-0 flex">
          <TransitionChild
            as="template"
            enter="transition ease-in-out duration-200 transform"
            enter-from="-translate-x-full"
            enter-to="translate-x-0"
            leave="transition ease-in-out duration-150 transform"
            leave-from="translate-x-0"
            leave-to="-translate-x-full"
          >
            <DialogPanel class="relative flex w-full max-w-xs">
              <div class="block w-full">
                <AppSidebar class="!flex h-full !w-full" />
              </div>
            </DialogPanel>
          </TransitionChild>
        </div>
      </Dialog>
    </TransitionRoot>

    <div class="flex min-w-0 flex-1 flex-col">
      <AppTopbar>
        <template #title-text>
          <slot name="title">Veltiq</slot>
        </template>
      </AppTopbar>

      <main class="flex-1 overflow-y-auto px-4 py-6 sm:px-6 lg:px-8">
        <div class="mx-auto max-w-6xl">
          <slot />
        </div>
      </main>
    </div>
  </div>
</template>
