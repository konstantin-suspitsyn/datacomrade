package userdomainrolesapiv1

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

	getReq  *userdomainrolesv1.GetUserDomainRolesesRequest
	getResp *userdomainrolesv1.GetUserDomainRolesesResponse

	createReq  *userdomainrolesv1.CreateUserDomainRolesRequest
	createResp *userdomainrolesv1.CreateUserDomainRolesResponse

	updateReq   *userdomainrolesv1.UpdateUserDomainRolesByIdRequest
	deleteReq   *userdomainrolesv1.DeleteUserDomainRolesByIdRequest
	undeleteReq *userdomainrolesv1.UndeleteUserDomainRolesByIdRequest
}

func (s *stubClient) GetUserDomainRoleses(_ context.Context, in *userdomainrolesv1.GetUserDomainRolesesRequest, _ ...grpc.CallOption) (*userdomainrolesv1.GetUserDomainRolesesResponse, error) {
	s.calls++
	s.getReq = in

	if s.err != nil {
		return nil, s.err
	}

	return s.getResp, nil
}

func (s *stubClient) CreateUserDomainRoles(_ context.Context, in *userdomainrolesv1.CreateUserDomainRolesRequest, _ ...grpc.CallOption) (*userdomainrolesv1.CreateUserDomainRolesResponse, error) {
	s.calls++
	s.createReq = in

	if s.err != nil {
		return nil, s.err
	}

	return s.createResp, nil
}

func (s *stubClient) UpdateUserDomainRolesById(_ context.Context, in *userdomainrolesv1.UpdateUserDomainRolesByIdRequest, _ ...grpc.CallOption) (*userdomainrolesv1.UpdateUserDomainRolesByIdResponse, error) {
	s.calls++
	s.updateReq = in

	if s.err != nil {
		return nil, s.err
	}

	return &userdomainrolesv1.UpdateUserDomainRolesByIdResponse{}, nil
}

func (s *stubClient) DeleteUserDomainRolesById(_ context.Context, in *userdomainrolesv1.DeleteUserDomainRolesByIdRequest, _ ...grpc.CallOption) (*userdomainrolesv1.DeleteUserDomainRolesByIdResponse, error) {
	s.calls++
	s.deleteReq = in

	if s.err != nil {
		return nil, s.err
	}

	return &userdomainrolesv1.DeleteUserDomainRolesByIdResponse{}, nil
}

func (s *stubClient) UndeleteUserDomainRolesById(_ context.Context, in *userdomainrolesv1.UndeleteUserDomainRolesByIdRequest, _ ...grpc.CallOption) (*userdomainrolesv1.UndeleteUserDomainRolesByIdResponse, error) {
	s.calls++
	s.undeleteReq = in

	if s.err != nil {
		return nil, s.err
	}

	return &userdomainrolesv1.UndeleteUserDomainRolesByIdResponse{}, nil
}

