# Veltiq — Frontend Plan

## 0. О проекте

**Veltiq** — SaaS-продукт для малых розничных сетей. Помогает сокращать неликвидный товар на основе аналитики чековых книг за выбранные периоды: пользователь загружает выгрузку чеков, бэкенд парсит и считает метрики, фронт показывает сводку, графики и рекомендации по товарам-кандидатам на списание/распродажу.

Проект учебный (университет, диплом / курсовой предпоследнего курса).

**Стек:**
- Backend: Go + Gin, хексагональная архитектура.
- Frontend (этот документ): Nuxt 3 + Tailwind CSS + Headless UI.

---

## 1. Состояние бэкенда на момент написания плана

### 1.1. Архитектура

```
backend/
├─ cmd/main.go
└─ internal/
   ├─ app/bootstrap.go
   ├─ core/
   │  ├─ domain/        # errors, import, metrics, receipt, report, tenant, user
   │  ├─ ports/         # интерфейсы (analytics, stores, parser, token_manager, ...)
   │  ├─ service/       # auth_service, import_service, report_service
   │  └─ orchestrator/  # Runner для импортов
   ├─ infrastructure/
   │  ├─ config/
   │  ├─ cookies/       # Manager — кладёт/читает access/refresh из httpOnly cookies
   │  ├─ http/
   │  │  ├─ handlers/api/v1/   # auth.go, import.go, profile.go
   │  │  ├─ middleware/        # CORS, RequestID, Recovery, RequestLogger,
   │  │  │                     # MaxBodyBytes, AuthRateLimit, AuthMiddleware, TenantMiddleware
   │  │  ├─ router/router.go
   │  │  ├─ httperrors/
   │  │  └─ server/
   │  ├─ jwt/           # access + refresh, разделение типов токенов
   │  ├─ logging/
   │  ├─ repository/
   │  └─ security/
   └─ modules/
      ├─ analytics/
      └─ parser/
```

### 1.2. Что уже реализовано

Полностью работает:

| Метод   | Путь                  | Что делает                                            |
|---------|-----------------------|-------------------------------------------------------|
| POST    | `/api/v1/register`    | Регистрация (`email`, `password >= 8`)                |
| POST    | `/api/v1/login`       | Логин → ставит httpOnly cookies (access + refresh)    |
| GET     | `/api/v1/auth/session`| Проверка сессии. Возвращает `valid`, `user_id`, `tenant_id`, опц. `access_expired` |
| POST    | `/api/v1/auth/refresh`| Обмен refresh → новая пара токенов                    |
| POST    | `/api/v1/auth/logout` | Чистит cookies                                        |

Скаффолд (handlers есть, бизнес-логика может быть частичной):

| Метод | Путь                              | Что делает                       |
|-------|-----------------------------------|----------------------------------|
| GET   | `/api/v1/profile`                 | Профиль текущего пользователя    |
| POST  | `/api/v1/imports`                 | Загрузка файла чековой книги     |
| GET   | `/api/v1/imports/:id/status`      | Статус импорта (для polling)     |
| GET   | `/api/v1/imports/:id/report`      | Получение отчёта по импорту      |

Все `protected` маршруты идут через `AuthMiddleware` + `TenantMiddleware` (tenant выводится из userID, мультитенантность встроена в домен).

### 1.3. Ключевые архитектурные факты для фронта

- **JWT лежит в httpOnly cookies** — фронт **не видит** токен в JS. Это XSS-стойко, но требует:
  - все запросы с `credentials: 'include'`;
  - CORS на бэке с `AllowCredentials: true` и явным origin (не `*`);
  - проверка сессии — только через `/auth/session`, нельзя декодировать JWT на клиенте.
- **Refresh-флоу**: на 401 фронт должен один раз дёрнуть `/auth/refresh` и повторить запрос. `/auth/session` уже умеет fallback на refresh-токен, возвращая `access_expired: true`.
- **Tenant** проставляется бэком автоматически из userID — фронт явно его никуда не передаёт.
- **Доменный тип отчёта** (`backend/internal/core/domain/report.go`):

  ```go
  type ReportSummary struct {
      ReceiptsCount int    `json:"receiptsCount"`
      TotalSum      int    `json:"totalSum"`
      CashSum       int    `json:"cashSum"`
      CardSum       int    `json:"cardSum"`
      IsStub        bool   `json:"isStub"`
      Note          string `json:"note,omitempty"`
  }
  ```

  Зеркалить во фронт-типы в [shared/types/report.ts](shared/types/report.ts).

