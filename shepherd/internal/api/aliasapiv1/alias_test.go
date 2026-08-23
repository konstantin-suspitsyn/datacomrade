package aliasapiv1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/pagination"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// subject — sub из claims токена. Хендлеры обязаны подставлять его в
// external_id запроса к Metadata Service, что бы ни пришло в теле.
const subject = "5f0a0e1e-0000-4000-8000-000000000001"

// createdAt/updatedAt — фиксированные метки, чтобы сравнивать JSON дословно.
var (
	createdAt = time.Date(2026, 8, 18, 10, 0, 0, 0, time.UTC)
	updatedAt = time.Date(2026, 8, 18, 11, 30, 0, 0, time.UTC)
)

// stubClient подменяет gRPC-клиент Metadata Service: запоминает запрос, с
// которым его позвали, считает вызовы и отдаёт заданный ответ либо ошибку.
// Счётчик нужен для проверок «до бэкенда не дошли»: 400/401 обязаны
// отрабатывать до похода в gRPC.
type stubClient struct {
	calls int
	err   error

	getReq  *tablesv1.GetAliasesRequest
	getResp *tablesv1.GetAliasesResponse

	getDeletedReq  *tablesv1.GetAliasesDeletedRequest
	getDeletedResp *tablesv1.GetAliasesDeletedResponse

	getByIdReq  *tablesv1.GetAliasByIdRequest
	getByIdResp *tablesv1.GetAliasByIdResponse

	createReq  *tablesv1.CreateAliasRequest
	createResp *tablesv1.CreateAliasResponse

	updateReq   *tablesv1.UpdateAliasByIdRequest
	deleteReq   *tablesv1.DeleteAliasByIdRequest
	undeleteReq *tablesv1.UndeleteAliasByIdRequest
}

func (s *stubClient) GetAliases(_ context.Context, in *tablesv1.GetAliasesRequest, _ ...grpc.CallOption) (*tablesv1.GetAliasesResponse, error) {
	s.calls++
	s.getReq = in

	if s.err != nil {
		return nil, s.err
	}

	return s.getResp, nil
}

func (s *stubClient) GetAliasesDeleted(_ context.Context, in *tablesv1.GetAliasesDeletedRequest, _ ...grpc.CallOption) (*tablesv1.GetAliasesDeletedResponse, error) {
	s.calls++
	s.getDeletedReq = in

	if s.err != nil {
		return nil, s.err
	}

	return s.getDeletedResp, nil
}

func (s *stubClient) GetAliasById(_ context.Context, in *tablesv1.GetAliasByIdRequest, _ ...grpc.CallOption) (*tablesv1.GetAliasByIdResponse, error) {
	s.calls++
	s.getByIdReq = in

	if s.err != nil {
		return nil, s.err
	}

	return s.getByIdResp, nil
}

func (s *stubClient) CreateAlias(_ context.Context, in *tablesv1.CreateAliasRequest, _ ...grpc.CallOption) (*tablesv1.CreateAliasResponse, error) {
	s.calls++
	s.createReq = in

	if s.err != nil {
		return nil, s.err
	}

	return s.createResp, nil
}

func (s *stubClient) UpdateAliasById(_ context.Context, in *tablesv1.UpdateAliasByIdRequest, _ ...grpc.CallOption) (*tablesv1.UpdateAliasByIdResponse, error) {
	s.calls++
	s.updateReq = in

	if s.err != nil {
		return nil, s.err
	}

	return &tablesv1.UpdateAliasByIdResponse{}, nil
}

func (s *stubClient) DeleteAliasById(_ context.Context, in *tablesv1.DeleteAliasByIdRequest, _ ...grpc.CallOption) (*tablesv1.DeleteAliasByIdResponse, error) {
	s.calls++
	s.deleteReq = in

	if s.err != nil {
		return nil, s.err
	}

	return &tablesv1.DeleteAliasByIdResponse{}, nil
}

func (s *stubClient) UndeleteAliasById(_ context.Context, in *tablesv1.UndeleteAliasByIdRequest, _ ...grpc.CallOption) (*tablesv1.UndeleteAliasByIdResponse, error) {
	s.calls++
	s.undeleteReq = in

	if s.err != nil {
		return nil, s.err
	}

	return &tablesv1.UndeleteAliasByIdResponse{}, nil
}

