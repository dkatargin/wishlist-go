# Wishlist App

Приложение для создания и управления списками желаний с интеграцией Telegram
(авторизация + Mini App). Позиции можно добавлять вручную или **ссылкой на товар** —
фоновый воркер сам распарсит название, картинку и цену с Яндекс.Маркета, Ozon и
Wildberries.

## Возможности

- **Списки желаний** и позиции в них.
- **Парсинг товаров по ссылке** — Яндекс.Маркет, Ozon, Wildberries (асинхронно, через воркер).
- **Резервирование подарков** — гость может «забронировать» позицию; владелец списка брони не видит.
- **Избранное** — сохранённые чужие списки.
- **Telegram-авторизация** (Mini App `initData`) + публичный гостевой просмотр шаренного списка.

## Технологический стек

- **Backend**: Go 1.25, Gin, GORM
- **Frontend**: React 19, Vite, Material-UI (пакетный менеджер — Bun)
- **База данных**: PostgreSQL 16
- **Очередь**: RabbitMQ 3.13
- **Воркер**: Go (фоновый парсинг товаров)
- **Контейнеризация**: Docker, Docker Compose

## Архитектура

Один Go-модуль `wishlist-go` в каталоге `backend/`, организованный по Clean / Hexagonal
архитектуре: зависимости направлены внутрь, к `domain`; ничто в `domain` не импортирует фреймворк.

```
wishlist-go/
├── backend/                      # Go, модуль `wishlist-go`
│   ├── cmd/api/                  # точка входа API     → app.NewAPIApp
│   ├── cmd/worker/               # точка входа воркера → app.NewWorkerApp
│   └── internal/
│       ├── app/                  # DI-сборка (infra → repo → usecase → handlers → routes)
│       ├── domain/               # сущности + интерфейсы репозиториев (порты)
│       ├── usecase/              # бизнес-логика (по агрегату на пакет)
│       ├── repository/postgres/  # GORM-реализации интерфейсов domain
│       ├── delivery/http/        # Gin: handlers, middleware, dto
│       ├── delivery/worker/      # RabbitMQ consumer
│       └── infrastructure/       # config, database, queue, crawler
├── frontend/                     # React 19 + Vite (Bun)
├── docker/                       # docker-compose.{dev,}.yml + .env.{dev,prod}
└── docs/                         # архитектура и планы
```

Подробнее об архитектуре — в `docs/superpowers/ARCHITECTURE.md`, гайд для разработки — в `CLAUDE.md`.

## Быстрый старт (Docker, режим разработки)

### Требования

- Docker >= 20.10, Docker Compose >= 2.0
- GNU Make

### Конфигурация

