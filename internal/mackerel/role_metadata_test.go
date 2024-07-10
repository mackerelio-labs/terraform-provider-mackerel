package mackerel

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

func Test_ReadRoleMetadata(t *testing.T) {
	t.Parallel()

	defaultClient := func(serviceName, roleName, namespace string) (*mackerel.RoleMetaDataResp, error) {
		if serviceName != "service" || roleName != "role" || namespace != "namespace" {
			return nil, fmt.Errorf("no metadata found")
		}
		return &mackerel.RoleMetaDataResp{
			RoleMetaData: map[string]any{"v": 1},
		}, nil
	}

	cases := map[string]struct {
		inClient      roleMetadataReaderFunc
		inServiceName string
		inRoleName    string
		inNamespace   string

		wants RoleMetadataModel
	}{
		"basic": {
			inClient:      defaultClient,
			inServiceName: "service",
			inRoleName:    "role",
			inNamespace:   "namespace",

			wants: RoleMetadataModel{
				ID:           types.StringValue("service:role/namespace"),
				ServiceName:  types.StringValue("service"),
				RoleName:     types.StringValue("role"),
				Namespace:    types.StringValue("namespace"),
				MetadataJSON: jsontypes.NewNormalizedValue(`{"v":1}`),
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			data, err := readRoleMetadata(tt.inClient, tt.inServiceName, tt.inRoleName, tt.inNamespace)
			if err != nil {
				t.Errorf("unexpected error: %+v", err)
				return
			}

			if diff := cmp.Diff(data, tt.wants); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func Test_ImportRoleMetadata(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		inID string

		wants   RoleMetadataModel
		wantErr bool
	}{
		"valid": {
			inID: "service:role/namespace",

			wants: RoleMetadataModel{
				ID:          types.StringValue("service:role/namespace"),
				ServiceName: types.StringValue("service"),
				RoleName:    types.StringValue("role"),
				Namespace:   types.StringValue("namespace"),
			},
		},
		"invalid": {
			inID: "invalidid",

			wantErr: true,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			data, err := ImportRoleMetadata(tt.inID)
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error: %+v", err)
			}
			if err != nil {
				return
			}

			if diff := cmp.Diff(data, tt.wants); diff != "" {
				t.Error(diff)
			}
		})

	}
}

func Test_RoleMetadata_Create(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in       RoleMetadataModel
		inClient roleMetadataUpdatorFunc

		wants RoleMetadataModel
	}{
		"basic": {
			in: RoleMetadataModel{
				ServiceName:  types.StringValue("service"),
				RoleName:     types.StringValue("role"),
				Namespace:    types.StringValue("namespace"),
				MetadataJSON: jsontypes.NewNormalizedValue(`{"v":1}`),
			},
			inClient: func(serviceName, roleName, namespace string, metadata mackerel.RoleMetaData) error {
				if serviceName != "service" {
					return fmt.Errorf("unexpected service name: %s", serviceName)
				}
				if roleName != "role" {
					return fmt.Errorf("unexpected role name: %s", roleName)
				}
				if namespace != "namespace" {
					return fmt.Errorf("unexpected namespace: %s", namespace)
				}
				if diff := cmp.Diff(metadata, map[string]any{"v": 1.}); diff != "" {
					return fmt.Errorf("unexpected metadata: %s", diff)
				}
				return nil
			},

			wants: RoleMetadataModel{
				ID:           types.StringValue("service:role/namespace"),
				ServiceName:  types.StringValue("service"),
				RoleName:     types.StringValue("role"),
				Namespace:    types.StringValue("namespace"),
				MetadataJSON: jsontypes.NewNormalizedValue(`{"v":1}`),
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			data := tt.in
			if err := data.create(tt.inClient); err != nil {
				t.Errorf("unexpected error: %+v", err)
				return
			}

			if diff := cmp.Diff(data, tt.wants); diff != "" {
				t.Error(diff)
			}
		})
	}
}

type roleMetadataReaderFunc func(serviceName, roleName, namespace string) (*mackerel.RoleMetaDataResp, error)

func (f roleMetadataReaderFunc) GetRoleMetaData(serviceName, roleName, namespace string) (*mackerel.RoleMetaDataResp, error) {
	return f(serviceName, roleName, namespace)
}

type roleMetadataUpdatorFunc func(serviceName, roleName, namespace string, metadata mackerel.RoleMetaData) error

func (f roleMetadataUpdatorFunc) PutRoleMetaData(serviceName, roleName, namespace string, metadata mackerel.RoleMetaData) error {
	return f(serviceName, roleName, namespace, metadata)
}
