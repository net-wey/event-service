# Event Service

REST API сервис для управления мероприятиями, участниками и площадками.

## Стек

- Go 1.25
- PostgreSQL 18
- Docker / Docker Compose
- Swagger (swaggo)
- GitHub Actions (CI/CD)

## Структура

- `app/src` - исходный код Go-приложения
- `app/Dockerfile` - Dockerfile приложения
- `db` - инициализация и Dockerfile БД
- `.github/workflows/ci-cd.yml` - CI/CD pipeline
- `.golangci.yml` - конфиг линтера

## Локальный запуск

### 1. Через Docker Compose

```bash
make up
```

Проверка:

- API: `http://localhost:8080/api/health`
- Swagger UI: `http://localhost:8080/swagger/index.html`

Остановить:

```bash
make down
```

Полная очистка (с volume):

```bash
make clean
```

### 2. Локально без Docker (только приложение)

Из директории `app/src`:

```bash
go mod download
go run .
```

Нужны переменные окружения для подключения к БД (см. `app/src/internal/config/config.go`) или запущенная БД из `docker-compose.yml`.

## Тесты и coverage

Из `app/src`:

```bash
go test ./internal/service/... -covermode=atomic -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Порог в CI: **минимум 50%**. Если ниже, pipeline падает.

## Линтер

Используется `golangci-lint` с конфигом `.golangci.yml`.

Локальный запуск (из `app/src`):

```bash
golangci-lint run --config ../../.golangci.yml
```

## CI/CD (GitHub Actions)

Pipeline запускается:

- при `pull_request` (проверки качества)
- при `push` в `main/master` (проверки + публикация образа)

Стадии:

1. `build` - `go build ./...`
2. `lint` - `golangci-lint`
3. `test` - unit-тесты слоя `internal/service` + расчёт coverage + проверка порога 50% + artifact отчёта
4. `docker_build` - сборка Docker-образа с тегом `<branch>-<short_sha>`
5. `docker_push` - push в Docker Hub (только на push в main/master)

## Docker Hub secrets

В репозитории GitHub нужно задать secrets:

- `DOCKER_USERNAME`
- `DOCKER_TOKEN`

Токен рекомендуется создать в Docker Hub (Access Token) и использовать вместо пароля.

## Почему pipeline падает

Pipeline корректно падает, если:

- приложение не собирается (`build`)
- линтер находит ошибки (`lint`)
- тесты не проходят (`test`)
- coverage ниже 50% (`test`)
- Docker-образ не собирается (`docker_build`)
