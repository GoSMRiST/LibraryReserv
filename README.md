# 📦 Reservations Service

Сервис бронирования книг в библиотечной системе.
Отвечает за создание и управление резервами книг пользователями.

Реализован на **Go** с использованием **REST** и **gRPC**, построен по принципам **чистой архитектуры**.

---

## 🚀 Возможности

* 📚 Бронирование книг
* ❌ Отмена бронирования
* 📋 Получение списка резервов
* 🔗 Интеграция с другими сервисами (например, Books Service)
* 🔐 Авторизация через middleware
* 🗄️ Работа с БД через repository слой
* Работа с Docker

---

## 🛠️ Технологии

* **Go**
* **gRPC**
* **REST API**
* **PostgreSQL**
* **SQL migrations**
* **Docker**

---

## 📂 Структура проекта

```bash
.
├── cmd/app                # Точка входа (main.go)
├── internal/
│   ├── app/               # Инициализация приложения
│   │   ├── grpc/
│   │   └── rest/
│   ├── config/            # Конфигурация
│   ├── core/              # Модели и бизнес-логика
│   ├── middleware/        # Middleware
│   ├── repository/        # Работа с БД
│   ├── services/          # Бизнес-логика
│   └── transport/         # REST / gRPC обработчики
├── migrations/            # SQL миграции
├── go.mod
└── go.sum
└── Dockerfile
```

---

## ⚙️ Конфигурация

Создайте `.env` файл:

```env
DB_HOST=host.docker.internal
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=DataBasePostgresPass
DB_NAME=BooksReserv
```

---

## 🗄️ База данных

Используется PostgreSQL.

Миграции находятся в папке:

```bash
migrations/
```

---

## ▶️ Запуск

```bash
go run cmd/app/main.go
```

---

## 🌐 API

### REST

Пример эндпоинтов:

* `POST /reservations` — создать бронь
* `GET /reservations` — получить список броней
* `PATCH /reservations/:id/close` — закрыть бронь

---

### gRPC

* Используется для взаимодействия с другими сервисами
* Реализация: `internal/transport/grpc`

---

## 🔐 Аутентификация

Реализована через middleware.
Доступ к эндпоинтам ограничен с использованием токенов.

---

## 🔗 Интеграции

Сервис может взаимодействовать с:

* 📚 **Books Service** — проверка существования и доступности книги

---

## 🧪 Пример запроса

```bash
curl -X POST http://localhost:8080/reservations \
  -H "Content-Type: application/json" \
  -d '{"book_id":1,"user_id":42}'
```

---

## 👤 Автор

* GitHub: https://github.com/GoSMRiST
