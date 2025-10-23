package provider_test

import (
	"context"
	"fmt"
	"testing"

	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
)

func Test_MackerelMonitorDataSource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	req := fwdatasource.SchemaRequest{}
	resp := fwdatasource.SchemaResponse{}
	provider.NewMackerelMonitorDataSource().Schema(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %+v", resp.Diagnostics)
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema validation diagnostics: %+v", diags)
	}
}

func TestAccDataSourceMackerelMonitorHostMetric(t *testing.T) {
	dsName := "data.mackerel_monitor.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-monitor-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelMonitorConfigHostMetric(rand, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "name", name),
					resource.TestCheckResourceAttr(dsName, "memo", "This monitor is managed by Terraform."),
					resource.TestCheckResourceAttr(dsName, "is_mute", "true"),
					resource.TestCheckResourceAttr(dsName, "notification_interval", "30"),
					resource.TestCheckResourceAttr(dsName, "host_metric.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dsName, "host_metric.0.metric", "disk%"),
						resource.TestCheckResourceAttr(dsName, "host_metric.0.operator", ">"),
						resource.TestCheckResourceAttr(dsName, "host_metric.0.warning", "70"),
						resource.TestCheckResourceAttr(dsName, "host_metric.0.critical", "90"),
						resource.TestCheckResourceAttr(dsName, "host_metric.0.duration", "3"),
						resource.TestCheckResourceAttr(dsName, "host_metric.0.max_check_attempts", "5"),
						resource.TestCheckResourceAttr(dsName, "host_metric.0.scopes.#", "2"),
						resource.TestCheckResourceAttr(dsName, "host_metric.0.exclude_scopes.#", "2"),
					),
					resource.TestCheckResourceAttr(dsName, "connectivity.#", "0"),
					resource.TestCheckResourceAttr(dsName, "service_metric.#", "0"),
					resource.TestCheckResourceAttr(dsName, "external.#", "0"),
					resource.TestCheckResourceAttr(dsName, "expression.#", "0"),
					resource.TestCheckResourceAttr(dsName, "anomaly_detection.#", "0"),
					resource.TestCheckResourceAttr(dsName, "query.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceMackerelMonitorConnectivity(t *testing.T) {
	dsName := "data.mackerel_monitor.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-monitor-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelMonitorConfigConnectivity(rand, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "name", name),
					resource.TestCheckResourceAttr(dsName, "memo", "This monitor is managed by Terraform."),
					resource.TestCheckResourceAttr(dsName, "is_mute", "true"),
					resource.TestCheckResourceAttr(dsName, "notification_interval", "30"),
					resource.TestCheckResourceAttr(dsName, "host_metric.#", "0"),
					resource.TestCheckResourceAttr(dsName, "connectivity.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dsName, "connectivity.0.scopes.#", "2"),
						resource.TestCheckResourceAttr(dsName, "connectivity.0.exclude_scopes.#", "2"),
						resource.TestCheckResourceAttr(dsName, "connectivity.0.alert_status_on_gone", "WARNING"),
					),
					resource.TestCheckResourceAttr(dsName, "service_metric.#", "0"),
					resource.TestCheckResourceAttr(dsName, "external.#", "0"),
					resource.TestCheckResourceAttr(dsName, "expression.#", "0"),
					resource.TestCheckResourceAttr(dsName, "anomaly_detection.#", "0"),
					resource.TestCheckResourceAttr(dsName, "query.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceMackerelMonitorServiceMetric(t *testing.T) {
	dsName := "mackerel_monitor.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-monitor service_metric %s", rand)
	serviceName := fmt.Sprintf("tf-service-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelMonitorConfigServiceMetric(serviceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "name", name),
					resource.TestCheckResourceAttr(dsName, "memo", "This monitor is managed by Terraform."),
					resource.TestCheckResourceAttr(dsName, "is_mute", "true"),
					resource.TestCheckResourceAttr(dsName, "notification_interval", "30"),
					resource.TestCheckResourceAttr(dsName, "host_metric.#", "0"),
					resource.TestCheckResourceAttr(dsName, "connectivity.#", "0"),
					resource.TestCheckResourceAttr(dsName, "service_metric.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dsName, "service_metric.0.service", serviceName),
						resource.TestCheckResourceAttr(dsName, "service_metric.0.metric", "custom.access.2xx_ratio"),
						resource.TestCheckResourceAttr(dsName, "service_metric.0.operator", "<"),
						resource.TestCheckResourceAttr(dsName, "service_metric.0.warning", "99.9"),
						resource.TestCheckResourceAttr(dsName, "service_metric.0.critical", "99.99"),
						resource.TestCheckResourceAttr(dsName, "service_metric.0.duration", "3"),
						resource.TestCheckResourceAttr(dsName, "service_metric.0.max_check_attempts", "5"),
						resource.TestCheckResourceAttr(dsName, "service_metric.0.missing_duration_warning", "10"),
						resource.TestCheckResourceAttr(dsName, "service_metric.0.missing_duration_critical", "10080"),
					),
					resource.TestCheckResourceAttr(dsName, "external.#", "0"),
					resource.TestCheckResourceAttr(dsName, "expression.#", "0"),
					resource.TestCheckResourceAttr(dsName, "anomaly_detection.#", "0"),
					resource.TestCheckResourceAttr(dsName, "query.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceMackerelMonitorExternal(t *testing.T) {
	dsName := "data.mackerel_monitor.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-monitor-%s", rand)
	serviceName := fmt.Sprintf("tf-service-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelMonitorConfigExternal(serviceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "name", name),
					resource.TestCheckResourceAttr(dsName, "memo", "This monitor is managed by Terraform."),
					resource.TestCheckResourceAttr(dsName, "is_mute", "true"),
					resource.TestCheckResourceAttr(dsName, "notification_interval", "30"),
					resource.TestCheckResourceAttr(dsName, "host_metric.#", "0"),
					resource.TestCheckResourceAttr(dsName, "connectivity.#", "0"),
					resource.TestCheckResourceAttr(dsName, "service_metric.#", "0"),
					resource.TestCheckResourceAttr(dsName, "external.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dsName, "external.0.method", "POST"),
						resource.TestCheckResourceAttr(dsName, "external.0.url", "https://terraform-provider-mackerel.test/"),
						resource.TestCheckResourceAttr(dsName, "external.0.max_check_attempts", "3"),
						resource.TestCheckResourceAttr(dsName, "external.0.service", serviceName),
						resource.TestCheckResourceAttr(dsName, "external.0.response_time_critical", "3000"),
						resource.TestCheckResourceAttr(dsName, "external.0.response_time_warning", "2000"),
						resource.TestCheckResourceAttr(dsName, "external.0.response_time_duration", "3"),
						resource.TestCheckResourceAttr(dsName, "external.0.request_body", "foo=bar"),
						resource.TestCheckResourceAttr(dsName, "external.0.contains_string", "blah blah blah"),
						resource.TestCheckResourceAttr(dsName, "external.0.certification_expiration_critical", "7"),
						resource.TestCheckResourceAttr(dsName, "external.0.certification_expiration_warning", "14"),
						resource.TestCheckResourceAttr(dsName, "external.0.skip_certificate_verification", "true"),
						resource.TestCheckResourceAttr(dsName, "external.0.headers.%", "1"),
						resource.TestCheckResourceAttr(dsName, "external.0.headers.Cache-Control", "no-cache"),
						resource.TestCheckResourceAttr(dsName, "external.0.follow_redirect", "true"),
						resource.TestCheckResourceAttr(dsName, "external.0.expected_status_code", "200"),
					),
					resource.TestCheckResourceAttr(dsName, "expression.#", "0"),
					resource.TestCheckResourceAttr(dsName, "anomaly_detection.#", "0"),
					resource.TestCheckResourceAttr(dsName, "query.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceMackerelMonitorExpression(t *testing.T) {
	dsName := "data.mackerel_monitor.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-monitor-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelMonitorExpression(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "name", name),
					resource.TestCheckResourceAttr(dsName, "memo", "This monitor is managed by Terraform."),
					resource.TestCheckResourceAttr(dsName, "is_mute", "true"),
					resource.TestCheckResourceAttr(dsName, "notification_interval", "30"),
					resource.TestCheckResourceAttr(dsName, "host_metric.#", "0"),
					resource.TestCheckResourceAttr(dsName, "connectivity.#", "0"),
					resource.TestCheckResourceAttr(dsName, "service_metric.#", "0"),
					resource.TestCheckResourceAttr(dsName, "external.#", "0"),
					resource.TestCheckResourceAttr(dsName, "expression.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dsName, "expression.0.expression", "max(role(my-service:db, loadavg5))"),
						resource.TestCheckResourceAttr(dsName, "expression.0.operator", ">"),
						resource.TestCheckResourceAttr(dsName, "expression.0.warning", "0.7"),
						resource.TestCheckResourceAttr(dsName, "expression.0.critical", "0.9"),
						resource.TestCheckResourceAttr(dsName, "expression.0.evaluate_backward_minutes", "3"),
					),
					resource.TestCheckResourceAttr(dsName, "anomaly_detection.#", "0"),
					resource.TestCheckResourceAttr(dsName, "query.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceMackerelMonitorAnomalyDetection(t *testing.T) {
	dsName := "data.mackerel_monitor.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-monitor-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelMonitorConfigAnomalyDetection(rand, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "name", name),
					resource.TestCheckResourceAttr(dsName, "memo", "This monitor is managed by Terraform."),
					resource.TestCheckResourceAttr(dsName, "is_mute", "true"),
					resource.TestCheckResourceAttr(dsName, "notification_interval", "30"),
					resource.TestCheckResourceAttr(dsName, "host_metric.#", "0"),
					resource.TestCheckResourceAttr(dsName, "connectivity.#", "0"),
					resource.TestCheckResourceAttr(dsName, "service_metric.#", "0"),
					resource.TestCheckResourceAttr(dsName, "external.#", "0"),
					resource.TestCheckResourceAttr(dsName, "expression.#", "0"),
					resource.TestCheckResourceAttr(dsName, "anomaly_detection.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dsName, "anomaly_detection.0.warning_sensitivity", "insensitive"),
						resource.TestCheckResourceAttr(dsName, "anomaly_detection.0.critical_sensitivity", "normal"),
						resource.TestCheckResourceAttr(dsName, "anomaly_detection.0.max_check_attempts", "5"),
						resource.TestCheckResourceAttr(dsName, "anomaly_detection.0.training_period_from", "1577836800"),
						resource.TestCheckResourceAttr(dsName, "anomaly_detection.0.scopes.#", "1"),
					),
					resource.TestCheckResourceAttr(dsName, "query.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceMackerelMonitorQuery(t *testing.T) {
	dsName := "data.mackerel_monitor.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-monitor-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelMonitorConfigQuery(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "name", name),
					resource.TestCheckResourceAttr(dsName, "memo", "This monitor is managed by Terraform."),
					resource.TestCheckResourceAttr(dsName, "is_mute", "true"),
					resource.TestCheckResourceAttr(dsName, "notification_interval", "30"),
					resource.TestCheckResourceAttr(dsName, "host_metric.#", "0"),
					resource.TestCheckResourceAttr(dsName, "connectivity.#", "0"),
					resource.TestCheckResourceAttr(dsName, "service_metric.#", "0"),
					resource.TestCheckResourceAttr(dsName, "external.#", "0"),
					resource.TestCheckResourceAttr(dsName, "expression.#", "0"),
					resource.TestCheckResourceAttr(dsName, "anomaly_detection.#", "0"),
					resource.TestCheckResourceAttr(dsName, "query.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dsName, "query.0.query", "container.cpu.utilization{k8s.deployment.name=\"httpbin\"}"),
						resource.TestCheckResourceAttr(dsName, "query.0.legend", "cpu.utilization {{k8s.node.name}}"),
						resource.TestCheckResourceAttr(dsName, "query.0.operator", ">"),
						resource.TestCheckResourceAttr(dsName, "query.0.warning", "70"),
						resource.TestCheckResourceAttr(dsName, "query.0.critical", "90"),
						resource.TestCheckResourceAttr(dsName, "query.0.evaluate_backward_minutes", "3"),
					),
				),
			},
		},
	})
}

func testAccDataSourceMackerelMonitorConfigHostMetric(rand, name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "scoped" {
  name = "tf-service-%s-scoped"
}

resource "mackerel_role" "not_scoped" {
  service = mackerel_service.scoped.name
  name = "tf-role-%s-not_scoped"
}

resource "mackerel_service" "not_scoped" {
  name = "tf-service-%s-not-scoped"
}

resource "mackerel_role" "scoped" {
  service = mackerel_service.not_scoped.name
  name = "tf-role-%s-scoped"
}

resource "mackerel_monitor" "foo" {
  name = "%s"
  memo = "This monitor is managed by Terraform."
  is_mute = true
  notification_interval = 30
  host_metric {
    metric = "disk%%"
    operator = ">"
    warning = "70"
    critical = "90"
    duration = 3
    max_check_attempts = 5
    scopes = [
      mackerel_service.scoped.name,
      mackerel_role.scoped.id]
    exclude_scopes = [
      mackerel_service.not_scoped.name,
      mackerel_role.not_scoped.id]
  }
}

data "mackerel_monitor" "foo" {
  id = mackerel_monitor.foo.id
}
`, rand, rand, rand, rand, name)
}

func testAccDataSourceMackerelMonitorConfigConnectivity(rand, name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "scoped" {
  name = "tf-service-%s-scoped"
}

resource "mackerel_role" "not_scoped" {
  service = mackerel_service.scoped.name
  name = "tf-role-%s-not-scoped"
}

resource "mackerel_service" "not_scoped" {
  name = "tf-service-%s-not-scoped"
}

resource "mackerel_role" "scoped" {
  service = mackerel_service.not_scoped.name
  name = "tf-role-%s-scoped"
}

resource "mackerel_monitor" "foo" {
  name = "%s"
  memo = "This monitor is managed by Terraform."
  is_mute = true
  notification_interval = 30
  connectivity {
    scopes = [
      mackerel_service.scoped.name,
      mackerel_role.scoped.id]
    exclude_scopes = [
      mackerel_service.not_scoped.name,
      mackerel_role.not_scoped.id]
    alert_status_on_gone = "WARNING"
  }
}

data "mackerel_monitor" "foo" {
  id = mackerel_monitor.foo.id
}
`, rand, rand, rand, rand, name)
}

