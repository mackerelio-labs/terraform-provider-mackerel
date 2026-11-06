package provider_test

import (
	"context"
	"testing"

	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
)

func Test_MackerelIPRangesDataSource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwdatasource.SchemaRequest{}
	resp := &fwdatasource.SchemaResponse{}
	provider.NewMackerelIPRangesDataSource().Schema(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Errorf("schema method: %+v", resp.Diagnostics)
		return
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Errorf("schema validation: %+v", diags)
	}
}

func TestAccDataSourceMackerelIPRanges(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "mackerel_ip_ranges" "this" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mackerel_ip_ranges.this", "id", "https://mackerel.io/meta/ip-ranges.json"),
					resource.TestCheckResourceAttrSet("data.mackerel_ip_ranges.this", "incoming.0"),
					resource.TestCheckResourceAttrSet("data.mackerel_ip_ranges.this", "incoming_ipv6.0"),
					resource.TestCheckResourceAttrSet("data.mackerel_ip_ranges.this", "outgoing.0"),
					resource.TestCheckResourceAttrSet("data.mackerel_ip_ranges.this", "outgoing_ipv6.0"),
				),
			},
		},
	})
}
