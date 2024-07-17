package mackerel

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccDataSourceMackerelService_withMemo(t *testing.T) {
	name := fmt.Sprintf("tf-service-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelServiceConfig_withMemo(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mackerel_service.foo", "id", name),
					resource.TestCheckResourceAttr("data.mackerel_service.foo", "name", name),
					resource.TestCheckResourceAttr("data.mackerel_service.foo", "memo", "This service is managed by Terraform."),
				),
			},
		},
	})
}

func testAccDataSourceMackerelServiceConfig_withMemo(name string) string {
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

func TestAccDataSourceMackerelService_noMemo(t *testing.T) {
	resourceName := "data.mackerel_service.foo"
	name := fmt.Sprintf("tf-service-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelServiceConfig_noMemo(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("memo"), knownvalue.StringExact("")),
				},
			},
		},
	})
}
func testAccDataSourceMackerelServiceConfig_noMemo(name string) string {
	return `
resource "mackerel_service" "foo" {
  name = "` + name + `"
}

data "mackerel_service" "foo" {
  name = mackerel_service.foo.id
}`
}

func TestAccDataSourceMackerelService_NotMatchAnyService(t *testing.T) {
	name := fmt.Sprintf("tf-service-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`data "mackerel_service" "foo" { name = "%s" }`, name),
				// FIXME: error message should not be tested
				ExpectError: regexp.MustCompile(fmt.Sprintf(`the name '%s' does not match any service in mackerel\.io`, name)),
			},
		},
	})
}
