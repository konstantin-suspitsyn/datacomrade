package tablesapiv1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"google.golang.org/grpc/codes"
)

// stubCatalogue подменяет сервисный слой обеими выборками сразу: алиас по id и
// постраничный список колонок с фильтрами. Запоминает, с чем его позвали, и
// считает вызовы отдельно — «до счётчика не дошли» здесь такая же значимая
// проверка, как и «до базы не дошли».
type stubCatalogue struct {
	aliasCalls int
	aliasID    int64
	alias      tables_model.DcAlias
	aliasErr   error

	filterCalls  int
	filterParams tables_model.GetColumnCatsOptionalFiltersParams
	totalItems   int64
	filterErr    error
}

func (s *stubCatalogue) GetAliasById(_ context.Context, id int64) (tables_model.DcAlias, error) {
	s.aliasCalls++
	s.aliasID = id

	if s.aliasErr != nil {
		return tables_model.DcAlias{}, s.aliasErr
	}

	return s.alias, nil
}

func (s *stubCatalogue) GetColumnCatsOptionalFilters(
	_ context.Context,
	params tables_model.GetColumnCatsOptionalFiltersParams,
) ([]tables_model.DcColumnCat, tables_model.CountGetColumnCatsOptionalFiltersRow, error) {
	s.filterCalls++
	s.filterParams = params

	if s.filterErr != nil {
		return nil, tables_model.CountGetColumnCatsOptionalFiltersRow{}, s.filterErr
	}

	return nil, tables_model.CountGetColumnCatsOptionalFiltersRow{TotalItems: s.totalItems}, nil
}

func aliasGuardWith(catalogue catalogueReader) *AliasGuard {
	// TablesApiV1 намеренно nil: все проверяемые здесь случаи возвращаются до
	// делегирования встроенному методу, и обращение к нему было бы ошибкой.
	return &AliasGuard{catalogue: catalogue}
}

// validUpdateAliasByIdRequest — запрос, который проходит валидацию целиком:
// придирки должны приходить только от сверки с колонками.
func validUpdateAliasByIdRequest() *tablesv1.UpdateAliasByIdRequest {
	return &tablesv1.UpdateAliasByIdRequest{
		Id:           7,
		Name:         "выручка",
		Description:  "сумма продаж без НДС",
		IsDeleted:    false,
		ColumnTypeId: 3,
		ExternalId:   "00000000-0000-4000-8000-000000000001",
	}
}

// TestCheckColumnsKeepTypeSameType — тип не меняется: переименование и правка
// описания не должны стоить ни одного лишнего запроса, поэтому колонки в этом
// случае не считаются вовсе.
func TestCheckColumnsKeepTypeSameType(t *testing.T) {
	catalogue := &stubCatalogue{alias: tables_model.DcAlias{ID: 7, ColumnTypeID: 3}, totalItems: 42}

	if err := aliasGuardWith(catalogue).checkColumnsKeepType(context.Background(), 7, 3); err != nil {
		t.Fatalf("checkColumnsKeepType() = %v, ожидался nil", err)
	}

	if catalogue.aliasCalls != 1 {
		t.Errorf("выборок алиаса = %d, ожидалась 1", catalogue.aliasCalls)
	}
	if catalogue.aliasID != 7 {
		t.Errorf("id алиаса = %d, ожидался 7", catalogue.aliasID)
	}
	if catalogue.filterCalls != 0 {
		t.Errorf("подсчётов колонок = %d, ожидалось 0", catalogue.filterCalls)
	}
}

// TestCheckColumnsKeepTypeNoColumns — тип меняется, но менять его не под кем:
// свободный алиас перетипизировать можно.
func TestCheckColumnsKeepTypeNoColumns(t *testing.T) {
	catalogue := &stubCatalogue{alias: tables_model.DcAlias{ID: 7, ColumnTypeID: 3}, totalItems: 0}

	if err := aliasGuardWith(catalogue).checkColumnsKeepType(context.Background(), 7, 9); err != nil {
		t.Fatalf("checkColumnsKeepType() = %v, ожидался nil", err)
	}

	if catalogue.filterCalls != 1 {
		t.Fatalf("подсчётов колонок = %d, ожидался 1", catalogue.filterCalls)
	}
}

// TestCheckColumnsKeepTypeFiltersByAlias — счётчик обязан спрашивать колонки
// именно этого алиаса. Незаполненный фильтр посчитал бы все колонки каталога,
// и тогда правило запретило бы менять тип у любого алиаса — молча и всегда.
func TestCheckColumnsKeepTypeFiltersByAlias(t *testing.T) {
	catalogue := &stubCatalogue{alias: tables_model.DcAlias{ID: 7, ColumnTypeID: 3}}

	if err := aliasGuardWith(catalogue).checkColumnsKeepType(context.Background(), 7, 9); err != nil {
		t.Fatalf("checkColumnsKeepType() = %v, ожидался nil", err)
	}

	got := catalogue.filterParams
	if got.AliasID != (sql.NullInt64{Int64: 7, Valid: true}) {
		t.Errorf("alias_id в фильтре = %+v, ожидался 7", got.AliasID)
	}
	if got.PageLimit < 1 || got.Page < 1 {
		t.Errorf("страница = %d/%d, ожидались положительные значения", got.Page, got.PageLimit)
	}

	// Остальные фильтры обязаны остаться незаданными: любой из них сузил бы
	// выборку и часть привязанных колонок осталась бы незамеченной.
	for name, valid := range map[string]bool{
		"table_id":            got.TableID.Valid,
		"column_type_id":      got.ColumnTypeID.Valid,
		"calculation_type_id": got.CalculationTypeID.Valid,
		"show_in_ui":          got.ShowInUi.Valid,
	} {
		if valid {
			t.Errorf("фильтр %s задан, ожидался пустым", name)
		}
	}
}

