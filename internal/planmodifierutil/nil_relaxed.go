package planmodifierutil

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NilRelaxedMap() planmodifier.Map {
	return nilRelaxedModifier{}
}

func NilRelaxedSet() planmodifier.Set {
	return nilRelaxedModifier{}
}

func NilRelaxedList() planmodifier.List {
	return nilRelaxedModifier{}
}

type nilRelaxedModifier struct{}

const desctiprion = "For compatibility with the states created by SDK provider, Terraform consider nil and zero values to be same."

func (nilRelaxedModifier) Description(context.Context) string {
	return desctiprion
}

func (nilRelaxedModifier) MarkdownDescription(context.Context) string {
	return desctiprion
}

func (nilRelaxedModifier) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	if req.PlanValue.IsUnknown() {
		resp.PlanValue = types.MapNull(req.PlanValue.ElementType(ctx))
	} else if req.PlanValue.IsNull() && len(req.StateValue.Elements()) == 0 {
		resp.PlanValue = req.StateValue
	}
}

func (nilRelaxedModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if req.PlanValue.IsUnknown() {
		resp.PlanValue = types.SetNull(req.PlanValue.ElementType(ctx))
	} else if req.PlanValue.IsNull() && len(req.StateValue.Elements()) == 0 {
		resp.PlanValue = req.StateValue
	}
}

func (nilRelaxedModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.PlanValue.IsUnknown() {
		resp.PlanValue = types.ListNull(req.PlanValue.ElementType(ctx))
	} else if req.PlanValue.IsNull() && len(req.StateValue.Elements()) == 0 {
		resp.PlanValue = req.StateValue
	}
}
