# dc.alias

Короткие псевдонимы (alias) для колонок и их описаний.
Это то, что будет показываться пользователям в UI.
Alias используется для того, чтобы понять одна и та же сущность в таблице или нет. Если к колонкам в разных таблицах присоединен один alias, значит эта колонка одинаковая. Допускаются агрегации, не допускаются фильтрации в таблицах

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK записи alias. |
| `name` | `varchar(255)` | no | `none` | Уникальное имя alias. |
| `description` | `varchar(1000)` | no | `none` | Описание alias. |
| `created_at` | `timestamp` | no | `now()` | Дата создания. |
| `updated_at` | `timestamp` | yes | `now()` | Дата обновления. |
| `is_deleted` | `bool` | no | `false` | Флаг мягкого удаления. |
| `user_id` | `int8` | no | `none` | Владелец записи. |

## Keys & Constraints
- Primary key: `alias_pk` (`id`).
- Unique: `alias_name_unique` (`name`).

## Foreign Keys
- `alias_user_fk`: `user_id` -> [[dc.user]](`dc."user".id`) `ON DELETE RESTRICT ON UPDATE CASCADE`.

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000002_create_table_alias.sql`
- FK/Alter: `datacatalogue/db/migrations/000000018_add_fk_alias.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.user]] [[dc.column_cat]]
