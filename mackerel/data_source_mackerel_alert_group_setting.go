package mackerel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func dataSourceMackerelAlertGroupSetting() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMackerelAlertGroupSettingRead,

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

func dataSourceMackerelAlertGroupSettingRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Get("id").(string)

	client := meta.(*mackerel.Client)

	group, err := client.GetAlertGroupSetting(id)
	if err != nil {
		return err
	}
	d.SetId(group.ID)
	return flattenAlertGroupSetting(group, d)
}
