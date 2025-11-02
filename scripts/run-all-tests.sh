#!/bin/bash

# Comprehensive test script для всех окружений
# Использование: ./scripts/run-all-tests.sh [dev|debug|prod]

set -e

ENV=${1:-debug}
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== Запуск тестов для $ENV окружения ===${NC}\n"

# Определение compose файла
case $ENV in
    dev)
        COMPOSE_FILE="docker/docker-compose.dev.yml"
        ;;
    debug)
        COMPOSE_FILE="docker/docker-compose.debug.yml"
        ;;
    prod)
        COMPOSE_FILE="docker/docker-compose.yml"
        ;;
    *)
        echo -e "${RED}Неизвестное окружение: $ENV${NC}"
        echo "Использование: $0 [dev|debug|prod]"
        exit 1
        ;;
esac

# Функция для запуска теста
run_test() {
    local name=$1
    local command=$2

    echo -e "\n${YELLOW}▶ $name${NC}"
    if eval "$command"; then
        echo -e "${GREEN}✓ $name - PASSED${NC}"
        return 0
    else
        echo -e "${RED}✗ $name - FAILED${NC}"
        return 1
    fi
}

# Счетчики
total_tests=0
passed_tests=0

# Проверка что окружение запущено
echo -e "${BLUE}Проверка окружения...${NC}"
if ! docker compose -f "$COMPOSE_FILE" ps | grep -q "Up"; then
    echo -e "${YELLOW}Окружение не запущено. Запускаем...${NC}"
    docker compose -f "$COMPOSE_FILE" up -d
    sleep 10
fi

# Backend тесты
echo -e "\n${BLUE}=== Backend Tests ===${NC}"

((total_tests++))
if run_test "Backend: Health Check" \
    "curl -sf http://localhost:8081/api/health || curl -sf http://localhost:8080/api/health"; then
    ((passed_tests++))
fi

((total_tests++))
if run_test "Backend: Go Tests" \
    "cd backend && go test -v ./... -short"; then
    ((passed_tests++))
fi

((total_tests++))
if run_test "Backend: Go Vet" \
    "cd backend && go vet ./..."; then
    ((passed_tests++))
fi

# Frontend тесты
echo -e "\n${BLUE}=== Frontend Tests ===${NC}"

((total_tests++))
if run_test "Frontend: Access Test" \
    "curl -sf http://localhost:3002 || curl -sf http://localhost:3001 || curl -sf http://localhost:3000"; then
    ((passed_tests++))
fi

((total_tests++))
if run_test "Frontend: Build Check" \
    "cd frontend && bun install && bun run build"; then
    ((passed_tests++))
fi

# Database тесты
echo -e "\n${BLUE}=== Database Tests ===${NC}"

((total_tests++))
if run_test "Database: Connection Test" \
    "docker compose -f $COMPOSE_FILE exec -T postgres pg_isready -U wishlist_user"; then
    ((passed_tests++))
fi

((total_tests++))
if run_test "Database: Schema Check" \
    "docker compose -f $COMPOSE_FILE exec -T postgres psql -U wishlist_user -d wishlist -c '\dt'"; then
    ((passed_tests++))
fi

# Integration тесты
echo -e "\n${BLUE}=== Integration Tests ===${NC}"

((total_tests++))
if run_test "Integration: Backend -> Database" \
    "curl -sf http://localhost:8081/api/accounts || curl -sf http://localhost:8080/api/accounts"; then
    ((passed_tests++))
fi

# Performance тесты (опционально)
if [ "$ENV" = "debug" ]; then
    echo -e "\n${BLUE}=== Performance Tests ===${NC}"

    ((total_tests++))
    if run_test "Performance: Response Time" \
        "time curl -sf http://localhost:8081/api/health -w '\nTime: %{time_total}s\n'"; then
        ((passed_tests++))
    fi
fi

# Итоговый отчет
echo -e "\n${BLUE}=== Итоговый отчет ===${NC}"
echo -e "Тестов пройдено: ${passed_tests}/${total_tests}"

if [ $passed_tests -eq $total_tests ]; then
    echo -e "\n${GREEN}✓ Все тесты пройдены успешно!${NC}"
    exit 0
else
    failed=$((total_tests - passed_tests))
    echo -e "\n${RED}✗ Провалено тестов: ${failed}${NC}"
    exit 1
fi

