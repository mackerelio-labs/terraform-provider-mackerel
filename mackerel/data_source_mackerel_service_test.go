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

func TestAccDataSourceMackerelService(t *testing.T) {
	t.Parallel()
	resourceName := "data.mackerel_service.foo"
	cases := map[string]func() []resource.TestStep{
		"withMemo": func() []resource.TestStep {
			name := fmt.Sprintf("tf-service-%s", acctest.RandString(5))
			return []resource.TestStep{{
				Config: `
resource "mackerel_service" "foo" {
  name = "` + name + `"
  memo = "This service is managed by Terraform."
}

data "mackerel_service" "foo" {
  name = mackerel_service.foo.id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "memo", "This service is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "0"),
				),
			}}
		},
		"noMemo": func() []resource.TestStep {
			name := fmt.Sprintf("tf-service-%s", acctest.RandString(5))
			return []resource.TestStep{{
				Config: `
resource "mackerel_service" "foo" {
  name = "` + name + `"
}

data "mackerel_service" "foo" {
  name = mackerel_service.foo.id
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("memo"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("roles"), knownvalue.SetExact(nil)),
				},
			}}
		},
		"not match any service": func() []resource.TestStep {
			name := fmt.Sprintf("tf-service-%s", acctest.RandString(5))
			return []resource.TestStep{{
				Config: fmt.Sprintf(`data "mackerel_service" "foo" { name = "%s" }`, name),
				// FIXME: error message should not be tested
				ExpectError: regexp.MustCompile(fmt.Sprintf(`the name '%s' does not match any service in mackerel\.io`, name)),
			}}
		},
	}

	for name, f := range cases {
		t.Run(name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t) },
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				Steps:                    f(),
			})
		})
	}
}
