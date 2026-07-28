# dc.domain_cat

Каталог бизнес-доменов.

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK. |
| `domain_name` | `varchar(100)` | no | `none` | Уникальное имя домена. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `user_id` | `int8` | no | `none` | Пользователь. |

## Keys & Constraints
- Primary key: `domain_cat_pk` (`id`).
- Unique: `domain_cat_name_unique` (`domain_name`).

## Foreign Keys
- `domain_cat_user_fk`: `user_id` -> [[dc.user]](`dc."user".id`).

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000005_create_table_domain_cat.sql`
- FK/Alter: `datacatalogue/db/migrations/000000016_add_fk_domain_cat.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.user]] [[dc.table_cat]] [[dc.domains_domain_roles]]
