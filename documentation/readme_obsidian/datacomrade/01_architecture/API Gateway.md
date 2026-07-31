2026-07-30 16:20

Status:

Tags: [[architecture]] [[security]]

# API Gateway

Единственный сервис, открытый наружу. Реализован в репозитории как микросервис `shepherd`. REST/JSON на границе (роутер `go-chi`), gRPC внутрь — к [[Metadata Service]] и будущим сервисам.

## Задачи

- REST-хендлеры вручную поверх сгенерированных gRPC-клиентов (не `grpc-gateway` и не gRPC-Web) — полный контроль над формой запроса/ответа и возможностью агрегировать несколько бэкендов в один ответ фронтенду.
- Валидация access-token Keycloak: **локальная проверка JWT по JWKS** (ключи кешируются в Redis, TTL по `Cache-Control` Keycloak — см. [[Кеширование Redis]]), без token introspection на каждый запрос.
- Токены — stateless Bearer: фронтенд логинится в Keycloak напрямую и сам хранит access token, Gateway токенов и сессий не хранит.
- Извлечение `sub` и ролей из claims, проброс `sub` (=`external_id`) и резолвленного `user_id` во внутренние gRPC-вызовы через metadata — см. «Автосоздание пользователя» ниже.
- RBAC по трём realm-ролям Keycloak — `admin` / `maintainer` / `viewer` (константы `RoleAdmin`/`RoleMaintainer`/`RoleViewer` в `shepherd/internal/middleware/rbac`; названия буква в букву совпадают с ролями в Keycloak). Middleware `RequireRole(...)` вешается на конкретный роут/группу роутов — это гейт **на уровне маршрута** («может ли вообще дёрнуть этот метод»). Не путать с [[dc.domain_roles]]/[[dc.table_roles]] — те решают, какие домены/таблицы видны внутри разрешённых методов, и живут в [[Metadata Service]]/[[Query Builder Service]] (см. [[Разграничение доступа]]).

  Сопоставление роль → группа методов Data Catalog:
  - `admin` — просмотр пользователей и назначение им прав: методы `userdomainrolesapiv1` (`datacatalogue/internal/api/userdomainrolesapiv1`).
  - `maintainer` — просмотр и использование каталога: методы `tablesapiv1` (`datacatalogue/internal/api/tablesapiv1` — host / database_cat / domain_cat / table_cat / column_cat и остальные справочники).
  - `viewer` — базовая роль без специальных прав сверх обычной аутентификации.
  - Создание пользователя (`UserService.CreateUser`, срабатывает через `POST /login`, см. «Автосоздание пользователя» ниже) нарочно не гейтится ролью — доступно любому аутентифицированному пользователю независимо от роли.

  На момент написания у Shepherd REST-роуты для `userdomainrolesapiv1`/`tablesapiv1` ещё не проброшены наружу (реализованы только `GET /v1/me` и `POST /v1/login`, ни один из них не гейтится ролью) — `RequireRole` объявлен и готов к использованию, но ни на один роут пока не навешан. Конкретное сопоставление роут → роль фиксируется по мере появления этих эндпоинтов в OpenAPI-контракте.
- CORS: фронтенд — отдельный домен и порт, allowlist origin'ов из конфига. Credentials не нужны, токен идёт в заголовке, не в cookie.
- Агрегация ответов для фронтенда (BFF): экран конструктора собирает дерево доступных таблиц/полей одним вызовом, дёргая несколько gRPC-клиентов параллельно и склеивая ответ в хендлере.
- Rate limiting — token bucket на Redis, ключ `rl:{user_id}:{window}` (см. [[Кеширование Redis]]).
- Старт трейса: генерация `trace_id`, если его нет во входящем запросе — см. [[Логирование]]. **Зависимость:** `platform/pkg/logger` пока не даёт способа положить `trace_id`/`user_id` в контекст снаружи пакета (см. [[Замечания к реализации]] п.4) — без экспортированных сеттеров обогащение логов Gateway работать не будет.

## Автосоздание пользователя

[[dc.user]] — локальное зеркало Keycloak-пользователя в [[Metadata Service]], находится по `external_id` (=`sub`). Раньше существование этой строки гарантировал middleware `EnsureUser`, вызывавшийся на **каждый** запрос под `/v1` — лишний поход в Redis/gRPC там, где `user_id` не нужен. Сейчас резолв явный и вынесен из middleware-цепочки в `shepherd/internal/ensureuser.Resolver`, у которого два метода с разной семантикой, вызываемые напрямую из конкретных хендлеров:

- **`POST /v1/login`** (`loginapiv1.Login`) → `Resolver.GetOrCreate` — единственное место, которое умеет писать в `dc.user`:
  1. Кеш `user:known:{external_id}` в Redis (см. [[Кеширование Redis]]) — есть запись → возвращаем `user_id`.
  2. Нет в кеше → `UserService.GetUserByExternalId`. Нашли → кешируем, возвращаем.
  3. `NotFound` → `UserService.CreateUser(name=<из claims>, external_id=sub)`. Успех → кешируем, возвращаем.
  4. Ошибка создания из-за гонки (два параллельных первых запроса одного нового пользователя одновременно бьют в уникальный `external_id`) → один повторный `GetUserByExternalId`; не нашли — 500.
- **`GET /v1/me`** (`meapiv1.GetMe`) → `Resolver.Resolve` — только читает (кеш → `GetUserByExternalId`); если записи нет, отдаёт **404** и просит вызвать `POST /login` — сам ничего не создаёт.

Фронтенд обязан вызвать `POST /v1/login` один раз сразу после получения токена от Keycloak (это и есть флоу логина в терминах Gateway); все остальные защищённые роуты монтируются только под `authMiddleware.Authenticate`, без резолва `user_id` за кулисами.

`CreateUser` в [[Metadata Service]] сейчас не идемпотентен (обычный `INSERT` без `ON CONFLICT`) — шаг 4 в `GetOrCreate` компенсирует это на стороне Gateway. Заведено отдельным пунктом в [[Замечания к реализации]].

## REST-контракт: OpenAPI / Swagger

Contract-first, по той же логике, что `shared/proto` + `buf` задают gRPC-границы:

- `shepherd/api/openapi/v1/openapi.yaml` — рукописный контракт, источник истины.
- `oapi-codegen` генерирует `internal/api/openapiv1/{types.go,server.gen.go}` — модели и `ServerInterface` под `go-chi`. Реализация хендлеров обязана покрывать весь интерфейс целиком — расхождение спеки и кода не компилируется.
- Рантайм-валидация входящих запросов по спеке — `oapi-codegen/nethttp-middleware`, отрабатывает до хендлера.
- `GET /openapi.yaml` отдаёт спеку как есть, `GET /docs` — Swagger UI; оба включены только вне production (`SHEPHERD_ENABLE_DOCS`).

## Почему проверка прав не только здесь

Gateway не знает, какие именно таблицы затрагивает конкретный запрос: он видит DTO конструктора, а не итоговый SQL. Полное разрешение прав требует резолва полей через [[Metadata Service]], поэтому вторая (основная) проверка выполняется в Builder перед сборкой SQL. Это defense in depth, а не дублирование.

### References
- [[Архитектура]]
- [[Разграничение доступа]]
- [[Логирование]]
- [[Кеширование Redis]]
- [[Технологический стек]]
- [[Замечания к реализации]]
