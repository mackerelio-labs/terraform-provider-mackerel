package mackerel

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/typeutil"
	"github.com/mackerelio/mackerel-client-go"
)

func Test_Monitor_toModel(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in mackerel.Monitor

		wants   MonitorModel
		wantErr bool
	}{
		"host/basic": {
			in: &mackerel.MonitorHostMetric{
				ID:   "5dx2AtermgS",
				Name: "tf-monitor-host",
				Type: "host",

				Metric:           "cpu_sys",
				Operator:         ">",
				Warning:          float64Pointer(75),
				Duration:         1,
				MaxCheckAttempts: 1,
			},
			wants: MonitorModel{
				ID:                   types.StringValue("5dx2AtermgS"),
				Name:                 types.StringValue("tf-monitor-host"),
				Memo:                 types.StringValue(""),
				IsMute:               types.BoolValue(false),
				NotificationInterval: types.Int64Value(0),
				HostMetricMonitor: []MonitorHostMetric{{
					MetricName:       types.StringValue("cpu_sys"),
					Operator:         types.StringValue(">"),
					Warning:          typeutil.NewFloatStringValue("75"),
					Critical:         typeutil.NewFloatStringValue(""),
					Duration:         types.Int64Value(1),
					MaxCheckAttempts: types.Int64Value(1),
				}},
			},
		},
		"host/scoped": {
			in: &mackerel.MonitorHostMetric{
				ID:                   "5dx2Av8M8um",
				Name:                 "tf-monitor-host-scoped",
				Memo:                 "This monitor is managed by Terraform.",
				Type:                 "host",
				IsMute:               true,
				NotificationInterval: 30,

				Metric:           "cpu.usr",
				Operator:         ">",
				Warning:          float64Pointer(70),
				Critical:         float64Pointer(90),
				Duration:         3,
				MaxCheckAttempts: 5,

				Scopes:        []string{"tf-svc-scope", "tf-svc-scope: tf-role-scope"},
				ExcludeScopes: []string{"tf-svc-exscope", "tf-svc-exscope: tf-role-exscope"},
			},
			wants: MonitorModel{
				ID:                   types.StringValue("5dx2Av8M8um"),
				Name:                 types.StringValue("tf-monitor-host-scoped"),
				Memo:                 types.StringValue("This monitor is managed by Terraform."),
				IsMute:               types.BoolValue(true),
				NotificationInterval: types.Int64Value(30),
				HostMetricMonitor: []MonitorHostMetric{{
					MetricName:       types.StringValue("cpu.usr"),
					Operator:         types.StringValue(">"),
					Warning:          typeutil.NewFloatStringValue("70"),
					Critical:         typeutil.NewFloatStringValue("90"),
					Duration:         types.Int64Value(3),
					MaxCheckAttempts: types.Int64Value(5),
					Scopes:           []string{"tf-svc-scope", "tf-svc-scope:tf-role-scope"},
					ExcludeScopes:    []string{"tf-svc-exscope", "tf-svc-exscope:tf-role-exscope"},
				}},
			},
		},
		"service/basic": {
			in: &mackerel.MonitorServiceMetric{
				ID:   "5dxHNTkGv3G",
				Name: "tf-monitor-svc",
				Type: "service",

				Service:          "tf-svc",
				Metric:           "custom.access.2xx_ratio",
				Operator:         "<",
				Warning:          float64Pointer(99.9),
				Duration:         1,
				MaxCheckAttempts: 1,
			},
			wants: MonitorModel{
				ID:                   types.StringValue("5dxHNTkGv3G"),
				Name:                 types.StringValue("tf-monitor-svc"),
				Memo:                 types.StringValue(""),
				IsMute:               types.BoolValue(false),
				NotificationInterval: types.Int64Value(0),
				ServiceMetricMonitor: []MonitorServiceMetric{{
					ServiceName:             types.StringValue("tf-svc"),
					MetricName:              types.StringValue("custom.access.2xx_ratio"),
					Operator:                types.StringValue("<"),
					Warning:                 typeutil.NewFloatStringValue("99.9"),
					Critical:                typeutil.NewFloatStringValue(""),
					Duration:                types.Int64Value(1),
					MaxCheckAttempts:        types.Int64Value(1),
					MissingDurationWarning:  types.Int64Value(0),
					MissingDurationCritical: types.Int64Value(0),
				}},
			},
		},
		"service/full": {
			in: &mackerel.MonitorServiceMetric{
				ID:                   "5dxHNT34wH7",
				Name:                 "tf-monitor-svc-full",
				Memo:                 "This monitor is managed by Terraform.",
				Type:                 "service",
				IsMute:               true,
				NotificationInterval: 30,

				Service:                 "tf-svc",
				Metric:                  "custom.access.5xx_ratio",
				Operator:                ">",
				Warning:                 float64Pointer(99.9),
				Critical:                float64Pointer(99.99),
				Duration:                3,
				MaxCheckAttempts:        5,
				MissingDurationWarning:  10,
				MissingDurationCritical: 10080,
			},
			wants: MonitorModel{
				ID:                   types.StringValue("5dxHNT34wH7"),
				Name:                 types.StringValue("tf-monitor-svc-full"),
				Memo:                 types.StringValue("This monitor is managed by Terraform."),
				IsMute:               types.BoolValue(true),
				NotificationInterval: types.Int64Value(30),
				ServiceMetricMonitor: []MonitorServiceMetric{{
					ServiceName:             types.StringValue("tf-svc"),
					MetricName:              types.StringValue("custom.access.5xx_ratio"),
					Operator:                types.StringValue(">"),
					Warning:                 typeutil.NewFloatStringValue("99.9"),
					Critical:                typeutil.NewFloatStringValue("99.99"),
					Duration:                types.Int64Value(3),
					MaxCheckAttempts:        types.Int64Value(5),
					MissingDurationWarning:  types.Int64Value(10),
					MissingDurationCritical: types.Int64Value(10080),
				}},
			},
		},
		"expression/basic": {
			in: &mackerel.MonitorExpression{
				ID:   "5dxWMxdx8w1",
				Name: "tf-monitor-expr",
				Type: "expression",

				Expression: "max(role(my-service:db, loadavg5))",
				Operator:   ">",
				Warning:    float64Pointer(0.7),
			},
			wants: MonitorModel{
				ID:                   types.StringValue("5dxWMxdx8w1"),
				Name:                 types.StringValue("tf-monitor-expr"),
				Memo:                 types.StringValue(""),
				IsMute:               types.BoolValue(false),
				NotificationInterval: types.Int64Value(0),
				ExpressionMonitor: []MonitorExpression{{
					Expression: types.StringValue("max(role(my-service:db, loadavg5))"),
					Operator:   types.StringValue(">"),
					Warning:    typeutil.NewFloatStringValue("0.7"),
					Critical:   typeutil.NewFloatStringValue(""),
				}},
			},
		},
		"expression/full": {
			in: &mackerel.MonitorExpression{
				ID:                   "5dxWMxdx8w1",
				Name:                 "tf-monitor-expr-full",
				Type:                 "expression",
				Memo:                 "This monitor is managed by Terraform.",
				IsMute:               true,
				NotificationInterval: 30,

				Expression: "max(role(my-service:db, loadavg5))",
				Operator:   ">",
				Warning:    float64Pointer(0.7),
				Critical:   float64Pointer(0.9),
			},
			wants: MonitorModel{
				ID:                   types.StringValue("5dxWMxdx8w1"),
				Name:                 types.StringValue("tf-monitor-expr-full"),
				Memo:                 types.StringValue("This monitor is managed by Terraform."),
				IsMute:               types.BoolValue(true),
				NotificationInterval: types.Int64Value(30),
				ExpressionMonitor: []MonitorExpression{{
					Expression: types.StringValue("max(role(my-service:db, loadavg5))"),
					Operator:   types.StringValue(">"),
					Warning:    typeutil.NewFloatStringValue("0.7"),
					Critical:   typeutil.NewFloatStringValue("0.9"),
				}},
			},
		},
		"query/basic": {
			in: &mackerel.MonitorQuery{
				ID:   "5dQpDFsYzrS",
				Name: "tf-monitor-query",
				Type: "query",

				Query:    `sum by (k8s.node.name) (container.cpu.utilization{k8s.deployment.name="nginx"})`,
				Legend:   "nginx cpu utilization on {{k8s.node.name}}",
				Operator: ">",
				Warning:  float64Pointer(0.7),
			},
			wants: MonitorModel{
				ID:                   types.StringValue("5dQpDFsYzrS"),
				Name:                 types.StringValue("tf-monitor-query"),
				Memo:                 types.StringValue(""),
				IsMute:               types.BoolValue(false),
				NotificationInterval: types.Int64Value(0),

				QueryMonitor: []MonitorQuery{{
					Query:    types.StringValue(`sum by (k8s.node.name) (container.cpu.utilization{k8s.deployment.name="nginx"})`),
					Legend:   types.StringValue("nginx cpu utilization on {{k8s.node.name}}"),
					Operator: types.StringValue(">"),
					Warning:  typeutil.NewFloatStringValue("0.7"),
					Critical: typeutil.NewFloatStringValue(""),
				}},
			},
		},
		"query/full": {
			in: &mackerel.MonitorQuery{
				ID:                   "5dQpDFKqCpJ",
				Name:                 "tf-monitor-query-full",
				Type:                 "query",
				Memo:                 "This monitor is managed by Terraform.",
				IsMute:               true,
				NotificationInterval: 30,

				Query:    `sum by (k8s.node.name) (container.cpu.utilization{k8s.deployment.name="nginx"})`,
				Legend:   "nginx cpu utilization on {{k8s.node.name}}",
				Operator: ">",
				Warning:  float64Pointer(0.7),
				Critical: float64Pointer(0.9),
			},
			wants: MonitorModel{
				ID:                   types.StringValue("5dQpDFKqCpJ"),
				Name:                 types.StringValue("tf-monitor-query-full"),
				Memo:                 types.StringValue("This monitor is managed by Terraform."),
				IsMute:               types.BoolValue(true),
				NotificationInterval: types.Int64Value(30),

				QueryMonitor: []MonitorQuery{{
					Query:    types.StringValue(`sum by (k8s.node.name) (container.cpu.utilization{k8s.deployment.name="nginx"})`),
					Legend:   types.StringValue("nginx cpu utilization on {{k8s.node.name}}"),
					Operator: types.StringValue(">"),
					Warning:  typeutil.NewFloatStringValue("0.7"),
					Critical: typeutil.NewFloatStringValue("0.9"),
				}},
			},
		},
		"connectivity/basic": {
			in: &mackerel.MonitorConnectivity{
				ID:   "5dQBLCbvquy",
				Name: "tf-monitor-connectivity",
				Type: "connectivity",

				AlertStatusOnGone: "CRITICAL",
			},
			wants: MonitorModel{
				ID:                   types.StringValue("5dQBLCbvquy"),
				Name:                 types.StringValue("tf-monitor-connectivity"),
				Memo:                 types.StringValue(""),
				IsMute:               types.BoolValue(false),
				NotificationInterval: types.Int64Value(0),

				ConnectivityMonitor: []MonitorConnectivity{{
					AlertStatusOnGone: types.StringValue("CRITICAL"),
				}},
			},
		},
		"connectivity/full": {
			in: &mackerel.MonitorConnectivity{
				ID:                   "5dQBLEkUD6s",
				Name:                 "tf-monitor-connectivity-full",
				Type:                 "connectivity",
				Memo:                 "This monitor is managed by Terraform.",
				IsMute:               true,
				NotificationInterval: 30,

				AlertStatusOnGone: "WARNING",
				Scopes:            []string{"tf-svc-scope", "tf-svc-scope: tf-role-scope"},
				ExcludeScopes:     []string{"tf-svc-exscope", "tf-svc-exscope: tf-role-exscope"},
			},
			wants: MonitorModel{
				ID:                   types.StringValue("5dQBLEkUD6s"),
				Name:                 types.StringValue("tf-monitor-connectivity-full"),
				Memo:                 types.StringValue("This monitor is managed by Terraform."),
				IsMute:               types.BoolValue(true),
				NotificationInterval: types.Int64Value(30),

				ConnectivityMonitor: []MonitorConnectivity{{
					AlertStatusOnGone: types.StringValue("WARNING"),
					Scopes:            []string{"tf-svc-scope", "tf-svc-scope:tf-role-scope"},
					ExcludeScopes:     []string{"tf-svc-exscope", "tf-svc-exscope:tf-role-exscope"},
				}},
			},
		},
		"external/basic": {
			in: &mackerel.MonitorExternalHTTP{
				ID:   "5dQKsg6fxvS",
				Name: "tf-monitor-external",
				Type: "external",

				Method:           "GET",
				URL:              "https://terraform-provider-mackerel.test",
				MaxCheckAttempts: 1,
			},
			wants: MonitorModel{
				ID:                   types.StringValue("5dQKsg6fxvS"),
				Name:                 types.StringValue("tf-monitor-external"),
				Memo:                 types.StringValue(""),
				IsMute:               types.BoolValue(false),
				NotificationInterval: types.Int64Value(0),

				ExternalMonitor: []MonitorExternal{{
					Method:                          types.StringValue("GET"),
					URL:                             types.StringValue("https://terraform-provider-mackerel.test"),
					MaxCheckAttempts:                types.Int64Value(1),
					ServiceName:                     types.StringValue(""),
					ResponseTimeCritical:            types.Float64Value(0.0),
					ResponseTimeWarning:             types.Float64Value(0.0),
					ResponseTimeDuration:            types.Int64Value(0),
					RequestBody:                     types.StringValue(""),
					ContainsString:                  types.StringValue(""),
					CertificationExpirationCritical: types.Int64Value(0),
					CertificationExpirationWarning:  types.Int64Value(0),
					SkipCertificateVerification:     types.BoolValue(false),
					Headers:                         nil,
					FollowRedirect:                  types.BoolValue(false),
				}},
			},
		},
		"external/full": {
			in: &mackerel.MonitorExternalHTTP{
				ID:                   "5dQKsiUxvf9",
				Name:                 "tf-monitor-external-full",
				Memo:                 "This monitor is managed by Terraform.",
				Type:                 "external",
				IsMute:               true,
				NotificationInterval: 30,

				Method:                          "POST",
				URL:                             "https://terraform-provider-mackerel.test/",
				MaxCheckAttempts:                3,
				Service:                         "tf-test-svc",
				ResponseTimeCritical:            float64Pointer(3000),
				ResponseTimeWarning:             float64Pointer(2000),
				ResponseTimeDuration:            uint64Pointer(3),
				RequestBody:                     "foo=bar",
				ContainsString:                  "blah blah blah",
				CertificationExpirationCritical: uint64Pointer(7),
				CertificationExpirationWarning:  uint64Pointer(14),
				SkipCertificateVerification:     true,
				FollowRedirect:                  true,
				Headers: []mackerel.HeaderField{
					{Name: "Cache-Control", Value: "no-cache"},
				},
			},
			wants: MonitorModel{
				ID:                   types.StringValue("5dQKsiUxvf9"),
				Name:                 types.StringValue("tf-monitor-external-full"),
				Memo:                 types.StringValue("This monitor is managed by Terraform."),
				IsMute:               types.BoolValue(true),
				NotificationInterval: types.Int64Value(30),

				ExternalMonitor: []MonitorExternal{{
					Method:                          types.StringValue("POST"),
					URL:                             types.StringValue("https://terraform-provider-mackerel.test/"),
					MaxCheckAttempts:                types.Int64Value(3),
					ServiceName:                     types.StringValue("tf-test-svc"),
					ResponseTimeCritical:            types.Float64Value(3000),
					ResponseTimeWarning:             types.Float64Value(2000),
					ResponseTimeDuration:            types.Int64Value(3),
					RequestBody:                     types.StringValue("foo=bar"),
					ContainsString:                  types.StringValue("blah blah blah"),
					CertificationExpirationCritical: types.Int64Value(7),
					CertificationExpirationWarning:  types.Int64Value(14),
					SkipCertificateVerification:     types.BoolValue(true),
					Headers: map[string]string{
						"Cache-Control": "no-cache",
					},
					FollowRedirect: types.BoolValue(true),
				}},
			},
		},
		"anomaly-detection/basic": {
			in: &mackerel.MonitorAnomalyDetection{
				ID:   "5ebKY6keJaW",
				Name: "tf-monitor-anomaly-detection",
				Type: "anomalyDetection",

				WarningSensitivity: "insensitive",
				MaxCheckAttempts:   3,
				Scopes:             []string{"tf-svc:tf-role"},
			},
			wants: MonitorModel{
				ID:                   types.StringValue("5ebKY6keJaW"),
				Name:                 types.StringValue("tf-monitor-anomaly-detection"),
				Memo:                 types.StringValue(""),
				IsMute:               types.BoolValue(false),
				NotificationInterval: types.Int64Value(0),
				AnomalyDetectionMonitor: []MonitorAnomalyDetection{{
					WarningSensitivity:  types.StringValue("insensitive"),
					CriticalSensitivity: types.StringValue(""),
					MaxCheckAttempts:    types.Int64Value(3),
					TrainingPeriodFrom:  types.Int64Value(0),
					Scopes:              []string{"tf-svc:tf-role"},
				}},
			},
		},
		"anomaly-detection/full": {
			in: &mackerel.MonitorAnomalyDetection{
				ID:                   "5ebKY5KmNd1",
				Name:                 "tf-monitor-anomaly-detection-full",
				Memo:                 "This monitor is managed by Terraform.",
				Type:                 "anomalyDetection",
				IsMute:               true,
				NotificationInterval: 30,

				WarningSensitivity:  "insensitive",
				CriticalSensitivity: "normal",
				MaxCheckAttempts:    5,
				TrainingPeriodFrom:  1577836800,
				Scopes:              []string{"tf-svc:tf-role"},
			},
			wants: MonitorModel{
				ID:                   types.StringValue("5ebKY5KmNd1"),
				Name:                 types.StringValue("tf-monitor-anomaly-detection-full"),
				Memo:                 types.StringValue("This monitor is managed by Terraform."),
				IsMute:               types.BoolValue(true),
				NotificationInterval: types.Int64Value(30),
				AnomalyDetectionMonitor: []MonitorAnomalyDetection{{
					WarningSensitivity:  types.StringValue("insensitive"),
					CriticalSensitivity: types.StringValue("normal"),
					MaxCheckAttempts:    types.Int64Value(5),
					TrainingPeriodFrom:  types.Int64Value(1577836800),
					Scopes:              []string{"tf-svc:tf-role"},
				}},
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m, err := newMonitor(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error: %+v", err)
			}
			if err != nil {
				return
			}

			if diff := cmp.Diff(m, tt.wants); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func Test_Monitor_toMackerelMonitor(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in    MonitorModel
		wants mackerel.Monitor
	}{
		"host/basic": {
			in: MonitorModel{
				ID:                   types.StringValue("5dx2AtermgS"),
				Name:                 types.StringValue("tf-monitor-host"),
				Memo:                 types.StringValue(""),
				IsMute:               types.BoolValue(false),
				NotificationInterval: types.Int64Value(0),
				HostMetricMonitor: []MonitorHostMetric{{
					MetricName:       types.StringValue("cpu_sys"),
					Operator:         types.StringValue(">"),
					Warning:          typeutil.NewFloatStringValue("75"),
					Critical:         typeutil.NewFloatStringValue(""),
					Duration:         types.Int64Value(1),
					MaxCheckAttempts: types.Int64Value(1),
				}},
			},
			wants: &mackerel.MonitorHostMetric{
				ID:   "5dx2AtermgS",
				Name: "tf-monitor-host",
				Type: "host",

				Metric:           "cpu_sys",
				Operator:         ">",
				Warning:          float64Pointer(75),
				Duration:         1,
				MaxCheckAttempts: 1,
			},
		},
		"host/scoped": {
			in: MonitorModel{
				ID:                   types.StringValue("5dx2Av8M8um"),
				Name:                 types.StringValue("tf-monitor-host-scoped"),
				Memo:                 types.StringValue("This monitor is managed by Terraform."),
				IsMute:               types.BoolValue(true),
				NotificationInterval: types.Int64Value(30),
				HostMetricMonitor: []MonitorHostMetric{{
					MetricName:       types.StringValue("cpu.usr"),
					Operator:         types.StringValue(">"),
					Warning:          typeutil.NewFloatStringValue("70"),
					Critical:         typeutil.NewFloatStringValue("90"),
					Duration:         types.Int64Value(3),
					MaxCheckAttempts: types.Int64Value(5),
					Scopes:           []string{"tf-svc-scope", "tf-svc-scope:tf-role-scope"},
					ExcludeScopes:    []string{"tf-svc-exscope", "tf-svc-exscope:tf-role-exscope"},
				}},
			},
			wants: &mackerel.MonitorHostMetric{
				ID:                   "5dx2Av8M8um",
				Name:                 "tf-monitor-host-scoped",
				Memo:                 "This monitor is managed by Terraform.",
				Type:                 "host",
				IsMute:               true,
				NotificationInterval: 30,

				Metric:           "cpu.usr",
				Operator:         ">",
				Warning:          float64Pointer(70),
				Critical:         float64Pointer(90),
				Duration:         3,
				MaxCheckAttempts: 5,

				Scopes:        []string{"tf-svc-scope", "tf-svc-scope:tf-role-scope"},
				ExcludeScopes: []string{"tf-svc-exscope", "tf-svc-exscope:tf-role-exscope"},
			},
		},
		"service/basic": {
			in: MonitorModel{
				ID:                   types.StringValue("5dxHNTkGv3G"),
				Name:                 types.StringValue("tf-monitor-svc"),
				Memo:                 types.StringValue(""),
				IsMute:               types.BoolValue(false),
				NotificationInterval: types.Int64Value(0),
				ServiceMetricMonitor: []MonitorServiceMetric{{
					ServiceName:             types.StringValue("tf-svc"),
					MetricName:              types.StringValue("custom.access.2xx_ratio"),
					Operator:                types.StringValue("<"),
					Warning:                 typeutil.NewFloatStringValue("99.9"),
					Critical:                typeutil.NewFloatStringValue(""),
					Duration:                types.Int64Value(1),
					MaxCheckAttempts:        types.Int64Value(1),
					MissingDurationWarning:  types.Int64Value(0),
					MissingDurationCritical: types.Int64Value(0),
				}},
			},
			wants: &mackerel.MonitorServiceMetric{
				ID:   "5dxHNTkGv3G",
				Name: "tf-monitor-svc",
				Type: "service",

				Service:          "tf-svc",
				Metric:           "custom.access.2xx_ratio",
				Operator:         "<",
				Warning:          float64Pointer(99.9),
				Duration:         1,
				MaxCheckAttempts: 1,
			},
		},
		"service/full": {
			in: MonitorModel{
				ID:                   types.StringValue("5dxHNT34wH7"),
				Name:                 types.StringValue("tf-monitor-svc-full"),
				Memo:                 types.StringValue("This monitor is managed by Terraform."),
				IsMute:               types.BoolValue(true),
				NotificationInterval: types.Int64Value(30),
				ServiceMetricMonitor: []MonitorServiceMetric{{
					ServiceName:             types.StringValue("tf-svc"),
					MetricName:              types.StringValue("custom.access.5xx_ratio"),
					Operator:                types.StringValue(">"),
					Warning:                 typeutil.NewFloatStringValue("99.9"),
					Critical:                typeutil.NewFloatStringValue("99.99"),
					Duration:                types.Int64Value(3),
					MaxCheckAttempts:        types.Int64Value(5),
					MissingDurationWarning:  types.Int64Value(10),
					MissingDurationCritical: types.Int64Value(10080),
				}},
			},
			wants: &mackerel.MonitorServiceMetric{
				ID:                   "5dxHNT34wH7",
				Name:                 "tf-monitor-svc-full",
				Memo:                 "This monitor is managed by Terraform.",
				Type:                 "service",
				IsMute:               true,
				NotificationInterval: 30,

				Service:                 "tf-svc",
				Metric:                  "custom.access.5xx_ratio",
				Operator:                ">",
				Warning:                 float64Pointer(99.9),
				Critical:                float64Pointer(99.99),
				Duration:                3,
				MaxCheckAttempts:        5,
				MissingDurationWarning:  10,
				MissingDurationCritical: 10080,
			},
		},
		"expression/basic": {
			in: MonitorModel{
				ID:                   types.StringValue("5dxWMxdx8w1"),
				Name:                 types.StringValue("tf-monitor-expr"),
				Memo:                 types.StringValue(""),
				IsMute:               types.BoolValue(false),
				NotificationInterval: types.Int64Value(0),
				ExpressionMonitor: []MonitorExpression{{
					Expression: types.StringValue("max(role(my-service:db, loadavg5))"),
					Operator:   types.StringValue(">"),
					Warning:    typeutil.NewFloatStringValue("0.7"),
					Critical:   typeutil.NewFloatStringValue(""),
				}},
			},
			wants: &mackerel.MonitorExpression{
				ID:   "5dxWMxdx8w1",
				Name: "tf-monitor-expr",
				Type: "expression",

				Expression: "max(role(my-service:db, loadavg5))",
				Operator:   ">",
				Warning:    float64Pointer(0.7),
			},
		},
		"expression/full": {
			in: MonitorModel{
				ID:                   types.StringValue("5dxWMxdx8w1"),
				Name:                 types.StringValue("tf-monitor-expr-full"),
				Memo:                 types.StringValue("This monitor is managed by Terraform."),
				IsMute:               types.BoolValue(true),
				NotificationInterval: types.Int64Value(30),
				ExpressionMonitor: []MonitorExpression{{
					Expression: types.StringValue("max(role(my-service:db, loadavg5))"),
					Operator:   types.StringValue(">"),
					Warning:    typeutil.NewFloatStringValue("0.7"),
					Critical:   typeutil.NewFloatStringValue("0.9"),
				}},
			},
			wants: &mackerel.MonitorExpression{
				ID:                   "5dxWMxdx8w1",
				Name:                 "tf-monitor-expr-full",
				Type:                 "expression",
				Memo:                 "This monitor is managed by Terraform.",
				IsMute:               true,
				NotificationInterval: 30,

				Expression: "max(role(my-service:db, loadavg5))",
				Operator:   ">",
				Warning:    float64Pointer(0.7),
				Critical:   float64Pointer(0.9),
			},
		},
		"query/basic": {
			in: MonitorModel{
				ID:                   types.StringValue("5dQpDFsYzrS"),
				Name:                 types.StringValue("tf-monitor-query"),
				Memo:                 types.StringValue(""),
				IsMute:               types.BoolValue(false),
				NotificationInterval: types.Int64Value(0),

				QueryMonitor: []MonitorQuery{{
					Query:    types.StringValue(`sum by (k8s.node.name) (container.cpu.utilization{k8s.deployment.name="nginx"})`),
					Legend:   types.StringValue("nginx cpu utilization on {{k8s.node.name}}"),
					Operator: types.StringValue(">"),
					Warning:  typeutil.NewFloatStringValue("0.7"),
					Critical: typeutil.NewFloatStringValue(""),
				}},
			},
			wants: &mackerel.MonitorQuery{
				ID:   "5dQpDFsYzrS",
				Name: "tf-monitor-query",
				Type: "query",

				Query:    `sum by (k8s.node.name) (container.cpu.utilization{k8s.deployment.name="nginx"})`,
				Legend:   "nginx cpu utilization on {{k8s.node.name}}",
				Operator: ">",
				Warning:  float64Pointer(0.7),
			},
		},
		"query/full": {
			in: MonitorModel{
				ID:                   types.StringValue("5dQpDFKqCpJ"),
				Name:                 types.StringValue("tf-monitor-query-full"),
				Memo:                 types.StringValue("This monitor is managed by Terraform."),
				IsMute:               types.BoolValue(true),
				NotificationInterval: types.Int64Value(30),

				QueryMonitor: []MonitorQuery{{
					Query:    types.StringValue(`sum by (k8s.node.name) (container.cpu.utilization{k8s.deployment.name="nginx"})`),
					Legend:   types.StringValue("nginx cpu utilization on {{k8s.node.name}}"),
					Operator: types.StringValue(">"),
					Warning:  typeutil.NewFloatStringValue("0.7"),
					Critical: typeutil.NewFloatStringValue("0.9"),
				}},
			},
			wants: &mackerel.MonitorQuery{
				ID:                   "5dQpDFKqCpJ",
				Name:                 "tf-monitor-query-full",
				Type:                 "query",
				Memo:                 "This monitor is managed by Terraform.",
				IsMute:               true,
				NotificationInterval: 30,

				Query:    `sum by (k8s.node.name) (container.cpu.utilization{k8s.deployment.name="nginx"})`,
				Legend:   "nginx cpu utilization on {{k8s.node.name}}",
				Operator: ">",
				Warning:  float64Pointer(0.7),
				Critical: float64Pointer(0.9),
			},
		},
		"connectivity/basic": {
			in: MonitorModel{
				ID:                   types.StringValue("5dQBLCbvquy"),
				Name:                 types.StringValue("tf-monitor-connectivity"),
				Memo:                 types.StringValue(""),
				IsMute:               types.BoolValue(false),
				NotificationInterval: types.Int64Value(0),

				ConnectivityMonitor: []MonitorConnectivity{{
					AlertStatusOnGone: types.StringValue("CRITICAL"),
				}},
			},
			wants: &mackerel.MonitorConnectivity{
				ID:   "5dQBLCbvquy",
				Name: "tf-monitor-connectivity",
				Type: "connectivity",

				AlertStatusOnGone: "CRITICAL",
			},
		},
		"connectivity/full": {
			in: MonitorModel{
				ID:                   types.StringValue("5dQBLEkUD6s"),
				Name:                 types.StringValue("tf-monitor-connectivity-full"),
				Memo:                 types.StringValue("This monitor is managed by Terraform."),
				IsMute:               types.BoolValue(true),
				NotificationInterval: types.Int64Value(30),

				ConnectivityMonitor: []MonitorConnectivity{{
					AlertStatusOnGone: types.StringValue("WARNING"),
					Scopes:            []string{"tf-svc-scope", "tf-svc-scope:tf-role-scope"},
					ExcludeScopes:     []string{"tf-svc-exscope", "tf-svc-exscope:tf-role-exscope"},
				}},
			},
			wants: &mackerel.MonitorConnectivity{
				ID:                   "5dQBLEkUD6s",
				Name:                 "tf-monitor-connectivity-full",
				Type:                 "connectivity",
				Memo:                 "This monitor is managed by Terraform.",
				IsMute:               true,
				NotificationInterval: 30,

				AlertStatusOnGone: "WARNING",
				Scopes:            []string{"tf-svc-scope", "tf-svc-scope:tf-role-scope"},
				ExcludeScopes:     []string{"tf-svc-exscope", "tf-svc-exscope:tf-role-exscope"},
			},
		},
		"external/basic": {
			in: MonitorModel{
				ID:                   types.StringValue("5dQKsg6fxvS"),
				Name:                 types.StringValue("tf-monitor-external"),
				Memo:                 types.StringValue(""),
				IsMute:               types.BoolValue(false),
				NotificationInterval: types.Int64Value(0),

				ExternalMonitor: []MonitorExternal{{
					Method:                          types.StringValue("GET"),
					URL:                             types.StringValue("https://terraform-provider-mackerel.test"),
					MaxCheckAttempts:                types.Int64Value(1),
					ServiceName:                     types.StringValue(""),
					ResponseTimeCritical:            types.Float64Value(0.0),
					ResponseTimeWarning:             types.Float64Value(0.0),
					ResponseTimeDuration:            types.Int64Value(0),
					RequestBody:                     types.StringValue(""),
					ContainsString:                  types.StringValue(""),
					CertificationExpirationCritical: types.Int64Value(0),
					CertificationExpirationWarning:  types.Int64Value(0),
					SkipCertificateVerification:     types.BoolValue(false),
					Headers:                         nil,
					FollowRedirect:                  types.BoolValue(false),
				}},
			},
			wants: &mackerel.MonitorExternalHTTP{
				ID:   "5dQKsg6fxvS",
				Name: "tf-monitor-external",
				Type: "external",

				Method:           "GET",
				URL:              "https://terraform-provider-mackerel.test",
				MaxCheckAttempts: 1,
				Headers:          []mackerel.HeaderField{},
			},
		},
		"external/full": {
			in: MonitorModel{
				ID:                   types.StringValue("5dQKsiUxvf9"),
				Name:                 types.StringValue("tf-monitor-external-full"),
				Memo:                 types.StringValue("This monitor is managed by Terraform."),
				IsMute:               types.BoolValue(true),
				NotificationInterval: types.Int64Value(30),

				ExternalMonitor: []MonitorExternal{{
					Method:                          types.StringValue("POST"),
					URL:                             types.StringValue("https://terraform-provider-mackerel.test/"),
					MaxCheckAttempts:                types.Int64Value(3),
					ServiceName:                     types.StringValue("tf-test-svc"),
					ResponseTimeCritical:            types.Float64Value(3000),
					ResponseTimeWarning:             types.Float64Value(2000),
					ResponseTimeDuration:            types.Int64Value(3),
					RequestBody:                     types.StringValue("foo=bar"),
					ContainsString:                  types.StringValue("blah blah blah"),
					CertificationExpirationCritical: types.Int64Value(7),
					CertificationExpirationWarning:  types.Int64Value(14),
					SkipCertificateVerification:     types.BoolValue(true),
					Headers: map[string]string{
						"Cache-Control": "no-cache",
					},
					FollowRedirect: types.BoolValue(true),
				}},
			},
			wants: &mackerel.MonitorExternalHTTP{
				ID:                   "5dQKsiUxvf9",
				Name:                 "tf-monitor-external-full",
				Memo:                 "This monitor is managed by Terraform.",
				Type:                 "external",
				IsMute:               true,
				NotificationInterval: 30,

				Method:                          "POST",
				URL:                             "https://terraform-provider-mackerel.test/",
				MaxCheckAttempts:                3,
				Service:                         "tf-test-svc",
				ResponseTimeCritical:            float64Pointer(3000),
				ResponseTimeWarning:             float64Pointer(2000),
				ResponseTimeDuration:            uint64Pointer(3),
				RequestBody:                     "foo=bar",
				ContainsString:                  "blah blah blah",
				CertificationExpirationCritical: uint64Pointer(7),
				CertificationExpirationWarning:  uint64Pointer(14),
				SkipCertificateVerification:     true,
				FollowRedirect:                  true,
				Headers: []mackerel.HeaderField{
					{Name: "Cache-Control", Value: "no-cache"},
				},
			},
		},
		"anomaly-detection/basic": {
			in: MonitorModel{
				ID:                   types.StringValue("5ebKY6keJaW"),
				Name:                 types.StringValue("tf-monitor-anomaly-detection"),
				Memo:                 types.StringValue(""),
				IsMute:               types.BoolValue(false),
				NotificationInterval: types.Int64Value(0),
				AnomalyDetectionMonitor: []MonitorAnomalyDetection{{
					WarningSensitivity:  types.StringValue("insensitive"),
					CriticalSensitivity: types.StringValue(""),
					MaxCheckAttempts:    types.Int64Value(3),
					TrainingPeriodFrom:  types.Int64Value(0),
					Scopes:              []string{"tf-svc:tf-role"},
				}},
			},
			wants: &mackerel.MonitorAnomalyDetection{
				ID:   "5ebKY6keJaW",
				Name: "tf-monitor-anomaly-detection",
				Type: "anomalyDetection",

				WarningSensitivity: "insensitive",
				MaxCheckAttempts:   3,
				Scopes:             []string{"tf-svc:tf-role"},
			},
		},
		"anomaly-detection/full": {
			in: MonitorModel{
				ID:                   types.StringValue("5ebKY5KmNd1"),
				Name:                 types.StringValue("tf-monitor-anomaly-detection-full"),
				Memo:                 types.StringValue("This monitor is managed by Terraform."),
				IsMute:               types.BoolValue(true),
				NotificationInterval: types.Int64Value(30),
				AnomalyDetectionMonitor: []MonitorAnomalyDetection{{
					WarningSensitivity:  types.StringValue("insensitive"),
					CriticalSensitivity: types.StringValue("normal"),
					MaxCheckAttempts:    types.Int64Value(5),
					TrainingPeriodFrom:  types.Int64Value(1577836800),
					Scopes:              []string{"tf-svc:tf-role"},
				}},
			},
			wants: &mackerel.MonitorAnomalyDetection{
				ID:                   "5ebKY5KmNd1",
				Name:                 "tf-monitor-anomaly-detection-full",
				Memo:                 "This monitor is managed by Terraform.",
				Type:                 "anomalyDetection",
				IsMute:               true,
				NotificationInterval: 30,

				WarningSensitivity:  "insensitive",
				CriticalSensitivity: "normal",
				MaxCheckAttempts:    5,
				TrainingPeriodFrom:  1577836800,
				Scopes:              []string{"tf-svc:tf-role"},
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if diff := cmp.Diff(tt.in.mackerelMonitor(), tt.wants); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func float64Pointer(f float64) *float64 {
	return &f
}

func uint64Pointer(u64 uint64) *uint64 {
	return &u64
}
