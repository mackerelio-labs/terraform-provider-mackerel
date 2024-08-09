package typeutil

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type (
	FloatStringType struct {
		basetypes.StringType
	}
	FloatString struct {
		basetypes.StringValue
	}
)

var (
	_ basetypes.StringTypable = (*FloatStringType)(nil)
)

func (t FloatStringType) String() string {
	return "typeutil.FloatStringType"
}

func (t FloatStringType) ValueType(ctx context.Context) attr.Value {
	return FloatString{}
}

func (t FloatStringType) Equal(o attr.Type) bool {
	other, ok := o.(FloatStringType)
	if !ok {
		return false
	}
	return t.StringType.Equal(other.StringType)
}

func (t FloatStringType) Validate(ctx context.Context, in tftypes.Value, path path.Path) (diags diag.Diagnostics) {
	if in.Type() == nil {
		return diags
	}

	if !in.Type().Is(tftypes.String) {
		msg := fmt.Sprintf("expected String value, but got %T with value: %v", in, in)
		diags.AddAttributeError(
			path,
			"Float String Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
				"Please report the following to the provider developer:\n\n"+msg,
		)
		return diags
	}

	if !in.IsKnown() || in.IsNull() {
		return diags
	}

	var valueString string
	if err := in.As(&valueString); err != nil {
		diags.AddAttributeError(
			path,
			"Float String Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
				"Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return diags
	}

	// For legacy reasons, empty strings are treated as a kind of empty value.
	if valueString == "" {
		return diags
	}

	if _, err := strconv.ParseFloat(valueString, 64); err != nil {
		diags.AddAttributeError(
			path,
			"Invalid Float String Value",
			"A string value was provided that is not valid float string format.\n\n"+
				"Given value: "+valueString+"\n"+
				err.Error(),
		)
		return diags
	}

	return diags
}

func (t FloatStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return FloatString{StringValue: in}, nil
}

func (t FloatStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	return FloatString{StringValue: stringValue}, nil
}

var (
	_ basetypes.StringValuable                   = (*FloatString)(nil)
	_ basetypes.Float64Valuable                  = (*FloatString)(nil)
	_ basetypes.StringValuableWithSemanticEquals = (*FloatString)(nil)
)

func NewFloatStringNull() FloatString {
	return FloatString{StringValue: basetypes.NewStringNull()}
}

func NewFloatStringUnknown() FloatString {
	return FloatString{StringValue: basetypes.NewStringUnknown()}
}

func NewFloatStringValue(value string) FloatString {
	return FloatString{StringValue: basetypes.NewStringValue(value)}
}

func NewFloatStringPointerValue(value *string) FloatString {
	return FloatString{StringValue: basetypes.NewStringPointerValue(value)}
}

func (v FloatString) Type(_ context.Context) attr.Type {
	return FloatStringType{}
}

func (v FloatString) ToFloat64Value(_ context.Context) (basetypes.Float64Value, diag.Diagnostics) {
	if v.StringValue.IsUnknown() {
		return basetypes.NewFloat64Unknown(), nil
	}
	if v.StringValue.IsNull() {
		return basetypes.NewFloat64Null(), nil
	}

	var diags diag.Diagnostics
	valueString := v.ValueString()
	f, err := strconv.ParseFloat(valueString, 64)
	if err != nil {
		diags.AddError(
			"Conversion Error",
			"An unexected string value was received while converting to float value. "+
				"Please report this to the provider developers.\n\n"+
				"Got value: "+valueString+"\n"+
				err.Error(),
		)
		return basetypes.Float64Value{}, diags
	}
	return basetypes.NewFloat64Value(f), nil
}

func (v FloatString) ValueFloat64() float64 {
	fp := v.ValueFloat64Pointer()
	if fp == nil {
		return 0.0
	}
	return *fp
}

func (v FloatString) ValueFloat64Pointer() *float64 {
	if v.StringValue.IsUnknown() || v.StringValue.IsNull() {
		return nil
	}
	f, err := strconv.ParseFloat(v.ValueString(), 64)
	if err != nil {
		return nil
	}
	return &f
}

func (v FloatString) Equal(o attr.Value) bool {
	other, ok := o.(FloatString)
	if !ok {
		return false
	}
	return v.StringValue.Equal(other.StringValue)
}

func (v FloatString) StringSemanticEquals(ctx context.Context, otherValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	otherValue, ok := otherValuable.(FloatString)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected: "+fmt.Sprintf("%T", v)+"\n"+
				"Got: "+fmt.Sprintf("%T", otherValuable),
		)
		return false, diags
	}

	valueFloat, diagConv := v.ToFloat64Value(ctx)
	diags.Append(diagConv...)
	otherValueFloat, diagConv := otherValue.ToFloat64Value(ctx)
	diags.Append(diagConv...)
	if diags.HasError() {
		return false, diags
	}

	return valueFloat.Equal(otherValueFloat), nil
}
