package tables

import (
	"errors"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// calculationTypeFieldErrors достаёт из ошибки список полей с претензиями.
func calculationTypeFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateGetCalculationTypes(t *testing.T) {
	tests := []struct {
		name    string
		order   string
		wantErr bool
	}{
		{name: "empty order", order: ""},
		{name: "ASC", order: "ASC"},
		{name: "DESC", order: "DESC"},
		{name: "invalid order", order: "sideways", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetCalculationTypes(&tablesv1.GetCalculationTypesRequest{Order: tt.order})
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetCalculationTypes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetCalculationTypesNil(t *testing.T) {
	if err := ValidateGetCalculationTypes(nil); err == nil {
		t.Error("ValidateGetCalculationTypes(nil) = nil, want error")
	}
}
