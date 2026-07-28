# dc.table_cat

Каталог таблиц источников данных.

## Columns

| Name            | Type            | Nullable | Default | Description                                                                             |
| --------------- | --------------- | -------- | ------- | --------------------------------------------------------------------------------------- |
| `id`            | `bigserial`     | no       | `none`  | PK.                                                                                     |
| `name`          | `varchar(128)`  | no       | `none`  | Уникальное имя таблицы.                                                                 |
| `description`   | `varchar(2000)` | no       | `none`  | Описание.                                                                               |
| `schema_id`     | `int8`          | no       | `none`  | Схема.                                                                                  |
| `table_type_id` | `int8`          | no       | `none`  | Тип таблицы.                                                                            |
| `domain_id`     | `int8`          | no       | `none`  | Бизнес-домен.                                                                           |
| `is_deleted`    | `bool`          | no       | `false` | Флаг удаления.                                                                          |
| `created_at`    | `timestamp`     | no       | `now()` | Создание.                                                                               |
| `updated_at`    | `timestamp`     | no       | `now()` | Обновление.                                                                             |
| `user_id`       | `int8`          | no       | `none`  | Пользователь.                                                                           |
| `id_dict_get`   | `bool`          | no       | `false` | Только для ClickHouse. Чтобы использовать getDict вместо join. В остальном случае false |

## Keys & Constraints
- Primary key: `table_cat_pk` (`id`).
- Unique: `table_cat_name_unique` (`name`).

## Foreign Keys
- `table_cat_user_fk`: `user_id` -> [[dc.user]](`dc."user".id`).
- `table_cat_schema_cat_fk`: `schema_id` -> [[dc.schema_cat]](`dc.schema_cat.id`).
- `table_cat_table_type_fk`: `table_type_id` -> [[dc.table_type]](`dc.table_type.id`).
- `table_cat_domain_cat_fk`: `domain_id` -> [[dc.domain_cat]](`dc.domain_cat.id`).

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000009_create_table_cat.sql`
- FK/Alter: `datacatalogue/db/migrations/000000017_add_fk_table_cat.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.user]] [[dc.schema_cat]] [[dc.table_type]] [[dc.domain_cat]] [[dc.column_cat]] [[dc.tables_table_roles]]
