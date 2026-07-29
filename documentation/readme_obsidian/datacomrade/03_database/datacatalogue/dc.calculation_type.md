# dc.calculation_type

Справочник типов вычислений для колонок и зависимых расчетов.
Здесь содержатся типы вычислений, которые, в принципе бывают. Например, сумма, кол-во и т. д.

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK типа расчета. |
| `name` | `varchar(52)` | no | `none` | Уникальное имя типа. |
| `description` | `varchar(1000)` | no | `none` | Описание типа. |
| `created_at` | `timestamp` | no | `now()` | Дата создания. |
| `updated_at` | `timestamp` | no | `now()` | Дата обновления. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |

## Keys & Constraints
- Primary key: `calculation_type_pk` (`id`).
- Unique: `calculation_type_name_unique` (`name`).

## Foreign Keys
- Исходящие FK отсутствуют.

## Indexes
- Явно не создаются отдельными миграциями.

## Seed data
- Миграция `000000034_insert_calculation_types.sql` заполняет справочник значениями `sum`, `count`, `count_distinct`, `avg`, `max`, `min`.

## Migration source
- Create: `datacatalogue/db/migrations/000000023_create_table_calculation_type.sql`
- Seed: `datacatalogue/db/migrations/000000034_insert_calculation_types.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.column_cat]] [[dc.following_calculation]] [[dc.database_calculation]]
