package mackerel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func dataSourceMackerelNotificationGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMackerelNotificationGroupRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notification_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"child_notification_group_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"child_channel_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"monitor": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     monitorResource,
			},
			"service": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     serviceResource,
			},
		},
	}
}

func dataSourceMackerelNotificationGroupRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id := d.Get("id").(string)

	client := m.(*mackerel.Client)

	groups, err := client.FindNotificationGroups()
	if err != nil {
		return diag.FromErr(err)
	}
	var group *mackerel.NotificationGroup
	for _, g := range groups {
		if g.ID == id {
			group = g
			break
		}
	}
	if group == nil {
		return diag.Errorf("the ID '%s' does not match any notification-group in mackerel.io", id)
	}
	d.SetId(group.ID)
	if err := flattenNotificationGroup(group, d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}
