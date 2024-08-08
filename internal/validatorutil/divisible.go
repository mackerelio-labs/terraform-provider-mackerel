package validatorutil

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type divisibleValidator struct {
	divisor int64
}

var _ validator.Int64 = (*divisibleValidator)(nil)

func IntDivisibleBy(divisor int64) validator.Int64 {
	return &divisibleValidator{
		divisor: divisor,
	}
}

func (v *divisibleValidator) Description(context.Context) string {
	return fmt.Sprintf("integer which is divisible by %d", v.divisor)
}

func (v *divisibleValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *divisibleValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	iv := req.ConfigValue.ValueInt64()
	if iv%v.divisor != 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Indivisible integer",
			fmt.Sprintf("expected to be divisible by %d, got: %d", v.divisor, iv),
		)
		return
	}
}
