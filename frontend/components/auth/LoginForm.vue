<script setup lang="ts">
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { loginSchema } from '~/shared/schemas/auth'

const { login } = useAuth()
const toast = useToast()
const route = useRoute()

const submitting = ref(false)

const { handleSubmit, defineField, errors } = useForm({
  validationSchema: toTypedSchema(loginSchema),
  initialValues: { email: '', password: '' },
})

const [email, emailAttrs] = defineField('email')
const [password, passwordAttrs] = defineField('password')

const onSubmit = handleSubmit(async (values) => {
  submitting.value = true
  const err = await login(values.email, values.password)
  submitting.value = false
  if (err) return

  toast.success('Добро пожаловать!')
  const redirect = (route.query.redirect as string) || '/app'
  await navigateTo(redirect)
})
</script>

<template>
  <form class="space-y-5" novalidate @submit.prevent="onSubmit">
    <header class="space-y-1">
      <h1 class="text-xl font-semibold text-surface-900">Вход в Veltiq</h1>
      <p class="text-sm text-surface-500">
        Войдите, чтобы продолжить работу с аналитикой чеков.
      </p>
    </header>

    <VInput
      v-model="email"
      v-bind="emailAttrs"
      type="email"
      label="Email"
      autocomplete="email"
      placeholder="you@example.com"
      icon-left="lucide:mail"
      :error="errors.email"
      required
    />

    <VInput
      v-model="password"
      v-bind="passwordAttrs"
      type="password"
      label="Пароль"
      autocomplete="current-password"
      placeholder="••••••••"
      icon-left="lucide:lock"
      :error="errors.password"
      required
    />

    <VButton type="submit" block :loading="submitting">Войти</VButton>

    <p class="text-center text-sm text-surface-500">
      Нет аккаунта?
      <NuxtLink to="/register" class="font-medium text-brand-600 hover:underline">
        Зарегистрироваться
      </NuxtLink>
    </p>
  </form>
</template>
