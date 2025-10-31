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

func Test_MackerelRoleResource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwresource.SchemaRequest{}
	resp := fwresource.SchemaResponse{}
	provider.NewMackerelRoleResource().Schema(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema: %+v", resp.Diagnostics)
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema validation: %+v", diags)
	}
}

func TestAccMackerelRole(t *testing.T) {
	resourceName := "mackerel_role.bar"
	rand := acctest.RandString(5)
	serviceName := fmt.Sprintf("tf-service-%s", rand)
	name := fmt.Sprintf("tf-role-%s", rand)
	nameUpdated := fmt.Sprintf("tf-rol-%s-updated", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		CheckDestroy:             testAccCheckMackerelRoleDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelRoleConfig(serviceName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelRoleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "service", serviceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "memo", ""),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelRoleConfigUpdated(serviceName, nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelRoleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "service", serviceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "memo", fmt.Sprintf("%s is managed by Terraform", nameUpdated)),
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

func testAccCheckMackerelRoleDestroy(s *terraform.State) error {
	client := mackerelClient()
	for _, r := range s.RootModule().Resources {
		if r.Type != "mackerel_role" {
			continue
		}

		services, err := client.FindServices()
		if err != nil {
			return err
		}
		for _, service := range services {
			if service.Name != r.Primary.Attributes["service"] {
				continue
			}
			for _, role := range service.Roles {
				if role == r.Primary.Attributes["name"] {
					return fmt.Errorf("mackerel role still exists: %s", r.Primary.ID)
				}
			}
		}
	}
	return nil
}

func testAccCheckMackerelRoleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("role not found from resources: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no role ID is set")
		}

		client := mackerelClient()
		roles, err := client.FindRoles(rs.Primary.Attributes["service"])
		if err != nil {
			return err
		}
		for _, role := range roles {
			if role.Name == rs.Primary.Attributes["name"] {
				return nil
			}
		}

		return fmt.Errorf("role not found from mackerel: %s", rs.Primary.ID)
	}
}

func testAccMackerelRoleConfig(serviceName, name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
}

resource "mackerel_role" "bar" {
  service = mackerel_service.foo.id
  name = "%s"
}
`, serviceName, name)
}

func testAccMackerelRoleConfigUpdated(serviceName, name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
}

resource "mackerel_role" "bar" {
  service = mackerel_service.foo.name
  name = "%s"
  memo = "%s is managed by Terraform"
}
`, serviceName, name, name)
}
