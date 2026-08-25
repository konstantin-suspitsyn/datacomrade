package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	pkglogger "github.com/konstantin-suspitsyn/datacomrade/platform/pkg/logger"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/rbac"
	"github.com/redis/go-redis/v9"
)

func TestMain(m *testing.M) {
	// Обработчики логируют через глобальный логгер: без инициализации
	// pkglogger.Warn разыменует nil.
	pkglogger.SetNopLogger()
	m.Run()
}

// stubPinger подменяет Redis: возвращает заданную ошибку пинга.
type stubPinger struct {
	err error
}

func (s stubPinger) Ping(context.Context) *redis.StatusCmd {
	return redis.NewStatusResult("PONG", s.err)
}

// TestHandleReadyz фиксирует главное свойство пробы: недоступный Redis
// деградирует ответ, но не выводит инстанс из балансировки (всегда 200) —
// см. Кеширование Redis.md.
func TestHandleReadyz(t *testing.T) {
	tests := []struct {
		name    string
		pingErr error
		want    healthBody
	}{
		{
			name:    "redis доступен",
			pingErr: nil,
			want:    healthBody{Status: "ok", Redis: "up"},
		},
		{
			name:    "redis недоступен",
			pingErr: errors.New("dial tcp: connection refused"),
			want:    healthBody{Status: "degraded", Redis: "down"},
		},
		{
			name:    "redis вернул redis.Nil",
			pingErr: redis.Nil,
			want:    healthBody{Status: "degraded", Redis: "down"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/readyz", nil)

			handleReadyz(stubPinger{err: tt.pingErr})(rec, req)

			assertHealthResponse(t, rec, tt.want)
		})
	}
}

// TestHandleHealthz проверяет, что проба живости не утверждает ничего о кеше:
// поле redis в её ответе отсутствует.
func TestHandleHealthz(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	handleHealthz(rec, req)

	assertHealthResponse(t, rec, healthBody{Status: "ok"})

	var raw map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &raw); err != nil {
		t.Fatalf("тело не разобралось как JSON: %v", err)
	}
	if _, ok := raw["redis"]; ok {
		t.Errorf("healthz не должен сообщать состояние redis, тело = %s", rec.Body.String())
	}
}

func assertHealthResponse(t *testing.T, rec *httptest.ResponseRecorder, want healthBody) {
	t.Helper()

	if rec.Code != http.StatusOK {
		t.Errorf("статус = %d, ожидался %d", rec.Code, http.StatusOK)
	}

	const wantContentType = "application/json; charset=utf-8"
	if got := rec.Header().Get("Content-Type"); got != wantContentType {
		t.Errorf("Content-Type = %q, ожидался %q", got, wantContentType)
	}

	var got healthBody
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("тело не разобралось как JSON: %v (тело = %s)", err, rec.Body.String())
	}
	if got != want {
		t.Errorf("тело = %+v, ожидалось %+v", got, want)
	}
}

// ungatedMethods — методы, намеренно доступные любому аутентифицированному
// пользователю независимо от realm-роли: POST /login создаёт dc.user при
// первом входе (гейт по роли сделал бы вход невозможным для нового
// пользователя), GET /me отдаёт его же identity. См. комментарий про
// CreateUser в middleware/rbac.
//
// Всё остальное на api обязано быть за requireCatalogRole или
// requireUserRolesRole — это проверяет TestAPICatalogMethodsRequireRole.
var ungatedMethods = map[string]struct{}{
	"Login": {},
	"GetMe": {},

	// ListOlapAliases — доступ каждому пользователю уже ограничен на уровне
	// SQL (dc.get_table_ids_by_external_user_and_role, роль can_read), а не
	// ролью Gateway: нужна только аутентификация, как у Login/GetMe.
	"ListOlapAliases": {},
}

// adminOnlyMethods — методы выдачи прав пользователям (dc.user_domain_roles и
// dc.user_table_roles). Они стоят за requireUserRolesRole, где достаточной
// ролью является только admin: для maintainer ожидаемый исход — 403, а не
// проход к хендлеру. См. mapping роль → группа методов в rbac.go.
var adminOnlyMethods = map[string]struct{}{
	"ListUserDomainRoles":         {},
	"CreateUserDomainRoles":       {},
	"UpdateUserDomainRolesById":   {},
	"DeleteUserDomainRolesById":   {},
	"UndeleteUserDomainRolesById": {},

	"ListUserTableRoles":         {},
	"FilterUserTableRoles":       {},
	"CreateUserTableRoles":       {},
	"GetUserTableRolesById":      {},
	"UpdateUserTableRolesById":   {},
	"DeleteUserTableRolesById":   {},
	"UndeleteUserTableRolesById": {},
}

