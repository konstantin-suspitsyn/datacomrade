# dc.domain_cat

Каталог бизнес-доменов.

## Фиксированный домен `no_domain`

Строка `domain_name = 'no_domain'` засеяна миграцией `000000048_insert_no_domain.sql`
и присутствует в базе всегда — как справочные роли `can_read`/`can_write`/`can_grant`
и пользователь `Lomonosov M.`.

Это домен для таблиц, не относящихся ни к одной бизнес-области: общие словари,
календари, географические и прочие сквозные справочники. Доступ к таблицам этого
домена открыт **всем** пользователям и не требует записи в [[dc.user_domain_roles]]
(см. [[Разграничение доступа]]). Удалять или переименовывать `no_domain` нельзя —
на него опирается логика эффективного доступа.

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
- Seed: `datacatalogue/db/migrations/000000048_insert_no_domain.sql` (`no_domain`)

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.user]] [[dc.table_cat]] [[dc.user_domain_roles]]
