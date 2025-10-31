# Техническое задание для фронтенд-разработки
## Система мониторинга склада Smart Warehouse

## Общая информация
**Базовый URL:** `http://localhost:3000` (development)  
**Тип аутентификации:** Bearer Token  
**Формат данных:** JSON

---

## Эндпоинты API

### 1. Авторизация пользователя
**Метод:** `POST /api/auth/login`  

**Request Body:**
{
  "email": "string",
  "password": "string"
}
Response Success (200):
{
  "token": "jwt_token_string",
  "user": {
    "id": "stng",
    "name": "striring",
    "role": "admin|operator|viewer"
  }
}


2. WebSocket для real-time обновлений
Метод: `GET /api/ws/dashboard`
Headers: Authorization: Barear jwt_token

{
    "data": {
        "id": "RB-003",
        "status": "active",
        "battery_level": 90,
        "last_update": "2025-10-30T12:40:17.631898Z",
        "current_zone": "A",
        "current_row": 5,
        "current_shelf": 3
    },
    "type": "robot_update"
}
или 
{
  "type": "inventory_alert",
  "data": {
    "product_id": "TEL-4567",
    "product_name": "Роутер RT-AC68U",
    "current_quantity": 5,
    "zone": "A",
    "row": 12,
    "shelf": 3,
    "status": "CRITICAL",
    "alert_type": "scanned", // или "predicted"
    "timestamp": "2025-10-28T14:30:00Z",
    "message": "Критический остаток! Требуется пополнение."
  }
}

4. Получение текущего состояния дашборда
Метод: GET /api/dashboard/current
Headers: Аутентификация: User Token

Response 200
json
{
  "robots": [
    {
        "id": "IMPORT_SERVICE",
        "status": "active",
        "battery_level": 100,
        "last_update": "0001-01-01T00:00:00Z",
        "current_zone": "",
        "current_row": 0,
        "current_shelf": 0
    },
  ],
  "recent_scans": [
    {
        "id": 332,
        "robot_id": "RB-003",
        "product_id": "TEL-4567",
        "quantity": 72,
        "zone": "A",
        "row_number": 4,
        "shelf_number": 4,
        "status": "OK",
        "scanned_at": "2025-10-30T12:41:17.650267Z",
        "created_at": "2025-10-30T12:41:17.663829Z",
        "robot": {
            "id": "RB-003",
            "status": "active",
            "battery_level": 90,
            "last_update": "0001-01-01T00:00:00Z",
            "current_zone": "A",
            "current_row": 5,
            "current_shelf": 4
        },
        "product": {
                "id": "TEL-8901",
                "name": "Модем DSL-2640U",
                "category": "network",
                "min_stock": 5,
                "optimal_stock": 50
            }
        },
  ],
  "statistics": {
        "total_check": 335,
        "unique_products": 5,
        "find _discrepancies": 0,
        "average_time": 0
    }
}

5. Получение исторических данных
Метод: GET /api/inventory/history
Headers: Аутентификация: User Token

Query Parameters:

from (optional): Дата начала периода (YYYY-MM-DD)
to (optional): Дата окончания периода (YYYY-MM-DD)
zone (optional): Фильтр по зоне (A, B, C, D, E)
status (optional): Фильтр по статусу (OK, LOW_STOCK, CRITICAL)

пример localhost:3000/api/inventory/history?from=2025-10-30 12:08:17&to=2025-10-30 12:10:17&zone=A&status=OK

Response 200
json
{
    "total": 0,
    "items": [
        {
            "id": 20,
            "robot_id": "RB-002",
            "product_id": "TEL-8901",
            "quantity": 62,
            "zone": "A",
            "row_number": 1,
            "shelf_number": 2,
            "status": "OK",
            "scanned_at": "2025-10-30T12:09:17.193272Z",
            "created_at": "2025-10-30T12:09:17.20532Z",
            "robot": {
                "id": "RB-002",
                "status": "active",
                "battery_level": 87,
                "last_update": "0001-01-01T00:00:00Z",
                "current_zone": "A",
                "current_row": 6,
                "current_shelf": 2
            },
            "product": {
                "id": "TEL-8901",
                "name": "Модем DSL-2640U",
                "category": "network",
                "min_stock": 5,
                "optimal_stock": 50
            }
        },
    ],
    "pagination": {
        "limit": 50,
        "offset": 0
    }
}


6. Загрузка инвентаризации через CSV
Метод: POST /api/inventory/import
Headers:Аутентификация: User Token
        Content-Type: multipart/form-data

Request Body:

file: CSV файл с данными инвентаризации

Response Success (200):

json
{
  "success": 145,
  "failed": 5,
  "errors": [
    {
      "row": 3,
      "error": "Invalid product ID",
      "data": "INVALID-123"
    }
  ]
}


7. AI-прогнозирование запасов
Метод: POST /api/ai/predict
Headers: Аутентификация: User Token

Request Body:

json
{
  "period_days": 7,
  "categories": ["networking", "cables"],
}


Response Success (200):

json
{
  "predictions": [
    {
            "product_id": "TEL-4567",
            "prediction_date": "01.01.2025",
            "days_until_stockout": 10,
            "recommended_order": 90,
            "confidence_score": 0.8
    },
  ],
  "confidence": 0.85
}

9.
Метод: GET /api/export/excel?ids=1,2,3
Headers: Аутентификация: User Token
Response: Binary file stream