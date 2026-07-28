# dc.domains_domain_roles

Таблица связей доменов и ролей доменов (many-to-many).

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK. |
| `domain_cat_id` | `bigint` | no | `none` | Ссылка на домен. |
| `domain_roles_id` | `bigint` | no | `none` | Ссылка на роль домена. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |

## Keys & Constraints
- Primary key: `domains_domain_roles_pk` (`id`).

## Foreign Keys
- `domains_domain_roles_domain_cat_fk`: `domain_cat_id` -> [[dc.domain_cat]](`dc.domain_cat.id`).
- `domains_domain_roles_domain_roles_fk`: `domain_roles_id` -> [[dc.domain_roles]](`dc.domain_roles.id`).

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000034_create_table_domains_domain_roles.sql`
- FK/Alter: `datacatalogue/db/migrations/000000035_add_fk_domains_domain_roles.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.domain_cat]] [[dc.domain_roles]]
