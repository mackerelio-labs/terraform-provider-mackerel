package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelServiceMetadata(t *testing.T) {
	resourceName := "mackerel_service_metadata.foo"
	rand := acctest.RandString(5)
	rServiceName := fmt.Sprintf("tf-%s", rand)
	rNamespace := fmt.Sprintf("tf-namespace-%s", rand)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelServiceMetadataDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelServiceMetadataConfig(rServiceName, rNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelServiceMetadataExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "service", rServiceName),
					resource.TestCheckResourceAttr(resourceName, "namespace", rNamespace),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelServiceMetadataConfigUpdated(rServiceName, rNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelServiceMetadataExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "service", rServiceName),
					resource.TestCheckResourceAttr(resourceName, "namespace", rNamespace),
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

func testAccCheckMackerelServiceMetadataDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*mackerel.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "mackerel_service_metadata" {
			continue
		}

		if _, err := client.GetServiceMetaData(r.Primary.Attributes["service"], r.Primary.Attributes["namespace"]); err == nil {
			return fmt.Errorf("service metadata still exists")
		}
	}
	return nil
}

func testAccCheckMackerelServiceMetadataExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("service_metadata not found from resources: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no service_metadata ID is set")
		}

		client := testAccProvider.Meta().(*mackerel.Client)
		_, err := client.GetServiceMetaData(rs.Primary.Attributes["service"], rs.Primary.Attributes["namespace"])
		if err != nil {
			return fmt.Errorf("err: %s", err)
		}
		return nil
	}
}

func testAccMackerelServiceMetadataConfig(serviceName, namespace string) string {
	// language=HCL
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
	name = "%s"
}

resource "mackerel_service_metadata" "foo" {
	service = mackerel_service.foo.id
	namespace = "%s"
	metadata_json = jsonencode({
		id = 1
	})
}
`, serviceName, namespace)
}

func testAccMackerelServiceMetadataConfigUpdated(serviceName, namespace string) string {
	// language=HCL
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
	name = "%s"
}

resource "mackerel_service_metadata" "foo" {
	service = mackerel_service.foo.id
	namespace = "%s"
	metadata_json = jsonencode({
		id = 2
	})
}
`, serviceName, namespace)
}
