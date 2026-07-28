# dc.table_type

Справочник типов таблиц.
На данный момент, бывают 2 видов:
1. fact
2. dimension

API для измения делать не буду. Это мини-справочник. Один раз забили и хватит

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK. |
| `name` | `varchar(128)` | no | `none` | Уникальный тип. |
| `description` | `varchar(1000)` | no | `none` | Описание. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `user_id` | `int8` | no | `none` | Пользователь. |

## Keys & Constraints
- Primary key: `table_type_pk` (`id`).
- Unique: `table_type_name_unique` (`name`).

## Foreign Keys
- `table_type_user_fk`: `user_id` -> [[dc.user]](`dc."user".id`).

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000008_create_table _table_type.sql`
- FK/Alter: `datacatalogue/db/migrations/000000014_add_fk_table_type_user.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.user]] [[dc.table_cat]]
