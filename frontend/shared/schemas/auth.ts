import { z } from 'zod'

export const emailSchema = z
  .string({ required_error: 'Введите email' })
  .min(1, 'Введите email')
  .email('Некорректный email')

export const passwordSchema = z
  .string({ required_error: 'Введите пароль' })
  .min(8, 'Минимум 8 символов')

export const loginSchema = z.object({
  email: emailSchema,
  password: z
    .string({ required_error: 'Введите пароль' })
    .min(1, 'Введите пароль'),
})

export const registerSchema = z
  .object({
    email: emailSchema,
    password: passwordSchema,
    passwordConfirm: z
      .string({ required_error: 'Повторите пароль' })
      .min(1, 'Повторите пароль'),
  })
  .refine((data) => data.password === data.passwordConfirm, {
    message: 'Пароли не совпадают',
    path: ['passwordConfirm'],
  })

export type LoginFormValues = z.infer<typeof loginSchema>
export type RegisterFormValues = z.infer<typeof registerSchema>
