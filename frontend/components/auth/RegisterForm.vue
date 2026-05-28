<script setup lang="ts">
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { registerSchema } from '~/shared/schemas/auth'

const { register, login } = useAuth()
const toast = useToast()

const submitting = ref(false)

const { handleSubmit, defineField, errors } = useForm({
  validationSchema: toTypedSchema(registerSchema),
  initialValues: { email: '', password: '', passwordConfirm: '' },
})

const [email, emailAttrs] = defineField('email')
const [password, passwordAttrs] = defineField('password')
const [passwordConfirm, passwordConfirmAttrs] = defineField('passwordConfirm')

const onSubmit = handleSubmit(async (values) => {
  submitting.value = true
  const regErr = await register(values.email, values.password)
  if (regErr) {
    submitting.value = false
    return
  }
  const loginErr = await login(values.email, values.password)
  submitting.value = false
  if (loginErr) {
    toast.warning('Аккаунт создан', {
      description: 'Войдите вручную для продолжения.',
    })
    await navigateTo('/login')
    return
  }
  toast.success('Аккаунт создан!')
  await navigateTo('/app')
})
</script>

<template>
  <form class="space-y-5" novalidate @submit.prevent="onSubmit">
    <header class="space-y-1">
      <h1 class="text-xl font-semibold text-surface-900">Регистрация</h1>
      <p class="text-sm text-surface-500">
        Создайте аккаунт, чтобы загружать чеки и видеть аналитику.
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
      autocomplete="new-password"
      placeholder="Минимум 8 символов"
      icon-left="lucide:lock"
      :error="errors.password"
      hint="Минимум 8 символов"
      required
    />

    <VInput
      v-model="passwordConfirm"
      v-bind="passwordConfirmAttrs"
      type="password"
      label="Повторите пароль"
      autocomplete="new-password"
      icon-left="lucide:lock"
      :error="errors.passwordConfirm"
      required
    />

    <VButton type="submit" block :loading="submitting">Создать аккаунт</VButton>

    <p class="text-center text-sm text-surface-500">
      Уже есть аккаунт?
      <NuxtLink to="/login" class="font-medium text-brand-600 hover:underline">
        Войти
      </NuxtLink>
    </p>
  </form>
</template>
