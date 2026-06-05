Для данных задач использовалась локальная база данных PostgreSQL с паролем `secret` и портом `5432`
Для данных задач нужно скачать `protoc`: https://github.com/protocolbuffers/protobuf/releases 

# Базовый гараж (Задача 1) `../garage`

Связка **protobuf → кодогенерация → gRPC‑сервер → REST‑шлюз**.

## Что сделано

- Реализован gRPC‑сервис `CarService` с методами `CreateCar` и `GetCar`.
- Через grpc‑gateway автоматически сгенерирован REST API для тех же методов.
- Хранение данных — in‑memory map (данные живут, пока работает сервер).
- Подключена серверная gRPC Reflection — можно вызывать методы без указания `.proto` файла.
- Код разделён на слои: `domain`, `service`, `repository`, `transport` (чистая архитектура).

## Быстрый старт

1. Инструменты (один раз)

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

2. Склонируйте/откройте проект, скачайте зависимости

```bash
go mod tidy
```

3. Генерация кода

```bash
protoc -I api/proto -I third_party \
  --go_out=api/gen --go_opt=paths=source_relative \
  --go-grpc_out=api/gen --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=api/gen --grpc-gateway_opt=paths=source_relative \
  api/proto/garage/v1/car.proto
```

4. Запустить сервер

```bash
go run cmd/server/main.go
```

В консоли появится:
```text
gRPC на :9090
HTTP на :8080
```

## Тестирование

1. gRPC (порт 9090)

- Используйте `grpcurl`. Благодаря включённой рефлексии можно не указывать proto‑файл.

**Создать автомобиль:**
```powershell
grpcurl -plaintext -d '{"vin":"VIN001","brand":"Toyota","year":2020}' localhost:9090 garage.v1.CarService/CreateCar
```

**Ожидаемый ответ:**
```json
{
  "car": {
    "vin": "VIN001",
    "brand": "Toyota",
    "year": 2020
  }
}
```

**Получить автомобиль:**
```powershell
grpcurl -plaintext -d '{"vin":"VIN001"}' localhost:9090 garage.v1.CarService/GetCar
```

Ответ будет таким же.
Повторный вызов с тем же VIN вернёт ошибку `AlreadyExists`.
Вызов с несуществующим VIN вернёт `NotFound`.

2. REST (порт 8080)

Используйте `curl` (доступен в PowerShell).

**Создать автомобиль:**
```powershell
curl -X POST http://localhost:8080/api/v1/cars -H "Content-Type: application/json" -d '{"vin":"VIN002","brand":"Honda","year":2022}'
```

**Ожидаемый ответ:**
```json
{"car":{"vin":"VIN002","brand":"Honda","year":2022}}
```

**Получить автомобиль:**
```powershell
curl http://localhost:8080/api/v1/cars/VIN002
```

Ответ — тот же объект.
При отсутствии автомобиля HTTP‑шлюз вернёт `404 Not Found`.



# Сервис-агрегатор документов (Задача 2) `../documents`

Переход от in‑memory к реляционной БД, работа со связанными сущностями «Пользователь» и «Документ» (One‑to‑Many).  
gRPC + REST через grpc‑gateway, SQL‑запросы генерируются `sqlc`, миграции PostgreSQL.

## Что сделано

- Реализованы gRPC‑методы `CreateUser`, `AddDocumentToUser`, `ListUserDocuments`.
- REST‑эндпоинты:
  - `POST /api/v1/users`
  - `POST /api/v1/users/{user_id}/documents`
  - `GET /api/v1/users/{user_id}/documents`
- Хранение — PostgreSQL.
- Миграции создают таблицы `users` и `documents` с внешним ключом.
- SQL‑запросы полностью генерируются **sqlc**.
- Бизнес‑логика проверяет существование пользователя перед добавлением документа.
- Код разделён на слои (domain, service, repository, transport).
- Включена gRPC Reflection.

## Быстрый старт

1. Инструменты (один раз)

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

2. Склонируйте/откройте проект, скачайте зависимости

```bash
go mod tidy
```

3. Генерация кода

```bash
protoc -I api/proto -I third_party \
  --go_out=api/gen --go_opt=paths=source_relative \
  --go-grpc_out=api/gen --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=api/gen --grpc-gateway_opt=paths=source_relative \
  api/proto/documents/v1/service.proto
```

```bash
sqlc generate
```

## Тестирование

### Запуск PostgreSQL (локально)

```bash
psql -U postgres -c "CREATE DATABASE documents"
psql -U postgres documents < migrations/001_init.up.sql
```

