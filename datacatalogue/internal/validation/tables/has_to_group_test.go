package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateHasToGroupRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateHasToGroupRequest() *tablesv1.CreateHasToGroupRequest {
	return &tablesv1.CreateHasToGroupRequest{
		ColumnIdA:   100,
		ColumnIdB:   101,
		Description: "description-2",
		UserId:      103,
	}
}

func validUpdateHasToGroupByIdRequest() *tablesv1.UpdateHasToGroupByIdRequest {
	return &tablesv1.UpdateHasToGroupByIdRequest{
		Id:          42,
		ColumnIdA:   100,
		ColumnIdB:   101,
		Description: "description-2",
		UserId:      103,
	}
}

// hasToGroupFieldErrors достаёт из ошибки список полей с претензиями.
func hasToGroupFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateHasToGroup(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateHasToGroupRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateHasToGroupRequest) {}},
		{name: "zero column_id_a", mutate: func(r *tablesv1.CreateHasToGroupRequest) { r.ColumnIdA = 0 }, wantField: "column_id_a"},
		{name: "negative column_id_a", mutate: func(r *tablesv1.CreateHasToGroupRequest) { r.ColumnIdA = -1 }, wantField: "column_id_a"},
		{name: "zero column_id_b", mutate: func(r *tablesv1.CreateHasToGroupRequest) { r.ColumnIdB = 0 }, wantField: "column_id_b"},
		{name: "negative column_id_b", mutate: func(r *tablesv1.CreateHasToGroupRequest) { r.ColumnIdB = -1 }, wantField: "column_id_b"},
		{name: "empty description", mutate: func(r *tablesv1.CreateHasToGroupRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *tablesv1.CreateHasToGroupRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.CreateHasToGroupRequest) {
			r.Description = strings.Repeat("a", hasToGroupDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "zero user_id", mutate: func(r *tablesv1.CreateHasToGroupRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "negative user_id", mutate: func(r *tablesv1.CreateHasToGroupRequest) { r.UserId = -1 }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateHasToGroupRequest()
			tt.mutate(req)

			err := ValidateCreateHasToGroup(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateHasToGroup() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateHasToGroup() = nil, want error on %q", tt.wantField)
			}

			fields := hasToGroupFieldErrors(t, err)

			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}

			// Порча одного поля не должна задевать остальные.
			if len(fields) != 1 {
				t.Errorf("errors on %d fields, want only %q: %v", len(fields), tt.wantField, fields)
			}
		})
	}
}

func TestValidateUpdateHasToGroupById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateHasToGroupByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateHasToGroupByIdRequest) {}},
		{name: "zero id", mutate: func(r *tablesv1.UpdateHasToGroupByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "negative id", mutate: func(r *tablesv1.UpdateHasToGroupByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "zero column_id_a", mutate: func(r *tablesv1.UpdateHasToGroupByIdRequest) { r.ColumnIdA = 0 }, wantField: "column_id_a"},
		{name: "negative column_id_a", mutate: func(r *tablesv1.UpdateHasToGroupByIdRequest) { r.ColumnIdA = -1 }, wantField: "column_id_a"},
		{name: "zero column_id_b", mutate: func(r *tablesv1.UpdateHasToGroupByIdRequest) { r.ColumnIdB = 0 }, wantField: "column_id_b"},
		{name: "negative column_id_b", mutate: func(r *tablesv1.UpdateHasToGroupByIdRequest) { r.ColumnIdB = -1 }, wantField: "column_id_b"},
		{name: "empty description", mutate: func(r *tablesv1.UpdateHasToGroupByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *tablesv1.UpdateHasToGroupByIdRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.UpdateHasToGroupByIdRequest) {
			r.Description = strings.Repeat("a", hasToGroupDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "zero user_id", mutate: func(r *tablesv1.UpdateHasToGroupByIdRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "negative user_id", mutate: func(r *tablesv1.UpdateHasToGroupByIdRequest) { r.UserId = -1 }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateHasToGroupByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateHasToGroupById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateHasToGroupById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateHasToGroupById() = nil, want error on %q", tt.wantField)
			}

			fields := hasToGroupFieldErrors(t, err)

			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}

			// Порча одного поля не должна задевать остальные.
			if len(fields) != 1 {
				t.Errorf("errors on %d fields, want only %q: %v", len(fields), tt.wantField, fields)
			}
		})
	}
}

// Ровно граничная длина проходит: varchar(n) допускает n символов.
func TestValidateCreateHasToGroupAtVarcharLimit(t *testing.T) {
	req := validCreateHasToGroupRequest()
	req.Description = strings.Repeat("a", hasToGroupDescriptionMaxLen)

	if err := ValidateCreateHasToGroup(req); err != nil {
		t.Errorf("ValidateCreateHasToGroup() = %v, want nil at exactly %d chars", err, hasToGroupDescriptionMaxLen)
	}
}

// Длина считается в символах, а не в байтах: кириллица занимает по 2 байта.
func TestValidateCreateHasToGroupCyrillicAtVarcharLimit(t *testing.T) {
	req := validCreateHasToGroupRequest()
	req.Description = strings.Repeat("я", hasToGroupDescriptionMaxLen)

	if err := ValidateCreateHasToGroup(req); err != nil {
		t.Errorf("ValidateCreateHasToGroup() = %v, want nil at exactly %d cyrillic chars", err, hasToGroupDescriptionMaxLen)
	}
}

func TestValidateCreateHasToGroupCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateHasToGroup(&tablesv1.CreateHasToGroupRequest{})

	if err == nil {
		t.Fatal("ValidateCreateHasToGroup() = nil, want errors")
	}

	fields := hasToGroupFieldErrors(t, err)

	wantFields := []string{"column_id_a", "column_id_b", "description", "user_id"}

	for _, field := range wantFields {
		if len(fields[field]) == 0 {
			t.Errorf("no error on %q", field)
		}
	}

	if len(fields) != len(wantFields) {
		t.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)
	}
}

func TestValidateCreateHasToGroupNil(t *testing.T) {
	if err := ValidateCreateHasToGroup(nil); err == nil {
		t.Error("ValidateCreateHasToGroup(nil) = nil, want error")
	}
}

func TestValidateUpdateHasToGroupByIdNil(t *testing.T) {
	if err := ValidateUpdateHasToGroupById(nil); err == nil {
		t.Error("ValidateUpdateHasToGroupById(nil) = nil, want error")
	}
}
