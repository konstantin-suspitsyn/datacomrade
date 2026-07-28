# dc.database_type

Справочник типов СУБД.

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK. |
| `name` | `varchar(128)` | no | `none` | Уникальное имя типа. |
| `db_version` | `varchar(512)` | no | `none` | Версия/семейство СУБД. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `user_id` | `int8` | no | `none` | Пользователь. |

## Keys & Constraints
- Primary key: `database_type_pk` (`id`).
- Unique: `database_type_name_unique` (`name`).

## Foreign Keys
- Исходящие FK не заданы миграциями.

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000004_create_table_database_type.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.database_cat]]
