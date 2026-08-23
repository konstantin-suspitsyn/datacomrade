package columncatapiv1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// subject — sub из claims токена: автор изменения берётся оттуда, а не из тела.
const subject = "5f0a0e1e-0000-4000-8000-000000000001"

// Весь файл про одно поле — alias_id. Колонка каталога может быть без алиаса,
// и в контракте это указатель: null и отсутствие поля означают «алиаса нет»,
// а 0 — не значение, а потерянный указатель. Компилятор такую потерю не
// поймает: и в REST-схеме, и в proto тип один и тот же *int64, и подстановка
// геттера вместо поля собирается молча.

// stubClient подменяет gRPC-клиент Metadata Service: запоминает запрос, с
// которым его позвали, и отдаёт заданный ответ либо ошибку.
type stubClient struct {
	calls int
	err   error

	getByIdReq  *tablesv1.GetColumnCatByIdRequest
	getByIdResp *tablesv1.GetColumnCatByIdResponse

	createReq  *tablesv1.CreateColumnCatRequest
	createResp *tablesv1.CreateColumnCatResponse

	updateReq *tablesv1.UpdateColumnCatByIdRequest
}

func (s *stubClient) GetColumnCatById(_ context.Context, in *tablesv1.GetColumnCatByIdRequest, _ ...grpc.CallOption) (*tablesv1.GetColumnCatByIdResponse, error) {
	s.calls++
	s.getByIdReq = in

	if s.err != nil {
		return nil, s.err
	}

	return s.getByIdResp, nil
}

func (s *stubClient) CreateColumnCat(_ context.Context, in *tablesv1.CreateColumnCatRequest, _ ...grpc.CallOption) (*tablesv1.CreateColumnCatResponse, error) {
	s.calls++
	s.createReq = in

	if s.err != nil {
		return nil, s.err
	}

	return s.createResp, nil
}

func (s *stubClient) UpdateColumnCatById(_ context.Context, in *tablesv1.UpdateColumnCatByIdRequest, _ ...grpc.CallOption) (*tablesv1.UpdateColumnCatByIdResponse, error) {
	s.calls++
	s.updateReq = in

	if s.err != nil {
		return nil, s.err
	}

	return &tablesv1.UpdateColumnCatByIdResponse{}, nil
}

// Остальные методы клиента в этих проверках не участвуют, но нужны, чтобы
// заглушка удовлетворяла интерфейсу целиком.
func (s *stubClient) GetColumnCats(context.Context, *tablesv1.GetColumnCatsRequest, ...grpc.CallOption) (*tablesv1.GetColumnCatsResponse, error) {
	return &tablesv1.GetColumnCatsResponse{}, nil
}

func (s *stubClient) GetColumnCatsDeleted(context.Context, *tablesv1.GetColumnCatsDeletedRequest, ...grpc.CallOption) (*tablesv1.GetColumnCatsDeletedResponse, error) {
	return &tablesv1.GetColumnCatsDeletedResponse{}, nil
}

func (s *stubClient) GetColumnCatsOptionalFilters(context.Context, *tablesv1.GetColumnCatsOptionalFiltersRequest, ...grpc.CallOption) (*tablesv1.GetColumnCatsOptionalFiltersResponse, error) {
	return &tablesv1.GetColumnCatsOptionalFiltersResponse{}, nil
}

func (s *stubClient) GetColumnCatsSearchName(context.Context, *tablesv1.GetColumnCatsSearchNameRequest, ...grpc.CallOption) (*tablesv1.GetColumnCatsSearchNameResponse, error) {
	return &tablesv1.GetColumnCatsSearchNameResponse{}, nil
}

func (s *stubClient) DeleteColumnCatById(context.Context, *tablesv1.DeleteColumnCatByIdRequest, ...grpc.CallOption) (*tablesv1.DeleteColumnCatByIdResponse, error) {
	return &tablesv1.DeleteColumnCatByIdResponse{}, nil
}

func (s *stubClient) UndeleteColumnCatById(context.Context, *tablesv1.UndeleteColumnCatByIdRequest, ...grpc.CallOption) (*tablesv1.UndeleteColumnCatByIdResponse, error) {
	return &tablesv1.UndeleteColumnCatByIdResponse{}, nil
}

