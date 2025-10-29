# ================================
# Build Stage
# ================================
FROM node:20-alpine AS build

WORKDIR /app

# Копируем файлы зависимостей
COPY package*.json ./

# Устанавливаем зависимости
RUN npm ci --prefer-offline --no-audit && \
    npm cache clean --force

# Копируем весь проект
COPY . .

# Аргументы сборки для переменных окружения
ARG VITE_API_URL=/api
ENV VITE_API_URL=$VITE_API_URL

# Сборка приложения
RUN npm run build

# ================================
# Production Stage
# ================================
FROM nginx:1.25-alpine

# Метаданные образа
LABEL maintainer="Smart Warehouse Team"
LABEL description="Frontend for Smart Warehouse Application"

# Копируем собранное приложение из build stage
COPY --from=build /app/dist /usr/share/nginx/html

# Копируем конфигурации nginx
COPY nginx-main.conf /etc/nginx/nginx.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Создаем директории для работы nginx без root прав
RUN mkdir -p /var/cache/nginx/client_temp && \
    mkdir -p /tmp/nginx/client_body && \
    mkdir -p /tmp/nginx/proxy && \
    mkdir -p /tmp/nginx/fastcgi && \
    mkdir -p /tmp/nginx/uwsgi && \
    mkdir -p /tmp/nginx/scgi && \
    chown -R nginx:nginx /var/cache/nginx && \
    chown -R nginx:nginx /usr/share/nginx/html && \
    chown -R nginx:nginx /tmp/nginx && \
    chmod -R 755 /usr/share/nginx/html && \
    touch /tmp/nginx.pid && \
    chown nginx:nginx /tmp/nginx.pid

# Переключаемся на непривилегированного пользователя
USER nginx

# Открываем порт 80
EXPOSE 80

# Healthcheck для проверки доступности приложения
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost/ || exit 1

# Запускаем nginx
CMD ["nginx", "-g", "daemon off;"]