---

## 2. Фронтенд: технологический стек

| Категория        | Выбор                                          | Зачем                                                |
|------------------|------------------------------------------------|------------------------------------------------------|
| Фреймворк        | **Nuxt 3** (Vue 3 + Vite)                      | Composition API, auto-imports, нормальный TS         |
| Стили            | **Tailwind CSS** (`@nuxtjs/tailwindcss`)       | Скорость, утилитарность                              |
| UI-примитивы     | **Headless UI** (`@headlessui/vue`)            | Доступные модалки, dropdown, comboboxes              |
| UI fallback      | `radix-vue` / `reka-ui`                        | Сложные компоненты (datepicker, tabs с URL-стейтом)  |
| State (UI/auth)  | **Pinia**                                      | Глобальный стор для auth и UI                        |
| State (server)   | **TanStack Query** (`@tanstack/vue-query`)     | Кэш, refetch, polling статусов импорта               |
| Формы            | **VeeValidate** + **Zod** + `@vee-validate/zod`| Валидация, типобезопасные схемы                      |
| Утилиты          | **VueUse**                                     | Debounce, intersection observer, storage             |
| Графики          | **Chart.js** + **vue-chartjs**                 | Просто и достаточно (ECharts — если нужно мощнее)    |
| Даты             | **dayjs**                                      | Лёгкий, удобный                                      |
| Иконки           | **Nuxt Icon** (`@nuxt/icon` + iconify)         | Не тянуть несколько библиотек                        |
| Тесты            | **Vitest** (unit) + **Playwright** (e2e)       | Базовый auth-флоу для защиты                         |
| Линтер           | ESLint + Prettier (через `@nuxt/eslint`)       |                                                      |

**Чего НЕ берём:**
- Никаких `nuxt-auth` / `@sidebase/nuxt-auth` — кастомный cookie-флоу пишется руками.
- Без SSR-fetch для приватных страниц — `/app/**` идёт в SPA-режиме.
- Без openapi-генератора типов — для учебного проекта overkill, типы пишем руками.

---

## 3. Архитектурные решения

### 3.1. Режим рендеринга — гибрид

```ts
// nuxt.config.ts
routeRules: {
  '/':         { prerender: true },  // лендинг
  '/login':    { ssr: true },
  '/register': { ssr: true },
  '/app/**':   { ssr: false },        // приватная зона — чистый SPA
}
```

Лендинг и auth-страницы — SSR/prerender (SEO, скорость). Приватная зона — SPA, чтобы не возиться с форвардом cookies на сервере.

### 3.2. Авторизация

- Источник правды о юзере = ответ `GET /auth/session`.
- На старте SPA глобальный middleware `auth.global.ts` зовёт `/auth/session`:
  - `valid:true` → пишем юзера в `useAuthStore`;
  - `access_expired:true` → дёргаем `/auth/refresh`, повторяем `/auth/session`;
  - `valid:false` и роут приватный → `navigateTo('/login')`.
- В `$fetch` интерсептор:
  - все запросы — `credentials: 'include'`;
  - на 401 (кроме `/auth/refresh`) — пробуем refresh один раз, повторяем запрос;
  - если refresh упал → `authStore.clear()` + `navigateTo('/login')`.

### 3.3. API-слой

Всё общение с бэком — через тонкие модули в [shared/api/](shared/api/). В компонентах **не** используем `$fetch` напрямую и **не** хардкодим URL.

### 3.4. Состояние

Жёсткое разделение:

| Тип данных                       | Где                          |
|----------------------------------|------------------------------|
| Серверные (импорты, отчёты, …)   | TanStack Query (`useQuery`)  |
| Auth / current user / tenant     | Pinia (`useAuthStore`)       |
| UI-стейт (сайдбар, модалки)      | local `ref` или `useUIStore` |
| Формы                            | VeeValidate field state      |

