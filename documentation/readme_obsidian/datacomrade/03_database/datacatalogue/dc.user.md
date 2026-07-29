# dc.user

Каталог пользователей/сервисных пользователей (SQL-таблица `dc."user"`).

## Columns

| Name          | Type           | Nullable | Default | Description                         |
| ------------- | -------------- | -------- | ------- | ------------------------------------ |
| `id`          | `bigserial`    | no       | `none`  | PK.                                  |
| `name`        | `varchar(512)` | no       | `none`  | Уникальное имя пользователя.         |
| `created_at`  | `timestamp`    | no       | `now()` | Создание.                            |
| `updated_at`  | `timestamp`    | no       | `now()` | Обновление.                          |
| `is_deleted`  | `bool`         | no       | `false` | Флаг удаления.                       |
| `external_id` | `uuid`         | no       | `none`  | Subject (sub) из Keycloak.           |

> Ранее в таблицу добавлялась колонка `incoming_user_id` (ID пользователя из микросервиса с пользователями), но она была удалена в той же серии миграций и в текущей схеме отсутствует.

## Keys & Constraints
- Primary key: `user_pk` (`id`).
- Unique: `user_name_unique` (`name`).
- Unique: `user_external_id_unique` (`external_id`).

## Foreign Keys
- Исходящие FK отсутствуют.

## Indexes
- `user_external_id_idx` по `external_id`.

## Seed data
- Миграция `000000037_insert_into_user.sql` создаёт технического пользователя `Lomonosov M.` с `external_id = '00000000-0000-0000-0000-000000000000'`. Он используется как `user_id` для строк-справочников, вставляемых другими seed-миграциями (см. [[dc.table_type]], [[dc.column_type]]).

## Migration source
- Create: `datacatalogue/db/migrations/000000007_create_table_user.sql`
- Alter: `datacatalogue/db/migrations/000000033_add_external_id_to_user.sql` — добавление `external_id`, уникальный индекс.
- Историческое: `datacatalogue/db/migrations/000000035_add_incoming_user_id.sql` и `datacatalogue/db/migrations/000000036_delete_incoming_id.sql` — добавление и последующее удаление колонки `incoming_user_id` (в текущей схеме её нет).
- Seed: `datacatalogue/db/migrations/000000037_insert_into_user.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.alias]] [[dc.database_cat]] [[dc.table_type]] [[dc.schema_cat]] [[dc.domain_cat]] [[dc.table_cat]] [[dc.column_type]] [[dc.column_cat]] [[dc.group_levels]] [[dc.following_calculation]] [[dc.user_domain_roles]] [[dc.user_table_roles]]
