2026-07-28 21:10

Status:

Tags: [[new_microservice]] [[db]] [[workflow]]

# standard_crud

Первый шаг цепочки генерации: из `schema.sql` получаем `query.sql` — по 8 запросов на каждую таблицу.

- Инструкция: `documentation/dev_instructions/crud/standard_crud.md`
- Вход: `datacatalogue/db/sqlc/<модель>/schema.sql`
- Выход: `datacatalogue/db/sqlc/<модель>/query.sql`

Схему для sqlc готовит [[sqlc_generation]], после этого шага запускается `sqlc generate`, а дальше идёт [[proto_based_on_crud]].

## Восемь запросов на таблицу

`<Entity>` — имя таблицы в PascalCase единственного числа (`table_cat` → `TableCat`, `group_levels` → `GroupLevel`). `<Entity>s` — множественное число по правилам английского (`Alias` → `Aliases`).

| # | Запрос | Тип | Что делает |
| --- | --- | --- | --- |
| 1 | `Get<Entity>ById` | `:one` | активная запись по id |
| 2 | `Get<Entity>s` | `:many` | все активные, `ORDER BY id` |
| 3 | `GetDeleted<Entity>ById` | `:one` | удалённая запись по id |
| 4 | `GetDeleted<Entity>s` | `:many` | все удалённые, `ORDER BY id` |
| 5 | `Create<Entity>` | `:one` | вставка, `RETURNING *` |
| 6 | `Update<Entity>ById` | `:one` | обновление по id, `RETURNING *` |
| 7 | `Delete<Entity>ById` | `:exec` | мягкое удаление |
| 8 | `Undelete<Entity>ById` | `:exec` | восстановление |

Сверх восьми добавляются выборки по внешним ключам вида `Get<Entity>sBy<Fk>` — например `GetColumnCatsByTableId`. Сейчас таких восемь на весь каталог.

## Правила, о которых легко забыть

**Удаление везде мягкое.** Физического `DELETE` в каталоге нет: `is_deleted` переключается в `true`, строка остаётся. Поэтому на каждую таблицу приходится два запроса на чтение активных записей и два — на чтение удалённых.

**`Delete` и `Undelete` не фильтруют по `is_deleted`.** Если добавить фильтр, восстановить запись станет невозможно: `Undelete` не найдёт строку, которую сам же должен вернуть к жизни.

**Времена ставит SQL, а не приложение.** `created_at` и `updated_at` проставляются через `now()` прямо в запросе. `updated_at` обновляется и при `Update`, и при `Delete`, и при `Undelete`.

**`id` в `INSERT` не передаётся** — работает последовательность `bigserial`. `is_deleted` при вставке всегда литерал `false`.

**В `UPDATE` меняются все колонки, кроме `id`, `created_at` и `is_deleted`.** `user_id` считается изменяемым: это тот, кто правил последним. Порядок параметров жёсткий — `$1` всегда `id`, дальше поля по порядку из схемы.

**Зарезервированные слова квотируются** — `"name"`, `"level"` и подобные пишутся в двойных кавычках.

## Проверка

Количество запросов должно равняться числу таблиц, умноженному на 8, плюс FK-выборки:

```bash
grep -c "^-- name:" datacatalogue/db/sqlc/tables_model/query.sql
```

Сейчас по трём моделям: 128 запросов на 15 таблиц в `tables_model`, 8 на одну таблицу в `user_model`, 48 на шесть таблиц в `user_domain_roles`.

## Осторожно при перегенерации

Если `query.sql` уже существовал, часть старых имён может исчезнуть или измениться — Go-код, который их использует, перестанет компилироваться. После генерации обязательно `sqlc generate`, а затем полная перегенерация слоёв по [[generate_gprpc_go]].

### References
- [[sqlc_generation]]
- [[proto_based_on_crud]]
- [[generate_gprpc_go]]
- [[Goose]]
- [[README]]