// TestListAliasesPagination фиксирует контракт постраничности: непереданные
// order/page_limit подменяются умолчаниями shepherd, переданные уходят в gRPC
// как есть, а order вне ASC/DESC — 400, а не молчаливый дефолт.
func TestListAliasesPagination(t *testing.T) {
	order := "ASC"
	page := int32(3)
	pageLimit := int32(25)
	badOrder := "asc"

	tests := []struct {
		name          string
		params        openapiv1.ListAliasesParams
		wantStatus    int
		wantCalls     int
		wantOrder     string
		wantPage      int32
		wantPageLimit int32
	}{
		{
			name:          "умолчания вместо непереданных параметров",
			params:        openapiv1.ListAliasesParams{},
			wantStatus:    http.StatusOK,
			wantCalls:     1,
			wantOrder:     pagination.DefaultOrder,
			wantPage:      0,
			wantPageLimit: pagination.DefaultPageLimit,
		},
		{
			name:          "переданные параметры уходят в gRPC",
			params:        openapiv1.ListAliasesParams{Order: &order, Page: &page, PageLimit: &pageLimit},
			wantStatus:    http.StatusOK,
			wantCalls:     1,
			wantOrder:     order,
			wantPage:      page,
			wantPageLimit: pageLimit,
		},
		{
			name:       "order вне ASC/DESC отсекается до gRPC",
			params:     openapiv1.ListAliasesParams{Order: &badOrder},
			wantStatus: http.StatusBadRequest,
			wantCalls:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubClient{getResp: &tablesv1.GetAliasesResponse{}}
			rec := httptest.NewRecorder()

			New(client).ListAliases(rec, newRequest(t, http.MethodGet, "", nil), tt.params)

			assertStatus(t, rec, tt.wantStatus)
			if client.calls != tt.wantCalls {
				t.Fatalf("вызовов gRPC = %d, ожидалось %d", client.calls, tt.wantCalls)
			}
			if tt.wantCalls == 0 {
				return
			}

			if got := client.getReq.GetOrder(); got != tt.wantOrder {
				t.Errorf("order = %q, ожидался %q", got, tt.wantOrder)
			}
			if got := client.getReq.GetPage(); got != tt.wantPage {
				t.Errorf("page = %d, ожидалась %d", got, tt.wantPage)
			}
			if got := client.getReq.GetPageLimit(); got != tt.wantPageLimit {
				t.Errorf("page_limit = %d, ожидался %d", got, tt.wantPageLimit)
			}
		})
	}
}

// TestListDeletedAliasesPagination — у выборки удалённых свой тип запроса, и
// умолчания ей нужно подставлять отдельно: общего кода с ListAliases нет.
func TestListDeletedAliasesPagination(t *testing.T) {
	order := "ASC"
	page := int32(2)
	pageLimit := int32(10)
	badOrder := "DESCENDING"

	tests := []struct {
		name          string
		params        openapiv1.ListDeletedAliasesParams
		wantStatus    int
		wantCalls     int
		wantOrder     string
		wantPage      int32
		wantPageLimit int32
	}{
		{
			name:          "умолчания вместо непереданных параметров",
			params:        openapiv1.ListDeletedAliasesParams{},
			wantStatus:    http.StatusOK,
			wantCalls:     1,
			wantOrder:     pagination.DefaultOrder,
			wantPage:      0,
			wantPageLimit: pagination.DefaultPageLimit,
		},
		{
			name:          "переданные параметры уходят в gRPC",
			params:        openapiv1.ListDeletedAliasesParams{Order: &order, Page: &page, PageLimit: &pageLimit},
			wantStatus:    http.StatusOK,
			wantCalls:     1,
			wantOrder:     order,
			wantPage:      page,
			wantPageLimit: pageLimit,
		},
		{
			name:       "order вне ASC/DESC отсекается до gRPC",
			params:     openapiv1.ListDeletedAliasesParams{Order: &badOrder},
			wantStatus: http.StatusBadRequest,
			wantCalls:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubClient{getDeletedResp: &tablesv1.GetAliasesDeletedResponse{}}
			rec := httptest.NewRecorder()

			New(client).ListDeletedAliases(rec, newRequest(t, http.MethodGet, "", nil), tt.params)

			assertStatus(t, rec, tt.wantStatus)
			if client.calls != tt.wantCalls {
				t.Fatalf("вызовов gRPC = %d, ожидалось %d", client.calls, tt.wantCalls)
			}
			if tt.wantCalls == 0 {
				return
			}

			if got := client.getDeletedReq.GetOrder(); got != tt.wantOrder {
				t.Errorf("order = %q, ожидался %q", got, tt.wantOrder)
			}
			if got := client.getDeletedReq.GetPage(); got != tt.wantPage {
				t.Errorf("page = %d, ожидалась %d", got, tt.wantPage)
			}
			if got := client.getDeletedReq.GetPageLimit(); got != tt.wantPageLimit {
				t.Errorf("page_limit = %d, ожидался %d", got, tt.wantPageLimit)
			}
		})
	}
}

