package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelMonitor_HostMetric(t *testing.T) {
	resourceName := "mackerel_monitor.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-monitor host_metric %s", rand)
	nameUpdated := fmt.Sprintf("tf-monitor host_metric %s updated", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMackerelMonitorDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelMonitorConfigHostMetric(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "memo", ""),
					resource.TestCheckResourceAttr(resourceName, "is_mute", "false"),
					resource.TestCheckResourceAttr(resourceName, "notification_interval", "0"),
					resource.TestCheckResourceAttr(resourceName, "host_metric.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "host_metric.0.metric", "cpu.sys"),
						resource.TestCheckResourceAttr(resourceName, "host_metric.0.operator", ">"),
						resource.TestCheckResourceAttr(resourceName, "host_metric.0.warning", "75"),
						resource.TestCheckResourceAttr(resourceName, "host_metric.0.critical", ""),
						resource.TestCheckResourceAttr(resourceName, "host_metric.0.duration", "1"),
						resource.TestCheckResourceAttr(resourceName, "host_metric.0.max_check_attempts", "1"),
						resource.TestCheckResourceAttr(resourceName, "host_metric.0.scopes.#", "0"),
						resource.TestCheckResourceAttr(resourceName, "host_metric.0.exclude_scopes.#", "0"),
					),
					resource.TestCheckResourceAttr(resourceName, "connectivity.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "service_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "external.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "expression.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "anomaly_detection.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "query.#", "0"),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelMonitorConfigHostMetricUpdated(rand, nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "memo", "This monitor is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "is_mute", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification_interval", "30"),
					resource.TestCheckResourceAttr(resourceName, "host_metric.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "host_metric.0.metric", "cpu.usr"),
						resource.TestCheckResourceAttr(resourceName, "host_metric.0.operator", ">"),
						resource.TestCheckResourceAttr(resourceName, "host_metric.0.warning", "70"),
						resource.TestCheckResourceAttr(resourceName, "host_metric.0.critical", "90"),
						resource.TestCheckResourceAttr(resourceName, "host_metric.0.duration", "3"),
						resource.TestCheckResourceAttr(resourceName, "host_metric.0.max_check_attempts", "5"),
						resource.TestCheckResourceAttr(resourceName, "host_metric.0.scopes.#", "2"),
						resource.TestCheckResourceAttr(resourceName, "host_metric.0.exclude_scopes.#", "2"),
					),
					resource.TestCheckResourceAttr(resourceName, "connectivity.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "service_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "external.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "expression.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "anomaly_detection.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "query.#", "0"),
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

func TestAccMackerelMonitor_Connectivity(t *testing.T) {
	resourceName := "mackerel_monitor.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-monitor connectivity %s", rand)
	nameUpdated := fmt.Sprintf("tf-monitor connectivity %s updated", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMackerelMonitorDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelMonitorConfigConnectivity(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "memo", ""),
					resource.TestCheckResourceAttr(resourceName, "is_mute", "false"),
					resource.TestCheckResourceAttr(resourceName, "notification_interval", "0"),
					resource.TestCheckResourceAttr(resourceName, "host_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "connectivity.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "connectivity.0.scopes.#", "0"),
						resource.TestCheckResourceAttr(resourceName, "connectivity.0.exclude_scopes.#", "0"),
						resource.TestCheckResourceAttr(resourceName, "connectivity.0.alert_status_on_gone", "CRITICAL"),
					),
					resource.TestCheckResourceAttr(resourceName, "service_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "external.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "expression.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "anomaly_detection.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "query.#", "0"),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelMonitorConfigConnectivityUpdated(rand, nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "memo", "This monitor is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "is_mute", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification_interval", "30"),
					resource.TestCheckResourceAttr(resourceName, "host_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "connectivity.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "connectivity.0.scopes.#", "2"),
						resource.TestCheckResourceAttr(resourceName, "connectivity.0.exclude_scopes.#", "2"),
						resource.TestCheckResourceAttr(resourceName, "connectivity.0.alert_status_on_gone", "WARNING"),
					),
					resource.TestCheckResourceAttr(resourceName, "service_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "external.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "expression.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "anomaly_detection.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "query.#", "0"),
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

func TestAccMackerelMonitor_ServiceMetric(t *testing.T) {
	resourceName := "mackerel_monitor.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-monitor service_metric %s", rand)
	nameUpdated := fmt.Sprintf("tf-monitor service_metric %s updated", rand)
	serviceName := fmt.Sprintf("tf-service-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMackerelMonitorDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelMonitorConfigServiceMetric(serviceName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "memo", ""),
					resource.TestCheckResourceAttr(resourceName, "is_mute", "false"),
					resource.TestCheckResourceAttr(resourceName, "notification_interval", "0"),
					resource.TestCheckResourceAttr(resourceName, "host_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "connectivity.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "service_metric.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "service_metric.0.service", serviceName),
						resource.TestCheckResourceAttr(resourceName, "service_metric.0.metric", "custom.access.2xx_ratio"),
						resource.TestCheckResourceAttr(resourceName, "service_metric.0.operator", "<"),
						resource.TestCheckResourceAttr(resourceName, "service_metric.0.warning", "99.9"),
						resource.TestCheckResourceAttr(resourceName, "service_metric.0.critical", ""),
						resource.TestCheckResourceAttr(resourceName, "service_metric.0.duration", "1"),
						resource.TestCheckResourceAttr(resourceName, "service_metric.0.max_check_attempts", "1"),
						resource.TestCheckResourceAttr(resourceName, "service_metric.0.missing_duration_warning", "0"),
						resource.TestCheckResourceAttr(resourceName, "service_metric.0.missing_duration_critical", "0"),
					),
					resource.TestCheckResourceAttr(resourceName, "external.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "expression.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "anomaly_detection.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "query.#", "0"),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelMonitorConfigServiceMetricUpdated(serviceName, nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "memo", "This monitor is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "is_mute", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification_interval", "30"),
					resource.TestCheckResourceAttr(resourceName, "host_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "connectivity.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "service_metric.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "service_metric.0.service", serviceName),
						resource.TestCheckResourceAttr(resourceName, "service_metric.0.metric", "custom.access.5xx_ratio"),
						resource.TestCheckResourceAttr(resourceName, "service_metric.0.operator", "<"),
						resource.TestCheckResourceAttr(resourceName, "service_metric.0.warning", "99.9"),
						resource.TestCheckResourceAttr(resourceName, "service_metric.0.critical", "99.99"),
						resource.TestCheckResourceAttr(resourceName, "service_metric.0.duration", "3"),
						resource.TestCheckResourceAttr(resourceName, "service_metric.0.max_check_attempts", "5"),
						resource.TestCheckResourceAttr(resourceName, "service_metric.0.missing_duration_warning", "10"),
						resource.TestCheckResourceAttr(resourceName, "service_metric.0.missing_duration_critical", "10080"),
					),
					resource.TestCheckResourceAttr(resourceName, "external.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "expression.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "anomaly_detection.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "query.#", "0"),
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

func TestAccMackerelMonitor_External(t *testing.T) {
	resourceName := "mackerel_monitor.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-monitor external %s", rand)
	nameUpdated := fmt.Sprintf("tf-monitor external %s updated", rand)
	serviceName := fmt.Sprintf("tf-service-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMackerelMonitorDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelMonitorConfigExternal(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "memo", ""),
					resource.TestCheckResourceAttr(resourceName, "is_mute", "false"),
					resource.TestCheckResourceAttr(resourceName, "notification_interval", "0"),
					resource.TestCheckResourceAttr(resourceName, "host_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "connectivity.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "service_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "external.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "external.0.method", "GET"),
						resource.TestCheckResourceAttr(resourceName, "external.0.url", "https://terraform-provider-mackerel.test/"),
						resource.TestCheckResourceAttr(resourceName, "external.0.max_check_attempts", "1"),
						resource.TestCheckResourceAttr(resourceName, "external.0.service", ""),
						resource.TestCheckResourceAttr(resourceName, "external.0.response_time_critical", "0"),
						resource.TestCheckResourceAttr(resourceName, "external.0.response_time_warning", "0"),
						resource.TestCheckResourceAttr(resourceName, "external.0.response_time_duration", "0"),
						resource.TestCheckResourceAttr(resourceName, "external.0.request_body", ""),
						resource.TestCheckResourceAttr(resourceName, "external.0.contains_string", ""),
						resource.TestCheckResourceAttr(resourceName, "external.0.certification_expiration_critical", "0"),
						resource.TestCheckResourceAttr(resourceName, "external.0.certification_expiration_warning", "0"),
						resource.TestCheckResourceAttr(resourceName, "external.0.skip_certificate_verification", "false"),
						resource.TestCheckResourceAttr(resourceName, "external.0.headers.%", "0"),
						resource.TestCheckResourceAttr(resourceName, "external.0.follow_redirect", "false"),
					),
					resource.TestCheckResourceAttr(resourceName, "expression.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "anomaly_detection.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "query.#", "0"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectSensitiveValue(resourceName, tfjsonpath.New("external").AtSliceIndex(0).AtMapKey("headers")),
				},
			},
			// Test: Update
			{
				Config: testAccMackerelMonitorConfigExternalUpdated(serviceName, nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "memo", "This monitor is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "is_mute", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification_interval", "30"),
					resource.TestCheckResourceAttr(resourceName, "host_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "connectivity.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "service_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "external.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "external.0.method", "POST"),
						resource.TestCheckResourceAttr(resourceName, "external.0.url", "https://terraform-provider-mackerel.test/"),
						resource.TestCheckResourceAttr(resourceName, "external.0.max_check_attempts", "3"),
						resource.TestCheckResourceAttr(resourceName, "external.0.service", serviceName),
						resource.TestCheckResourceAttr(resourceName, "external.0.response_time_critical", "3000"),
						resource.TestCheckResourceAttr(resourceName, "external.0.response_time_warning", "2000"),
						resource.TestCheckResourceAttr(resourceName, "external.0.response_time_duration", "3"),
						resource.TestCheckResourceAttr(resourceName, "external.0.request_body", "foo=bar"),
						resource.TestCheckResourceAttr(resourceName, "external.0.contains_string", "blah blah blah"),
						resource.TestCheckResourceAttr(resourceName, "external.0.certification_expiration_critical", "7"),
						resource.TestCheckResourceAttr(resourceName, "external.0.certification_expiration_warning", "14"),
						resource.TestCheckResourceAttr(resourceName, "external.0.skip_certificate_verification", "true"),
						resource.TestCheckResourceAttr(resourceName, "external.0.headers.%", "1"),
						resource.TestCheckResourceAttr(resourceName, "external.0.headers.Cache-Control", "no-cache"),
						resource.TestCheckResourceAttr(resourceName, "external.0.follow_redirect", "true"),
						resource.TestCheckResourceAttr(resourceName, "external.0.expected_status_code", "200"),
					),
					resource.TestCheckResourceAttr(resourceName, "expression.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "anomaly_detection.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "query.#", "0"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectSensitiveValue(resourceName, tfjsonpath.New("external").AtSliceIndex(0).AtMapKey("headers")),
				},
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

func TestAccMackerelMonitor_Expression(t *testing.T) {
	resourceName := "mackerel_monitor.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-monitor expression %s", rand)
	nameUpdated := fmt.Sprintf("tf-monitor expression %s updated", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMackerelMonitorDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelMonitorConfigExpression(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "memo", ""),
					resource.TestCheckResourceAttr(resourceName, "is_mute", "false"),
					resource.TestCheckResourceAttr(resourceName, "notification_interval", "0"),
					resource.TestCheckResourceAttr(resourceName, "host_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "connectivity.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "service_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "external.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "expression.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "expression.0.expression", "max(role(my-service:db, loadavg5))"),
						resource.TestCheckResourceAttr(resourceName, "expression.0.operator", ">"),
						resource.TestCheckResourceAttr(resourceName, "expression.0.warning", "0.7"),
						resource.TestCheckResourceAttr(resourceName, "expression.0.evaluate_backward_minutes", "0"),
					),
					resource.TestCheckResourceAttr(resourceName, "anomaly_detection.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "query.#", "0"),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelMonitorConfigExpressionUpdated(nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "memo", "This monitor is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "is_mute", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification_interval", "30"),
					resource.TestCheckResourceAttr(resourceName, "host_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "connectivity.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "service_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "external.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "expression.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "expression.0.expression", "max(role(my-service:db, loadavg5))"),
						resource.TestCheckResourceAttr(resourceName, "expression.0.operator", ">"),
						resource.TestCheckResourceAttr(resourceName, "expression.0.warning", "0.7"),
						resource.TestCheckResourceAttr(resourceName, "expression.0.critical", "0.9"),
						resource.TestCheckResourceAttr(resourceName, "expression.0.evaluate_backward_minutes", "0"),
					),
					resource.TestCheckResourceAttr(resourceName, "anomaly_detection.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "query.#", "0"),
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

func TestAccMackerelMonitor_AnomalyDetection(t *testing.T) {
	resourceName := "mackerel_monitor.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-monitor anomaly_detection %s", rand)
	nameUpdated := fmt.Sprintf("tf-monitor anomaly_detection %s updated", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMackerelMonitorDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelMonitorConfigAnomalyDetection(rand, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "memo", ""),
					resource.TestCheckResourceAttr(resourceName, "is_mute", "false"),
					resource.TestCheckResourceAttr(resourceName, "notification_interval", "0"),
					resource.TestCheckResourceAttr(resourceName, "host_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "connectivity.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "service_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "external.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "expression.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "anomaly_detection.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "anomaly_detection.0.warning_sensitivity", "insensitive"),
						resource.TestCheckResourceAttr(resourceName, "anomaly_detection.0.critical_sensitivity", ""),
						resource.TestCheckResourceAttr(resourceName, "anomaly_detection.0.max_check_attempts", "3"),
						resource.TestCheckResourceAttr(resourceName, "anomaly_detection.0.training_period_from", "0"),
						resource.TestCheckResourceAttr(resourceName, "anomaly_detection.0.scopes.#", "1"),
					),
					resource.TestCheckResourceAttr(resourceName, "query.#", "0"),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelMonitorConfigAnomalyDetectionUpdated(rand, nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "memo", "This monitor is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "is_mute", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification_interval", "30"),
					resource.TestCheckResourceAttr(resourceName, "host_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "connectivity.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "service_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "external.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "expression.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "anomaly_detection.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "anomaly_detection.0.warning_sensitivity", "insensitive"),
						resource.TestCheckResourceAttr(resourceName, "anomaly_detection.0.critical_sensitivity", "normal"),
						resource.TestCheckResourceAttr(resourceName, "anomaly_detection.0.max_check_attempts", "5"),
						resource.TestCheckResourceAttr(resourceName, "anomaly_detection.0.training_period_from", "1577836800"),
						resource.TestCheckResourceAttr(resourceName, "anomaly_detection.0.scopes.#", "1"),
					),
					resource.TestCheckResourceAttr(resourceName, "query.#", "0"),
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

func TestAccMackerelMonitor_Query(t *testing.T) {
	resourceName := "mackerel_monitor.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-monitor query %s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMackerelMonitorDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelMonitorConfigQuery(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "memo", ""),
					resource.TestCheckResourceAttr(resourceName, "is_mute", "false"),
					resource.TestCheckResourceAttr(resourceName, "notification_interval", "0"),
					resource.TestCheckResourceAttr(resourceName, "host_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "connectivity.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "service_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "external.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "expression.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "anomaly_detection.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "query.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "query.0.query", "container.cpu.utilization{k8s.deployment.name=\"httpbin\"}"),
						resource.TestCheckResourceAttr(resourceName, "query.0.legend", "cpu.utilization {{k8s.node.name}}"),
						resource.TestCheckResourceAttr(resourceName, "query.0.operator", ">"),
						resource.TestCheckResourceAttr(resourceName, "query.0.warning", "70"),
						resource.TestCheckResourceAttr(resourceName, "query.0.critical", ""),
						resource.TestCheckResourceAttr(resourceName, "query.0.evaluate_backward_minutes", "2"),
					),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelMonitorConfigQueryUpdated(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "memo", "This monitor is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "is_mute", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification_interval", "30"),
					resource.TestCheckResourceAttr(resourceName, "host_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "connectivity.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "service_metric.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "external.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "expression.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "anomaly_detection.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "query.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "query.0.query", "container.cpu.utilization{k8s.deployment.name=\"nginx\"}"),
						resource.TestCheckResourceAttr(resourceName, "query.0.legend", "cpu.utilization {{k8s.node.name}}"),
						resource.TestCheckResourceAttr(resourceName, "query.0.operator", ">"),
						resource.TestCheckResourceAttr(resourceName, "query.0.warning", "70"),
						resource.TestCheckResourceAttr(resourceName, "query.0.critical", "90"),
						resource.TestCheckResourceAttr(resourceName, "query.0.evaluate_backward_minutes", "2"),
					),
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

func testAccCheckMackerelMonitorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*mackerel.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "mackerel_monitor" {
			continue
		}

		if _, err := client.GetMonitor(r.Primary.ID); err == nil {
			return fmt.Errorf("monitor still exists: %s", r.Primary.ID)
		}
	}

	return nil
}

func testAccCheckMackerelMonitorExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("monitor not found from resources: %s", n)
		}

		if r.Primary.ID == "" {
			return fmt.Errorf("no monitor ID is set")
		}

		client := testAccProvider.Meta().(*mackerel.Client)
		if _, err := client.GetMonitor(r.Primary.ID); err != nil {
			return err
		}

		return nil
	}
}

func testAccMackerelMonitorConfigHostMetric(name string) string {
	return fmt.Sprintf(`
resource "mackerel_monitor" "foo" {
  name = "%s"
  host_metric {
    metric = "cpu.sys"
    operator = ">"
    warning = 75
    duration = 1
  }
}
`, name)
}

func testAccMackerelMonitorConfigHostMetricUpdated(rand, name string) string {
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
    metric = "cpu.usr"
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
`, rand, rand, rand, rand, name)
}

func testAccMackerelMonitorConfigConnectivity(name string) string {
	return fmt.Sprintf(`
resource "mackerel_monitor" "foo" {
  name = "%s"
  connectivity {}
}
`, name)
}

func testAccMackerelMonitorConfigConnectivityUpdated(rand, name string) string {
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
`, rand, rand, rand, rand, name)
}

func testAccMackerelMonitorConfigServiceMetric(serviceName, name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
}

resource "mackerel_monitor" "foo" {
  name = "%s"
  service_metric {
    service = mackerel_service.foo.name
    duration = 1
    metric = "custom.access.2xx_ratio"
    operator = "<"
    warning = "99.9"
  }
}
`, serviceName, name)
}

func testAccMackerelMonitorConfigServiceMetricUpdated(serviceName, name string) string {
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
    metric = "custom.access.5xx_ratio"
    operator = "<"
    warning = "99.9"
    critical = "99.99"
    max_check_attempts = 5
    missing_duration_warning = 10
    missing_duration_critical = 10080
  }
}
`, serviceName, name)
}

