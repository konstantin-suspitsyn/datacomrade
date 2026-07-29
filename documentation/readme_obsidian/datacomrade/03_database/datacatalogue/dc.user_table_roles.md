# dc.user_table_roles

Связка пользователя, роли таблицы и конкретной таблицы (three-way grant): кому, какая роль и на какую таблицу выдана.

Ранее модель прав строилась через отдельную таблицу-связку [[dc.tables_table_roles]] (таблица ↔ роль таблицы) плюс `user_table_roles` (пользователь ↔ роль таблицы). В текущей схеме `dc.tables_table_roles` удалена, а колонка `table_id` добавлена прямо в `user_table_roles`, так что связь "пользователь — роль — таблица" хранится одной строкой.

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK. |
| `user_id` | `bigint` | no | `none` | Пользователь. |
| `table_roles_id` | `bigint` | no | `none` | Роль таблицы. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |
| `table_id` | `bigint` | no | `none` | Таблица, на которую выдана роль. |

## Keys & Constraints
- Primary key: `user_table_roles_pk` (`id`).

## Foreign Keys
- `user_table_roles_user_fk`: `user_id` -> [[dc.user]](`dc."user".id`) `ON UPDATE CASCADE ON DELETE RESTRICT`.
- `user_table_roles_table_roles_fk`: `table_roles_id` -> [[dc.table_roles]](`dc.table_roles.id`) `ON UPDATE CASCADE ON DELETE RESTRICT`.
- `user_table_roles_table_cat_id_fk`: `table_id` -> [[dc.table_cat]](`dc.table_cat.id`).

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000032_create_table_user_table_roles.sql` (включает все FK, включая `table_id`, без отдельной alter-миграции).

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.user]] [[dc.table_roles]] [[dc.table_cat]]
