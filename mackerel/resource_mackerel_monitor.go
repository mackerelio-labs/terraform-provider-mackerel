package mackerel

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelMonitorCreate,
		Read:   resourceMackerelMonitorRead,
		Update: resourceMackerelMonitorUpdate,
		Delete: resourceMackerelMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"memo": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_mute": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"notification_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"host_metric": {
				Type:         schema.TypeList,
				Optional:     true,
				ExactlyOneOf: []string{"host_metric", "connectivity", "service_metric", "external", "expression", "anomaly_detection"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric": {
							Type:     schema.TypeString,
							Required: true,
						},
						"operator": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{">", "<"}, false),
						},
						"warning": {
							Type:         schema.TypeFloat,
							Optional:     true,
							AtLeastOneOf: []string{"host_metric.0.warning", "host_metric.0.critical"},
						},
						"critical": {
							Type:         schema.TypeFloat,
							Optional:     true,
							AtLeastOneOf: []string{"host_metric.0.warning", "host_metric.0.critical"},
						},
						"duration": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 10),
						},
						"max_check_attempts": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(1, 10),
						},
						"scopes": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"exclude_scopes": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"connectivity": {
				Type:         schema.TypeList,
				Optional:     true,
				ExactlyOneOf: []string{"host_metric", "connectivity", "service_metric", "external", "expression", "anomaly_detection"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scopes": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"exclude_scopes": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"service_metric": {
				Type:         schema.TypeList,
				Optional:     true,
				ExactlyOneOf: []string{"host_metric", "connectivity", "service_metric", "external", "expression", "anomaly_detection"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service": {
							Type:     schema.TypeString,
							Required: true,
						},
						"metric": {
							Type:     schema.TypeString,
							Required: true,
						},
						"operator": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{">", "<"}, false),
						},
						"warning": {
							Type:         schema.TypeFloat,
							Optional:     true,
							AtLeastOneOf: []string{"service_metric.0.warning", "service_metric.0.critical"},
						},
						"critical": {
							Type:         schema.TypeFloat,
							Optional:     true,
							AtLeastOneOf: []string{"service_metric.0.warning", "service_metric.0.critical"},
						},
						"duration": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"max_check_attempts": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(1, 10),
						},
						"missing_duration_warning": {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateFunc: validation.All(
								validation.IntBetween(10, 7*24*60),
								validation.IntDivisibleBy(10),
							),
						},
						"missing_duration_critical": {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateFunc: validation.All(
								validation.IntBetween(10, 7*24*60),
								validation.IntDivisibleBy(10),
							),
						},
					},
				},
			},
			"external": {
				Type:         schema.TypeList,
				Optional:     true,
				ExactlyOneOf: []string{"host_metric", "connectivity", "service_metric", "external", "expression", "anomaly_detection"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"method": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"GET", "POST", "PUT", "DELETE"}, false),
						},
						"url": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsURLWithHTTPorHTTPS,
						},
						"max_check_attempts": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(1, 10),
						},
						"service": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"response_time_critical": {
							Type:         schema.TypeFloat,
							Optional:     true,
							RequiredWith: []string{"external.0.service"},
						},
						"response_time_warning": {
							Type:         schema.TypeFloat,
							Optional:     true,
							RequiredWith: []string{"external.0.service"},
						},
						"response_time_duration": {
							Type:     schema.TypeInt,
							Optional: true,
							// Default:      1,
							RequiredWith: []string{"external.0.service"},
							ValidateFunc: validation.IntBetween(1, 10),
						},
						"request_body": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"contains_string": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"certification_expiration_critical": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"certification_expiration_warning": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"skip_certificate_verification": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"headers": {
							Type:      schema.TypeMap,
							Optional:  true,
							Sensitive: true,
							Elem:      &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"expression": {
				Type:         schema.TypeList,
				Optional:     true,
				ExactlyOneOf: []string{"host_metric", "connectivity", "service_metric", "external", "expression", "anomaly_detection"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expression": {
							Type:     schema.TypeString,
							Required: true,
						},
						"operator": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{">", "<"}, false),
						},
						"warning": {
							Type:         schema.TypeFloat,
							Optional:     true,
							AtLeastOneOf: []string{"expression.0.warning", "expression.0.critical"},
						},
						"critical": {
							Type:         schema.TypeFloat,
							Optional:     true,
							AtLeastOneOf: []string{"expression.0.warning", "expression.0.critical"},
						},
					},
				},
			},
			"anomaly_detection": {
				Type:         schema.TypeList,
				Optional:     true,
				ExactlyOneOf: []string{"host_metric", "connectivity", "service_metric", "external", "expression", "anomaly_detection"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"warning_sensitivity": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: []string{"anomaly_detection.0.warning_sensitivity", "anomaly_detection.0.critical_sensitivity"},
							ValidateFunc: validation.StringInSlice([]string{"insensitive", "normal", "sensitive"}, false),
						},
						"critical_sensitivity": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: []string{"anomaly_detection.0.warning_sensitivity", "anomaly_detection.0.critical_sensitivity"},
							ValidateFunc: validation.StringInSlice([]string{"insensitive", "normal", "sensitive"}, false),
						},
						"max_check_attempts": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      3,
							ValidateFunc: validation.IntBetween(1, 10),
						},
						"training_period_from": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"scopes": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceMackerelMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	monitor, err := client.CreateMonitor(expandMonitor(d))
	if err != nil {
		return err
	}
	d.SetId(monitor.MonitorID())
	return resourceMackerelMonitorRead(d, meta)
}

func resourceMackerelMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	monitor, err := client.GetMonitor(d.Id())
	if err != nil {
		return err
	}
	return flattenMonitor(monitor, d)
}

func resourceMackerelMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	monitor, err := client.UpdateMonitor(d.Id(), expandMonitor(d))
	if err != nil {
		return err
	}
	d.SetId(monitor.MonitorID())
	return resourceMackerelMonitorRead(d, meta)
}

func resourceMackerelMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	_, err := client.DeleteMonitor(d.Id())
	return err
}

func expandMonitor(d *schema.ResourceData) mackerel.Monitor {
	var monitor mackerel.Monitor
	if _, ok := d.GetOk("host_metric"); ok {
		monitor = expandMonitorHostMetric(d)
	}
	if _, ok := d.GetOk("connectivity"); ok {
		monitor = expandMonitorConnectivity(d)
	}
	if _, ok := d.GetOk("service_metric"); ok {
		monitor = expandMonitorServiceMetric(d)
	}
	if _, ok := d.GetOk("external"); ok {
		monitor = expandMonitorExternalHTTP(d)
	}
	if _, ok := d.GetOk("expression"); ok {
		monitor = expandMonitorExpression(d)
	}
	if _, ok := d.GetOk("anomaly_detection"); ok {
		monitor = expandMonitorAnomalyDetection(d)
	}
	return monitor
}

func flattenMonitor(monitor mackerel.Monitor, d *schema.ResourceData) error {
	var err error
	if v, ok := monitor.(*mackerel.MonitorHostMetric); ok {
		err = flattenMonitorHostMetric(v, d)
	}
	if v, ok := monitor.(*mackerel.MonitorConnectivity); ok {
		err = flattenMonitorConnectivity(v, d)
	}
	if v, ok := monitor.(*mackerel.MonitorServiceMetric); ok {
		err = flattenMonitorServiceMetric(v, d)
	}
	if v, ok := monitor.(*mackerel.MonitorExternalHTTP); ok {
		err = flattenMonitorExternalHTTP(v, d)
	}
	if v, ok := monitor.(*mackerel.MonitorExpression); ok {
		err = flattenMonitorExpression(v, d)
	}
	if v, ok := monitor.(*mackerel.MonitorAnomalyDetection); ok {
		err = flattenMonitorAnomalyDetection(v, d)
	}
	return err
}

func expandMonitorHostMetric(d *schema.ResourceData) *mackerel.MonitorHostMetric {
	monitor := &mackerel.MonitorHostMetric{
		Name:                 d.Get("name").(string),
		Memo:                 d.Get("memo").(string),
		Type:                 "host",
		IsMute:               d.Get("is_mute").(bool),
		NotificationInterval: uint64(d.Get("notification_interval").(int)),
		Metric:               d.Get("host_metric.0.metric").(string),
		Operator:             d.Get("host_metric.0.operator").(string),
		Warning:              nil,
		Critical:             nil,
		Duration:             uint64(d.Get("host_metric.0.duration").(int)),
		MaxCheckAttempts:     uint64(d.Get("host_metric.0.max_check_attempts").(int)),
		Scopes:               expandStringListFromSet(d.Get("host_metric.0.scopes").(*schema.Set)),
		ExcludeScopes:        expandStringListFromSet(d.Get("host_metric.0.exclude_scopes").(*schema.Set)),
	}
	if warning, ok := d.GetOk("host_metric.0.warning"); ok {
		warning := warning.(float64)
		monitor.Warning = &warning
	}
	if critical, ok := d.GetOk("host_metric.0.critical"); ok {
		critical := critical.(float64)
		monitor.Critical = &critical
	}
	return monitor
}

