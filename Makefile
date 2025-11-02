.PHONY: help dev prod logs clean test build install deps health restart
export COMPOSE_PROJECT_NAME=wishlist
# Переменные
DOCKER_COMPOSE := docker compose
DEV_COMPOSE := $(DOCKER_COMPOSE) -f docker/docker-compose.dev.yml
PROD_COMPOSE := $(DOCKER_COMPOSE) -f docker/docker-compose.yml

# Цвета для красивого вывода
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m # No Color

## help: Показать эту справку
help:
	@echo "$(GREEN)Доступные команды:$(NC)"
	@echo ""
	@echo "$(YELLOW)Управление окружениями:$(NC)"
	@echo "  make dev            - Запустить development окружение (hot-reload)"
	@echo "  make prod           - Запустить production окружение"
	@echo ""
	@echo "$(YELLOW)Управление сервисами:$(NC)"
	@echo "  make up ENV=dev     - Запустить окружение (dev/prod)"
	@echo "  make down ENV=dev   - Остановить окружение"
	@echo "  make restart ENV=dev - Перезапустить окружение"
	@echo "  make logs ENV=dev   - Показать логи (опционально SERVICE=backend)"
	@echo ""
	@echo "$(YELLOW)Разработка:$(NC)"
	@echo "  make build          - Собрать все образы"
	@echo "  make rebuild        - Пересобрать образы без кэша"
	@echo "  make shell SVC=backend - Зайти в контейнер"
	@echo "  make exec SVC=backend CMD='go test' - Выполнить команду"
	@echo ""
	@echo "$(YELLOW)База данных:$(NC)"
	@echo "  make db-shell       - Подключиться к PostgreSQL"
	@echo "  make db-migrate     - Запустить миграции"
	@echo "  make db-reset       - Сбросить базу данных"
	@echo ""
	@echo "$(YELLOW)Тестирование:$(NC)"
	@echo "  make test           - Запустить тесты"
	@echo "  make test-backend   - Тесты backend"
	@echo "  make test-frontend  - Тесты frontend"
	@echo ""
	@echo "$(YELLOW)Очистка:$(NC)"
	@echo "  make clean          - Остановить и удалить все контейнеры"
	@echo "  make clean-volumes  - Удалить volumes (БД будет очищена!)"
	@echo "  make clean-all      - Полная очистка (контейнеры + volumes + образы)"
	@echo ""
	@echo "$(YELLOW)Мониторинг:$(NC)"
	@echo "  make health         - Проверить здоровье сервисов"
	@echo "  make ps             - Статус контейнеров"
	@echo "  make stats          - Статистика использования ресурсов"

## dev: Запустить development окружение с hot-reload
dev:
	@echo "$(GREEN)Запуск development окружения...$(NC)"
	@echo "$(YELLOW)Backend Delve debugger будет доступен на порту 2345$(NC)"
	@echo "$(YELLOW)Frontend доступен на http://localhost:3002$(NC)"
	@$(DEV_COMPOSE) up --build

## dev-detached: Запустить dev в фоне
dev-d:
	@echo "$(GREEN)Запуск development окружения в фоновом режиме...$(NC)"
	@$(DEV_COMPOSE) up -d --build
	@make health ENV=dev

## prod: Запустить production окружение
prod:
	@echo "$(GREEN)Запуск production окружения...$(NC)"
	@$(PROD_COMPOSE) up -d --build
	@make health ENV=prod

## down: Остановить все окружения
down:
	@echo "$(YELLOW)Остановка всех окружений...$(NC)"
	@$(DEV_COMPOSE) down 2>/dev/null || true
	@$(PROD_COMPOSE) down 2>/dev/null || true

## logs: Показать логи (использование: make logs ENV=dev SERVICE=backend)
logs:
	@if [ "$(ENV)" = "dev" ]; then \
		if [ -z "$(SERVICE)" ]; then \
			$(DEV_COMPOSE) logs -f; \
		else \
			$(DEV_COMPOSE) logs -f $(SERVICE); \
		fi \
	elif [ "$(ENV)" = "prod" ]; then \
		if [ -z "$(SERVICE)" ]; then \
			$(PROD_COMPOSE) logs -f; \
		else \
			$(PROD_COMPOSE) logs -f $(SERVICE); \
		fi \
	else \
		echo "$(RED)Укажите ENV=dev|prod$(NC)"; \
	fi

## restart: Перезапустить сервис
restart:
	@if [ "$(ENV)" = "dev" ]; then \
		$(DEV_COMPOSE) restart $(SERVICE); \
	elif [ "$(ENV)" = "prod" ]; then \
		$(PROD_COMPOSE) restart $(SERVICE); \
	else \
		echo "$(RED)Укажите ENV=dev|prod$(NC)"; \
	fi

## build: Собрать образы
build:
	@echo "$(GREEN)Сборка образов...$(NC)"
	@$(PROD_COMPOSE) build

## rebuild: Пересобрать образы без кэша
rebuild:
	@echo "$(GREEN)Пересборка образов без кэша...$(NC)"
	@$(PROD_COMPOSE) build --no-cache