func testAccMackerelMonitorConfigExternal(name string) string {
	return fmt.Sprintf(`
resource "mackerel_monitor" "foo" {
  name = "%s"
  external {
    method = "GET"
    url = "https://terraform-provider-mackerel.test/"
  }
}
`, name)
}

func testAccMackerelMonitorConfigExternalUpdated(serviceName, name string) string {
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
`, serviceName, name)
}

func testAccMackerelMonitorConfigExpression(name string) string {
	return fmt.Sprintf(`
resource "mackerel_monitor" "foo" {
  name = "%s"
  expression {
    expression = "max(role(my-service:db, loadavg5))"
    operator = ">"
    warning = "0.7"
  }
}
`, name)
}

func testAccMackerelMonitorConfigExpressionUpdated(name string) string {
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
  }
}
`, name)
}

func testAccMackerelMonitorConfigAnomalyDetection(rand, name string) string {
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
  anomaly_detection {
    warning_sensitivity = "insensitive"
    scopes = [mackerel_role.foo.id]
  }
}
`, rand, rand, name)
}

func testAccMackerelMonitorConfigAnomalyDetectionUpdated(rand, name string) string {
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
`, rand, rand, name)
}

func testAccMackerelMonitorConfigQuery(name string) string {
	return fmt.Sprintf(`
resource "mackerel_monitor" "foo" {
  name = "%s"
  query {
    query = "container.cpu.utilization{k8s.deployment.name=\"httpbin\"}"
    legend = "cpu.utilization {{k8s.node.name}}"
    operator = ">"
    warning = "70"
    evaluate_backward_minutes = 2
  }
}
`, name)
}

func testAccMackerelMonitorConfigQueryUpdated(name string) string {
	return fmt.Sprintf(`
resource "mackerel_monitor" "foo" {
  name = "%s"
  memo = "This monitor is managed by Terraform."
  is_mute = true
  notification_interval = 30
  query {
    query = "container.cpu.utilization{k8s.deployment.name=\"nginx\"}"
    legend = "cpu.utilization {{k8s.node.name}}"
    operator = ">"
    warning = "70"
    critical = "90"
    evaluate_backward_minutes = 2
  }
}
`, name)
}
