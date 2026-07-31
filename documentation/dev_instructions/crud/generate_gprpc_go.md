# Генерация Go-слоёв gRPC по query.sql и .proto

Инструкция для Claude и для человека: как получить слои `converter`, `validation`,
`service` и `api` из уже существующих `query.sql`, sqlc-моделей и `.proto`.

Завершает цепочку:
[standard_crud.md](standard_crud.md) → [proto_based_on_crud.md](proto_based_on_crud.md) → **этот файл**.

## Что это делает

Конвейер из трёх скриптов в [deploy/generators](../../../deploy/generators):

| Шаг | Скрипт | Вход | Выход |
|---|---|---|---|
| 1 | [parse.py](../../../deploy/generators/parse.py) | `schema.sql`, `query.sql`, `models.go`, `query.sql.go`, `*.pb.go` | `build/parsed.json` |
| 2 | [resolve.py](../../../deploy/generators/resolve.py) | `build/parsed.json` | `build/model.json` |
| 3 | [gen.py](../../../deploy/generators/gen.py) | `build/model.json` | Go-файлы в `internal/` |

Настройки — в [crud_config.json](../../../deploy/generators/crud_config.json),
общие функции — в [genconfig.py](../../../deploy/generators/genconfig.py).

На таблицу генерируется 8 стандартных RPC плюс по одному на каждую
дополнительную выборку из `query.sql`. Дополнительных видов два:

| Вид | Имя запроса | Тип | Что отдаёт |
|---|---|---|---|
| FK-выборка | `Get<E>sBy<Fk>` | `:many` | список строк по неуникальной колонке |
| Выборка по уникальной колонке | `Get<E>By<Column>` | `:one` | одну строку |

Различает их генератор по имени: множественное число сущности — список,
единственное — одна строка. `Get<E>By<Column>` с `:many` считается ошибкой.

Для каждой таблицы — по файлу в каждом из четырёх слоёв.

## Запуск

Из корня репозитория, python 3 без сторонних зависимостей:

```bash
python deploy/generators/parse.py && python deploy/generators/resolve.py && python deploy/generators/gen.py
```

Затем обязательно отформатировать — генератор не выравнивает структуры,
это делает `gofmt`:

```bash
cd datacatalogue && gofmt -w internal/ && go build ./... && go vet ./... && go test ./internal/...
```

Промежуточные `parsed.json` и `model.json` лежат в `deploy/generators/build/`
и в git не попадают. Их полезно открыть, если генератор упал: там видно,
что именно он вычитал из источников.

## Что перезаписывается, а что нет

**gen.py перезаписывает без спроса** — правки в этих файлах будут потеряны:

- `internal/converter/<домен>/<таблица>.go` и `_test.go`
- `internal/validation/<домен>/<таблица>.go` и `_test.go`
- `internal/validation/<домен>/limits.go`
- `internal/service/<домен>_service/<таблица>.go`
- `internal/api/<домен>apiv1/<таблица>.go`

**Написано руками, генератор не трогает** — сюда можно вносить правки:

- `internal/converter/convert.go` — хелперы `TimeToProto`, `NullTimeToProto` и прочие
- `internal/validation/validation.go` — `ValidateID` для запросов, где в теле только id
- `internal/api/apierror/apierror.go` — перевод ошибок в коды gRPC
- `internal/api/*/init.go`, `internal/service/*/init.go`, `internal/service/services/init.go`
- `internal/utils/custom_errors/crud.go`, `internal/db/db.go`
- `cmd/main.go` — регистрация сервисов

Если понадобилась логика сверх шаблона (например, проверка прав в конкретном
RPC), её нельзя дописывать в сгенерированный файл. Варианты: положить рядом
в отдельный файл того же пакета либо поменять шаблон в `gen.py`.

## Как поменять папки

Все пути и имена — в [crud_config.json](../../../deploy/generators/crud_config.json),
менять код для этого не нужно. Пути задаются **от корня репозитория, через прямой слэш**;
корень скрипты вычисляют сами от своего расположения, абсолютных путей нигде нет.

Верхний уровень:

| Ключ | Что это | Сейчас |
|---|---|---|
| `module` | Go-модуль репозитория | `github.com/konstantin-suspitsyn/datacomrade` |
| `internal` | куда генерировать слои | `datacatalogue/internal` |
| `sqlc_root` | где лежат `schema.sql` и `query.sql` | `datacatalogue/db/sqlc` |
| `proto_go_root` | где лежат сгенерированные `*.pb.go` | `shared/pkg/proto` |
| `shared_imports` | пути импорта рукописных пакетов | см. файл |

На каждый домен в `domains` (домен = один sqlc-пакет = один gRPC-сервис):

