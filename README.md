# Smart Warehouse - Система умного складского учёта

Система автоматизированной инвентаризации складов с использованием роботов и AI для прогнозирования остатков товаров.

## Технологический стек

### Frontend
- React 18 + TypeScript
- Material-UI (MUI)
- Redux Toolkit
- Vite
- WebSocket для real-time обновлений

### Backend
- Go (Golang)
- Gin Web Framework
- GORM (PostgreSQL)
- WebSocket
- GigaChat API (Sber AI)

### Инфраструктура
- Docker & Docker Compose
- PostgreSQL 15
- Redis 7
- Nginx

## Требования

Перед запуском убедитесь, что у вас установлены:

- **Docker Desktop** (версия 20.10+)
- **Docker Compose** (версия 2.0+)
- **Git**

### Проверка установки

```bash
docker --version
docker-compose --version
```

## Быстрый старт

### 1. Клонирование репозитория

```bash
git clone <repository-url>
cd Monkey_Team
```

### 2. Настройка переменных окружения

Создайте файл `.env` в корне проекта:

```bash
cp .env.example .env
```

**ВАЖНО**: Отредактируйте `.env` и укажите ваши данные для GigaChat API (если планируете использовать AI-прогнозы).

### 3. Запуск проекта

Все компоненты запускаются одной командой:

```bash
docker-compose up -d
```

Эта команда:
- Скачает необходимые образы Docker
- Соберёт frontend и backend
- Запустит PostgreSQL и Redis
- Применит миграции базы данных
- Запустит все сервисы

### 4. Проверка статуса

Убедитесь, что все контейнеры запущены:

```bash
docker ps
```

Должны быть запущены 4 контейнера:
- `monkey_team-frontend-1` (порт 80)
- `monkey_team-backend-1` (порт 3000)
- `monkey_team-postgres-1` (порт 5432)
- `monkey_team-redis-1` (порт 6379)

### 5. Доступ к приложению

Откройте браузер и перейдите по адресу:

**http://localhost**

## Структура проекта

```
Monkey_Team/
├── frontend/               # React приложение
│   ├── src/
│   │   ├── components/    # React компоненты
│   │   ├── pages/         # Страницы
│   │   ├── services/      # API и WebSocket сервисы
│   │   ├── store/         # Redux store
│   │   └── types/         # TypeScript типы
│   ├── Dockerfile         # Frontend Docker образ
│   └── nginx.conf         # Nginx конфигурация
│
├── backend/               # Go приложение
│   ├── cmd/
│   │   └── app/          # Точка входа
│   ├── internal/
│   │   ├── config/       # Конфигурация
│   │   ├── delivery/     # HTTP handlers
│   │   ├── repository/   # Database layer
│   │   ├── service/      # Business logic
│   │   └── server/       # HTTP server
│   ├── migrations/       # SQL миграции
│   └── Dockerfile        # Backend Docker образ
│
├── docker-compose.yml    # Docker Compose конфигурация
├── .env                  # Переменные окружения (не коммитится)
└── .env.example          # Пример переменных окружения
```

## Переменные окружения

### База данных (PostgreSQL)
```env
DB_HOST=postgres
DB_PORT=5432
DB_NAME=warehouse_db
DB_USER=warehouse_user
DB_PASSWORD=secure_password
DATABASE_URL=postgresql://warehouse_user:secure_password@postgres:5432/warehouse_db?sslmode=disable
```

### Redis
```env
REDIS_URL=redis://redis:6379
```

### Backend
```env
JWT_SECRET=your_jwt_secret_key_here
```

### GigaChat AI (опционально)
```env
GIGACHAT_CLIENT_ID=your_client_id
GIGACHAT_CLIENT_SECRET=your_client_secret
GIGACHAT_SCOPE=GIGACHAT_API_PERS
AI_SERVICE=sber
```

### Frontend
```env
VITE_API_URL=http://localhost:3000/api
VITE_WS_URL=ws://localhost:3000
```

## Разработка

### Запуск в режиме разработки

Для разработки frontend:
```bash
cd frontend
npm install
npm run dev
```

Для разработки backend:
```bash
cd backend
go mod download
go run cmd/app/main.go
```

### Остановка всех сервисов

