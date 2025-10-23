package provider_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

func TestAccMackerelAWSIntegrationIAMRole(t *testing.T) {
	resourceName := "mackerel_aws_integration.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-aws-integration-%s", rand)
	nameUpdated := fmt.Sprintf("tf-aws-integration-%s-updated", rand)

	externalID := os.Getenv("EXTERNAL_ID")
	if externalID == "" {
		t.Skip("EXTERNAL_ID must be set for acceptance tests")
	}
	awsRoleArn := os.Getenv("AWS_ROLE_ARN")
	if awsRoleArn == "" {
		t.Skip("AWS_ROLE_ARN must be set for acceptance tests")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		CheckDestroy:             testAccCheckMackerelAWSIntegrationDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccSourceMackerelAWSIntegrationConfigIAMRole(rand, name, awsRoleArn, externalID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelAWSIntegrationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "memo", "This aws integration is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "key", ""),
					resource.TestCheckResourceAttr(resourceName, "role_arn", awsRoleArn),
					resource.TestCheckResourceAttr(resourceName, "external_id", externalID),
					resource.TestCheckResourceAttr(resourceName, "region", "ap-northeast-1"),
					resource.TestCheckResourceAttr(resourceName, "included_tags", "Name:staging-server,Environment:staging"),
					resource.TestCheckResourceAttr(resourceName, "excluded_tags", "Name:develop-server,Environment:develop"),
					resource.TestCheckResourceAttr(resourceName, "alb.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rds.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "nlb.#", "1"),
				),
			},
			// Test: Update
			{
				Config: testAccSourceMackerelAWSIntegrationConfigIAMRoleUpdate(rand, nameUpdated, awsRoleArn, externalID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelAWSIntegrationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "memo", "This aws integration is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "key", ""),
					resource.TestCheckResourceAttr(resourceName, "role_arn", awsRoleArn),
					resource.TestCheckResourceAttr(resourceName, "external_id", externalID),
					resource.TestCheckResourceAttr(resourceName, "region", "ap-northeast-1"),
					resource.TestCheckResourceAttr(resourceName, "included_tags", "Name:production-server,Environment:production"),
					resource.TestCheckResourceAttr(resourceName, "excluded_tags", "Name:staging-server,Environment:staging"),
					resource.TestCheckResourceAttr(resourceName, "alb.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rds.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "nlb.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ec2.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "lambda.#", "1"),
				),
			},
			// Test: Import
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_key"},
			},
		},
	})
}

func TestAccMackerelAWSIntegrationCredentials(t *testing.T) {
	resourceName := "mackerel_aws_integration.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-aws-integration-%s", rand)
	nameUpdated := fmt.Sprintf("tf-aws-integration-%s-updated", rand)

	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	if awsAccessKeyID == "" {
		t.Skip("AWS_ACCESS_KEY_ID must be set for acceptance tests")
	}
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if awsSecretAccessKey == "" {
		t.Skip("AWS_SECRET_ACCESS_KEY must be set for acceptance tests")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		CheckDestroy:             testAccCheckMackerelAWSIntegrationDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccSourceMackerelAWSIntegrationConfigCredential(rand, name, awsAccessKeyID, awsSecretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelAWSIntegrationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "memo", "This aws integration is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "key", awsAccessKeyID),
					resource.TestCheckResourceAttr(resourceName, "role_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "external_id", ""),
					resource.TestCheckResourceAttr(resourceName, "region", "ap-northeast-1"),
					resource.TestCheckResourceAttr(resourceName, "included_tags", "Name:staging-server,Environment:staging"),
					resource.TestCheckResourceAttr(resourceName, "excluded_tags", "Name:develop-server,Environment:develop"),
					resource.TestCheckResourceAttr(resourceName, "alb.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rds.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "nlb.#", "1"),
				),
			},
			// Test: Update
			{
				Config: testAccSourceMackerelAWSIntegrationConfigCredentialUpdate(rand, nameUpdated, awsAccessKeyID, awsSecretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelAWSIntegrationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "memo", "This aws integration is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "key", awsAccessKeyID),
					resource.TestCheckResourceAttr(resourceName, "role_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "external_id", ""),
					resource.TestCheckResourceAttr(resourceName, "region", "ap-northeast-1"),
					resource.TestCheckResourceAttr(resourceName, "included_tags", "Name:production-server,Environment:production"),
					resource.TestCheckResourceAttr(resourceName, "excluded_tags", "Name:staging-server,Environment:staging"),
					resource.TestCheckResourceAttr(resourceName, "alb.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rds.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "nlb.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ec2.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "lambda.#", "1"),
				),
			},
			// Test: Import
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_key"},
			},
		},
	})
}