// TestListAliasesBody проверяет форму ответа: имена полей — то, на что смотрит
// фронтенд, и переименование в конвертере он молча не переживёт. Отдельно
// важен column_type_id: он появился позже остальных колонок, и потерять его в
// toAlias компилятор не помешает — поле просто уедет нулём.
func TestListAliasesBody(t *testing.T) {
	client := &stubClient{getResp: &tablesv1.GetAliasesResponse{
		Data: []*tablesv1.Alias{{
			Id:           7,
			Name:         "выручка",
			Description:  "сумма продаж без НДС",
			CreatedAt:    timestamppb.New(createdAt),
			UpdatedAt:    timestamppb.New(updatedAt),
			IsDeleted:    false,
			UserId:       42,
			ColumnTypeId: 3,
		}},
		Pagination: &tablesv1.Pagination{Page: 1, PerPage: 100, TotalItems: 1, TotalPages: 1},
	}}
	rec := httptest.NewRecorder()

	New(client).ListAliases(rec, newRequest(t, http.MethodGet, "", nil), openapiv1.ListAliasesParams{})

	assertStatus(t, rec, http.StatusOK)
	assertContentTypeJSON(t, rec)

	var got openapiv1.AliasPage
	decode(t, rec, &got)

	if len(got.Data) != 1 {
		t.Fatalf("строк в data = %d, ожидалась 1", len(got.Data))
	}

	want := openapiv1.Alias{
		Id:           7,
		Name:         "выручка",
		Description:  "сумма продаж без НДС",
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		IsDeleted:    false,
		UserId:       42,
		ColumnTypeId: 3,
	}
	assertAlias(t, got.Data[0], want)

	wantPagination := openapiv1.Pagination{Page: 1, PerPage: 100, TotalItems: 1, TotalPages: 1}
	if got.Pagination != wantPagination {
		t.Errorf("pagination = %+v, ожидалась %+v", got.Pagination, wantPagination)
	}
}

// TestListDeletedAliasesBody — выборка удалённых идёт через тот же конвертер,
// но свой ответ: проверяем, что страница собирается и is_deleted доезжает true.
func TestListDeletedAliasesBody(t *testing.T) {
	client := &stubClient{getDeletedResp: &tablesv1.GetAliasesDeletedResponse{
		Data: []*tablesv1.Alias{{
			Id:           8,
			Name:         "маржа",
			Description:  "устаревший алиас",
			CreatedAt:    timestamppb.New(createdAt),
			UpdatedAt:    timestamppb.New(updatedAt),
			IsDeleted:    true,
			UserId:       42,
			ColumnTypeId: 1,
		}},
		Pagination: &tablesv1.Pagination{Page: 1, PerPage: 100, TotalItems: 1, TotalPages: 1},
	}}
	rec := httptest.NewRecorder()

	New(client).ListDeletedAliases(rec, newRequest(t, http.MethodGet, "", nil), openapiv1.ListDeletedAliasesParams{})

	assertStatus(t, rec, http.StatusOK)
	assertContentTypeJSON(t, rec)

	var got openapiv1.AliasPage
	decode(t, rec, &got)

	if len(got.Data) != 1 {
		t.Fatalf("строк в data = %d, ожидалась 1", len(got.Data))
	}

	assertAlias(t, got.Data[0], openapiv1.Alias{
		Id:           8,
		Name:         "маржа",
		Description:  "устаревший алиас",
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		IsDeleted:    true,
		UserId:       42,
		ColumnTypeId: 1,
	})
}