func testAccDataSourceMackerelMonitorConfigServiceMetric(serviceName, name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
}

resource "mackerel_monitor" "foo" {
  name = "%s"
  memo = "This monitor is managed by Terraform."
  is_mute = true
  notification_interval = 30
  service_metric {
    service = mackerel_service.foo.name
    duration = 3
    metric = "custom.access.2xx_ratio"
    operator = "<"
    warning = "99.9"
    critical = "99.99"
    max_check_attempts = 5
    missing_duration_warning = 10
    missing_duration_critical = 10080
  }
}

data "mackerel_monitor" "foo" {
  id = mackerel_monitor.foo.id
}
`, serviceName, name)
}

func testAccDataSourceMackerelMonitorConfigExternal(serviceName, name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
}

resource "mackerel_monitor" "foo" {
  name = "%s"
  memo = "This monitor is managed by Terraform."
  is_mute = true
  notification_interval = 30
  external {
    method = "POST"
    url = "https://terraform-provider-mackerel.test/"
    max_check_attempts = 3
    service = mackerel_service.foo.name
    response_time_critical = 3000
    response_time_warning = 2000
    response_time_duration = 3
    request_body = "foo=bar"
    contains_string = "blah blah blah"
    certification_expiration_critical = 7
    certification_expiration_warning = 14
    skip_certificate_verification = true
    headers = {
      Cache-Control = "no-cache"
    }
    follow_redirect = true
    expected_status_code = 200
  }
}

data "mackerel_monitor" "foo" {
  id = mackerel_monitor.foo.id
}
`, serviceName, name)
}

