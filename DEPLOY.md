# Veltiq — деплой на VPS

Развёртывание через Docker Compose. Открываем наружу только **80-й порт**, всё остальное живёт во внутренней Docker-сети — это не мешает уже запущенному xray на 443.

## Что получится

```
http://147.45.232.247
        │
        ▼
   nginx :80   ← единственный публичный порт
        │
   ┌────┴─────┐
   │          │
/api/v1     /*
   │          │
backend:8080  frontend:3000 (Nuxt SSR)
   │
postgres:5432 (внутри сети, наружу не светим)
```

## 1. Подготовка VPS

```bash
# зайди на сервер
ssh root@147.45.232.247

# обнови систему
apt update && apt upgrade -y

# установи Docker + Compose plugin
curl -fsSL https://get.docker.com | sh
apt install -y docker-compose-plugin

# проверь
docker --version
docker compose version
```

Открой 80-й порт в файрволе, если включён `ufw`:

```bash
ufw status
# если active — добавь:
ufw allow 80/tcp
```

Если используешь облачный файрвол хостера (Selectel/TimeWeb/etc) — открой 80-й там же.

## 2. Подтянуть проект

```bash
cd /opt
git clone https://github.com/<твой-юзер>/Veltiq.git veltiq
cd veltiq
```

## 3. Создать `.env` для docker-compose

```bash
cp .env.example .env
# открой и поменяй DB_PASSWORD на что-то длинное и случайное
nano .env
```

Сгенерировать пароль одной командой:

```bash
openssl rand -base64 24
```

## 4. Создать прод-конфиг бэка

```bash
cp deploy/backend.config.example.yaml deploy/backend.config.yaml
nano deploy/backend.config.yaml
```

Что обязательно поменять:

| Поле                  | Значение                                                |
|-----------------------|---------------------------------------------------------|
| `db.password`         | то же, что DB_PASSWORD в `.env`                         |
| `jwt.secret_key`      | длинная случайная строка: `openssl rand -hex 48`         |

Остальное можно не трогать — `host: postgres` указывает на контейнер Postgres внутри Docker-сети.

## 5. Собрать и запустить

```bash
docker compose up -d --build
```

Первая сборка минут 5–10 (тянет образы Go/Node, билдит Nuxt). Дальше — секунды.

Проверь, что всё поднялось:

```bash
docker compose ps
```

Должно быть 4 сервиса в статусе `running`/`healthy`:

```
NAME                IMAGE                STATUS
veltiq-backend-1    veltiq-backend       Up
veltiq-frontend-1   veltiq-frontend      Up
veltiq-nginx-1      nginx:1.27-alpine    Up
veltiq-postgres-1   postgres:16-alpine   Up (healthy)
```

## 6. Открыть в браузере

```
http://147.45.232.247
```

Должен открыться лендинг. Зарегистрируйся, залогинься, загрузи `sample-receipts.csv` со страницы импорта.

## Полезные команды

```bash
# логи всех сервисов
docker compose logs -f

# логи только бэка
docker compose logs -f backend

# перезапустить один сервис (например, после правки конфига)
docker compose restart backend

# обновить код и пересобрать
git pull
docker compose up -d --build

# остановить всё
docker compose down

# остановить и удалить тома (СНЕСЁТ БД)
docker compose down -v

# подключиться к Postgres
docker compose exec postgres psql -U veltiq -d veltiq
```

## Если что-то сломалось

**`docker compose ps` показывает один из сервисов как `Exited`:**

```bash
docker compose logs <название_сервиса>
```

**Backend падает с ошибкой подключения к БД:**
- Проверь, что `db.password` в `deploy/backend.config.yaml` совпадает с `DB_PASSWORD` в `.env`.
- Postgres мог не успеть подняться — `docker compose restart backend`.

**Сайт открывается, но регистрация даёт 500:**
- В логе бэка должна быть строка `running database auto migration` — иначе таблиц нет.
- Глянь `docker compose logs backend` сразу после `up`.

**Браузер не открывает сайт:**
- `curl -I http://localhost` прямо на VPS — отвечает nginx?
- Снаружи: проверь, что у хостера в облачном файрволе открыт 80.

## Что НЕ задеплоено

- **HTTPS** — для защиты диплома HTTP по IP сойдёт. Если позже захочешь HTTPS: либо вынеси xray на нестандартный порт и подними certbot, либо поставь сайт за Cloudflare Tunnel (бесплатно, без публичных портов).
- **Бэкап БД** — для учебного проекта не критично, но `docker compose exec postgres pg_dump -U veltiq veltiq > backup.sql` сделает дамп.

## Безопасность для прода (если будешь развивать после защиты)

1. Поставь HTTPS (Cloudflare Tunnel — самый простой путь при занятом 443).
2. После HTTPS поменяй `cookies.secure: true` в `backend.config.yaml`.
3. Не коммить `deploy/backend.config.yaml` и `.env` в git — они уже в `.gitignore`.
4. `ufw` оставь активным, открыты только 22 (SSH) и 80.