// TestListAliasesGRPCError — ошибка бэкенда доезжает до клиента со своим
// кодом, а не схлопывается в 500.
func TestListAliasesGRPCError(t *testing.T) {
	client := &stubClient{err: status.Error(codes.Unavailable, "datacatalogue недоступен")}
	rec := httptest.NewRecorder()

	New(client).ListAliases(rec, newRequest(t, http.MethodGet, "", nil), openapiv1.ListAliasesParams{})

	assertStatus(t, rec, http.StatusServiceUnavailable)
	if got := decodeError(t, rec).Code; got != "unavailable" {
		t.Errorf("code = %q, ожидался %q", got, "unavailable")
	}
}

// TestGetAliasById — чтение по id: id уходит в gRPC, строка возвращается
// целиком. Токен этой ручке не нужен — она ничего не пишет.
func TestGetAliasById(t *testing.T) {
	client := &stubClient{getByIdResp: &tablesv1.GetAliasByIdResponse{
		Alias: &tablesv1.Alias{
			Id:           7,
			Name:         "выручка",
			Description:  "сумма продаж без НДС",
			CreatedAt:    timestamppb.New(createdAt),
			UpdatedAt:    timestamppb.New(updatedAt),
			IsDeleted:    false,
			UserId:       42,
			ColumnTypeId: 3,
		},
	}}
	rec := httptest.NewRecorder()

	New(client).GetAliasById(rec, newRequest(t, http.MethodGet, "", nil), 7)

	assertStatus(t, rec, http.StatusOK)
	assertContentTypeJSON(t, rec)

	if got := client.getByIdReq.GetId(); got != 7 {
		t.Errorf("id = %d, ожидался 7", got)
	}

	var got openapiv1.Alias
	decode(t, rec, &got)

	assertAlias(t, got, openapiv1.Alias{
		Id:           7,
		Name:         "выручка",
		Description:  "сумма продаж без НДС",
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		IsDeleted:    false,
		UserId:       42,
		ColumnTypeId: 3,
	})
}

// TestGetAliasByIdNotFound — несуществующий id обязан отдавать 404 с телом
// ошибки, а не пустую строку со 200.
func TestGetAliasByIdNotFound(t *testing.T) {
	client := &stubClient{err: status.Error(codes.NotFound, "dc.alias id = 7")}
	rec := httptest.NewRecorder()

	New(client).GetAliasById(rec, newRequest(t, http.MethodGet, "", nil), 7)

	assertStatus(t, rec, http.StatusNotFound)
	if got := decodeError(t, rec).Code; got != "not_found" {
		t.Errorf("code = %q, ожидался %q", got, "not_found")
	}
}

// TestCreateAliasUsesClaims — ключевое свойство всех пишущих ручек: автор
// изменения берётся из токена. Тело намеренно содержит external_id — его в
// схеме нет, и подсунуть чужой id через него нельзя. Заодно проверяется, что
// column_type_id из тела доезжает до gRPC: без него бэкенд ответит 400 на
// Int64ID, и потеря поля выглядела бы как ошибка валидации, а не маппинга.
func TestCreateAliasUsesClaims(t *testing.T) {
	client := &stubClient{createResp: &tablesv1.CreateAliasResponse{
		Alias: &tablesv1.Alias{
			Id:           7,
			Name:         "выручка",
			Description:  "сумма продаж без НДС",
			CreatedAt:    timestamppb.New(createdAt),
			UpdatedAt:    timestamppb.New(updatedAt),
			UserId:       42,
			ColumnTypeId: 3,
		},
	}}
	rec := httptest.NewRecorder()

	body := `{"name":"выручка","description":"сумма продаж без НДС","column_type_id":3,"external_id":"подделка"}`
	New(client).CreateAlias(rec, newRequest(t, http.MethodPost, body, claims()))

	assertStatus(t, rec, http.StatusCreated)

	if got := client.createReq.GetExternalId(); got != subject {
		t.Errorf("external_id = %q, ожидался sub из claims %q", got, subject)
	}
	if got := client.createReq.GetName(); got != "выручка" {
		t.Errorf("name = %q, ожидалось %q", got, "выручка")
	}
	if got := client.createReq.GetDescription(); got != "сумма продаж без НДС" {
		t.Errorf("description = %q, ожидалось %q", got, "сумма продаж без НДС")
	}
	if got := client.createReq.GetColumnTypeId(); got != 3 {
		t.Errorf("column_type_id = %d, ожидался 3", got)
	}

	var created openapiv1.Alias
	decode(t, rec, &created)
	if created.Id != 7 || created.ColumnTypeId != 3 {
		t.Errorf("тело ответа = %+v, ожидалась созданная строка id=7 с column_type_id=3", created)
	}
}

