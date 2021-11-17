package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMackerelServiceMetricNames(t *testing.T) {
	name := fmt.Sprintf("tf-service-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelServiceMetricNamesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mackerel_service_metric_names.foo", "id", name+":"),
					resource.TestCheckResourceAttr("data.mackerel_service_metric_names.foo", "name", name),
				),
			},
		},
	})
}

func testAccDataSourceMackerelServiceMetricNamesConfig(name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
  memo = "This service is managed by Terraform."
}

data "mackerel_service_metric_names" "foo" {
  name = mackerel_service.foo.name
}
`, name)
}
