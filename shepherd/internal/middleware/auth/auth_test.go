package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func TestHasRole(t *testing.T) {
	tests := []struct {
		name  string
		have  []string
		want  []string
		match bool
	}{
		{name: "единственная роль совпала", have: []string{"admin"}, want: []string{"admin"}, match: true},
		{name: "совпала одна из перечисленных", have: []string{"maintainer"}, want: []string{"admin", "maintainer"}, match: true},
		{name: "совпала не первая из имеющихся", have: []string{"viewer", "admin"}, want: []string{"admin"}, match: true},
		{name: "роль не совпала", have: []string{"viewer"}, want: []string{"admin", "maintainer"}, match: false},
		{name: "ролей нет", have: nil, want: []string{"admin"}, match: false},
		{name: "регистр важен", have: []string{"Admin"}, want: []string{"admin"}, match: false},
		// Вызов без аргументов не должен случайно пропускать всех: пустой
		// список запрошенных ролей означает «ни одна не подходит».
		{name: "запрошенных ролей нет", have: []string{"admin"}, want: nil, match: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (Claims{Roles: tt.have}).HasRole(tt.want...); got != tt.match {
				t.Errorf("HasRole(%v) при ролях %v = %v, ожидалось %v", tt.want, tt.have, got, tt.match)
			}
		})
	}
}

func TestClaimsContext(t *testing.T) {
	t.Run("кладём и достаём", func(t *testing.T) {
		want := Claims{Subject: "sub-1", Name: "иван", Email: "ivan@example.com", Roles: []string{"admin"}}

		got, ok := ClaimsFromContext(ContextWithClaims(context.Background(), want))

		if !ok {
			t.Fatal("claims не найдены в контексте")
		}
		if got.Subject != want.Subject || got.Name != want.Name || got.Email != want.Email {
			t.Errorf("claims = %+v, ожидались %+v", got, want)
		}
		if len(got.Roles) != 1 || got.Roles[0] != "admin" {
			t.Errorf("roles = %v, ожидались [admin]", got.Roles)
		}
	})

	t.Run("в пустом контексте claims нет", func(t *testing.T) {
		if _, ok := ClaimsFromContext(context.Background()); ok {
			t.Error("claims нашлись в пустом контексте")
		}
	})
}

func TestBearerToken(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
	}{
		{name: "обычный токен", header: "Bearer abc.def.ghi", want: "abc.def.ghi"},
		{name: "лишние пробелы обрезаются", header: "Bearer   abc.def.ghi  ", want: "abc.def.ghi"},
		{name: "заголовка нет", header: "", want: ""},
		{name: "другая схема", header: "Basic dXNlcjpwYXNz", want: ""},
		{name: "без пробела после схемы", header: "Bearerabc.def.ghi", want: ""},
		{name: "пустой токен после схемы", header: "Bearer ", want: ""},
		// Схема в HTTP регистронезависима, но здесь сравнение точное.
		// Тест фиксирует текущее поведение: Keycloak и наш фронтенд шлют
		// ровно "Bearer", а расширение разбора — осознанное решение, а не
		// случайная правка.
		{name: "нижний регистр схемы не принимается", header: "bearer abc.def.ghi", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/hosts", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}

			if got := bearerToken(req); got != tt.want {
				t.Errorf("bearerToken() = %q, ожидался %q", got, tt.want)
			}
		})
	}
}

