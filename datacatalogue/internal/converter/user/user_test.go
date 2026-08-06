package user

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_model"
	userv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user/v1"
)

var (
	userCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	userUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testUserRow — строка dc.user со значениями, различимыми между полями.
func testUserRow() user_model.DcUser {
	return user_model.DcUser{
		ID:         100,
		Name:       "name-0",
		CreatedAt:  userCreatedAt,
		UpdatedAt:  userUpdatedAt,
		IsDeleted:  false,
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000006"),
	}
}

func TestUserToProto(t *testing.T) {
	row := testUserRow()
	got := UserToProto(row)

	if got == nil {
		t.Fatal("UserToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetName() != row.Name {
		t.Errorf("Name = %q, want %q", got.GetName(), row.Name)
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

	if got.GetExternalId() != row.ExternalID.String() {
		t.Errorf("ExternalId = %q, want %q", got.GetExternalId(), row.ExternalID.String())
	}

}

func TestUsersToProto(t *testing.T) {
	first := testUserRow()

	second := testUserRow()
	second.ID = 999
	second.Name = "second-value"

	tests := []struct {
		name    string
		input   []user_model.DcUser
		wantLen int
	}{
		{name: "two rows", input: []user_model.DcUser{first, second}, wantLen: 2},
		{name: "empty slice", input: []user_model.DcUser{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UsersToProto(tt.input)

			if got == nil {
				t.Fatal("UsersToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestToCreateUserParams(t *testing.T) {
	req := &userv1.CreateUserRequest{
		Name:       "name-0",
		ExternalId: "00000000-0000-4000-8000-000000000002",
	}

	want := user_model.CreateUserParams{
		Name:       "name-0",
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000002"),
	}

	if got := ToCreateUserParams(req); got != want {
		t.Errorf("ToCreateUserParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateUserParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateUserParams(nil); got != (user_model.CreateUserParams{}) {
		t.Errorf("ToCreateUserParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateUserByIdParams(t *testing.T) {
	req := &userv1.UpdateUserByIdRequest{
		Id:         100,
		Name:       "name-0",
		ExternalId: "00000000-0000-4000-8000-000000000003",
	}

	want := user_model.UpdateUserByIdParams{
		ID:         100,
		Name:       "name-0",
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000003"),
	}

	if got := ToUpdateUserByIdParams(req); got != want {
		t.Errorf("ToUpdateUserByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateUserByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUpdateUserByIdParams(nil); got != (user_model.UpdateUserByIdParams{}) {
		t.Errorf("ToUpdateUserByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToGetUserByExternalIdArg(t *testing.T) {
	req := &userv1.GetUserByExternalIdRequest{ExternalId: "00000000-0000-4000-8000-000000007001"}

	if got := ToGetUserByExternalIdArg(req); got != uuid.MustParse("00000000-0000-4000-8000-000000007001") {
		t.Errorf("ToGetUserByExternalIdArg() = %v, want %v", got, uuid.MustParse("00000000-0000-4000-8000-000000007001"))
	}
}

func TestToGetUserByExternalIdArgNil(t *testing.T) {
	if got := ToGetUserByExternalIdArg(nil); got != uuid.Nil {
		t.Errorf("ToGetUserByExternalIdArg(nil) = %v, want zero value", got)
	}
}
