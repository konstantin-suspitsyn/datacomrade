package apierror

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestWrite(t *testing.T) {
	rec := httptest.NewRecorder()

	Write(rec, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")

	if rec.Code != http.StatusBadRequest {
		t.Errorf("статус = %d, ожидался %d", rec.Code, http.StatusBadRequest)
	}

	const wantContentType = "application/json; charset=utf-8"
	if got := rec.Header().Get("Content-Type"); got != wantContentType {
		t.Errorf("Content-Type = %q, ожидался %q", got, wantContentType)
	}

	want := Body{Code: "invalid_argument", Message: "order должен быть ASC или DESC"}
	if got := decodeBody(t, rec); got != want {
		t.Errorf("тело = %+v, ожидалось %+v", got, want)
	}
}

// TestWriteGRPCMapsCodes фиксирует таблицу соответствий gRPC → HTTP: это
// внешний контракт REST-периметра, на него завязан фронтенд.
func TestWriteGRPCMapsCodes(t *testing.T) {
	tests := []struct {
		name       string
		code       codes.Code
		wantStatus int
		wantCode   string
	}{
		{"InvalidArgument", codes.InvalidArgument, http.StatusBadRequest, "invalid_argument"},
		{"NotFound", codes.NotFound, http.StatusNotFound, "not_found"},
		{"AlreadyExists", codes.AlreadyExists, http.StatusConflict, "already_exists"},
		{"PermissionDenied", codes.PermissionDenied, http.StatusForbidden, "permission_denied"},
		{"Unauthenticated", codes.Unauthenticated, http.StatusUnauthorized, "unauthenticated"},
		{"DeadlineExceeded", codes.DeadlineExceeded, http.StatusGatewayTimeout, "deadline_exceeded"},
		{"Unavailable", codes.Unavailable, http.StatusServiceUnavailable, "unavailable"},
		// Всё неперечисленное схлопывается в 500: детали внутренних сбоев
		// наружу не выносим.
		{"Internal", codes.Internal, http.StatusInternalServerError, "internal"},
		{"Unknown", codes.Unknown, http.StatusInternalServerError, "internal"},
		{"ResourceExhausted", codes.ResourceExhausted, http.StatusInternalServerError, "internal"},
		{"Aborted", codes.Aborted, http.StatusInternalServerError, "internal"},
		{"Unimplemented", codes.Unimplemented, http.StatusInternalServerError, "internal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			WriteGRPC(rec, status.Error(tt.code, "сообщение бэкенда"))

			if rec.Code != tt.wantStatus {
				t.Errorf("статус = %d, ожидался %d", rec.Code, tt.wantStatus)
			}

			want := Body{Code: tt.wantCode, Message: "сообщение бэкенда"}
			if got := decodeBody(t, rec); got != want {
				t.Errorf("тело = %+v, ожидалось %+v", got, want)
			}
		})
	}
}

// TestWriteGRPCWrappedStatus — хендлеры каталога отдают ошибку клиента как
// есть, но если её когда-нибудь обернут через %w, код статуса обязан
// сохраниться, а не схлопнуться в 500.
func TestWriteGRPCWrappedStatus(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteGRPC(rec, fmt.Errorf("get host by id: %w", status.Error(codes.NotFound, "строка не найдена")))

	if rec.Code != http.StatusNotFound {
		t.Errorf("статус = %d, ожидался %d", rec.Code, http.StatusNotFound)
	}
	if got := decodeBody(t, rec).Code; got != "not_found" {
		t.Errorf("код = %q, ожидался not_found", got)
	}
}

// TestWriteGRPCPlainError — ошибка не от gRPC (например, сбой в самом
// Gateway) отдаётся как 500, а её текст попадает в message.
func TestWriteGRPCPlainError(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteGRPC(rec, errors.New("что-то сломалось"))

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("статус = %d, ожидался %d", rec.Code, http.StatusInternalServerError)
	}

	want := Body{Code: "internal", Message: "что-то сломалось"}
	if got := decodeBody(t, rec); got != want {
		t.Errorf("тело = %+v, ожидалось %+v", got, want)
	}
}

// TestWriteGRPCNilError фиксирует неочевидный край: status.FromError(nil)
// возвращает ok=true с кодом OK, который не перечислен в mapCode и потому
// схлопывается в 500 с пустым message. Вызывать WriteGRPC на nil-ошибке
// нельзя — тест закрепляет, что это не тихий 200.
func TestWriteGRPCNilError(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteGRPC(rec, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("статус = %d, ожидался %d", rec.Code, http.StatusInternalServerError)
	}
	if got := decodeBody(t, rec).Code; got != "internal" {
		t.Errorf("код = %q, ожидался internal", got)
	}
}

func decodeBody(t *testing.T, rec *httptest.ResponseRecorder) Body {
	t.Helper()

	var body Body
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("тело не разобралось как JSON: %v (тело = %s)", err, rec.Body.String())
	}

	return body
}
