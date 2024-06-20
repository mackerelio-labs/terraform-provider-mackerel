package provider_test

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
	"github.com/mackerelio/mackerel-client-go"
)

func Test_MackerelServiceResource_schema(t *testing.T) {
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

func TestAcc_MackerelServiceResource_basic(t *testing.T) {
	client := newClient(t)
	resourceName := "mackerel_service.foo"
	name := acctest.RandomWithPrefix("tf")
	nameUpdated := acctest.RandomWithPrefix("tf-updated")
	memoUpdated := fmt.Sprintf("%s is managed by Terraform.", nameUpdated)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             checkMackerelServiceDestroy(client),
		Steps: []resource.TestStep{
			// Create
			{
				Config: mackerelServiceResourceConfig_basic(name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("id")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(name)),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("memo"), knownvalue.Null()),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					expectMackerelService(client, resourceName, mackerel.Service{
						Name: name,
						Memo: "",
					}),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("memo"), knownvalue.Null()),
				},
			},
			// Update
			{
				Config: mackerelServiceResourceConfig_withMemo(nameUpdated, memoUpdated),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("id")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(nameUpdated)),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("memo"), knownvalue.StringExact(memoUpdated)),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					expectMackerelService(client, resourceName, mackerel.Service{
						Name: nameUpdated,
						Memo: memoUpdated,
					}),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.StringExact(nameUpdated)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(nameUpdated)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("memo"), knownvalue.StringExact(memoUpdated)),
				},
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func mackerelServiceResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
}`, name)
}

func mackerelServiceResourceConfig_withMemo(name string, memo string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
  memo = %q
}`, name, memo)
}

func checkMackerelServiceDestroy(client *mackerel.Client) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		services, err := client.FindServices()
		if err != nil {
			return err
		}
		for _, r := range s.RootModule().Resources {
			if r.Type != "mackerel_service" {
				continue
			}
			if slices.ContainsFunc(services, func(svc *mackerel.Service) bool {
				return svc.Name == r.Primary.ID
			}) {
				return fmt.Errorf("service still exists: %s", r.Primary.ID)
			}
		}
		return nil
	}
}

func expectMackerelService(client *mackerel.Client, resourceAddr string, wants mackerel.Service) statecheck.StateCheck {
	return stateCheckFunc(func(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
		resource, err := findStateResource(req.State, resourceAddr)
		if err != nil {
			resp.Error = err
			return
		}
		idAny, err := tfjsonpath.Traverse(resource.AttributeValues, tfjsonpath.New("id"))
		if err != nil {
			resp.Error = err
			return
		}
		id, ok := idAny.(string)
		if !ok {
			resp.Error = fmt.Errorf("id must be a string")
			return
		}

		services, err := client.FindServices()
		if err != nil {
			resp.Error = err
			return
		}

		serviceIdx := slices.IndexFunc(services, func(svc *mackerel.Service) bool {
			return svc.Name == id
		})
		if serviceIdx < 0 {
			resp.Error = fmt.Errorf("service not found from mackerel: %s", id)
			return
		}

		if diff := cmp.Diff(
			*services[serviceIdx],
			wants,
			cmpopts.IgnoreFields(mackerel.Service{}, "Roles"),
		); diff != "" {
			resp.Error = fmt.Errorf("unexpected service. diff: %s", diff)
			return
		}
	})
}
