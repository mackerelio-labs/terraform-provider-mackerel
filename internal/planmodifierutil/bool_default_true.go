package planmodifierutil

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BoolDefaultTrue returns a plan modifier that sets the bool value to true
// if the config value is null (not set by user).
func BoolDefaultTrue() planmodifier.Bool {
	return boolDefaultTrueModifier{}
}

type boolDefaultTrueModifier struct{}

func (boolDefaultTrueModifier) Description(context.Context) string {
	return "If the value is not set, defaults to true."
}

func (boolDefaultTrueModifier) MarkdownDescription(context.Context) string {
	return "If the value is not set, defaults to true."
}

func (m boolDefaultTrueModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	// If the config value is null (not set by user), set it to true
	if req.ConfigValue.IsNull() {
		resp.PlanValue = types.BoolValue(true)
	}
}
