package validation

import (
	"errors"
	"math"
	"testing"

	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

func TestValidateID(t *testing.T) {
	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{name: "positive", id: 1},
		{name: "large", id: math.MaxInt64},
		{name: "zero", id: 0, wantErr: true},
		{name: "negative", id: -1, wantErr: true},
		{name: "min int64", id: math.MinInt64, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateID(tt.id)

			if !tt.wantErr {
				if err != nil {
					t.Fatalf("ValidateID(%d) = %v, want nil", tt.id, err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateID(%d) = nil, want error", tt.id)
			}

			var validationErr *validator.ValidationError
			if !errors.As(err, &validationErr) {
				t.Fatalf("error = %v, want *validator.ValidationError", err)
			}

			if len(validationErr.Errors["id"]) == 0 {
				t.Errorf("no error on \"id\": %v", validationErr.Errors)
			}
		})
	}
}
