package mackerel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func dataSourceMackerelMonitor() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMackerelMonitorRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"memo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_mute": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"notification_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"host_metric": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"warning": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"critical": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"duration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_check_attempts": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scopes": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"exclude_scopes": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"connectivity": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scopes": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"exclude_scopes": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"service_metric": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metric": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"warning": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"critical": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"duration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_check_attempts": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"missing_duration_warning": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"missing_duration_critical": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"external": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_check_attempts": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"service": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"response_time_critical": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"response_time_warning": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"response_time_duration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"request_body": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"contains_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certification_expiration_critical": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"certification_expiration_warning": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"skip_certificate_verification": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"headers": {
							Type:      schema.TypeMap,
							Computed:  true,
							Sensitive: true,
							Elem:      &schema.Schema{Type: schema.TypeString},
						},
						"follow_redirect": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"expression": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expression": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"warning": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"critical": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"anomaly_detection": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"warning_sensitivity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"critical_sensitivity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_check_attempts": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"training_period_from": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scopes": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceMackerelMonitorRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Get("id").(string)

	client := m.(*mackerel.Client)
	monitor, err := client.GetMonitor(id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(monitor.MonitorID())
	return flattenMonitor(monitor, d)
}
