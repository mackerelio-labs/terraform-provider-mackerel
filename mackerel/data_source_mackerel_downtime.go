package mackerel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func dataSourceMackerelDowntime() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMackerelDowntimeRead,

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
			"start": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"recurrence": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"weekdays": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"until": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"service_scopes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"service_exclude_scopes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"role_scopes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"role_exclude_scopes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"monitor_scopes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"monitor_exclude_scopes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceMackerelDowntimeRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Get("id").(string)

	client := m.(*mackerel.Client)

	downtimes, err := client.FindDowntimes()
	if err != nil {
		return diag.FromErr(err)
	}
	var downtime *mackerel.Downtime
	for _, dt := range downtimes {
		if dt.ID == id {
			downtime = dt
			break
		}
	}
	if downtime == nil {
		return diag.Errorf("the ID '%s' does not match any downtime in mackerel.io", id)
	}
	d.SetId(downtime.ID)
	return flattenDowntime(downtime, d)
}