func testAccDataSourceMackerelMonitorExpression(name string) string {
	return fmt.Sprintf(`
resource "mackerel_monitor" "foo" {
  name = "%s"
  memo = "This monitor is managed by Terraform."
  is_mute = true
  notification_interval = 30
  expression {
    expression = "max(role(my-service:db, loadavg5))"
    operator = ">"
    warning = "0.7"
    critical = "0.9"
    evaluate_backward_minutes = 3
  }
}

data "mackerel_monitor" "foo" {
  id = mackerel_monitor.foo.id
}
`, name)
}

func testAccDataSourceMackerelMonitorConfigAnomalyDetection(rand, name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "tf-service-%s"
}

resource "mackerel_role" "foo" {
  service = mackerel_service.foo.name
  name = "tf-role-%s"
}

resource "mackerel_monitor" "foo" {
  name = "%s"
  memo = "This monitor is managed by Terraform."
  is_mute = true
  notification_interval = 30
  anomaly_detection {
    warning_sensitivity = "insensitive"
    critical_sensitivity = "normal"
    max_check_attempts = 5
    training_period_from = 1577836800
    scopes = [mackerel_role.foo.id]
  }
}

data "mackerel_monitor" "foo" {
  id = mackerel_monitor.foo.id
}
`, rand, rand, name)
}

func testAccDataSourceMackerelMonitorConfigQuery(name string) string {
	return fmt.Sprintf(`
resource "mackerel_monitor" "foo" {
  name = "%s"
  memo = "This monitor is managed by Terraform."
  is_mute = true
  notification_interval = 30
  query {
    query = "container.cpu.utilization{k8s.deployment.name=\"httpbin\"}"
    legend = "cpu.utilization {{k8s.node.name}}"
    operator = ">"
    warning = "70"
    critical = "90"
    evaluate_backward_minutes = 3
  }
}

data "mackerel_monitor" "foo" {
  id = mackerel_monitor.foo.id
}
`, name)
}
