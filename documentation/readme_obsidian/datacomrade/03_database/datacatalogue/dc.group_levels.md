# dc.group_levels

Иерархия группировок/уровней между колонками.

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK. |
| `column_id` | `int8` | no | `none` | Дочерняя колонка. |
| `parent_column_id` | `int8` | no | `none` | Родительская колонка. |
| `level` | `int2` | no | `none` | Уровень вложенности. |
| `description` | `varchar(1000)` | no | `none` | Описание связи. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |
| `user_id` | `int8` | no | `none` | Пользователь. |

## Keys & Constraints
- Primary key: `group_levels_pk` (`id`).

## Foreign Keys
- `group_levels_column_cat_fk`: `column_id` -> [[dc.column_cat]](`dc.column_cat.id`).
- `group_levels_column_cat_parent_fk`: `parent_column_id` -> [[dc.column_cat]](`dc.column_cat.id`).
- `group_levels_user_fk`: `user_id` -> [[dc.user]](`dc."user".id`).

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000012_create_table_group_levels.sql`
- FK/Alter: `datacatalogue/db/migrations/000000022_add_fk_group_levels.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.column_cat]] [[dc.user]]
