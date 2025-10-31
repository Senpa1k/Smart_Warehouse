# Чеклист перед Git Push

## ✅ Проверьте перед коммитом

### Безопасность
- [ ] `.env` добавлен в `.gitignore`
- [ ] `.env` НЕ закоммичен (проверьте `git status`)
- [ ] Все секретные ключи удалены из кода
- [ ] `.env.example` не содержит реальных credentials

### Документация
- [ ] README.md актуален
- [ ] QUICKSTART.md содержит правильные команды
- [ ] .env.example содержит все нужные переменные
- [ ] Ссылки в документации работают

### Код
- [ ] Проект собирается без ошибок: `docker-compose build`
- [ ] Все контейнеры запускаются: `docker-compose up -d`
- [ ] Frontend доступен на http://localhost
- [ ] Backend API работает
- [ ] Миграции БД применяются автоматически

### Git
- [ ] `.gitignore` настроен правильно
- [ ] Нет больших файлов (>10MB)
- [ ] Нет скомпилированных бинарников
- [ ] Нет node_modules в репозитории

## Команды для проверки

```bash
# Проверка, что .env не попадет в git
git status | grep ".env"  # Должно быть пусто!

# Проверка размера файлов
find . -type f -size +10M

# Тестовый запуск
docker-compose down -v
docker-compose build
docker-compose up -d
docker ps  # Все 4 контейнера должны быть UP

# Тест приложения
curl http://localhost  # Должен вернуть HTML
```

## Git команды для первого push

```bash
# Инициализация (если еще не сделано)
git init

# Добавить все файлы
git add .

# Проверить что добавляется (убедитесь что .env НЕТ в списке!)
git status

# Первый коммит
git commit -m "Initial commit: Smart Warehouse system

- Frontend: React + TypeScript + MUI
- Backend: Go + Gin + PostgreSQL
- Docker Compose для быстрого запуска
- Полная документация и инструкции"

# Привязать к remote репозиторию
git remote add origin <your-repo-url>

# Push
git push -u origin main
```

## Для новых разработчиков

После вашего push, другие разработчики смогут:

```bash
git clone <your-repo-url>
cd Monkey_Team
cp .env.example .env
docker-compose up -d
```

И всё заработает автоматически!

## Важные файлы для коммита

✅ Должны быть в Git:
- README.md
- QUICKSTART.md
- DEPLOYMENT.md
- .env.example
- .gitignore
- docker-compose.yml
- frontend/ (весь код)
- backend/ (весь код)

❌ НЕ должны быть в Git:
- .env (настоящий)
- node_modules/
- dist/
- build/
- *.log
- postgres_data/
- Скомпилированные бинарники
