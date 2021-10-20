package mackerel

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMackerelAWSIntegrationIAMRole(t *testing.T) {
	dsName := "data.mackerel_aws_integration.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-aws-integration-%s", rand)
	externalID := os.Getenv("EXTERNAL_ID")
	awsRoleArn := os.Getenv("AWS_ROLE_ARN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelAWSIntegrationConfigIAMRole(rand, name, awsRoleArn, externalID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "name", name),
					resource.TestCheckResourceAttr(dsName, "memo", "This aws integration is managed by Terraform."),
					resource.TestCheckResourceAttr(dsName, "key", ""),
					resource.TestCheckResourceAttr(dsName, "role_arn", awsRoleArn),
					resource.TestCheckResourceAttr(dsName, "external_id", externalID),
					resource.TestCheckResourceAttr(dsName, "region", "ap-northeast-1"),
					resource.TestCheckResourceAttr(dsName, "included_tags", "Name:staging-server,Environment:staging"),
					resource.TestCheckResourceAttr(dsName, "excluded_tags", "Name:develop-server,Environment:develop"),
					resource.TestCheckResourceAttr(dsName, "alb.#", "1"),
					resource.TestCheckResourceAttr(dsName, "rds.#", "1"),
					resource.TestCheckResourceAttr(dsName, "nlb.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceMackerelAWSIntegrationCredential(t *testing.T) {
	dsName := "data.mackerel_aws_integration.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-aws-integration-%s", rand)
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelAWSIntegrationConfigCredential(rand, name, awsAccessKeyID, awsSecretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "name", name),
					resource.TestCheckResourceAttr(dsName, "memo", "This aws integration is managed by Terraform."),
					resource.TestCheckResourceAttr(dsName, "key", awsAccessKeyID),
					resource.TestCheckResourceAttr(dsName, "role_arn", ""),
					resource.TestCheckResourceAttr(dsName, "external_id", ""),
					resource.TestCheckResourceAttr(dsName, "region", "ap-northeast-1"),
					resource.TestCheckResourceAttr(dsName, "included_tags", "Name:staging-server,Environment:staging"),
					resource.TestCheckResourceAttr(dsName, "excluded_tags", "Name:develop-server,Environment:develop"),
					resource.TestCheckResourceAttr(dsName, "alb.#", "1"),
					resource.TestCheckResourceAttr(dsName, "rds.#", "1"),
					resource.TestCheckResourceAttr(dsName, "nlb.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceMackerelAWSIntegrationConfigIAMRole(rand, name, roleArn, externalID string) string {
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
    enable           = true
    role             = "${mackerel_service.include.name}: ${mackerel_role.include.name}"
    excluded_metrics = ["rds.cpu.used"]
  }

  nlb {
    enable = false
  }
}

data "mackerel_aws_integration" "foo" {
  id = mackerel_aws_integration.foo.id
}
`, rand, rand, name, roleArn, externalID)
}

func testAccDataSourceMackerelAWSIntegrationConfigCredential(rand, name, awsAccessKeyID, awsSecretAccessKey string) string {
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

data "mackerel_aws_integration" "foo" {
  id = mackerel_aws_integration.foo.id
}
`, rand, rand, name, awsAccessKeyID, awsSecretAccessKey)
}