---

## 4. Файловая структура

```
frontend/
├─ app.vue
├─ error.vue
├─ nuxt.config.ts
├─ tailwind.config.ts
├─ tsconfig.json
├─ package.json
├─ .env.example                      # NUXT_PUBLIC_API_BASE=http://localhost:8080
│
├─ assets/
│  └─ css/
│     └─ main.css                    # @tailwind base/components/utilities + css-vars темы
│
├─ public/
│  ├─ favicon.svg
│  └─ og.png
│
├─ layouts/
│  ├─ default.vue                    # публичный (лендинг)
│  ├─ auth.vue                       # логин/регистрация (центрированная карточка)
│  └─ app.vue                        # приватная зона: sidebar + topbar
│
├─ middleware/
│  ├─ auth.global.ts                 # проверка сессии при первом заходе
│  ├─ guest.ts                       # для /login,/register — редирект если уже залогинен
│  └─ require-auth.ts                # для /app/** — редирект на /login если нет сессии
│
├─ pages/
│  ├─ index.vue                      # лендинг
│  ├─ login.vue
│  ├─ register.vue
│  ├─ forgot-password.vue            # заглушка под будущее
│  └─ app/
│     ├─ index.vue                   # дашборд (сводка по тенанту)
│     ├─ imports/
│     │  ├─ index.vue                # список импортов
│     │  ├─ new.vue                  # загрузка файла чековой книги
│     │  └─ [id].vue                 # статус + переход к отчёту
│     ├─ reports/
│     │  ├─ index.vue                # список отчётов
│     │  └─ [id].vue                 # детальный отчёт + графики + советы
│     ├─ inventory/                  # неликвид — главная фича продукта
│     │  └─ index.vue
│     ├─ settings/
│     │  ├─ profile.vue
│     │  ├─ organization.vue         # настройки тенанта
│     │  └─ billing.vue              # заглушка
│     └─ help.vue
│
├─ components/
│  ├─ ui/                            # дизайн-система: обёртки над HeadlessUI + Tailwind
│  │  ├─ VButton.vue
│  │  ├─ VInput.vue
│  │  ├─ VSelect.vue
│  │  ├─ VModal.vue
│  │  ├─ VDropdown.vue
│  │  ├─ VTable.vue                  # generic table со слотами
│  │  ├─ VBadge.vue
│  │  ├─ VCard.vue
│  │  ├─ VEmptyState.vue
│  │  ├─ VSkeleton.vue
│  │  └─ VToast.vue
│  ├─ layout/
│  │  ├─ AppSidebar.vue
│  │  ├─ AppTopbar.vue
│  │  └─ UserMenu.vue
│  ├─ auth/
│  │  ├─ LoginForm.vue
│  │  └─ RegisterForm.vue
│  ├─ imports/
│  │  ├─ ImportUploader.vue          # drag&drop + прогресс
│  │  ├─ ImportStatusBadge.vue
│  │  └─ ImportList.vue
│  ├─ reports/
│  │  ├─ ReportSummaryCards.vue      # 4 карточки: totalSum/cashSum/cardSum/receiptsCount
│  │  ├─ ReportSalesChart.vue
│  │  ├─ ReportIlliquidTable.vue     # главная таблица — товары-кандидаты
│  │  └─ ReportInsights.vue          # список советов от аналитики
│  └─ dashboard/
│     ├─ KpiCard.vue
│     └─ ActivityFeed.vue
│
├─ composables/
│  ├─ useAuth.ts                     # login/logout/session — фасад над api/auth + Pinia
│  ├─ useImportPolling.ts            # useQuery с refetchInterval, пока status != done
│  ├─ useReport.ts
│  ├─ useToast.ts
│  └─ useBreakpoints.ts              # из VueUse
│
├─ stores/
│  ├─ auth.ts                        # user, tenantId, isAuthenticated, hydrate()
│  └─ ui.ts                          # sidebar open/closed, theme
│
├─ shared/
│  ├─ api/
│  │  ├─ client.ts                   # $fetch.create({ baseURL, credentials:'include' })
│  │  ├─ auth.ts                     # login, register, logout, session, refresh
│  │  ├─ imports.ts                  # upload, status, report
│  │  ├─ profile.ts
│  │  └─ errors.ts                   # нормализация ошибок бэка
│  ├─ types/
│  │  ├─ api.ts                      # ApiError, Paginated<T>
│  │  ├─ auth.ts                     # SessionResponse, User
│  │  ├─ import.ts                   # ImportStatus enum, Import
│  │  └─ report.ts                   # ReportSummary (зеркало backend domain)
│  ├─ schemas/                       # zod-схемы форм
│  │  ├─ auth.ts
│  │  └─ import.ts
│  └─ utils/
│     ├─ format.ts                   # formatMoney, formatDate
│     └─ result.ts                   # тип-обёртка для api-ответов
│
├─ plugins/
│  ├─ api.ts                         # подкидывает $api в nuxtApp
│  ├─ vue-query.ts                   # QueryClient + persister (опц.)
│  └─ toast.ts
│
└─ tests/
   ├─ unit/
   │  └─ format.spec.ts
   └─ e2e/
      └─ auth.spec.ts                # Playwright: register → login → logout
```

