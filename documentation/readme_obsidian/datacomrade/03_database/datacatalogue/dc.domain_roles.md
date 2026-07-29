# dc.domain_roles

Справочник ролей на уровне домена.

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK. |
| `name` | `varchar(128)` | no | `none` | Уникальная роль. |
| `description` | `varchar(2000)` | no | `none` | Описание роли. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |

## Keys & Constraints
- Primary key: `domain_roles_pk` (`id`).
- Unique: `domain_roles_name_unique` (`name`).

## Foreign Keys
- Исходящие FK отсутствуют.

## Indexes
- Явно не создаются отдельными миграциями.

## Seed data
- Миграция `000000039_insert_domain_roles.sql` заполняет справочник тремя ролями: `can_read`, `can_write`, `can_grant`.

## Migration source
- Create: `datacatalogue/db/migrations/000000029_create_table_domain_roles.sql`
- Seed: `datacatalogue/db/migrations/000000039_insert_domain_roles.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.user_domain_roles]]
