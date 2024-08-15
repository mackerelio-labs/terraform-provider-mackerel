package validatorutil_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/validatorutil"
)

func Test_Validator_Divisible(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		inDivisor int64
		inVal     types.Int64
		wantErr   bool
	}{
		"divisible": {
			inDivisor: 19,
			inVal:     types.Int64Value(57),
		},
		"indivisible": {
			inDivisor: 10,
			inVal:     types.Int64Value(101),
			wantErr:   true,
		},
	}

	ctx := context.Background()
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := validator.Int64Request{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    tt.inVal,
			}
			resp := validator.Int64Response{}

			validatorutil.IntDivisibleBy(tt.inDivisor).ValidateInt64(ctx, req, &resp)

			for _, d := range resp.Diagnostics {
				assertDiagMatchPathExpr(t, d, path.MatchRoot("test"))
			}

			if resp.Diagnostics.HasError() != tt.wantErr {
				t.Errorf("unexpected error: %+v", resp.Diagnostics.Errors())
			}
		})
	}
}
