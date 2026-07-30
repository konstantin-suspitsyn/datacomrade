// Package requestlog генерирует/пробрасывает trace_id и логирует каждый
// HTTP-запрос через platform/pkg/logger — см. Логирование.md.
package requestlog

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	pkglogger "github.com/konstantin-suspitsyn/datacomrade/platform/pkg/logger"
	"go.uber.org/zap"
)

const traceIDHeader = "X-Trace-Id"

// Middleware — самый внешний слой цепочки: должен вешаться раньше auth,
// чтобы trace_id был даже у отклонённых на аутентификации запросов.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get(traceIDHeader)
		if traceID == "" {
			traceID = uuid.NewString()
		}
		w.Header().Set(traceIDHeader, traceID)

		ctx := pkglogger.ContextWithTraceID(r.Context(), traceID)

		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r.WithContext(ctx))

		pkglogger.Info(ctx, "http request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", rec.status),
			zap.Duration("duration", time.Since(start)),
		)
	})
}

// statusRecorder запоминает статус ответа для лога — http.ResponseWriter сам
// его не отдаёт.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