func flattenMonitorHostMetric(monitor *mackerel.MonitorHostMetric, d *schema.ResourceData) error {
	d.Set("name", monitor.Name)
	d.Set("memo", monitor.Memo)
	d.Set("is_mute", monitor.IsMute)
	d.Set("notification_interval", monitor.NotificationInterval)
	normalizedScopes := make([]string, 0, len(monitor.Scopes))
	for _, s := range monitor.Scopes {
		normalizedScopes = append(normalizedScopes, strings.ReplaceAll(s, " ", ""))
	}
	normalizedExcludeScopes := make([]string, 0, len(monitor.ExcludeScopes))
	for _, s := range monitor.ExcludeScopes {
		normalizedExcludeScopes = append(normalizedExcludeScopes, strings.ReplaceAll(s, " ", ""))
	}
	d.Set("host_metric", []map[string]interface{}{
		{
			"metric":             monitor.Metric,
			"operator":           monitor.Operator,
			"warning":            monitor.Warning,
			"critical":           monitor.Critical,
			"duration":           monitor.Duration,
			"max_check_attempts": monitor.MaxCheckAttempts,
			"scopes":             flattenStringListToSet(normalizedScopes),
			"exclude_scopes":     flattenStringListToSet(normalizedExcludeScopes),
		},
	})
	return nil
}

func expandMonitorConnectivity(d *schema.ResourceData) *mackerel.MonitorConnectivity {
	monitor := &mackerel.MonitorConnectivity{
		Name:                 d.Get("name").(string),
		Memo:                 d.Get("memo").(string),
		Type:                 "connectivity",
		IsMute:               d.Get("is_mute").(bool),
		NotificationInterval: uint64(d.Get("notification_interval").(int)),
		Scopes:               expandStringListFromSet(d.Get("connectivity.0.scopes").(*schema.Set)),
		ExcludeScopes:        expandStringListFromSet(d.Get("connectivity.0.exclude_scopes").(*schema.Set)),
	}
	return monitor
}

func flattenMonitorConnectivity(monitor *mackerel.MonitorConnectivity, d *schema.ResourceData) error {
	d.Set("name", monitor.Name)
	d.Set("memo", monitor.Memo)
	d.Set("is_mute", monitor.IsMute)
	d.Set("notification_interval", monitor.NotificationInterval)
	normalizedScopes := make([]string, 0, len(monitor.Scopes))
	for _, s := range monitor.Scopes {
		normalizedScopes = append(normalizedScopes, strings.ReplaceAll(s, " ", ""))
	}
	normalizedExcludeScopes := make([]string, 0, len(monitor.ExcludeScopes))
	for _, s := range monitor.ExcludeScopes {
		normalizedExcludeScopes = append(normalizedExcludeScopes, strings.ReplaceAll(s, " ", ""))
	}
	d.Set("connectivity", []map[string]interface{}{
		{
			"scopes":         flattenStringListToSet(normalizedScopes),
			"exclude_scopes": flattenStringListToSet(normalizedExcludeScopes),
		},
	})
	return nil
}

func expandMonitorServiceMetric(d *schema.ResourceData) *mackerel.MonitorServiceMetric {
	monitor := &mackerel.MonitorServiceMetric{
		Name:                    d.Get("name").(string),
		Memo:                    d.Get("memo").(string),
		Type:                    "service",
		IsMute:                  d.Get("is_mute").(bool),
		NotificationInterval:    uint64(d.Get("notification_interval").(int)),
		Service:                 d.Get("service_metric.0.service").(string),
		Metric:                  d.Get("service_metric.0.metric").(string),
		Operator:                d.Get("service_metric.0.operator").(string),
		Warning:                 nil,
		Critical:                nil,
		Duration:                uint64(d.Get("service_metric.0.duration").(int)),
		MaxCheckAttempts:        uint64(d.Get("service_metric.0.max_check_attempts").(int)),
		MissingDurationWarning:  uint64(d.Get("service_metric.0.missing_duration_warning").(int)),
		MissingDurationCritical: uint64(d.Get("service_metric.0.missing_duration_critical").(int)),
	}
	if warning, ok := d.GetOk("service_metric.0.warning"); ok {
		warning := warning.(float64)
		monitor.Warning = &warning
	}
	if critical, ok := d.GetOk("service_metric.0.critical"); ok {
		critical := critical.(float64)
		monitor.Critical = &critical
	}
	return monitor
}

