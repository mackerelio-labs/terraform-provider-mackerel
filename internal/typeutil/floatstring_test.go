package typeutil_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/typeutil"
)

func Test_FloatStringType_Validate(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in      tftypes.Value
		wantErr bool
	}{
		"empty struct": {
			in: tftypes.Value{},
		},
		"null": {
			in: tftypes.NewValue(tftypes.String, nil),
		},
		"unknown": {
			in: tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		},
		"empty string": {
			in: tftypes.NewValue(tftypes.String, ""),
		},
		"valid float string": {
			in: tftypes.NewValue(tftypes.String, ".1"),
		},
		"invalid float string": {
			in:      tftypes.NewValue(tftypes.String, "xyz"),
			wantErr: true,
		},
		"wrong type": {
			in:      tftypes.NewValue(tftypes.Number, .1),
			wantErr: true,
		},
	}

	ctx := context.Background()
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			diags := typeutil.FloatStringType{}.Validate(ctx, tt.in, path.Root("test"))

			if diags.HasError() != tt.wantErr {
				t.Errorf("unexpected diags: %+v", diags)
			}
		})
	}
}

func Test_FloatStringType_FromTerraform(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in      tftypes.Value
		wants   attr.Value
		wantErr bool
	}{
		"value": {
			in:    tftypes.NewValue(tftypes.String, ".1"),
			wants: typeutil.NewFloatStringValue(".1"),
		},
		"unknown": {
			in:    tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			wants: typeutil.NewFloatStringUnknown(),
		},
		"null": {
			in:    tftypes.NewValue(tftypes.String, nil),
			wants: typeutil.NewFloatStringNull(),
		},
		"empty string": {
			in:    tftypes.NewValue(tftypes.String, ""),
			wants: typeutil.NewFloatStringValue(""),
		},
		"wrong type": {
			in:      tftypes.NewValue(tftypes.Number, .1),
			wantErr: true,
		},
	}

	ctx := context.Background()
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := typeutil.FloatStringType{}.ValueFromTerraform(ctx, tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error: %+v", err)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(got, tt.wants); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func Test_FloatString_StringSemanticEquals(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in        typeutil.FloatString
		inOther   basetypes.StringValuable
		wantMatch bool
		wantErr   bool
	}{
		"strong equal": {
			in:        typeutil.NewFloatStringValue("0.1"),
			inOther:   typeutil.NewFloatStringValue("0.1"),
			wantMatch: true,
		},
		"point": {
			in:        typeutil.NewFloatStringValue("0.1"),
			inOther:   typeutil.NewFloatStringValue(".1"),
			wantMatch: true,
		},
		"exp": {
			in:        typeutil.NewFloatStringValue("0.1"),
			inOther:   typeutil.NewFloatStringValue("1e-1"),
			wantMatch: true,
		},
		"wrong type": {
			in:      typeutil.NewFloatStringValue("0.1"),
			inOther: basetypes.NewStringValue("0.1"),
			wantErr: true,
		},
	}

	ctx := context.Background()
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			match, diags := tt.in.StringSemanticEquals(ctx, tt.inOther)
			if diags.HasError() != tt.wantErr {
				t.Errorf("unexpected diags: %+v", diags)
			}
			if diags.HasError() {
				return
			}
			if match != tt.wantMatch {
				t.Error("unexpected matching result")
			}
		})
	}
}

func Test_FloatString_ValueFloat64(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in    typeutil.FloatString
		wants float64
	}{
		"value": {
			in:    typeutil.NewFloatStringValue("0.1"),
			wants: 0.1,
		},
		"unknown": {
			in:    typeutil.NewFloatStringUnknown(),
			wants: 0.0,
		},
		"null": {
			in:    typeutil.NewFloatStringNull(),
			wants: 0.0,
		},
		"invalid": {
			in:    typeutil.NewFloatStringValue("invalid"),
			wants: 0.0,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tt.in.ValueFloat64() != tt.wants {
				t.Error("unmatched")
			}
		})
	}
}