// TestGetColumnCatByIdAliasIdOptional — форма ответа. Колонка с алиасом отдаёт
// число, колонка без алиаса не отдаёт поле вовсе: по схеме alias_id
// необязателен и nullable, и подстановка нуля выглядела бы как ссылка на
// алиас с id 0.
func TestGetColumnCatByIdAliasIdOptional(t *testing.T) {
	aliasID := int64(5)

	tests := []struct {
		name     string
		aliasID  *int64
		wantJSON string
		wantPtr  *int64
	}{
		{name: "алиас задан", aliasID: &aliasID, wantJSON: `"alias_id":5`, wantPtr: &aliasID},
		{name: "алиаса нет", aliasID: nil, wantJSON: "", wantPtr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubClient{getByIdResp: &tablesv1.GetColumnCatByIdResponse{
				ColumnCat: &tablesv1.ColumnCat{Id: 7, TableId: 3, Name: "amount", AliasId: tt.aliasID, ColumnTypeId: 2},
			}}
			rec := httptest.NewRecorder()

			New(client).GetColumnCatById(rec, newRequest(t, http.MethodGet, "", nil), 7)

			assertStatus(t, rec, http.StatusOK)

			if got := client.getByIdReq.GetId(); got != 7 {
				t.Errorf("id = %d, ожидался 7", got)
			}

			// Сырое тело, а не только разобранная структура: отсутствие поля и
			// его нулевое значение после Unmarshal неразличимы не будут, но
			// фронтенд читает именно JSON.
			body := rec.Body.String()
			if tt.wantJSON != "" && !strings.Contains(body, tt.wantJSON) {
				t.Errorf("тело = %s, ожидалось вхождение %s", body, tt.wantJSON)
			}
			if tt.wantJSON == "" && strings.Contains(body, "alias_id") {
				t.Errorf("тело = %s, ожидалось отсутствие alias_id", body)
			}

			var got openapiv1.ColumnCat
			decode(t, rec, &got)

			switch {
			case tt.wantPtr == nil && got.AliasId != nil:
				t.Errorf("alias_id = %d, ожидался nil", *got.AliasId)
			case tt.wantPtr != nil && got.AliasId == nil:
				t.Errorf("alias_id = nil, ожидался %d", *tt.wantPtr)
			case tt.wantPtr != nil && *got.AliasId != *tt.wantPtr:
				t.Errorf("alias_id = %d, ожидался %d", *got.AliasId, *tt.wantPtr)
			}
		})
	}
}