// TestRealmRolesFromToken — самый хрупкий кусок пакета: realm_access
// разбирается вручную через map[string]any/[]any, потому что jwx не умеет
// анмаршалить приватный claim в структуру. Если разбор сломается, роли станут
// пустыми и весь каталог начнёт молча отвечать 403.
//
// Токены здесь подписываются и парсятся по-настоящему, а не собираются в
// памяти: только после round-trip через JSON приватные claim'ы приобретают те
// самые типы any, ради которых написан ручной разбор.
func TestRealmRolesFromToken(t *testing.T) {
	tests := []struct {
		name        string
		realmAccess any
		omit        bool
		want        []string
	}{
		{
			name:        "обычный набор ролей",
			realmAccess: map[string]any{"roles": []any{"admin", "maintainer"}},
			want:        []string{"admin", "maintainer"},
		},
		{
			name:        "нестроковые элементы пропускаются",
			realmAccess: map[string]any{"roles": []any{"admin", 42, true, "viewer"}},
			want:        []string{"admin", "viewer"},
		},
		{
			name:        "пустой список ролей",
			realmAccess: map[string]any{"roles": []any{}},
			want:        []string{},
		},
		{name: "claim отсутствует", omit: true, want: nil},
		{name: "claim не объект", realmAccess: "admin", want: nil},
		{name: "claim — число", realmAccess: 42, want: nil},
		{name: "объект без поля roles", realmAccess: map[string]any{"other": "value"}, want: nil},
		{name: "roles не массив", realmAccess: map[string]any{"roles": "admin"}, want: nil},
		{name: "roles — объект", realmAccess: map[string]any{"roles": map[string]any{"a": 1}}, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := map[string]any{}
			if !tt.omit {
				claims["realm_access"] = tt.realmAccess
			}

			tok := signAndParse(t, primaryKey(t), tokenSpec{claims: claims})

			got := realmRolesFromToken(tok)

			if tt.want == nil && got != nil {
				t.Fatalf("роли = %v, ожидался nil", got)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("роли = %v, ожидались %v", got, tt.want)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("роль[%d] = %q, ожидалась %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestClaimsFromToken(t *testing.T) {
	t.Run("все поля заполнены", func(t *testing.T) {
		tok := signAndParse(t, primaryKey(t), tokenSpec{
			subject: "0d1c9f4e-0000-0000-0000-000000000001",
			claims: map[string]any{
				"preferred_username": "ivan",
				"email":              "ivan@example.com",
				"realm_access":       map[string]any{"roles": []any{"admin"}},
			},
		})

		got := claimsFromToken(tok)

		if got.Subject != "0d1c9f4e-0000-0000-0000-000000000001" {
			t.Errorf("Subject = %q", got.Subject)
		}
		if got.Name != "ivan" {
			t.Errorf("Name = %q, ожидалось ivan", got.Name)
		}
		if got.Email != "ivan@example.com" {
			t.Errorf("Email = %q", got.Email)
		}
		if len(got.Roles) != 1 || got.Roles[0] != "admin" {
			t.Errorf("Roles = %v, ожидались [admin]", got.Roles)
		}
	})

	// Keycloak может не отдавать email (не заполнен в профиле) или
	// preferred_username. Отсутствие необязательных claim'ов не должно
	// ломать разбор — ensureuser умеет откатываться на sub.
	t.Run("необязательные claim'ы отсутствуют", func(t *testing.T) {
		tok := signAndParse(t, primaryKey(t), tokenSpec{subject: "sub-без-профиля"})

		got := claimsFromToken(tok)

		if got.Subject != "sub-без-профиля" {
			t.Errorf("Subject = %q", got.Subject)
		}
		if got.Name != "" {
			t.Errorf("Name = %q, ожидалась пустая строка", got.Name)
		}
		if got.Email != "" {
			t.Errorf("Email = %q, ожидалась пустая строка", got.Email)
		}
		if got.Roles != nil {
			t.Errorf("Roles = %v, ожидался nil", got.Roles)
		}
	})
}

// TestAuthenticate проверяет middleware целиком: подпись по JWKS, issuer, срок
// действия. JWKS отдаёт httptest-сервер, Redis не задействован (fetcher умеет
// работать без него).
func TestAuthenticate(t *testing.T) {
	key := primaryKey(t)
	jwksSrv := newJWKSServer(t, key)
	issuer := jwksSrv.URL + "/realms/datacomrade"

	mw := NewMiddleware(NewJWKSFetcher(nil, issuer), issuer)

	t.Run("валидный токен пропускается", func(t *testing.T) {
		token := sign(t, key, tokenSpec{
			issuer:  issuer,
			subject: "sub-1",
			claims: map[string]any{
				"preferred_username": "ivan",
				"realm_access":       map[string]any{"roles": []any{"maintainer"}},
			},
		})

		var got Claims
		var called bool
		rec := serve(t, mw, token, func(_ http.ResponseWriter, r *http.Request) {
			called = true
			got, _ = ClaimsFromContext(r.Context())
		})

		if !called {
			t.Fatal("хендлер не был вызван")
		}
		if rec.Code != http.StatusOK {
			t.Errorf("статус = %d, ожидался %d", rec.Code, http.StatusOK)
		}
		if got.Subject != "sub-1" || got.Name != "ivan" {
			t.Errorf("claims = %+v", got)
		}
		if !got.HasRole("maintainer") {
			t.Errorf("роль maintainer не доехала до контекста: %v", got.Roles)
		}
	})

	rejections := []struct {
		name       string
		token      string
		wantStatus int
		wantCode   string
	}{
		{
			name:       "без токена",
			token:      "",
			wantStatus: http.StatusUnauthorized,
			wantCode:   "unauthenticated",
		},
		{
			name:       "мусор вместо токена",
			token:      "не-jwt",
			wantStatus: http.StatusUnauthorized,
			wantCode:   "invalid_token",
		},
		{
			name:       "чужой issuer",
			token:      sign(t, key, tokenSpec{issuer: "https://evil.example.com/realms/other", subject: "sub-1"}),
			wantStatus: http.StatusUnauthorized,
			wantCode:   "invalid_token",
		},
		{
			name: "истёкший токен",
			token: sign(t, key, tokenSpec{
				issuer:  issuer,
				subject: "sub-1",
				expires: time.Now().Add(-time.Hour),
			}),
			wantStatus: http.StatusUnauthorized,
			wantCode:   "invalid_token",
		},
		{
			// Тот же kid, но другой ключ: проверяется именно подпись, а не
			// наличие подходящего идентификатора в JWKS.
			name:       "подпись чужим ключом",
			token:      sign(t, impostorKey(t), tokenSpec{issuer: issuer, subject: "sub-1"}),
			wantStatus: http.StatusUnauthorized,
			wantCode:   "invalid_token",
		},
	}

	for _, tt := range rejections {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			rec := serve(t, mw, tt.token, func(http.ResponseWriter, *http.Request) { called = true })

			if called {
				t.Error("хендлер был вызван, хотя запрос должен быть отклонён")
			}
			assertErrorBody(t, rec, tt.wantStatus, tt.wantCode)
		})
	}
}

// TestAuthenticateJWKSUnavailable — Keycloak недоступен, ключей взять негде.
// Это 503, а не 401: токен пользователя может быть в полном порядке, проблема
// на нашей стороне, и клиенту следует повторить запрос, а не переходить к
// повторному логину.
func TestAuthenticateJWKSUnavailable(t *testing.T) {
	key := primaryKey(t)

	deadSrv := httptest.NewServer(http.NotFoundHandler())
	issuer := deadSrv.URL + "/realms/datacomrade"
	deadSrv.Close()

	mw := NewMiddleware(NewJWKSFetcher(nil, issuer), issuer)
	token := sign(t, key, tokenSpec{issuer: issuer, subject: "sub-1"})

	called := false
	rec := serve(t, mw, token, func(http.ResponseWriter, *http.Request) { called = true })

	if called {
		t.Error("хендлер был вызван без проверенного токена")
	}
	assertErrorBody(t, rec, http.StatusServiceUnavailable, "jwks_unavailable")
}

// --- вспомогательное ---

// signingKey — пара ключей для тестов: private подписывает токены, public
// уезжает в JWKS.
type signingKey struct {
	private jwk.Key
	public  jwk.Key
}

// Генерация RSA дорогая, а ключи между тестами не меняются — считаем их один
// раз на весь пакет. Оба ключа получают одинаковый kid, чтобы тест «подпись
// чужим ключом» проверял именно проверку подписи.
const testKeyID = "test-key"

var loadKeys = sync.OnceValues(func() ([2]signingKey, error) {
	var keys [2]signingKey
	for i := range keys {
		raw, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return keys, err
		}

		private, err := jwk.Import(raw)
		if err != nil {
			return keys, err
		}
		if err := setKeyMeta(private); err != nil {
			return keys, err
		}

		public, err := jwk.PublicKeyOf(private)
		if err != nil {
			return keys, err
		}
		if err := setKeyMeta(public); err != nil {
			return keys, err
		}

		keys[i] = signingKey{private: private, public: public}
	}

	return keys, nil
})

func setKeyMeta(key jwk.Key) error {
	if err := key.Set("kid", testKeyID); err != nil {
		return err
	}

	return key.Set("alg", jwa.RS256().String())
}

func primaryKey(t *testing.T) signingKey {
	t.Helper()

	keys, err := loadKeys()
	if err != nil {
		t.Fatalf("не удалось подготовить тестовые ключи: %v", err)
	}

	return keys[0]
}

// impostorKey — ключ, которого нет в JWKS, но с тем же kid.
func impostorKey(t *testing.T) signingKey {
	t.Helper()

	keys, err := loadKeys()
	if err != nil {
		t.Fatalf("не удалось подготовить тестовые ключи: %v", err)
	}

	return keys[1]
}

// tokenSpec описывает выпускаемый токен. Нулевой expires означает «валиден час».
type tokenSpec struct {
	issuer  string
	subject string
	expires time.Time
	claims  map[string]any
}

func sign(t *testing.T, key signingKey, spec tokenSpec) string {
	t.Helper()

	expires := spec.expires
	if expires.IsZero() {
		expires = time.Now().Add(time.Hour)
	}

	builder := jwt.NewBuilder().IssuedAt(time.Now().Add(-time.Minute)).Expiration(expires)
	if spec.issuer != "" {
		builder = builder.Issuer(spec.issuer)
	}
	if spec.subject != "" {
		builder = builder.Subject(spec.subject)
	}
	for name, value := range spec.claims {
		builder = builder.Claim(name, value)
	}

	tok, err := builder.Build()
	if err != nil {
		t.Fatalf("не удалось собрать токен: %v", err)
	}

	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.RS256(), key.private))
	if err != nil {
		t.Fatalf("не удалось подписать токен: %v", err)
	}

	return string(signed)
}

