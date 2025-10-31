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

func Test_MackerelRoleMetadataDataSource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	req := fwdatasource.SchemaRequest{}
	resp := fwdatasource.SchemaResponse{}
	provider.NewMackerelRoleMetadataDataSource().Schema(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %+v", resp.Diagnostics)
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema validation diagnostics: %+v", diags)
	}
}

func TestAccDataSourceMackerelRoleMetadata(t *testing.T) {
	dsName := "data.mackerel_role_metadata.foo"
	rand := acctest.RandString(5)
	service := fmt.Sprintf("tf-service-%s", rand)
	role := fmt.Sprintf("tf-role-%s", rand)
	namespace := fmt.Sprintf("tf-namespace-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelRoleMetadataConfig(service, role, namespace),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsName, "id", fmt.Sprintf("%s:%s/%s", service, role, namespace)),
					resource.TestCheckResourceAttr(dsName, "service", service),
					resource.TestCheckResourceAttr(dsName, "role", role),
					resource.TestCheckResourceAttr(dsName, "namespace", namespace),
					resource.TestCheckResourceAttr(dsName, "metadata_json", `{"id":1}`),
				),
			},
		},
	})
}

func testAccDataSourceMackerelRoleMetadataConfig(service, role, namespace string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
}

resource "mackerel_role" "foo" {
  service = mackerel_service.foo.name
  name = "%s"
}

resource "mackerel_role_metadata" "foo" {
  service = mackerel_role.foo.service
  role = mackerel_role.foo.name
  namespace = "%s"
  metadata_json = jsonencode({
    id = 1
  })
}

data "mackerel_role_metadata" "foo" {
  depends_on = [mackerel_role_metadata.foo]
  service = mackerel_role_metadata.foo.service
  role = mackerel_role_metadata.foo.role
  namespace = mackerel_role_metadata.foo.namespace
}`, service, role, namespace)
}
