# Veltiq Backend

Go 1.26 · Gin · GORM · PostgreSQL · JWT в httpOnly cookies.

REST API для SaaS-продукта по сокращению неликвидного товара: принимает чековые книги розничных магазинов в CSV, разбирает их, считает аналитику и выдаёт рекомендации по товарам-кандидатам на распродажу/списание.

## Стек

| Слой | Технология |
|---|---|
| Язык | Go 1.26 |
| HTTP-роутер | [gin-gonic/gin](https://github.com/gin-gonic/gin) |
| ORM / БД | [GORM](https://gorm.io/) + Postgres 16 |
| Конфиг | [spf13/viper](https://github.com/spf13/viper) (YAML) |
| Auth | [golang-jwt/jwt v5](https://github.com/golang-jwt/jwt) · access+refresh JWT в httpOnly cookies |
| Пароли | bcrypt |
| ID | UUID v4 |
| Логгер | `log/slog` (stdlib) — JSON в stdout |

## Архитектура

Гексагональная / Ports & Adapters:

```
cmd/
└─ main.go                # точка входа

internal/
├─ app/
│  └─ bootstrap.go        # сборка зависимостей и DI
│
├─ core/                  # доменный слой, не зависит от HTTP/БД
│  ├─ domain/             # User, Import, Receipt, Report, ImportStatus, ошибки
│  ├─ ports/              # интерфейсы (что core ожидает от инфраструктуры)
│  ├─ service/            # бизнес-логика (AuthService, ImportService, ReportService)
│  └─ orchestrator/       # Runner — фасад поверх сервисов для HTTP-слоя
│
├─ infrastructure/        # реализация портов
│  ├─ config/             # загрузка YAML через Viper
│  ├─ cookies/            # установка/чтение access/refresh cookies
│  ├─ http/
│  │  ├─ handlers/api/v1/ # Gin-хэндлеры (auth, profile, import)
│  │  ├─ middleware/      # CORS, RequestID, Recovery, Logger,
│  │  │                   # MaxBodyBytes, AuthRateLimit, AuthMiddleware, TenantMiddleware
│  │  ├─ httperrors/      # маппинг доменных ошибок → HTTP-коды
│  │  ├─ router/          # сборка маршрутов
│  │  └─ server/          # gracefull HTTP-сервер
│  ├─ jwt/                # access + refresh, разделение типов токенов
│  ├─ logging/            # slog-обёртка
│  ├─ repository/postgres/# GORM-репозитории (User, Import, Receipt, Workflow)
│  └─ security/           # bcrypt
│
└─ modules/               # «функциональные» модули
   ├─ parser/             # парсинг чековых книг (CSV)
   └─ analytics/          # подсчёт сводки и поиск неликвида
```

**Правила зависимостей:**
- `core/*` не импортирует ничего из `infrastructure/*` и `modules/*` напрямую — общается через `ports`.
- `infrastructure/*` и `modules/*` реализуют `ports`.
- `app/bootstrap.go` склеивает всё в граф зависимостей.

Такая структура позволяет менять БД/HTTP-фреймворк/парсер без правки бизнес-логики.

## Авторизация

JWT-пара лежит в **httpOnly cookies** — фронт их не видит в JS, что даёт устойчивость к XSS:

| Cookie | Назначение | TTL по умолчанию |
|---|---|---|
| `access_token` | Авторизация всех приватных запросов | 15 мин |
| `refresh_token` | Получение новой пары без перелогина | 7 дней |

**Параметры cookies:** `HttpOnly=true`, `Path=/`, `SameSite=Lax`, `Secure` управляется `cookies.secure` в конфиге.

**Поток:**

```
POST /register      → bcrypt-хэш + UUID в users
POST /login         → проверка пароля → выписка пары токенов → Set-Cookie
GET  /auth/session  → проверка access; если просрочен — попытка через refresh
POST /auth/refresh  → обмен refresh → новая пара токенов
POST /auth/logout   → MaxAge=-1 для обеих cookies
```

`AuthMiddleware` дёргает access из cookies, верифицирует через JWT-менеджер, кладёт `userID` в контекст. `TenantMiddleware` вычисляет `tenantID` детерминированно из `userID` (см. `domain.TenantIDFromUserID`) — каждый пользователь == свой тенант. Все запросы к импортам/отчётам фильтруются по `tenantID`.

## API

База: `/api/v1`. Все ответы — JSON. Все ошибки маппятся через `httperrors.Write` в формат `{ "error": "..." }` с HTTP-кодом.

### Публичные эндпоинты

| Метод | Путь | Тело | Ответ | Что делает |
|---|---|---|---|---|
| `POST` | `/register` | `{ email, password }` | `201 { "status": "created" }` | Регистрация. `password >= 8` |
| `POST` | `/login` | `{ email, password }` | `200 { "status": "authenticated" }` + Set-Cookie | Логин |
| `GET` | `/auth/session` | — | `200 { valid, user_id, tenant_id, access_expired? }` или `401` | Проверка сессии |
| `POST` | `/auth/refresh` | — | `200 { "status": "refreshed" }` + Set-Cookie | Обновить access по refresh |
| `POST` | `/auth/logout` | — | `200 { "status": "logged_out" }` + Set-Cookie clear | Выход |

`/register` и `/login` дополнительно проходят через **rate-limit** (`auth.rate_limit_per_minute` в конфиге, дефолт 30/мин по IP).

### Приватные эндпоинты (требуют access cookie)

| Метод | Путь | Ответ | Что делает |
|---|---|---|---|
| `GET` | `/profile` | `{ id, email?, tenantId? }` | Профиль текущего юзера |
| `GET` | `/imports?limit=N` | `{ items: ImportRecord[], total, cursor }` | Список импортов тенанта (свежие сверху) |
| `POST` | `/imports` | `201 { id, status }` | Загрузка чековой книги. multipart `file` или сырой body |
| `GET` | `/imports/:id/status` | `{ id, status, errorCode?, createdAt, updatedAt }` | Статус (для polling) |
| `GET` | `/imports/:id/report` | `{ importId, status, data: ReportSummary, createdAt, updatedAt }` | Сводный отчёт |
| `GET` | `/imports/:id/insights` | `{ importId, generatedAt, items: IlliquidItem[] }` | Рекомендации по неликвиду |

### Доменные типы

```go
type ImportStatus = "pending" | "processing" | "partial_failed" | "done" | "failed"

type ReportSummary struct {
    receiptsCount int
    totalSum      int  // в основных денежных единицах (₽)
    cashSum       int
    cardSum       int
    isStub        bool
    note          string // опционально
}

type IlliquidItem struct {
    sku                string
    name               string
    category           string  // опционально
    salesQuantity      int     // сколько штук продано за период
    daysWithoutSale    int
    lastSaleAt         string  // RFC3339
    recommendation     "discount"|"bundle"|"writeoff"|"monitor"
    recommendationNote string
}
```

### Формат CSV для загрузки

Парсер ждёт UTF-8 CSV с заголовком:

```
receipt_id,date,store_id,payment,sku,product_name,category,quantity,unit_price
1001,2026-04-01 12:34,1,card,SKU-001,Хлеб белый,Хлебобулочные,1,45
1001,2026-04-01 12:34,1,card,SKU-005,Молоко 2.5%,Молочные,2,80
1002,2026-04-01 13:10,1,cash,SKU-001,Хлеб белый,Хлебобулочные,1,45
```

| Колонка | Тип | Описание |
|---|---|---|
| `receipt_id` | string | Внешний ID чека (для группировки позиций одного чека) |
| `date` | `YYYY-MM-DD HH:MM` | Дата и время продажи |
| `store_id` | int | Номер магазина |
| `payment` | `cash`/`card` | Способ оплаты |
| `sku` | string | Артикул товара |
| `product_name` | string | Наименование |
| `category` | string | Категория (опционально) |
| `quantity` | int | Кол-во штук |
| `unit_price` | int | Цена за штуку в основных денежных единицах |

Одна строка = одна позиция в чеке. Парсер группирует строки по `receipt_id`+`date`, формируя чек с агрегатом и позициями.

Пример: [frontend/public/sample-receipts.csv](../frontend/public/sample-receipts.csv).

## База данных

Postgres 16. Используется GORM с `AutoMigrate` — миграции выполняются при старте сервера, если `db.auto_migrate: true` в конфиге.

### Таблицы

| Таблица | Зачем | Ключевые поля |
|---|---|---|
| `users` | Аккаунты | `id (uuid)`, `email (unique)`, `password_hash`, `created_at` |
| `imports` | Загруженные книги | `id (uuid)`, `tenant_id`, `document_id (sha256)`, `status`, `error_code`, `created_at`, `updated_at` |
| `receipts` | Агрегаты чеков | `id`, `import_id`, `tenant_id`, `external_id`, `store_id`, `payment_type`, `sum`, `cashier`, `created_at` |
| `receipt_items` | Позиции чеков | `id`, `import_id`, `tenant_id`, `receipt_external_id`, `sku`, `name`, `category`, `quantity`, `unit_price`, `total_price`, `payment_type`, `sold_at` |

Уникальный индекс `(tenant_id, document_id)` в `imports` гарантирует **идемпотентность загрузки** — повторная загрузка того же файла возвращает уже существующий импорт.

## Конфигурация

`CONFIG_PATH` указывает на YAML. Пример: [config.yaml](config.yaml).

```yaml
env: local                    # local | prod — влияет на дефолт cookies.secure
http_server:
  address: "0.0.0.0:8080"
  timeout: 30s
  idle_timeout: 60s
  max_body_bytes: 26214400    # 25 МБ — лимит на upload
db:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: veltiq
  sslmode: disable
  auto_migrate: true          # AutoMigrate при старте
jwt:
  secret_key: "длинная-случайная-строка"  # openssl rand -hex 48
  access_ttl: 15m
  refresh_ttl: 168h           # 7 дней
cookies:
  secure: false               # true для HTTPS-прода
cors:
  allowed_origins:
    - "http://localhost:3000"
    - "http://127.0.0.1:3000"
auth:
  rate_limit_per_minute: 30   # rate-limit для /register и /login по IP
```

**Важно:** Viper Unmarshal по умолчанию использует `mapstructure`-теги, но в `loader.go` мы переопределяем `c.TagName = "yaml"`, поэтому работают существующие `yaml:"..."`-теги.

## Локальный запуск

### Требования

- Go 1.26+
- Postgres 16 (локально или в Docker)

### Подготовка БД

```powershell
# через Docker
docker run --name veltiq-pg -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=veltiq -p 5432:5432 -d postgres:16-alpine
```

### Запуск

```powershell
cd backend
$env:CONFIG_PATH = "$PWD\config.yaml"
go run ./cmd
```

В логе должно появиться:

```
running database auto migration
[GIN-debug] POST   /api/v1/register ...
server started on 0.0.0.0:8080
```

### Сборка бинарника

```powershell
go build -o veltiq.exe ./cmd
```

## Тесты

```powershell
go test ./...
```

## Деплой в Docker

См. корневой [DEPLOY.md](../DEPLOY.md) и [docker-compose.yml](../docker-compose.yml). Bёрстка: multistage `Dockerfile` → distroless `nonroot` образ ~15 МБ.

## Известные ограничения / TODO

- **CSRF**: без явного токена. Защищаемся `SameSite=Lax` + httpOnly cookies. Для бóльших гарантий стоит добавить double-submit token при переходе на cross-site сценарии.
- **Refresh-rotation**: текущий refresh не одноразовый. Для прода стоит выдавать новый refresh при каждом `/auth/refresh` и хранить чёрный список отозванных.
- **Multi-tenancy**: tenant сейчас = user. Для команд нужно завести отдельную таблицу `tenants` и `tenant_users`.
- **Тесты**: integration-тесты на репозитории через `testcontainers-go` пока не написаны.
