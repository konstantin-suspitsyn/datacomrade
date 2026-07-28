# dc.column_type

Справочник типов колонок.

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK. |
| `name` | `varchar(128)` | no | `none` | Уникальное имя типа. |
| `description` | `varchar(1000)` | no | `none` | Описание типа. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `user_id` | `int8` | no | `none` | Пользователь. |

## Keys & Constraints
- Primary key: `column_type_pk` (`id`).
- Unique: `column_type_name_unique` (`name`).

## Foreign Keys
- `column_type_user_fk`: `user_id` -> [[dc.user]](`dc."user".id`).

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000011_create_table_column_type.sql`
- FK/Alter: `datacatalogue/db/migrations/000000019_add_fk_column_type.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.user]] [[dc.column_cat]]
