#!/bin/bash

# GoBazaar Development Setup Script
# Этот скрипт настраивает среду разработки для проекта

set -e

echo "🚀 Настройка среды разработки GoBazaar..."

# Проверка наличия Go
if ! command -v go &> /dev/null; then
    echo "❌ Go не установлен. Пожалуйста, установите Go 1.23 или выше."
    exit 1
fi

# Проверка версии Go
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.23"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "❌ Требуется Go $REQUIRED_VERSION или выше. Установлена версия: $GO_VERSION"
    exit 1
fi

echo "✅ Go версия $GO_VERSION"

# Проверка наличия Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker не установлен. Пожалуйста, установите Docker."
    exit 1
fi

echo "✅ Docker установлен"

# Проверка наличия Docker Compose
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose не установлен. Пожалуйста, установите Docker Compose."
    exit 1
fi

echo "✅ Docker Compose установлен"

# Установка зависимостей
echo "📦 Установка зависимостей..."
make deps

# Установка инструментов разработки
echo "🔧 Установка инструментов разработки..."
make tools

# Создание pre-commit hooks
echo "🪝 Установка pre-commit hooks..."
make install-hooks

# Создание .env файла, если его нет
if [ ! -f ".env" ]; then
    echo "📝 Создание .env файла..."
    cat > .env << EOF
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gobazaar

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# NATS
NATS_URL=nats://localhost:4222

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# Stripe
STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret

# Service URLs
AUTH_SERVICE_URL=localhost:8080
PRODUCT_SERVICE_URL=localhost:8081
CART_SERVICE_URL=localhost:8082
ORDER_SERVICE_URL=localhost:8083
PAYMENT_SERVICE_URL=localhost:8084
EOF
    echo "✅ .env файл создан"
else
    echo "⚠️  .env файл уже существует"
fi

# Форматирование кода
echo "🎨 Форматирование кода..."
make fmt

# Запуск тестов
echo "🧪 Запуск тестов..."
make test-quick

echo ""
echo "🎉 Настройка завершена!"
echo ""
echo "Полезные команды:"
echo "  make help          - показать все доступные команды"
echo "  make build         - собрать все сервисы"
echo "  make test          - запустить тесты с покрытием"
echo "  make docker-up     - запустить все сервисы в Docker"
echo "  make status        - показать статус проекта"
echo ""
echo "Для начала работы:"
echo "  1. Запустите: make docker-up"
echo "  2. Откройте: http://localhost:8000"
echo "" 