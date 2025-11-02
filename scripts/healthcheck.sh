#!/bin/bash

# Скрипт для мониторинга здоровья сервисов
# Использование: ./scripts/healthcheck.sh [dev|debug|prod]

set -e

ENV=${1:-debug}
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}=== Health Check для $ENV окружения ===${NC}\n"

# Определение портов в зависимости от окружения
case $ENV in
    dev)
        BACKEND_PORT=8080
        FRONTEND_PORT=3001
        POSTGRES_PORT=5432
        ;;
    debug)
        BACKEND_PORT=8081
        FRONTEND_PORT=3002
        POSTGRES_PORT=5433
        ;;
    prod)
        BACKEND_PORT=8080
        FRONTEND_PORT=3000
        POSTGRES_PORT=5432
        ;;
    *)
        echo -e "${RED}Неизвестное окружение: $ENV${NC}"
        echo "Использование: $0 [dev|debug|prod]"
        exit 1
        ;;
esac

check_service() {
    local name=$1
    local url=$2
    local description=$3

    echo -n "Проверка $description... "
    if curl -s -f -o /dev/null --max-time 5 "$url"; then
        echo -e "${GREEN}✓ OK${NC}"
        return 0
    else
        echo -e "${RED}✗ FAILED${NC}"
        return 1
    fi
}

check_tcp_port() {
    local name=$1
    local host=$2
    local port=$3
    local description=$4

    echo -n "Проверка $description (TCP $host:$port)... "
    if nc -z -w5 "$host" "$port" 2>/dev/null; then
        echo -e "${GREEN}✓ OK${NC}"
        return 0
    else
        echo -e "${RED}✗ FAILED${NC}"
        return 1
    fi
}

# Счетчики
total=0
passed=0

# PostgreSQL
((total++))
if check_tcp_port "postgres" "localhost" "$POSTGRES_PORT" "PostgreSQL"; then
    ((passed++))
fi

# Backend API
((total++))
if check_service "backend" "http://localhost:$BACKEND_PORT/api/health" "Backend API"; then
    ((passed++))
fi

# Frontend
((total++))
if check_service "frontend" "http://localhost:$FRONTEND_PORT" "Frontend"; then
    ((passed++))
fi

# Для debug окружения проверяем дополнительные сервисы
if [ "$ENV" = "debug" ]; then
    # Adminer
    ((total++))
    if check_service "adminer" "http://localhost:8083" "Adminer (DB UI)"; then
        ((passed++))
    fi

    # Delve
    ((total++))
    if check_tcp_port "delve" "localhost" "2345" "Delve Debugger"; then
        ((passed++))
    fi
fi

# Итоговый отчет
echo -e "\n${GREEN}=== Результаты ===${NC}"
echo -e "Прошло проверку: ${GREEN}$passed${NC}/$total"

if [ $passed -eq $total ]; then
    echo -e "${GREEN}Все сервисы работают нормально!${NC}"
    exit 0
else
    echo -e "${RED}Некоторые сервисы недоступны!${NC}"
    exit 1
fi

