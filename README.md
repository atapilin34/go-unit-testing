# URL Shortener

Учебный HTTP-сервис на Go для сокращения длинных URL и обратного перехода по короткому идентификатору.

## Возможности

- Сокращение URL через `POST /shorten`.
- Генерация короткого ID длиной 8 символов.
- Использование букв латинского алфавита в обоих регистрах и цифр.
- Проверка URL с поддержкой только `http` и `https`.
- Защита общего хранилища через `sync.RWMutex`.
- Проверка коллизий коротких идентификаторов.
- Редирект `302 Found` через `GET /{shortID}`.
- Ответ `404 Not Found` для неизвестного ID.
- Unit-тесты бизнес-логики.
- HTTP-тесты обработчиков через `httptest`.
- Table-driven тесты с использованием `t.Run`.

## Требования

- Go 1.22 или новее.

Go 1.22+ нужен для маршрутов с методом и wildcard-параметром:

```go
POST /shorten
GET /{shortID}
```

а также для получения параметра пути:

```go
r.PathValue("shortID")
```

## Структура проекта

```text
.
├── go.mod
├── main.go
├── shortener.go
├── shortener_test.go
├── handlers_test.go
└── README.md
```

## Запуск

Если файл `go.mod` ещё не создан:

```powershell
go mod init example
```

Запустить сервис:

```powershell
go run .
```

После запуска сервер будет доступен по адресу:

```text
http://localhost:8080
```

Ожидаемое сообщение:

```text
server started on http://localhost:8080
```

## API

### Создание короткой ссылки

```http
POST http://localhost:8080/shorten
Content-Type: application/json
```

Тело запроса:

```json
{
  "url": "https://example.com/very/long/path?user=123"
}
```

Успешный ответ:

```http
201 Created
Content-Type: application/json
```

```json
{
  "short_url": "aB7xQ2mK",
  "original_url": "https://example.com/very/long/path?user=123"
}
```

`short_url` — это короткий идентификатор. Для перехода по ссылке добавьте его к адресу сервиса:

```text
http://localhost:8080/aB7xQ2mK
```

### Получение оригинального URL

```http
GET http://localhost:8080/aB7xQ2mK
```

При существующем идентификаторе сервис возвращает:

```http
302 Found
Location: https://example.com/very/long/path?user=123
```

Если идентификатор не найден:

```http
404 Not Found
```

## Проверка в Postman

### POST-запрос

1. Запустите приложение командой `go run .`.
2. Создайте запрос с методом `POST`.
3. Укажите адрес:

```text
http://localhost:8080/shorten
```

4. Откройте **Body → raw → JSON**.
5. Передайте:

```json
{
  "url": "https://example.com/very/long/path"
}
```

6. Нажмите **Send**.
7. Скопируйте `short_url` из ответа.

### GET-запрос

Создайте запрос:

```text
GET http://localhost:8080/<short_url>
```

Например:

```text
GET http://localhost:8080/aB7xQ2mK
```

Postman может автоматически перейти по ответу `302`. Чтобы увидеть сам редирект, отключите настройку **Automatically follow redirects**.

## Проверка через PowerShell

Создать короткую ссылку:

```powershell
$body = @{
    url = "https://example.com/very/long/path?user=123"
} | ConvertTo-Json

$response = Invoke-RestMethod `
    -Method Post `
    -Uri "http://localhost:8080/shorten" `
    -ContentType "application/json" `
    -Body $body

$response
```

Проверить переход:

```powershell
Invoke-WebRequest `
    -Uri "http://localhost:8080/$($response.short_url)" `
    -MaximumRedirection 0
```

## Тестирование

Запустить все тесты:

```powershell
go test -v
```

Запустить только тесты из конкретного файла вместе с исходниками:

```powershell
go test -v handlers_test.go main.go shortener.go
```

Запустить проверку гонок:

```powershell
go test -race -v
```

Проверить покрытие:

```powershell
go test -cover
```

Создать профиль покрытия:

```powershell
go test -coverprofile=coverage.out
```

Посмотреть покрытие по функциям:

```powershell
go tool cover -func=coverage.out
```

Открыть HTML-отчёт:

```powershell
go tool cover -html=coverage.out
```

## Что проверяют тесты

### `shortener_test.go`

Проверяет бизнес-логику:

- валидный HTTP URL;
- валидный HTTPS URL;
- URL с путём и query-параметрами;
- невалидный URL;
- пустую строку;
- генерацию короткого ID;
- получение оригинального URL;
- ошибку при неизвестном ID.

### `handlers_test.go`

Проверяет HTTP-слой через `httptest`:

- успешный `POST /shorten`;
- ошибку при невалидном URL;
- ошибку при пустом URL;
- ошибку при некорректном JSON;
- успешный `GET /{shortID}`;
- статус `302`;
- заголовок `Location`;
- статус `404` для неизвестного ID.

## Хранение данных

Ссылки хранятся в памяти процесса в структуре:

```go
map[string]string
```

Доступ к карте защищён через:

```go
sync.RWMutex
```

`RLock` используется для чтения, а `Lock` — для записи и проверки уникальности ID.

Текущая реализация имеет ограничения:

- данные удаляются после перезапуска приложения;
- ссылки не сохраняются в базе данных;
- несколько экземпляров сервиса не используют общее хранилище;
- отсутствует срок действия ссылки.

Для production-версии следует использовать базу данных, например PostgreSQL или MySQL, и добавить уникальный индекс для `short_id`.

## Генерация ID

Короткий идентификатор состоит из символов:

```text
abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789
```

Для получения случайных данных используется `crypto/rand`. После генерации ID проверяется его наличие в `map`. Если такой ID уже существует, создаётся новый, поэтому существующая ссылка не перезаписывается.

## Валидация URL

Поддерживаются только абсолютные URL со схемами:

```text
http://
https://
```

Примеры корректных URL:

```text
http://example.com
https://google.com/search?q=test
https://example.com:8080/api/items
```

Примеры некорректных значений:

```text
example.com
not-a-url

```

## Лицензия

Учебный проект. Лицензия не задана.
