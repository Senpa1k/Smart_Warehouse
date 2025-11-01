
FROM node:20-alpine AS build

WORKDIR /app

COPY package*.json ./

RUN if [ -f package-lock.json ]; then \
        npm ci --prefer-offline --no-audit; \
    else \
        npm install --prefer-offline --no-audit; \
    fi && \
    npm cache clean --force

COPY . .

ARG VITE_API_URL=/api
ENV VITE_API_URL=$VITE_API_URL


RUN npm run build


FROM nginx:1.25-alpine

LABEL maintainer="Smart Warehouse Team"
LABEL description="Frontend for Smart Warehouse Application"

COPY --from=build /app/dist /usr/share/nginx/html


COPY nginx-main.conf /etc/nginx/nginx.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf

RUN mkdir -p /var/cache/nginx/client_temp \
    /tmp/nginx/client_body \
    /tmp/nginx/proxy \
    /tmp/nginx/fastcgi \
    /tmp/nginx/uwsgi \
    /tmp/nginx/scgi && \
    chown -R nginx:nginx /var/cache/nginx \
    /usr/share/nginx/html \
    /tmp/nginx && \
    chmod -R 755 /usr/share/nginx/html && \
    touch /tmp/nginx.pid && \
    chown nginx:nginx /tmp/nginx.pid

USER nginx

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost/ || exit 1

CMD ["nginx", "-g", "daemon off;"]
