package provider_test

import (
	"context"
	"testing"

	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
)

func Test_MackerelServiceMetadataDataSource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwdatasource.SchemaRequest{}
	resp := &fwdatasource.SchemaResponse{}
	provider.NewMackerelServiceMetadataDataSource().Schema(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Errorf("schema method: %+v", resp.Diagnostics)
		return
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Errorf("schema validation: %+v", diags)
	}
}
