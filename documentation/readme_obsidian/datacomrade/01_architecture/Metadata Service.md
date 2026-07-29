2026-07-27 20:45

Status:

Tags: [[architecture]] [[datacatalogue]] [[db]]

# Metadata Service

Реализован в репозитории как микросервис `datacatalogue` (схема Postgres `dc`, gRPC API: `TableService`, `UserService`, `CreateUserDomainRolesService`). Единственный источник правды о том, какие данные существуют и кому они доступны.

## Что хранит

**Подключения и структура** — [[dc.host]], [[dc.database_type]], [[dc.database_cat]], [[dc.schema_cat]], [[dc.table_type]], [[dc.table_cat]], [[dc.column_cat]], [[dc.column_type]], [[dc.alias]].

**Вычисления и группировки** — [[dc.calculation_type]], [[dc.database_calculation]], [[dc.following_calculation]], [[dc.group_levels]], [[dc.has_to_gpoup]].

**Домены, роли, пользователи** — [[dc.domain_cat]], [[dc.domain_roles]], [[dc.user_domain_roles]], [[dc.table_roles]], [[dc.user_table_roles]], [[dc.user]].

Полный индекс таблиц: [[README]] в `03_database/datacatalogue`.

## Ключевые свойства модели

- Таблица всегда принадлежит ровно одному бизнес-домену (`table_cat.domain_id` — обязательный FK). На этом держится доменная ветка прав.
- `alias` связывает «одинаковые по смыслу» колонки в разных таблицах. Это зачаток семантического слоя: по колонкам с общим alias допускаются агрегации, но не фильтрация.
- Секреты подключений в `dc` **не хранятся**. [[dc.host]] хранит имена переменных окружения (`host_env`, `port_env`, `username_env`, `password_env`), а значения подставляются в рантайме из окружения сервиса — см. [[Работа с env файлами]].
- Пользователь в [[dc.user]] — локальное зеркало с полем `external_id` (Subject/`sub` из Keycloak), по которому находится строка. Аутентификация здесь не живёт, только привязка прав.

## API

Административные методы (создание и правка host / database_cat / domain_cat / table_cat / column_cat, назначение ролей) доступны только роли `admin`. Читающий API отдаёт обычному пользователю только разрешённое ему поддерево каталога — на нём строится экран выбора полей.

## Кеш

Каталог и права — самые горячие данные системы и читаются на каждом открытии конструктора. Кешируются в Redis, инвалидация — по событию записи, а не только по TTL. Подробнее: [[Кеширование Redis]].

### References
- [[Разграничение доступа]]
- [[Кеширование Redis]]
- [[Замечания к реализации]]
