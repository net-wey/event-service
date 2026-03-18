# Event Service

REST API сервис для управления мероприятиями, участниками и площадками.

## Технологии

- Go 1.25, chi v5
- PostgreSQL 18.2
- Docker, Docker Compose
- Swagger (swaggo)

## Быстрый старт

```bash
make up
```

## Swagger UI

После запуска документация доступна по адресу:

```
http://localhost:8080/swagger/index.html
```

## Переменные окружения

| Переменная   | Описание             | По умолчанию |
|--------------|----------------------|-------------|
| APP_PORT     | Порт приложения      | 8080        |
| DB_HOST      | Хост БД              | db          |
| DB_PORT      | Порт БД              | 5432        |
| DB_USER      | Пользователь БД      | postgres    |
| DB_PASSWORD  | Пароль БД            | (из .env)   |
| DB_NAME      | Имя базы данных      | eventdb     |

## Makefile

```bash
make up       # Запуск
make down     # Остановка
make clean    # Остановка + удаление volume
make restart  # Перезапуск приложения
```
