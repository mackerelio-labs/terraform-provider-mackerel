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
		"unknown to null": {
			request: planmodifier.MapRequest{
				PlanValue: types.MapUnknown(types.StringType),
			},
			expected: &planmodifier.MapResponse{
				PlanValue: types.MapNull(types.StringType),
			},
		},
		"no updates when plan is null and state is empty": {
			request: planmodifier.MapRequest{
				PlanValue:  types.MapNull(types.StringType),
				StateValue: types.MapValueMust(types.StringType, map[string]attr.Value{}),
			},
			expected: &planmodifier.MapResponse{
				PlanValue: types.MapValueMust(types.StringType, map[string]attr.Value{}),
			},
		},
		"passthrough otherwise": {
			request: planmodifier.MapRequest{
				PlanValue: types.MapValueMust(types.StringType, map[string]attr.Value{"a": types.StringValue("a")}),
			},
			expected: &planmodifier.MapResponse{
				PlanValue: types.MapValueMust(types.StringType, map[string]attr.Value{"a": types.StringValue("a")}),
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
		"unknown to null": {
			request: planmodifier.SetRequest{
				PlanValue: types.SetUnknown(types.StringType),
			},
			expected: &planmodifier.SetResponse{
				PlanValue: types.SetNull(types.StringType),
			},
		},
		"no updates when plan is null and state is empty": {
			request: planmodifier.SetRequest{
				PlanValue:  types.SetNull(types.StringType),
				StateValue: types.SetValueMust(types.StringType, []attr.Value{}),
			},
			expected: &planmodifier.SetResponse{
				PlanValue: types.SetValueMust(types.StringType, []attr.Value{}),
			},
		},
		"passthrough otherwise": {
			request: planmodifier.SetRequest{
				PlanValue: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("a")}),
			},
			expected: &planmodifier.SetResponse{
				PlanValue: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("a")}),
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

func TestNilRelaxedList(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		request  planmodifier.ListRequest
		expected *planmodifier.ListResponse
	}{
		"unknown to null": {
			request: planmodifier.ListRequest{
				PlanValue: types.ListUnknown(types.StringType),
			},
			expected: &planmodifier.ListResponse{
				PlanValue: types.ListNull(types.StringType),
			},
		},
		"no updates when plan is null and state is empty": {
			request: planmodifier.ListRequest{
				PlanValue:  types.ListNull(types.StringType),
				StateValue: types.ListValueMust(types.StringType, []attr.Value{}),
			},
			expected: &planmodifier.ListResponse{
				PlanValue: types.ListValueMust(types.StringType, []attr.Value{}),
			},
		},
		"passthrough otherwise": {
			request: planmodifier.ListRequest{
				PlanValue: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("a")}),
			},
			expected: &planmodifier.ListResponse{
				PlanValue: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("a")}),
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &planmodifier.ListResponse{
				PlanValue: tt.request.PlanValue,
			}

			planmodifierutil.NilRelaxedList().PlanModifyList(context.Background(), tt.request, resp)

			if diff := cmp.Diff(tt.expected, resp); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}
