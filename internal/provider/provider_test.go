package provider_test

import (
	"context"
	"testing"

	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
	sdkmackerel "github.com/mackerelio-labs/terraform-provider-mackerel/mackerel"
)

var (
	protoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
		"mackerel": providerserver.NewProtocol5WithError(provider.New()),
	}
	protoV5SDKProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
		"mackerel": func() (tfprotov5.ProviderServer, error) {
			return sdkmackerel.Provider().GRPCProvider(), nil
		},
	}
)

func preCheck(t *testing.T) {
	t.Helper()
}

func stepNoPlanInFramework(config string) resource.TestStep {
	return resource.TestStep{
		Config:                   config,
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		ConfigPlanChecks: resource.ConfigPlanChecks{
			PreApply: []plancheck.PlanCheck{
				plancheck.ExpectEmptyPlan(),
			},
		},
	}
}

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

func TestMackerelProvider_Configure_WithConfigOnly(t *testing.T) {
	// Clear environment variables
	t.Setenv("MACKEREL_API_KEY", "")
	t.Setenv("MACKEREL_APIKEY", "")

	ctx := context.Background()
	p := provider.New()

	// Get schema
	schemaReq := fwprovider.SchemaRequest{}
	schemaResp := &fwprovider.SchemaResponse{}
	p.Schema(ctx, schemaReq, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("Schema error: %v", schemaResp.Diagnostics)
	}

	// Create config with api_key set
	configValue := tftypes.NewValue(
		tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"api_key":  tftypes.String,
				"api_base": tftypes.String,
			},
		},
		map[string]tftypes.Value{
			"api_key":  tftypes.NewValue(tftypes.String, "test_api_key_from_config"),
			"api_base": tftypes.NewValue(tftypes.String, nil),
		},
	)

	req := fwprovider.ConfigureRequest{
		Config: tfsdk.Config{
			Schema: schemaResp.Schema,
			Raw:    configValue,
		},
	}
	resp := &fwprovider.ConfigureResponse{}

	p.Configure(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("Expected no error, but got: %v", resp.Diagnostics)
	}
	if resp.ResourceData == nil {
		t.Error("Expected ResourceData to be set, but got nil")
	}
}
