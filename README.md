# DDD_example
Implementation of domain-driven design logic

---

## 🚀 Быстрый старт

### Требования
- Go 1.24+
- PostgreSQL 15+
- Docker & Docker Compose (опционально)

### Запуск с Docker (рекомендуется)

```bash
# Клонировать репозиторий
git clone https://github.com/akonovalovdev/DDD_example.git
cd DDD_example

# Запустить всё одной командой (сборка + БД + миграции + сервер)
docker-compose up -d --build

# Проверить статус контейнеров
docker-compose ps

# Проверить что сервер работает
curl http://localhost:8080/health
```

> **Примечание:** При запуске автоматически:
> - Поднимается PostgreSQL
> - Применяются миграции (goose)
> - Создаётся пользователь с id=1 и балансом 1000.00
> - Запускается сервер на порту 8080

### Запуск локально

```bash
# 1. Клонировать репозиторий
git clone https://github.com/akonovalovdev/DDD_example.git
cd DDD_example

# 2. Создать .env файл (или экспортировать переменные)
cp .env.example .env

# 3. Запустить PostgreSQL
docker run -d --name postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=skinport \
  -p 5432:5432 \
  postgres:15-alpine

# 4. Установить goose (если не установлен)
go install github.com/pressly/goose/v3/cmd/goose@latest

# 5. Применить миграции
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/skinport?sslmode=disable"
make migrate-up

# 6. Запустить сервер
make run
```

---

## 📋 ENV переменные

| Переменная | Описание | Значение по умолчанию |
|------------|----------|----------------------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:postgres@localhost:5432/skinport?sslmode=disable` |
| `SERVER_PORT` | Порт HTTP сервера | `8080` |
| `SERVER_READ_TIMEOUT` | Таймаут чтения запроса | `10s` |
| `SERVER_WRITE_TIMEOUT` | Таймаут записи ответа | `10s` |
| `SERVER_SHUTDOWN_TIMEOUT` | Таймаут graceful shutdown | `5s` |
| `DB_MAX_OPEN_CONNS` | Макс. открытых соединений к БД | `25` |
| `DB_MAX_IDLE_CONNS` | Макс. idle соединений к БД | `5` |
| `DB_CONN_MAX_LIFETIME` | Время жизни соединения | `5m` |
| `CACHE_TTL` | Время жизни кэша | `5m` |
| `SKINPORT_API_URL` | URL Skinport API | `https://api.skinport.com/v1` |
| `SKINPORT_TIMEOUT` | Таймаут запросов к Skinport | `30s` |
| `LOG_LEVEL` | Уровень логирования | `info` |
| `LOG_FORMAT` | Формат логов | `json` |

---

## 📡 API Endpoints

### Health Check
```bash
curl http://localhost:8080/health
```
**Response:** `{"status":"ok"}`

---

### GET /items
Получение списка предметов Skinport с минимальными ценами (tradable и non-tradable)

```bash
curl -X GET http://localhost:8080/items
```

**Response:**
```json
[
  {
    "market_hash_name": "AK-47 | Redline (Field-Tested)",
    "currency": "USD",
    "suggested_price": "15.23",
    "item_page": "https://skinport.com/item/...",
    "market_page": "https://skinport.com/market/...",
    "tradable_min_price": "12.50",
    "non_tradable_min_price": "10.20",
    "max_price": "25.00",
    "mean_price": "18.75",
    "quantity": 150,
    "created_at": 1609459200,
    "updated_at": 1609459200
  }
]
```

---

### POST /users/{id}/withdraw
Списание баланса пользователя

```bash
curl -X POST http://localhost:8080/users/1/withdraw \
  -H "Content-Type: application/json" \
  -d '{"amount": "100.00"}'
```

**Response (успех):**
```json
{
  "success": true,
  "transaction_id": "550e8400-e29b-41d4-a716-446655440000",
  "balance_before": "1000.00",
  "balance_after": "900.00"
}
```

**Response (недостаточно средств):**
```json
{
  "error": "insufficient balance"
}
```

**Response (пользователь не найден):**
```json
{
  "error": "user not found"
}
```

---

### GET /users/{id}/balance
Получение текущего баланса пользователя

```bash
curl -X GET http://localhost:8080/users/1/balance
```

**Response:**
```json
{
  "user_id": 1,
  "balance": "900.00"
}
```

---

## 🛠 Makefile команды

