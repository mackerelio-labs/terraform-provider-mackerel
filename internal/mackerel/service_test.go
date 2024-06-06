package mackerel

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

func Test_ServiceNameValidator(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		val       types.String
		wantError bool
	}{
		"valid": {
			val: types.StringValue("service1"),
		},
		"too short": {
			val:       types.StringValue("a"),
			wantError: true,
		},
		"too long": {
			val:       types.StringValue("toooooooooooooooooooo-looooooooooooooooooooooooooooooooooooooooooong"),
			wantError: true,
		},
		"invalid char": {
			val:       types.StringValue("v('ω')v"),
			wantError: true,
		},
	}

	ctx := context.Background()
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			req := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    tt.val,
			}
			resp := &validator.StringResponse{}
			ServiceNameValidator().ValidateString(ctx, req, resp)

			hasError := resp.Diagnostics.HasError()
			if hasError != tt.wantError {
				if tt.wantError {
					t.Error("expected to have errors, but got no error")
				} else {
					t.Errorf("unexpected error: %+v", resp.Diagnostics.Errors())
				}
			}
		})
	}
}

type serviceFinderFunc func() ([]*mackerel.Service, error)

func (f serviceFinderFunc) FindServices() ([]*mackerel.Service, error) {
	return f()
}

func Test_ReadService(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		inClient serviceFinderFunc
		inName   string
		want     *ServiceModel
		wantFail bool
	}{
		"success": {
			inClient: func() ([]*mackerel.Service, error) {
				return []*mackerel.Service{
					{
						Name: "service0",
					},
					{
						Name: "service1",
						Memo: "memo",
					},
				}, nil
			},
			inName: "service1",
			want: &ServiceModel{
				ID:   types.StringValue("service1"),
				Name: types.StringValue("service1"),
				Memo: types.StringValue("memo"),
			},
		},
		"no service": {
			inClient: func() ([]*mackerel.Service, error) {
				return []*mackerel.Service{
					{
						Name: "service0",
					},
					{
						Name: "service1",
						Memo: "memo",
					},
				}, nil
			},
			inName:   "service2",
			wantFail: true,
		},
	}

	ctx := context.Background()
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s, err := readServiceInner(ctx, tt.inClient, tt.inName)
			if err != nil {
				if !tt.wantFail {
					t.Errorf("unexpected error: %+v", err)
				}
				return
			}
			if tt.wantFail {
				t.Errorf("unexpected success")
			}

			if diff := cmp.Diff(tt.want, s); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}
