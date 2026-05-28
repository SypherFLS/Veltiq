<script setup lang="ts">
import {
  Menu,
  MenuButton,
  MenuItem,
  MenuItems,
} from '@headlessui/vue'

const { user, logout } = useAuth()

const initials = computed(() => {
  const email = user.value?.email ?? user.value?.id ?? '?'
  return email.slice(0, 2).toUpperCase()
})
</script>

<template>
  <Menu as="div" class="relative">
    <MenuButton
      class="flex items-center gap-2 rounded-xl border border-surface-200 bg-white px-2 py-1.5 text-sm font-medium text-surface-700 hover:bg-surface-50 focus-ring"
    >
      <span
        class="flex size-7 items-center justify-center rounded-lg bg-brand-100 text-xs font-semibold text-brand-700"
      >
        {{ initials }}
      </span>
      <span class="hidden max-w-[140px] truncate sm:inline">
        {{ user?.email ?? 'Пользователь' }}
      </span>
      <Icon name="lucide:chevron-down" class="size-4 text-surface-400" />
    </MenuButton>

    <transition
      enter-active-class="transition duration-100 ease-out"
      enter-from-class="transform opacity-0 scale-95"
      enter-to-class="transform opacity-100 scale-100"
      leave-active-class="transition duration-75 ease-in"
      leave-from-class="transform opacity-100 scale-100"
      leave-to-class="transform opacity-0 scale-95"
    >
      <MenuItems
        class="absolute right-0 mt-2 w-56 origin-top-right overflow-hidden rounded-xl border border-surface-100 bg-white shadow-elevated focus:outline-none"
      >
        <div class="border-b border-surface-100 px-3 py-2 text-xs text-surface-500">
          {{ user?.email ?? user?.id }}
        </div>
        <MenuItem v-slot="{ active }">
          <NuxtLink
            to="/app/settings/profile"
            :class="[
              'flex items-center gap-2 px-3 py-2 text-sm',
              active ? 'bg-surface-50 text-surface-900' : 'text-surface-700',
            ]"
          >
            <Icon name="lucide:user" class="size-4" />
            Профиль
          </NuxtLink>
        </MenuItem>
        <MenuItem v-slot="{ active }">
          <button
            type="button"
            :class="[
              'flex w-full items-center gap-2 px-3 py-2 text-sm',
              active ? 'bg-red-50 text-red-600' : 'text-red-600',
            ]"
            @click="logout"
          >
            <Icon name="lucide:log-out" class="size-4" />
            Выйти
          </button>
        </MenuItem>
      </MenuItems>
    </transition>
  </Menu>
</template>
