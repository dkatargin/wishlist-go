#!/bin/bash

# Скрипт для быстрого запуска debug окружения
# Использование: ./scripts/debug-start.sh

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Wishlist Debug Environment ===${NC}\n"

# Проверка наличия Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Docker не установлен!${NC}"
    exit 1
fi

# Проверка наличия docker compose
if ! docker compose version &> /dev/null; then
    echo -e "${RED}Docker Compose не установлен!${NC}"
    exit 1
fi

# Создание .env.debug если не существует
if [ ! -f .env.debug ]; then
    echo -e "${YELLOW}Создание .env.debug из шаблона...${NC}"
    cp .env.debug .env.debug.local
    echo -e "${GREEN}✓ .env.debug создан${NC}"
fi

# Очистка старых контейнеров
echo -e "\n${YELLOW}Остановка старых контейнеров...${NC}"
docker compose -f docker/docker-compose.debug.yml down 2>/dev/null || true

# Запуск
echo -e "\n${GREEN}Запуск debug окружения...${NC}\n"
docker compose -f docker/docker-compose.debug.yml --env-file .env.debug up --build -d

# Ожидание запуска
echo -e "\n${YELLOW}Ожидание запуска сервисов...${NC}"
sleep 5

# Проверка здоровья
echo -e "\n${GREEN}=== Статус сервисов ===${NC}"
docker compose -f docker/docker-compose.debug.yml ps

echo -e "\n${GREEN}=== Доступные сервисы ===${NC}"
echo -e "${YELLOW}Backend:${NC}        http://localhost:8081"
echo -e "${YELLOW}Frontend:${NC}       http://localhost:3002"
echo -e "${YELLOW}Adminer (DB UI):${NC} http://localhost:8083"
echo -e "${YELLOW}Delve Debugger:${NC}  localhost:2345"
echo -e "${YELLOW}PostgreSQL:${NC}     localhost:5433"

echo -e "\n${GREEN}=== Полезные команды ===${NC}"
echo -e "Логи всех сервисов:    ${YELLOW}make logs ENV=debug${NC}"
echo -e "Логи backend:          ${YELLOW}make logs ENV=debug SERVICE=backend${NC}"
echo -e "Перезапуск backend:    ${YELLOW}make restart ENV=debug SERVICE=backend${NC}"
echo -e "Подключение к БД:      ${YELLOW}make db-shell${NC}"
echo -e "Остановка:             ${YELLOW}make down${NC}"

echo -e "\n${GREEN}=== Debug окружение запущено! ===${NC}\n"

# Опционально показать логи
read -p "Показать логи? [y/N] " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    docker compose -f docker/docker-compose.debug.yml logs -f
fi

