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

func Test_MackerelServiceMetadataResource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwresource.SchemaRequest{}
	resp := &fwresource.SchemaResponse{}
	provider.NewMackerelServiceMetadataResource().Schema(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Errorf("schema method: %+v", resp.Diagnostics)
		return
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Errorf("schema validation: %+v", diags)
	}
}

func TestAccMackerelServiceMetadata(t *testing.T) {
	resourceName := "mackerel_service_metadata.foo"
	rand := acctest.RandString(5)
	serviceName := fmt.Sprintf("tf-%s", rand)
	namespace := fmt.Sprintf("tf-namespace-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		CheckDestroy:             testAccCheckMackerelServiceMetadataDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelServiceMetadataConfig(serviceName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelServiceMetadataExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "service", serviceName),
					resource.TestCheckResourceAttr(resourceName, "namespace", namespace),
					resource.TestCheckResourceAttr(resourceName, "metadata_json", `{"id":1}`),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelServiceMetadataConfigUpdated(serviceName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelServiceMetadataExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "service", serviceName),
					resource.TestCheckResourceAttr(resourceName, "namespace", namespace),
					resource.TestCheckResourceAttr(resourceName, "metadata_json", `{"id":2}`),
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
	client := mackerelClient()
	for _, r := range s.RootModule().Resources {
		if r.Type != "mackerel_service_metadata" {
			continue
		}

		service := r.Primary.Attributes["service"]
		namespace := r.Primary.Attributes["namespace"]
		if _, err := client.GetServiceMetaData(service, namespace); err == nil {
			return fmt.Errorf("service metadata still exists: %s:%s", service, namespace)
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

		client := mackerelClient()
		_, err := client.GetServiceMetaData(rs.Primary.Attributes["service"], rs.Primary.Attributes["namespace"])
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccMackerelServiceMetadataConfig(serviceName, namespace string) string {
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
