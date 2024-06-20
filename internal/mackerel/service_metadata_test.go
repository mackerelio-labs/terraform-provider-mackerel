package mackerel

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

func Test_ReadServiceMetadata(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		inClient serviceMetadataGetterFunc
		in       ServiceMetadataModel

		wants    ServiceMetadataModel
		wantFail bool
	}{
		"from service name and namespace": {
			inClient: func(s, ns string) (*mackerel.ServiceMetaDataResp, error) {
				if s != "service0" || ns != "data0" {
					return nil, fmt.Errorf("no metadata found")
				}
				return &mackerel.ServiceMetaDataResp{
					ServiceMetaData: map[string]any{
						"foo": "bar",
					},
				}, nil
			},
			in: ServiceMetadataModel{
				ID:          types.StringUnknown(),
				ServiceName: types.StringValue("service0"),
				Namespace:   types.StringValue("data0"),
			},

			wants: ServiceMetadataModel{
				ID:           types.StringValue("service0/data0"),
				ServiceName:  types.StringValue("service0"),
				Namespace:    types.StringValue("data0"),
				MetadataJSON: jsontypes.NewNormalizedValue(`{"foo":"bar"}`),
			},
		},
		"from id": {
			inClient: func(s, ns string) (*mackerel.ServiceMetaDataResp, error) {
				if s != "service0" || ns != "data0" {
					return nil, fmt.Errorf("no metadata found")
				}
				return &mackerel.ServiceMetaDataResp{
					ServiceMetaData: map[string]any{
						"foo": "bar",
					},
				}, nil
			},
			in: ServiceMetadataModel{
				ID: types.StringValue("service0/data0"),
			},

			wants: ServiceMetadataModel{
				ID:           types.StringValue("service0/data0"),
				ServiceName:  types.StringValue("service0"),
				Namespace:    types.StringValue("data0"),
				MetadataJSON: jsontypes.NewNormalizedValue(`{"foo":"bar"}`),
			},
		},
	}

	ctx := context.Background()
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual, err := readServiceMetadataInner(ctx, tt.inClient, tt.in)
			if (err != nil) != tt.wantFail {
				if tt.wantFail {
					t.Errorf("unexpected success")
				} else {
					t.Errorf("unexpected error: %+v", err)
				}
				return
			}

			if diff := cmp.Diff(tt.wants, actual); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}

}

type serviceMetadataGetterFunc func(string, string) (*mackerel.ServiceMetaDataResp, error)

func (f serviceMetadataGetterFunc) GetServiceMetaData(serviceName, namespace string) (*mackerel.ServiceMetaDataResp, error) {
	return f(serviceName, namespace)
}

func Test_ServiceMetadata_Validate(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in ServiceMetadataModel

		wantError   bool
		wantErrorIn path.Expressions
	}{
		"valid": {
			in: ServiceMetadataModel{
				ID:          types.StringValue("service/namespace"),
				ServiceName: types.StringValue("service"),
				Namespace:   types.StringValue("namespace"),
			},
		},
		"invalid id syntax": {
			in: ServiceMetadataModel{
				ID: types.StringValue("service,namespace"),
			},
			wantError:   true,
			wantErrorIn: path.Expressions{path.MatchRoot("id")},
		},
		"unmatched service": {
			in: ServiceMetadataModel{
				ID:          types.StringValue("service0/namespace"),
				ServiceName: types.StringValue("service1"),
				Namespace:   types.StringValue("namespace"),
			},
			wantError:   true,
			wantErrorIn: path.Expressions{path.MatchRoot("id")},
		},
		"unmatched namespace": {
			in: ServiceMetadataModel{
				ID:          types.StringValue("service/namespace0"),
				ServiceName: types.StringValue("service"),
				Namespace:   types.StringValue("namespace1"),
			},
			wantError:   true,
			wantErrorIn: path.Expressions{path.MatchRoot("id")},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			diags := tt.in.Validate(path.Empty())
			for _, d := range diags {
				if d.Severity() != diag.SeverityError {
					continue
				}
				dwp, ok := d.(diag.DiagnosticWithPath)
				if ok {
					p := dwp.Path()
					if slices.ContainsFunc(tt.wantErrorIn, func(expr path.Expression) bool {
						return expr.Matches(p)
					}) {
						continue
					}
				} else if tt.wantError {
					continue
				}
				t.Errorf("unexpected error: %v", d)
			}
		})
	}
}
