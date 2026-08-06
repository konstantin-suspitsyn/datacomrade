# Генерация Go-слоёв gRPC для доменов без табличного CRUD

Инструкция для Claude и для человека: как получить слои `converter`, `validation`,
`service` и `api` для доменов вроде `auth_logic` — sqlc-пакетов, где запросы не
раскладываются на одну таблицу с восемью стандартными операциями (см.
[standard_crud.md](standard_crud.md)), а являются сквозными выборками поверх
нескольких таблиц: `GetTableIdsByExternalUserIdAndRoles`, агрегаты, джойны и т.п.

Это отдельная ветка от основного конвейера
[standard_crud.md](standard_crud.md) → [proto_based_on_crud.md](proto_based_on_crud.md) →
[generate_gprpc_go.md](generate_gprpc_go.md). `query.sql` и `.proto` для таких
доменов тоже пишет SG Buddy (<https://github.com/konstantin-suspitsyn/sg_buddy>) —
запросы там просто настраиваются вручную, а не восьмёркой на таблицу.
Используй ту тройку файлов, если
домен раскладывается на таблицы с `Create<E>`; используй этот файл, если нет.
Оба конвейера читают один и тот же `crud_config.json`, пишут в один и тот же
`internal/`, не пересекаются по путям и могут запускаться независимо.

## Когда использовать этот конвейер, а не generate_gprpc_go.md

`resolve.py` (основной конвейер) требует на каждую таблицу ровно один `Create<E>`
плюс 8 стандартных операций. Если в `query.sql` домена такого нет — запросы не
привязаны к одной таблице, возвращают голые скаляры/списки id, а не строку
целиком — `resolve.py` упадёт с `Fail("ожидался ровно один Create...")`. Это
признак того, что нужен этот файл, а не правка `query.sql` под стандартную форму.

## Конвейер

Три скрипта в [deploy/generators](../../../deploy/generators), по образу
основного конвейера, но единица сборки — один запрос `query.sql`, а не одна
таблица:

| Шаг | Скрипт | Вход | Выход |
|---|---|---|---|
| 1 | [parse_custom.py](../../../deploy/generators/parse_custom.py) | `schema.sql` (все таблицы домена), `query.sql` (плоский список), `query.sql.go`, `*.pb.go` | `build/parsed_custom.json` |
| 2 | [resolve_custom.py](../../../deploy/generators/resolve_custom.py) | `build/parsed_custom.json` | `build/model_custom.json` |
| 3 | [gen_custom.py](../../../deploy/generators/gen_custom.py) | `build/model_custom.json` | Go-файлы в `internal/` |

Настройки — в `crud_config.json`, ключ `custom_domains` (список, той же формы
метаданных домена, что и `domains` основного конвейера: `sqlc`, `proto_pkg`,
`conv_dir`, `service_dir`, `api_dir` и т.д. — см.
[generate_gprpc_go.md](generate_gprpc_go.md#как-поменять-папки)).

На домен генерируется по одному файлу на слой (не по файлу на таблицу — здесь
нет сущности, только независимые запросы): `internal/<conv_dir>/<sqlc>.go`,
`internal/<service_dir>/<sqlc>.go`, `internal/<api_dir>/<sqlc>.go`, плюс тесты
для converter и validation и `limits.go`.

### Запуск

Из корня репозитория, python 3 без сторонних зависимостей:

```bash
python deploy/generators/parse_custom.py && python deploy/generators/resolve_custom.py && python deploy/generators/gen_custom.py
```

Затем как обычно:

```bash
cd datacatalogue && gofmt -w internal/ && go build ./... && go vet ./... && go test ./internal/...
```

Требуется, чтобы `.proto` домена уже был сгенерирован в Go (`task proto:gen`) —
конвейер читает `*.pb.go`, а не сам `.proto`.

## Что проверяет resolve_custom.py

Как и `resolve.py`, падает при любом расхождении, а не додумывает:

- у каждого запроса `query.sql` есть метод репозитория (`query.sql.go`);
- аргумент репозитория — ровно один, и это sqlc-структура `...Params`
  (форма "один скаляр без обёртки" пока не поддержана — падает с понятным
  сообщением, а не молча генерирует не то);
- каждому полю `...Params` находится поле в proto-сообщении `<Запрос>Request`
  (сопоставление по нормализованному имени, как в `resolve.py`);
- имя колонки из proto-тега есть хоть в одной таблице `schema.sql` домена;
- у `<Запрос>Response` ровно одно поле;
- возвращаемый sqlc-тип (`[]T` или `T`) соответствует типу поля ответа по
  `repeated`/не-`repeated` и по правилу из таблицы типов ниже.

### column_hints — когда имя колонки неоднозначно

Запросы здесь — джойны нескольких таблиц, поэтому колонка не привязана к одной
таблице. Если колонка с таким именем встречается в разных таблицах домена
**с разной границей `varchar`**, конвейер не будет гадать — упадёт с просьбой
указать нужную таблицу явно в конфиге:

```json
"column_hints": {
  "GetTableIdsByExternalUserIdAndRoles": { "name": "dc.domain_roles" },
  "GetTableIdsByUserIdAndRoles": { "name": "dc.domain_roles" }
}
```

Ключ верхнего уровня — имя запроса, внутри — колонка → таблица, откуда брать
границу `varchar` для валидации. Нужно только когда границы расходятся; если
одноимённая колонка везде одной длины, `column_hints` не нужен.

## Соответствие типов

То же самое, что в [generate_gprpc_go.md](generate_gprpc_go.md#правила-зашитые-в-генератор),
плюс то же самое применяется поэлементно к `repeated`-полям ответа (список
`int64` → `repeated int64`, конвертер по элементу не нужен, если типы совпадают
один в один; если нет — генерируется цикл с поэлементным преобразованием).
Список пар размечен в `resolve_custom.py` (`TYPE_RULES`) — там же и расширять
при добавлении нового типа.

## Что генерируется

На каждый запрос:

- **converter**: `To<Запрос>Params(req) <repo>.<Запрос>Params` — собирает
  параметры из gRPC-запроса; `<Запрос>ToProto(rows) *<proto>.<Запрос>Response` —
  оборачивает результат репозитория в ответ.
- **validation**: `Validate<Запрос>(req) error` — по одной проверке на поле
  (`StringVarchar`/`Int64ID`/`StringUUID`, `bool` не проверяется), плюс
  `limits.go` с границами `varchar` по колонкам.
- **service**: `(s *<ServiceType>) <Запрос>(ctx, params) (<T>, error)` — прямой
  вызов репозитория с оборачиванием ошибки. Без обработки `sql.ErrNoRows`:
  как и списочные запросы основного конвейера, эти запросы не считают
  «ничего не найдено» ошибкой.
- **api**: `(a *<ApiType>) <Запрос>(ctx, req) (*<Resp>, error)` —
  валидация → конвертация запроса → вызов сервиса → конвертация ответа.

Комментарии к функциям не описывают бизнес-смысл запроса (генератор не может
его знать) — только то, что функция оборачивает `<sqlc>.<Запрос>`. Если нужен
содержательный комментарий, дописать его руками поверх — `gen_custom.py`
перезаписывает файл целиком при повторном запуске, поэтому такая правка не
переживёт следующую генерацию не глядя (как и в основном конвейере).

## Что перезаписывается, а что нет

Такие же правила, как в [generate_gprpc_go.md](generate_gprpc_go.md#что-перезаписывается-а-что-нет):

**gen_custom.py перезаписывает без спроса**: `internal/<conv_dir>/<sqlc>.go` и
`_test.go`, `internal/<valid_dir>/<sqlc>.go`, `_test.go` и `limits.go`,
`internal/<service_dir>/<sqlc>.go`, `internal/<api_dir>/<sqlc>.go`.

**Написано руками для нового домена, генератор не трогает**:
`internal/service/<домен>_service/init.go`, `internal/api/<домен>apiv1/init.go`,
плюс регистрация в `internal/service/services/init.go` и `cmd/main.go` (по
образцу уже подключённых доменов). Это ровно то, чего не хватает после первого
запуска конвейера для нового домена — без этого будет `undefined:
<ServiceType>` / `undefined: <ApiType>` при сборке.

## Известные ограничения

- Поддержана только форма аргументов «один параметр — sqlc-структура
  `...Params`» (это всё, что встречается в `custom_domains` сейчас). Запрос
  с одним голым скалярным аргументом без обёртки или с нулём аргументов упадёт
  явно, а не сгенерирует что-то наугад — для такой формы правила ещё не
  написаны, дописывать в `resolve_custom.py`/`gen_custom.py` по мере надобности.
- Один файл на домен, не на запрос и не на таблицу: подходит, пока запросов в
  домене немного (единицы-десятки). Если станет много, дробить по смыслу —
  это ручная правка `gen_custom.py` (например, по признаку из имени запроса),
  а не то, что конвейер решает сам.
- `varchar`-граница ищется по всем таблицам домена; если она неоднозначна —
  обязателен `column_hints` (см. выше), автоматической эвристики выбора
  таблицы нет и не должно быть.

## Проверки после генерации

Те же две сверки, что в
[generate_gprpc_go.md](generate_gprpc_go.md#проверки-после-генерации), с теми же
именами домена/пакета/типа:

```bash
grep -ho "^  rpc [A-Za-z]*" shared/proto/auth_logic/v1/*.proto | sed 's/.*rpc //' | sort > /tmp/rpc.txt
grep -ho "^func (a \*AuthLogicApiV1) [A-Za-z]*" datacatalogue/internal/api/authlogicapiv1/*.go | sed 's/.*) //' | sort > /tmp/impl.txt
diff /tmp/rpc.txt /tmp/impl.txt && echo MATCH

grep -ho "^func (q \*Queries) [A-Za-z]*" datacatalogue/internal/repository/auth_logic/query.sql.go | sed 's/.*) //' | sort > /tmp/r.txt
grep -ho "^func (s \*AuthLogicService) [A-Za-z]*" datacatalogue/internal/service/auth_logic_service/*.go | sed 's/.*) //' | sort > /tmp/s.txt
diff /tmp/r.txt /tmp/s.txt && echo MATCH
```

На момент написания: домен `auth_logic`, 2 запроса, 2 RPC, сборка и тесты
проходят (`go build ./...`, `go vet ./...`, `go test ./internal/...`).

## Известный найденный баг источников (не в этом конвейере)

`parse_schema()` в `parse.py` (основной конвейер) и в `parse_custom.py` ищет
границу `varchar` только по спеллингу `character varying(n)`. Схема
`auth_logic/schema.sql` пишет `varchar(n)` — эквивалентно для Postgres, но
другим текстом. В `parse_custom.py` это уже исправлено (регулярка принимает
оба варианта); `parse.py` пока нет — если когда-нибудь `schema.sql` основного
конвейера напишут через `varchar(n)`, граница длины молча пропадёт. Чинить
там же, тем же способом.