Префикс `V*` для UI-компонентов — чтобы auto-imports Nuxt не конфликтовали и доменные компоненты сразу отличались от примитивов.

---

## 5. Карта страниц (MVP)

| Роут                      | Layout | Доступ      | Что показывает                                                |
|---------------------------|--------|-------------|---------------------------------------------------------------|
| `/`                       | default| public      | Hero, 3 фичи, CTA                                             |
| `/login`                  | auth   | guest       | Форма логина                                                  |
| `/register`               | auth   | guest       | Форма регистрации                                             |
| `/app`                    | app    | auth        | 4 KPI-карточки, последние импорты, быстрый upload             |
| `/app/imports`            | app    | auth        | Список импортов со статусами                                  |
| `/app/imports/new`        | app    | auth        | Drag&drop загрузки файла чековой книги                        |
| `/app/imports/:id`        | app    | auth        | Статус + прогресс, кнопка «Открыть отчёт» при done            |
| `/app/reports`            | app    | auth        | Список отчётов                                                |
| `/app/reports/:id`        | app    | auth        | `ReportSummary` + график + таблица неликвида + советы         |
| `/app/inventory`          | app    | auth        | Общий вид по неликвиду (по всем отчётам тенанта)              |
| `/app/settings/profile`   | app    | auth        | Email пользователя, смена пароля (когда появится на бэке)     |
| `/app/settings/organization` | app | auth        | Настройки тенанта (заглушка)                                  |
| `/app/settings/billing`   | app    | auth        | Заглушка                                                      |

**Фокус на защите:** страница отчёта (`/app/reports/:id`) — комиссию интересует именно домен. Туда вкладывать больше всего UI-усилий.

---

## 6. Критичные потоки данных

### 6.1. Polling статуса импорта

```
Upload файла → бэк вернул { importId }
            → router.push(`/app/imports/${id}`)
            → useImportPolling(id):
                useQuery({
                  queryKey: ['import', id],
                  queryFn: () => api.imports.status(id),
                  refetchInterval: data => data?.status === 'done' || data?.status === 'failed' ? false : 3000,
                })
            → когда status === 'done' — кнопка «Открыть отчёт» → /app/reports/:id
```

### 6.2. Refresh-токен (в plugins/api.ts)

```
onResponseError:
  if (response.status === 401 && !request.url.endsWith('/auth/refresh') && !retried):
    await api.auth.refresh()      // обновит cookies на бэке
    return $fetch(originalRequest) // повтор один раз
  else:
    authStore.clear()
    navigateTo('/login')
```

### 6.3. Загрузка файла

- `FormData` через `$fetch` с `credentials: 'include'`.
- **Не выставлять `Content-Type` вручную** — браузер сам поставит multipart boundary.
- Прогресс — через XHR-обёртку (нативный fetch прогресса не даёт удобно).

---

## 7. Roadmap по неделям (~10–15 ч/нед)

