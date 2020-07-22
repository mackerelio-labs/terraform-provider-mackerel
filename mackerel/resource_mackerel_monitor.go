package mackerel

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				Elem:         &schema.Resource{},
			},
			"external": {
				Type:         schema.TypeList,
				Optional:     true,
				ExactlyOneOf: []string{"host_metric", "connectivity", "service_metric", "external", "expression", "anomaly_detection"},
				MaxItems:     1,
				Elem:         &schema.Resource{},
			},
			"expression": {
				Type:         schema.TypeList,
				Optional:     true,
				ExactlyOneOf: []string{"host_metric", "connectivity", "service_metric", "external", "expression", "anomaly_detection"},
				MaxItems:     1,
				Elem:         &schema.Resource{},
			},
			"anomaly_detection": {
				Type:         schema.TypeList,
				Optional:     true,
				ExactlyOneOf: []string{"host_metric", "connectivity", "service_metric", "external", "expression", "anomaly_detection"},
				MaxItems:     1,
				Elem:         &schema.Resource{},
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
	switch v := monitor.(type) {
	case *mackerel.MonitorHostMetric:
		return flattenMonitorHostMetric(v, d)
	case *mackerel.MonitorConnectivity:
		return flattenMonitorConnectivity(v, d)
	case *mackerel.MonitorServiceMetric:
		return flattenMonitorServiceMetric(v, d)
	case *mackerel.MonitorExternalHTTP:
		return flattenMonitorExternalHTTP(v, d)
	case *mackerel.MonitorExpression:
		return flattenMonitorExpression(v, d)
	case *mackerel.MonitorAnomalyDetection:
		return flattenMonitorAnomalyDetection(v, d)
	default:
		return fmt.Errorf("unknown monitor type")
	}
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

// todo
func expandMonitorServiceMetric(d *schema.ResourceData) *mackerel.MonitorServiceMetric {
	monitor := &mackerel.MonitorServiceMetric{}

	return monitor
}

// todo
func flattenMonitorServiceMetric(monitor *mackerel.MonitorServiceMetric, d *schema.ResourceData) error {
	return nil
}

// todo
func expandMonitorExternalHTTP(d *schema.ResourceData) *mackerel.MonitorExternalHTTP {
	monitor := &mackerel.MonitorExternalHTTP{}

	return monitor
}

// todo
func flattenMonitorExternalHTTP(monitor *mackerel.MonitorExternalHTTP, d *schema.ResourceData) error {
	return nil
}

// todo
func expandMonitorExpression(d *schema.ResourceData) *mackerel.MonitorExpression {
	monitor := &mackerel.MonitorExpression{}

	return monitor
}

// todo
func flattenMonitorExpression(monitor *mackerel.MonitorExpression, d *schema.ResourceData) error {
	return nil
}

// todo
func expandMonitorAnomalyDetection(d *schema.ResourceData) *mackerel.MonitorAnomalyDetection {
	monitor := &mackerel.MonitorAnomalyDetection{}

	return monitor
}

// todo
func flattenMonitorAnomalyDetection(monitor *mackerel.MonitorAnomalyDetection, d *schema.ResourceData) error {
	return nil
}
