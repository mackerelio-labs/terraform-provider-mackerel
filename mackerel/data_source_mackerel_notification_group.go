package mackerel

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func dataSourceMackerelNotificationGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMackerelNotificationGroupRead,

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

func dataSourceMackerelNotificationGroupRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Get("id").(string)

	client := meta.(*mackerel.Client)

	groups, err := client.FindNotificationGroups()
	if err != nil {
		return err
	}
	var group *mackerel.NotificationGroup
	for _, g := range groups {
		if g.ID == id {
			group = g
			break
		}
	}
	if group == nil {
		return fmt.Errorf("the ID '%s' does not match any notification-group in mackerel.io", id)
	}
	d.SetId(group.ID)
	return flattenNotificationGroup(group, d)
}
