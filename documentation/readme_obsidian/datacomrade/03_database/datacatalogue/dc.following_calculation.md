# dc.following_calculation

Зависимые вычисления, привязанные к колонкам.

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK. |
| `column_cat_id` | `bigint` | no | `none` | Ссылка на колонку. |
| `calculation_type_id` | `bigint` | no | `none` | Тип вычисления. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |
| `user_id` | `bigint` | no | `none` | Пользователь. |

## Keys & Constraints
- Primary key: `following_calculation_pk` (`id`).

## Foreign Keys
- `following_calculation_column_cat_fk`: `column_cat_id` -> [[dc.column_cat]](`dc.column_cat.id`).
- `following_calculation_calculation_type_fk`: `calculation_type_id` -> [[dc.calculation_type]](`dc.calculation_type.id`).
- `following_calculation_user_fk`: `user_id` -> [[dc.user]](`dc."user".id`).

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000026_create_table_following_calculation.sql`
- FK/Alter: `datacatalogue/db/migrations/00000028_add_fk_following_calculation.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.column_cat]] [[dc.calculation_type]] [[dc.user]]
