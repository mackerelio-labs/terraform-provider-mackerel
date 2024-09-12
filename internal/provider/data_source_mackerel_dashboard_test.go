package provider_test

import (
	"context"
	"testing"

	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
)

func Test_MackerelDashboardDataSource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwdatasource.SchemaRequest{}
	resp := fwdatasource.SchemaResponse{}
	if provider.NewMackerelDashboardDataSource().Schema(ctx, req, &resp); resp.Diagnostics.HasError() {
		t.Fatalf("schema method diagnostics: %+v", resp.Diagnostics)
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema validation diagnostics: %+v", diags)
	}
}