| Ключ | Что это | Пример |
|---|---|---|
| `sqlc` | имя sqlc-пакета; оно же имя папки в `sqlc_root` и в `internal/repository` | `tables_model` |
| `proto_pkg` / `proto_file` | где взять сгенерированный proto-код | `tables/v1` / `tables.pb.go` |
| `proto_alias` / `proto_import` | как импортировать proto-пакет | `tablesv1` |
| `repo_import` / `repo_alias` | как импортировать sqlc-репозиторий | `tables_model` |
| `conv_pkg` | имя Go-пакета для converter и validation | `tables` |
| `conv_dir` / `valid_dir` | папки внутри `internal` | `converter/tables` |
| `service_dir` / `service_pkg` / `service_type` | папка, пакет и тип сервиса | `service/tables_service` |
| `service_field` | поле в `ServiceLayer` | `TablesService` |
| `repo_field` | поле репозитория внутри сервиса | `TablesRepository` |
| `api_dir` / `api_pkg` / `api_type` | папка, пакет и тип хендлеров | `api/tablesapiv1` |
| `api_recv` | буква получателя методов | `t` |

Конфиг можно не трогать в репозитории, а подсунуть свой:

```bash
python deploy/generators/gen.py --config path/to/my_config.json --work-dir path/to/build
```

### Если просите Claude поменять папки

Достаточно назвать домен и то, что меняется. Claude правит `crud_config.json`,
перезапускает конвейер и прогоняет проверки из раздела ниже. Формулировки, которые
понимаются однозначно:

- «в домене `tables_model` переложи converter в `internal/mapper/tables`»
- «переименуй пакет api для `user_domain_roles` в `udrapiv1`»
- «генерируй всё не в `datacatalogue/internal`, а в `datacatalogue/internal/gen`»
- «добавь новый домен: sqlc-пакет `audit_model`, proto `audit/v1`, сервис `AuditService`»
- «в `user_model` поменяй получателя методов api с `u` на `a`»

Старые файлы при переезде **не удаляются автоматически** — Claude должен их снести
отдельно, иначе в пакете останутся дубли функций и сборка упадёт. Об этом стоит
сказать прямо: «переложи и удали старые».

## Проверки после генерации

Зелёные тесты сами по себе ничего не доказывают: встроенный
`Unimplemented<Сервис>Server` молча подставляет заглушку вместо пропущенного
метода, поэтому пропуск RPC компиляцию не ломает. Нужны две сверки по списку.

Обе команды запускаются **из корня репозитория** и считают нулевой результат
провалом: иначе запуск не из того каталога напечатал бы `MATCH` при `0 = 0`.

Каждый RPC из `.proto` имеет хендлер:

```bash
for d in "tables:tablesapiv1:TablesApiV1" "user:userapiv1:UserApiV1" "user_domain_roles:userdomainrolesapiv1:UserDomainRolesApiV1"; do proto=$(echo $d|cut -d: -f1); pkg=$(echo $d|cut -d: -f2); typ=$(echo $d|cut -d: -f3); grep -ho "^  rpc [A-Za-z]*" shared/proto/$proto/v1/*.proto 2>/dev/null | sed 's/.*rpc //' | sort > /tmp/rpc.txt; grep -ho "^func ([a-z] \*$typ) [A-Za-z]*" datacatalogue/internal/api/$pkg/*.go 2>/dev/null | sed "s/.*) //" | sort > /tmp/impl.txt; n=$(wc -l < /tmp/rpc.txt); printf "%-24s rpc=%s impl=%s " $pkg $n $(wc -l < /tmp/impl.txt); if [ "$n" -eq 0 ]; then echo "ПУСТО — не тот каталог?"; elif diff -q /tmp/rpc.txt /tmp/impl.txt > /dev/null; then echo MATCH; else echo РАСХОЖДЕНИЕ; fi; done
```

Каждый метод репозитория обёрнут сервисом:

```bash
for d in "tables_model:tables_service:TablesService" "user_model:user_service:UserService" "user_domain_roles:user_domain_roles_service:UserDomainRolesService"; do repo=$(echo $d|cut -d: -f1); sd=$(echo $d|cut -d: -f2); typ=$(echo $d|cut -d: -f3); grep -ho "^func (q \*Queries) [A-Za-z]*" datacatalogue/internal/repository/$repo/query.sql.go 2>/dev/null | sed 's/.*) //' | sort > /tmp/r.txt; grep -ho "^func (s \*$typ) [A-Za-z]*" datacatalogue/internal/service/$sd/*.go 2>/dev/null | sed 's/.*) //' | sort > /tmp/s.txt; n=$(wc -l < /tmp/r.txt); printf "%-30s repo=%s service=%s " $sd $n $(wc -l < /tmp/s.txt); if [ "$n" -eq 0 ]; then echo "ПУСТО — не тот каталог?"; elif diff -q /tmp/r.txt /tmp/s.txt > /dev/null; then echo MATCH; else echo РАСХОЖДЕНИЕ; fi; done
```

