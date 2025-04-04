package provider_test

import (
	"context"
	"os"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
)

func Test_MackerelAWSIntegrationResource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwresource.SchemaRequest{}
	resp := fwresource.SchemaResponse{}
	if provider.NewMackerelAWSIntegrationResource().Schema(ctx, req, &resp); resp.Diagnostics.HasError() {
		t.Fatalf("schema method diagnostics: %+v", resp.Diagnostics)
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema validation diagnostics: %+v", diags)
	}
}

func TestAccCompat_MackerelAWSIntegrationResource(t *testing.T) {
	t.Parallel()

	roleARN := os.Getenv("AWS_ROLE_ARN")
	if roleARN == "" {
		t.Skip("AWS_ROLE_ARN must be set for acceptance tests")
	}
	externalID := os.Getenv("EXTERNAL_ID")
	if externalID == "" {
		t.Skip("EXTERNAL_ID must be set for acceptance tests")
	}

	cases := map[string]struct {
		config func(name string) string
	}{
		"basic": {
			config: func(name string) string {
				return `
resource "mackerel_service" "service" {
  name = "` + name + `-service"
}
resource "mackerel_role" "role" {
  service = mackerel_service.service.name
  name = "` + name + `-role"
}
resource "mackerel_aws_integration" "aws_integration" {
  name = "` + name + `"
  role_arn = "` + roleARN + `"
  external_id = "` + externalID + `"
  region = "ap-northeast-1"
  included_tags = "Name:staging-server,Environment:staging"
  excluded_tags = "Name:develop-server,Environment:develop"

  ec2 {
    enable = true
    role = "${mackerel_service.service.name}: ${mackerel_role.role.name}"
  }
  alb {
    enable = true
    role = "${mackerel_service.service.name}: ${mackerel_role.role.name}"
  }
}`
			},
		},
	}

	for testName, tt := range cases {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			name := acctest.RandomWithPrefix("tf-test-compat-aws-integration")
			config := tt.config(name)

			resource.Test(t, resource.TestCase{
				PreCheck: func() { preCheck(t) },

				Steps: []resource.TestStep{
					{
						ProtoV5ProviderFactories: protoV5SDKProviderFactories,
						Config:                   config,
					},
					stepNoPlanInFramework(config),
					{
						ProtoV5ProviderFactories: protoV5SDKProviderFactories,
						Config:                   config,
					},
					stepNoPlanInFramework(config),
				},
			})
		})
	}
}