Пароль для БД: `secret`

### Запуск Сервера

```bash
go run cmd/server/main.go
```

1. gRPC (порт 9090)

**Создание пользователя:**
```bash
grpcurl -plaintext -d '{"name":"Alice"}' localhost:9090 documents.v1.DocumentService/CreateUser
```

**Ответ:**
```json
{
  "user": {
    "id":"uuid",
    "name":"Alice"
  }
}
```

**Добавление документа (замените user_id):**
```bash
grpcurl -plaintext -d '{"user_id":"<uuid>", "document":{"title":"Doc1"}}' localhost:9090 documents.v1.DocumentService/AddDocumentToUser
```

**Ответ:**
```json
{
  "document": {
    "id": "udid",
    "title": "Doc1",
    "userId": "uuid"
  }
}
```

**Список документов (замените user_id):**
```bash
grpcurl -plaintext -d '{"user_id":"<uuid>"}' localhost:9090 documents.v1.DocumentService/ListUserDocuments
```

**Ответ:**
```json
{
  "documents": [
    {
      "id": "udid",
      "title": "Doc1",
      "userId": "uuid"
    }
  ]
}
```

2. REST (порт 8080)

**Создание пользователя:**
```bash
curl -X POST http://localhost:8080/api/v1/users -H "Content-Type: application/json" -d '{"name":"Bob"}'
```

**Ответ:**
```json
{"user":{"id":"uuid","name":"Bob"}}
```

**Добавление документа (замените user_id):**
```bash
curl -X POST http://localhost:8080/api/v1/users/<uuid>/documents -H "Content-Type: application/json" -d '{"document":{"title":"Doc2"}}'
```

**Ответ:**
```json
{"document":{"id":"udid","title":"","userId":"uuid"}}
```

**Список документов (замените user_id):**
```bash
curl http://localhost:8080/api/v1/users/<uuid>/documents
```

**Ответ:**
```json
{"documents":[{"id":"udid","title":"","userId":"uuid"}]}
```



# Конвейер автомобиля (Задача 3) `../assembly`

Развитие примера с автомобилем. Есть три сущности: Car (VIN, Brand, Year), Engine (ID, Horsepower), Transmission (ID, Type).
Нужно атомарно связать автомобиль с двигателем и коробкой передач, а также получать полное описание автомобиля с характеристиками двигателя и КПП через вложенную структуру.

## Что сделано

- gRPC‑метод `AssembleCar` – атомарная связка автомобиля с двигателем и КПП в одной транзакции.
- gRPC‑метод `GetCarSpec` – возвращает `CarSpec` с вложенными объектами `Car`, `Engine`, `Transmission`.
- REST‑эндпоинты (через grpc‑gateway):
  - `POST /api/v1/cars/assemble`
  - `GET /api/v1/cars/{vin}/spec`
- Вспомогательные методы `CreateEngine` и `CreateTransmission` для наполнения справочников.
- Хранение – PostgreSQL, миграции.
- Бизнес‑логика проверяет существование пользователя перед добавлением документа.
- Код разделён на слои (domain, service, repository, transport).
- Включена gRPC Reflection.

## Быстрый старт

1. Инструменты (один раз)

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

2. Склонируйте/откройте проект, скачайте зависимости

```bash
go mod tidy
```

3. Генерация кода

```bash
protoc -I api/proto -I third_party \
  --go_out=api/gen --go_opt=paths=source_relative \
  --go-grpc_out=api/gen --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=api/gen --grpc-gateway_opt=paths=source_relative \
  api/proto/assembly/v1/service.proto
```

```bash
sqlc generate
```

## Тестирование

### Запуск PostgreSQL (локально)

```bash
psql -U postgres -c "CREATE DATABASE assembly"
psql -U postgres assembly < migrations/001_init.up.sql
```

Пароль для БД: `secret`

### Запуск Сервера

```bash
go run cmd/server/main.go
```

1. Создаём двигатель и коробку передач

**gRPC:**
```bash
grpcurl -plaintext -d '{"id":"eng1","horsepower":200}' localhost:9090 assembly.v1.AssemblyService/CreateEngine
grpcurl -plaintext -d '{"id":"trans1","type":"Automatic"}' localhost:9090 assembly.v1.AssemblyService/CreateTransmission
```

**Ответ:**
```json
{
  "id": "eng1",
  "horsepower": 200
}
{
  "id": "trans1",
  "type": "Automatic"
}
```

