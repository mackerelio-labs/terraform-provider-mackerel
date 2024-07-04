package provider_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/mackerel"
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

// For acceptance tests

func protoV5ProviderFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"mackerel": providerserver.NewProtocol5WithError(provider.New()),
	}
}

func preCheck(t *testing.T) {
	t.Helper()

	// Currently, do nothing
}

func newClient(t testing.TB) *mackerel.Client {
	t.Helper()

	config := mackerel.NewClientConfigFromEnv()
	client, err := config.NewClient()
	if err != nil {
		if errors.Is(err, mackerel.ErrNoAPIKey) {
			t.Fatal("MACKEREL_API_KEY or MACKEREL_APIKEY is required for acceptance tests.")
		} else {
			t.Fatalf("Failed to create Mackerel client: %+v", err)
		}
	}
	return client
}

// state check helpers
type stateCheckFunc func(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse)

func (sc stateCheckFunc) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	sc(ctx, req, resp)
}

func findStateResource(state *tfjson.State, resoruceAddress string) (*tfjson.StateResource, error) {
	if state == nil {
		return nil, fmt.Errorf("state is nil")
	}
	if state.Values == nil {
		return nil, fmt.Errorf("state does not contain any state values")
	}
	if state.Values.RootModule == nil {
		return nil, fmt.Errorf("state does not contain a root module")
	}
	for _, r := range state.Values.RootModule.Resources {
		if r.Address == resoruceAddress {
			return r, nil
		}
	}
	return nil, fmt.Errorf("%s - Resource not found in state", resoruceAddress)
}
