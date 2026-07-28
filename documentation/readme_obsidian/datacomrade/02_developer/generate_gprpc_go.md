2026-07-28 21:30

Status:

Tags: [[new_microservice]] [[workflow]] [[validation]] [[datacatalogue]]

# generate_gprpc_go

Третий, последний шаг цепочки: по `query.sql`, sqlc-моделям и `.proto` генерируются четыре слоя Go — `converter`, `validation`, `service`, `api`.

- Инструкция: `documentation/dev_instructions/crud/generate_gprpc_go.md`
- Скрипты: `deploy/generators`
- Выход: `datacatalogue/internal/{converter,validation,service,api}`

Цепочка целиком: [[sqlc_generation]] → [[standard_crud]] → [[proto_based_on_crud]] → **этот шаг**.

## Зачем

184 RPC на 22 таблицы — это около 135 файлов почти одинакового кода. Руками такой объём набирается с опечатками, которые компилятор не ловит: перепутанные местами `host_env` и `port_env` — оба `string`, сборка пройдёт. Генератор берёт имена полей из уже сгенерированных sqlc-структур и proto-сообщений, поэтому перепутать их не может.

## Конвейер

| Шаг | Скрипт | Вход | Выход |
| --- | --- | --- | --- |
| 1 | `parse.py` | `schema.sql`, `query.sql`, `models.go`, `query.sql.go`, `*.pb.go` | `build/parsed.json` |
| 2 | `resolve.py` | `parsed.json` | `build/model.json` |
| 3 | `gen.py` | `model.json` | Go-файлы в `internal/` |

Настройки — в `deploy/generators/crud_config.json`. Python 3, сторонних зависимостей нет.

```bash
python deploy/generators/parse.py && python deploy/generators/resolve.py && python deploy/generators/gen.py
```

Дальше обязательно форматирование — генератор не выравнивает структуры, это делает `gofmt`:

```bash
cd datacatalogue && gofmt -w internal/ && go build ./... && go vet ./... && go test ./internal/...
```

## Как устроены слои

Поток данных на примере `CreateHost`:

```
req → validation.ValidateCreateHost(req) → converter.ToCreateHostParams(req)
    → service.CreateHost → repo.CreateHost → DcHost
    → converter.HostToProto → CreateHostResponse
```

Промежуточного слоя DTO нет: конвертер отображает proto прямо в sqlc-типы и обратно. Сервис работает с `CreateHostParams` и `DcHost`.

**converter** — чистые функции без ошибок и походов в базу, по четыре на таблицу: `<E>ToProto`, `<E>sToProto`, `ToCreate<E>Params`, `ToUpdate<E>ByIdParams`.

**validation** — проверки полей поверх [[Валидатор полей]]. Границы длин собираются в `limits.go` прямо из `schema.sql`.

**service** — обёртка над репозиторием: перевод `sql.ErrNoRows` в `ErrNotFound` и проверка существования записи перед мягким удалением.

**api** — хендлеры gRPC, один метод на один RPC.

## Что перезаписывается

`gen.py` **затирает без спроса**, правки в этих файлах будут потеряны:

- `internal/converter/<домен>/<таблица>.go` и `_test.go`
- `internal/validation/<домен>/<таблица>.go`, `_test.go`, `limits.go`
- `internal/service/<домен>_service/<таблица>.go`
- `internal/api/<домен>apiv1/<таблица>.go`

**Написано руками, генератор не трогает:**

- `internal/converter/convert.go` — хелперы `TimeToProto`, `NullTimeToProto`
- `internal/validation/validation.go` — `ValidateID`
- `internal/api/apierror/apierror.go` — перевод ошибок в коды gRPC
- все `init.go`, `internal/utils/custom_errors/crud.go`, `internal/db/db.go`, `cmd/main.go`

Логику сверх шаблона нельзя дописывать в сгенерированный файл: либо отдельный файл того же пакета, либо правка шаблона в `gen.py`.

## Правила валидации

| Колонка | Проверка |
| --- | --- |
| `varchar(n) NOT NULL` | `StringVarchar` — непустая и не длиннее n **символов**, не байт |
| `bigint` | `Int64ID` — строго больше нуля |
| `smallint` | `Int32Between` в границах `int16` |
| `boolean` | не проверяется |

Про `smallint` стоит пояснить. В proto это `int32`, в базе — `int2`. Без проверки границ значение вроде 40000 прошло бы валидацию и упёрлось бы уже в Postgres. Единственная такая колонка сейчас — `level` в [[dc.group_levels]].

Существование записи по внешнему ключу **не проверяется**: битый `user_id` дойдёт до базы и вернётся ошибкой вставки. Это осознанное решение, а не недосмотр.

## Проверки после генерации

Зелёные тесты сами по себе ничего не доказывают. Встроенный `Unimplemented<Сервис>Server` молча подставляет заглушку вместо пропущенного метода, поэтому забытый RPC компиляцию не ломает. Нужны две сверки по списку — команды приведены в инструкции:

1. каждому RPC из `.proto` соответствует хендлер;
2. каждый метод репозитория обёрнут сервисом.

Обе запускаются из корня репозитория и считают нулевой результат провалом — иначе запуск не из того каталога напечатал бы `MATCH` при `0 = 0`.

Сейчас: 184 RPC (128 / 8 / 48), 22 таблицы, 851 тест.

## Генератор падает, а не додумывает

Конвейер останавливается с ошибкой при любом расхождении между источниками. Это осознанно: молча сгенерированный неверный маппинг ищется гораздо дольше, чем упавший скрипт. Проверяется, что набор запросов таблицы полон, каждому запросу есть метод репозитория, каждому полю sqlc — поле в proto, имя колонки из proto-тега есть в `schema.sql`, у каждого `Response` ровно одно поле и для каждой пары типов есть правило.

Так нашлось расхождение в [[dc.has_to_gpoup]]: колонка `description` объявлена `NOT NULL`, а поле в `.proto` было помечено `optional`.

## Когда перезапускать

| Что изменилось | Что сделать до генерации |
| --- | --- |
| `schema.sql` | пересобрать `query.sql` по [[standard_crud]], затем `sqlc generate` |
| `query.sql` | `sqlc generate`, обновить `.proto` по [[proto_based_on_crud]] |
| `.proto` | `task proto:gen` |
| только шаблон кода | ничего, достаточно `gen.py` |

## Как поменять папки

Все пути и имена пакетов — в `crud_config.json`, менять код не нужно. Пути задаются от корня репозитория; корень скрипты вычисляют сами от своего расположения, абсолютных путей нигде нет. Можно подсунуть свой конфиг:

```bash
python deploy/generators/gen.py --config path/to/my_config.json --work-dir path/to/build
```

При переезде слоя в другую папку старые файлы **не удаляются автоматически** — их нужно снести вручную, иначе в пакете останутся дубли функций и сборка упадёт.

## Ограничения

- Состав полей `Create` и `Update` берётся из sqlc-структур параметров, то есть задаётся `query.sql`, а не `.proto`. Поле, которого нет в `INSERT`, будет проигнорировано молча.
- Пагинации нет, списочные запросы отдают всё.
- Тестами покрыты только конвертеры и валидация — чистые функции без базы. Сервисы и хендлеры не покрыты.
- `Delete` и `Undelete` делают лишний `SELECT` перед `UPDATE`, потому что запросы объявлены как `:exec` и не сообщают число затронутых строк.

### References
- [[standard_crud]]
- [[proto_based_on_crud]]
- [[sqlc_generation]]
- [[Валидатор полей]]
- [[Taskdetail]]
- [[README]]
