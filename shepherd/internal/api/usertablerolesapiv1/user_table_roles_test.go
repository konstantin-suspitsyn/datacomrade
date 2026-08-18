package usertablerolesapiv1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
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

// createdAt/updatedAt — фиксированные метки, чтобы сравнивать ответ дословно.
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

	getReq  *userdomainrolesv1.GetUserTableRolesesRequest
	getResp *userdomainrolesv1.GetUserTableRolesesResponse

	filterReq  *userdomainrolesv1.GetUserTableRolesesFilteredRequest
	filterResp *userdomainrolesv1.GetUserTableRolesesFilteredResponse

	getByIdReq  *userdomainrolesv1.GetUserTableRolesByIdRequest
	getByIdResp *userdomainrolesv1.GetUserTableRolesByIdResponse

	createReq  *userdomainrolesv1.CreateUserTableRolesRequest
	createResp *userdomainrolesv1.CreateUserTableRolesResponse

	updateReq   *userdomainrolesv1.UpdateUserTableRolesByIdRequest
	deleteReq   *userdomainrolesv1.DeleteUserTableRolesByIdRequest
	undeleteReq *userdomainrolesv1.UndeleteUserTableRolesByIdRequest
}

func (s *stubClient) GetUserTableRoleses(_ context.Context, in *userdomainrolesv1.GetUserTableRolesesRequest, _ ...grpc.CallOption) (*userdomainrolesv1.GetUserTableRolesesResponse, error) {
	s.calls++
	s.getReq = in

	if s.err != nil {
		return nil, s.err
	}

	return s.getResp, nil
}

func (s *stubClient) GetUserTableRolesesFiltered(_ context.Context, in *userdomainrolesv1.GetUserTableRolesesFilteredRequest, _ ...grpc.CallOption) (*userdomainrolesv1.GetUserTableRolesesFilteredResponse, error) {
	s.calls++
	s.filterReq = in

	if s.err != nil {
		return nil, s.err
	}

	return s.filterResp, nil
}

func (s *stubClient) GetUserTableRolesById(_ context.Context, in *userdomainrolesv1.GetUserTableRolesByIdRequest, _ ...grpc.CallOption) (*userdomainrolesv1.GetUserTableRolesByIdResponse, error) {
	s.calls++
	s.getByIdReq = in

	if s.err != nil {
		return nil, s.err
	}

	return s.getByIdResp, nil
}

func (s *stubClient) CreateUserTableRoles(_ context.Context, in *userdomainrolesv1.CreateUserTableRolesRequest, _ ...grpc.CallOption) (*userdomainrolesv1.CreateUserTableRolesResponse, error) {
	s.calls++
	s.createReq = in

	if s.err != nil {
		return nil, s.err
	}

	return s.createResp, nil
}

func (s *stubClient) UpdateUserTableRolesById(_ context.Context, in *userdomainrolesv1.UpdateUserTableRolesByIdRequest, _ ...grpc.CallOption) (*userdomainrolesv1.UpdateUserTableRolesByIdResponse, error) {
	s.calls++
	s.updateReq = in

	if s.err != nil {
		return nil, s.err
	}

	return &userdomainrolesv1.UpdateUserTableRolesByIdResponse{}, nil
}

func (s *stubClient) DeleteUserTableRolesById(_ context.Context, in *userdomainrolesv1.DeleteUserTableRolesByIdRequest, _ ...grpc.CallOption) (*userdomainrolesv1.DeleteUserTableRolesByIdResponse, error) {
	s.calls++
	s.deleteReq = in

	if s.err != nil {
		return nil, s.err
	}

	return &userdomainrolesv1.DeleteUserTableRolesByIdResponse{}, nil
}

func (s *stubClient) UndeleteUserTableRolesById(_ context.Context, in *userdomainrolesv1.UndeleteUserTableRolesByIdRequest, _ ...grpc.CallOption) (*userdomainrolesv1.UndeleteUserTableRolesByIdResponse, error) {
	s.calls++
	s.undeleteReq = in

	if s.err != nil {
		return nil, s.err
	}

	return &userdomainrolesv1.UndeleteUserTableRolesByIdResponse{}, nil
}