// TestAPIImplementsServerInterface ловит расхождение между набором методов api
// и сгенерированным ServerInterface. Нужен как страховка для обхода
// рефлексией ниже: тот перебирает методы api, и лишний экспортированный метод
// (или пропущенный интерфейсный) молча изменил бы список проверяемых ручек.
func TestAPIImplementsServerInterface(t *testing.T) {
	apiType := reflect.TypeOf(api{})
	ifaceType := reflect.TypeOf((*openapiv1.ServerInterface)(nil)).Elem()

	if apiType.NumMethod() != ifaceType.NumMethod() {
		t.Errorf("у api %d методов, у ServerInterface %d — рефлексия проверит не тот набор ручек",
			apiType.NumMethod(), ifaceType.NumMethod())
	}

	for name := range ungatedMethods {
		if _, ok := apiType.MethodByName(name); !ok {
			t.Errorf("метод %s из ungatedMethods не найден на api — опечатка в имени", name)
		}
	}

	for name := range adminOnlyMethods {
		if _, ok := apiType.MethodByName(name); !ok {
			t.Errorf("метод %s из adminOnlyMethods не найден на api — опечатка в имени", name)
		}
	}
}

// TestAPICatalogMethodsRequireRole — главный тест обвязки: каждый метод
// каталога на api должен быть обёрнут в requireCatalogRole, а метод из
// adminOnlyMethods — в requireUserRolesRole. Обёртки проставлены руками у
// сотни делегатов, и забытая у нового метода — молчаливая дыра в доступе,
// которую больше ничто не поймает.
func TestAPICatalogMethodsRequireRole(t *testing.T) {
	scenarios := []struct {
		name   string
		claims *auth.Claims
		check  func(t *testing.T, method string, res callResult)
		// checkAdminOnly — ожидание для методов из adminOnlyMethods, если оно
		// отличается от check. nil — сценарий одинаков для всех методов.
		checkAdminOnly func(t *testing.T, method string, res callResult)
	}{
		{
			name:   "без claims",
			claims: nil,
			check: func(t *testing.T, method string, res callResult) {
				t.Helper()
				if res.panicked {
					t.Fatalf("%s: гейт пропустил запрос без claims", method)
				}
				if res.code != http.StatusUnauthorized {
					t.Errorf("%s: статус = %d, ожидался %d", method, res.code, http.StatusUnauthorized)
				}
			},
		},
		{
			name:   "роль viewer",
			claims: &auth.Claims{Subject: "sub", Roles: []string{rbac.RoleViewer}},
			check:  assertGateBlocked,
		},
		{
			name:           "роль maintainer",
			claims:         &auth.Claims{Subject: "sub", Roles: []string{rbac.RoleMaintainer}},
			check:          assertGatePassed,
			checkAdminOnly: assertGateBlocked,
		},
		{
			name:   "роль admin",
			claims: &auth.Claims{Subject: "sub", Roles: []string{rbac.RoleAdmin}},
			check:  assertGatePassed,
		},
	}

	apiType := reflect.TypeOf(api{})
	for i := range apiType.NumMethod() {
		method := apiType.Method(i)
		if _, ok := ungatedMethods[method.Name]; ok {
			continue
		}

		_, adminOnly := adminOnlyMethods[method.Name]

		t.Run(method.Name, func(t *testing.T) {
			for _, sc := range scenarios {
				check := sc.check
				if adminOnly && sc.checkAdminOnly != nil {
					check = sc.checkAdminOnly
				}

				t.Run(sc.name, func(t *testing.T) {
					check(t, method.Name, callAPIMethod(method, sc.claims))
				})
			}
		})
	}
}

