package planmodifierutil_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/planmodifierutil"
)

func TestBoolDefaultTrue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		configValue types.Bool
		expected    types.Bool
	}{
		"null_config_sets_true": {
			configValue: types.BoolNull(),
			expected:    types.BoolValue(true),
		},
		"explicit_true_remains_true": {
			configValue: types.BoolValue(true),
			expected:    types.BoolValue(true),
		},
		"explicit_false_remains_false": {
			configValue: types.BoolValue(false),
			expected:    types.BoolValue(false),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			req := planmodifier.BoolRequest{
				ConfigValue: testCase.configValue,
				PlanValue:   testCase.configValue,
			}
			resp := &planmodifier.BoolResponse{
				PlanValue: testCase.configValue,
			}

			planmodifierutil.BoolDefaultTrue().PlanModifyBool(ctx, req, resp)

			if !resp.PlanValue.Equal(testCase.expected) {
				t.Errorf("expected %v, got %v", testCase.expected, resp.PlanValue)
			}
		})
	}
}
