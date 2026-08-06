package userdomainroles

import (
	"testing"
	"time"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
)

var (
	domainRoleCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	domainRoleUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testDomainRoleRow — строка dc.domain_roles со значениями, различимыми между полями.
func testDomainRoleRow() user_domain_roles.DcDomainRole {
	return user_domain_roles.DcDomainRole{
		ID:          100,
		Name:        "name-0",
		Description: "description-0",
		CreatedAt:   domainRoleCreatedAt,
		UpdatedAt:   domainRoleUpdatedAt,
		IsDeleted:   false,
	}
}

func TestDomainRoleToProto(t *testing.T) {
	row := testDomainRoleRow()
	got := DomainRoleToProto(row)

	if got == nil {
		t.Fatal("DomainRoleToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetName() != row.Name {
		t.Errorf("Name = %q, want %q", got.GetName(), row.Name)
	}

	if got.GetDescription() != row.Description {
		t.Errorf("Description = %q, want %q", got.GetDescription(), row.Description)
	}

	if !got.GetCreatedAt().AsTime().Equal(row.CreatedAt) {
		t.Errorf("CreatedAt = %v, want %v", got.GetCreatedAt().AsTime(), row.CreatedAt)
	}

	if !got.GetUpdatedAt().AsTime().Equal(row.UpdatedAt) {
		t.Errorf("UpdatedAt = %v, want %v", got.GetUpdatedAt().AsTime(), row.UpdatedAt)
	}

	if got.GetIsDeleted() != row.IsDeleted {
		t.Errorf("IsDeleted = %v, want %v", got.GetIsDeleted(), row.IsDeleted)
	}

}

func TestDomainRolesToProto(t *testing.T) {
	first := testDomainRoleRow()

	second := testDomainRoleRow()
	second.ID = 999
	second.Name = "second-value"

	tests := []struct {
		name    string
		input   []user_domain_roles.DcDomainRole
		wantLen int
	}{
		{name: "two rows", input: []user_domain_roles.DcDomainRole{first, second}, wantLen: 2},
		{name: "empty slice", input: []user_domain_roles.DcDomainRole{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DomainRolesToProto(tt.input)

			if got == nil {
				t.Fatal("DomainRolesToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestToCreateDomainRoleParams(t *testing.T) {
	req := &userdomainrolesv1.CreateDomainRoleRequest{
		Name:        "name-0",
		Description: "description-0",
	}

	want := user_domain_roles.CreateDomainRoleParams{
		Name:        "name-0",
		Description: "description-0",
	}

	if got := ToCreateDomainRoleParams(req); got != want {
		t.Errorf("ToCreateDomainRoleParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateDomainRoleParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateDomainRoleParams(nil); got != (user_domain_roles.CreateDomainRoleParams{}) {
		t.Errorf("ToCreateDomainRoleParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateDomainRoleByIdParams(t *testing.T) {
	req := &userdomainrolesv1.UpdateDomainRoleByIdRequest{
		Id:          100,
		Name:        "name-0",
		Description: "description-0",
	}

	want := user_domain_roles.UpdateDomainRoleByIdParams{
		ID:          100,
		Name:        "name-0",
		Description: "description-0",
	}

	if got := ToUpdateDomainRoleByIdParams(req); got != want {
		t.Errorf("ToUpdateDomainRoleByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateDomainRoleByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUpdateDomainRoleByIdParams(nil); got != (user_domain_roles.UpdateDomainRoleByIdParams{}) {
		t.Errorf("ToUpdateDomainRoleByIdParams(nil) = %+v, want zero value", got)
	}
}
