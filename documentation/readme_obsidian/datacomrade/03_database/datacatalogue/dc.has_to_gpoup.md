# dc.has_to_group

Техническая таблица связей колонок (имя содержит опечатку в исходной схеме).

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK. |
| `column_id_a` | `int8` | no | `none` | Первая колонка в связи. |
| `column_id_b` | `int8` | no | `none` | Вторая колонка в связи. |
| `description` | `varchar(1000)` | yes | `none` | Описание связи. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `user_id` | `int8` | no | `none` | Пользователь. |

## Keys & Constraints
- Primary key: `has_to_group_pk` (`id`).

## Foreign Keys
- Исходящие FK отсутствуют в миграциях.

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000001_create_table_has_to_group.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.column_cat]] [[dc.user]]
