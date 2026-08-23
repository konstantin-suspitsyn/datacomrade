package tablesapiv1

// Логика сверх шаблона генератора: два правила про alias_id колонки каталога.
// Первое связывает две таблицы (тип колонки и тип её алиаса), а gen.py умеет
// только правила внутри одной; второе живёт внутри запроса (show_in_ui
// требует алиас), но шаблон не знает и таких связей между полями. Дописать
// их в column_cat.go нельзя — тот перезаписывается без спроса (см. "Что
// перезаписывается, а что нет" в
// documentation/dev_instructions/crud/generate_gprpc_go.md), поэтому проверки
// живут в обёртке: ColumnCatGuard встраивает *TablesApiV1, перекрывает две
// пишущие ручки и делегирует остальные 5 RPC сервиса встроенному типу.
// Регистрируется он вместо TablesApiV1 в cmd/main.go.

import (
	"context"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// aliasReader — всё, что обёртке нужно от сервисного слоя: одна выборка.
// Узкий интерфейс вместо *TablesService, чтобы правило проверялось тестами
// без базы (репозиторий там — конкретный *tables_model.Queries).
type aliasReader interface {
	GetAliasById(ctx context.Context, id int64) (tables_model.DcAlias, error)
}

// ColumnCatGuard — реализация ColumnCatServiceServer, которая перед записью
// требует алиас у видимой в интерфейсе колонки и сверяет column_type_id
// колонки с column_type_id её алиаса.
type ColumnCatGuard struct {
	*TablesApiV1

	aliases aliasReader
}

// NewColumnCatGuard оборачивает уже собранные хендлеры: сервисный слой у
// обёртки тот же самый, отдельного состояния у неё нет.
func NewColumnCatGuard(api *TablesApiV1) *ColumnCatGuard {
	return &ColumnCatGuard{TablesApiV1: api, aliases: api.services.TablesService}
}

// Проверка на этапе компиляции: встроенный *TablesApiV1 доносит остальные
// методы сервиса и mustEmbed-заглушку, так что обёртка удовлетворяет
// интерфейсу целиком.
var _ tablesv1.ColumnCatServiceServer = (*ColumnCatGuard)(nil)

// CreateColumnCat создаёт строку dc.column_cat, проверив её связь с алиасом.
func (g *ColumnCatGuard) CreateColumnCat(ctx context.Context, req *tablesv1.CreateColumnCatRequest) (*tablesv1.CreateColumnCatResponse, error) {
	// Валидация повторяется во встроенном методе, но она чистая и дешёвая, а
	// запускать ради битого запроса чтение алиаса из базы незачем.
	if err := tablesvalidation.ValidateCreateColumnCat(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := checkAliasRequiredForUI(req.GetShowInUi(), req.AliasId); err != nil {
		return nil, err
	}

	if err := g.checkAliasColumnType(ctx, req.AliasId, req.GetColumnTypeId()); err != nil {
		return nil, err
	}

	return g.TablesApiV1.CreateColumnCat(ctx, req)
}

// UpdateColumnCatById обновляет строку dc.column_cat, проверив её связь с алиасом.
func (g *ColumnCatGuard) UpdateColumnCatById(ctx context.Context, req *tablesv1.UpdateColumnCatByIdRequest) (*tablesv1.UpdateColumnCatByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateColumnCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := checkAliasRequiredForUI(req.GetShowInUi(), req.AliasId); err != nil {
		return nil, err
	}

	if err := g.checkAliasColumnType(ctx, req.AliasId, req.GetColumnTypeId()); err != nil {
		return nil, err
	}

	return g.TablesApiV1.UpdateColumnCatById(ctx, req)
}

// checkAliasRequiredForUI требует алиас у колонки, показываемой в интерфейсе.
//
// alias_id стал необязательным (миграция 000000046): колонку можно завести
// без алиаса. Но именно алиас несёт бизнес-имя колонки, и без него в UI
// показывать нечего — так что show_in_ui = true без alias_id это ошибка
// клиента, а не пустая подпись на экране.
//
// Проверка чистая, до базы не ходит, поэтому идёт раньше сверки типов: на
// таком запросе читать алиас всё равно нечем. aliasID — поле запроса, а не
// геттер: nil («алиаса нет») обязан отличаться от переданного нуля, который
// отбивает валидация.
func checkAliasRequiredForUI(showInUI bool, aliasID *int64) error {
	if !showInUI || aliasID != nil {
		return nil
	}

	v := validator.New()
	v.AddError("alias_id", "обязателен при show_in_ui = true: колонке без алиаса нечего показывать в интерфейсе")

	return apierror.Wrap(v.Err())
}

// checkAliasColumnType требует, чтобы колонка и её алиас описывали один и тот
// же тип данных: алиас несёт column_type_id, и разойтись эти два значения не
// должны — иначе одна и та же колонка описана двумя разными типами.
//
// aliasID — поле запроса, а не геттер: alias_id необязателен, и nil («алиаса
// нет») обязан отличаться от переданного нуля, который отбивает валидация.
// Алиаса нет в базе — ошибка сервиса (NotFound), несовпадение — ошибка
// клиента (InvalidArgument).
func (g *ColumnCatGuard) checkAliasColumnType(ctx context.Context, aliasID *int64, columnTypeID int64) error {
	if aliasID == nil {
		return nil
	}

	alias, err := g.aliases.GetAliasById(ctx, *aliasID)
	if err != nil {
		return apierror.Wrap(err)
	}

	if alias.ColumnTypeID == columnTypeID {
		return nil
	}

	v := validator.New()
	v.AddError("column_type_id", fmt.Sprintf(
		"должен совпадать с column_type_id алиаса %d (%d)", *aliasID, alias.ColumnTypeID,
	))

	return apierror.Wrap(v.Err())
}
