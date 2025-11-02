# Wishlist App

Приложение для создания и управления списками желаний с интеграцией Telegram.

## Технологический стек

- **Backend**: Go 1.25, Gin Framework
- **Frontend**: React 19, Vite, Material-UI
- **База данных**: PostgreSQL 16
- **Message Queue**: RabbitMQ 3.13
- **Worker**: Go 1.25 (для обработки фоновых задач)
- **Контейнеризация**: Docker, Docker Compose

## Структура проекта

```
wishlist-go/
├── backend/          # Go backend сервер
├── worker/           # Go worker для обработки фоновых задач
├── frontend/         # React frontend приложение
├── docker/           # Docker конфигурации
├── config.yaml       # Конфигурация для локальной разработки
└── config.docker.yaml # Конфигурация для Docker
```

## Быстрый старт с Docker

### Предварительные требования

- Docker >= 20.10
- Docker Compose >= 2.0

### Запуск приложения

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd wishlist-go
```

2. Настройте переменные окружения (опционально):
Отредактируйте `config.docker.yaml` и укажите ваши значения для:
- `telegram.bot_token` - токен Telegram бота
- `sentry.dsn` - DSN для Sentry (если используете)

3. Запустите все сервисы:
```bash
docker-compose up -d
```

4. Проверьте статус:
```bash
docker-compose ps
```

### Доступ к приложению

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **PostgreSQL**: localhost:5432
- **RabbitMQ Management UI**: http://localhost:15672 (по умолчанию: wishlist_user / rabbitpassword)
- **RabbitMQ**: localhost:5672

### Остановка приложения

```bash
docker-compose down
```

Для удаления данных БД:
```bash
docker-compose down -v
```

## Разработка без Docker

### Backend

1. Установите Go 1.25+
2. Установите PostgreSQL
3. Настройте `config.yaml`
4. Запустите сервер:
```bash
cd backend
go mod download
go run server.go
```

### Frontend

1. Установите Bun или Node.js
2. Установите зависимости:
```bash
cd frontend
bun install  # или npm install
```
3. Запустите dev-сервер:
```bash
bun run dev  # или npm run dev
```

## Сборка для продакшена

### Backend отдельно
```bash
docker build -f Dockerfile.backend -t wishlist-backend .
```

### Frontend отдельно
```bash
docker build -f Dockerfile.frontend -t wishlist-frontend \
  --build-arg VITE_BACKEND_HOST=your-api-host \
  --build-arg VITE_BACKEND_PORT=443 \
  --build-arg VITE_BACKEND_SCHEME=https \
  .
```

## Переменные окружения

### Backend (config.yaml)
- `server.host` - хост сервера (по умолчанию: 0.0.0.0 в Docker)
- `server.port` - порт сервера (по умолчанию: 8080)
- `database.*` - параметры подключения к PostgreSQL
- `rabbitmq.host` - хост RabbitMQ (по умолчанию: wishlist-rabbitmq-develop в Docker)
- `rabbitmq.port` - порт RabbitMQ (по умолчанию: 5672)
- `rabbitmq.user` - пользователь RabbitMQ
- `rabbitmq.password` - пароль RabbitMQ
- `rabbitmq.vhost` - virtual host RabbitMQ
- `telegram.bot_token` - токен Telegram бота
- `sentry.dsn` - DSN для мониторинга ошибок

### Frontend (build args)
- `VITE_BACKEND_HOST` - хост backend API
- `VITE_BACKEND_PORT` - порт backend API
- `VITE_BACKEND_SCHEME` - протокол (http/https)
- `APP_VERSION` - версия приложения

## API Endpoints

- `GET /api/health` - проверка здоровья сервиса
- `POST /api/account/login` - авторизация
- `GET /api/wishlists` - получить списки желаний
- `POST /api/wishlists` - создать новый список
- `GET /api/wishlists/:id` - получить конкретный список
- `PUT /api/wishlists/:id` - обновить список
- `DELETE /api/wishlists/:id` - удалить список

## RabbitMQ и Worker

Приложение использует RabbitMQ для асинхронной обработки событий. Backend выступает в роли producer (отправитель сообщений), а Worker - consumer (получатель и обработчик сообщений).

### Архитектура

1. **Backend** - отправляет сообщения в очередь при определенных событиях:
   - Создание wishlist (`wishlist_created`)
   - Создание wishitem (`wishitem_created`)
   - Обновление account (`account_updated`)

2. **Worker** - получает сообщения из очереди и обрабатывает их:
   - Отправка уведомлений
   - Обновление кэша
   - Синхронизация данных
   - Другая фоновая обработка

### Формат сообщений

```json
{
  "type": "wishlist_created",
  "payload": {
    "wishlist_id": "uuid",
    "owner_id": 123,
    "name": "My Wishlist"
  },
  "timestamp": "2025-11-02T12:00:00Z"
}
```

### Использование в коде

Пример отправки сообщения из backend:

```go
import "wishlist-go/internal/queue"

// Отправка сообщения в очередь
if queue.Client != nil {
    err := queue.Client.PublishMessage("wishlist_created", map[string]interface{}{
        "wishlist_id": wishlist.ID,
        "owner_id":    userID,
        "name":        wishlist.Name,
    })
    if err != nil {
        log.Printf("Failed to publish message: %v", err)
    }
}
```

Worker автоматически получает и обрабатывает эти сообщения.

### Мониторинг очереди

RabbitMQ Management UI доступен по адресу http://localhost:15672
- Логин: wishlist_user
- Пароль: rabbitpassword

В интерфейсе можно:
- Просматривать очереди и их состояние
- Мониторить количество сообщений
- Отслеживать производительность
- Управлять соединениями


## Мониторинг и логи

### Просмотр логов
```bash
# Все сервисы
docker-compose logs -f

# Конкретный сервис
docker-compose logs -f backend
docker-compose logs -f worker
docker-compose logs -f frontend
docker-compose logs -f postgres
docker-compose logs -f rabbitmq
```

### Проверка здоровья
```bash
# Backend
curl http://localhost:8080/api/health

# Frontend
curl http://localhost:3000/health
```

## Troubleshooting

### Backend не подключается к БД
- Убедитесь, что PostgreSQL запущен и готов принимать соединения
- Проверьте параметры подключения в `config.docker.yaml`
- Проверьте логи: `docker-compose logs postgres`

### Frontend не может подключиться к Backend
- Проверьте, что backend запущен: `curl http://localhost:8080/api/health`
- Убедитесь, что переменные окружения правильно настроены
- Проверьте сетевые настройки Docker

### Ошибки при сборке
- Очистите Docker кэш: `docker-compose build --no-cache`
- Удалите старые образы: `docker system prune -a`

## Лицензия

MIT

