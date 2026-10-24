# 🏭 Умный Склад (Smart Warehouse)

## 📖 О проекте

Умный Склад — это комплексное веб-решение для автоматизации складских операций. Система объединяет мониторинг автономных роботов в реальном времени, AI-прогнозирование остатков и аналитику исторических данных в едином интерфейсе.

## 🏗 Архитектура системы

```mermaid
flowchart TB
    %% Frontend Layer
    subgraph Frontend["Frontend Layer<br/>React + TypeScript"]
        A[React App] 
        B[MUI Components]
        C[Redux Store]
        D[WebSocket Client]
        
        A --> B
        A --> C
        A --> D
    end

    %% Backend Layer
    subgraph Backend["Backend Layer<br/>Go + Gin"]
        E[API Server]
        F[JWT Auth]
        G[WebSocket Manager]
        H[REST API]
        I[GigaChat Integration]
        
        E --> F
        E --> G
        E --> H
        E --> I
    end

    %% Data Layer
    subgraph Data["Data Layer"]
        J[PostgreSQL<br/>Основная БД]
        K[Redis<br/>Кэш & WebSocket]
    end

    %% External Services
    subgraph External["External Services"]
        L[Robot Emulator<br/>Python]
        M[GigaChat AI<br/>Sber API]
    end

    %% Connections
    A <-->|HTTP/REST| E
    D <-->|WebSocket| G
    E <-->|SQL| J
    E <-->|Pub/Sub| K
    E <-->|REST API| M
    L -->|Robot Data| E

    %% Styling
    classDef frontend fill:#61dafb,color:#000
    classDef backend fill:#00add8,color:#000
    classDef database fill:#336791,color:#fff
    classDef cache fill:#d82c20,color:#fff
    classDef external fill:#3776ab,color:#fff
    classDef ai fill:#ff6b35,color:#fff
    
    class A,B,C,D frontend
    class E,F,G,H,I backend
    class J database
    class K cache
    class L external
    class M ai
```
## 🚀 Быстрый старт
### 1. Клонирование и настройка
```bash
git clone https://github.com/Senpa1k/Smart_Warehouse.git
cd Smart_Warehouse
cp .env_example .env
```
## 2. Запуск приложения
```bash
docker-compose up --build -d
```
## 3. Открытие в браузере
Адрес: http://localhost

Тестовые учетные данные:

```text
Email: test1@gmail.com
Password: test2025
```

## 🛠 Технологический стек
### Backend (Go)
* **Go 1.24 + Gin** - высокопроизводительный API сервер

* **GORM** - работа с PostgreSQL

* **JWT** - аутентификация

* **Redis** - кэш и WebSocket

* **GigaChat API** - AI-прогнозирование

### Frontend (TypeScript)
* **React 18 + TypeScript** - пользовательский интерфейс

* **Vite** - сборка и разработка

* **Material-UI** - компоненты дизайна

* **Redux Toolkit** - управление состоянием

* **WebSocket** - real-time обновления

### Инфраструктура
* **PostgreSQL** - основная БД

* **Docker Compose** - контейнеризация

* **Nginx** - веб-сервер

* **Python** - эмулятор роботов

## 📁 Ключевые компоненты
### Backend структура
```text
backend/
├── cmd/app/                 # Точка входа приложения
├── internal/
│   ├── delivery/http/       # HTTP хендлеры (REST API)
│   ├── service/             # Бизнес-логика
│   ├── repository/          # Доступ к данным (PostgreSQL)
│   └── config/              # Конфигурация
├── migrations/              # Миграции БД
└── go.mod                  # Зависимости
```
### Frontend структура
```text
frontend/
├── src/
│   ├── pages/              # Страницы приложения
│   │   ├── LoginPage.tsx   # Авторизация
│   │   ├── DashboardPage.tsx # Мониторинг
│   │   └── HistoryPage.tsx # История
│   ├── components/         # UI компоненты
│   ├── store/              # Redux store
│   └── services/           # API клиенты
└── package.json           # Зависимости
```
## 🎯 Основные функции
### 🔐 Авторизация
* JWT аутентификация

* Валидация email/пароля

* Разграничение прав доступа

### 📊 Дашборд реального времени
* Интерактивная карта склада с роботами

* Статистика активных роботов и проверок

* Таблица последних сканирований

* WebSocket обновления

### 🤖 AI-аналитика
* Интеграция с GigaChat API

* Прогноз остатков на 7 дней

* Рекомендации по пополнению запасов

### 📈 Исторические данные
* Фильтрация по периоду, зонам, статусам

* Сортировка и пагинация

* Экспорт в Excel

* Визуализация трендов

### 📁 Загрузка данных
* Drag & Drop CSV импорт

* Валидация формата файлов

* Предпросмотр данных

## 🔌 API Endpoints

### Основные эндпоинты

| Метод | Путь | Назначение |
|-------|------|------------|
| `POST` | `/api/auth/login` | Авторизация пользователя |
| `POST` | `/api/robots/data` | Прием данных от роботов |
| `GET` | `/api/dashboard/current` | Текущее состояние склада |
| `GET` | `/api/inventory/history` | Исторические данные |
| `POST` | `/api/inventory/import` | Загрузка CSV файлов |
| `POST` | `/api/ai/predict` | Получение AI-прогнозов |

### WebSocket события

| Событие | Назначение |
|---------|------------|
| `robot_update` | Обновление позиций и статусов роботов |
| `new_scan` | Новые сканирования товаров |
| `inventory_alert` | Предупреждения о критических остатках |

## 🤖 Эмулятор роботов
Автономный эмулятор на Python для тестирования системы:

* **5 виртуальных роботов** с уникальными ID

* **Реалистичное перемещение** по зонам склада (A-E)

* **Сканирование товаров** с разными статусами

* **Симуляция батареи** (расход и зарядка)

* **Настраиваемый интервал** обновления (10 сек)

## 🐳 Управление контейнерами
### Основные команды
```bash
# Запуск всех сервисов
docker-compose up -d

# Просмотр логов
docker-compose logs -f backend
docker-compose logs -f frontend

# Остановка
docker-compose down

# Пересборка
docker-compose up --build -d
```
### Доступные сервисы
* **Frontend:** http://localhost (порт 80)

* **Backend API:** http://localhost:3000

* **PostgreSQL:** localhost:5432

* **Redis:** localhost:6379

## 🔧 Разработка
### Локальный запуск Backend
```bash
cd backend
go mod download
go run cmd/app/main.go
```

### Проверка Backend
Команды выполняются из каталога `backend`. Для `make test` нужен CGO. Линтер `golangci-lint` v2.14.0 требует Go 1.26+; `make install` установит его в стандартный каталог Go (`GOBIN` или `GOPATH/bin`).

```bash
cd backend
make install   # зависимости и закреплённая версия golangci-lint
make format    # gofumpt и goimports
make lint      # tidy, конфигурация, форматирование и статический анализ
make test      # тесты с race detector и отчётом покрытия
make build
```
### Локальный запуск Frontend
```bash
cd frontend
npm install
npm run dev  # http://localhost:3001
```
