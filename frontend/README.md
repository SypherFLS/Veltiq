# Veltiq Frontend

Nuxt 3 + Tailwind CSS + Headless UI + Pinia + TanStack Query.

## Требования

- Node.js 20+
- pnpm / npm / yarn (примеры — для npm)
- Запущенный бэкенд Veltiq (по умолчанию `http://localhost:8080`)

## Быстрый старт

```powershell
cp .env.example .env
npm install
npm run dev
```

Открой `http://localhost:3000`.

## Скрипты

| Команда             | Что делает                                |
|---------------------|-------------------------------------------|
| `npm run dev`       | Запуск dev-сервера                        |
| `npm run build`     | Production-сборка                         |
| `npm run preview`   | Просмотр production-сборки                |
| `npm run typecheck` | Проверка типов через `vue-tsc`            |
| `npm run lint`      | Линтинг                                   |
| `npm run lint:fix`  | Линтинг + автофикс                        |
| `npm test`          | Юнит-тесты (Vitest)                       |
| `npm run test:e2e`  | E2E-тесты (Playwright)                    |

## Структура

```
shared/        # типы, API-клиент, схемы валидации, утилиты
stores/        # Pinia (auth, ui, toast)
composables/   # useAuth, useToast, useImportPolling, useReport
middleware/    # auth.global, guest, require-auth
plugins/       # 01.api.ts (ofetch + refresh-флоу), 02.vue-query.ts
components/
  ui/          # дизайн-система: V*
  layout/      # AppSidebar, AppTopbar, UserMenu
  auth/        # LoginForm, RegisterForm
  imports/     # ImportStatusBadge, ...
  reports/     # ReportSummaryCards, ...
layouts/       # default (лендинг), auth (формы), app (приватная зона)
pages/         # /, /login, /register, /app/**
```

## Поток авторизации

1. На клиенте `middleware/auth.global.ts` вызывает `GET /api/v1/auth/session`.
2. Если `access_expired`, дёргается `POST /api/v1/auth/refresh` и сессия перепроверяется.
3. Все `$fetch`-вызовы идут с `credentials: 'include'`; на 401 `plugins/01.api.ts` пробует refresh один раз и повторяет запрос.
4. Если refresh упал — `useAuthStore.clear()` + редирект на `/login`.

## Конфигурация

`NUXT_PUBLIC_API_BASE` в `.env` — URL бэкенда. Должен быть тем же origin, что указан в `CORS.AllowedOrigins` у бэка (с `AllowCredentials: true`).