**REST:**
```bash
curl -X POST http://localhost:8080/api/v1/engines -H "Content-Type: application/json" -d '{"id":"eng1","horsepower":200}'
curl -X POST http://localhost:8080/api/v1/transmissions -H "Content-Type: application/json" -d '{"id":"trans1","type":"Automatic"}'
```

**Ответ:**
```json
{"code":2, "message":"pq: повторяющееся значение ключа нарушает ограничение уникальности \"engines_pkey\" (23505)", "details":[]}{"code":2, "message":"pq: повторяющееся значение ключа нарушает ограничение уникальности \"transmissions_pkey\" (23505)", "details":[]}
```

2. Собираем автомобиль

**gRPC:**
```bash
grpcurl -plaintext -d '{"vin":"VIN123","brand":"Toyota","year":2023,"engine_id":"eng1","transmission_id":"trans1"}' localhost:9090 assembly.v1.AssemblyService/AssembleCar
```

**Ответ:**
```json
{
  "spec": {
    "car": {
      "vin": "VIN123",
      "brand": "Toyota",
      "year": 2023
    },
    "engine": {
      "id": "eng1",
      "horsepower": 200
    },
    "transmission": {
      "id": "trans1",
      "type": "Automatic"
    }
  }
}
```

**REST:**
```bash
curl -X POST http://localhost:8080/api/v1/cars/assemble -H "Content-Type: application/json" -d '{"vin":"VIN456","brand":"Honda","year":2022,"engine_id":"eng1","transmission_id":"trans1"}'
```

**Ответ:**
```json
{"spec":{"car":{"vin":"VIN456", "brand":"Honda", "year":2022}, "engine":{"id":"eng1", "horsepower":200}, "transmission":{"id":"trans1", "type":"Automatic"}}}
```

3. Получаем спецификацию

```bash
grpcurl -plaintext -d '{"vin":"VIN123"}' localhost:9090 assembly.v1.AssemblyService/GetCarSpec
curl http://localhost:8080/api/v1/cars/VIN456/spec
```

**Ответ:**
```json
{
  "car": {
    "vin": "VIN123",
    "brand": "Toyota",
    "year": 2023
  },
  "engine": {
    "id": "eng1",
    "horsepower": 200
  },
  "transmission": {
    "id": "trans1",
    "type": "Automatic"
  }
}
{"car":{"vin":"VIN456", "brand":"Honda", "year":2022}, "engine":{"id":"eng1", "horsepower":200}, "transmission":{"id":"trans1", "type":"Automatic"}}
```



# Складской учет (Задача 4) `../warehouse`

Асинхронная обработка заказов с использованием очереди сообщений.
Создание заказа (синхронный gRPC/REST) не резервирует товар сразу, а отправляет сообщение в очередь.
Фоновый воркер читает очередь, имитирует резервирование и обновляет статус заказа.

## Что сделано

- Синхронная часть: `CreateOrder` сохраняет заказ со статусом `PENDING` в БД и отправляет его ID в брокер. Ответ возвращается сразу.
- Асинхронная часть: горутина‑воркер подписывается на очередь, получает ID заказа, эмулирует обработку (задержка 10-12 секунд) и переводит заказ в статус `RESERVED` или `FAILED`.
- Проверка статуса: `GetOrderStatus` просто читает заказ из БД.
- Хранение: PostgreSQL (таблицы `orders` и `products`).
- Очередь: in‑memory каналы Go (реализует интерфейс `MessageBroker`). При желании легко заменяется на RabbitMQ.
- gRPC + REST: стандартная связка через grpc‑gateway.

## Быстрый старт

1. Инструменты (один раз)

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

2. Склонируйте/откройте проект, скачайте зависимости

```bash
go mod tidy
```

3. Генерация кода

```bash
protoc -I api/proto -I third_party \
  --go_out=api/gen --go_opt=paths=source_relative \
  --go-grpc_out=api/gen --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=api/gen --grpc-gateway_opt=paths=source_relative \
  api/proto/warehouse/v1/service.proto
```

## Тестирование

### Запуск PostgreSQL (локально)

```bash
psql -U postgres -c "CREATE DATABASE warehouse"
psql -U postgres warehouse < migrations/001_init.up.sql
```

Пароль для БД: `secret`

### Запуск Сервера

```bash
go run cmd/server/main.go
```

1. REST тестирование

**Создание заказа:**
```bash
curl -X POST http://localhost:8080/api/v1/orders \
     -H "Content-Type: application/json" \
     -d '{"product_id":"prod-1", "quantity":2}'
```

