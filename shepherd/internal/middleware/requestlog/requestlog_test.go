package requestlog

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	pkglogger "github.com/konstantin-suspitsyn/datacomrade/platform/pkg/logger"
)

func TestMain(m *testing.M) {
	// Middleware логирует через глобальный логгер: без инициализации
	// pkglogger.Info разыменует nil.
	pkglogger.SetNopLogger()
	m.Run()
}

// TestMiddlewareGeneratesTraceID — запрос без X-Trace-Id получает свой
// идентификатор, и он уезжает клиенту в ответе: без этого связать жалобу
// пользователя с логами невозможно.
func TestMiddlewareGeneratesTraceID(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/hosts", nil)

	Middleware(okHandler()).ServeHTTP(rec, req)

	got := rec.Header().Get(traceIDHeader)
	if got == "" {
		t.Fatal("заголовок X-Trace-Id отсутствует в ответе")
	}
	if _, err := uuid.Parse(got); err != nil {
		t.Errorf("X-Trace-Id = %q не является UUID: %v", got, err)
	}
}

// TestMiddlewarePreservesTraceID — если фронтенд или вышестоящий прокси уже
// прислал trace_id, цепочка обязана его сохранить, а не начать свою.
func TestMiddlewarePreservesTraceID(t *testing.T) {
	const incoming = "1a2b3c4d-0000-0000-0000-000000000000"

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/hosts", nil)
	req.Header.Set(traceIDHeader, incoming)

	Middleware(okHandler()).ServeHTTP(rec, req)

	if got := rec.Header().Get(traceIDHeader); got != incoming {
		t.Errorf("X-Trace-Id = %q, ожидался пришедший %q", got, incoming)
	}
}

// TestMiddlewareSetsTraceIDBeforeHandler — заголовок должен быть выставлен до
// вызова хендлера. Иначе хендлер, успевший записать статус, зафиксирует ответ
// без trace_id: после WriteHeader правки заголовков теряются.
func TestMiddlewareSetsTraceIDBeforeHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/hosts", nil)

	var seen string
	Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		seen = w.Header().Get(traceIDHeader)
		w.WriteHeader(http.StatusCreated)
	})).ServeHTTP(rec, req)

	if seen == "" {
		t.Error("на входе в хендлер X-Trace-Id ещё не выставлен")
	}
	if got := rec.Header().Get(traceIDHeader); got != seen {
		t.Errorf("клиенту ушёл X-Trace-Id %q, а хендлер видел %q", got, seen)
	}
}

// TestMiddlewarePassesThroughResponse — обёртка не должна вмешиваться в ответ
// хендлера.
func TestMiddlewarePassesThroughResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/hosts", nil)

	Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		_, _ = w.Write([]byte("тело"))
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusTeapot {
		t.Errorf("статус = %d, ожидался %d", rec.Code, http.StatusTeapot)
	}
	if rec.Body.String() != "тело" {
		t.Errorf("тело = %q, ожидалось %q", rec.Body.String(), "тело")
	}
}

// TestStatusRecorder — statusRecorder существует только чтобы в лог попал
// настоящий статус ответа: http.ResponseWriter его не отдаёт.
func TestStatusRecorder(t *testing.T) {
	t.Run("без WriteHeader остаётся 200", func(t *testing.T) {
		rec := httptest.NewRecorder()
		sr := &statusRecorder{ResponseWriter: rec, status: http.StatusOK}

		_, _ = sr.Write([]byte("тело"))

		if sr.status != http.StatusOK {
			t.Errorf("status = %d, ожидался %d", sr.status, http.StatusOK)
		}
	})

	t.Run("запоминает и пробрасывает статус", func(t *testing.T) {
		rec := httptest.NewRecorder()
		sr := &statusRecorder{ResponseWriter: rec, status: http.StatusOK}

		sr.WriteHeader(http.StatusForbidden)

		if sr.status != http.StatusForbidden {
			t.Errorf("запомненный status = %d, ожидался %d", sr.status, http.StatusForbidden)
		}
		if rec.Code != http.StatusForbidden {
			t.Errorf("статус не дошёл до ответа: %d", rec.Code)
		}
	})
}

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