// TestListUserTableRolesPagination фиксирует контракт постраничности:
// непереданные order/page_limit подменяются умолчаниями shepherd, переданные
// уходят в gRPC как есть, а order вне ASC/DESC — 400, а не молчаливый дефолт.
func TestListUserTableRolesPagination(t *testing.T) {
	order := "ASC"
	page := int32(3)
	pageLimit := int32(25)
	badOrder := "по возрастанию"

	tests := []struct {
		name          string
		params        openapiv1.ListUserTableRolesParams
		wantStatus    int
		wantCalls     int
		wantOrder     string
		wantPage      int32
		wantPageLimit int32
	}{
		{
			name:          "умолчания вместо непереданных параметров",
			params:        openapiv1.ListUserTableRolesParams{},
			wantStatus:    http.StatusOK,
			wantCalls:     1,
			wantOrder:     pagination.DefaultOrder,
			wantPage:      0,
			wantPageLimit: pagination.DefaultPageLimit,
		},
		{
			name:          "переданные параметры уходят в gRPC",
			params:        openapiv1.ListUserTableRolesParams{Order: &order, Page: &page, PageLimit: &pageLimit},
			wantStatus:    http.StatusOK,
			wantCalls:     1,
			wantOrder:     order,
			wantPage:      page,
			wantPageLimit: pageLimit,
		},
		{
			name:       "order вне ASC/DESC отсекается до gRPC",
			params:     openapiv1.ListUserTableRolesParams{Order: &badOrder},
			wantStatus: http.StatusBadRequest,
			wantCalls:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubClient{getResp: &userdomainrolesv1.GetUserTableRolesesResponse{}}
			rec := httptest.NewRecorder()

			New(client).ListUserTableRoles(rec, newRequest(t, http.MethodGet, "", nil), tt.params)

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

// TestListUserTableRolesBody — форма строки списка. Главное здесь:
// table_roles_id, который в proto зовётся table_roles_id_2 (имя от sqlc для
// второй одноимённой колонки в джойне) — наружу он обязан уйти под обычным
// именем.
func TestListUserTableRolesBody(t *testing.T) {
	dcUserId := int64(42)
	tableCatId := int64(100)

	client := &stubClient{getResp: &userdomainrolesv1.GetUserTableRolesesResponse{
		Data: []*userdomainrolesv1.GetUserTableRolesesRow{{
			Id:                    7,
			UserId:                42,
			CreatedAt:             timestamppb.New(createdAt),
			UpdatedAt:             timestamppb.New(updatedAt),
			TableId:               100,
			DcUserId:              &dcUserId,
			DcUserName:            "Пётр",
			DcUserExternalId:      subject,
			TableRolesId_2:        3,
			TableRolesName:        "table_reader",
			TableRolesDescription: "чтение таблицы",
			TableCatId:            &tableCatId,
			TableCatName:          "orders",
		}},
		Pagination: &userdomainrolesv1.Pagination{Page: 1, PerPage: 100, TotalItems: 1, TotalPages: 1},
	}}
	rec := httptest.NewRecorder()

	New(client).ListUserTableRoles(rec, newRequest(t, http.MethodGet, "", nil), openapiv1.ListUserTableRolesParams{})

	assertStatus(t, rec, http.StatusOK)
	assertContentTypeJSON(t, rec)

	var got openapiv1.UserTableRolesPage
	decode(t, rec, &got)

	if len(got.Data) != 1 {
		t.Fatalf("строк в data = %d, ожидалась 1", len(got.Data))
	}

	row := got.Data[0]
	if row.TableRolesId != 3 {
		t.Errorf("table_roles_id = %d, ожидался 3 (переносится из table_roles_id_2)", row.TableRolesId)
	}
	if row.DcUserId == nil || *row.DcUserId != dcUserId {
		t.Errorf("dc_user_id = %v, ожидался %d", row.DcUserId, dcUserId)
	}
	if row.TableCatId == nil || *row.TableCatId != tableCatId {
		t.Errorf("table_cat_id = %v, ожидался %d", row.TableCatId, tableCatId)
	}
	if row.DcUserName != "Пётр" || row.DcUserExternalId != subject || row.TableCatName != "orders" {
		t.Errorf("джойн-поля = %q/%q/%q", row.DcUserName, row.DcUserExternalId, row.TableCatName)
	}
	if !row.CreatedAt.Equal(createdAt) || !row.UpdatedAt.Equal(updatedAt) {
		t.Errorf("метки времени = %v/%v, ожидались %v/%v", row.CreatedAt, row.UpdatedAt, createdAt, updatedAt)
	}

	wantPagination := openapiv1.Pagination{Page: 1, PerPage: 100, TotalItems: 1, TotalPages: 1}
	if got.Pagination != wantPagination {
		t.Errorf("pagination = %+v, ожидалась %+v", got.Pagination, wantPagination)
	}
}

// TestListUserTableRolesOmitsNullJoins — левый джойн может не найти строку:
// dc_user_id/table_cat_id тогда не должны появиться в JSON нулями, иначе
// фронтенд примет отсутствие связи за ссылку на id = 0.
func TestListUserTableRolesOmitsNullJoins(t *testing.T) {
	client := &stubClient{getResp: &userdomainrolesv1.GetUserTableRolesesResponse{
		Data: []*userdomainrolesv1.GetUserTableRolesesRow{{
			Id:        7,
			UserId:    42,
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
			TableId:   100,
		}},
		Pagination: &userdomainrolesv1.Pagination{},
	}}
	rec := httptest.NewRecorder()

	New(client).ListUserTableRoles(rec, newRequest(t, http.MethodGet, "", nil), openapiv1.ListUserTableRolesParams{})

	assertStatus(t, rec, http.StatusOK)

	var raw struct {
		Data []map[string]any `json:"data"`
	}
	decode(t, rec, &raw)

	if len(raw.Data) != 1 {
		t.Fatalf("строк в data = %d, ожидалась 1", len(raw.Data))
	}

	for _, key := range []string{"dc_user_id", "table_cat_id"} {
		if _, ok := raw.Data[0][key]; ok {
			t.Errorf("%s присутствует в JSON, хотя джойн ничего не нашёл: %v", key, raw.Data[0][key])
		}
	}
}

// TestFilterUserTableRoles — фильтры необязательны и независимы:
// непереданный уходит в gRPC как отсутствующий (nil) и не сужает выборку.
func TestFilterUserTableRoles(t *testing.T) {
	userName := "Пётр"
	tableCatId := int64(100)

	t.Run("непереданные фильтры уходят как nil", func(t *testing.T) {
		client := &stubClient{filterResp: &userdomainrolesv1.GetUserTableRolesesFilteredResponse{}}
		rec := httptest.NewRecorder()

		New(client).FilterUserTableRoles(rec, newRequest(t, http.MethodGet, "", nil), openapiv1.FilterUserTableRolesParams{})

		assertStatus(t, rec, http.StatusOK)

		req := client.filterReq
		if req.DcUserName != nil || req.DcUserExternalId != nil || req.TableRolesName != nil ||
			req.TableCatId != nil || req.TableCatName != nil {
			t.Errorf("непереданные фильтры доехали не как nil: %+v", req)
		}
		if got := req.GetOrder(); got != pagination.DefaultOrder {
			t.Errorf("order = %q, ожидался %q", got, pagination.DefaultOrder)
		}
	})

	t.Run("переданные фильтры уходят в gRPC", func(t *testing.T) {
		client := &stubClient{filterResp: &userdomainrolesv1.GetUserTableRolesesFilteredResponse{}}
		rec := httptest.NewRecorder()

		params := openapiv1.FilterUserTableRolesParams{DcUserName: &userName, TableCatId: &tableCatId}
		New(client).FilterUserTableRoles(rec, newRequest(t, http.MethodGet, "", nil), params)

		assertStatus(t, rec, http.StatusOK)

		if got := client.filterReq.GetDcUserName(); got != userName {
			t.Errorf("dc_user_name = %q, ожидался %q", got, userName)
		}
		if got := client.filterReq.GetTableCatId(); got != tableCatId {
			t.Errorf("table_cat_id = %d, ожидался %d", got, tableCatId)
		}
		// Остальные фильтры не передавались — сузить выборку они не должны.
		if client.filterReq.TableRolesName != nil || client.filterReq.DcUserExternalId != nil {
			t.Errorf("непереданные фильтры доехали не как nil: %+v", client.filterReq)
		}
	})

	t.Run("order вне ASC/DESC отсекается до gRPC", func(t *testing.T) {
		badOrder := "asc"
		client := &stubClient{}
		rec := httptest.NewRecorder()

		params := openapiv1.FilterUserTableRolesParams{Order: &badOrder}
		New(client).FilterUserTableRoles(rec, newRequest(t, http.MethodGet, "", nil), params)

		assertStatus(t, rec, http.StatusBadRequest)
		if client.calls != 0 {
			t.Errorf("вызовов gRPC = %d, ожидалось 0", client.calls)
		}
	})
}

// TestFilterUserTableRolesBody — у отфильтрованной строки нет table_roles_id:
// запрос с фильтрами его не отдаёт, поэтому и схема ответа отдельная.
func TestFilterUserTableRolesBody(t *testing.T) {
	client := &stubClient{filterResp: &userdomainrolesv1.GetUserTableRolesesFilteredResponse{
		Data: []*userdomainrolesv1.GetUserTableRolesesFilteredRow{{
			Id:                    7,
			UserId:                42,
			CreatedAt:             timestamppb.New(createdAt),
			UpdatedAt:             timestamppb.New(updatedAt),
			TableId:               100,
			DcUserName:            "Пётр",
			DcUserExternalId:      subject,
			TableRolesName:        "table_reader",
			TableRolesDescription: "чтение таблицы",
			TableCatName:          "orders",
		}},
		Pagination: &userdomainrolesv1.Pagination{Page: 2, PerPage: 50, TotalItems: 60, TotalPages: 2},
	}}
	rec := httptest.NewRecorder()

	New(client).FilterUserTableRoles(rec, newRequest(t, http.MethodGet, "", nil), openapiv1.FilterUserTableRolesParams{})

	assertStatus(t, rec, http.StatusOK)

	var got openapiv1.UserTableRolesFilteredPage
	decode(t, rec, &got)

	if len(got.Data) != 1 {
		t.Fatalf("строк в data = %d, ожидалась 1", len(got.Data))
	}

	row := got.Data[0]
	if row.Id != 7 || row.UserId != 42 || row.TableId != 100 {
		t.Errorf("ключевые поля строки = %d/%d/%d, ожидались 7/42/100", row.Id, row.UserId, row.TableId)
	}
	if row.TableRolesName != "table_reader" || row.TableCatName != "orders" {
		t.Errorf("джойн-поля = %q/%q", row.TableRolesName, row.TableCatName)
	}

	wantPagination := openapiv1.Pagination{Page: 2, PerPage: 50, TotalItems: 60, TotalPages: 2}
	if got.Pagination != wantPagination {
		t.Errorf("pagination = %+v, ожидалась %+v", got.Pagination, wantPagination)
	}
}

// TestGetUserTableRolesById — чтение по id: id уходит в gRPC, строка
// разворачивается в UserTableRolesDetail.
func TestGetUserTableRolesById(t *testing.T) {
	dcUserId := int64(42)

	client := &stubClient{getByIdResp: &userdomainrolesv1.GetUserTableRolesByIdResponse{
		UserTableRoles: &userdomainrolesv1.GetUserTableRolesByIdRow{
			Id:                    7,
			UserId:                42,
			TableRolesId:          3,
			CreatedAt:             timestamppb.New(createdAt),
			UpdatedAt:             timestamppb.New(updatedAt),
			TableId:               100,
			DcUserId:              &dcUserId,
			DcUserName:            "Пётр",
			DcUserExternalId:      subject,
			TableRolesName:        "table_reader",
			TableRolesDescription: "чтение таблицы",
			TableCatName:          "orders",
		},
	}}
	rec := httptest.NewRecorder()

	New(client).GetUserTableRolesById(rec, newRequest(t, http.MethodGet, "", nil), 7)

	assertStatus(t, rec, http.StatusOK)
	assertContentTypeJSON(t, rec)

	if got := client.getByIdReq.GetId(); got != 7 {
		t.Errorf("id в gRPC-запросе = %d, ожидался 7", got)
	}

	var got openapiv1.UserTableRolesDetail
	decode(t, rec, &got)

	if got.Id != 7 || got.TableRolesId != 3 || got.TableId != 100 {
		t.Errorf("строка = %+v, ожидались id=7, table_roles_id=3, table_id=100", got)
	}
	if got.DcUserId == nil || *got.DcUserId != dcUserId {
		t.Errorf("dc_user_id = %v, ожидался %d", got.DcUserId, dcUserId)
	}
	// table_cat_id в ответе не пришёл — наружу тоже не выдумывается.
	if got.TableCatId != nil {
		t.Errorf("table_cat_id = %v, ожидался nil", *got.TableCatId)
	}
}

// TestGetUserTableRolesByIdNotFound — несуществующий id обязан отдавать 404 с
// кодом not_found, а не 500.
func TestGetUserTableRolesByIdNotFound(t *testing.T) {
	client := &stubClient{err: status.Error(codes.NotFound, "dc.user_table_roles id = 7")}
	rec := httptest.NewRecorder()

	New(client).GetUserTableRolesById(rec, newRequest(t, http.MethodGet, "", nil), 7)

	assertStatus(t, rec, http.StatusNotFound)
	if got := decodeError(t, rec).Code; got != "not_found" {
		t.Errorf("code = %q, ожидался %q", got, "not_found")
	}
}

// TestCreateUserTableRolesUsesClaims — автор изменения берётся из токена.
// Тело намеренно содержит external_id: в схеме его нет, и подсунуть чужой id
// через него нельзя.
func TestCreateUserTableRolesUsesClaims(t *testing.T) {
	client := &stubClient{createResp: &userdomainrolesv1.CreateUserTableRolesResponse{
		UserTableRoles: &userdomainrolesv1.UserTableRoles{
			Id:           7,
			UserId:       42,
			TableRolesId: 3,
			CreatedAt:    timestamppb.New(createdAt),
			UpdatedAt:    timestamppb.New(updatedAt),
			TableId:      100,
			UpdatedById:  42,
		},
	}}
	rec := httptest.NewRecorder()

	body := `{"user_id":42,"table_roles_id":3,"table_id":100,"external_id":"подделка"}`
	New(client).CreateUserTableRoles(rec, newRequest(t, http.MethodPost, body, claims()))

	assertStatus(t, rec, http.StatusCreated)

	if got := client.createReq.GetExternalId(); got != subject {
		t.Errorf("external_id = %q, ожидался sub из claims %q", got, subject)
	}
	if got := client.createReq.GetUserId(); got != 42 {
		t.Errorf("user_id = %d, ожидался 42", got)
	}
	if got := client.createReq.GetTableRolesId(); got != 3 {
		t.Errorf("table_roles_id = %d, ожидался 3", got)
	}
	if got := client.createReq.GetTableId(); got != 100 {
		t.Errorf("table_id = %d, ожидался 100", got)
	}

	var created openapiv1.UserTableRoles
	decode(t, rec, &created)
	if created.Id != 7 || created.UserId != 42 || created.TableId != 100 {
		t.Errorf("тело ответа = %+v, ожидалась созданная строка id=7", created)
	}
}

// TestCreateUserTableRolesRejects — запрос без claims или с битым телом
// заворачивается до похода в Metadata Service.
func TestCreateUserTableRolesRejects(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		claims     *auth.Claims
		wantStatus int
		wantCode   string
	}{
		{
			name:       "без claims",
			body:       `{"user_id":42,"table_roles_id":3,"table_id":100}`,
			claims:     nil,
			wantStatus: http.StatusUnauthorized,
			wantCode:   "unauthenticated",
		},
		{
			name:       "битый JSON",
			body:       `{"user_id":`,
			claims:     claims(),
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_argument",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubClient{}
			rec := httptest.NewRecorder()

			New(client).CreateUserTableRoles(rec, newRequest(t, http.MethodPost, tt.body, tt.claims))

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

// TestUpdateUserTableRolesById — id берётся из пути, автор из claims, тело
// отдаётся пустым с 204.
func TestUpdateUserTableRolesById(t *testing.T) {
	client := &stubClient{}
	rec := httptest.NewRecorder()

	body := `{"user_id":42,"table_roles_id":9,"table_id":100}`
	New(client).UpdateUserTableRolesById(rec, newRequest(t, http.MethodPut, body, claims()), 7)

	assertStatus(t, rec, http.StatusNoContent)
	assertEmptyBody(t, rec)

	if got := client.updateReq.GetId(); got != 7 {
		t.Errorf("id = %d, ожидался 7", got)
	}
	if got := client.updateReq.GetTableRolesId(); got != 9 {
		t.Errorf("table_roles_id = %d, ожидался 9", got)
	}
	if got := client.updateReq.GetExternalId(); got != subject {
		t.Errorf("external_id = %q, ожидался sub из claims %q", got, subject)
	}
}

// TestUpdateUserTableRolesByIdRejects — та же пара отказов, что и у Create.
func TestUpdateUserTableRolesByIdRejects(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		claims     *auth.Claims
		wantStatus int
	}{
		{"без claims", `{"user_id":42,"table_roles_id":9,"table_id":100}`, nil, http.StatusUnauthorized},
		{"битый JSON", `{`, claims(), http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubClient{}
			rec := httptest.NewRecorder()

			New(client).UpdateUserTableRolesById(rec, newRequest(t, http.MethodPut, tt.body, tt.claims), 7)

			assertStatus(t, rec, tt.wantStatus)
			if client.calls != 0 {
				t.Errorf("вызовов gRPC = %d, ожидалось 0", client.calls)
			}
		})
	}
}

// TestDeleteUserTableRolesByIdNotFound — сервис проверяет существование строки
// заранее и на несуществующем id отдаёт NotFound; хендлер обязан перенести это
// в 404, а не в 204.
func TestDeleteUserTableRolesByIdNotFound(t *testing.T) {
	client := &stubClient{err: status.Error(codes.NotFound, "dc.user_table_roles id = 7")}
	rec := httptest.NewRecorder()

	New(client).DeleteUserTableRolesById(rec, newRequest(t, http.MethodDelete, "", claims()), 7)

	assertStatus(t, rec, http.StatusNotFound)
}

// TestDeleteUndeleteUserTableRolesById — обе ручки принимают только id и
// автора: тела запроса у них нет, ответ 204.
func TestDeleteUndeleteUserTableRolesById(t *testing.T) {
	t.Run("delete", func(t *testing.T) {
		client := &stubClient{}
		rec := httptest.NewRecorder()

		New(client).DeleteUserTableRolesById(rec, newRequest(t, http.MethodDelete, "", claims()), 7)

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

		New(client).UndeleteUserTableRolesById(rec, newRequest(t, http.MethodPost, "", claims()), 7)

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

// TestDeleteUndeleteUserTableRolesByIdWithoutClaims — без токена ни одна из
// них до бэкенда не доходит.
func TestDeleteUndeleteUserTableRolesByIdWithoutClaims(t *testing.T) {
	tests := []struct {
		name string
		call func(h *Handler, w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "delete",
			call: func(h *Handler, w http.ResponseWriter, r *http.Request) { h.DeleteUserTableRolesById(w, r, 7) },
		},
		{
			name: "undelete",
			call: func(h *Handler, w http.ResponseWriter, r *http.Request) { h.UndeleteUserTableRolesById(w, r, 7) },
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

	req := httptest.NewRequest(method, "/v1/user-table-roles", strings.NewReader(body))
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
