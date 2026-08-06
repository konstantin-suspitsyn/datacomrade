2026-07-28 21:20

Status:

Tags: [[new_microservice]] [[workflow]] [[architecture]]

# proto_based_on_crud

Второй шаг цепочки: `.proto`, точно повторяющий набор CRUD-операций.

Отдельного действия здесь нет: контракт пишет та же программа **SG Buddy** (https://github.com/konstantin-suspitsyn/sg_buddy), что и `query.sql`, за одно сохранение и из того же `schema.json`. Скрипты `deploy/generators/proto_gen_*.py`, которые делали это раньше, удалены.

- Инструкция: `documentation/dev_instructions/crud/proto_based_on_crud.md`
- Вход: `schema.json` (настройки) и `schema.sql` (типы и состав полей)
- Выход: `shared/proto/<домен>/v1/<домен>.proto`, путь задаётся ключом `save_proto_path`

Предыдущий шаг — [[standard_crud]], следующий — [[generate_gprpc_go]].

## Главное правило

**Каждому запросу в `query.sql` соответствует ровно один RPC с тем же именем.** Количество RPC равно количеству запросов, ни лишних, ни пропущенных. Сейчас это 184 RPC: 128 в `tables.proto`, 8 в `user.proto`, 48 в `user_domain_roles.proto`.

Правило не косметическое. Именно на нём держится генератор из [[generate_gprpc_go]]: он сопоставляет запрос, метод репозитория и RPC по имени, и падает, если тройка не сходится.

## Сообщения на таблицу

**Сущность** — одна message на всю строку таблицы. Все `Response`, возвращающие данные, ссылаются на неё, а не дублируют поля.

```proto
// Host is a full row of dc.host.
message Host {
  int64 id = 1;
  string name = 2;
  ...
}
```

**Request и Response на каждый RPC:**

| RPC | Request | Response |
| --- | --- | --- |
| `Get<E>ById` | `int64 id = 1;` | `<E> <snake_e> = 1;` |
| `Get<E>s` | пустое тело `{}` | `repeated <E> <snake_plural> = 1;` |
| `GetDeleted<E>ById` | `int64 id = 1;` | `<E> <snake_e> = 1;` |
| `GetDeleted<E>s` | пустое тело `{}` | `repeated <E> <snake_plural> = 1;` |
| `Get<E>sBy<Fk>` | `int64 <fk_field> = 1;` | `repeated <E> <snake_plural> = 1;` |
| `Create<E>` | параметры `INSERT` | `<E> <snake_e> = 1;` |
| `Update<E>ById` | `id` + параметры `UPDATE` | `<E> <snake_e> = 1;` |
| `Delete<E>ById` | `int64 id = 1;` | `google.protobuf.Empty empty = 1;` |
| `Undelete<E>ById` | `int64 id = 1;` | `google.protobuf.Empty empty = 1;` |

В `Create` и `Update` **не входят** `created_at`, `updated_at`, `is_deleted` и `id` для `Create` — их выставляет SQL.

Списочные `Request` делаются пустыми сообщениями `message ...Request {}`, а не через `google.protobuf.Empty`. Это задел на пагинацию: в пустое сообщение поля добавить можно, в `Empty` — нет.

## Соответствие типов

| Postgres | proto |
| --- | --- |
| `bigint` / `int8` | `int64` |
| `integer` / `int4` | `int32` |
| `smallint` / `int2` | `int32` |
| `varchar` / `text` | `string` |
| `boolean` | `bool` |
| `timestamp` / `timestamptz` | `google.protobuf.Timestamp` |
| `numeric` | `string`, чтобы не терять точность |

Nullable-колонка → `optional <type>`. Для message-типов вроде `Timestamp` presence есть по умолчанию, `optional` не нужен.

**Осторожно с `optional`.** Ставить его нужно строго по nullability колонки в `schema.sql`. Расхождение уже случалось: `dc.has_to_group.description` объявлен `NOT NULL`, а в `.proto` поле было помечено `optional`. Контракт обещал клиенту, что поле можно не передавать, хотя база такую строку не приняла бы. Нашёл это генератор из [[generate_gprpc_go]] — он сверяет типы и падает на несоответствии.

## Комментарии обязательны

`buf lint` с правилом `COMMENT_MESSAGE` требует непустой комментарий у **каждой** message, включая `Response`. Комментарий пишется строкой выше `message`, начинается с имени сообщения и содержит ссылку на исходный запрос:

```proto
// GetHostByIdResponse returns a single active dc.host row (query GetHostById).
message GetHostByIdResponse {
  Host host = 1;
}
```

Комментарии к каждому `rpc` не добавляются — блок `service` становится нечитаемым.

## Как генерировать

Никак: контракт приезжает из SG Buddy вместе с `query.sql`. При 15+ таблицах это около 270 сообщений — руками их не набрать без опечаток, и именно поэтому нумерация полей и имена сообщений отданы программе.

Править `.proto` руками бессмысленно: следующее сохранение в программе перезапишет файл целиком. Всё, что нужно изменить, меняется в настройках схемы.

## Проверка

Имена RPC совпадают с именами запросов один в один:

```bash
grep -o "^-- name: [A-Za-z]*" query.sql | sed 's/-- name: //' | sort > /tmp/q.txt
grep -o "^  rpc [A-Za-z]*" tables.proto | sed 's/  rpc //' | sort > /tmp/r.txt
diff /tmp/q.txt /tmp/r.txt && echo MATCH
```

Затем `buf lint` и генерация Go-кода:

```bash
task proto:gen
```

Задача поднимает `buf` и плагины в `./bin`, прогоняет `buf lint` и `buf generate` — см. [[Taskdetail]].

## Осторожно при перегенерации

Перегенерация ломает API: имена RPC меняются (`GetHost` → `GetHostById`), ответы становятся вложенными (`Host host = 1;` вместо плоских полей). Хендлеры, конвертеры и клиентов придётся переписывать под новый контракт — для слоёв каталога это делается автоматически через [[generate_gprpc_go]].

### References
- [[standard_crud]]
- [[generate_gprpc_go]]
- [[Taskdetail]]
- [[Технологический стек]]
- [[API Gateway]]
