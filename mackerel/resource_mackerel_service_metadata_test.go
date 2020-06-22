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
	serviceName := fmt.Sprintf("tf-service-%s", acctest.RandString(5))
	namespace := fmt.Sprintf("tf-ns-%s", serviceName)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil, // todo
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMackerelServiceMetadataConfig(serviceName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelServiceMetadataExists("mackerel_service_metadata.foo"),
					resource.TestCheckResourceAttr(
						"mackerel_service_metadata.foo", "service", serviceName),
					resource.TestCheckResourceAttr(
						"mackerel_service_metadata.foo", "namespace", namespace),
				),
			},
		},
	})
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
		_, err := client.GetServiceMetaData(rs.Primary.Attributes["service"], rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("err: %s", err)
		}
		return nil
	}
}

func testAccCheckMackerelServiceMetadataConfig(serviceName, namespace string) string {
	// language=HCL
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
	name = "%s"
}

resource "mackerel_service_metadata" "foo" {
	service = "${mackerel_service.foo.id}"
	namespace = "%s"
	metadata_json = jsonencode({
		int = 1
		string = "foo bar baz"
		array = ["1", true, 1]
	})
}
`, serviceName, namespace)
}