### Неделя 1 — фундамент
- `npx nuxi init`, установить Tailwind, Pinia, VueQuery, Headless UI, ESLint + Prettier.
- Каркас layouts (`default`, `auth`, `app`), пустые pages.
- `shared/api/client.ts` с интерсепторами (CORS, credentials, базовая обработка ошибок).
- Pinia store `auth` + middleware `auth.global.ts`.

### Неделя 2 — auth end-to-end
- Формы login/register: VeeValidate + Zod.
- Обработка ошибок бэка (показ через тост-систему).
- Refresh-флоу в интерсепторе.
- Logout с очисткой Query cache (`queryClient.clear()`).

### Неделя 3 — UI-кит и app shell
- `VButton`, `VInput`, `VModal`, `VTable`, `VBadge`, `VEmptyState`, `VSkeleton`.
- `AppSidebar` + `AppTopbar` + `UserMenu`.
- Скелетон дашборда с моками.

### Неделя 4 — импорты
- Drag&drop uploader.
- Список импортов, статусы (`ImportStatusBadge`), страница импорта с polling.

### Неделя 5 — отчёты
- `ReportSummaryCards` от `ReportSummary` (зеркало бэка).
- График через `vue-chartjs`.
- Таблица неликвида (структура готова, данные могут быть мокаными).

### Неделя 6 — полировка для защиты
- Лендинг.
- Пустые состояния, скелетоны загрузки, тосты.
- Responsive (минимум tablet+).
- README с скринами, базовый Playwright-тест auth-флоу, демо-видео 1–2 мин.

---

## 8. Чего НЕ делать (анти-паттерны)

- Не писать свои tooltip/dropdown/dialog — Headless UI закрывает всё это.
- Не делать i18n с первого дня — сначала русский хардкодом.
- Не пытаться идеально типизировать ответы бэка через openapi-генератор.
- Не делать тёмную тему до защиты (Tailwind `dark:` внедряется за вечер позже).
- Не вкладываться в анимации/parallax на лендинге — комиссии важен домен.
- Не пихать всё в Pinia — серверные данные принадлежат TanStack Query.
- Не использовать `useFetch` в шаблонах для бизнес-эндпоинтов — только через `shared/api/*`.

---

## 9. Что согласовать с бэкендом

- **Единый формат ошибок**: `{ code, message, details? }`. Сейчас разные хэндлеры отвечают по-разному (где `error`, где `status`).
- **Пагинация** списков импортов/отчётов: cursor (`?limit&cursor`) или page (`?page&size`) — выбрать одно.
- **Эндпоинт списка импортов** (`GET /imports`) — отсутствует в `router.go`, фронту нужен.
- **Эндпоинт инсайтов по неликвиду** — отдельный (`GET /imports/:id/insights`) или часть report? Предложение: отдельный, чтобы фронт мог рендерить таблицу независимо от summary и со своим состоянием загрузки.
- **CORS**: убедиться, что `AllowedOrigins` включает `http://localhost:3000` и `AllowCredentials: true`.
- **Эндпоинт смены пароля / выхода со всех устройств** — в roadmap.

---

## 10. Полезные ссылки на бэкенд (для сверки контрактов)

- Роутер: [backend/internal/infrastructure/http/router/router.go](../backend/internal/infrastructure/http/router/router.go)
- Auth-хэндлер: [backend/internal/infrastructure/http/handlers/api/v1/auth.go](../backend/internal/infrastructure/http/handlers/api/v1/auth.go)
- Import-хэндлер: [backend/internal/infrastructure/http/handlers/api/v1/import.go](../backend/internal/infrastructure/http/handlers/api/v1/import.go)
- Profile-хэндлер: [backend/internal/infrastructure/http/handlers/api/v1/profile.go](../backend/internal/infrastructure/http/handlers/api/v1/profile.go)
- Доменные модели: [backend/internal/core/domain/](../backend/internal/core/domain/)
- Cookie-менеджер: [backend/internal/infrastructure/cookies/](../backend/internal/infrastructure/cookies/)
- JWT-логика: [backend/internal/infrastructure/jwt/](../backend/internal/infrastructure/jwt/)