// TestCreateColumnCatAliasId — форма запроса. Три записи клиента означают одно
// и то же для двух из них: null и пропуск поля — «алиаса нет». Ноль остаётся
// значением и уезжает в бэкенд как есть: отбивать его — дело валидации
// Metadata Service, а не гейтвея.
func TestCreateColumnCatAliasId(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantSet bool
		want    int64
	}{
		{
			name:    "алиас передан",
			body:    `{"table_id":3,"name":"amount","alias_id":5,"column_type_id":2,"description":"сумма","calculation_type_id":1,"show_in_ui":true}`,
			wantSet: true,
			want:    5,
		},
		// show_in_ui = false у случаев без алиаса не декорация: с true такое
		// тело отбивает aliasSetForUI, и до проверки формы alias_id дело бы
		// не дошло — см. TestColumnCatAliasRequiredForUI.
		{
			name: "alias_id: null снимает алиас",
			body: `{"table_id":3,"name":"amount","alias_id":null,"column_type_id":2,"description":"сумма","calculation_type_id":1,"show_in_ui":false}`,
		},
		{
			name: "поле не передано",
			body: `{"table_id":3,"name":"amount","column_type_id":2,"description":"сумма","calculation_type_id":1,"show_in_ui":false}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubClient{createResp: &tablesv1.CreateColumnCatResponse{
				ColumnCat: &tablesv1.ColumnCat{Id: 7},
			}}
			rec := httptest.NewRecorder()

			New(client).CreateColumnCat(rec, newRequest(t, http.MethodPost, tt.body, claims()))

			assertStatus(t, rec, http.StatusCreated)

			got := client.createReq.AliasId
			switch {
			case !tt.wantSet && got != nil:
				t.Errorf("alias_id = %d, ожидался nil", *got)
			case tt.wantSet && got == nil:
				t.Errorf("alias_id = nil, ожидался %d", tt.want)
			case tt.wantSet && *got != tt.want:
				t.Errorf("alias_id = %d, ожидался %d", *got, tt.want)
			}

			if got := client.createReq.GetExternalId(); got != subject {
				t.Errorf("external_id = %q, ожидался sub из claims %q", got, subject)
			}
		})
	}
}

// TestUpdateColumnCatByIdClearsAlias — снятие алиаса у существующей колонки:
// именно ради этого случая alias_id в схеме nullable, а не просто
// необязательный.
func TestUpdateColumnCatByIdClearsAlias(t *testing.T) {
	client := &stubClient{}
	rec := httptest.NewRecorder()

	// Вместе с алиасом снимается и показ в интерфейсе: колонку без алиаса
	// показывать нечем, и с show_in_ui = true такое тело отбил бы гейтвей.
	body := `{"table_id":3,"name":"amount","alias_id":null,"column_type_id":2,"description":"сумма","calculation_type_id":1,"show_in_ui":false}`
	New(client).UpdateColumnCatById(rec, newRequest(t, http.MethodPut, body, claims()), 7)

	assertStatus(t, rec, http.StatusNoContent)

	if got := client.updateReq.AliasId; got != nil {
		t.Errorf("alias_id = %d, ожидался nil", *got)
	}
	if got := client.updateReq.GetId(); got != 7 {
		t.Errorf("id = %d, ожидался 7", got)
	}
	if got := client.updateReq.GetExternalId(); got != subject {
		t.Errorf("external_id = %q, ожидался sub из claims %q", got, subject)
	}
}

// TestCreateColumnCatAliasRules — правило связи колонки с алиасом живёт в
// Metadata Service (ColumnCatGuard), гейтвей его только доносит. Обе ошибки
// описаны в openapi.yaml для этой ручки, и схлопывание любой из них в 500
// сделало бы их неотличимыми для клиента.
func TestCreateColumnCatAliasRules(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{
			name:       "column_type_id не совпал с алиасом",
			err:        status.Error(codes.InvalidArgument, "column_type_id: должен совпадать с column_type_id алиаса 5 (2)"),
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_argument",
		},
		{
			name:       "алиаса с таким id нет",
			err:        status.Error(codes.NotFound, "dc.alias id = 5"),
			wantStatus: http.StatusNotFound,
			wantCode:   "not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubClient{err: tt.err}
			rec := httptest.NewRecorder()

			body := `{"table_id":3,"name":"amount","alias_id":5,"column_type_id":9,"description":"сумма","calculation_type_id":1,"show_in_ui":true}`
			New(client).CreateColumnCat(rec, newRequest(t, http.MethodPost, body, claims()))

			assertStatus(t, rec, tt.wantStatus)
			if got := decodeError(t, rec).Code; got != tt.wantCode {
				t.Errorf("code = %q, ожидался %q", got, tt.wantCode)
			}
		})
	}
}

// TestColumnCatAliasRequiredForUI — правило гейтвея: show_in_ui = true без
// alias_id отбивается здесь, до бэкенда. Проверяются обе пишущие ручки и обе
// формы «алиаса нет» — null и пропуск поля: для JSON это разные тела, а для
// правила одно и то же.
func TestColumnCatAliasRequiredForUI(t *testing.T) {
	const (
		aliasNull    = `{"table_id":3,"name":"amount","alias_id":null,"column_type_id":2,"description":"сумма","calculation_type_id":1,"show_in_ui":true}`
		aliasMissing = `{"table_id":3,"name":"amount","column_type_id":2,"description":"сумма","calculation_type_id":1,"show_in_ui":true}`
	)

	tests := []struct {
		name string
		body string
		call func(t *testing.T, h *Handler, rec *httptest.ResponseRecorder, body string)
	}{
		{name: "create, alias_id: null", body: aliasNull, call: createColumnCat},
		{name: "create, поле не передано", body: aliasMissing, call: createColumnCat},
		{name: "update, alias_id: null", body: aliasNull, call: updateColumnCat},
		{name: "update, поле не передано", body: aliasMissing, call: updateColumnCat},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubClient{}
			rec := httptest.NewRecorder()

			tt.call(t, New(client), rec, tt.body)

			assertStatus(t, rec, http.StatusBadRequest)

			body := decodeError(t, rec)
			if body.Code != "invalid_argument" {
				t.Errorf("code = %q, ожидался %q", body.Code, "invalid_argument")
			}
			for _, want := range []string{"alias_id", "show_in_ui"} {
				if !strings.Contains(body.Message, want) {
					t.Errorf("message = %q, ожидалось упоминание %q", body.Message, want)
				}
			}

			if client.calls != 0 {
				t.Errorf("вызовов gRPC = %d, ожидалось 0", client.calls)
			}
		})
	}
}

// TestColumnCatAliasOptionalWithoutUI — обратная сторона правила: скрытая
// колонка без алиаса законна и уходит в бэкенд как есть, с alias_id = nil.
func TestColumnCatAliasOptionalWithoutUI(t *testing.T) {
	const body = `{"table_id":3,"name":"amount","column_type_id":2,"description":"сумма","calculation_type_id":1,"show_in_ui":false}`

	tests := []struct {
		name       string
		call       func(t *testing.T, h *Handler, rec *httptest.ResponseRecorder, body string)
		wantStatus int
		aliasID    func(c *stubClient) *int64
	}{
		{
			name:       "create",
			call:       createColumnCat,
			wantStatus: http.StatusCreated,
			aliasID:    func(c *stubClient) *int64 { return c.createReq.AliasId },
		},
		{
			name:       "update",
			call:       updateColumnCat,
			wantStatus: http.StatusNoContent,
			aliasID:    func(c *stubClient) *int64 { return c.updateReq.AliasId },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubClient{createResp: &tablesv1.CreateColumnCatResponse{
				ColumnCat: &tablesv1.ColumnCat{Id: 7},
			}}
			rec := httptest.NewRecorder()

			tt.call(t, New(client), rec, body)

			assertStatus(t, rec, tt.wantStatus)

			if client.calls != 1 {
				t.Fatalf("вызовов gRPC = %d, ожидался 1", client.calls)
			}
			if got := tt.aliasID(client); got != nil {
				t.Errorf("alias_id = %d, ожидался nil", *got)
			}
		})
	}
}

// createColumnCat и updateColumnCat сводят две ручки с разными сигнатурами к
// одной, чтобы таблицы выше не дублировали случаи по разу на ручку.
func createColumnCat(t *testing.T, h *Handler, rec *httptest.ResponseRecorder, body string) {
	t.Helper()

	h.CreateColumnCat(rec, newRequest(t, http.MethodPost, body, claims()))
}

func updateColumnCat(t *testing.T, h *Handler, rec *httptest.ResponseRecorder, body string) {
	t.Helper()

	h.UpdateColumnCatById(rec, newRequest(t, http.MethodPut, body, claims()), 7)
}

// TestCreateColumnCatRejects — без токена и с битым телом до бэкенда не идём.
func TestCreateColumnCatRejects(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		claims     *auth.Claims
		wantStatus int
	}{
		{name: "без claims", body: `{"table_id":3}`, claims: nil, wantStatus: http.StatusUnauthorized},
		{name: "битый JSON", body: `{"table_id":`, claims: claims(), wantStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubClient{}
			rec := httptest.NewRecorder()

			New(client).CreateColumnCat(rec, newRequest(t, http.MethodPost, tt.body, tt.claims))

			assertStatus(t, rec, tt.wantStatus)
			if client.calls != 0 {
				t.Errorf("вызовов gRPC = %d, ожидалось 0", client.calls)
			}
		})
	}
}

func claims() *auth.Claims {
	return &auth.Claims{Subject: subject, Roles: []string{"admin"}}
}

// newRequest собирает запрос с телом и, если claims переданы, с ними в
// контексте — так же, как это делает authMiddleware в проде.
func newRequest(t *testing.T, method, body string, c *auth.Claims) *http.Request {
	t.Helper()

	req := httptest.NewRequest(method, "/v1/column-cats", strings.NewReader(body))
	if c != nil {
		req = req.WithContext(auth.ContextWithClaims(req.Context(), *c))
	}

	return req
}

func assertStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()

	if rec.Code != want {
		t.Fatalf("статус = %d, ожидался %d (тело = %s)", rec.Code, want, rec.Body.String())
	}
}

func decode(t *testing.T, rec *httptest.ResponseRecorder, into any) {
	t.Helper()

	if err := json.Unmarshal(rec.Body.Bytes(), into); err != nil {
		t.Fatalf("тело не разобралось как JSON: %v (тело = %s)", err, rec.Body.String())
	}
}

func decodeError(t *testing.T, rec *httptest.ResponseRecorder) apierror.Body {
	t.Helper()

	var body apierror.Body
	decode(t, rec, &body)

	return body
}
