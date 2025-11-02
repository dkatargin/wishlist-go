# Wishlist App

Приложение для создания и управления списками желаний с интеграцией Telegram.

## Технологический стек

- **Backend**: Go 1.25, Gin Framework
- **Frontend**: React 19, Vite, Material-UI
- **База данных**: PostgreSQL 16
- **Контейнеризация**: Docker, Docker Compose

## Структура проекта

```
wishlist-go/
├── backend/          # Go backend сервер
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

## Мониторинг и логи

### Просмотр логов
```bash
# Все сервисы
docker-compose logs -f

# Конкретный сервис
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f postgres
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

[Укажите вашу лицензию]

## Контакты

[Укажите контактную информацию]