На момент написания: 185 RPC — 128 `tables`, 9 `user`, 48 `user_domain_roles`;
22 таблицы; 861 тест.

## Когда перезапускать

| Что изменилось | Что делать до генерации |
|---|---|
| `schema.sql` | пересобрать `query.sql` по [standard_crud.md](standard_crud.md), затем `sqlc generate` |
| `query.sql` | `sqlc generate`, затем обновить `.proto` по [proto_based_on_crud.md](proto_based_on_crud.md) |
| `.proto` | `task proto:gen` (buf lint + generate) |
| только шаблон кода | ничего, достаточно `gen.py` |

Полный цикл после правки схемы:

```bash
task proto:gen && python deploy/generators/parse.py && python deploy/generators/resolve.py && python deploy/generators/gen.py
```

## Правила, зашитые в генератор

Конвейер **падает с ошибкой** при любом расхождении между источниками и ничего
не додумывает. Это осознанно: молча сгенерированный неверный маппинг найти
гораздо дороже, чем упавший скрипт. Он проверяет, что

- набор запросов таблицы — ровно 8 стандартных плюс дополнительные выборки, без лишних и пропущенных;
- каждому запросу соответствует метод репозитория;
- каждому полю sqlc-структуры находится поле в proto-сообщении;
- имя колонки из proto-тега есть в `schema.sql`;
- у каждого `Response` ровно одно поле;
- для каждой пары типов sqlc→proto есть правило.

Именно так нашлось расхождение в `dc.has_to_group.description`: колонка
`NOT NULL`, а поле в `.proto` было помечено `optional`.

Соответствие типов:

| sqlc | proto | Как переносится |
|---|---|---|
| `int64`, `string`, `bool` | тот же | напрямую |
| `int16` | `int32` | явное приведение в обе стороны |
| `time.Time` | `*timestamppb.Timestamp` | `converter.TimeToProto` |
| `sql.NullTime` | `*timestamppb.Timestamp` | `converter.NullTimeToProto`, NULL → `nil` |
| `uuid.UUID` | `string` | `converter.UUIDToProto` / `converter.ProtoToUUID` |

В protobuf нет типа для UUID, поэтому колонка `uuid` едет строкой в каноническом
виде 8-4-4-4-12. Обратное преобразование не возвращает ошибку (конвертеры —
чистые функции): нераспознанная строка становится `uuid.Nil`, а формат
проверяется раньше, на валидации.

Отдельный случай — колонка `<col>_id` (`bigint`), которая в `Create`/`Update`
заполняется через `(SELECT u.id FROM dc."user" u WHERE u.external_id = $N)`
(см. раздел про колонки `*_id` через внешний id в [standard_crud.md](standard_crud.md)).
sqlc называет такой параметр не по колонке `<col>_id`, а по колонке сравнения —
всегда `ExternalID uuid.UUID`. `resolve.py` ищет в proto-сообщении ещё не
занятое поле с суффиксом `_external_id` (`EXTERNAL_ID_SUFFIX`) и сверяет, что
`<col>_id` есть в `schema.sql`; сам генератор кода вставляет туда
`converter.ProtoToUUID`/`converter.UUIDToProto`, как для обычного `uuid.UUID`.

Правила валидации по типу колонки:

| Колонка | Проверка |
|---|---|
| `varchar(n) NOT NULL` | `StringVarchar` — непустая и не длиннее n **символов**, не байт |
| `bigint` | `Int64ID` — строго больше нуля |
| `smallint` | `Int32Between` в границах `int16` — иначе значение молча обрежет Postgres |
| `uuid` | `StringUUID` — непустая строка канонического вида 8-4-4-4-12 |
| `boolean` | не проверяется, любое значение допустимо |

Границы длин собираются в `limits.go` каждого пакета validation прямо из
`schema.sql`, руками их править не нужно.

Что **не** проверяется: существование записи по внешнему ключу. Битый `user_id`
дойдёт до базы и вернётся как ошибка вставки — это осознанное решение,
а не недосмотр.

## Известные ограничения

- Поля `Create`/`Update` берутся из sqlc-структур параметров, то есть состав
  полей задаётся `query.sql`, а не `.proto`. Если в `.proto` есть поле, которого
  нет в `INSERT`, оно будет проигнорировано молча.
- Пагинации нет: списочные запросы отдают всё. Списочные `Request` в `.proto`
  оставлены непустыми сообщениями именно чтобы её потом добавить.
- Сервисы и хендлеры тестами не покрыты — тесты генерируются только для
  конвертеров и валидации, то есть для чистых функций без базы.
- `Delete`/`Undelete` делают лишний `SELECT` перед `UPDATE`, потому что запросы
  объявлены как `:exec` и не сообщают число затронутых строк. Если переделать
  их на `:execrows`, проверку из шаблона сервиса нужно убрать.
