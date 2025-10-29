# ================================
# Build Stage
# ================================
FROM node:20-alpine AS build

WORKDIR /app

# Копируем package.json и package-lock.json для кэширования слоя с зависимостями
COPY package.json package-lock.json ./

# Устанавливаем зависимости (этот слой будет закэширован, если package.json не изменился)
RUN npm ci --only=production=false && \
    npm cache clean --force

# Копируем остальные файлы проекта
COPY . .

# Аргументы сборки для переменных окружения
ARG VITE_API_URL=/api
ENV VITE_API_URL=$VITE_API_URL

# Сборка приложения
RUN npm run build

# Удаляем исходники и зависимости разработки для уменьшения размера
RUN rm -rf src node_modules

# ================================
# Production Stage
# ================================
FROM nginx:1.25-alpine

# Метаданные образа
LABEL maintainer="Smart Warehouse Team"
LABEL description="Frontend for Smart Warehouse Application"

# Копируем собранное приложение из build stage
COPY --from=build /app/dist /usr/share/nginx/html

# Копируем конфигурацию nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Создаем директорию для PID файла nginx
RUN mkdir -p /var/cache/nginx/client_temp && \
    chown -R nginx:nginx /var/cache/nginx && \
    chown -R nginx:nginx /usr/share/nginx/html && \
    chmod -R 755 /usr/share/nginx/html

# Переключаемся на непривилегированного пользователя
USER nginx

# Открываем порт 80
EXPOSE 80

# Healthcheck для проверки доступности приложения
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost/ || exit 1

# Запускаем nginx
CMD ["nginx", "-g", "daemon off;"]
