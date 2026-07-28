# Генерация .proto по query.sql

Инструкция для Claude: как построить `.proto` файл, полностью соответствующий
CRUD-запросам из `query.sql`. Дополняет [standard_crud.md](standard_crud.md) —
сначала генерируется `query.sql`, потом по нему `.proto`.

## Входные данные

- `query.sql` — источник истины по списку операций (имена RPC берутся отсюда).
- `schema.sql` — источник истины по типам и составу полей.

Больше никуда не смотреть. Если в задании названы конкретные папки — работать
только в них.

## Главное правило

**Каждому запросу в `query.sql` соответствует ровно один RPC с тем же именем.**
Количество RPC = количество запросов. Никаких лишних и никаких пропущенных.

## Структура файла

```
syntax = "proto3";
package <domain>.v1;
import "google/protobuf/empty.proto";
import "google/protobuf/timestamp.proto";
option go_package = "...";

service <Domain>Service { ... }   // все RPC, сгруппированы по таблицам
// далее блоки сообщений по таблицам, порядок как в schema.sql
```

Перед блоком каждой таблицы — разделитель:

```proto
// =========================================================
// dc.<table_name>
// =========================================================
```

Внутри блока RPC идут в порядке: чтение (`Get*`), выборки по FK, затем запись
(`Create`, `Update`, `Delete`, `Undelete`).

## Сообщения на таблицу

### 1. Сущность — одна message на всю строку таблицы

```proto
// Host is a full row of dc.host.
message Host {
  int64 id = 1;
  string name = 2;
  ...
}
```

Все `Response`, возвращающие данные, ссылаются на неё, а не дублируют поля.

### 2. Request / Response на каждый RPC

| RPC | Request | Response |
|---|---|---|
| `Get<E>ById` | `int64 id = 1;` | `<E> <snake_e> = 1;` |
| `Get<E>s` | пустое тело `{}` | `repeated <E> <snake_plural> = 1;` |
| `GetDeleted<E>ById` | `int64 id = 1;` | `<E> <snake_e> = 1;` |
| `GetDeleted<E>s` | пустое тело `{}` | `repeated <E> <snake_plural> = 1;` |
| `Get<E>sBy<Fk>` | `int64 <fk_field> = 1;` | `repeated <E> <snake_plural> = 1;` |
| `Create<E>` | параметры `INSERT` | `<E> <snake_e> = 1;` |
| `Update<E>ById` | `id` + параметры `UPDATE` | `<E> <snake_e> = 1;` |
| `Delete<E>ById` | `int64 id = 1;` | `google.protobuf.Empty empty = 1;` |
| `Undelete<E>ById` | `int64 id = 1;` | `google.protobuf.Empty empty = 1;` |

- Имя message = имя RPC + `Request` / `Response` (требование `buf lint`).
- В `Create`/`Update` request **не входят** `id` (для Create), `created_at`,
  `updated_at`, `is_deleted` — они выставляются в SQL.
- Списочные Request делать пустыми `message ...Request {}`, а не через
  `google.protobuf.Empty` — чтобы потом можно было добавить пагинацию.

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

Nullable-колонка → `optional <type>` (для message-типов вроде `Timestamp`
presence есть по умолчанию, `optional` не нужен).

## Комментарии — обязательны

`buf lint` (правило `COMMENT_MESSAGE`) требует непустой комментарий **у каждой**
message, включая Response. Комментарий пишется строкой выше `message`, начинается
с имени сообщения и содержит ссылку на исходный запрос:

```proto
// GetHostByIdResponse returns a single active dc.host row (query GetHostById).
message GetHostByIdResponse {
  Host host = 1;
}
```

Шаблоны формулировок:

- сущность — `<E> is a full row of <table>.`
- `Get...ByIdRequest` — `asks for a single active <table> row by id`
- `Get...ByIdResponse` — `returns a single active <table> row`
- `Get...sRequest` / `Response` — `asks for / returns every active <table> row`
- `GetDeleted...` — то же, но `soft deleted` вместо `active`
- `Get...sBy<Fk>...` — `asks for / returns active <table> rows filtered by <fk_field>`
- `Create...Request` — `carries the fields needed to insert a <table> row`
- `Create...Response` — `returns the created <table> row`
- `Update...ByIdRequest` — `carries the fields to update on a <table> row`
- `Update...ByIdResponse` — `returns the updated <table> row`
- `Delete...ByIdRequest` — `asks to soft delete a <table> row by id, setting is_deleted = true`
- `Undelete...ByIdRequest` — `asks to restore a soft deleted <table> row by id, setting is_deleted = false`
- `Delete/Undelete...ByIdResponse` — `carries no payload`

Комментарии к каждому `rpc` не добавлять — блок `service` становится
нечитаемым. Добавить только если включено правило `COMMENT_RPC`.

## Как генерировать

При 15+ таблицах это ~270 сообщений — руками не набирать, будут опечатки.
Написать одноразовый Python-скрипт в скратчпаде: список таблиц с полями
(строка, поля Create, список FK) → генерация текста → запись `.proto`.
Так правки схемы применяются перегенерацией, а нумерация полей и имена
сообщений гарантированно согласованы.

## Проверка после генерации

1. Имена RPC совпадают с именами запросов один в один:

```bash
grep -o "^-- name: [A-Za-z]*" query.sql | sed 's/-- name: //' | sort > /tmp/q.txt; grep -o "^  rpc [A-Za-z]*" tables.proto | sed 's/  rpc //' | sort > /tmp/r.txt; diff /tmp/q.txt /tmp/r.txt && echo MATCH
```

2. У каждой message есть комментарий:

```bash
python -c "import io;l=io.open('tables.proto',encoding='utf-8').read().split('\n');print(sum(1 for i,x in enumerate(l) if x.startswith('message ') and not l[i-1].strip().startswith('//')))"
```

Должен вывести `0`.

3. Прогнать `buf lint` и генерацию кода — в окружении может не быть `protoc`,
   тогда синтаксис проверяется только на стороне пользователя.

## Внимание

Перегенерация ломает API: имена RPC меняются (`GetHost` → `GetHostById`),
ответы становятся вложенными (`Host host = 1;` вместо плоских полей).
Хендлеры, конвертеры и клиентов придётся переписывать под новый контракт.
