package tablesapiv1

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// stubAliases подменяет сервисный слой: отдаёт заданный алиас либо ошибку и
// запоминает, с каким id его позвали. Счётчик нужен для проверок «до базы не
// дошли»: без alias_id и на битом запросе выборки быть не должно.
type stubAliases struct {
	calls int
	gotID int64

	alias tables_model.DcAlias
	err   error
}

func (s *stubAliases) GetAliasById(_ context.Context, id int64) (tables_model.DcAlias, error) {
	s.calls++
	s.gotID = id

	if s.err != nil {
		return tables_model.DcAlias{}, s.err
	}

	return s.alias, nil
}

func guardWith(aliases aliasReader) *ColumnCatGuard {
	// TablesApiV1 намеренно nil: все проверяемые здесь случаи возвращаются до
	// делегирования встроенному методу, и обращение к нему было бы ошибкой.
	return &ColumnCatGuard{aliases: aliases}
}

func ptr[T any](v T) *T {
	return &v
}

// validCreateColumnCatRequest — запрос, который проходит валидацию целиком:
// придирки должны приходить только от сверки с алиасом.
func validCreateColumnCatRequest() *tablesv1.CreateColumnCatRequest {
	return &tablesv1.CreateColumnCatRequest{
		TableId:           100,
		Name:              "revenue",
		AliasId:           ptr(int64(7)),
		ColumnTypeId:      3,
		Description:       "сумма продаж",
		CalculationTypeId: 5,
		ShowInUi:          true,
		ExternalId:        "00000000-0000-4000-8000-000000000001",
	}
}

func validUpdateColumnCatByIdRequest() *tablesv1.UpdateColumnCatByIdRequest {
	return &tablesv1.UpdateColumnCatByIdRequest{
		Id:                11,
		TableId:           100,
		Name:              "revenue",
		AliasId:           ptr(int64(7)),
		ColumnTypeId:      3,
		Description:       "сумма продаж",
		CalculationTypeId: 5,
		ShowInUi:          true,
		ExternalId:        "00000000-0000-4000-8000-000000000001",
	}
}

func assertCode(t *testing.T, err error, want codes.Code) *status.Status {
	t.Helper()

	if err == nil {
		t.Fatalf("error = nil, ожидался %s", want)
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("error = %v, ожидался статус gRPC", err)
	}
	if st.Code() != want {
		t.Fatalf("код = %s, ожидался %s (сообщение = %q)", st.Code(), want, st.Message())
	}

	return st
}

// TestCheckAliasRequiredForUI — таблица всех четырёх сочетаний флага и
// алиаса: запрет ровно один, показ без алиаса.
func TestCheckAliasRequiredForUI(t *testing.T) {
	tests := []struct {
		name     string
		showInUI bool
		aliasID  *int64
		wantErr  bool
	}{
		{name: "показ с алиасом", showInUI: true, aliasID: ptr(int64(7))},
		{name: "показ без алиаса", showInUI: true, aliasID: nil, wantErr: true},
		{name: "без показа с алиасом", showInUI: false, aliasID: ptr(int64(7))},
		{name: "без показа без алиаса", showInUI: false, aliasID: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkAliasRequiredForUI(tt.showInUI, tt.aliasID)

			if !tt.wantErr {
				if err != nil {
					t.Fatalf("checkAliasRequiredForUI() = %v, ожидался nil", err)
				}

				return
			}

			st := assertCode(t, err, codes.InvalidArgument)
			for _, want := range []string{"alias_id", "show_in_ui"} {
				if !strings.Contains(st.Message(), want) {
					t.Errorf("сообщение = %q, ожидалось упоминание %q", st.Message(), want)
				}
			}
		})
	}
}