func flattenMonitorServiceMetric(monitor *mackerel.MonitorServiceMetric, d *schema.ResourceData) error {
	d.Set("name", monitor.Name)
	d.Set("memo", monitor.Memo)
	d.Set("is_mute", monitor.IsMute)
	d.Set("notification_interval", monitor.NotificationInterval)
	d.Set("service_metric", []map[string]interface{}{
		{
			"service":                   monitor.Service,
			"metric":                    monitor.Metric,
			"operator":                  monitor.Operator,
			"warning":                   monitor.Warning,
			"critical":                  monitor.Critical,
			"duration":                  monitor.Duration,
			"max_check_attempts":        monitor.MaxCheckAttempts,
			"missing_duration_warning":  monitor.MissingDurationWarning,
			"missing_duration_critical": monitor.MissingDurationCritical,
		},
	})
	return nil
}

func expandMonitorExternalHTTP(d *schema.ResourceData) *mackerel.MonitorExternalHTTP {
	monitor := &mackerel.MonitorExternalHTTP{
		Name:                            d.Get("name").(string),
		Memo:                            d.Get("memo").(string),
		Type:                            "external",
		IsMute:                          d.Get("is_mute").(bool),
		NotificationInterval:            uint64(d.Get("notification_interval").(int)),
		Method:                          d.Get("external.0.method").(string),
		URL:                             d.Get("external.0.url").(string),
		MaxCheckAttempts:                uint64(d.Get("external.0.max_check_attempts").(int)),
		Service:                         d.Get("external.0.service").(string),
		ResponseTimeCritical:            nil,
		ResponseTimeWarning:             nil,
		ResponseTimeDuration:            nil,
		RequestBody:                     d.Get("external.0.request_body").(string),
		ContainsString:                  d.Get("external.0.contains_string").(string),
		CertificationExpirationCritical: nil,
		CertificationExpirationWarning:  nil,
		SkipCertificateVerification:     d.Get("external.0.skip_certificate_verification").(bool),
		Headers:                         []mackerel.HeaderField{},
	}
	if responseTimeCritical, ok := d.GetOk("external.0.response_time_critical"); ok {
		responseTimeCritical := responseTimeCritical.(float64)
		monitor.ResponseTimeCritical = &responseTimeCritical
	}
	if responseTimeWarning, ok := d.GetOk("external.0.response_time_warning"); ok {
		responseTimeWarning := responseTimeWarning.(float64)
		monitor.ResponseTimeWarning = &responseTimeWarning
	}
	if responseTimeDuration, ok := d.GetOk("external.0.response_time_duration"); ok {
		responseTimeDuration := uint64(responseTimeDuration.(int))
		monitor.ResponseTimeDuration = &responseTimeDuration
	}
	if certificationExpirationCritical, ok := d.GetOk("external.0.certification_expiration_critical"); ok {
		certificationExpirationCritical := uint64(certificationExpirationCritical.(int))
		monitor.CertificationExpirationCritical = &certificationExpirationCritical
	}
	if certificationExpirationWarning, ok := d.GetOk("external.0.certification_expiration_warning"); ok {
		certificationExpirationWarning := uint64(certificationExpirationWarning.(int))
		monitor.CertificationExpirationWarning = &certificationExpirationWarning
	}
	if headers, ok := d.GetOk("external.0.headers"); ok {
		for name, value := range headers.(map[string]interface{}) {
			monitor.Headers = append(monitor.Headers, mackerel.HeaderField{Name: name, Value: value.(string)})
		}
	}
	return monitor
}

