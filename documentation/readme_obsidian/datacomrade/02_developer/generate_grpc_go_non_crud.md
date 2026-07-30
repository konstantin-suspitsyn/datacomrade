2026-07-30 22:00

Status:

Tags: [[new_microservice]] [[workflow]] [[validation]] [[datacatalogue]]

# generate_grpc_go_non_crud

Боковая ветка от основной цепочки [[standard_crud]] → [[proto_based_on_crud]] → [[generate_gprpc_go]] — для доменов, которые не раскладываются на таблицы с `Create<E>`: сквозные выборки поверх нескольких таблиц, голые скаляры и списки id вместо строки целиком. Пример — `auth_logic`: `GetTableIdsByExternalUserIdAndRoles`, `GetTableIdsByUserIdAndRoles`.

- Инструкция: `documentation/dev_instructions/crud/generate_grpc_go_non_crud.md`
- Скрипты: `deploy/generators/{parse_custom,resolve_custom,gen_custom}.py`
- Выход: `datacatalogue/internal/{converter,validation,service,api}` — по одному файлу на домен, а не на таблицу

## Когда сюда, а не в [[generate_gprpc_go]]

`resolve.py` основного конвейера требует на каждую таблицу ровно один `Create<E>` плюс 8 стандартных операций и падает, если такого нет. Если запросы `query.sql` домена не привязаны к одной таблице — это сигнал использовать этот конвейер, а не пытаться подогнать `query.sql` под стандартную форму.

Оба конвейера читают один `crud_config.json` и пишут в один `internal/`, не пересекаются по путям и запускаются независимо.

## Конвейер

| Шаг | Скрипт | Вход | Выход |
| --- | --- | --- | --- |
| 1 | `parse_custom.py` | `schema.sql` (все таблицы домена сразу), `query.sql` (плоский список запросов, без блоков по таблицам), `query.sql.go`, `*.pb.go` | `build/parsed_custom.json` |
| 2 | `resolve_custom.py` | `parsed_custom.json` | `build/model_custom.json` |
| 3 | `gen_custom.py` | `model_custom.json` | Go-файлы в `internal/` |

Настройки — в `crud_config.json`, ключ `custom_domains` (та же форма метаданных, что и `domains` основного конвейера: `sqlc`, `proto_pkg`, `conv_dir`, `service_dir`, `api_dir` и т.д.).

```bash
python deploy/generators/parse_custom.py && python deploy/generators/resolve_custom.py && python deploy/generators/gen_custom.py
```

Дальше как обычно:

```bash
cd datacatalogue && gofmt -w internal/ && go build ./... && go vet ./... && go test ./internal/...
```

`.proto` домена должен быть уже сгенерирован в Go (`task proto:gen`) — конвейер читает `*.pb.go`, не сам `.proto`.

## Что проверяет resolve_custom.py

Падает при любом расхождении, как и [[generate_gprpc_go]]:

- у каждого запроса есть метод репозитория;
- аргумент репозитория — ровно один, и это sqlc-структура `...Params` (форма «голый скаляр без обёртки» пока не поддержана — падает явно, а не додумывает);
- каждому полю `...Params` находится поле в `<Запрос>Request` по нормализованному имени;
- имя колонки из proto-тега есть хоть в одной таблице `schema.sql` домена;
- у `<Запрос>Response` ровно одно поле;
- тип sqlc (`[]T` или `T`) соответствует типу ответа по `repeated` и по правилу из таблицы типов.

### column_hints — когда колонка неоднозначна

Запросы здесь — джойны нескольких таблиц, колонка не привязана к одной. Если одноимённая колонка встречается в разных таблицах домена с разной границей `varchar`, конвейер не гадает — падает с просьбой указать таблицу явно:

```json
"column_hints": {
  "GetTableIdsByExternalUserIdAndRoles": { "name": "dc.domain_roles" },
  "GetTableIdsByUserIdAndRoles": { "name": "dc.domain_roles" }
}
```

