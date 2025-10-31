# Быстрый старт - Smart Warehouse

## За 3 минуты до запуска

### 1. Установите Docker Desktop
```bash
# Проверьте установку
docker --version
docker-compose --version
```

### 2. Клонируйте и запустите
```bash
git clone <repository-url>
cd Monkey_Team

# Скопируйте конфигурацию
cp .env.example .env

# Запустите всё одной командой
docker-compose up -d
```

### 3. Откройте приложение
Браузер: **http://localhost**

---

## Основные команды

```bash
# Запустить
docker-compose up -d

# Остановить
docker-compose down

# Посмотреть логи
docker-compose logs -f

# Статус контейнеров
docker ps

# Пересобрать после изменений
docker-compose down && docker-compose build && docker-compose up -d
```

## Структура портов

- **Frontend**: http://localhost (порт 80)
- **Backend API**: http://localhost:3000
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

## Подключение к БД

```bash
docker exec -it monkey_team-postgres-1 psql -U warehouse_user -d warehouse_db
```

Или через клиент:
- Host: localhost
- Port: 5432
- Database: warehouse_db
- User: warehouse_user
- Password: secure_password

## Что дальше?

Читайте полную документацию в [README.md](README.md)
