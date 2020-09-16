package mackerel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func dataSourceMackerelAlertGroupSetting() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMackerelAlertGroupSettingRead,

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
			"service_scopes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"role_scopes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"monitor_scopes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"notification_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceMackerelAlertGroupSettingRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id := d.Get("id").(string)

	client := m.(*mackerel.Client)

	group, err := client.GetAlertGroupSetting(id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(group.ID)
	if err := flattenAlertGroupSetting(group, d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}
