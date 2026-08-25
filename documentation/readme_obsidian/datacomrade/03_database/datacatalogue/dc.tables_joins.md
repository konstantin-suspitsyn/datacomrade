# dc.tables_joins

Пары таблиц каталога, которые разрешено соединять (join) при построении запроса: левая таблица + правая таблица.

Строка задаёт только сам факт связи между двумя записями [[dc.table_cat]]. Условие соединения (по каким колонкам идёт join, тип join) в этой таблице не хранится — в миграции `000000047` таких колонок нет.

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK. |
| `table_cat_id_left` | `bigint` | no | `none` | Левая таблица связи. |
| `table_cat_id_right` | `bigint` | no | `none` | Правая таблица связи. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `user_id` | `bigint` | no | `none` | Пользователь. |

## Keys & Constraints
- Primary key: `tables_joins_pk` (`id`).
- Unique: `tables_joins_unique` (`table_cat_id_left`, `table_cat_id_right`).

## Foreign Keys
- `table_cat_left_fk`: `table_cat_id_left` -> [[dc.table_cat]](`dc.table_cat.id`).
- `table_cat_right_fk`: `table_cat_id_right` -> [[dc.table_cat]](`dc.table_cat.id`).
- `tables_joins_user_fk`: `user_id` -> [[dc.user]](`dc."user".id`) `ON DELETE RESTRICT ON UPDATE CASCADE`.

## Indexes
- Явно не создаются отдельными миграциями. Уникальный констрейнт `tables_joins_unique` даёт индекс по паре `(left, right)`; отдельного индекса по `table_cat_id_right` нет, поэтому поиск связей «справа налево» идёт seq scan'ом.

## Поведение и ограничения схемы
- Связь направленная: пары `(A, B)` и `(B, A)` — две разные строки, уникальный констрейнт их не схлопывает. Если связь считается симметричной, дубль нужно отсекать на уровне приложения.
- Self-join не запрещён: `table_cat_id_left = table_cat_id_right` вставится без ошибки, `CHECK` на неравенство в схеме нет.
- Удаление логическое (`is_deleted`), но `tables_joins_unique` его не учитывает: пометив пару удалённой, вставить её заново не получится — конфликт по уникальности. Повторное включение связи делается через `is_deleted = false` на существующей строке.

## Migration source
- Create: `datacatalogue/db/migrations/000000047_create_join_table.sql` (включает оба FK, без отдельной alter-миграции).

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.table_cat]] [[dc.user]]
