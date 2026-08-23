package tablesapiv1

// Вторая половина того же правила, что стережёт ColumnCatGuard: колонка
// каталога и её алиас обязаны описывать один тип данных. Сверки при записи
// колонки для этого мало — тип можно развести и с другой стороны, сменив
// column_type_id у алиаса, на который уже кто-то ссылается. Отдельный файл и
// отдельная обёртка по той же причине, что и у ColumnCatGuard: alias.go
// перезаписывается генератором без спроса (см. "Что перезаписывается, а что
// нет" в documentation/dev_instructions/crud/generate_gprpc_go.md).
// AliasGuard встраивает *TablesApiV1, перекрывает одну ручку и делегирует
// остальные 6 RPC сервиса встроенному типу; регистрируется вместо TablesApiV1
// в cmd/main.go.

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// catalogueReader — всё, что обёртке нужно от сервисного слоя: прочитать
// алиас и пересчитать ссылающиеся на него колонки. Узкий интерфейс вместо
// *TablesService, чтобы правило проверялось тестами без базы.
type catalogueReader interface {
	GetAliasById(ctx context.Context, id int64) (tables_model.DcAlias, error)
	GetColumnCatsOptionalFilters(
		ctx context.Context,
		params tables_model.GetColumnCatsOptionalFiltersParams,
	) ([]tables_model.DcColumnCat, tables_model.CountGetColumnCatsOptionalFiltersRow, error)
}

// AliasGuard — реализация AliasServiceServer, которая не даёт сменить тип
// алиаса под уже привязанными к нему колонками.
type AliasGuard struct {
	*TablesApiV1

	catalogue catalogueReader
}

// NewAliasGuard оборачивает уже собранные хендлеры: сервисный слой у обёртки
// тот же самый, отдельного состояния у неё нет.
func NewAliasGuard(api *TablesApiV1) *AliasGuard {
	return &AliasGuard{TablesApiV1: api, catalogue: api.services.TablesService}
}

// Проверка на этапе компиляции: встроенный *TablesApiV1 доносит остальные
// методы сервиса и mustEmbed-заглушку, так что обёртка удовлетворяет
// интерфейсу целиком.
var _ tablesv1.AliasServiceServer = (*AliasGuard)(nil)

// UpdateAliasById обновляет строку dc.alias, не позволяя сменой column_type_id
// расстроить уже привязанные колонки.
func (g *AliasGuard) UpdateAliasById(ctx context.Context, req *tablesv1.UpdateAliasByIdRequest) (*tablesv1.UpdateAliasByIdResponse, error) {
	// Валидация повторяется во встроенном методе, но она чистая и дешёвая, а
	// запускать ради битого запроса чтение из базы незачем.
	if err := tablesvalidation.ValidateUpdateAliasById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := g.checkColumnsKeepType(ctx, req.GetId(), req.GetColumnTypeId()); err != nil {
		return nil, err
	}

	return g.TablesApiV1.UpdateAliasById(ctx, req)
}

// checkColumnsKeepType пропускает обновление, если тип алиаса не меняется
// (переименование и правка описания не должны ничего стоить) либо менять его
// не под кем. Иначе это ошибка клиента: сменить тип алиаса и типы его колонок
// нужно вместе, а одним запросом такого нет — обновление колонок идёт своими
// ручками, каждая со своей сверкой в ColumnCatGuard.
//
// Алиаса нет в базе — ошибка сервиса (NotFound), как и во встроенном методе:
// сверка идёт первой, поэтому 404 на несуществующий id приходит отсюда.
func (g *AliasGuard) checkColumnsKeepType(ctx context.Context, aliasID, columnTypeID int64) error {
	alias, err := g.catalogue.GetAliasById(ctx, aliasID)
	if err != nil {
		return apierror.Wrap(err)
	}

	if alias.ColumnTypeID == columnTypeID {
		return nil
	}

	attached, err := g.countAttachedColumns(ctx, aliasID)
	if err != nil {
		return apierror.Wrap(err)
	}

	if attached == 0 {
		return nil
	}

	v := validator.New()
	v.AddError("column_type_id", fmt.Sprintf(
		"нельзя сменить тип алиаса с %d на %d: на алиас ссылаются активные колонки (%d), их column_type_id нужно менять вместе с ним",
		alias.ColumnTypeID, columnTypeID, attached,
	))

	return apierror.Wrap(v.Err())
}

// countAttachedColumns считает активные колонки, ссылающиеся на алиас.
//
// Отдельного запроса под такой счётчик нет, и заводить его в query.sql ради
// одной проверки не стоит: он приезжает из SG Buddy и там же будет затёрт.
// Постраничная выборка с фильтрами уже умеет всё нужное — она фильтрует по
// alias_id и по is_deleted = false, а её сервисный метод дёргает счётчик
// первым и при нулевом счётчике второго запроса не делает. PageLimit = 1
// поэтому не про размер страницы, а про то, что страница нам не нужна вовсе:
// читается только TotalItems.
//
// Мягко удалённые колонки в счёт не идут — их и не должно: правило про
// каталог, каким его видят клиенты. Обратная сторона в том, что удалённая
// колонка помнит свой прежний column_type_id, и восстановление её после
// смены типа алиаса тип разведёт. Закрыть это тем же приёмом нельзя:
// выборки удалённой строки по id у dc.column_cat нет, а заводить её ради
// одной проверки — та же правка query.sql, которую затрёт SG Buddy.
func (g *AliasGuard) countAttachedColumns(ctx context.Context, aliasID int64) (int64, error) {
	_, count, err := g.catalogue.GetColumnCatsOptionalFilters(ctx, tables_model.GetColumnCatsOptionalFiltersParams{
		AliasID:   sql.NullInt64{Int64: aliasID, Valid: true},
		Order:     "ASC",
		Page:      1,
		PageLimit: 1,
	})
	if err != nil {
		return 0, err
	}

	return count.TotalItems, nil
}
