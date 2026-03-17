# Auth Service — сервис регистрации и входа

##  Быстрый старт

### 1. Подготовка
```bash
# Скопировать пример .env
cp .env.example .env

# Отредактировать .env под себя (пароль БД, секреты и т.д.)
```

### 2. Запуск
```bash
# Скачать зависимости
go mod tidy

# Запустить сервер
go run main.go
```
Сервер запустится на `http://localhost:8080`

---

##  Команды для работы

### Зависимости
```bash
# Установить всё, что нужно
go mod tidy

# Добавить новый пакет
go get github.com/какой-то/пакет

# Обновить все зависимости
go get -u ./...
```

### Сборка
```bash
# Собрать бинарник
go build -o auth-service

# Запустить собранный бинарник
./auth-service
```

## Docker

### Собрать образ
```bash
docker build -t auth-service .
```

### Запустить контейнер
```bash
docker run -p 8080:8080 --env-file .env auth-service
```

### Остановить
```bash
docker ps # найти CONTAINER ID
docker stop <CONTAINER_ID>
```


## Переменные окружения (.env)

```ini
# База данных
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=123
DB_NAME=auth_service

# JWT токены
JWT_ACCESS_SECRET=секрет1
JWT_REFRESH_SECRET=секрет2
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h

# Сервер
SERVER_PORT=8080
```



##  Что внутри

- **Регистрация** — проверка пароля, хэш в БД, выдача двух токенов
- **Логин** — проверка email/пароля, новые токены
- **Refresh** — обновление пары по refresh токену
- **Logout** — выход (клиент сам удаляет токены)
- **Me** — данные пользователя по access токену
- **Middleware** — защита маршрутов, прокидывание user_id в контекст

---

