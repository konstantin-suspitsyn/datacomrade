package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// hostFieldErrors достаёт из ошибки список полей с претензиями.
func hostFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateHostRequest() *tablesv1.CreateHostRequest {
	return &tablesv1.CreateHostRequest{
		Name:        "name-0",
		Description: "description-1",
		HostEnv:     "host-env-2",
		PortEnv:     "port-env-3",
		UsernameEnv: "username-env-4",
		PasswordEnv: "password-env-5",
		ExternalId:  "00000000-0000-4000-8000-000000000007",
	}
}

func TestValidateCreateHost(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateHostRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateHostRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.CreateHostRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.CreateHostRequest) { r.Name = strings.Repeat("a", hostNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.CreateHostRequest) { r.Description = "" }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.CreateHostRequest) { r.Description = strings.Repeat("a", hostDescriptionMaxLen+1) }, wantField: "description"},
		{name: "empty host_env", mutate: func(r *tablesv1.CreateHostRequest) { r.HostEnv = "" }, wantField: "host_env"},
		{name: "host_env too long", mutate: func(r *tablesv1.CreateHostRequest) { r.HostEnv = strings.Repeat("a", hostHostEnvMaxLen+1) }, wantField: "host_env"},
		{name: "empty port_env", mutate: func(r *tablesv1.CreateHostRequest) { r.PortEnv = "" }, wantField: "port_env"},
		{name: "port_env too long", mutate: func(r *tablesv1.CreateHostRequest) { r.PortEnv = strings.Repeat("a", hostPortEnvMaxLen+1) }, wantField: "port_env"},
		{name: "empty username_env", mutate: func(r *tablesv1.CreateHostRequest) { r.UsernameEnv = "" }, wantField: "username_env"},
		{name: "username_env too long", mutate: func(r *tablesv1.CreateHostRequest) { r.UsernameEnv = strings.Repeat("a", hostUsernameEnvMaxLen+1) }, wantField: "username_env"},
		{name: "empty password_env", mutate: func(r *tablesv1.CreateHostRequest) { r.PasswordEnv = "" }, wantField: "password_env"},
		{name: "password_env too long", mutate: func(r *tablesv1.CreateHostRequest) { r.PasswordEnv = strings.Repeat("a", hostPasswordEnvMaxLen+1) }, wantField: "password_env"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateHostRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateHostRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateHostRequest()
			tt.mutate(req)

			err := ValidateCreateHost(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateHost() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateHost() = nil, want error on %q", tt.wantField)
			}

			fields := hostFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateCreateHostNil(t *testing.T) {
	if err := ValidateCreateHost(nil); err == nil {
		t.Error("ValidateCreateHost(nil) = nil, want error")
	}
}

func validValidateUpdateHostByIdRequest() *tablesv1.UpdateHostByIdRequest {
	return &tablesv1.UpdateHostByIdRequest{
		Id:          100,
		Name:        "name-1",
		Description: "description-2",
		HostEnv:     "host-env-3",
		PortEnv:     "port-env-4",
		UsernameEnv: "username-env-5",
		PasswordEnv: "password-env-6",
		ExternalId:  "00000000-0000-4000-8000-000000000008",
	}
}

func TestValidateUpdateHostById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateHostByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateHostByIdRequest) {}},
		{name: "zero id", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "empty name", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.Name = strings.Repeat("a", hostNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.Description = strings.Repeat("a", hostDescriptionMaxLen+1) }, wantField: "description"},
		{name: "empty host_env", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.HostEnv = "" }, wantField: "host_env"},
		{name: "host_env too long", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.HostEnv = strings.Repeat("a", hostHostEnvMaxLen+1) }, wantField: "host_env"},
		{name: "empty port_env", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.PortEnv = "" }, wantField: "port_env"},
		{name: "port_env too long", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.PortEnv = strings.Repeat("a", hostPortEnvMaxLen+1) }, wantField: "port_env"},
		{name: "empty username_env", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.UsernameEnv = "" }, wantField: "username_env"},
		{name: "username_env too long", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.UsernameEnv = strings.Repeat("a", hostUsernameEnvMaxLen+1) }, wantField: "username_env"},
		{name: "empty password_env", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.PasswordEnv = "" }, wantField: "password_env"},
		{name: "password_env too long", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.PasswordEnv = strings.Repeat("a", hostPasswordEnvMaxLen+1) }, wantField: "password_env"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUpdateHostByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateHostById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateHostById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateHostById() = nil, want error on %q", tt.wantField)
			}

			fields := hostFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUpdateHostByIdNil(t *testing.T) {
	if err := ValidateUpdateHostById(nil); err == nil {
		t.Error("ValidateUpdateHostById(nil) = nil, want error")
	}
}

func validValidateDeleteHostByIdRequest() *tablesv1.DeleteHostByIdRequest {
	return &tablesv1.DeleteHostByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateDeleteHostById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.DeleteHostByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.DeleteHostByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.DeleteHostByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.DeleteHostByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.DeleteHostByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateDeleteHostByIdRequest()
			tt.mutate(req)

			err := ValidateDeleteHostById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateDeleteHostById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateDeleteHostById() = nil, want error on %q", tt.wantField)
			}

			fields := hostFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateDeleteHostByIdNil(t *testing.T) {
	if err := ValidateDeleteHostById(nil); err == nil {
		t.Error("ValidateDeleteHostById(nil) = nil, want error")
	}
}

func validValidateUndeleteHostByIdRequest() *tablesv1.UndeleteHostByIdRequest {
	return &tablesv1.UndeleteHostByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateUndeleteHostById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UndeleteHostByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UndeleteHostByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.UndeleteHostByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UndeleteHostByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UndeleteHostByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUndeleteHostByIdRequest()
			tt.mutate(req)

			err := ValidateUndeleteHostById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUndeleteHostById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUndeleteHostById() = nil, want error on %q", tt.wantField)
			}

			fields := hostFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUndeleteHostByIdNil(t *testing.T) {
	if err := ValidateUndeleteHostById(nil); err == nil {
		t.Error("ValidateUndeleteHostById(nil) = nil, want error")
	}
}

func TestValidateGetHosts(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetHostsRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetHostsRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetHostsRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetHostsRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetHostsRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetHostsRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetHosts(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetHosts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetHostsNil(t *testing.T) {
	if err := ValidateGetHosts(nil); err == nil {
		t.Error("ValidateGetHosts(nil) = nil, want error")
	}
}

func TestValidateGetHostsSearchName(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetHostsSearchNameRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetHostsSearchNameRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetHostsSearchNameRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetHostsSearchNameRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetHostsSearchNameRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetHostsSearchNameRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetHostsSearchName(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetHostsSearchName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetHostsSearchNameNil(t *testing.T) {
	if err := ValidateGetHostsSearchName(nil); err == nil {
		t.Error("ValidateGetHostsSearchName(nil) = nil, want error")
	}
}

func TestValidateGetHostDeleted(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetHostDeletedRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetHostDeletedRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetHostDeletedRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetHostDeletedRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetHostDeletedRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetHostDeletedRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetHostDeleted(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetHostDeleted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetHostDeletedNil(t *testing.T) {
	if err := ValidateGetHostDeleted(nil); err == nil {
		t.Error("ValidateGetHostDeleted(nil) = nil, want error")
	}
}
