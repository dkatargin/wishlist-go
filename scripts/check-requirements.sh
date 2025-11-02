#!/bin/bash

# Скрипт проверки готовности системы для debug стенда
# Использование: ./scripts/check-requirements.sh

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Проверка требований для Debug Stand ===${NC}\n"

# Счетчики
total=0
passed=0
warnings=0

check_command() {
    local cmd=$1
    local name=$2
    local required=$3
    local install_hint=$4

    ((total++))
    echo -n "Проверка $name... "

    if command -v "$cmd" &> /dev/null; then
        version=$($cmd --version 2>&1 | head -n1)
        echo -e "${GREEN}✓ Установлен${NC} ($version)"
        ((passed++))
        return 0
    else
        if [ "$required" = "true" ]; then
            echo -e "${RED}✗ НЕ УСТАНОВЛЕН${NC}"
            echo -e "  ${YELLOW}Установка: $install_hint${NC}"
            return 1
        else
            echo -e "${YELLOW}⚠ Не установлен (опционально)${NC}"
            echo -e "  ${YELLOW}Установка: $install_hint${NC}"
            ((warnings++))
            return 0
        fi
    fi
}

check_port() {
    local port=$1
    local name=$2

    ((total++))
    echo -n "Проверка порта $port ($name)... "

    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1 ; then
        echo -e "${YELLOW}⚠ ЗАНЯТ${NC}"
        process=$(lsof -Pi :$port -sTCP:LISTEN | tail -n1)
        echo -e "  ${YELLOW}Процесс: $process${NC}"
        ((warnings++))
        return 1
    else
        echo -e "${GREEN}✓ Свободен${NC}"
        ((passed++))
        return 0
    fi
}

check_docker_running() {
    ((total++))
    echo -n "Проверка Docker daemon... "

    if docker info >/dev/null 2>&1; then
        echo -e "${GREEN}✓ Запущен${NC}"
        ((passed++))
        return 0
    else
        echo -e "${RED}✗ НЕ ЗАПУЩЕН${NC}"
        echo -e "  ${YELLOW}Запустите Docker Desktop${NC}"
        return 1
    fi
}

check_docker_compose_version() {
    ((total++))
    echo -n "Проверка Docker Compose v2... "

    if docker compose version >/dev/null 2>&1; then
        version=$(docker compose version --short)
        echo -e "${GREEN}✓ $version${NC}"
        ((passed++))
        return 0
    else
        echo -e "${RED}✗ Не найден${NC}"
        echo -e "  ${YELLOW}Обновите Docker Desktop до последней версии${NC}"
        return 1
    fi
}

check_memory() {
    ((total++))
    echo -n "Проверка доступной памяти... "

    # macOS
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # Получаем общую память в GB
        total_mem=$(sysctl -n hw.memsize | awk '{print int($1/1024/1024/1024)}')
        # Проверяем что Docker Desktop запущен и может использовать память
        if [ $total_mem -gt 8 ]; then
            echo -e "${GREEN}✓ ${total_mem}GB системной памяти${NC}"
            ((passed++))
        elif [ $total_mem -gt 4 ]; then
            echo -e "${YELLOW}⚠ ${total_mem}GB системной памяти (рекомендуется 8GB+)${NC}"
            ((warnings++))
        else
            echo -e "${RED}✗ Только ${total_mem}GB системной памяти${NC}"
            echo -e "  ${YELLOW}Рекомендуется минимум 4GB${NC}"
        fi
    # Linux
    else
        free_mem=$(free -m | awk '/Mem:/ {print $7}')
        if [ $free_mem -gt 2048 ]; then
            echo -e "${GREEN}✓ ${free_mem}MB доступно${NC}"
            ((passed++))
        elif [ $free_mem -gt 1024 ]; then
            echo -e "${YELLOW}⚠ ${free_mem}MB доступно (рекомендуется 4GB+)${NC}"
            ((warnings++))
        else
            echo -e "${RED}✗ Только ${free_mem}MB доступно${NC}"
            echo -e "  ${YELLOW}Закройте ненужные приложения${NC}"
        fi
    fi
}

check_disk_space() {
    ((total++))
    echo -n "Проверка свободного места на диске... "

    free_space=$(df -h . | awk 'NR==2 {print $4}' | sed 's/Gi*//')

    if [ "${free_space%.*}" -gt 5 ]; then
        echo -e "${GREEN}✓ ${free_space}GB свободно${NC}"
        ((passed++))
    else
        echo -e "${YELLOW}⚠ Только ${free_space}GB свободно (рекомендуется 10GB+)${NC}"
        ((warnings++))
    fi
}

# Основные проверки
echo -e "${BLUE}--- Обязательные зависимости ---${NC}"
check_command "docker" "Docker" "true" "https://docs.docker.com/get-docker/"
check_docker_running
check_docker_compose_version
check_command "make" "Make" "true" "brew install make (macOS) или apt-get install make (Linux)"

echo -e "\n${BLUE}--- Опциональные зависимости (для локальной разработки) ---${NC}"
check_command "go" "Go" "false" "https://go.dev/dl/"
check_command "bun" "Bun" "false" "curl -fsSL https://bun.sh/install | bash"
check_command "node" "Node.js" "false" "https://nodejs.org/"
check_command "psql" "PostgreSQL Client" "false" "brew install postgresql (macOS)"
check_command "golangci-lint" "golangci-lint" "false" "brew install golangci-lint"

echo -e "\n${BLUE}--- Проверка портов ---${NC}"
check_port 8081 "Backend Debug"
check_port 3002 "Frontend Debug"
check_port 5433 "PostgreSQL Debug"
check_port 8083 "Adminer"
check_port 2345 "Delve Debugger"

echo -e "\n${BLUE}--- Системные ресурсы ---${NC}"
check_memory
check_disk_space

# Итоговый отчет
echo -e "\n${BLUE}=== Итоговый отчет ===${NC}"
echo -e "Проверок пройдено: ${GREEN}$passed${NC}/$total"

if [ $warnings -gt 0 ]; then
    echo -e "Предупреждений: ${YELLOW}$warnings${NC}"
fi

if [ $passed -eq $total ]; then
    echo -e "\n${GREEN}✓ Все проверки пройдены! Система готова к запуску debug стенда.${NC}"
    echo -e "\n${BLUE}Следующие шаги:${NC}"
    echo -e "  1. ${YELLOW}make debug-d${NC}     - Запустить debug окружение"
    echo -e "  2. ${YELLOW}make health${NC}      - Проверить здоровье сервисов"
    echo -e "  3. ${YELLOW}make logs ENV=debug${NC} - Посмотреть логи"
    exit 0
elif [ $((passed + warnings)) -eq $total ]; then
    echo -e "\n${YELLOW}⚠ Есть предупреждения, но можно продолжить.${NC}"
    echo -e "\n${BLUE}Запустить debug окружение:${NC}"
    echo -e "  ${YELLOW}make debug-d${NC}"
    exit 0
else
    echo -e "\n${RED}✗ Требуется установка обязательных зависимостей!${NC}"
    echo -e "\n${BLUE}После установки повторите проверку:${NC}"
    echo -e "  ${YELLOW}./scripts/check-requirements.sh${NC}"
    exit 1
fi