// TestColumnCatGuardRequiresAliasForUI — обе пишущие ручки отбивают показ без
// алиаса, и до базы такой запрос не доходит: сверять column_type_id не с чем.
func TestColumnCatGuardRequiresAliasForUI(t *testing.T) {
	tests := []struct {
		name string
		call func(g *ColumnCatGuard) error
	}{
		{
			name: "create",
			call: func(g *ColumnCatGuard) error {
				req := validCreateColumnCatRequest()
				req.ShowInUi = true
				req.AliasId = nil
				_, err := g.CreateColumnCat(context.Background(), req)

				return err
			},
		},
		{
			name: "update",
			call: func(g *ColumnCatGuard) error {
				req := validUpdateColumnCatByIdRequest()
				req.ShowInUi = true
				req.AliasId = nil
				_, err := g.UpdateColumnCatById(context.Background(), req)

				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aliases := &stubAliases{}

			err := tt.call(guardWith(aliases))

			st := assertCode(t, err, codes.InvalidArgument)
			if !strings.Contains(st.Message(), "alias_id") {
				t.Errorf("сообщение = %q, ожидалось упоминание %q", st.Message(), "alias_id")
			}
			if aliases.calls != 0 {
				t.Errorf("выборок алиаса = %d, ожидалось 0", aliases.calls)
			}
		})
	}
}

// TestColumnCatGuardAllowsHiddenColumnWithoutAlias — обратная сторона правила:
// скрытая колонка без алиаса законна, и обёртка её не задерживает. Встроенный
// *TablesApiV1 здесь nil, поэтому проход дальше виден как паника — её и ловим,
// чтобы отличить «пропустили» от «отбили».
func TestColumnCatGuardAllowsHiddenColumnWithoutAlias(t *testing.T) {
	tests := []struct {
		name string
		call func(g *ColumnCatGuard) error
	}{
		{
			name: "create",
			call: func(g *ColumnCatGuard) error {
				req := validCreateColumnCatRequest()
				req.ShowInUi = false
				req.AliasId = nil
				_, err := g.CreateColumnCat(context.Background(), req)

				return err
			},
		},
		{
			name: "update",
			call: func(g *ColumnCatGuard) error {
				req := validUpdateColumnCatByIdRequest()
				req.ShowInUi = false
				req.AliasId = nil
				_, err := g.UpdateColumnCatById(context.Background(), req)

				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aliases := &stubAliases{}
			delegated := false

			func() {
				defer func() {
					// Паника здесь означает обращение к nil-*TablesApiV1, то
					// есть обе проверки запрос пропустили — именно этого и ждём.
					if recover() != nil {
						delegated = true
					}
				}()

				if err := tt.call(guardWith(aliases)); err != nil {
					t.Errorf("вызов = %v, ожидался проход дальше", err)
				}
			}()

			if !delegated {
				t.Error("запрос не дошёл до встроенного метода")
			}
			if aliases.calls != 0 {
				t.Errorf("выборок алиаса = %d, ожидалось 0", aliases.calls)
			}
		})
	}
}

// TestCheckAliasColumnTypeSkipsAbsentAlias — alias_id необязателен: без него
// сверять нечего, и выборки быть не должно вовсе.
func TestCheckAliasColumnTypeSkipsAbsentAlias(t *testing.T) {
	aliases := &stubAliases{}

	if err := guardWith(aliases).checkAliasColumnType(context.Background(), nil, 3); err != nil {
		t.Fatalf("checkAliasColumnType() = %v, ожидался nil", err)
	}

	if aliases.calls != 0 {
		t.Errorf("выборок алиаса = %d, ожидалось 0", aliases.calls)
	}
}

// TestCheckAliasColumnTypeMatches — совпадение пропускает запрос дальше, и id
// в выборку уходит именно тот, что пришёл в запросе.
func TestCheckAliasColumnTypeMatches(t *testing.T) {
	aliases := &stubAliases{alias: tables_model.DcAlias{ID: 7, ColumnTypeID: 3}}

	if err := guardWith(aliases).checkAliasColumnType(context.Background(), ptr(int64(7)), 3); err != nil {
		t.Fatalf("checkAliasColumnType() = %v, ожидался nil", err)
	}

	if aliases.calls != 1 {
		t.Fatalf("выборок алиаса = %d, ожидалась 1", aliases.calls)
	}
	if aliases.gotID != 7 {
		t.Errorf("id алиаса = %d, ожидался 7", aliases.gotID)
	}
}

// TestCheckAliasColumnTypeMismatch — расхождение это ошибка клиента, и в
// сообщении обязано быть видно, какое значение ожидалось.
func TestCheckAliasColumnTypeMismatch(t *testing.T) {
	aliases := &stubAliases{alias: tables_model.DcAlias{ID: 7, ColumnTypeID: 9}}

	err := guardWith(aliases).checkAliasColumnType(context.Background(), ptr(int64(7)), 3)

	st := assertCode(t, err, codes.InvalidArgument)
	for _, want := range []string{"column_type_id", "7", "9"} {
		if !strings.Contains(st.Message(), want) {
			t.Errorf("сообщение = %q, ожидалось упоминание %q", st.Message(), want)
		}
	}
}

