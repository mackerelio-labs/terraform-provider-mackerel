package mackerel

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceMackerelServiceMetadata(t *testing.T) {
	dsName := "data.mackerel_service_metadata.foo"
	rand := acctest.RandString(5)
	serviceName := fmt.Sprintf("tf-service-%s", rand)
	namespace := fmt.Sprintf("tf-namespace-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelServiceMetadataConfig(serviceName, namespace),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsName, "service", serviceName),
					resource.TestCheckResourceAttr(dsName, "namespace", namespace),
					resource.TestCheckResourceAttr(dsName, "metadata_json", `{"id":1}`),
				),
			},
		},
	})
}

func testAccDataSourceMackerelServiceMetadataConfig(serviceName, namespace string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
}

resource "mackerel_service_metadata" "foo" {
  service = mackerel_service.foo.name
  namespace = "%s"
  metadata_json = jsonencode({
    id = 1
  })
}

data "mackerel_service_metadata" "foo" {
  service = mackerel_service_metadata.foo.service
  namespace = mackerel_service_metadata.foo.namespace
}
`, serviceName, namespace)
}

func TestAccDataSourceMackerelServiceMetadata_NoMatchReturnsError(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceMackerelServiceMetadataConfigNoMatchReturnsError,
				ExpectError: regexp.MustCompile(`API request failed: Service not found.`),
			},
		},
	})
}

const testAccDataSourceMackerelServiceMetadataConfigNoMatchReturnsError = `
data "mackerel_service_metadata" "foo" {
  service = "not-found"
  namespace = "not-found"
}
`
