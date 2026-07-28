// Package apierror переводит ошибки внутренних слоёв в статусы gRPC.
//
// Пакет общий для всех трёх api-пакетов, чтобы одинаковая ошибка не
// превращалась в разные коды в разных сервисах.
package apierror

import (
	"errors"

	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Wrap подбирает код gRPC по типу ошибки:
//
//	*validator.ValidationError → InvalidArgument
//	customerrors.ErrNotFound   → NotFound
//	всё остальное              → Internal
//
// Для nil возвращается nil, поэтому вызов можно ставить безусловно.
func Wrap(err error) error {
	if err == nil {
		return nil
	}

	var validationErr *validator.ValidationError
	if errors.As(err, &validationErr) {
		return status.Error(codes.InvalidArgument, validationErr.Error())
	}

	if errors.Is(err, customerrors.ErrNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}

	return status.Error(codes.Internal, err.Error())
}
