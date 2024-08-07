package provider_test

import (
	"context"
	"testing"

	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
)

func Test_MackerelRoleDataSource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	req := fwdatasource.SchemaRequest{}
	resp := fwdatasource.SchemaResponse{}
	provider.NewMackerelRoleDataSource().Schema(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema: %+v", resp.Diagnostics)
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema validation: %+v", diags)
	}
}
