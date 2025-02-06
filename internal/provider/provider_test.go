package provider_test

import (
	"context"
	"testing"

	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
	sdkmackerel "github.com/mackerelio-labs/terraform-provider-mackerel/mackerel"
)

var (
	protoV5FrameworkProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
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
		ProtoV5ProviderFactories: protoV5FrameworkProviderFactories,
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
