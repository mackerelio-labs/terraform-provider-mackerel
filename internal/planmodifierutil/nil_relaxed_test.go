package planmodifierutil_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/planmodifierutil"
)

func TestNilRelaxedMap(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		request  planmodifier.MapRequest
		expected *planmodifier.MapResponse
	}{
		"non-empty": {
			request: planmodifier.MapRequest{
				StateValue: types.MapValueMust(types.StringType, map[string]attr.Value{
					"a": types.StringValue("b"),
				}),
				PlanValue: types.MapValueMust(types.StringType, map[string]attr.Value{
					"a": types.StringValue("c"),
				}),
				ConfigValue: types.MapValueMust(types.StringType, map[string]attr.Value{
					"a": types.StringValue("c"),
				}),
			},
			expected: &planmodifier.MapResponse{
				PlanValue: types.MapValueMust(types.StringType, map[string]attr.Value{
					"a": types.StringValue("c"),
				}),
			},
		},
		"empty-null": {
			request: planmodifier.MapRequest{
				StateValue:  types.MapValueMust(types.StringType, map[string]attr.Value{}),
				PlanValue:   types.MapNull(types.StringType),
				ConfigValue: types.MapNull(types.StringType),
			},
			expected: &planmodifier.MapResponse{
				PlanValue: types.MapValueMust(types.StringType, map[string]attr.Value{}),
			},
		},
		"null-empty": {
			request: planmodifier.MapRequest{
				StateValue:  types.MapNull(types.StringType),
				PlanValue:   types.MapValueMust(types.StringType, map[string]attr.Value{}),
				ConfigValue: types.MapValueMust(types.StringType, map[string]attr.Value{}),
			},
			expected: &planmodifier.MapResponse{
				PlanValue: types.MapNull(types.StringType),
			},
		},
		"unknown": {
			request: planmodifier.MapRequest{
				StateValue:  types.MapNull(types.StringType),
				PlanValue:   types.MapUnknown(types.StringType),
				ConfigValue: types.MapUnknown(types.StringType),
			},
			expected: &planmodifier.MapResponse{
				PlanValue: types.MapUnknown(types.StringType),
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &planmodifier.MapResponse{
				PlanValue: tt.request.PlanValue,
			}

			planmodifierutil.NilRelaxedMap().PlanModifyMap(context.Background(), tt.request, resp)

			if diff := cmp.Diff(tt.expected, resp); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func TestNilRelaxedSet(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		request  planmodifier.SetRequest
		expected *planmodifier.SetResponse
	}{
		"empty-null": {
			request: planmodifier.SetRequest{
				StateValue:  types.SetValueMust(types.StringType, []attr.Value{}),
				PlanValue:   types.SetNull(types.StringType),
				ConfigValue: types.SetNull(types.StringType),
			},
			expected: &planmodifier.SetResponse{
				PlanValue: types.SetValueMust(types.StringType, []attr.Value{}),
			},
		},
		"null-empty": {
			request: planmodifier.SetRequest{
				StateValue:  types.SetNull(types.StringType),
				PlanValue:   types.SetValueMust(types.StringType, []attr.Value{}),
				ConfigValue: types.SetValueMust(types.StringType, []attr.Value{}),
			},
			expected: &planmodifier.SetResponse{
				PlanValue: types.SetNull(types.StringType),
			},
		},
		"unknown": {
			request: planmodifier.SetRequest{
				StateValue:  types.SetNull(types.StringType),
				PlanValue:   types.SetUnknown(types.StringType),
				ConfigValue: types.SetUnknown(types.StringType),
			},
			expected: &planmodifier.SetResponse{
				PlanValue: types.SetUnknown(types.StringType),
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &planmodifier.SetResponse{
				PlanValue: tt.request.PlanValue,
			}

			planmodifierutil.NilRelaxedSet().PlanModifySet(context.Background(), tt.request, resp)

			if diff := cmp.Diff(tt.expected, resp); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}