## shell: Зайти в контейнер (использование: make shell SVC=backend ENV=dev)
shell:
	@if [ -z "$(SVC)" ]; then \
		echo "$(RED)Укажите SVC=backend|frontend|postgres$(NC)"; \
		exit 1; \
	fi
	@if [ "$(ENV)" = "dev" ]; then \
		$(DEV_COMPOSE) exec $(SVC) /bin/sh; \
	elif [ "$(ENV)" = "prod" ]; then \
		$(PROD_COMPOSE) exec $(SVC) /bin/sh; \
	else \
		echo "$(RED)Укажите ENV=dev|prod$(NC)"; \
	fi

## exec: Выполнить команду в контейнере
exec:
	@if [ -z "$(SVC)" ] || [ -z "$(CMD)" ]; then \
		echo "$(RED)Использование: make exec SVC=backend CMD='go test' ENV=dev$(NC)"; \
		exit 1; \
	fi
	@if [ "$(ENV)" = "dev" ]; then \
		$(DEV_COMPOSE) exec $(SVC) $(CMD); \
	elif [ "$(ENV)" = "prod" ]; then \
		$(PROD_COMPOSE) exec $(SVC) $(CMD); \
	else \
		echo "$(RED)Укажите ENV=dev|prod$(NC)"; \
	fi

## db-shell: Подключиться к PostgreSQL
db-shell:
	@echo "$(GREEN)Подключение к PostgreSQL...$(NC)"
	@$(DEV_COMPOSE) exec postgres psql -U wishlist_user -d wishlist

## db-reset: Сбросить базу данных
db-reset:
	@echo "$(YELLOW)Сброс базы данных...$(NC)"
	@$(DEV_COMPOSE) exec postgres psql -U wishlist_user -d wishlist -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	@echo "$(GREEN)База данных сброшена$(NC)"

## test: Запустить все тесты
test: test-backend test-frontend

## test-backend: Тесты backend
test-backend:
	@echo "$(GREEN)Запуск backend тестов...$(NC)"
	@cd backend && go test -v ./...

## test-frontend: Тесты frontend
test-frontend:
	@echo "$(GREEN)Запуск frontend тестов...$(NC)"
	@cd frontend && bun test

## clean: Остановить и удалить контейнеры
clean:
	@echo "$(YELLOW)Очистка контейнеров...$(NC)"
	@$(DEV_COMPOSE) down --remove-orphans
	@$(PROD_COMPOSE) down --remove-orphans

## clean-volumes: Удалить volumes
clean-volumes:
	@echo "$(RED)ВНИМАНИЕ: Все данные БД будут удалены!$(NC)"
	@read -p "Продолжить? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		$(DEV_COMPOSE) down -v; \
		$(PROD_COMPOSE) down -v; \
		echo "$(GREEN)Volumes удалены$(NC)"; \
	fi

## clean-all: Полная очистка
clean-all: clean clean-volumes
	@echo "$(YELLOW)Удаление образов...$(NC)"
	@docker rmi wishlist-backend wishlist-frontend 2>/dev/null || true
	@echo "$(GREEN)Полная очистка завершена$(NC)"

## health: Проверить здоровье сервисов
health:
	@echo "$(GREEN)Проверка здоровья сервисов...$(NC)"
	@echo ""
	@echo "$(YELLOW)PostgreSQL:$(NC)"
	@curl -f http://localhost:5432 2>/dev/null && echo "✓ Доступен" || echo "✗ Недоступен"
	@echo ""
	@echo "$(YELLOW)Backend API:$(NC)"
	@curl -f http://localhost:8080/api/health 2>/dev/null && echo "✓ Доступен" || echo "✗ Недоступен"
	@echo ""
	@echo "$(YELLOW)Frontend:$(NC)"
	@curl -f http://localhost:3000 2>/dev/null && echo "✓ Доступен" || echo "✗ Недоступен" || \
	 curl -f http://localhost:3001 2>/dev/null && echo "✓ Доступен (dev)" || echo "✗ Недоступен"

## ps: Статус контейнеров
ps:
	@echo "$(GREEN)Статус контейнеров:$(NC)"
	@docker ps -a --filter "name=wishlist-" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

## stats: Статистика использования ресурсов
stats:
	@docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}" \
		$$(docker ps --filter "name=wishlist-" -q)

## install: Установить зависимости локально
install:
	@echo "$(GREEN)Установка зависимостей...$(NC)"
	@cd backend && go mod download
	@cd frontend && bun install

## fmt: Форматирование кода
fmt:
	@echo "$(GREEN)Форматирование кода...$(NC)"
	@cd backend && go fmt ./...
	@cd frontend && bun run format 2>/dev/null || echo "Formatter не настроен"

## lint: Линтинг кода
lint:
	@echo "$(GREEN)Линтинг кода...$(NC)"
	@cd backend && golangci-lint run ./... 2>/dev/null || echo "$(YELLOW)golangci-lint не установлен$(NC)"
	@cd frontend && bun run lint 2>/dev/null || echo "$(YELLOW)ESLint не настроен$(NC)"

