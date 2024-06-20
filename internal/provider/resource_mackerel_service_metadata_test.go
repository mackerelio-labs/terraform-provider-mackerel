package provider_test

import (
	"context"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
)

func Test_MackerelServiceMetadataResource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwresource.SchemaRequest{}
	resp := &fwresource.SchemaResponse{}
	provider.NewMackerelServiceMetadataResource().Schema(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Errorf("schema method: %+v", resp.Diagnostics)
		return
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Errorf("schema validation: %+v", diags)
	}
}
