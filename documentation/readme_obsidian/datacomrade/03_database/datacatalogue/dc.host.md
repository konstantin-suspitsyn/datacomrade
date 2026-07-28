# dc.host

Каталог хостов подключений к БД.

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK. |
| `name` | `varchar(255)` | no | `none` | Уникальное имя хоста. |
| `description` | `varchar(1000)` | no | `none` | Описание хоста. |
| `host_env` | `varchar(255)` | no | `none` | Env ключ host. |
| `port_env` | `varchar(255)` | no | `none` | Env ключ port. |
| `username_env` | `varchar(255)` | no | `none` | Env ключ username. |
| `password_env` | `varchar(255)` | no | `none` | Env ключ password. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `user_id` | `int8` | no | `none` | Пользователь. |

## Keys & Constraints
- Primary key: `host_pk` (`id`).
- Unique: `host_name_unique`, `host_env_unique`, `host_port_env_unique`, `host_user_name_unique`, `host_password_env_unique`.

## Foreign Keys
- Исходящие FK отсутствуют в миграциях.

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000006_create_table_domain_cat.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.database_cat]]