**Ответ:**
```json
{"order_id":"<UUID>", "status":"PENDING"}
```

**Проверка статуса сразу и через 10-12 секунд:**
```bash
curl http://localhost:8080/api/v1/orders/<UUID>/status
```

**Ответ сразу:**
```json
{"orderId":"<UUID>", "status":"PENDING"}
```

**Ответ через 10-12 секунд:**
```json
{"orderId":"<UUID>", "status":"RESERVED"}
```

2. gRPC тестирование

**Создание заказа:**
```bash
grpcurl -plaintext -d '{"product_id":"prod-1","quantity":1}' \
        localhost:9090 warehouse.v1.WarehouseService/CreateOrder
```

**Ответ:**
```json
{
  "orderId": "<UUID>",
  "status": "PENDING"
}
```

**Проверка статуса сразу и через 10-12 секунд:**
```bash
grpcurl -plaintext -d '{"order_id":"<UUID>"}' \
        localhost:9090 warehouse.v1.WarehouseService/GetOrderStatus
```

**Ответ сразу:**
```json
{
  "orderId": "<UUID>",
  "status": "PENDING"
}
```

**Ответ через 10-12 секунд:**
```json
{
  "orderId": "<UUID>",
  "status": "RESERVED"
}
```



# Умный дом (Задача 5) `../devices`

Есть устройства, которые публикуют показания (Reading). Клиенты могут подписываться на поток показаний конкретного устройства в реальном времени.
Сущности: Device (упрощённо – просто идентификатор), Reading (device_id, value, timestamp).

## Что сделано

- `POST /readings` (REST) – публикация нового показания.
- `MonitorReadings(deviceID)` – gRPC Server Streaming, через который клиент получает поток новых показаний.
- Единый перехватчик аутентификации по API-ключу в обоих протоколах.

**Цель:** Использование gRPC streams и унификация middleware для REST и gRPC.

## Быстрый старт

1. Инструменты (один раз)

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

2. Склонируйте/откройте проект, скачайте зависимости

```bash
go mod tidy
```

3. Генерация кода

```bash
protoc -I api/proto -I third_party \
  --go_out=api/gen --go_opt=paths=source_relative \
  --go-grpc_out=api/gen --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=api/gen --grpc-gateway_opt=paths=source_relative \
  api/proto/devices/v1/service.proto
```

## Тестирование

### Запуск Сервера

```bash
go run cmd/server/main.go
```

1. Публикация показаний через REST (с API-ключом)

```bash
curl -X POST http://localhost:8080/api/v1/readings \
     -H "X-API-Key: secret-api-key" \
     -H "Content-Type: application/json" \
     -d '{"device_id":"sensor-1","value":23.5}'
```
Без ключа или с неверным ключом вернётся 401/403.

2. Подписка на поток через gRPC (с API-ключом в метаданных)

```bash
grpcurl -plaintext -H "x-api-key: secret-api-key" \
        -d '{"device_id":"sensor-1"}' \
        localhost:9090 devices.v1.DeviceService/MonitorReadings
```
Этот вызов зависнет в ожидании данных. 

В другом терминале опубликуйте ещё несколько показаний:
```bash
curl -X POST http://localhost:8080/api/v1/readings \
     -H "X-API-Key: secret-api-key" \
     -H "Content-Type: application/json" \
     -d '{"device_id":"sensor-1","value":24.0}'
```

В первом терминале вы увидите поток JSON-объектов с новыми показаниями. Для остановки нажмите `Ctrl+C`.

3. REST‑стрим через grpc‑gateway

```bash
curl -N -H "X-API-Key: secret-api-key" \
     http://localhost:8080/api/v1/devices/sensor-1/stream
```

Этот вызов зависнет в ожидании данных. 

В другом терминале опубликуйте ещё несколько показаний:
```bash
curl -X POST http://localhost:8080/api/v1/readings \
     -H "X-API-Key: secret-api-key" \
     -H "Content-Type: application/json" \
     -d '{"device_id":"sensor-1","value":25.0}'
```

В первом терминале вы увидите поток строк с JSON-объектами по мере публикации. Для остановки нажмите `Ctrl+C`.

4. Проверка авторизации

Попробуйте вызвать `MonitorReadings` без ключа:
```bash
grpcurl -plaintext -d '{"device_id":"sensor-1"}' localhost:9090 devices.v1.DeviceService/MonitorReadings
```

Получите ошибку `Unauthenticated`.