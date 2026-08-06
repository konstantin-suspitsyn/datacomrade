# Генерация .proto — SG Buddy

`.proto` не пишется руками и не собирается скриптами репозитория. Его делает
та же программа, что и `query.sql`, — **SG Buddy**,
<https://github.com/konstantin-suspitsyn/sg_buddy>, из того же `schema.json`.

Второй шаг цепочки:
[standard_crud.md](standard_crud.md) → **этот файл** → [generate_gprpc_go.md](generate_gprpc_go.md).

Скрипты `deploy/generators/proto_gen_*.py`, которые делали это раньше, удалены.

## Как запускается

Отдельного действия нет: программа пишет `query.sql` и `.proto` за одно
сохранение. Куда класть контракт и что писать в его шапке, спрашивается при
настройке схемы и хранится в `schema.json`:

| Ключ `schema.json` | Что это |
|---|---|
| `save_proto_path` | куда положить `.proto` |
| `proto_package` | `package` контракта |
| `go_package` | `option go_package` |

Пустыми эти ключи оставлять нельзя — контракт с пустым `package` не соберётся.

После генерации:

```bash
task proto:lint && task proto:gen
```

## Что получается

Главное правило прежнее: **каждому запросу `query.sql` соответствует ровно один
RPC с тем же именем**. На нём держится генератор Go-слоёв — он сопоставляет
запрос, метод репозитория и RPC по имени и падает, если тройка не сходится.

Форма контракта:

- на таблицу — сообщение со всеми её колонками: это строка, какой её вернут
  `SELECT *` и `RETURNING *`;
- на таблицу — свой сервис `<Таблица>Service`;
- на запрос — пара `<Имя>Request` / `<Имя>Response`;
- поля Request — параметры готового SQL в порядке появления: и те, что пришли
  из настроек, и те, что написаны в значении руками;
- ответ зависит от аннотации: `:exec` — пустое сообщение, `:one` — одна строка,
  `:many` — `repeated`; у `DELETE` строки нет никогда;
- выборка с явным списком колонок получает своё сообщение строки `<Имя>Row` —
  в ответе видно ровно то, что выгружается, а не вся таблица;
- постраничность и сортировка описываются целиком — см.
  [pagination_proto.md](pagination_proto.md);
- nullable-колонка и необязательный фильтр (`sqlc.narg`) дают `optional`-поле.

Тип параметра программа ищет по порядку: колонка своей таблицы, приведение из
SQL (`@id::uuid`), колонка того же имени в других таблицах схемы. Не нашлось
нигде — поле станет `string`, а сама программа напишет об этом в списке
проблем. Такое место нужно чинить в настройках, а не в `.proto`.

## Маппинг типов

| Postgres | proto |
|---|---|
| `bigint` / `int8` | `int64` |
| `integer` / `int4` | `int32` |
| `smallint` / `int2` | `int32` |
| `character varying` / `text` | `string` |
| `boolean` | `bool` |
| `timestamp` / `timestamptz` | `google.protobuf.Timestamp` |
| `numeric` | `string` (без потери точности) |
| `uuid` | `string` (канонический вид 8-4-4-4-12) |

**Осторожно с `optional`.** Он обязан соответствовать nullability колонки в
`schema.sql`. Расхождение уже случалось: `dc.has_to_group.description` объявлен
`NOT NULL`, а поле в контракте было `optional` — контракт обещал клиенту, что
поле можно не передавать, хотя база такую строку не приняла бы. Нашёл это
генератор Go-слоёв, он сверяет типы и падает на несоответствии.

## Параметр через внешний id пользователя

Если колонка `<col>_id` заполняется подзапросом
`(SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)` (см.
[standard_crud.md](standard_crud.md)), то поле в `Request` называется не
`<col>_id`, а `<col>_external_id`, тип `string` (uuid). Колонка `user_id` даёт
поле `user_external_id`, `updated_by_id` — `updated_by_external_id`.

Это не косметика: соглашение зашито в `EXTERNAL_ID_SUFFIX`/`pair_fields`
в [deploy/generators/resolve.py](../../../deploy/generators/resolve.py) —
генератор Go-слоёв ищет sqlc-параметр `ExternalID uuid.UUID` и сопоставляет его
с полем proto по этому суффиксу, а не по имени. Назвать иначе — сломать шаг
[generate_gprpc_go.md](generate_gprpc_go.md).

## Проверки

Имена RPC совпадают с именами запросов один в один:

```bash
grep -o "^-- name: [A-Za-z]*" datacatalogue/db/sqlc/tables_model/query.sql | sed 's/-- name: //' | sort > /tmp/q.txt; grep -o "rpc [A-Za-z]*" shared/proto/tables/v1/tables.proto | sed 's/rpc //' | sort > /tmp/r.txt; diff /tmp/q.txt /tmp/r.txt && echo MATCH
```

Затем `task proto:lint` и `task proto:gen`.

## Осторожно при перегенерации

Перегенерация меняет контракт: имена RPC и форма ответов приезжают из настроек,
а не из того, что было в файле раньше. Хендлеры, конвертеры и клиентов придётся
приводить к новому контракту — для слоёв каталога это делается автоматически
через [generate_gprpc_go.md](generate_gprpc_go.md), для внешних потребителей —
руками.

Править сам `.proto` бессмысленно: следующее сохранение в программе перезапишет
файл. Всё, что нужно поменять, меняется в настройках схемы.
