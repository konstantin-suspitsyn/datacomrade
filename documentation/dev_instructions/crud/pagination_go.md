# Пагинация: шаг 3, Go-слои микросервиса

`converter`, `validation`, `service` и `api` для постраничных выборок делает
конвейер репозитория ([generate_gprpc_go.md](generate_gprpc_go.md)) — он в
SG Buddy не входит и остаётся здесь.

Третий шаг цепочки пагинации:
[pagination_sql.md](pagination_sql.md) → [pagination_proto.md](pagination_proto.md) → **этот файл**.

Ниже — план. Как конвейер определяет пагинацию и опциональные операции
на практике (домен `tables_model`) — [optional_standard_ops.md](optional_standard_ops.md).

## Вход

- `query.sql` и sqlc-код: `task sqlc:gen` уже прогнан;
- `.proto` и `*.pb.go`: `task proto:gen` уже прогнан;
- и то и другое сгенерировано SG Buddy — руками эти файлы не правятся.

Постраничность не заводит новых имён: выборка `GetAliases` остаётся собой,
просто получает параметры страницы и парный счётчик `CountGetAliases`.
Отдельного RPC у счётчика нет.

## Что дописать в конвейер (один раз)

| Скрипт | Что добавить |
|---|---|
| `parse.py` | ничего: `Get<P>Params` и `CountGet<P>Row` уже попадают в `param_structs` |
| `resolve.py` | считать `Count*` не самостоятельным запросом, а счётчиком своей выборки; ослабить «в Response ровно одно поле» — у постраничного их три (`rows`, `total_rows`, `total_pages`); поля `page_limit`, `page_offset` и `order_by_*` не гонять через сопоставление с колонками `schema.sql` — колонок с такими именами нет |
| `gen.py` | шаблоны четырёх слоёв (ниже) |

Без этих правок конвейер **упадёт на первой же постраничной выборке**, и это
правильное поведение: лишний `CountGetAliases` ломает разбор набора запросов
таблицы, а у ответа три поля вместо одного.

`resolve.py` обязан падать, а не догадываться: есть постраничная выборка без
парного `Count<имя>` (или наоборот), в ответе не хватает счётчика, у выборки
есть `page_limit`, но нет `ORDER BY` → ошибка с именем таблицы.

## Границы страницы

Значения по умолчанию и потолок живут в рукописном коде — в настройках схемы
их нет, потому что это политика сервиса, а не описание запроса:

```go
const (
	DefaultPageSize = 50
	MaxPageSize     = 200
)
```

Место — `internal/validation/validation.go` (рукописный, генератор его не
трогает). Шлюз держит те же числа у себя.

## Шаблоны генерируемого кода

**converter** — пустой размер страницы заменяется умолчанием здесь, чтобы
ответ и запрос к базе видели одно и то же число:

```go
// ToGetAliasesParams собирает параметры страницы dc.alias из запроса gRPC.
func ToGetAliasesParams(req *tablesv1.GetAliasesRequest) tables_model.GetAliasesParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	return tables_model.GetAliasesParams{
		PageLimit:  limit,
		PageOffset: req.GetPageOffset(),
	}
}
```

Если у выборки есть необязательные сортировки, в параметры едут булевы поля
`OrderByName`, `OrderByCreatedAt` и так далее — по одному на отмеченную колонку,
имена совпадают с полями proto.

**validation** — проверять нечего, кроме границ страницы: сортировка приходит
булевыми флагами, и «неизвестное значение» здесь невозможно по построению.

```go
// ValidateGetAliases проверяет запрос страницы dc.alias.
func ValidateGetAliases(req *tablesv1.GetAliasesRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int32Between("page_limit", req.GetPageLimit(), 0, validation.MaxPageSize)
	v.Int32Min("page_offset", req.GetPageOffset(), 0)

	return v.Err()
}
```

Нижняя граница `page_limit` — 0, а не 1: ноль означает «размер по умолчанию»
и заменяется в конвертере.

**service** — счётчик первым; при нулевом результате второй запрос не делается:

```go
// GetAliases возвращает страницу строк dc.alias и её счётчики.
func (s *TablesService) GetAliases(
	ctx context.Context,
	params tables_model.GetAliasesParams,
) ([]tables_model.DcAlias, tables_model.CountGetAliasesRow, error) {
	count, err := s.TablesRepository.CountGetAliases(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetAliasesRow{}, fmt.Errorf("count dc.alias: %w", err)
	}

	if count.TotalRows == 0 {
		return []tables_model.DcAlias{}, count, nil
	}

	rows, err := s.TablesRepository.GetAliases(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetAliasesRow{}, fmt.Errorf("get dc.alias page: %w", err)
	}

	return rows, count, nil
}
```

Это единственный метод сервиса, который дёргает репозиторий дважды. Пустая
страница — не ошибка: `sql.ErrNoRows` здесь не обрабатывается.

**api**:

```go
// GetAliases отдаёт страницу строк dc.alias.
func (t *TablesApiV1) GetAliases(ctx context.Context, req *tablesv1.GetAliasesRequest) (*tablesv1.GetAliasesResponse, error) {
	if err := tablesvalidation.ValidateGetAliases(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetAliasesParams(req)

	rows, count, err := t.services.TablesService.GetAliases(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetAliasesResponse{
		Rows:       tablesconv.AliasesToProto(rows),
		TotalRows:  count.TotalRows,
		TotalPages: count.TotalPages,
	}, nil
}
```

## Запуск

Из корня репозитория:

```bash
python deploy/generators/parse.py && python deploy/generators/resolve.py && python deploy/generators/gen.py
```

Затем как обычно:

```bash
cd datacatalogue && gofmt -w internal/ && go build ./... && go vet ./... && go test ./internal/...
```

## Проверки

Счётчики не остались без обёртки, а выборки — без счётчиков:

```bash
grep -c "^-- name: Count" datacatalogue/db/sqlc/tables_model/query.sql; grep -c "total_pages" shared/proto/tables/v1/tables.proto
```

Числа обязаны совпасть. Плюс обе сверки из
[generate_gprpc_go.md](generate_gprpc_go.md#проверки-после-генерации) — с учётом
того, что метод `Count<имя>` репозитория обёрнут не своим методом сервиса, а
вызовом внутри постраничного.

## Шлюз

До фронтенда пагинация доезжает через `shepherd`: эндпоинт списка принимает
`?page_limit=&page_offset=` плюс флаги сортировки и отдаёт `rows` вместе с
`total_rows` и `total_pages`.

`shepherd/api/openapi/v1/openapi.yaml` пишется руками — соответствие контракту
проверяется глазами. Генератор спеки шлюза из настроек схемы в SG Buddy пока
не входит.
