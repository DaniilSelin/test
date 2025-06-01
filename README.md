# Quotes API Service

Простой REST API-сервис на Go для хранения и управления цитатами.

## Функциональные требования

Добавление новой цитаты
POST /quotes
Пример:

    curl -X POST http://localhost:8080/quotes \
      -H "Content-Type: application/json" \
      -d '{"author":"Confucius","quote":"Life is simple, but we insist on making it complicated."}'

Получение всех цитат
GET /quotes
Пример:

    curl http://localhost:8080/quotes

Получение случайной цитаты
GET /quotes/random
Пример:

    curl http://localhost:8080/quotes/random

Фильтрация по автору
GET /quotes?author=<имя>
Пример:

    curl "http://localhost:8080/quotes?author=Confucius"

Удаление цитаты по ID
DELETE /quotes/{id}
Пример:

    curl -X DELETE http://localhost:8080/quotes/1

## Запуск

Необходимые зависимости

    Go 1.23+

    PostgreSQL (если запускаете без Docker)

    Docker и Docker Compose (по желанию)

## Конфигурация

В файле config/config.yml задать настройки базы данных и логгера.

## Проект использует слоистую архитектуру:

    database: подключение и миграции (internal/database)

    models: описание модели Quote (internal/models/quote.go)

    repository: реализация работы с БД (internal/repository/quote_repository.go)

    service: бизнес-логика (internal/service/quote_service.go)

    transport/http/api: HTTP-хэндлеры и маршруты (internal/transport/http/api)

    logger: настройка zap-логгера (internal/logger)

    errdefs: стандартные ошибки (internal/errdefs/errdefs.go)

## Запуск локально (без Docker)

Установить зависимости:

    go mod download

Поднять локальную PostgreSQL и создать БД. Посмотреть перед запуском в config/config.yml

Запустить сервис:

    go run cmd/main.go

Либо воспользоваться makefile-ом

    make run
    make stop # остановаить сервер

Чтобы корректно завершить работу сервиса (gracefull shutdown), достаточно отправить ему сигнал SIGINT
    
    sudo kill -SIGINT $(sudo lsof -ti:<порт из конфигурации>)

### Запуск через Docker Compose

В корне проекта:

    docker-compose up -d --build

После старта контейнеров сервис доступен на http://localhost.

## Тесты

Ручные интеграционные тесты через curl (примеры выше).

Юнит-тесты для слоя репозитория (работают с реальной БД).
Для запуска:

    go test ./internal/repository/ -v

Мок-тесты для слоя сервиса (не требуют подключения к БД).
Для запуска:

    go test ./internal/service/ -v
