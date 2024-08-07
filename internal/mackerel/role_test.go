package mackerel

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

func Test_Role_ReadRole(t *testing.T) {
	t.Parallel()

	defaultClient := func(service string) ([]*mackerel.Role, error) {
		if service != "service0" {
			return nil, fmt.Errorf("service not found")
		}
		return []*mackerel.Role{
			{Name: "role0", Memo: "memo"},
			{Name: "role1"},
		}, nil
	}

	cases := map[string]struct {
		inService string
		inRole    string
		inClient  roleFinderFunc

		wants   RoleModel
		wantErr bool
	}{
		"valid": {
			inService: "service0",
			inRole:    "role0",
			inClient:  defaultClient,

			wants: RoleModel{
				ID:          types.StringValue("service0:role0"),
				ServiceName: types.StringValue("service0"),
				RoleName:    types.StringValue("role0"),
				Memo:        types.StringValue("memo"),
			},
		},
		"no role": {
			inService: "service0",
			inRole:    "role-not-exists",
			inClient:  defaultClient,

			wantErr: true,
		},
	}

	ctx := context.Background()
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			data, err := readRoleInner(ctx, tt.inClient, tt.inService, tt.inRole)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("unexpected error: %+v", err)
				}
				return
			} else if tt.wantErr {
				t.Error("expected error, but got no error")
				return
			}

			if diff := cmp.Diff(data, tt.wants); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func Test_RoleModel_Read(t *testing.T) {
	t.Parallel()

	defaultClient := roleFinderFunc(func(service string) ([]*mackerel.Role, error) {
		if service != "service0" {
			return nil, fmt.Errorf("service not found")
		}
		return []*mackerel.Role{
			{Name: "role0", Memo: "memo"},
			{Name: "role1"},
		}, nil
	})

	cases := map[string]struct {
		in       RoleModel
		inClient roleFinderFunc

		wants   RoleModel
		wantErr bool
	}{
		"from id": {
			in:       RoleModel{ID: types.StringValue("service0:role0")},
			inClient: defaultClient,

			wants: RoleModel{
				ID:          types.StringValue("service0:role0"),
				ServiceName: types.StringValue("service0"),
				RoleName:    types.StringValue("role0"),
				Memo:        types.StringValue("memo"),
			},
		},
		"invalid id": {
			in:       RoleModel{ID: types.StringValue("invalid id")},
			inClient: defaultClient,

			wantErr: true,
		},
		"from service and name": {
			in: RoleModel{
				ServiceName: types.StringValue("service0"),
				RoleName:    types.StringValue("role0"),
			},
			inClient: defaultClient,

			wants: RoleModel{
				ID:          types.StringValue("service0:role0"),
				ServiceName: types.StringValue("service0"),
				RoleName:    types.StringValue("role0"),
				Memo:        types.StringValue("memo"),
			},
		},
	}

	ctx := context.Background()
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := tt.in
			if err := m.readInner(ctx, tt.inClient); err != nil {
				if !tt.wantErr {
					t.Errorf("unexpected error: %+v", err)
				}
				return
			} else if tt.wantErr {
				t.Error("expected error, but got no error")
				return
			}

			if diff := cmp.Diff(m, tt.wants); diff != "" {
				t.Error(diff)
			}
		})
	}
}

type roleFinderFunc func(string) ([]*mackerel.Role, error)

func (f roleFinderFunc) FindRoles(serviceName string) ([]*mackerel.Role, error) {
	return f(serviceName)
}
