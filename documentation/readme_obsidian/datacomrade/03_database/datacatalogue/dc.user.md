# dc.user

Каталог пользователей/сервисных пользователей (SQL-таблица `dc."user"`).

## Columns

| Name               | Type           | Nullable | Default | Description                                      |
| ------------------ | -------------- | -------- | ------- | ------------------------------------------------ |
| `id`               | `bigserial`    | no       | `none`  | PK.                                              |
| `name`             | `varchar(512)` | no       | `none`  | Уникальное имя пользователя.                     |
| `created_at`       | `timestamp`    | no       | `now()` | Создание.                                        |
| `updated_at`       | `timestamp`    | no       | `now()` | Обновление.                                      |
| `is_deleted`       | `bool`         | no       | `false` | Флаг удаления.                                   |
| `incoming_user_id` | `bigint`       | no       | -1      | ID пользователя из микросервиса с пользователями |

## Keys & Constraints
- Primary key: `user_pk` (`id`).
- Unique: `user_name_unique` (`name`).

## Foreign Keys
- Исходящие FK отсутствуют.

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000007_create_table_user.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.alias]] [[dc.database_cat]] [[dc.table_type]] [[dc.schema_cat]] [[dc.domain_cat]] [[dc.table_cat]] [[dc.column_type]] [[dc.column_cat]] [[dc.group_levels]] [[dc.following_calculation]] [[dc.user_domain_roles]] [[dc.user_table_roles]]
