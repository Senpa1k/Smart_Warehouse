# Smart Warehouse API Documentation

## Base URL
- **Docker (Production):** `http://localhost/api` (через nginx прокси)
- **Direct Backend:** `http://localhost:3000/api`

## Authentication

Все защищенные эндпоинты требуют JWT токен в заголовке:
```
Authorization: Bearer <token>
```

---

## 🔐 Auth Endpoints

### 1. Регистрация пользователя
**POST** `/api/auth/sign-up`

**Request Body:**
```json
{
  "name": "Admin User",
  "email": "admin@warehouse.com",
  "password": "admin123",
  "role": "admin"  // "operator" | "admin" | "viewer"
}
```

**Response:**
```json
{
  "id": 2
}
```

---

### 2. Вход в систему
**POST** `/api/auth/login`

**Request Body:**
```json
{
  "email": "admin@warehouse.com",
  "password": "admin123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 2,
    "name": "Admin User",
    "role": "admin"
  }
}
```

---

### 3. Выход из системы
**POST** `/api/auth/logout`

**Headers:** `Authorization: Bearer <token>`

---

## 📊 Dashboard Endpoints

### 4. Получить данные дашборда
**GET** `/api/dashboard/current`

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "robots": [
    {
      "id": "RB-001",
      "status": "active",
      "battery_level": 85,
      "last_update": "2025-10-30T08:00:00Z",
      "current_zone": "A",
      "current_row": 5,
      "current_shelf": 12
    }
  ],
  "recent_scans": [
    {
      "id": 1,
      "robot_id": "RB-001",
      "product_id": "TEL-4567",
      "quantity": 45,
      "zone": "A",
      "row_number": 5,
      "shelf_number": 12,
      "status": "OK",
      "scanned_at": "2025-10-30T08:30:00Z",
      "created_at": "2025-10-30T08:30:00Z",
      "robot": { /* Robot object */ },
      "product": { /* Product object */ }
    }
  ],
  "statistics": {
    "active_robots": 5,
    "total_robots": 5,
    "items_checked_today": 150,
    "critical_items": 2,
    "avg_battery": 71.4
  }
}
```

---

## 🤖 Robot Endpoints

### 5. Отправить данные от робота
**POST** `/api/robots/data`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "robot_id": "RB-001",
  "battery_level": 85,
  "current_zone": "A",
  "current_row": 5,
  "current_shelf": 12,
  "scans": [
    {
      "product_id": "TEL-4567",
      "quantity": 45,
      "zone": "A",
      "row": 5,
      "shelf": 12
    }
  ]
}
```

---

## 📦 Inventory Endpoints

### 6. Получить историю инвентаризации
**GET** `/api/inventory/history`

**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**
- `from` - Дата начала (ISO 8601)
- `to` - Дата окончания (ISO 8601)
- `zone` - Зоны (через запятую): `A,B,C`
- `status` - Статусы (через запятую): `OK,LOW_STOCK,CRITICAL`
- `search` - Поисковый запрос
- `page` - Номер страницы (default: 1)
- `pageSize` - Размер страницы (default: 20)

**Example:**
```
GET /api/inventory/history?zone=A&status=CRITICAL&page=1&pageSize=20
```

**Response:**
```json
{
  "total": 100,
  "items": [
    {
      "id": 1,
      "robot_id": "RB-001",
      "product_id": "TEL-4567",
      "quantity": 45,
      "zone": "A",
      "row_number": 5,
      "shelf_number": 12,
      "status": "CRITICAL",
      "scanned_at": "2025-10-30T08:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 20,
    "totalPages": 5
  }
}
```

---

### 7. Импорт CSV файла
**POST** `/api/inventory/import`

**Headers:**
- `Authorization: Bearer <token>`
- `Content-Type: multipart/form-data`

**Request Body:**
```
FormData with file field
```

**Response:**
```json
{
  "success": 95,
  "failed": 5,
  "errors": [
    "Row 12: Invalid product ID",
    "Row 34: Missing quantity"
  ]
}
```

---

## 🤖 AI Prediction Endpoints

### 8. Получить предсказания AI
**POST** `/api/ai/predict`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "period_days": 7,
  "categories": []
}
```

**Response:**
```json
{
  "predictions": [
    {
      "product_id": "TEL-4567",
      "product_name": "Роутер RT-AC68U",
      "prediction_date": "2025-11-06",
      "days_until_stockout": 3,
      "recommended_order": 50,
      "confidence_score": 0.85
    }
  ],
  "confidence": 0.85
}
```

---

## 📤 Export Endpoints

### 9. Экспорт в Excel
**GET** `/api/export/excel`

**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**
- `ids` - ID записей через запятую: `1,2,3,4,5`

**Response:** Binary file (application/vnd.openxmlformats-officedocument.spreadsheetml.sheet)

---

## 🔌 WebSocket Endpoints

### 10. WebSocket для дашборда
**GET** `/api/ws/dashboard`

**Headers:**
- `Authorization: Bearer <token>`
- `Upgrade: websocket`

**Получаемые сообщения:**
```json
{
  "type": "robot_update",
  "data": {
    "id": "RB-001",
    "battery_level": 84,
    "status": "active"
  }
}
```

```json
{
  "type": "new_scan",
  "data": {
    "robot_id": "RB-001",
    "product_id": "TEL-4567",
    "quantity": 45,
    "status": "OK"
  }
}
```

---

## 📝 Error Responses

Все ошибки возвращаются в формате:
```json
{
  "message": "Error description"
}
```

**Коды ошибок:**
- `400` - Bad Request (неверные данные)
- `401` - Unauthorized (нет токена или токен невалидный)
- `403` - Forbidden (нет прав доступа)
- `404` - Not Found (ресурс не найден)
- `500` - Internal Server Error (ошибка сервера)

---

## 🧪 Примеры использования

### Curl Examples

**1. Логин:**
```bash
curl -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@warehouse.com","password":"admin123"}'
```

**2. Получить данные дашборда:**
```bash
curl -X GET http://localhost/api/dashboard/current \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

**3. Отправить данные от робота:**
```bash
curl -X POST http://localhost/api/robots/data \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "robot_id": "RB-001",
    "battery_level": 85,
    "current_zone": "A",
    "scans": [{"product_id": "TEL-4567", "quantity": 45}]
  }'
```

---

## 📚 TypeScript Types

Все типы данных доступны в: `frontend/src/types/index.ts`

---

## 🔑 Тестовый пользователь

- **Email:** `admin@warehouse.com`
- **Password:** `admin123`
- **Role:** `admin`
