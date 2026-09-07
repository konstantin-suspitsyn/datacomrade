# dc.column_cat

Каталог колонок бизнес-таблиц и их метаданных.

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | Идентификатор колонки. |
| `table_id` | `bigint` | no | `none` | Ссылка на таблицу. |
| `name (256)` | `varchar` | no | `none` | Имя колонки (quoted SQL-имя). |
| `alias_id` | `bigint` | yes | `none` | Ссылка на alias. Nullable с миграции `000000046` — колонка без бизнес-имени; при `show_in_ui = true` всё равно обязателен (`ColumnCatGuard`). |
| `column_type_id` | `bigint` | no | `none` | Ссылка на тип колонки. |
| `description` | `varchar(1000)` | no | `none` | Описание. |
| `calculation_type_id` | `bigint` | yes | `none` | Тип вычисления. Nullable с миграции `000000049` — у атрибута измерения / элемента словаря агрегата нет. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |
| `show_in_ui` | `bool` | no | `none` | Показывать в UI. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `user_id` | `bigint` | no | `none` | Пользователь. |

## Keys & Constraints
- Primary key: `column_cat_pk` (`id`) добавлен в alter-миграции.
- Unique: `column_cat_name_unique` (`name (256)`).

## Foreign Keys
- `column_cat_table_cat_fk`: `table_id` -> [[dc.table_cat]](`dc.table_cat.id`).
- `column_cat_alias_fk`: `alias_id` -> [[dc.alias]](`dc.alias.id`).
- `column_cat_column_type_fk`: `column_type_id` -> [[dc.column_type]](`dc.column_type.id`).
- `column_cat_user_fk`: `user_id` -> [[dc.user]](`dc."user".id`).
- `column_cat_calculation_type_fk`: `calculation_type_id` -> [[dc.calculation_type]](`dc.calculation_type.id`).

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000020_create_table_column_cat.sql`
- FK/Alter: `datacatalogue/db/migrations/000000021_add_fk_column_cat.sql`, `datacatalogue/db/migrations/000000027_add_fk_column_cat_1.sql`
- `alias_id` nullable: `datacatalogue/db/migrations/000000046_alter_column_cat_alias_id_nullable.sql`
- `calculation_type_id` nullable: `datacatalogue/db/migrations/000000049_alter_column_cat_calculation_type_id_nullable.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.table_cat]] [[dc.alias]] [[dc.column_type]] [[dc.calculation_type]] [[dc.user]] [[dc.group_levels]] [[dc.following_calculation]]
