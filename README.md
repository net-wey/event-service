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

## Makefile

```bash
make up       # Запуск
make down     # Остановка
make clean    # Остановка + удаление volume
make restart  # Перезапуск приложения
```
