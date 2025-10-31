package provider_test

import (
	"context"
	"fmt"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
)

func TestMackerelServiceResourceSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwresource.SchemaRequest{}
	resp := &fwresource.SchemaResponse{}
	provider.NewMackerelServiceResource().Schema(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema method diagnostics: %+v", resp.Diagnostics)
	}

	if diag := resp.Schema.ValidateImplementation(ctx); diag.HasError() {
		t.Fatalf("schema validation diagnostics: %+v", diag)
	}
}

func TestAccMackerelService(t *testing.T) {
	t.Parallel()
	resourceName := "mackerel_service.foo"
	cases := map[string]func() []resource.TestStep{
		"with memo": func() []resource.TestStep {
			config := func(name, memo string) string {
				return `
resource "mackerel_service" "foo" {
  name = "` + name + `"
  memo = "` + memo + `"
}`
			}
			rand := acctest.RandString(5)
			name := fmt.Sprintf("tf-%s", rand)
			nameUpdated := fmt.Sprintf("tf-updated-%s", rand)
			memo := fmt.Sprintf("%s is managed by Terraform.", name)
			memoUpdated := fmt.Sprintf("%s is managed by Terraform.", nameUpdated)
			return []resource.TestStep{
				// Test: Create
				{
					Config: config(name, memo),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckMackerelServiceExists(resourceName),
						resource.TestCheckResourceAttr(resourceName, "name", name),
						resource.TestCheckResourceAttr(resourceName, "memo", memo),
						resource.TestCheckResourceAttr(resourceName, "roles.#", "0"),
					),
				},
				// Test: Update
				{
					Config: config(nameUpdated, memoUpdated),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckMackerelServiceExists(resourceName),
						resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
						resource.TestCheckResourceAttr(resourceName, "memo", memoUpdated),
						resource.TestCheckResourceAttr(resourceName, "roles.#", "0"),
					),
				},
				// Test: Import
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateVerify: true,
				},
			}
		},
		"no memo": func() []resource.TestStep {
			config := func(name string) string {
				return `
resource "mackerel_service" "foo" {
  name = "` + name + `"
}`
			}
			rand := acctest.RandString(5)
			name := "tf-" + rand
			nameUpdated := "tf-updated-" + rand
			return []resource.TestStep{
				// Test: Create
				{
					Config: config(name),
					Check:  testAccCheckMackerelServiceExists(resourceName),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.StringExact(name)),
						statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(name)),
						statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("memo"), knownvalue.StringExact("")),
						statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("roles"), knownvalue.SetExact(nil)),
					},
				},
				// Test: Update
				{
					Config: config(nameUpdated),
					Check:  testAccCheckMackerelServiceExists(resourceName),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.StringExact(nameUpdated)),
						statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(nameUpdated)),
						statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("memo"), knownvalue.StringExact("")),
						statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("roles"), knownvalue.SetExact(nil)),
					},
				},
				// Test: Import
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateVerify: true,
				},
			}
		},
	}

	for name, f := range cases {
		t.Run(name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				PreCheck:                 func() { preCheck(t) },
				ProtoV5ProviderFactories: protoV5ProviderFactories,
				CheckDestroy:             testAccCheckMackerelServiceDestroy,
				Steps:                    f(),
			})
		})
	}
}

func testAccCheckMackerelServiceDestroy(s *terraform.State) error {
	client := mackerelClient()
	for _, r := range s.RootModule().Resources {
		if r.Type != "mackerel_service" {
			continue
		}

		services, err := client.FindServices()
		if err != nil {
			return err
		}
		for _, srv := range services {
			if srv.Name == r.Primary.ID {
				return fmt.Errorf("service still exists: %s", r.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckMackerelServiceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("service not found from resources: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no service ID is set")
		}

		client := mackerelClient()
		services, err := client.FindServices()
		if err != nil {
			return err
		}

		for _, srv := range services {
			if srv.Name == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("service not found from mackerel: %s", rs.Primary.ID)
	}
}