// signAndParse прогоняет токен через подпись и разбор, чтобы claim'ы получили
// те же типы, что и у настоящего токена Keycloak.
func signAndParse(t *testing.T, key signingKey, spec tokenSpec) jwt.Token {
	t.Helper()

	set := jwk.NewSet()
	if err := set.AddKey(key.public); err != nil {
		t.Fatalf("не удалось собрать набор ключей: %v", err)
	}

	tok, err := jwt.Parse([]byte(sign(t, key, spec)), jwt.WithKeySet(set))
	if err != nil {
		t.Fatalf("не удалось разобрать подписанный токен: %v", err)
	}

	return tok
}

// newJWKSServer поднимает подобие эндпоинта ключей Keycloak. Путь не важен:
// JWKSFetcher ходит по вычисленному адресу, а сервер отвечает на любой.
func newJWKSServer(t *testing.T, key signingKey) *httptest.Server {
	t.Helper()

	set := jwk.NewSet()
	if err := set.AddKey(key.public); err != nil {
		t.Fatalf("не удалось собрать набор ключей: %v", err)
	}

	body, err := json.Marshal(set)
	if err != nil {
		t.Fatalf("не удалось сериализовать JWKS: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	t.Cleanup(srv.Close)

	return srv
}

func serve(t *testing.T, mw *Middleware, token string, next http.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/hosts", nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	mw.Authenticate(next).ServeHTTP(rec, req)

	return rec
}

func assertErrorBody(t *testing.T, rec *httptest.ResponseRecorder, wantStatus int, wantCode string) {
	t.Helper()

	if rec.Code != wantStatus {
		t.Errorf("статус = %d, ожидался %d", rec.Code, wantStatus)
	}

	var body struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("тело не разобралось как JSON: %v (тело = %s)", err, rec.Body.String())
	}
	if body.Code != wantCode {
		t.Errorf("code = %q, ожидался %q", body.Code, wantCode)
	}
	if body.Message == "" {
		t.Error("message пустой — клиенту нечего показать")
	}
}
