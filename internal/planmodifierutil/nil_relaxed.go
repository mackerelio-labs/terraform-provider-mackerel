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

type nilRelaxedModifier struct{}

const desctiprion = "For compatibility with the states created by SDK provider, Terraform consider nil and zero values to be same."

func (_ nilRelaxedModifier) Description(context.Context) string {
	return desctiprion
}

func (_ nilRelaxedModifier) MarkdownDescription(context.Context) string {
	return desctiprion
}

func (_ nilRelaxedModifier) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	if req.PlanValue.IsUnknown() {
		resp.PlanValue = types.MapNull(req.PlanValue.ElementType(ctx))
	} else if req.PlanValue.IsNull() && len(req.StateValue.Elements()) == 0 {
		resp.PlanValue = req.StateValue
	}
}

func (_ nilRelaxedModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if req.PlanValue.IsUnknown() {
		resp.PlanValue = types.SetNull(req.PlanValue.ElementType(ctx))
	} else if req.PlanValue.IsNull() && len(req.StateValue.Elements()) == 0 {
		resp.PlanValue = req.StateValue
	}
}
