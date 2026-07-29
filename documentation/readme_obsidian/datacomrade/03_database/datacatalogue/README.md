# datacatalogue tables

Индекс заметок по таблицам схемы `dc`.

- [[dc.alias]]
- [[dc.calculation_type]]
- [[dc.column_cat]]
- [[dc.column_type]]
- [[dc.database_calculation]]
- [[dc.database_cat]]
- [[dc.database_type]]
- [[dc.domain_cat]]
- [[dc.domain_roles]]
- [[dc.following_calculation]]
- [[dc.group_levels]]
- [[dc.has_to_group]]
- [[dc.host]]
- [[dc.schema_cat]]
- [[dc.table_cat]]
- [[dc.table_roles]]
- [[dc.table_type]]
- [[dc.user]]
- [[dc.user_domain_roles]]
- [[dc.user_table_roles]]

## Модель прав доступа

Права на домены и таблицы выдаются напрямую через связки [[dc.user_domain_roles]] (пользователь + роль домена + домен) и [[dc.user_table_roles]] (пользователь + роль таблицы + таблица). Промежуточные таблицы `dc.domains_domain_roles` и `dc.tables_table_roles`, ранее связывавшие роли с доменами/таблицами отдельно от пользователя, из схемы удалены — колонки `domain_id`/`table_id` перенесены прямо в таблицы связок с пользователем.

## Seed-миграции

Часть справочников заполняется данными сразу в миграциях (а не через приложение):
- [[dc.calculation_type]] — `sum`, `count`, `count_distinct`, `avg`, `max`, `min`.
- [[dc.table_type]] — `fact`, `dimension`.
- [[dc.domain_roles]] / [[dc.table_roles]] — `can_read`, `can_write`, `can_grant`.
- [[dc.column_type]] — базовые типы колонок (`int`, `bigint`, `decimal`, `float`, `bool`, `string`, `text`, `date`, `datetime`, `datetime_tz`).

Для этих seed-миграций нужен технический пользователь — он создаётся в `000000037_insert_into_user.sql` (`Lomonosov M.`) и используется как `user_id` в справочниках, где это поле обязательно. См. [[dc.user]].
