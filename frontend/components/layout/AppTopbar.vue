<script setup lang="ts">
const ui = useUIStore()
</script>

<template>
  <header
    class="flex h-16 items-center gap-3 border-b border-surface-100 bg-white px-4 sm:px-6"
  >
    <button
      type="button"
      class="rounded-xl border border-surface-200 p-2 text-surface-700 lg:hidden focus-ring"
      aria-label="Открыть меню"
      @click="ui.openMobileSidebar"
    >
      <Icon name="lucide:menu" class="size-5" />
    </button>

    <div class="flex-1">
      <slot name="title">
        <h1 class="text-lg font-semibold text-surface-900">
          <slot name="title-text" />
        </h1>
      </slot>
    </div>

    <slot name="actions" />
    <UserMenu />
  </header>
</template>
