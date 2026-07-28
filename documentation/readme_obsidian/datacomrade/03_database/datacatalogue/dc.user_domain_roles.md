# dc.user_domain_roles

Связка пользователей и ролей доменов.

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK. |
| `user_id` | `bigint` | no | `none` | Пользователь. |
| `domain_roles_id` | `bigint` | no | `none` | Роль домена. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |

## Keys & Constraints
- Primary key: `user_domain_roles_pk` (`id`).

## Foreign Keys
- `user_domain_roles_user_fk`: `user_id` -> [[dc.user]](`dc."user".id`).
- `user_domain_roles_domain_roles_fk`: `domain_roles_id` -> [[dc.domain_roles]](`dc.domain_roles.id`).

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000031_create_table_user_domain_roles.sql`
- FK/Alter: `datacatalogue/db/migrations/000000037_add_fk_user_domain_roles.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.user]] [[dc.domain_roles]]
