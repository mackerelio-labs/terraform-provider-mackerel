package mackerel

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mackerelio/mackerel-client-go"
	"os"
	"testing"
)

func TestAccMackerelAWSIntegrationIAMRole(t *testing.T) {
	resourceName := "mackerel_aws_integration.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-aws-integration-%s", rand)
	nameUpdated := fmt.Sprintf("tf-aws-integration-%s-updated", rand)
	externalID := os.Getenv("EXTERNAL_ID")
	awsRoleArn := os.Getenv("AWS_ROLE_ARN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMackerelAWSIntegrationDestroy,
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
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsSecretKey := os.Getenv("AWS_SECRET_KEY")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMackerelAWSIntegrationDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccSourceMackerelAWSIntegrationConfigCredential(rand, name, awsAccessKey, awsSecretKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelAWSIntegrationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "memo", "This aws integration is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "key", awsAccessKey),
					resource.TestCheckResourceAttr(resourceName, "role_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "external_id", ""),
					resource.TestCheckResourceAttr(resourceName, "region", "ap-northeast-1"),
					resource.TestCheckResourceAttr(resourceName, "included_tags", "Name:staging-server,Environment:staging"),
					resource.TestCheckResourceAttr(resourceName, "excluded_tags", "Name:develop-server,Environment:develop"),
				),
			},
			// Test: Update
			{
				Config: testAccSourceMackerelAWSIntegrationConfigCredentialUpdate(rand, nameUpdated, awsAccessKey, awsSecretKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelAWSIntegrationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "memo", "This aws integration is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "key", awsAccessKey),
					resource.TestCheckResourceAttr(resourceName, "role_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "external_id", ""),
					resource.TestCheckResourceAttr(resourceName, "region", "ap-northeast-1"),
					resource.TestCheckResourceAttr(resourceName, "included_tags", "Name:production-server,Environment:production"),
					resource.TestCheckResourceAttr(resourceName, "excluded_tags", "Name:staging-server,Environment:staging"),
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
	client := testAccProvider.Meta().(*mackerel.Client)
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

		client := testAccProvider.Meta().(*mackerel.Client)
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
  name = "tf-role-%s-include"
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
    enable           = true
	role             = "${mackerel_service.include.name}: ${mackerel_role.include.name}"
    excluded_metrics = ["rds.cpu.used"]
 }

  nlb {
    enable           = true
	role             = "${mackerel_service.include.name}: ${mackerel_role.include.name}"
    excluded_metrics = []
  }
}
`, rand, rand, name, roleArn, externalID)
}

func testAccSourceMackerelAWSIntegrationConfigCredential(rand, name, awsAccessKey, awsSecretKey string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "include" {
  name = "tf-service-%s-include"
}

resource "mackerel_role" "include" {
  service = mackerel_service.include.name
  name = "tf-role-%s-include"
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
    enable           = true
	role             = "${mackerel_service.include.name}: ${mackerel_role.include.name}"
    excluded_metrics = []
  }
}
`, rand, rand, name, awsAccessKey, awsSecretKey)
}

func testAccSourceMackerelAWSIntegrationConfigIAMRoleUpdate(rand, name, roleArn, externalID string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "include" {
  name = "tf-service-%s-include"
}

resource "mackerel_role" "include" {
  service = mackerel_service.include.name
  name = "tf-role-%s-include"
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
    enable           = true
	role             = ""
    excluded_metrics = []
	retire_automatically = true
  }

  lambda {
    enable = true
  }
}
`, rand, rand, name, roleArn, externalID)
}

func testAccSourceMackerelAWSIntegrationConfigCredentialUpdate(rand, name, awsAccessKey, awsSecretKey string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "include" {
  name = "tf-service-%s-include"
}

resource "mackerel_role" "include" {
  service = mackerel_service.include.name
  name = "tf-role-%s-include"
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
    enable           = true
	role             = "${mackerel_service.include.name}: ${mackerel_role.include.name}"
    excluded_metrics = ["rds.cpu.used"]
 }

  nlb {
    enable           = true
    excluded_metrics = []
  }

  ec2 {
    enable           = true
	role             = ""
    excluded_metrics = []
	retire_automatically = true
  }

  lambda {
    enable = true
  }
}
`, rand, rand, name, awsAccessKey, awsSecretKey)
}