// TestCreateAliasRejects — запрос без claims или с битым телом заворачивается
// до похода в Metadata Service.
func TestCreateAliasRejects(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		claims     *auth.Claims
		wantStatus int
		wantCode   string
	}{
		{
			name:       "без claims",
			body:       `{"name":"выручка","description":"описание","column_type_id":3}`,
			claims:     nil,
			wantStatus: http.StatusUnauthorized,
			wantCode:   "unauthenticated",
		},
		{
			name:       "битый JSON",
			body:       `{"name":`,
			claims:     claims(),
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_argument",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubClient{}
			rec := httptest.NewRecorder()

			New(client).CreateAlias(rec, newRequest(t, http.MethodPost, tt.body, tt.claims))

			assertStatus(t, rec, tt.wantStatus)
			if client.calls != 0 {
				t.Errorf("вызовов gRPC = %d, ожидалось 0", client.calls)
			}
			if got := decodeError(t, rec).Code; got != tt.wantCode {
				t.Errorf("code = %q, ожидался %q", got, tt.wantCode)
			}
		})
	}
}

// TestCreateAliasInvalidArgument — column_type_id проверяет бэкенд
// (validation.Int64ID), а не gateway: его InvalidArgument обязан стать 400 с
// исходным сообщением, а не 500.
func TestCreateAliasInvalidArgument(t *testing.T) {
	client := &stubClient{err: status.Error(codes.InvalidArgument, "column_type_id: должно быть больше нуля")}
	rec := httptest.NewRecorder()

	body := `{"name":"выручка","description":"описание","column_type_id":0}`
	New(client).CreateAlias(rec, newRequest(t, http.MethodPost, body, claims()))

	assertStatus(t, rec, http.StatusBadRequest)

	got := decodeError(t, rec)
	if got.Code != "invalid_argument" {
		t.Errorf("code = %q, ожидался %q", got.Code, "invalid_argument")
	}
	if !strings.Contains(got.Message, "column_type_id") {
		t.Errorf("message = %q, ожидалось упоминание column_type_id", got.Message)
	}
}

// TestUpdateAliasById — id берётся из пути, автор из claims, остальные поля из
// тела; ответ 204 без тела.
func TestUpdateAliasById(t *testing.T) {
	client := &stubClient{}
	rec := httptest.NewRecorder()

	body := `{"name":"выручка нетто","description":"уточнили формулу","is_deleted":false,"column_type_id":5}`
	New(client).UpdateAliasById(rec, newRequest(t, http.MethodPut, body, claims()), 7)

	assertStatus(t, rec, http.StatusNoContent)
	assertEmptyBody(t, rec)

	if got := client.updateReq.GetId(); got != 7 {
		t.Errorf("id = %d, ожидался 7", got)
	}
	if got := client.updateReq.GetName(); got != "выручка нетто" {
		t.Errorf("name = %q, ожидалось %q", got, "выручка нетто")
	}
	if got := client.updateReq.GetDescription(); got != "уточнили формулу" {
		t.Errorf("description = %q, ожидалось %q", got, "уточнили формулу")
	}
	if got := client.updateReq.GetIsDeleted(); got {
		t.Errorf("is_deleted = %v, ожидалось false", got)
	}
	if got := client.updateReq.GetColumnTypeId(); got != 5 {
		t.Errorf("column_type_id = %d, ожидался 5", got)
	}
	if got := client.updateReq.GetExternalId(); got != subject {
		t.Errorf("external_id = %q, ожидался sub из claims %q", got, subject)
	}
}

// TestUpdateAliasByIdRejects — та же пара отказов, что и у Create.
func TestUpdateAliasByIdRejects(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		claims     *auth.Claims
		wantStatus int
	}{
		{"без claims", `{"name":"выручка","description":"описание","is_deleted":false,"column_type_id":3}`, nil, http.StatusUnauthorized},
		{"битый JSON", `{`, claims(), http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubClient{}
			rec := httptest.NewRecorder()

			New(client).UpdateAliasById(rec, newRequest(t, http.MethodPut, tt.body, tt.claims), 7)

			assertStatus(t, rec, tt.wantStatus)
			if client.calls != 0 {
				t.Errorf("вызовов gRPC = %d, ожидалось 0", client.calls)
			}
		})
	}
}

