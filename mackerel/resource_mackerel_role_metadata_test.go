package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelRoleMetadata(t *testing.T) {
	serviceName := fmt.Sprintf("tf-service-%s", acctest.RandString(5))
	roleName := fmt.Sprintf("tf-role-%s", acctest.RandString(5))
	namespace := fmt.Sprintf("tf-nt-%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil, // todo
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMackerelRoleMetadataConfig(serviceName, roleName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelRoleMetadataExists("mackerel_role_metadata.foo"),
					resource.TestCheckResourceAttr(
						"mackerel_role_metadata.foo", "service", serviceName),
					resource.TestCheckResourceAttr(
						"mackerel_role_metadata.foo", "role", roleName),
					resource.TestCheckResourceAttr(
						"mackerel_role_metadata.foo", "namespace", namespace),
					resource.TestCheckResourceAttr(
						"mackerel_role_metadata.foo", "metadata.TZ", "UTC"),
				),
			},
		},
	})
}

func testAccCheckMackerelRoleMetadataExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("role_metadata not found from resources: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no role_metadata ID is set")
		}

		client := testAccProvider.Meta().(*mackerel.Client)
		_, err := client.GetRoleMetaData(rs.Primary.Attributes["service"], rs.Primary.Attributes["role"], rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("err: %s", err)
		}
		return nil
	}
}

func testAccCheckMackerelRoleMetadataConfig(serviceName, roleName, namespace string) string {
	// language=HCL
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
    name = "%s"
}

resource "mackerel_role" "foo" {
    service = "${mackerel_service.foo.id}"
    name = "%s"
}

resource "mackerel_role_metadata" "foo" {
    service = "${mackerel_service.foo.id}"
    role = "${mackerel_role.foo.id}"
    namespace = "%s"
    metadata = {
        TZ = "UTC"
    }
}
`, serviceName, roleName, namespace)
}
