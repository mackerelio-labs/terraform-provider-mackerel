package provider_test

import (
	"context"
	"testing"

	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
)

func TestMackerelProvider_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwprovider.SchemaRequest{}
	resp := &fwprovider.SchemaResponse{}
	provider.New().Schema(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema: %+v", resp.Diagnostics)
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("Schema validation: %+v", diags)
	}
}

func TestMackerelProvider_apiKey(t *testing.T) {
	testCases := map[string]struct {
		MACKEREL_API_KEY string
		MACKEREL_APIKEY  string
	}{
		"MACKEREL_API_KEY": {
			MACKEREL_API_KEY: "apikey1",
		},
		"MACKEREL_APIKEY": {
			MACKEREL_APIKEY: "apikey1",
		},
	}

	ctx := context.Background()
	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("MACKEREL_API_KEY", tt.MACKEREL_API_KEY)
			t.Setenv("MACKEREL_APIKEY", tt.MACKEREL_APIKEY)

			p := provider.New()

			creq := newProviderConfigureRequest(ctx, nil)
			cresp := &fwprovider.ConfigureResponse{}
			p.Configure(ctx, creq, cresp)
			if cresp.Diagnostics.HasError() {
				t.Errorf("Configure: %+v", cresp.Diagnostics)
				return
			}
		})
	}
}

func newProviderConfigureRequest(
	ctx context.Context,
	c *provider.MackerelProviderModel,
) fwprovider.ConfigureRequest {
	if c == nil {
		c = &provider.MackerelProviderModel{}
	}
	p := provider.New()

	sreq := fwprovider.SchemaRequest{}
	sresp := &fwprovider.SchemaResponse{}
	p.Schema(ctx, sreq, sresp)
	schema := sresp.Schema

	state := tfsdk.State{
		Schema: schema,
	}
	state.Set(ctx, c)
	return fwprovider.ConfigureRequest{
		Config: tfsdk.Config{
			Schema: schema,
			Raw:    state.Raw,
		},
	}
}
