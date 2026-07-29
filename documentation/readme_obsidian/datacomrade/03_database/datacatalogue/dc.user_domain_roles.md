# dc.user_domain_roles

Связка пользователя, роли домена и конкретного домена (three-way grant): кому, какая роль и на какой домен выдана.

Ранее модель прав строилась через отдельную таблицу-связку [[dc.domains_domain_roles]] (домен ↔ роль домена) плюс `user_domain_roles` (пользователь ↔ роль домена). В текущей схеме `dc.domains_domain_roles` удалена, а колонка `domain_id` добавлена прямо в `user_domain_roles`, так что связь "пользователь — роль — домен" хранится одной строкой.

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK. |
| `user_id` | `bigint` | no | `none` | Пользователь. |
| `domain_roles_id` | `bigint` | no | `none` | Роль домена. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |
| `domain_id` | `bigint` | no | `none` | Домен, на который выдана роль. |

## Keys & Constraints
- Primary key: `user_domain_roles_pk` (`id`).

## Foreign Keys
- `user_domain_roles_user_fk`: `user_id` -> [[dc.user]](`dc."user".id`) `ON UPDATE CASCADE ON DELETE RESTRICT`.
- `user_domain_roles_domain_roles_fk`: `domain_roles_id` -> [[dc.domain_roles]](`dc.domain_roles.id`).
- `user_domain_roles_domain_cat_id_fk`: `domain_id` -> [[dc.domain_cat]](`dc.domain_cat.id`).

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000031_create_table_user_domain_roles.sql` (включает все FK, включая `domain_id`, без отдельной alter-миграции).

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.user]] [[dc.domain_roles]] [[dc.domain_cat]]
