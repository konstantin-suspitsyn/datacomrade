package tables

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

var (
	hostCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	hostUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testHostRow — строка dc.host со значениями, различимыми между полями.
func testHostRow() tables_model.DcHost {
	return tables_model.DcHost{
		ID:          100,
		Name:        "name-0",
		Description: "description-0",
		HostEnv:     "host-env-0",
		PortEnv:     "port-env-0",
		UsernameEnv: "username-env-0",
		PasswordEnv: "password-env-0",
		IsDeleted:   false,
		CreatedAt:   hostCreatedAt,
		UpdatedAt:   hostUpdatedAt,
		UserID:      110,
	}
}

func TestHostToProto(t *testing.T) {
	row := testHostRow()
	got := HostToProto(row)

	if got == nil {
		t.Fatal("HostToProto() = nil, want value")
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

	if got.GetHostEnv() != row.HostEnv {
		t.Errorf("HostEnv = %q, want %q", got.GetHostEnv(), row.HostEnv)
	}

	if got.GetPortEnv() != row.PortEnv {
		t.Errorf("PortEnv = %q, want %q", got.GetPortEnv(), row.PortEnv)
	}

	if got.GetUsernameEnv() != row.UsernameEnv {
		t.Errorf("UsernameEnv = %q, want %q", got.GetUsernameEnv(), row.UsernameEnv)
	}

	if got.GetPasswordEnv() != row.PasswordEnv {
		t.Errorf("PasswordEnv = %q, want %q", got.GetPasswordEnv(), row.PasswordEnv)
	}

	if got.GetIsDeleted() != row.IsDeleted {
		t.Errorf("IsDeleted = %v, want %v", got.GetIsDeleted(), row.IsDeleted)
	}

	if !got.GetCreatedAt().AsTime().Equal(row.CreatedAt) {
		t.Errorf("CreatedAt = %v, want %v", got.GetCreatedAt().AsTime(), row.CreatedAt)
	}

	if !got.GetUpdatedAt().AsTime().Equal(row.UpdatedAt) {
		t.Errorf("UpdatedAt = %v, want %v", got.GetUpdatedAt().AsTime(), row.UpdatedAt)
	}

	if got.GetUserId() != row.UserID {
		t.Errorf("UserId = %d, want %d", got.GetUserId(), row.UserID)
	}

}

func TestHostToProtoDeleted(t *testing.T) {
	row := testHostRow()
	row.IsDeleted = true

	if got := HostToProto(row); !got.GetIsDeleted() {
		t.Error("IsDeleted = false, want true")
	}
}

func TestHostsToProto(t *testing.T) {
	first := testHostRow()

	second := testHostRow()
	second.ID = 999
	second.Name = "second-value"

	tests := []struct {
		name    string
		input   []tables_model.DcHost
		wantLen int
	}{
		{name: "two rows", input: []tables_model.DcHost{first, second}, wantLen: 2},
		{name: "empty slice", input: []tables_model.DcHost{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HostsToProto(tt.input)

			// Пустой вход даёт пустой, а не nil-слайс.
			if got == nil {
				t.Fatal("HostsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestHostsToProtoKeepsOrder(t *testing.T) {
	first := testHostRow()
	second := testHostRow()
	second.Name = "second-value"

	got := HostsToProto([]tables_model.DcHost{first, second})

	if got[0].GetName() != first.Name {
		t.Errorf("[0] = %q, want %q", got[0].GetName(), first.Name)
	}

	if got[1].GetName() != second.Name {
		t.Errorf("[1] = %q, want %q", got[1].GetName(), second.Name)
	}
}

func TestToCreateHostParams(t *testing.T) {
	req := &tablesv1.CreateHostRequest{
		Name:           "name-0",
		Description:    "description-0",
		HostEnv:        "host-env-0",
		PortEnv:        "port-env-0",
		UsernameEnv:    "username-env-0",
		PasswordEnv:    "password-env-0",
		UserExternalId: "00000000-0000-4000-8000-000000000007",
	}

	want := tables_model.CreateHostParams{
		Name:        "name-0",
		Description: "description-0",
		HostEnv:     "host-env-0",
		PortEnv:     "port-env-0",
		UsernameEnv: "username-env-0",
		PasswordEnv: "password-env-0",
		ExternalID:  uuid.MustParse("00000000-0000-4000-8000-000000000007"),
	}

	if got := ToCreateHostParams(req); got != want {
		t.Errorf("ToCreateHostParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateHostParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateHostParams(nil); got != (tables_model.CreateHostParams{}) {
		t.Errorf("ToCreateHostParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateHostByIdParams(t *testing.T) {
	req := &tablesv1.UpdateHostByIdRequest{
		Id:             100,
		Name:           "name-0",
		Description:    "description-0",
		HostEnv:        "host-env-0",
		PortEnv:        "port-env-0",
		UsernameEnv:    "username-env-0",
		PasswordEnv:    "password-env-0",
		UserExternalId: "00000000-0000-4000-8000-000000000008",
	}

	want := tables_model.UpdateHostByIdParams{
		ID:          100,
		Name:        "name-0",
		Description: "description-0",
		HostEnv:     "host-env-0",
		PortEnv:     "port-env-0",
		UsernameEnv: "username-env-0",
		PasswordEnv: "password-env-0",
		ExternalID:  uuid.MustParse("00000000-0000-4000-8000-000000000008"),
	}

	if got := ToUpdateHostByIdParams(req); got != want {
		t.Errorf("ToUpdateHostByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateHostByIdParamsNil(t *testing.T) {
	if got := ToUpdateHostByIdParams(nil); got != (tables_model.UpdateHostByIdParams{}) {
		t.Errorf("ToUpdateHostByIdParams(nil) = %+v, want zero value", got)
	}
}
