

# Цитатник — REST API для управления цитатами на Go

Простой сервис для добавления, получения, фильтрации и удаления цитат через HTTP.

---

## Требования

- Go 1.24 или выше  
- (Опционально) [swag](https://github.com/swaggo/swag) для генерации Swagger документации

---

## Установка и запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/tajimyradov/quotes-api.git
cd quotes-api
````

2. Установите зависимости:

```bash
go mod tidy
```

3. Запустите сервер:

```bash
go run cmd/main.go
```

По умолчанию сервер доступен на `http://localhost:8080`

---

## API

### Добавить цитату

```
POST /quotes
Content-Type: application/json

{
  "author": "Confucius",
  "text": "Life is simple, but we insist on making it complicated."
}
```

### Получить все цитаты (с опциональным фильтром по автору)

```
GET /quotes
GET /quotes?author=Confucius
```

### Получить случайную цитату

```
GET /quotes/random
```

### Удалить цитату по ID

```
DELETE /quotes/{id}
```

---

## Swagger документация

Генерация документации:

```bash
swag init -g handlers/handler.go
```

Swagger UI доступен по адресу:

```
http://localhost:8080/swagger/index.html
```

---

## Тестирование

Запуск юнит-тестов:

```bash
go test ./...
```

---

## Логирование

Сервис логирует действия создания, получения и удаления цитат.