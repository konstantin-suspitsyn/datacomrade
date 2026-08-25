package olapexchange

import (
	"errors"
	"testing"

	olapexchangev1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/olap_exchange/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// authLogicFieldErrors достаёт из ошибки список полей с претензиями.
func authLogicFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validGetOlapAliasesRequest() *olapexchangev1.GetOlapAliasesRequest {
	return &olapexchangev1.GetOlapAliasesRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Order:      "ASC",
	}
}

func TestValidateGetOlapAliases(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*olapexchangev1.GetOlapAliasesRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*olapexchangev1.GetOlapAliasesRequest) {}},
		{name: "empty external_id", mutate: func(r *olapexchangev1.GetOlapAliasesRequest) { r.ExternalId = "" }, wantField: "external_id"},
		{name: "malformed external_id", mutate: func(r *olapexchangev1.GetOlapAliasesRequest) { r.ExternalId = "not-a-uuid" }, wantField: "external_id"},
		{name: "invalid order", mutate: func(r *olapexchangev1.GetOlapAliasesRequest) { r.Order = "sideways" }, wantField: "order"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validGetOlapAliasesRequest()
			tt.mutate(req)

			err := ValidateGetOlapAliases(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateGetOlapAliases() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateGetOlapAliases() = nil, want error on %q", tt.wantField)
			}

			fields := authLogicFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateGetOlapAliasesNil(t *testing.T) {
	if err := ValidateGetOlapAliases(nil); err == nil {
		t.Error("ValidateGetOlapAliases(nil) = nil, want error")
	}
}
