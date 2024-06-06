package mackerel

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Test_Mackerel_ServiceNameValidator(t *testing.T) {
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
