package provider_test

import (
	"context"
	"testing"

	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
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
