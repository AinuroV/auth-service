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

## Тестирование API в Postman

### 1. Регистрация
**POST** `http://localhost:8080/api/auth/register`
```json
{
    "email": "test@example.com",
    "password": "Test123!@#"
}
```

### 2. Вход в систему
**POST** `http://localhost:8080/api/auth/login`
```json
{
    "email": "test@example.com",
    "password": "Test123!@#"
}
```

### 3. Получение информации о пользователе
**GET** `http://localhost:8080/api/auth/me`  
**Headers:**
```
Authorization: Bearer <вставь_access_token_сюда>
```

### 4. Обновление токенов
**POST** `http://localhost:8080/api/auth/refresh`
```json
{
    "refresh_token": "<вставь_refresh_token_сюда>"
}
```

### 5. Выход из системы
**POST** `http://localhost:8080/api/auth/logout`
```json
{
    "refresh_token": "<вставь_refresh_token_сюда>"
}
```



### Ожидаемые статусы:
- `200 OK` — успех
- `201 Created` — создан новый пользователь
- `400 Bad Request` — неверный формат
- `401 Unauthorized` — не авторизован
- `409 Conflict` — email уже существует
- `404 Not Found` — несуществующий эндпоинт