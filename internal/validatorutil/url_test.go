package validatorutil_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/validatorutil"
)

func Test_Validator_URL(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		val       types.String
		wantError bool
	}{
		"https": {
			val: types.StringValue("https://example.com/path/to/resource?q="),
		},
		"http": {
			val: types.StringValue("http://example.com/path/to/resource?q="),
		},
		"empty": {
			val:       types.StringValue(""),
			wantError: true,
		},
		"not url": {
			val:       types.StringValue(":nonurlstring"),
			wantError: true,
		},
		"no host": {
			val:       types.StringValue("http:///path/to/resource"),
			wantError: true,
		},
		"invalid scheme": {
			val:       types.StringValue("ftp://example.com/path/to/resource"),
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
			validatorutil.IsURLWithHTTPorHTTPS().ValidateString(ctx, req, resp)

			for _, d := range resp.Diagnostics {
				assertDiagMatchPathExpr(t, d, path.MatchRoot("test"))
			}

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

func assertDiagMatchPathExpr(t testing.TB, d diag.Diagnostic, pathExpr path.Expression) bool {
	t.Helper()

	dp, ok := d.(diag.DiagnosticWithPath)
	if !ok {
		t.Errorf("expected to have a path, but got no path: %+v", d)
		return true
	}

	if !pathExpr.Matches(dp.Path()) {
		t.Errorf("expteted to have a path that matches to %s, but got: %+v", pathExpr.String(), dp.Path())
		return true
	}

	return false
}
