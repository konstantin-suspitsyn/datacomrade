package authlogic

import (
	"testing"

	"github.com/google/uuid"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/auth_logic"
	authlogicv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/auth_logic/v1"
)

func TestToGetTableIdsByExternalUserIdAndRolesParams(t *testing.T) {
	req := &authlogicv1.GetTableIdsByExternalUserIdAndRolesRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Name:       "name-1",
	}

	want := auth_logic.GetTableIdsByExternalUserIdAndRolesParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		Name:       "name-1",
	}

	if got := ToGetTableIdsByExternalUserIdAndRolesParams(req); got != want {
		t.Errorf("ToGetTableIdsByExternalUserIdAndRolesParams() = %+v, want %+v", got, want)
	}
}

func TestToGetTableIdsByExternalUserIdAndRolesParamsNil(t *testing.T) {
	if got := ToGetTableIdsByExternalUserIdAndRolesParams(nil); got != (auth_logic.GetTableIdsByExternalUserIdAndRolesParams{}) {
		t.Errorf("ToGetTableIdsByExternalUserIdAndRolesParams(nil) = %+v, want zero value", got)
	}
}

func TestGetTableIdsByExternalUserIdAndRolesToProto(t *testing.T) {
	got := GetTableIdsByExternalUserIdAndRolesToProto([]int64{1, 2, 3})

	if got == nil {
		t.Fatal("GetTableIdsByExternalUserIdAndRolesToProto() = nil, want value")
	}

	if len(got.GetTableIds()) != 3 {
		t.Errorf("len(TableIds) = %d, want 3", len(got.GetTableIds()))
	}
}

func TestToGetTableIdsByUserIdAndRolesParams(t *testing.T) {
	req := &authlogicv1.GetTableIdsByUserIdAndRolesRequest{
		UserId: 100,
		Name:   "name-1",
	}

	want := auth_logic.GetTableIdsByUserIdAndRolesParams{
		UserID: 100,
		Name:   "name-1",
	}

	if got := ToGetTableIdsByUserIdAndRolesParams(req); got != want {
		t.Errorf("ToGetTableIdsByUserIdAndRolesParams() = %+v, want %+v", got, want)
	}
}

func TestToGetTableIdsByUserIdAndRolesParamsNil(t *testing.T) {
	if got := ToGetTableIdsByUserIdAndRolesParams(nil); got != (auth_logic.GetTableIdsByUserIdAndRolesParams{}) {
		t.Errorf("ToGetTableIdsByUserIdAndRolesParams(nil) = %+v, want zero value", got)
	}
}

func TestGetTableIdsByUserIdAndRolesToProto(t *testing.T) {
	got := GetTableIdsByUserIdAndRolesToProto([]int64{1, 2, 3})

	if got == nil {
		t.Fatal("GetTableIdsByUserIdAndRolesToProto() = nil, want value")
	}

	if len(got.GetTableIds()) != 3 {
		t.Errorf("len(TableIds) = %d, want 3", len(got.GetTableIds()))
	}
}
