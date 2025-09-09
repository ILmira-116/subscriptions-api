# Subscriptions API

Сервис для управления онлайн-подписками пользователей.

---

## Технологии и стек

- **Язык:** Go (Golang)  
- **HTTP:** Chi  
- **Документация:** Swagger (сгенерирована автоматически)  
- **База данных:** PostgreSQL (с миграциями)  
- **Логирование:** кастомный логгер  
- **Конфигурация:** `.env` 
- **Контейнеризация:** Docker + Docker Compose  

---

## Основные методы API

- `GET /subscriptions` — список подписок с пагинацией  
- `POST /subscriptions` — создать подписку  
- `GET /subscriptions/{id}` — получить подписку по UUID  
- `PUT /subscriptions/{id}` — обновить подписку  
- `DELETE /subscriptions/{id}` — удалить подписку  
- `GET /subscriptions/summary` — подсчёт суммарной стоимости подписок  

> Подробные схемы запросов и ответов доступны в Swagger UI

---

## Настройка и запуск

1. Клонируйте репозиторий:  
```bash
git clone <репозиторий>
cd subscriptions-api
Создайте .env или config.yaml с настройками базы данных, логгера и порта.

Применение миграций

Перед запуском сервера необходимо накатить миграции:

docker-compose run --rm app ./subscriptions-api -migrate
Запуск HTTP-сервера

После применения миграций или если база уже готова, запускаем сервер:

docker-compose run --rm app ./subscriptions-api -serve

Приложение поднимет HTTP-сервер на порту 8080.

Также можно поднять оба контейнера сразу, но миграции при этом нужно запускать отдельно:

docker-compose up -d

2. Запустите сервис через Docker Compose:
```bash
Копировать код
docker-compose up --build

Swagger UI доступен на:

```bash
http://localhost:8080/swagger/index.html

Пример запроса на создание подписки

{
  "service": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025"
}
