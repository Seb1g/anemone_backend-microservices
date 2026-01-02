# Anemone Backend
```text

    XXXX   XXX   XX   XXXXXXXX   XXX   XXX    XXXXX     XXX    XX   XXXXXXXX
   XX XX   XXXX  XX   XX         XXXX XXXX   XX    XX   XXXX   XX   XX      
XXX   XX   XX XX XX   XXXXXXXX   XX XXX XX   XX    XX   XX XX  XX   XXXXXXXX
 XXXXXXXX  XX  XXXX   XXXXXXXX   XX  X  XX   XX    XX   XX  XX XX   XXXXXXXX
XX    XX   XX   XXX   XX         XX     XX   XX    XX   XX   XXXX   XX      
XX    XX   XX    XX   XXXXXXXX   XX     XX    XXXXX     XX    XXX   XXXXXXXX
                                                                            
                                                                            
                                                                            
   XX    XX   XXXXXXX   XX        XXXXXX    XXXXXXX   XX   XX   XXXXXXX     
   XX    XX   XXXXXXX   XX        XXXXXXX   XXXXXXX   XXX  XX   XXXXXXX     
   XXXXXXXX   XX   XX   XX        XX   XX     XXX     XXXX XX   XX          
   XXXXXXXX   XX   XX   XX        XX  XXX     XXX     XX XXXX   XX   XX     
   XX    XX   XXXXXXX   XXXXXXX   XXXXXXX   XXXXXXX   XX  XXX   XXXXXXX     
   XX    XX   XXXXXXX   XXXXXXX   XXXXX     XXXXXXX   XX   XX   XXXXXXX 

```


Anemone Backend — backend-платформа, развиваемая в сторону микросервисной архитектуры.
Проект написан на Go и использует PostgreSQL как основное хранилище данных.

---

## Текущий статус проекта

Проект находится в активной фазе разработки и рефакторинга.

На данный момент в репозитории присутствуют следующие сервисы:

* **auth** — аутентификация, JWT (access/refresh), middleware авторизации
* **notes** — заметки, папки, soft-delete, ownership-проверки
* **mail** — временные/постоянные email-адреса и inbox. Имеет ownership-проверки
* **catechize (quiz)** — хранение результатов квизов
* **kanban** — трекер задач, похожий на trello. Имеет ownership-проверки

Каждый сервис развивается как изолированный модуль с собственной бизнес-логикой.

---

## Архитектура

Внутри каждого сервиса используется Layered Architecture:

* API / Handlers — HTTP, JSON, валидация, ответы
* Services — бизнес-логика и use cases
* Repositories — работа с базой данных

Общие принципы:

* передача userID через context
* middleware для auth и ownership
* явные интерфейсы между слоями

---

## Структура проекта

```text
anemone-backend-microservices/
├── internal/
│   ├── auth/        # Аутентификация и JWT
│   ├── notes/       # Заметки и папки
│   ├── mail/        # Временная почта
│   ├── catechize/   # Квизы и результаты
│   ├── kanban/      # Канбан (WIP)
│
├── pkg/             # Общие утилиты
├── docker-compose.yaml
├── .env
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

---

## Технологии

* Go
* PostgreSQL
* gorilla/mux
* sqlx
* lib/pq
* godotenv
* golang-jwt

---

## Конфигурация

Каждый сервис использует собственный `.env` файл.
Переменные окружения загружаются при старте сервиса.

Файлы `.env` не должны попадать в репозиторий и добавлены в `.gitignore`.

Экземпляры .env файлов:

Для notes, catechize, kanban сервисов:
```env
DATABASE_URL=postgres://user:password@db-note:5432/note_db?sslmode=disable
PORT=8083
JWT_SECRET=secretkey
CORS_DEV=localhost:5173
CORS_PROD=https://example.com
```

Для auth сервиса:
```env
DATABASE_URL=postgres://user:password@db-auth:5432/auth_db?sslmode=disable
PORT=8080
JWT_SECRET=secretkey
CORS_DEV=localhost:5173
CORS_PROD=https://example.com
ACCESS_SECRET=secretkey
REFRESH_SECRET=secretkey
```

Для mail сервиса:
```env
DATABASE_URL=postgres://user:password@db-mail:5432/mail_db?sslmode=disable
PORT=8082
SMTP_PORT=25
JWT_SECRET=secretkey
CORS_DEV=localhost:5173
CORS_PROD=https://example.com
DOMAIN_NAME=localhost
```

p.s: JWT_SECRET должен быть везде одинаков

---

## Roadmap

* Унифицировать middleware ownership
* Расширить mail-сервис
* Расширить kanban-сервис
* Привести интерфейсы всех сервисов к единому стилю
* Разделить сервисы на отдельные deployable units