```bash
make run            # Запуск сервера
make build          # Сборка бинарника
make test           # Запуск тестов
make lint           # Линтинг кода
make migrate-up     # Применить миграции
make migrate-down   # Откатить миграции
make migrate-status # Статус миграций
make docker-build   # Собрать Docker образ
make docker-run     # Запустить через docker-compose
make docker-down    # Остановить docker-compose
make clean          # Очистить артефакты сборки
```

---

## 📁 Структура проекта (Clean Architecture + DDD Lite)

```
DDD_example/
├── cmd/
│   └── server/
│       └── main.go                 # Точка входа приложения
├── internal/
│   ├── domain/                     # СЛОЙ 1: Бизнес-логика (ядро)
│   │   ├── item/
│   │   │   ├── entity.go           # Сущность Item
│   │   │   └── errors.go           # Доменные ошибки
│   │   ├── user/
│   │   │   ├── entity.go           # Сущность User
│   │   │   ├── balance.go          # Методы работы с балансом
│   │   │   └── errors.go           # ErrInsufficientBalance, ErrUserNotFound
│   │   └── transaction/
│   │       └── entity.go           # Сущность Transaction (история)
│   ├── application/                # СЛОЙ 2: Use Cases
│   │   ├── item_service.go         # Логика получения items с кэшем
│   │   └── balance_service.go      # Логика списания баланса
│   ├── ports/                      # СЛОЙ 3: Интерфейсы (порты)
│   │   ├── input/                  # Входящие порты (use cases)
│   │   │   ├── item_service.go
│   │   │   └── balance_service.go
│   │   └── output/                 # Исходящие порты (репозитории)
│   │       ├── item_fetcher.go
│   │       ├── user_repository.go
│   │       ├── transaction_repository.go
│   │       └── cache.go
│   ├── adapters/                   # СЛОЙ 4: Адаптеры (реализации)
│   │   ├── http/
│   │   │   ├── server.go           # HTTP сервер
│   │   │   └── handlers/
│   │   │       ├── item_handler.go
│   │   │       └── balance_handler.go
│   │   ├── repository/
│   │   │   └── postgres/
│   │   │       ├── user_repository.go
│   │   │       └── transaction_repository.go
│   │   └── skinport/
│   │       └── client.go           # Клиент Skinport API
│   ├── config/
│   │   └── config.go               # Парсинг YAML + валидация ENV
│   └── pkg/
│       └── cache/
│           └── inmemory.go         # In-memory кэш
├── config/
│   └── config.yaml                 # Конфигурация приложения
├── migrations/                     # Goose миграции
│   ├── 001_create_users_table.sql
│   ├── 002_create_transactions_table.sql
│   └── 003_seed_user.sql
├── Makefile
├── go.mod
└── README.md
```

---

## 🗄 База данных

### Таблицы

**users**
| Поле | Тип | Описание |
|------|-----|----------|
| id | BIGSERIAL | Primary key |
| balance | DECIMAL(15,2) | Баланс пользователя |
| created_at | TIMESTAMP | Дата создания |
| updated_at | TIMESTAMP | Дата обновления |

**transactions**
| Поле | Тип | Описание |
|------|-----|----------|
| id | UUID | Primary key |
| user_id | BIGINT | FK на users |
| amount | DECIMAL(15,2) | Сумма операции |
| balance_before | DECIMAL(15,2) | Баланс до операции |
| balance_after | DECIMAL(15,2) | Баланс после операции |
| description | VARCHAR(255) | Описание операции |
| created_at | TIMESTAMP | Дата операции |

---

## 🏗 Архитектура

Проект реализован с использованием **Clean Architecture** и **DDD Lite** (Domain First) подхода:

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Handlers                          │
│                   (adapters/http/handlers)                  │
├─────────────────────────────────────────────────────────────┤
│                    Application Services                     │
│                      (application/)                         │
├─────────────────────────────────────────────────────────────┤
│                      Domain Layer                           │
│                       (domain/)                             │
├─────────────────────────────────────────────────────────────┤
│              Ports (Interfaces)    │    Adapters            │
│                 (ports/)           │  (adapters/repository) │
└─────────────────────────────────────────────────────────────┘
```

---

## 📝 Примечания

- **Кэширование**: Items кэшируются в памяти с TTL (по умолчанию 5 минут)
- **Транзакции**: Списание баланса выполняется в PostgreSQL транзакции с `SELECT ... FOR UPDATE`
- **Без ORM**: Используется чистый `database/sql` с raw SQL запросами
- **Decimal**: Для работы с денежными суммами используется `shopspring/decimal`
