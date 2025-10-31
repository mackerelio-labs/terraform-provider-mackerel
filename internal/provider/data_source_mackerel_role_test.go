package provider_test

import (
	"context"
	"fmt"
	"testing"

	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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

func TestAccDataSourceMackerelRole(t *testing.T) {
	dsName := "data.mackerel_role.foo"
	rand := acctest.RandString(5)
	serviceName := fmt.Sprintf("tf-service-%s", rand)
	name := fmt.Sprintf("tf-role-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelRoleConfig(serviceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsName, "id", fmt.Sprintf("%s:%s", serviceName, name)),
					resource.TestCheckResourceAttr(dsName, "service", serviceName),
					resource.TestCheckResourceAttr(dsName, "name", name),
					resource.TestCheckResourceAttr(dsName, "memo", "This role is managed by Terraform."),
				),
			},
		},
	})
}

func testAccDataSourceMackerelRoleConfig(serviceName, name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
}

resource "mackerel_role" "foo" {
  service = mackerel_service.foo.name
  name = "%s"
  memo = "This role is managed by Terraform."
}

data "mackerel_role" "foo" {
  depends_on = [mackerel_role.foo]
  service = mackerel_role.foo.service
  name = mackerel_role.foo.name
}
`, serviceName, name)
}