// TestCheckColumnsKeepTypeAttached — расхождение это ошибка клиента, и в
// сообщении обязано быть видно, что именно мешает: старый тип, новый и
// сколько колонок придётся переносить.
func TestCheckColumnsKeepTypeAttached(t *testing.T) {
	catalogue := &stubCatalogue{alias: tables_model.DcAlias{ID: 7, ColumnTypeID: 3}, totalItems: 4}

	err := aliasGuardWith(catalogue).checkColumnsKeepType(context.Background(), 7, 9)

	st := assertCode(t, err, codes.InvalidArgument)
	for _, want := range []string{"column_type_id", "3", "9", "4"} {
		if !strings.Contains(st.Message(), want) {
			t.Errorf("сообщение = %q, ожидалось упоминание %q", st.Message(), want)
		}
	}
}

// TestCheckColumnsKeepTypeAliasNotFound — несуществующий или мягко удалённый
// алиас: GetAliasById фильтрует is_deleted = false, значит оба случая NotFound.
// Сверка идёт раньше встроенного метода, поэтому 404 приходит отсюда.
func TestCheckColumnsKeepTypeAliasNotFound(t *testing.T) {
	catalogue := &stubCatalogue{aliasErr: fmt.Errorf("dc.alias id = 7: %w", customerrors.ErrNotFound)}

	err := aliasGuardWith(catalogue).checkColumnsKeepType(context.Background(), 7, 9)

	assertCode(t, err, codes.NotFound)
	if catalogue.filterCalls != 0 {
		t.Errorf("подсчётов колонок = %d, ожидалось 0", catalogue.filterCalls)
	}
}

// TestCheckColumnsKeepTypeServiceErrors — ошибка любой из двух выборок не
// должна маскироваться под претензию к запросу: клиент в них не виноват.
func TestCheckColumnsKeepTypeServiceErrors(t *testing.T) {
	tests := []struct {
		name      string
		catalogue *stubCatalogue
	}{
		{
			name:      "упала выборка алиаса",
			catalogue: &stubCatalogue{aliasErr: errors.New("соединение с базой потеряно")},
		},
		{
			name: "упал подсчёт колонок",
			catalogue: &stubCatalogue{
				alias:     tables_model.DcAlias{ID: 7, ColumnTypeID: 3},
				filterErr: errors.New("соединение с базой потеряно"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := aliasGuardWith(tt.catalogue).checkColumnsKeepType(context.Background(), 7, 9)

			assertCode(t, err, codes.Internal)
		})
	}
}

// TestAliasGuardRejectsTypeChange — ручка отбивает смену типа до делегирования:
// встроенный *TablesApiV1 здесь nil, и любой проход дальше уронил бы тест
// паникой, а не тихо записал строку.
func TestAliasGuardRejectsTypeChange(t *testing.T) {
	catalogue := &stubCatalogue{alias: tables_model.DcAlias{ID: 7, ColumnTypeID: 3}, totalItems: 1}

	req := validUpdateAliasByIdRequest()
	req.ColumnTypeId = 9

	_, err := aliasGuardWith(catalogue).UpdateAliasById(context.Background(), req)

	assertCode(t, err, codes.InvalidArgument)
	if catalogue.aliasID != 7 {
		t.Errorf("id алиаса = %d, ожидался 7", catalogue.aliasID)
	}
}

// TestAliasGuardValidatesBeforeReadingAlias — битый запрос не должен стоить
// похода в базу: валидация идёт раньше сверки.
func TestAliasGuardValidatesBeforeReadingAlias(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(req *tablesv1.UpdateAliasByIdRequest)
	}{
		{name: "пустое name", mutate: func(req *tablesv1.UpdateAliasByIdRequest) { req.Name = "" }},
		{name: "нулевой id", mutate: func(req *tablesv1.UpdateAliasByIdRequest) { req.Id = 0 }},
		{name: "нулевой column_type_id", mutate: func(req *tablesv1.UpdateAliasByIdRequest) { req.ColumnTypeId = 0 }},
		{name: "битый external_id", mutate: func(req *tablesv1.UpdateAliasByIdRequest) { req.ExternalId = "не-uuid" }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			catalogue := &stubCatalogue{alias: tables_model.DcAlias{ID: 7, ColumnTypeID: 3}}

			req := validUpdateAliasByIdRequest()
			tt.mutate(req)

			_, err := aliasGuardWith(catalogue).UpdateAliasById(context.Background(), req)

			assertCode(t, err, codes.InvalidArgument)
			if catalogue.aliasCalls != 0 {
				t.Errorf("выборок алиаса = %d, ожидалось 0", catalogue.aliasCalls)
			}
		})
	}
}

// TestAliasGuardNilRequest — nil до сверки тоже не доходит: у геттеров он
// безопасен, но читать под него алиас незачем.
func TestAliasGuardNilRequest(t *testing.T) {
	catalogue := &stubCatalogue{}

	_, err := aliasGuardWith(catalogue).UpdateAliasById(context.Background(), nil)

	assertCode(t, err, codes.InvalidArgument)
	if catalogue.aliasCalls != 0 {
		t.Errorf("выборок алиаса = %d, ожидалось 0", catalogue.aliasCalls)
	}
}
