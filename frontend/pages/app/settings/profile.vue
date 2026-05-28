<script setup lang="ts">
definePageMeta({
  layout: 'app',
  middleware: ['require-auth'],
})
useHead({ title: 'Профиль — Veltiq' })

const { user, logout } = useAuth()
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-semibold text-surface-900">Профиль</h1>
      <p class="text-sm text-surface-500">Данные вашего аккаунта.</p>
    </div>

    <VCard title="Учётная запись">
      <dl class="grid gap-4 sm:grid-cols-2">
        <div>
          <dt class="text-xs uppercase tracking-wide text-surface-500">Email</dt>
          <dd class="mt-1 text-sm text-surface-900">{{ user?.email ?? '—' }}</dd>
        </div>
        <div>
          <dt class="text-xs uppercase tracking-wide text-surface-500">ID пользователя</dt>
          <dd class="mt-1 text-sm text-surface-900">{{ user?.id ?? '—' }}</dd>
        </div>
      </dl>
    </VCard>

    <VCard title="Сессия">
      <p class="text-sm text-surface-600">
        Завершить текущую сессию и вернуться на страницу входа.
      </p>
      <div class="mt-4">
        <VButton variant="outline" icon="lucide:log-out" @click="logout">
          Выйти из аккаунта
        </VButton>
      </div>
    </VCard>
  </div>
</template>
