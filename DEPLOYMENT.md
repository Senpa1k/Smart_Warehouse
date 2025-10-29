# Production Deployment - Smart Warehouse

## Требования к серверу

### Минимальные требования
- **CPU**: 2 ядра
- **RAM**: 4 GB
- **Disk**: 20 GB SSD
- **OS**: Ubuntu 20.04+ / Debian 11+ / CentOS 8+
- **Docker**: 20.10+
- **Docker Compose**: 2.0+

### Рекомендуемые для production
- **CPU**: 4 ядра
- **RAM**: 8 GB
- **Disk**: 50 GB SSD
- **Backup**: автоматический backup БД

---

## Подготовка сервера

### 1. Установка Docker

```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Установка Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Проверка
docker --version
docker-compose --version
```

### 2. Настройка firewall

```bash
# Откройте необходимые порты
sudo ufw allow 22/tcp      # SSH
sudo ufw allow 80/tcp      # HTTP
sudo ufw allow 443/tcp     # HTTPS
sudo ufw enable
```

---

## Развертывание приложения

### 1. Клонирование репозитория

```bash
# Создайте директорию для проекта
sudo mkdir -p /opt/smart-warehouse
sudo chown $USER:$USER /opt/smart-warehouse
cd /opt/smart-warehouse

# Клонируйте репозиторий
git clone <repository-url> .
```

### 2. Конфигурация .env

```bash
cp .env.example .env
nano .env
```

**Обязательно измените**:

```env
# Сгенерируйте криптостойкий JWT секрет
JWT_SECRET=$(openssl rand -base64 32)

# Смените пароль БД
DB_PASSWORD=$(openssl rand -base64 24)

# Обновите DATABASE_URL с новым паролем
DATABASE_URL=postgresql://warehouse_user:NEW_PASSWORD@postgres:5432/warehouse_db?sslmode=require

# Включите SSL для БД
SSL_MODE=require

# Настройте GigaChat (если используете AI)
GIGACHAT_CLIENT_ID=your_production_client_id
GIGACHAT_CLIENT_SECRET=your_production_client_secret

# Production URLs (замените на ваш домен)
VITE_API_URL=https://yourdomain.com/api
VITE_WS_URL=wss://yourdomain.com
```

### 3. Настройка Docker Compose для production

Создайте `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  postgres:
    restart: always
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      POSTGRES_PASSWORD: ${DB_PASSWORD}

  redis:
    restart: always

  backend:
    restart: always
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  frontend:
    restart: always
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

volumes:
  postgres_data:
```

### 4. Запуск

```bash
# Соберите образы
docker-compose -f docker-compose.yml -f docker-compose.prod.yml build

# Запустите сервисы
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Проверьте статус
docker ps
docker-compose logs -f
```

---

## Настройка HTTPS (Nginx + Let's Encrypt)

### 1. Установка Nginx на хост

```bash
sudo apt update
sudo apt install nginx certbot python3-certbot-nginx
```

### 2. Конфигурация Nginx

Создайте `/etc/nginx/sites-available/smart-warehouse`:

```nginx
server {
    listen 80;
    server_name yourdomain.com www.yourdomain.com;

    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }

    location / {
        return 301 https://$server_name$request_uri;
    }
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com www.yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    # SSL настройки
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;

    location / {
        proxy_pass http://localhost:80;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # WebSocket support
    location /socket.io {
        proxy_pass http://localhost:80;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
    }
}
```

### 3. Получение SSL сертификата

```bash
# Активируйте конфигурацию
sudo ln -s /etc/nginx/sites-available/smart-warehouse /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx

# Получите сертификат
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# Настройте автообновление
sudo certbot renew --dry-run
```

---

## Backup и восстановление

### Автоматический backup БД

Создайте скрипт `/opt/smart-warehouse/backup.sh`:

```bash
#!/bin/bash
BACKUP_DIR="/opt/backups/smart-warehouse"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

# Backup PostgreSQL
docker exec monkey_team-postgres-1 pg_dump -U warehouse_user warehouse_db | gzip > $BACKUP_DIR/db_$DATE.sql.gz

# Удалить старые backup (старше 7 дней)
find $BACKUP_DIR -name "db_*.sql.gz" -mtime +7 -delete

echo "Backup completed: $DATE"
```

Настройте cron:

```bash
chmod +x /opt/smart-warehouse/backup.sh
crontab -e

# Добавьте строку (backup каждую ночь в 2:00)
0 2 * * * /opt/smart-warehouse/backup.sh >> /var/log/warehouse-backup.log 2>&1
```

### Восстановление из backup

```bash
# Остановите приложение
cd /opt/smart-warehouse
docker-compose down

# Восстановите БД
gunzip < /opt/backups/smart-warehouse/db_YYYYMMDD_HHMMSS.sql.gz | \
  docker exec -i monkey_team-postgres-1 psql -U warehouse_user -d warehouse_db

# Запустите приложение
docker-compose up -d
```

---

## Мониторинг

### Проверка здоровья сервисов

```bash
# Статус контейнеров
docker ps

# Использование ресурсов
docker stats

# Логи
docker-compose logs -f --tail=100
```

### Установка monitoring stack (опционально)

```bash
# Prometheus + Grafana для мониторинга
# Настройте согласно вашим требованиям
```

---

## Обновление приложения

```bash
cd /opt/smart-warehouse

# Backup перед обновлением
./backup.sh

# Получите новый код
git pull origin main

# Пересоберите и перезапустите
docker-compose -f docker-compose.yml -f docker-compose.prod.yml build
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Проверьте логи
docker-compose logs -f
```

---

## Безопасность

### Чеклист безопасности

- [ ] Измените все дефолтные пароли
- [ ] Используйте SSL/TLS для всех соединений
- [ ] Настройте firewall (только 80, 443, 22)
- [ ] Отключите root SSH доступ
- [ ] Настройте fail2ban
- [ ] Регулярные backup
- [ ] Обновляйте систему и Docker
- [ ] Используйте secrets manager для чувствительных данных
- [ ] Настройте логирование и мониторинг
- [ ] Ограничьте доступ к Docker socket

---

## Troubleshooting

### Проблемы с памятью

```bash
# Очистка Docker
docker system prune -a

# Увеличьте swap
sudo fallocate -l 4G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

### Проблемы с БД

```bash
# Перезапустите только БД
docker-compose restart postgres

# Проверьте логи БД
docker-compose logs postgres
```

### Откат к предыдущей версии

```bash
# Вернитесь к предыдущему коммиту
git log
git checkout <commit-hash>

# Пересоберите
docker-compose build
docker-compose up -d
```

---

## Контакты и поддержка

При возникновении проблем:
1. Проверьте логи: `docker-compose logs`
2. Проверьте документацию: [README.md](README.md)
3. Создайте issue в репозитории
4. Свяжитесь с командой разработки
