package tables

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
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

			if got == nil {
				t.Fatal("HostsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestToCreateHostParams(t *testing.T) {
	req := &tablesv1.CreateHostRequest{
		Name:        "name-0",
		Description: "description-0",
		HostEnv:     "host-env-0",
		PortEnv:     "port-env-0",
		UsernameEnv: "username-env-0",
		PasswordEnv: "password-env-0",
		ExternalId:  "00000000-0000-4000-8000-000000000007",
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
		Id:          100,
		Name:        "name-0",
		Description: "description-0",
		HostEnv:     "host-env-0",
		PortEnv:     "port-env-0",
		UsernameEnv: "username-env-0",
		PasswordEnv: "password-env-0",
		ExternalId:  "00000000-0000-4000-8000-000000000008",
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
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUpdateHostByIdParams(nil); got != (tables_model.UpdateHostByIdParams{}) {
		t.Errorf("ToUpdateHostByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToDeleteHostByIdParams(t *testing.T) {
	req := &tablesv1.DeleteHostByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.DeleteHostByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToDeleteHostByIdParams(req); got != want {
		t.Errorf("ToDeleteHostByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToDeleteHostByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToDeleteHostByIdParams(nil); got != (tables_model.DeleteHostByIdParams{}) {
		t.Errorf("ToDeleteHostByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToUndeleteHostByIdParams(t *testing.T) {
	req := &tablesv1.UndeleteHostByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.UndeleteHostByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToUndeleteHostByIdParams(req); got != want {
		t.Errorf("ToUndeleteHostByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUndeleteHostByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUndeleteHostByIdParams(nil); got != (tables_model.UndeleteHostByIdParams{}) {
		t.Errorf("ToUndeleteHostByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestGetHostsDefaultsPageLimit(t *testing.T) {
	got := ToGetHostsParams(&tablesv1.GetHostsRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetHostsDefaultsPage(t *testing.T) {
	got := ToGetHostsParams(&tablesv1.GetHostsRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetHostsKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetHostsParams(&tablesv1.GetHostsRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetHostsSearchNameDefaultsPageLimit(t *testing.T) {
	got := ToGetHostsSearchNameParams(&tablesv1.GetHostsSearchNameRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetHostsSearchNameDefaultsPage(t *testing.T) {
	got := ToGetHostsSearchNameParams(&tablesv1.GetHostsSearchNameRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetHostsSearchNameKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetHostsSearchNameParams(&tablesv1.GetHostsSearchNameRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetHostsSearchNamePassesFilterFields(t *testing.T) {
	req := &tablesv1.GetHostsSearchNameRequest{PageLimit: 10, Page: 1,
		SearchName: "SearchName-0",
	}

	got := ToGetHostsSearchNameParams(req)

	if got.SearchName != "SearchName-0" {
		t.Errorf("SearchName = %q, want %q", got.SearchName, "SearchName-0")
	}
}

func TestGetHostDeletedDefaultsPageLimit(t *testing.T) {
	got := ToGetHostDeletedParams(&tablesv1.GetHostDeletedRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetHostDeletedDefaultsPage(t *testing.T) {
	got := ToGetHostDeletedParams(&tablesv1.GetHostDeletedRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetHostDeletedKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetHostDeletedParams(&tablesv1.GetHostDeletedRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}
