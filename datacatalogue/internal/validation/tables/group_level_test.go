package tables

import (
	"errors"
	"math"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateGroupLevelRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateGroupLevelRequest() *tablesv1.CreateGroupLevelRequest {
	return &tablesv1.CreateGroupLevelRequest{
		ColumnId:       100,
		ParentColumnId: 101,
		Level:          12,
		Description:    "description-3",
		UserExternalId: "00000000-0000-4000-8000-000000000005",
	}
}

func validUpdateGroupLevelByIdRequest() *tablesv1.UpdateGroupLevelByIdRequest {
	return &tablesv1.UpdateGroupLevelByIdRequest{
		Id:             42,
		ColumnId:       100,
		ParentColumnId: 101,
		Level:          12,
		Description:    "description-3",
		UserExternalId: "00000000-0000-4000-8000-000000000005",
	}
}

// groupLevelFieldErrors достаёт из ошибки список полей с претензиями.
func groupLevelFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateGroupLevel(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateGroupLevelRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateGroupLevelRequest) {}},
		{name: "zero column_id", mutate: func(r *tablesv1.CreateGroupLevelRequest) { r.ColumnId = 0 }, wantField: "column_id"},
		{name: "negative column_id", mutate: func(r *tablesv1.CreateGroupLevelRequest) { r.ColumnId = -1 }, wantField: "column_id"},
		{name: "zero parent_column_id", mutate: func(r *tablesv1.CreateGroupLevelRequest) { r.ParentColumnId = 0 }, wantField: "parent_column_id"},
		{name: "negative parent_column_id", mutate: func(r *tablesv1.CreateGroupLevelRequest) { r.ParentColumnId = -1 }, wantField: "parent_column_id"},
		{name: "level over smallint", mutate: func(r *tablesv1.CreateGroupLevelRequest) { r.Level = math.MaxInt16 + 1 }, wantField: "level"},
		{name: "level under smallint", mutate: func(r *tablesv1.CreateGroupLevelRequest) { r.Level = math.MinInt16 - 1 }, wantField: "level"},
		{name: "empty description", mutate: func(r *tablesv1.CreateGroupLevelRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *tablesv1.CreateGroupLevelRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.CreateGroupLevelRequest) {
			r.Description = strings.Repeat("a", groupLevelDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateGroupLevelRequest) { r.UserExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateGroupLevelRequest) { r.UserExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "user_id without dashes", mutate: func(r *tablesv1.CreateGroupLevelRequest) { r.UserExternalId = "00000000000040008000000000000001" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateGroupLevelRequest()
			tt.mutate(req)

			err := ValidateCreateGroupLevel(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateGroupLevel() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateGroupLevel() = nil, want error on %q", tt.wantField)
			}

			fields := groupLevelFieldErrors(t, err)

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

func TestValidateUpdateGroupLevelById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateGroupLevelByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateGroupLevelByIdRequest) {}},
		{name: "zero id", mutate: func(r *tablesv1.UpdateGroupLevelByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "negative id", mutate: func(r *tablesv1.UpdateGroupLevelByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "zero column_id", mutate: func(r *tablesv1.UpdateGroupLevelByIdRequest) { r.ColumnId = 0 }, wantField: "column_id"},
		{name: "negative column_id", mutate: func(r *tablesv1.UpdateGroupLevelByIdRequest) { r.ColumnId = -1 }, wantField: "column_id"},
		{name: "zero parent_column_id", mutate: func(r *tablesv1.UpdateGroupLevelByIdRequest) { r.ParentColumnId = 0 }, wantField: "parent_column_id"},
		{name: "negative parent_column_id", mutate: func(r *tablesv1.UpdateGroupLevelByIdRequest) { r.ParentColumnId = -1 }, wantField: "parent_column_id"},
		{name: "level over smallint", mutate: func(r *tablesv1.UpdateGroupLevelByIdRequest) { r.Level = math.MaxInt16 + 1 }, wantField: "level"},
		{name: "level under smallint", mutate: func(r *tablesv1.UpdateGroupLevelByIdRequest) { r.Level = math.MinInt16 - 1 }, wantField: "level"},
		{name: "empty description", mutate: func(r *tablesv1.UpdateGroupLevelByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *tablesv1.UpdateGroupLevelByIdRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.UpdateGroupLevelByIdRequest) {
			r.Description = strings.Repeat("a", groupLevelDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateGroupLevelByIdRequest) { r.UserExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateGroupLevelByIdRequest) { r.UserExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "user_id without dashes", mutate: func(r *tablesv1.UpdateGroupLevelByIdRequest) { r.UserExternalId = "00000000000040008000000000000001" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateGroupLevelByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateGroupLevelById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateGroupLevelById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateGroupLevelById() = nil, want error on %q", tt.wantField)
			}

			fields := groupLevelFieldErrors(t, err)

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
func TestValidateCreateGroupLevelAtVarcharLimit(t *testing.T) {
	req := validCreateGroupLevelRequest()
	req.Description = strings.Repeat("a", groupLevelDescriptionMaxLen)

	if err := ValidateCreateGroupLevel(req); err != nil {
		t.Errorf("ValidateCreateGroupLevel() = %v, want nil at exactly %d chars", err, groupLevelDescriptionMaxLen)
	}
}

// Длина считается в символах, а не в байтах: кириллица занимает по 2 байта.
func TestValidateCreateGroupLevelCyrillicAtVarcharLimit(t *testing.T) {
	req := validCreateGroupLevelRequest()
	req.Description = strings.Repeat("я", groupLevelDescriptionMaxLen)

	if err := ValidateCreateGroupLevel(req); err != nil {
		t.Errorf("ValidateCreateGroupLevel() = %v, want nil at exactly %d cyrillic chars", err, groupLevelDescriptionMaxLen)
	}
}

func TestValidateCreateGroupLevelCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateGroupLevel(&tablesv1.CreateGroupLevelRequest{})

	if err == nil {
		t.Fatal("ValidateCreateGroupLevel() = nil, want errors")
	}

	fields := groupLevelFieldErrors(t, err)

	wantFields := []string{"column_id", "parent_column_id", "description", "user_id"}

	for _, field := range wantFields {
		if len(fields[field]) == 0 {
			t.Errorf("no error on %q", field)
		}
	}

	if len(fields) != len(wantFields) {
		t.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)
	}
}

func TestValidateCreateGroupLevelNil(t *testing.T) {
	if err := ValidateCreateGroupLevel(nil); err == nil {
		t.Error("ValidateCreateGroupLevel(nil) = nil, want error")
	}
}

func TestValidateUpdateGroupLevelByIdNil(t *testing.T) {
	if err := ValidateUpdateGroupLevelById(nil); err == nil {
		t.Error("ValidateUpdateGroupLevelById(nil) = nil, want error")
	}
}