// TestCheckAliasColumnTypeAliasNotFound — несуществующий или мягко удалённый
// алиас: GetAliasById фильтрует is_deleted = false, значит оба случая NotFound.
func TestCheckAliasColumnTypeAliasNotFound(t *testing.T) {
	aliases := &stubAliases{err: fmt.Errorf("dc.alias id = 7: %w", customerrors.ErrNotFound)}

	err := guardWith(aliases).checkAliasColumnType(context.Background(), ptr(int64(7)), 3)

	assertCode(t, err, codes.NotFound)
}

// TestCheckAliasColumnTypeServiceError — прочая ошибка сервиса не должна
// маскироваться под претензию к запросу.
func TestCheckAliasColumnTypeServiceError(t *testing.T) {
	aliases := &stubAliases{err: errors.New("соединение с базой потеряно")}

	err := guardWith(aliases).checkAliasColumnType(context.Background(), ptr(int64(7)), 3)

	assertCode(t, err, codes.Internal)
}

// TestColumnCatGuardRejectsMismatch — обе пишущие ручки отбивают расхождение
// до делегирования: встроенный *TablesApiV1 здесь nil, и любой проход дальше
// уронил бы тест паникой, а не тихо записал строку.
func TestColumnCatGuardRejectsMismatch(t *testing.T) {
	tests := []struct {
		name string
		call func(g *ColumnCatGuard) error
	}{
		{
			name: "create",
			call: func(g *ColumnCatGuard) error {
				_, err := g.CreateColumnCat(context.Background(), validCreateColumnCatRequest())
				return err
			},
		},
		{
			name: "update",
			call: func(g *ColumnCatGuard) error {
				_, err := g.UpdateColumnCatById(context.Background(), validUpdateColumnCatByIdRequest())
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aliases := &stubAliases{alias: tables_model.DcAlias{ID: 7, ColumnTypeID: 9}}

			err := tt.call(guardWith(aliases))

			assertCode(t, err, codes.InvalidArgument)
			if aliases.calls != 1 {
				t.Errorf("выборок алиаса = %d, ожидалась 1", aliases.calls)
			}
			if aliases.gotID != 7 {
				t.Errorf("id алиаса = %d, ожидался 7", aliases.gotID)
			}
		})
	}
}

// TestColumnCatGuardValidatesBeforeReadingAlias — битый запрос не должен
// стоить похода в базу: валидация идёт раньше сверки.
func TestColumnCatGuardValidatesBeforeReadingAlias(t *testing.T) {
	tests := []struct {
		name string
		call func(g *ColumnCatGuard) error
	}{
		{
			name: "create с пустым name",
			call: func(g *ColumnCatGuard) error {
				req := validCreateColumnCatRequest()
				req.Name = ""
				_, err := g.CreateColumnCat(context.Background(), req)

				return err
			},
		},
		{
			name: "create с alias_id = 0",
			call: func(g *ColumnCatGuard) error {
				req := validCreateColumnCatRequest()
				req.AliasId = ptr(int64(0))
				_, err := g.CreateColumnCat(context.Background(), req)

				return err
			},
		},
		{
			name: "update с нулевым id",
			call: func(g *ColumnCatGuard) error {
				req := validUpdateColumnCatByIdRequest()
				req.Id = 0
				_, err := g.UpdateColumnCatById(context.Background(), req)

				return err
			},
		},
		{
			name: "create с nil-запросом",
			call: func(g *ColumnCatGuard) error {
				_, err := g.CreateColumnCat(context.Background(), nil)
				return err
			},
		},
		{
			name: "update с nil-запросом",
			call: func(g *ColumnCatGuard) error {
				_, err := g.UpdateColumnCatById(context.Background(), nil)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aliases := &stubAliases{}

			err := tt.call(guardWith(aliases))

			assertCode(t, err, codes.InvalidArgument)
			if aliases.calls != 0 {
				t.Errorf("выборок алиаса = %d, ожидалось 0", aliases.calls)
			}
		})
	}
}
