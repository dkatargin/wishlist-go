#!/bin/bash

# Скрипт для сброса и пересоздания debug окружения
# Использование: ./scripts/reset-debug.sh

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${RED}=== ВНИМАНИЕ: Это удалит все данные debug окружения! ===${NC}\n"
read -p "Продолжить? [y/N] " -n 1 -r
echo

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Отменено"
    exit 0
fi

echo -e "\n${YELLOW}Остановка контейнеров...${NC}"
docker compose -f docker/docker-compose.debug.yml down -v

echo -e "${YELLOW}Удаление образов...${NC}"
docker rmi wishlist-backend-debug wishlist-frontend-debug 2>/dev/null || true

echo -e "${YELLOW}Очистка volumes...${NC}"
docker volume rm wishlist_postgres_debug wishlist_go_cache wishlist_frontend_node_modules 2>/dev/null || true

echo -e "${GREEN}Очистка завершена!${NC}"
echo -e "${YELLOW}Для запуска заново используйте: make debug${NC}"

