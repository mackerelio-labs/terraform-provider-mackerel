package provider_test

import (
	"context"
	"fmt"
	"testing"

	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
)

func Test_MackerelServiceDataSource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwdatasource.SchemaRequest{}
	resp := &fwdatasource.SchemaResponse{}
	provider.NewMackerelServiceDataSource().Schema(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema method diagnostics: %+v", resp.Diagnostics)
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema validation diagnostics: %+v", diags)
	}
}

func TestAcc_MackerelServiceDataSource_basic(t *testing.T) {
	name := fmt.Sprintf("tf-service-%s", acctest.RandString(5))
	resourceName := "data.mackerel_service.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: mackerelServiceDataSourceConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("memo"),
						knownvalue.StringExact("This service is managed by Terraform."),
					),
				},
			},
		},
	})
}

func mackerelServiceDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
  memo = "This service is managed by Terraform."
}

data "mackerel_service" "foo" {
  name = mackerel_service.foo.id
}
`, name)
}