Конфиг **только через переменные окружения** (12-factor) — никаких YAML-файлов.
Docker Compose читает значения из `docker/.env.dev` (для dev) и `docker/.env.prod` (для prod).
Эти файлы в `.gitignore` — создайте их и задайте как минимум обязательные переменные
(`POSTGRES_PASSWORD`, `RABBITMQ_PASSWORD`, `TELEGRAM_BOT_TOKEN`). Полный список — в разделе
[Переменные окружения](#переменные-окружения).

### Запуск

```bash
make dev      # сборка + запуск dev-стека с hot-reload (Air) и Delve-дебаггером
make dev-d    # то же, но в фоне (detached)
```

### Доступ к сервисам (dev)

- **Frontend** (Vite): http://localhost:3002
- **Backend API**: http://localhost:8081 (внутри контейнера — `:8080`)
- **Delve** (отладчик Go): `localhost:2345`
- **PostgreSQL**: `localhost:5433`
- **RabbitMQ**: `localhost:5672`
- **RabbitMQ Management UI**: http://localhost:15672 (логин/пароль — из `docker/.env.dev`)

### Остановка и очистка

```bash
make down            # остановить dev- и prod-стеки
make clean           # + удалить контейнеры/сети (--remove-orphans)
make clean-volumes   # + удалить тома (СОТРЁТ данные БД; спросит подтверждение)
make logs SERVICE=backend   # логи сервиса (по умолчанию — всех)
make db-shell        # psql в dev-базу
```

## Локальная разработка (без Docker)

### Backend

```bash
cd backend
go mod download

# конфиг через env: создайте gitignored .env.local с POSTGRES_HOST=localhost
# и опубликованными dev-портами (POSTGRES_PORT=5433, RABBITMQ_PORT=5672)
set -a && . ../.env.local && set +a

go run ./cmd/api       # API-сервер
go run ./cmd/worker    # воркер (в отдельном терминале)
```

### Frontend

```bash
cd frontend
bun install
bun run dev            # Vite dev-сервер
```

## Переменные окружения

`config.Load()` (`backend/internal/infrastructure/config`) заполняет конфиг из env по
тегам `env:`. Обязательные переменные при отсутствии валят старт (fail-fast).

| Переменная | Секция | По умолчанию | Обязательна |
|---|---|---|:---:|
| `SERVER_HOST` | server | `0.0.0.0` | |
| `SERVER_PORT` | server | `8080` | |
| `WORKER_HOST` | worker | `0.0.0.0` | |
| `WORKER_PORT` | worker | `8090` | |
| `POSTGRES_HOST` | database | `localhost` | |
| `POSTGRES_PORT` | database | `5432` | |
| `POSTGRES_USER` | database | `wishlist` | |
| `POSTGRES_PASSWORD` | database | — | ✅ |
| `POSTGRES_DB` | database | `wishlist` | |
| `TELEGRAM_BOT_TOKEN` | telegram | — | ✅ |
| `RABBITMQ_HOST` | rabbitmq | `localhost` | |
| `RABBITMQ_PORT` | rabbitmq | `5672` | |
| `RABBITMQ_USER` | rabbitmq | `guest` | |
| `RABBITMQ_PASSWORD` | rabbitmq | — | ✅ |
| `RABBITMQ_VHOST` | rabbitmq | `/` | |
| `SENTRY_DSN` | sentry | — | |
| `SENTRY_ENVIRONMENT` | sentry | `production` | |
| `SENTRY_RELEASE` | sentry | — | |
| `LOG_LEVEL` | logging | `info` | |
| `LOG_FILE` | logging | — | |

### Frontend (build-args, Vite)

| Переменная | По умолчанию |
|---|---|
| `VITE_BACKEND_HOST` | `wish.dimhost.ru` |
| `VITE_BACKEND_PORT` | `443` |
| `VITE_BACKEND_SCHEME` | `https` |
| `VITE_DEPLOYMENT_TYPE` | `production` |
| `VITE_AUTH_MOCKUP` | `` (пусто) |

## Тесты, форматирование, линт

```bash
make test            # backend + frontend
make test-backend    # cd backend && go test -v ./...
make test-frontend   # cd frontend && bun test
make fmt             # go fmt + bun run format
make lint            # golangci-lint + bun run lint
```

## API

Базовый префикс — `/api/v1`.

**Авторизация**: Telegram Mini App. Клиент шлёт заголовок `Authorization: tma <initData>`;
middleware валидирует HMAC `initData` (ключ выводится из токена бота), парсит пользователя
и **лениво создаёт аккаунт**. Отдельного эндпоинта логина нет.

### Публичные

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/api/v1/health` | проверка здоровья |
| `GET` | `/api/v1/share/:shareCode` | гостевой просмотр шаренного списка (auth опциональна; владельцу/гостю показываются разные данные о бронях) |

### Авторизованные

**Списки**

| Метод | Путь |
|---|---|
| `GET` | `/list` |
| `POST` | `/list` |
| `GET` | `/list/:listId` |
| `PATCH` | `/list/:listId` |
| `DELETE` | `/list/:listId` |

**Позиции**

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/list/:listId/wishes` | |
| `POST` | `/list/:listId/wishes` | создать вручную |
| `POST` | `/list/:listId/wishes/crawl` | создать по ссылке на товар (парсинг) |
| `GET` | `/list/:listId/wishes/:wishId` | |
| `PATCH` | `/list/:listId/wishes/:wishId` | |
| `DELETE` | `/list/:listId/wishes/:wishId` | |

**Резервирование подарков**

| Метод | Путь |
|---|---|
| `POST` | `/share/:shareCode/wishes/:wishId/reserve` |
| `DELETE` | `/share/:shareCode/wishes/:wishId/reserve` |
| `POST` | `/reservations/:reservationId/purchased` |
| `GET` | `/reservations` |

**Избранное** (`:id` = ShareCode списка)

| Метод | Путь |
|---|---|
| `GET` | `/favorites` |
| `POST` | `/wishlist/:id/favorite` |
| `DELETE` | `/wishlist/:id/favorite` |

**Аккаунт**

| Метод | Путь |
|---|---|
| `DELETE` | `/account` |

## Очередь и воркер

Для асинхронной работы используется RabbitMQ. Backend — producer, воркер — consumer.

Backend публикует JSON-сообщения в durable-очередь **`wishlist_tasks`**. Сейчас в системе
один тип задачи — **`crawl_product`**: он отправляется при создании позиции по ссылке
(`POST /list/:listId/wishes/crawl`), воркер забирает сообщение, дёргает краулер, парсит
товар и сохраняет `WishItem`.

Формат сообщения:

```json
{
  "type": "crawl_product",
  "payload": {
    "product_url": "https://market.yandex.ru/card/...",
    "wish_list_code": "<uuid>",
    "owner_id": 123
  },
  "timestamp": "2026-06-29T12:00:00Z"
}
```

## Краулер (парсинг товаров)

Серверные адаптеры под каждый источник; диспетчер выбирает адаптер по хосту ссылки.

| Источник | Хосты | Извлекает |
|---|---|---|
| Яндекс.Маркет | `market.yandex.ru`, `yandex.ru` | название, цена, картинка, описание |
| Ozon | `ozon.ru` | название, цена, картинка (через OpenGraph) |
| Wildberries | `wildberries.ru` | название, картинка (**без цены** — анти-бот) |

## Сборка для продакшена

### Frontend

```bash
docker build -f frontend/Dockerfile -t wishlist-frontend \
  --build-arg VITE_BACKEND_HOST=your-api-host \
  --build-arg VITE_BACKEND_PORT=443 \
  --build-arg VITE_BACKEND_SCHEME=https \
  frontend/
```

### Backend

> **Внимание:** production-стек (`make prod`, `docker/docker-compose.yml`) пока неполный —
> в нём нет сервиса RabbitMQ, и отсутствует production-Dockerfile для backend (есть только
> `backend/Dockerfile.develop` для разработки). Основной поддерживаемый сценарий —
> dev-стек (`make dev`).

## Troubleshooting

- **Backend не стартует** — проверьте, что заданы обязательные env (`POSTGRES_PASSWORD`,
  `RABBITMQ_PASSWORD`, `TELEGRAM_BOT_TOKEN`); при их отсутствии конфиг падает на старте.
- **Backend не подключается к БД** — убедитесь, что Postgres поднят (`make logs SERVICE=postgres`),
  и что хост/порт совпадают (в Docker — `postgres:5432`, локально — `localhost:5433`).
- **Frontend не видит backend** — проверьте `VITE_BACKEND_*` и доступность API
  (`curl http://localhost:8081/api/v1/health` в dev).
- **Пересборка с нуля** — `make clean && make dev` (или `docker system prune -a` для очистки кэша образов).

## Лицензия

[MIT](LICENSE)
