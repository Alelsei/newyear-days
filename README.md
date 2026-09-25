# newyear-days

Небольшой сервис на Go, который вычисляет количество календарных дней,
оставшихся от заданной даты до ближайшего следующего Нового года (1 января).

Логика расчёта полностью отделена от способа взаимодействия с программой:
пакет [`internal/daysuntil`](internal/daysuntil/daysuntil.go) не зависит ни
от HTTP, ни от какого-либо интерфейса пользователя и может быть протестирован
и переиспользован независимо. HTTP API в [`internal/api`](internal/api/handler.go)
— это лишь тонкая обёртка над этим пакетом.

## Требования

- Go 1.22 или новее ([инструкция по установке](https://go.dev/doc/install)).
- Внешние зависимости отсутствуют — используется только стандартная
  библиотека Go.

## Структура проекта

```
.
├── cmd/server/main.go          # точка входа HTTP-сервиса
├── internal/daysuntil/         # чистая бизнес-логика (без HTTP)
│   ├── daysuntil.go
│   └── daysuntil_test.go       # unit-тесты
├── internal/api/                # HTTP-обвязка над daysuntil
│   ├── handler.go
│   └── handler_test.go         # интеграционный тест HTTP API
└── .github/workflows/ci.yml    # GitHub Actions CI
```

## Сборка

```bash
go build ./...
```

Собрать исполняемый файл сервиса:

```bash
go build -o bin/server ./cmd/server
```

## Запуск

```bash
go run ./cmd/server
```

По умолчанию сервис слушает `:8080`. Порт можно переопределить переменной
окружения `PORT`:

```bash
PORT=9090 go run ./cmd/server
```

## Запуск тестов

```bash
go test ./...
```

С подробным выводом и покрытием:

```bash
go test ./... -v -race -cover
```

Проект содержит:

- **unit-тесты** для функции расчёта дней (`internal/daysuntil`), включая
  границы года, високосный год и даты вокруг 29 февраля;
- **интеграционный тест** HTTP API (`internal/api`), который поднимает
  реальный HTTP-сервер через `httptest.NewServer` и обращается к нему
  настоящими HTTP-запросами.

## Проверка качества кода

```bash
gofmt -l .        # список файлов с нестандартным форматированием (пусто = OK)
go vet ./...      # статический анализ
go mod tidy       # приводит go.mod/go.sum в соответствие с исходным кодом
go mod verify     # проверяет целостность зависимостей
```

Все эти проверки, а также сборка и тесты, автоматически выполняются в CI
(GitHub Actions) при пуше в `main` и при создании/обновлении Pull Request —
см. [`.github/workflows/ci.yml`](.github/workflows/ci.yml).

## HTTP API

### `GET /api/days-until-new-year`

Возвращает количество дней, оставшихся до ближайшего следующего
1 января.

#### Параметры запроса

| Параметр | Обязательный | Формат | Описание

| `date` | нет | `YYYY-MM-DD` | Дата, относительно которой выполняется расчёт. Если не указана — используется текущая дата (UTC).

Формат даты — `YYYY-MM-DD` (ISO 8601 / `2006-01-02` в терминах `time.Parse`
из стандартной библиотеки Go).

#### Успешный ответ — `200 OK`

```json
{
  "date": "2024-12-31",
  "days_until_new_year": 1
}
```

#### Ошибочный запрос — `400 Bad Request`

Возвращается, если параметр `date` присутствует, но не соответствует
формату `YYYY-MM-DD`:

```json
{
  "error": "invalid \"date\" parameter \"not-a-date\": expected format 2006-01-02"
}
```

#### Неподдерживаемый метод — `405 Method Not Allowed`

Возвращается для любых методов, кроме `GET`:

```json
{
  "error": "only GET is supported"
}
```

### Примеры запросов

```bash
# Относительно текущей даты
curl -s http://localhost:8080/api/days-until-new-year

# Относительно конкретной даты
curl -s "http://localhost:8080/api/days-until-new-year?date=2024-12-31"

# Некорректная дата -> 400
curl -si "http://localhost:8080/api/days-until-new-year?date=abc"
```

йцу
