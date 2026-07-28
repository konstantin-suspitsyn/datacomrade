# dc.database_calculation

Связка поддерживаемых типов вычислений с конкретной БД.

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK. |
| `database_cat_id` | `bigint` | no | `none` | Ссылка на БД. |
| `calculation_type_id` | `bigint` | no | `none` | Ссылка на тип вычисления. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |
| `user_id` | `bigint` | no | `none` | Пользователь. |

## Keys & Constraints
- Primary key: `database_calculation_pk` (`id`).

## Foreign Keys
- Исходящие FK не заданы отдельными миграциями (поля `database_cat_id`, `calculation_type_id`, `user_id` логически ссылочные).

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000024_create_table_database_calculation.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.database_cat]] [[dc.calculation_type]] [[dc.user]]
