2026-03-22 13:38

Status:

Tags: [[db]] [[datacatalogue]]

# Goose

Миграции `datacatalogue`. Общий порядок работы с миграциями по всему проекту — [[Миграции БД]].

## Запуск

```bash
task migrate:up:datacatalogue
```

Или вместе со всеми остальными микросервисами — `task migrate:up`.

Раньше миграции запускались руками из папки микросервиса (`goose up -env ../deploy/compose/datacatalogue/.env`) — так больше делать не нужно, задача сама переходит в `datacatalogue` и подхватывает `datacatalogue/.env`.

## Где что лежит

| Что | Где |
| --- | --- |
| файлы миграций | `datacatalogue/db/migrations` |
| настройки подключения (`GOOSE_*`) | `datacatalogue/.env`, генерится из `deploy/env/params.toml` — [[Работа с env файлами]] |
| история применённых миграций | таблица `tech.goose_migrations` в БД |

## Формат файла миграции

Нумерация последовательная, шестизначная: `000000045_add_alias_column_type_id.sql`. В файле обязательны обе секции — `up` и `down`:

```sql
-- +goose up

ALTER TABLE dc.alias ADD column_type_id bigint DEFAULT -1 NOT NULL;

-- +goose down

ALTER TABLE dc.alias DROP COLUMN IF EXISTS column_type_id;
```

После изменения схемы не забудь перегенерить код sqlc — [[sqlc_generation]].

### References
- [[Миграции БД]]
- [[Taskfile]]
- [[sqlc_generation]]
- [[Работа с env файлами]]
