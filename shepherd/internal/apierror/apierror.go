// Package apierror переводит ошибки в единый JSON-конверт HTTP-ответа
// и переносит коды gRPC от бэкендов на соответствующие HTTP-статусы.
package apierror

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Body — тело JSON-ответа об ошибке, единое для всего REST-периметра Shepherd.
type Body struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Write пишет JSON-ошибку с заданным HTTP-статусом.
func Write(w http.ResponseWriter, httpStatus int, code, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(Body{Code: code, Message: message})
}

// WriteGRPC подбирает HTTP-статус по коду ошибки gRPC от бэкенда и пишет JSON-ответ.
func WriteGRPC(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		Write(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}

	httpStatus, code := mapCode(st.Code())
	Write(w, httpStatus, code, st.Message())
}

func mapCode(c codes.Code) (int, string) {
	switch c {
	case codes.InvalidArgument:
		return http.StatusBadRequest, "invalid_argument"
	case codes.NotFound:
		return http.StatusNotFound, "not_found"
	case codes.AlreadyExists:
		return http.StatusConflict, "already_exists"
	case codes.PermissionDenied:
		return http.StatusForbidden, "permission_denied"
	case codes.Unauthenticated:
		return http.StatusUnauthorized, "unauthenticated"
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout, "deadline_exceeded"
	case codes.Unavailable:
		return http.StatusServiceUnavailable, "unavailable"
	default:
		return http.StatusInternalServerError, "internal"
	}
}