```bash
docker-compose down
```

### Остановка с удалением volumes (БД будет очищена)

```bash
docker-compose down -v
```

### Пересборка после изменений в коде

```bash
docker-compose down
docker-compose build
docker-compose up -d
```

### Пересборка конкретного сервиса

```bash
docker-compose build frontend
docker-compose up -d frontend
```

## Просмотр логов

Все сервисы:
```bash
docker-compose logs -f
```

Конкретный сервис:
```bash
docker-compose logs -f backend
docker-compose logs -f frontend
```

## База данных

### Применение миграций

Миграции применяются автоматически при первом запуске PostgreSQL контейнера.

Миграции находятся в `backend/migrations/`:
- `000001_init.up.sql` - создание таблиц
- `000002_seed_data.up.sql` - начальные данные

### Подключение к БД

Через Docker:
```bash
docker exec -it monkey_team-postgres-1 psql -U warehouse_user -d warehouse_db
```

Через локальный клиент:
```
Host: localhost
Port: 5432
Database: warehouse_db
User: warehouse_user
Password: secure_password
```

### Полезные SQL команды

```sql
-- Список таблиц
\dt

-- Структура таблицы
\d inventory_history

-- Проверка данных
SELECT * FROM robots LIMIT 5;
SELECT * FROM products LIMIT 5;
```

## API Endpoints

### Аутентификация
- `POST /api/auth/sign-up` - Регистрация
- `POST /api/auth/login` - Вход

### Роботы
- `POST /api/robots/data` - Отправка данных от робота
- `GET /api/ws/dashboard` - WebSocket для real-time обновлений

### Инвентаризация
- `POST /api/inventory/import` - Импорт данных из CSV
- `GET /api/inventory/history` - История инвентаризации
- `GET /api/export/excel` - Экспорт в Excel

### Dashboard
- `GET /api/dashboard/current` - Текущая информация склада

### AI прогнозы
- `POST /api/ai/predict` - Получить AI прогнозы остатков

## Устранение неполадок

### Порт уже занят

Если порт 80 или 3000 уже используется, измените в `docker-compose.yml`:

```yaml
frontend:
  ports:
    - "8080:80"  # Вместо 80:80

backend:
  ports:
    - "3001:3000"  # Вместо 3000:3000
```

Не забудьте обновить `VITE_API_URL` в `.env`!

### Контейнеры не запускаются

Проверьте логи:
```bash
docker-compose logs
```

Пересоберите без кэша:
```bash
docker-compose build --no-cache
docker-compose up -d
```

### Ошибка подключения к БД

Убедитесь, что PostgreSQL контейнер запущен и здоров:
```bash
docker ps
```

Проверьте логи PostgreSQL:
```bash
docker-compose logs postgres
```

### Проблемы с Docker Desktop на Windows

1. Убедитесь, что WSL2 включен
2. Перезапустите Docker Desktop
3. Очистите Docker кэш: `docker system prune -a`

## Архитектура системы

```
┌─────────────┐
│   Browser   │
└──────┬──────┘
       │ HTTP/WS
       ▼
┌─────────────┐
│   Nginx     │ (Frontend + Proxy)
│  Port: 80   │
└──────┬──────┘
       │ /api → backend:3000
       ▼
┌─────────────┐      ┌──────────────┐
│   Backend   │─────▶│  PostgreSQL  │
│ Go + Gin    │      │   Port 5432  │
│  Port 3000  │      └──────────────┘
└──────┬──────┘
       │
       ▼
┌─────────────┐
│    Redis    │
│  Port 6379  │
└─────────────┘
```

## Тестовые данные

После первого запуска в БД будут созданы тестовые данные:
- Пользователи
- Роботы
- Товары
- История инвентаризации

Проверьте файл `backend/migrations/000002_seed_data.up.sql`

## Производственное развёртывание

Для production:

1. Измените `JWT_SECRET` на криптостойкий случайный ключ
2. Смените пароли БД
3. Настройте HTTPS через reverse proxy (nginx/traefik)
4. Включите логирование и мониторинг
5. Настройте backup БД

## Лицензия

Proprietary - Ростелеком © 2025

## Поддержка

При возникновении проблем создайте issue в репозитории или свяжитесь с командой разработки.
