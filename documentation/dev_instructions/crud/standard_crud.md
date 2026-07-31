# Стандартный CRUD для sqlc

Инструкция для Claude: как генерировать `query.sql` по `schema.sql`.

## Задача

Из `schema.sql` (в той же папке) сгенерировать `query.sql` — по 8 запросов на каждую таблицу.
Смотреть только на `schema.sql`, ничего не додумывать из другого кода.

## Набор запросов на таблицу

`<Entity>` — имя таблицы в PascalCase единственного числа (`table_cat` → `TableCat`,
`group_levels` → `GroupLevel`). `<Entity>s` — множественное число по правилам английского
(`Alias` → `Aliases`, `SchemaCat` → `SchemaCats`).

| # | Имя запроса | Тип | Что делает |
|---|---|---|---|
| 1 | `Get<Entity>ById` | `:one` | активная запись по id (`is_deleted = false`) |
| 2 | `Get<Entity>s` | `:many` | все активные, `ORDER BY id` |
| 3 | `GetDeleted<Entity>ById` | `:one` | удалённая запись по id (`is_deleted = true`) |
| 4 | `GetDeleted<Entity>s` | `:many` | все удалённые, `ORDER BY id` |
| 5 | `Create<Entity>` | `:one` | вставка, `RETURNING *` |
| 6 | `Update<Entity>ById` | `:one` | обновление по id, `RETURNING *` |
| 7 | `Delete<Entity>ById` | `:exec` | мягкое удаление (`is_deleted = true`) |
| 8 | `Undelete<Entity>ById` | `:exec` | восстановление (`is_deleted = false`) |

## Дополнительные выборки

Сверх восьми стандартных запросов таблица может иметь выборки по другим
колонкам. Имя определяет, что запрос отдаёт, — генератор Go-слоёв опирается
именно на него:

| Имя запроса | Тип | Когда |
|---|---|---|
| `Get<Entity>sBy<Fk>` | `:many` | колонка неуникальна (внешний ключ) |
| `Get<Entity>By<Column>` | `:one` | на колонке UNIQUE |

```sql
-- name: GetUserByExternalId :one
SELECT *
FROM dc."user"
WHERE external_id = $1
AND is_deleted = false;
```

Фильтр `is_deleted = false` обязателен, как и в остальных выборках активных
строк. Такие запросы придумываются не сами по себе: их добавляют по задаче.

## Правила

- **id** в `INSERT` не передаётся — полагаемся на дефолт последовательности (`bigserial`).
- **created_at / updated_at** проставляются в SQL через `now()`, не параметром.
  При `UPDATE`, `Delete`, `Undelete` всегда обновляется `updated_at`.
- **is_deleted** при `INSERT` — всегда `false` литералом.
- **Изменяемые поля в `UPDATE`** — все колонки, кроме `id`, `created_at`, `is_deleted`.
  `user_id` считается изменяемым (кто последний правил).
- **Порядок параметров в `UPDATE`**: `$1` — всегда `id`, дальше поля по порядку из схемы.
- **Фильтры по `is_deleted`**:
  - `SELECT` активных и `UPDATE` — `AND is_deleted = false`;
  - `SELECT` удалённых — `AND is_deleted = true`;
  - `Delete` и `Undelete` — **без** фильтра по `is_deleted`, иначе восстановить запись
    будет невозможно.
- **Квотирование**: `"name"`, `"level"` и прочие зарезервированные слова — в двойных кавычках.
- Если в таблице нет какой-то колонки (например, `user_id` отсутствует
  в `dc.calculation_type`) — просто не включать её в запросы.
- Nullable-колонки не требуют особой обработки: sqlc сам сгенерирует nullable-тип.

## Колонки `*_id`, ссылающиеся на пользователя по внешнему id

Если у вызывающей стороны есть только внешний id пользователя (`dc."user".external_id`,
`uuid`), а не внутренний `bigint`-id, то колонка `<col>_id` в `Create`/`Update`
заполняется не параметром напрямую, а подзапросом:

```sql
<col>_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = $N)
```

Так сделано для `user_id` в `tables_model` (например
[query.sql](../../../datacatalogue/db/sqlc/tables_model/query.sql)) и для
`updated_by_id` в `user_domain_roles`/`user_table_roles`. Применять по задаче,
не для всех `*_id` подряд — обычные FK на другие таблицы (`table_id`, `alias_id`
и т.п.) остаются простым `$N`.

Это меняет тип параметра sqlc с `int64` на `uuid.UUID` (называется sqlc всегда
`ExternalID`, по колонке сравнения, а не по имени `<col>_id`) и требует
соответствующего поля `<col>_external_id` (`string`) в `.proto` —
см. раздел про этот спецкейс в [proto_based_on_crud.md](proto_based_on_crud.md).

## Оформление файла

Запросы группируются по таблицам, порядок таблиц — как в `schema.sql`.
Перед каждой таблицей — разделитель:

```sql
-- =========================================================
-- dc.<table_name>
-- =========================================================
```

Внутри блока порядок запросов — как в таблице выше (сначала чтение, потом запись).

## Шаблон блока

```sql
-- =========================================================
-- dc.host
-- =========================================================

-- name: GetHostById :one
SELECT *
FROM dc.host
WHERE id = $1
AND is_deleted = false;

-- name: GetHosts :many
SELECT *
FROM dc.host
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedHostById :one
SELECT *
FROM dc.host
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedHosts :many
SELECT *
FROM dc.host
WHERE is_deleted = true
ORDER BY id;

-- name: CreateHost :one
INSERT INTO dc.host
("name", description, host_env, port_env, username_env, password_env, is_deleted, created_at, updated_at, user_id)
VALUES($1, $2, $3, $4, $5, $6, false, now(), now(), $7)
RETURNING *;

-- name: UpdateHostById :one
UPDATE dc.host
SET "name"=$2, description=$3, host_env=$4, port_env=$5, username_env=$6, password_env=$7, user_id=$8, updated_at=now()
WHERE id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteHostById :exec
UPDATE dc.host
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteHostById :exec
UPDATE dc.host
SET is_deleted=false, updated_at=now()
WHERE id=$1;
```

## Проверка после генерации

Количество запросов должно быть равно `количество таблиц × 8` плюс
дополнительные выборки, если они есть:

```bash
grep -c "^-- name:" query.sql
```

## Внимание при перегенерации

Если `query.sql` уже существовал, часть старых имён запросов может исчезнуть или
измениться — Go-код, который их использует, перестанет компилироваться.
После генерации нужно прогнать `sqlc generate` и починить вызовы.
