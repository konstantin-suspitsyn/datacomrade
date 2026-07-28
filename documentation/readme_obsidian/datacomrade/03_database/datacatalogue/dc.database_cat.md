# dc.database_cat

Каталог баз данных и их подключений.

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `int8` | no | `none` | PK. |
| `name` | `varchar(255)` | no | `none` | Уникальное имя БД. |
| `host_id` | `int8` | no | `none` | Хост БД. |
| `database_type_id` | `int8` | no | `none` | Тип БД. |
| `description` | `varchar(1000)` | no | `none` | Описание. |
| `is_deleted` | `bool` | no | `false` | Удаление. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `user_id` | `int8` | no | `none` | Пользователь. |

## Keys & Constraints
- Primary key: `database_cat_pk` (`id`).
- Unique: `database_cat_name_unique` (`name`).

## Foreign Keys
- `database_cat_host_fk`: `host_id` -> [[dc.host]](`dc.host.id`).
- `database_cat_database_type_fk`: `database_type_id` -> [[dc.database_type]](`dc.database_type.id`).
- `database_cat_user_fk`: `user_id` -> [[dc.user]](`dc."user".id`).

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000003_create_table_database_cat.sql`
- FK/Alter: `datacatalogue/db/migrations/000000013_add_fk_database_cat.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.host]] [[dc.database_type]] [[dc.user]] [[dc.schema_cat]] [[dc.database_calculation]]