// TestUpdateAliasByIdNotFound — обновление несуществующей строки обязано
// отдавать 404, а не 204.
func TestUpdateAliasByIdNotFound(t *testing.T) {
	client := &stubClient{err: status.Error(codes.NotFound, "dc.alias id = 7")}
	rec := httptest.NewRecorder()

	body := `{"name":"выручка","description":"описание","is_deleted":false,"column_type_id":3}`
	New(client).UpdateAliasById(rec, newRequest(t, http.MethodPut, body, claims()), 7)

	assertStatus(t, rec, http.StatusNotFound)
}

// TestDeleteUndeleteAliasById — обе ручки принимают только id и автора: тела
// запроса у них нет, ответ 204.
func TestDeleteUndeleteAliasById(t *testing.T) {
	t.Run("delete", func(t *testing.T) {
		client := &stubClient{}
		rec := httptest.NewRecorder()

		New(client).DeleteAliasById(rec, newRequest(t, http.MethodDelete, "", claims()), 7)

		assertStatus(t, rec, http.StatusNoContent)
		assertEmptyBody(t, rec)
		if got := client.deleteReq.GetId(); got != 7 {
			t.Errorf("id = %d, ожидался 7", got)
		}
		if got := client.deleteReq.GetExternalId(); got != subject {
			t.Errorf("external_id = %q, ожидался sub из claims %q", got, subject)
		}
	})

	t.Run("undelete", func(t *testing.T) {
		client := &stubClient{}
		rec := httptest.NewRecorder()

		New(client).UndeleteAliasById(rec, newRequest(t, http.MethodPost, "", claims()), 7)

		assertStatus(t, rec, http.StatusNoContent)
		assertEmptyBody(t, rec)
		if got := client.undeleteReq.GetId(); got != 7 {
			t.Errorf("id = %d, ожидался 7", got)
		}
		if got := client.undeleteReq.GetExternalId(); got != subject {
			t.Errorf("external_id = %q, ожидался sub из claims %q", got, subject)
		}
	})
}

// TestDeleteUndeleteAliasByIdWithoutClaims — без токена ни одна из них до
// бэкенда не доходит.
func TestDeleteUndeleteAliasByIdWithoutClaims(t *testing.T) {
	tests := []struct {
		name string
		call func(h *Handler, w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "delete",
			call: func(h *Handler, w http.ResponseWriter, r *http.Request) { h.DeleteAliasById(w, r, 7) },
		},
		{
			name: "undelete",
			call: func(h *Handler, w http.ResponseWriter, r *http.Request) { h.UndeleteAliasById(w, r, 7) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubClient{}
			rec := httptest.NewRecorder()

			tt.call(New(client), rec, newRequest(t, http.MethodPost, "", nil))

			assertStatus(t, rec, http.StatusUnauthorized)
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

	req := httptest.NewRequest(method, "/v1/aliases", strings.NewReader(body))
	if c != nil {
		req = req.WithContext(auth.ContextWithClaims(req.Context(), *c))
	}

	return req
}

// assertAlias сравнивает строку целиком. Метки времени сверяются через Equal —
// после JSON-обхода у них другой location, и оператор != на этом сорвался бы.
func assertAlias(t *testing.T, got, want openapiv1.Alias) {
	t.Helper()

	if !got.CreatedAt.Equal(want.CreatedAt) || !got.UpdatedAt.Equal(want.UpdatedAt) {
		t.Errorf("метки времени = %v/%v, ожидались %v/%v", got.CreatedAt, got.UpdatedAt, want.CreatedAt, want.UpdatedAt)
	}

	got.CreatedAt, got.UpdatedAt = want.CreatedAt, want.UpdatedAt
	if got != want {
		t.Errorf("строка = %+v, ожидалась %+v", got, want)
	}
}

func assertStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()

	if rec.Code != want {
		t.Fatalf("статус = %d, ожидался %d (тело = %s)", rec.Code, want, rec.Body.String())
	}
}

func assertContentTypeJSON(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()

	const want = "application/json; charset=utf-8"
	if got := rec.Header().Get("Content-Type"); got != want {
		t.Errorf("Content-Type = %q, ожидался %q", got, want)
	}
}

func assertEmptyBody(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()

	if rec.Body.Len() != 0 {
		t.Errorf("тело = %s, ожидалось пустое", rec.Body.String())
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