func testAccCheckMackerelAWSIntegrationDestroy(s *terraform.State) error {
	client := mackerelClient()
	for _, r := range s.RootModule().Resources {
		if r.Type != "mackerel_aws_integration" {
			continue
		}

		integrations, err := client.FindAWSIntegrations()
		if err != nil {
			return err
		}
		for _, integration := range integrations {
			if integration.ID == r.Primary.ID {
				return fmt.Errorf("aws integration still exists: %s", r.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckMackerelAWSIntegrationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("aws integration not found from resources: %s", n)
		}

		if r.Primary.ID == "" {
			return fmt.Errorf("no aws integration ID is set")
		}

		client := mackerelClient()
		integrations, err := client.FindAWSIntegrations()
		if err != nil {
			return err
		}
		for _, integration := range integrations {
			if integration.ID == r.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("aws integration not found from mackerel: %s", r.Primary.ID)
	}
}

func testAccSourceMackerelAWSIntegrationConfigIAMRole(rand, name, roleArn, externalID string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "include" {
  name = "tf-service-%s-include"
}

resource "mackerel_role" "include" {
  service = mackerel_service.include.name
  name    = "tf-role-%s-include"
}

resource "mackerel_aws_integration" "foo" {
  name          = "%s"
  memo          = "This aws integration is managed by Terraform."
  key           = ""
  secret_key    = ""
  role_arn      = "%s"
  external_id   = "%s"
  region        = "ap-northeast-1"
  included_tags = "Name:staging-server,Environment:staging"
  excluded_tags = "Name:develop-server,Environment:develop"

  alb {
    enable           = true
    role             = "${mackerel_service.include.name}: ${mackerel_role.include.name}"
    excluded_metrics = ["alb.request.count", "alb.bytes.processed"]
  }

  rds {
    role             = "${mackerel_service.include.name}: ${mackerel_role.include.name}"
    excluded_metrics = ["rds.cpu.used"]
  }

  nlb {
    enable = false
  }
}
`, rand, rand, name, roleArn, externalID)
}

func testAccSourceMackerelAWSIntegrationConfigCredential(rand, name, awsAccessKeyID, awsSecretAccessKey string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "include" {
  name = "tf-service-%s-include"
}

resource "mackerel_role" "include" {
  service = mackerel_service.include.name
  name    = "tf-role-%s-include"
}

resource "mackerel_aws_integration" "foo" {
  name          = "%s"
  memo          = "This aws integration is managed by Terraform."
  key           = "%s"
  secret_key    = "%s"
  role_arn      = ""
  external_id   = ""
  region        = "ap-northeast-1"
  included_tags = "Name:staging-server,Environment:staging"
  excluded_tags = "Name:develop-server,Environment:develop"

  alb {
    enable           = true
    role             = "${mackerel_service.include.name}: ${mackerel_role.include.name}"
    excluded_metrics = ["alb.request.count", "alb.bytes.processed"]
  }

  rds {
    enable           = true
    role             = "${mackerel_service.include.name}: ${mackerel_role.include.name}"
    excluded_metrics = ["rds.cpu.used"]
  }

  nlb {
    enable = false
  }
}
`, rand, rand, name, awsAccessKeyID, awsSecretAccessKey)
}

func testAccSourceMackerelAWSIntegrationConfigIAMRoleUpdate(rand, name, roleArn, externalID string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "include" {
  name = "tf-service-%s-include"
}

resource "mackerel_role" "include" {
  service = mackerel_service.include.name
  name    = "tf-role-%s-include"
}

resource "mackerel_aws_integration" "foo" {
  name          = "%s"
  memo          = "This aws integration is managed by Terraform."
  key           = ""
  secret_key    = ""
  role_arn      = "%s"
  external_id   = "%s"
  region        = "ap-northeast-1"
  included_tags = "Name:production-server,Environment:production"
  excluded_tags = "Name:staging-server,Environment:staging"

  alb {
    enable           = true
    role             = "${mackerel_service.include.name}: ${mackerel_role.include.name}"
    excluded_metrics = ["alb.request.count", "alb.bytes.processed"]
  }

  rds {
    enable           = true
    role             = "${mackerel_service.include.name}: ${mackerel_role.include.name}"
    excluded_metrics = ["rds.cpu.used"]
  }

  nlb {
    enable           = true
    excluded_metrics = []
  }

  ec2 {
    enable               = true
    excluded_metrics     = []
    retire_automatically = true
  }

  lambda {
    enable = true
    retire_automatically = false
  }
}
`, rand, rand, name, roleArn, externalID)
}

func testAccSourceMackerelAWSIntegrationConfigCredentialUpdate(rand, name, awsAccessKeyID, awsSecretAccessKey string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "include" {
  name = "tf-service-%s-include"
}

resource "mackerel_role" "include" {
  service = mackerel_service.include.name
  name    = "tf-role-%s-include"
}

resource "mackerel_aws_integration" "foo" {
  name          = "%s"
  memo          = "This aws integration is managed by Terraform."
  key           = "%s"
  secret_key    = "%s"
  role_arn      = ""
  external_id   = ""
  region        = "ap-northeast-1"
  included_tags = "Name:production-server,Environment:production"
  excluded_tags = "Name:staging-server,Environment:staging"

  alb {
    enable           = true
    role             = "${mackerel_service.include.name}: ${mackerel_role.include.name}"
    excluded_metrics = ["alb.request.count", "alb.bytes.processed"]
  }

  rds {
    enable               = true
    role                 = "${mackerel_service.include.name}: ${mackerel_role.include.name}"
    excluded_metrics     = ["rds.cpu.used"]
    retire_automatically = true
  }

  nlb {
    enable           = true
    excluded_metrics = []
  }

  ec2 {
    enable               = true
    excluded_metrics     = []
    retire_automatically = true
  }

  lambda {
    enable = true
    retire_automatically = false
  }
}
`, rand, rand, name, awsAccessKeyID, awsSecretAccessKey)
}