// TestUngatedMethodsPassAnyRole — обратная проверка: /login и /me не должны
// отсекаться по роли. Пользователь с одним лишь viewer обязан дойти до
// хендлера.
func TestUngatedMethodsPassAnyRole(t *testing.T) {
	apiType := reflect.TypeOf(api{})
	claims := &auth.Claims{Subject: "sub", Roles: []string{rbac.RoleViewer}}

	for name := range ungatedMethods {
		t.Run(name, func(t *testing.T) {
			method, ok := apiType.MethodByName(name)
			if !ok {
				t.Fatalf("метод %s не найден на api", name)
			}

			assertGatePassed(t, name, callAPIMethod(method, claims))
		})
	}
}

// assertGatePassed проверяет, что запрос не отсечён на аутентификации или
// проверке роли. Что именно вернёт сам хендлер — предмет тестов его пакета.
func assertGatePassed(t *testing.T, method string, res callResult) {
	t.Helper()

	if res.panicked {
		return
	}
	if res.code == http.StatusUnauthorized || res.code == http.StatusForbidden {
		t.Errorf("%s: гейт отсёк запрос со статусом %d, хотя роли достаточно", method, res.code)
	}
}

// assertGateBlocked — обратная проверка: роли не хватает, до под-хендлера
// управление дойти не должно.
func assertGateBlocked(t *testing.T, method string, res callResult) {
	t.Helper()

	if res.panicked {
		t.Fatalf("%s: гейт пропустил запрос с недостаточной ролью", method)
	}
	if res.code != http.StatusForbidden {
		t.Errorf("%s: статус = %d, ожидался %d", method, res.code, http.StatusForbidden)
	}
}

// callResult — исход вызова метода api. panicked=true означает, что управление
// дошло до под-хендлера: у нулевого api{} все поля nil, поэтому делегат
// разыменовывает nil-указатель. Для этих тестов паника — надёжный признак
// «гейт запрос пропустил»: сам гейт паниковать не может, он пишет ответ и
// возвращается, не касаясь под-хендлеров.
type callResult struct {
	code     int
	panicked bool
}

// callAPIMethod вызывает метод api через reflect на нулевом api{}, подставляя
// нулевые значения params/id — на проверку роли они не влияют.
func callAPIMethod(method reflect.Method, claims *auth.Claims) (res callResult) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/probe", nil)
	if claims != nil {
		req = req.WithContext(auth.ContextWithClaims(req.Context(), *claims))
	}

	// In(0) — получатель api, In(1) и In(2) — w и r, дальше params или id.
	args := make([]reflect.Value, method.Type.NumIn())
	args[0] = reflect.ValueOf(api{})
	args[1] = reflect.ValueOf(rec)
	args[2] = reflect.ValueOf(req)
	for i := 3; i < method.Type.NumIn(); i++ {
		args[i] = reflect.Zero(method.Type.In(i))
	}

	defer func() {
		if r := recover(); r != nil {
			res = callResult{panicked: true}
		}
	}()

	method.Func.Call(args)

	return callResult{code: rec.Code}
}

// TestMountDocs фиксирует поверхность документации: ровно два GET-роута и
// ничего сверх того. main вешает их только при cfg.EnableDocs — в проде
// Swagger UI выключен, см. API Gateway.md.
func TestMountDocs(t *testing.T) {
	router := chi.NewRouter()
	mountDocs(router)

	t.Run("монтирует ровно два роута", func(t *testing.T) {
		var routes []string
		err := chi.Walk(router, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
			routes = append(routes, method+" "+route)
			return nil
		})
		if err != nil {
			t.Fatalf("обход роутов не удался: %v", err)
		}
		sort.Strings(routes)

		want := []string{"GET /docs", "GET /openapi.yaml"}
		if !reflect.DeepEqual(routes, want) {
			t.Errorf("роуты = %v, ожидались %v", routes, want)
		}
	})

	t.Run("docs отдаёт Swagger UI", func(t *testing.T) {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/docs", nil))

		if rec.Code != http.StatusOK {
			t.Errorf("статус = %d, ожидался %d", rec.Code, http.StatusOK)
		}

		const wantContentType = "text/html; charset=utf-8"
		if got := rec.Header().Get("Content-Type"); got != wantContentType {
			t.Errorf("Content-Type = %q, ожидался %q", got, wantContentType)
		}

		// Страница ссылается на второй роут: переименование спеки без правки
		// шаблона оставило бы Swagger UI с пустой схемой.
		if !strings.Contains(rec.Body.String(), `url: "/openapi.yaml"`) {
			t.Errorf("страница не ссылается на /openapi.yaml, тело = %s", rec.Body.String())
		}
	})
}