Нужно только когда границы расходятся между таблицами; если везде одна длина — `column_hints` не нужен.

## Что генерируется

На каждый запрос, в одном файле на домен (не на таблицу — сущности нет, есть только независимые запросы):

- **converter** — `To<Запрос>Params(req)` собирает параметры из gRPC-запроса, `<Запрос>ToProto(rows)` оборачивает результат репозитория в ответ.
- **validation** — `Validate<Запрос>(req) error`, по проверке на поле (`StringVarchar`/`Int64ID`/`StringUUID`), плюс `limits.go` с границами `varchar`.
- **service** — прямой вызов репозитория с оборачиванием ошибки, без обработки `sql.ErrNoRows` (как и списочные запросы основного конвейера).
- **api** — валидация → конвертация запроса → вызов сервиса → конвертация ответа.

Комментарии к функциям не описывают бизнес-смысл запроса — генератор его не знает, только то, что функция оборачивает `<sqlc>.<Запрос>`.

## Что перезаписывается, а что нет

Как и в [[generate_gprpc_go]]: `gen_custom.py` **затирает без спроса** `internal/<conv_dir>/<sqlc>.go` (+`_test.go`), `internal/<valid_dir>/<sqlc>.go` (+`_test.go`, `limits.go`), `internal/<service_dir>/<sqlc>.go`, `internal/<api_dir>/<sqlc>.go`.

**Написано руками для нового домена** (генератор не создаёт при первом запуске): `internal/service/<домен>_service/init.go`, `internal/api/<домен>apiv1/init.go`, плюс регистрация в `internal/service/services/init.go` и `cmd/main.go`. Без этого — `undefined: <ServiceType>` / `undefined: <ApiType>` при сборке.

## Ограничения

- Поддержана только форма «один параметр — sqlc-структура `...Params`». Запрос с голым скаляром без обёртки или без аргументов упадёт явно, а не сгенерирует что-то наугад.
- Один файл на домен — подходит, пока запросов немного. Если станет много, дробить по смыслу — ручная правка `gen_custom.py`, не то, что конвейер решает сам.
- Автоматической эвристики выбора таблицы для неоднозначной колонки нет и не должно быть — только явный `column_hints`.

## Найденный попутно баг

`parse_schema()` в `parse.py` (основной конвейер) ищет границу `varchar` только по спеллингу `character varying(n)`. `auth_logic/schema.sql` пишет эквивалентно, но иначе — `varchar(n)` — и граница молча пропадала. Починено в `parse_custom.py` (регулярка принимает оба варианта); в `parse.py` — ещё нет, задача на потом.

## Проверки после генерации

Те же две сверки, что и в [[generate_gprpc_go]], с именами домена `auth_logic`:

```bash
grep -ho "^  rpc [A-Za-z]*" shared/proto/auth_logic/v1/*.proto | sed 's/.*rpc //' | sort > /tmp/rpc.txt
grep -ho "^func (a \*AuthLogicApiV1) [A-Za-z]*" datacatalogue/internal/api/authlogicapiv1/*.go | sed 's/.*) //' | sort > /tmp/impl.txt
diff /tmp/rpc.txt /tmp/impl.txt && echo MATCH

grep -ho "^func (q \*Queries) [A-Za-z]*" datacatalogue/internal/repository/auth_logic/query.sql.go | sed 's/.*) //' | sort > /tmp/r.txt
grep -ho "^func (s \*AuthLogicService) [A-Za-z]*" datacatalogue/internal/service/auth_logic_service/*.go | sed 's/.*) //' | sort > /tmp/s.txt
diff /tmp/r.txt /tmp/s.txt && echo MATCH
```

Сейчас: домен `auth_logic`, 2 запроса, 2 RPC, сборка и тесты проходят.

### References
- [[standard_crud]]
- [[proto_based_on_crud]]
- [[generate_gprpc_go]]
- [[sqlc_generation]]
- [[Валидатор полей]]
- [[Taskdetail]]