// TestListUserDomainRolesPagination фиксирует контракт постраничности:
// непереданные order/page_limit подменяются умолчаниями shepherd, переданные
// уходят в gRPC как есть, а order вне ASC/DESC — 400, а не молчаливый дефолт.
func TestListUserDomainRolesPagination(t *testing.T) {
	order := "ASC"
	page := int32(3)
	pageLimit := int32(25)
	badOrder := "asc"

	tests := []struct {
		name          string
		params        openapiv1.ListUserDomainRolesParams
		wantStatus    int
		wantCalls     int
		wantOrder     string
		wantPage      int32
		wantPageLimit int32
	}{
		{
			name:          "умолчания вместо непереданных параметров",
			params:        openapiv1.ListUserDomainRolesParams{},
			wantStatus:    http.StatusOK,
			wantCalls:     1,
			wantOrder:     pagination.DefaultOrder,
			wantPage:      0,
			wantPageLimit: pagination.DefaultPageLimit,
		},
		{
			name:          "переданные параметры уходят в gRPC",
			params:        openapiv1.ListUserDomainRolesParams{Order: &order, Page: &page, PageLimit: &pageLimit},
			wantStatus:    http.StatusOK,
			wantCalls:     1,
			wantOrder:     order,
			wantPage:      page,
			wantPageLimit: pageLimit,
		},
		{
			name:       "order вне ASC/DESC отсекается до gRPC",
			params:     openapiv1.ListUserDomainRolesParams{Order: &badOrder},
			wantStatus: http.StatusBadRequest,
			wantCalls:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubClient{getResp: &userdomainrolesv1.GetUserDomainRolesesResponse{}}
			rec := httptest.NewRecorder()

			New(client).ListUserDomainRoles(rec, newRequest(t, http.MethodGet, "", nil), tt.params)

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

// TestListUserDomainRolesBody проверяет форму ответа: имена полей — то, на что
// смотрит фронтенд, и переименование в конвертере он молча не переживёт.
func TestListUserDomainRolesBody(t *testing.T) {
	client := &stubClient{getResp: &userdomainrolesv1.GetUserDomainRolesesResponse{
		Data: []*userdomainrolesv1.GetUserDomainRolesesRow{{
			Id:                     7,
			UserId:                 42,
			DomainRolesId:          3,
			CreatedAt:              timestamppb.New(createdAt),
			UpdatedAt:              timestamppb.New(updatedAt),
			IsDeleted:              false,
			DomainId:               11,
			UpdatedById:            42,
			DomainRolesName:        "domain_admin",
			DomainRolesDescription: "полные права на домен",
			DcUserName:             "Пётр",
		}},
		Pagination: &userdomainrolesv1.Pagination{Page: 1, PerPage: 100, TotalItems: 1, TotalPages: 1},
	}}
	rec := httptest.NewRecorder()

	New(client).ListUserDomainRoles(rec, newRequest(t, http.MethodGet, "", nil), openapiv1.ListUserDomainRolesParams{})

	assertStatus(t, rec, http.StatusOK)
	assertContentTypeJSON(t, rec)

	var got openapiv1.UserDomainRolesPage
	decode(t, rec, &got)

	if len(got.Data) != 1 {
		t.Fatalf("строк в data = %d, ожидалась 1", len(got.Data))
	}

	row := got.Data[0]
	want := openapiv1.UserDomainRolesRow{
		Id:                     7,
		UserId:                 42,
		DomainRolesId:          3,
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
		IsDeleted:              false,
		DomainId:               11,
		UpdatedById:            42,
		DomainRolesName:        "domain_admin",
		DomainRolesDescription: "полные права на домен",
		DcUserName:             "Пётр",
	}
	if !row.CreatedAt.Equal(want.CreatedAt) || !row.UpdatedAt.Equal(want.UpdatedAt) {
		t.Errorf("метки времени = %v/%v, ожидались %v/%v", row.CreatedAt, row.UpdatedAt, want.CreatedAt, want.UpdatedAt)
	}

	row.CreatedAt, row.UpdatedAt = want.CreatedAt, want.UpdatedAt
	if row != want {
		t.Errorf("строка = %+v, ожидалась %+v", row, want)
	}

	wantPagination := openapiv1.Pagination{Page: 1, PerPage: 100, TotalItems: 1, TotalPages: 1}
	if got.Pagination != wantPagination {
		t.Errorf("pagination = %+v, ожидалась %+v", got.Pagination, wantPagination)
	}
}

// TestListUserDomainRolesGRPCError — ошибка бэкенда доезжает до клиента со
// своим кодом, а не схлопывается в 500.
func TestListUserDomainRolesGRPCError(t *testing.T) {
	client := &stubClient{err: status.Error(codes.Unavailable, "datacatalogue недоступен")}
	rec := httptest.NewRecorder()

	New(client).ListUserDomainRoles(rec, newRequest(t, http.MethodGet, "", nil), openapiv1.ListUserDomainRolesParams{})

	assertStatus(t, rec, http.StatusServiceUnavailable)
	if got := decodeError(t, rec).Code; got != "unavailable" {
		t.Errorf("code = %q, ожидался %q", got, "unavailable")
	}
}

// TestCreateUserDomainRolesUsesClaims — ключевое свойство всех пишущих ручек:
// автор изменения берётся из токена. Тело намеренно содержит external_id —
// его в схеме нет, и подсунуть чужой id через него нельзя.
func TestCreateUserDomainRolesUsesClaims(t *testing.T) {
	client := &stubClient{createResp: &userdomainrolesv1.CreateUserDomainRolesResponse{
		UserDomainRoles: &userdomainrolesv1.UserDomainRoles{
			Id:            7,
			UserId:        42,
			DomainRolesId: 3,
			CreatedAt:     timestamppb.New(createdAt),
			UpdatedAt:     timestamppb.New(updatedAt),
			DomainId:      11,
			UpdatedById:   42,
		},
	}}
	rec := httptest.NewRecorder()

	body := `{"user_id":42,"domain_roles_id":3,"domain_id":11,"external_id":"подделка"}`
	New(client).CreateUserDomainRoles(rec, newRequest(t, http.MethodPost, body, claims()))

	assertStatus(t, rec, http.StatusCreated)

	if got := client.createReq.GetExternalId(); got != subject {
		t.Errorf("external_id = %q, ожидался sub из claims %q", got, subject)
	}
	if got := client.createReq.GetUserId(); got != 42 {
		t.Errorf("user_id = %d, ожидался 42", got)
	}
	if got := client.createReq.GetDomainRolesId(); got != 3 {
		t.Errorf("domain_roles_id = %d, ожидался 3", got)
	}
	if got := client.createReq.GetDomainId(); got != 11 {
		t.Errorf("domain_id = %d, ожидался 11", got)
	}

	var created openapiv1.UserDomainRoles
	decode(t, rec, &created)
	if created.Id != 7 || created.UserId != 42 || created.DomainId != 11 {
		t.Errorf("тело ответа = %+v, ожидалась созданная строка id=7", created)
	}
}

// TestCreateUserDomainRolesRejects — запрос без claims или с битым телом
// заворачивается до похода в Metadata Service.
func TestCreateUserDomainRolesRejects(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		claims     *auth.Claims
		wantStatus int
		wantCode   string
	}{
		{
			name:       "без claims",
			body:       `{"user_id":42,"domain_roles_id":3,"domain_id":11}`,
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

			New(client).CreateUserDomainRoles(rec, newRequest(t, http.MethodPost, tt.body, tt.claims))

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

// TestCreateUserDomainRolesConflict — повторная выдача той же роли упирается в
// уникальный индекс: AlreadyExists от бэкенда обязан стать 409.
func TestCreateUserDomainRolesConflict(t *testing.T) {
	client := &stubClient{err: status.Error(codes.AlreadyExists, "роль уже выдана")}
	rec := httptest.NewRecorder()

	body := `{"user_id":42,"domain_roles_id":3,"domain_id":11}`
	New(client).CreateUserDomainRoles(rec, newRequest(t, http.MethodPost, body, claims()))

	assertStatus(t, rec, http.StatusConflict)
}

// TestUpdateUserDomainRolesById — id берётся из пути, автор из claims, тело
// отдаётся пустым с 204.
func TestUpdateUserDomainRolesById(t *testing.T) {
	client := &stubClient{}
	rec := httptest.NewRecorder()

	body := `{"user_id":42,"domain_roles_id":9,"domain_id":11}`
	New(client).UpdateUserDomainRolesById(rec, newRequest(t, http.MethodPut, body, claims()), 7)

	assertStatus(t, rec, http.StatusNoContent)
	assertEmptyBody(t, rec)

	if got := client.updateReq.GetId(); got != 7 {
		t.Errorf("id = %d, ожидался 7", got)
	}
	if got := client.updateReq.GetDomainRolesId(); got != 9 {
		t.Errorf("domain_roles_id = %d, ожидался 9", got)
	}
	if got := client.updateReq.GetExternalId(); got != subject {
		t.Errorf("external_id = %q, ожидался sub из claims %q", got, subject)
	}
}

// TestUpdateUserDomainRolesByIdRejects — та же пара отказов, что и у Create.
func TestUpdateUserDomainRolesByIdRejects(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		claims     *auth.Claims
		wantStatus int
	}{
		{"без claims", `{"user_id":42,"domain_roles_id":9,"domain_id":11}`, nil, http.StatusUnauthorized},
		{"битый JSON", `{`, claims(), http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubClient{}
			rec := httptest.NewRecorder()

			New(client).UpdateUserDomainRolesById(rec, newRequest(t, http.MethodPut, tt.body, tt.claims), 7)

			assertStatus(t, rec, tt.wantStatus)
			if client.calls != 0 {
				t.Errorf("вызовов gRPC = %d, ожидалось 0", client.calls)
			}
		})
	}
}

// TestUpdateUserDomainRolesByIdNotFound — обновление несуществующей строки
// обязано отдавать 404, а не 204.
func TestUpdateUserDomainRolesByIdNotFound(t *testing.T) {
	client := &stubClient{err: status.Error(codes.NotFound, "dc.user_domain_roles id = 7")}
	rec := httptest.NewRecorder()

	body := `{"user_id":42,"domain_roles_id":9,"domain_id":11}`
	New(client).UpdateUserDomainRolesById(rec, newRequest(t, http.MethodPut, body, claims()), 7)

	assertStatus(t, rec, http.StatusNotFound)
}

// TestDeleteUndeleteUserDomainRolesById — обе ручки принимают только id и
// автора: тела запроса у них нет, ответ 204.
func TestDeleteUndeleteUserDomainRolesById(t *testing.T) {
	t.Run("delete", func(t *testing.T) {
		client := &stubClient{}
		rec := httptest.NewRecorder()

		New(client).DeleteUserDomainRolesById(rec, newRequest(t, http.MethodDelete, "", claims()), 7)

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

		New(client).UndeleteUserDomainRolesById(rec, newRequest(t, http.MethodPost, "", claims()), 7)

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

// TestDeleteUndeleteUserDomainRolesByIdWithoutClaims — без токена ни одна из
// них до бэкенда не доходит.
func TestDeleteUndeleteUserDomainRolesByIdWithoutClaims(t *testing.T) {
	tests := []struct {
		name string
		call func(h *Handler, w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "delete",
			call: func(h *Handler, w http.ResponseWriter, r *http.Request) { h.DeleteUserDomainRolesById(w, r, 7) },
		},
		{
			name: "undelete",
			call: func(h *Handler, w http.ResponseWriter, r *http.Request) { h.UndeleteUserDomainRolesById(w, r, 7) },
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

	req := httptest.NewRequest(method, "/v1/user-domain-roles", strings.NewReader(body))
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
