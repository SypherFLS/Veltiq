<script setup lang="ts">
interface NavItem {
  to: string
  label: string
  icon: string
}

const navItems: NavItem[] = [
  { to: '/app', label: 'Дашборд', icon: 'lucide:layout-dashboard' },
  { to: '/app/imports', label: 'Импорты', icon: 'lucide:upload-cloud' },
  { to: '/app/reports', label: 'Отчёты', icon: 'lucide:line-chart' },
  { to: '/app/inventory', label: 'Неликвид', icon: 'lucide:package-search' },
  { to: '/app/settings/profile', label: 'Настройки', icon: 'lucide:settings' },
]

const route = useRoute()

function isActive(to: string): boolean {
  if (to === '/app') return route.path === '/app'
  return route.path.startsWith(to)
}
</script>

<template>
  <aside class="hidden h-full w-64 shrink-0 flex-col border-r border-surface-100 bg-white lg:flex">
    <div class="flex h-16 items-center gap-2 border-b border-surface-100 px-5">
      <div class="flex size-9 items-center justify-center rounded-xl bg-brand-600 text-white">
        <Icon name="lucide:bar-chart-3" class="size-5" />
      </div>
      <span class="text-lg font-semibold text-surface-900">Veltiq</span>
    </div>

    <nav class="flex-1 overflow-y-auto p-3">
      <ul class="space-y-1">
        <li v-for="item in navItems" :key="item.to">
          <NuxtLink
            :to="item.to"
            :class="[
              'flex items-center gap-3 rounded-xl px-3 py-2 text-sm font-medium transition focus-ring',
              isActive(item.to)
                ? 'bg-brand-50 text-brand-700'
                : 'text-surface-700 hover:bg-surface-50 hover:text-surface-900',
            ]"
          >
            <Icon :name="item.icon" class="size-4" />
            {{ item.label }}
          </NuxtLink>
        </li>
      </ul>
    </nav>

    <div class="border-t border-surface-100 p-4 text-xs text-surface-400">
      Veltiq · MVP
    </div>
  </aside>
</template>
