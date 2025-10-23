package provider_test

import (
	"context"
	"fmt"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
)

func Test_MackerelRoleMetadataResource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	req := fwresource.SchemaRequest{}
	resp := fwresource.SchemaResponse{}
	provider.NewMackerelRoleMetadataResource().Schema(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %+v", resp.Diagnostics)
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema validation diagnostics: %+v", diags)
	}
}

func TestAccMackerelRoleMetadata(t *testing.T) {
	resourceName := "mackerel_role_metadata.foo"
	rand := acctest.RandString(5)
	rServiceName := fmt.Sprintf("tf-%s", rand)
	rRoleName := fmt.Sprintf("tf-%s-role", rand)
	rNamespace := fmt.Sprintf("tf-namespace-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		CheckDestroy:             testAccCheckMackerelRoleMetadataDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelRoleMetadataConfig(rServiceName, rRoleName, rNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelRoleMetadataExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "service", rServiceName),
					resource.TestCheckResourceAttr(resourceName, "role", rRoleName),
					resource.TestCheckResourceAttr(resourceName, "namespace", rNamespace),
					resource.TestCheckResourceAttr(resourceName, "metadata_json", `{"id":1}`),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelRoleMetadataConfigUpdated(rServiceName, rRoleName, rNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelRoleMetadataExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "service", rServiceName),
					resource.TestCheckResourceAttr(resourceName, "role", rRoleName),
					resource.TestCheckResourceAttr(resourceName, "namespace", rNamespace),
					resource.TestCheckResourceAttr(resourceName, "metadata_json", `{"id":2}`),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMackerelRoleMetadataDestroy(s *terraform.State) error {
	client := mackerelClient()
	for _, r := range s.RootModule().Resources {
		if r.Type != "mackerel_role_metadata" {
			continue
		}

		service := r.Primary.Attributes["service"]
		role := r.Primary.Attributes["role"]
		namespace := r.Primary.Attributes["namespace"]
		if _, err := client.GetRoleMetaData(service, role, namespace); err == nil {
			return fmt.Errorf("role metadata still exists: %s:%s/%s", service, role, namespace)
		}
	}
	return nil
}

func testAccCheckMackerelRoleMetadataExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("role metadata not found resources: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no role metadata ID is set")
		}

		client := mackerelClient()
		_, err := client.GetRoleMetaData(rs.Primary.Attributes["service"], rs.Primary.Attributes["role"], rs.Primary.Attributes["namespace"])
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccMackerelRoleMetadataConfig(serviceName, roleName, namespace string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
}

resource "mackerel_role" "foo" {
  service = mackerel_service.foo.id
  name = "%s"
}

resource "mackerel_role_metadata" "foo" {
  service = mackerel_service.foo.name
  role = mackerel_role.foo.name
  namespace = "%s"
  metadata_json = jsonencode({
    id = 1
  })
}
`, serviceName, roleName, namespace)
}

func testAccMackerelRoleMetadataConfigUpdated(serviceName, roleName, namespace string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
}

resource "mackerel_role" "foo" {
  service = mackerel_service.foo.id
  name = "%s"
}

resource "mackerel_role_metadata" "foo" {
  service = mackerel_service.foo.name
  role = mackerel_role.foo.name
  namespace = "%s"
  metadata_json = jsonencode({
    id = 2
  })
}
`, serviceName, roleName, namespace)
}
