2026-07-31 19:20

Status:

Tags: [[workflow]] [[must_install]] [[new_microservice]]

# Taskfile

Единая точка автоматизации проекта — корневой `Taskfile.yml`. Через [Task](https://taskfile.dev/) запускаются форматирование, линтинг, генерация кода (protobuf, OpenAPI, sqlc), генерация `.env` и поднятие/остановка docker-compose окружений. Установка самого Task — см. [[Taskdetail]].

## Как устроены задачи

Имена задач — `область:действие` (`proto:gen`, `sqlc:gen`, `env:gen`...); задачи без области (`format`, `lint`, `up-core`...) — общепроектные или per-компонентные, исторически без namespace.

В большинстве блоков есть **общая задача**, которую реально запускают руками (`*:gen`, `format`, `lint`, `up-all`...), и **под-таски**, которые она тянет за собой сама — руками их обычно не вызывают. Под-таски подключены двумя разными механизмами:

- **`deps: [...]`** — под-таск (обычно `*:install-plugins` / `install-*`) выполняется автоматически перед `cmds` самой задачи. У install-тасков есть блок `status:`, который через `test -x <бинарник>` пропускает переустановку, если инструмент уже стоит в `bin/`.
- **`cmds: - task: Y`** (или просто `- task Y` в shell-команде) — общая задача явно вызывает другие задачи из своего тела, одну за другой. Так работают `sqlc:gen` → `sqlc:gen:datacatalogue` и `up-all`/`down-all` → задачи `up-*`/`down-*`.

## Глобальные переменные (`vars`)

### Версии инструментов

| Переменная | Значение | Использует |
| --- | --- | --- |
| `GO_VERSION` | `1.25.0` | не используется ни в одной задаче — фиксирует версию Go проекта как документацию |
| `GOLANGCI_LINT_VERSION` | `v2.11.3` | `install-golangci-lint` |
| `GCI_VERSION` | `v0.14.0` | `install-formatters` |
| `GOFUMPT_VERSION` | `v0.9.2` | `install-formatters` |
| `BUF_VERSION` | `1.53.0` | `install-buf` (версия подставляется как `@v{{.BUF_VERSION}}`) |
| `PROTOC_GEN_GO_VERSION` | `v1.36.6` | `proto:install-plugins` |
| `PROTOC_GEN_GO_GRPC_VERSION` | `v1.5.1` | `proto:install-plugins` |
| `OAPI_CODEGEN_VERSION` | `v2.8.0` | `openapi:install-plugins` |
| `SQLC_VERSION` | `v1.31.1` | `sqlc:install-plugins` |

### Пути и списки

| Переменная | Значение | Использует |
| --- | --- | --- |
| `BIN_DIR` | `{{.ROOT_DIR}}/bin` | все install-таски — общая папка для скачанных бинарников |
| `GOLANGCI_LINT`, `GCI`, `GOFUMPT`, `BUF`, `PROTOC_GEN_GO`, `PROTOC_GEN_GO_GRPC`, `OAPI_CODEGEN`, `SQLC` | `{{.BIN_DIR}}/<бинарник>` | соответствующие install- и gen-таски |
| `GRPCURL` | `{{.BIN_DIR}}/grpcurl` | объявлена, но пока ни одна задача её не использует (зарезервирована) |
| `MODULES` | `datacatalogue shepherd shared` | `format`, `lint` — список Go-модулей, по которым идёт цикл |
| `SERVICES` | `datacatalogue shepherd shared` | объявлена, но пока ни одна задача её не использует (зарезервирована) |

## Форматирование и линтинг

### `format` — общая задача

Форматирует все `*.go` файлы модулей из `{{.MODULES}}` (кроме `*/mocks/*`): сперва `gofumpt -extra -w`, затем `gci write` с сортировкой импортов (`standard`, `default`, `prefix(https://github.com/konstantin-suspitsyn/datacomrade/)`).

Под-таск:

| Задача | Что делает |
| --- | --- |
| `install-formatters` | ставит `gofumpt` и `gci` в `bin/`, если их там ещё нет (`deps` у `format`) |

### `lint` — общая задача

Прогоняет `golangci-lint run <модуль>/... --config=.golangci.yml` по каждому модулю из `{{.MODULES}}`; ошибка любого модуля не прерывает остальные — итоговый код возврата ненулевой, если хоть один упал.

Под-таски (оба в `deps`, выполняются перед линтом; `install-golangci-lint` и `format` независимы друг от друга и запускаются параллельно):

| Задача | Что делает |
| --- | --- |
| `install-golangci-lint` | ставит `golangci-lint` в `bin/` |
| `format` | форматирует код (см. блок выше) перед проверкой линтером |

## Protobuf (`shared/proto`)

### `proto:gen` — общая задача

Генерирует Go-код из `.proto` через `buf generate` (`dir: shared/proto`).

Под-таски (все в `deps`, выполняются перед генерацией):

| Задача | Что делает |
| --- | --- |
| `install-buf` | ставит `buf` в `bin/` |
| `proto:install-plugins` | ставит `protoc-gen-go` и `protoc-gen-go-grpc` в `bin/` |
| `proto:lint` | `buf lint` — проверка `.proto` на стиль; сама тоже зависит от `install-buf` + `proto:install-plugins` и может запускаться отдельно |

## OpenAPI (Shepherd)

### `openapi:gen` — общая задача

Генерирует `types.go`/`server.gen.go` (модели + chi `ServerInterface`) из `shepherd/api/openapi/v1/openapi.yaml` — `dir: shepherd`.

Под-таск:

| Задача | Что делает |
| --- | --- |
| `openapi:install-plugins` | ставит `oapi-codegen` в `bin/` (`deps` у `openapi:gen`) |

## sqlc

### `sqlc:gen` — общая задача

Не генерирует код сама — просто вызывает под-таски вида `sqlc:gen:<микросервис>` (сейчас один — `sqlc:gen:datacatalogue`). При добавлении sqlc для нового микросервиса нужно завести его собственный `sqlc:gen:<микросервис>` и добавить вызов сюда — см. [[sqlc_generation]].

Под-таски:

| Задача | Что делает |
| --- | --- |
| `sqlc:install-plugins` | ставит `sqlc` в `bin/` |
| `sqlc:gen:datacatalogue` | `sqlc generate` в `datacatalogue/db/sqlc` (по `sqlc.yml`); сама зависит от `sqlc:install-plugins` |

## Генерация `.env`

### `env:gen` — общая задача

Запускает `deploy/python-scripts/generate_env.py deploy/env` — раскладывает `.env` по всем сервисам и compose-папкам из `deploy/env/params.toml` (+ `params.variables.toml` для `${VARIABLE}`-подстановок). Пути в команде относительные — от корня репозитория, где лежит `Taskfile.yml`.

Под-таск:

| Задача | Что делает |
| --- | --- |
| `env:install-requirements` | `python -m pip install -r requirements.txt` в `deploy/python-scripts` (`deps` у `env:gen`) |

## Docker Compose

Компонентные задачи `up-*`/`down-*` — каждая поднимает/останавливает один compose-стек и может запускаться отдельно:

| Задача | `dir` |
| --- | --- |
| `up-core` / `down-core` | `deploy/compose/core` |
| `up-datacatalogue` / `down-datacatalogue` | `deploy/compose/datacatalogue` |
| `up-keycloak` / `down-keycloak` | `deploy/compose/keycloak` |
| `up-redis` / `down-redis` | `deploy/compose/redis` |

### `up-all` / `down-all` — общие задачи

Просто вызывают компонентные задачи по очереди (`cmds: - task up-core` и т.д., не `deps`, то есть строго последовательно, не параллельно):

- `up-all`: `up-core` → `up-datacatalogue` → `up-keycloak` → `up-redis`.
- `down-all`: `down-redis` → `down-keycloak` → `down-datacatalogue` → `down-core` (обратный порядок относительно `up-all`).

### References
- [[Taskdetail]]
- [[sqlc_generation]]
- [[Работа с env файлами]]
- [[generate_gprpc_go]]
- [[Перезапуск микросервисов - рабочий процесс]]
