package userdomainroles

import (
	"testing"
	"time"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
)

var (
	domainsDomainRoleCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	domainsDomainRoleUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testDomainsDomainRoleRow — строка dc.domains_domain_roles со значениями, различимыми между полями.
func testDomainsDomainRoleRow() user_domain_roles.DcDomainsDomainRole {
	return user_domain_roles.DcDomainsDomainRole{
		ID:            100,
		DomainCatID:   101,
		DomainRolesID: 102,
		CreatedAt:     domainsDomainRoleCreatedAt,
		UpdatedAt:     domainsDomainRoleUpdatedAt,
		IsDeleted:     false,
	}
}

func TestDomainsDomainRoleToProto(t *testing.T) {
	row := testDomainsDomainRoleRow()
	got := DomainsDomainRoleToProto(row)

	if got == nil {
		t.Fatal("DomainsDomainRoleToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetDomainCatId() != row.DomainCatID {
		t.Errorf("DomainCatId = %d, want %d", got.GetDomainCatId(), row.DomainCatID)
	}

	if got.GetDomainRolesId() != row.DomainRolesID {
		t.Errorf("DomainRolesId = %d, want %d", got.GetDomainRolesId(), row.DomainRolesID)
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

func TestDomainsDomainRoleToProtoDeleted(t *testing.T) {
	row := testDomainsDomainRoleRow()
	row.IsDeleted = true

	if got := DomainsDomainRoleToProto(row); !got.GetIsDeleted() {
		t.Error("IsDeleted = false, want true")
	}
}

func TestDomainsDomainRolesToProto(t *testing.T) {
	first := testDomainsDomainRoleRow()

	second := testDomainsDomainRoleRow()
	second.ID = 999

	tests := []struct {
		name    string
		input   []user_domain_roles.DcDomainsDomainRole
		wantLen int
	}{
		{name: "two rows", input: []user_domain_roles.DcDomainsDomainRole{first, second}, wantLen: 2},
		{name: "empty slice", input: []user_domain_roles.DcDomainsDomainRole{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DomainsDomainRolesToProto(tt.input)

			// Пустой вход даёт пустой, а не nil-слайс.
			if got == nil {
				t.Fatal("DomainsDomainRolesToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestDomainsDomainRolesToProtoKeepsOrder(t *testing.T) {
	first := testDomainsDomainRoleRow()
	second := testDomainsDomainRoleRow()
	second.ID = 999

	got := DomainsDomainRolesToProto([]user_domain_roles.DcDomainsDomainRole{first, second})

	if got[0].GetId() != first.ID {
		t.Errorf("[0] = %d, want %d", got[0].GetId(), first.ID)
	}

	if got[1].GetId() != second.ID {
		t.Errorf("[1] = %d, want %d", got[1].GetId(), second.ID)
	}
}

func TestToCreateDomainsDomainRoleParams(t *testing.T) {
	req := &userdomainrolesv1.CreateDomainsDomainRoleRequest{
		DomainCatId:   100,
		DomainRolesId: 101,
	}

	want := user_domain_roles.CreateDomainsDomainRoleParams{
		DomainCatID:   100,
		DomainRolesID: 101,
	}

	if got := ToCreateDomainsDomainRoleParams(req); got != want {
		t.Errorf("ToCreateDomainsDomainRoleParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateDomainsDomainRoleParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateDomainsDomainRoleParams(nil); got != (user_domain_roles.CreateDomainsDomainRoleParams{}) {
		t.Errorf("ToCreateDomainsDomainRoleParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateDomainsDomainRoleByIdParams(t *testing.T) {
	req := &userdomainrolesv1.UpdateDomainsDomainRoleByIdRequest{
		Id:            100,
		DomainCatId:   101,
		DomainRolesId: 102,
	}

	want := user_domain_roles.UpdateDomainsDomainRoleByIdParams{
		ID:            100,
		DomainCatID:   101,
		DomainRolesID: 102,
	}

	if got := ToUpdateDomainsDomainRoleByIdParams(req); got != want {
		t.Errorf("ToUpdateDomainsDomainRoleByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateDomainsDomainRoleByIdParamsNil(t *testing.T) {
	if got := ToUpdateDomainsDomainRoleByIdParams(nil); got != (user_domain_roles.UpdateDomainsDomainRoleByIdParams{}) {
		t.Errorf("ToUpdateDomainsDomainRoleByIdParams(nil) = %+v, want zero value", got)
	}
}