func flattenMonitorExternalHTTP(monitor *mackerel.MonitorExternalHTTP, d *schema.ResourceData) error {
	d.Set("name", monitor.Name)
	d.Set("memo", monitor.Memo)
	d.Set("is_mute", monitor.IsMute)
	d.Set("notification_interval", monitor.NotificationInterval)
	external := map[string]interface{}{
		"method":                            monitor.Method,
		"url":                               monitor.URL,
		"max_check_attempts":                monitor.MaxCheckAttempts,
		"service":                           monitor.Service,
		"response_time_critical":            monitor.ResponseTimeCritical,
		"response_time_warning":             monitor.ResponseTimeWarning,
		"response_time_duration":            monitor.ResponseTimeDuration,
		"request_body":                      monitor.RequestBody,
		"contains_string":                   monitor.ContainsString,
		"certification_expiration_critical": monitor.CertificationExpirationCritical,
		"certification_expiration_warning":  monitor.CertificationExpirationWarning,
		"skip_certificate_verification":     monitor.SkipCertificateVerification,
		"headers":                           flattenHeaderField(monitor.Headers),
	}
	d.Set("external", []map[string]interface{}{external})
	return nil
}

func flattenHeaderField(fields []mackerel.HeaderField) map[string]interface{} {
	headers := make(map[string]interface{}, len(fields))
	for _, f := range fields {
		headers[f.Name] = f.Value
	}
	return headers
}

func expandMonitorExpression(d *schema.ResourceData) *mackerel.MonitorExpression {
	monitor := &mackerel.MonitorExpression{
		Name:                 d.Get("name").(string),
		Memo:                 d.Get("memo").(string),
		Type:                 "expression",
		IsMute:               d.Get("is_mute").(bool),
		NotificationInterval: uint64(d.Get("notification_interval").(int)),
		Expression:           d.Get("expression.0.expression").(string),
		Operator:             d.Get("expression.0.operator").(string),
		Warning:              nil,
		Critical:             nil,
	}
	if warning, ok := d.GetOk("expression.0.warning"); ok {
		warning := warning.(float64)
		monitor.Warning = &warning
	}
	if critical, ok := d.GetOk("expression.0.critical"); ok {
		critical := critical.(float64)
		monitor.Critical = &critical
	}
	return monitor
}

func flattenMonitorExpression(monitor *mackerel.MonitorExpression, d *schema.ResourceData) error {
	d.Set("name", monitor.Name)
	d.Set("memo", monitor.Memo)
	d.Set("is_mute", monitor.IsMute)
	d.Set("notification_interval", monitor.NotificationInterval)
	d.Set("expression", []map[string]interface{}{
		{
			"expression": monitor.Expression,
			"operator":   monitor.Operator,
			"warning":    monitor.Warning,
			"critical":   monitor.Critical,
		},
	})
	return nil
}

func expandMonitorAnomalyDetection(d *schema.ResourceData) *mackerel.MonitorAnomalyDetection {
	monitor := &mackerel.MonitorAnomalyDetection{
		Name:                 d.Get("name").(string),
		Memo:                 d.Get("memo").(string),
		Type:                 "anomalyDetection",
		IsMute:               d.Get("is_mute").(bool),
		NotificationInterval: uint64(d.Get("notification_interval").(int)),
		WarningSensitivity:   d.Get("anomaly_detection.0.warning_sensitivity").(string),
		CriticalSensitivity:  d.Get("anomaly_detection.0.critical_sensitivity").(string),
		TrainingPeriodFrom:   uint64(d.Get("anomaly_detection.0.training_period_from").(int)),
		MaxCheckAttempts:     uint64(d.Get("anomaly_detection.0.max_check_attempts").(int)),
		Scopes:               expandStringListFromSet(d.Get("anomaly_detection.0.scopes").(*schema.Set)),
	}
	return monitor
}

func flattenMonitorAnomalyDetection(monitor *mackerel.MonitorAnomalyDetection, d *schema.ResourceData) error {
	d.Set("name", monitor.Name)
	d.Set("memo", monitor.Memo)
	d.Set("is_mute", monitor.IsMute)
	d.Set("notification_interval", monitor.NotificationInterval)
	normalizedScopes := make([]string, 0, len(monitor.Scopes))
	for _, s := range monitor.Scopes {
		normalizedScopes = append(normalizedScopes, strings.ReplaceAll(s, " ", ""))
	}
	d.Set("anomaly_detection", []map[string]interface{}{
		{
			"warning_sensitivity":  monitor.WarningSensitivity,
			"critical_sensitivity": monitor.CriticalSensitivity,
			"training_period_from": monitor.TrainingPeriodFrom,
			"max_check_attempts":   monitor.MaxCheckAttempts,
			"scopes":               flattenStringListToSet(normalizedScopes),
		},
	})
	return nil
}
