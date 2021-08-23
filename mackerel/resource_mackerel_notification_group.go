package mackerel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mackerelio/mackerel-client-go"
)

var serviceResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
	},
}

var monitorResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"skip_default": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	},
}

func resourceMackerelNotificationGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMackerelNotificationGroupCreate,
		ReadContext:   resourceMackerelNotificationGroupRead,
		UpdateContext: resourceMackerelNotificationGroupUpdate,
		DeleteContext: resourceMackerelNotificationGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"notification_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"all", "critical"}, false),
				Default:      "all",
			},
			"child_notification_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"child_channel_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"monitor": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     monitorResource,
			},
			"service": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     serviceResource,
			},
		},
	}
}

func resourceMackerelNotificationGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mackerel.Client)
	group, err := client.CreateNotificationGroup(expandNotificationGroup(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(group.ID)
	return resourceMackerelNotificationGroupRead(ctx, d, m)
}

func resourceMackerelNotificationGroupRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mackerel.Client)
	groups, err := client.FindNotificationGroups()
	if err != nil {
		return diag.FromErr(err)
	}
	var group *mackerel.NotificationGroup
	for _, g := range groups {
		if g.ID == d.Id() {
			group = g
			break
		}
	}
	if group == nil {
		return diag.Errorf("the ID '%s' does not match any notification-group in mackerel.io", d.Id())
	}
	return flattenNotificationGroup(group, d)
}

func resourceMackerelNotificationGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mackerel.Client)
	group, err := client.UpdateNotificationGroup(d.Id(), expandNotificationGroup(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(group.ID)
	return resourceMackerelNotificationGroupRead(ctx, d, m)
}

func resourceMackerelNotificationGroupDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*mackerel.Client)
	_, err := client.DeleteNotificationGroup(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func expandNotificationGroup(d *schema.ResourceData) *mackerel.NotificationGroup {
	return &mackerel.NotificationGroup{
		Name:                      d.Get("name").(string),
		NotificationLevel:         mackerel.NotificationLevel(d.Get("notification_level").(string)),
		ChildNotificationGroupIDs: expandStringListFromSet(d.Get("child_notification_group_ids").(*schema.Set)),
		ChildChannelIDs:           expandStringListFromSet(d.Get("child_channel_ids").(*schema.Set)),
		Monitors:                  expandMonitorSet(d.Get("monitor").(*schema.Set)),
		Services:                  expandServiceSet(d.Get("service").(*schema.Set)),
	}
}

func expandMonitorSet(set *schema.Set) []*mackerel.NotificationGroupMonitor {
	monitors := make([]*mackerel.NotificationGroupMonitor, 0, set.Len())
	for _, monitor := range set.List() {
		monitor := monitor.(map[string]interface{})
		monitors = append(monitors, &mackerel.NotificationGroupMonitor{
			ID:          monitor["id"].(string),
			SkipDefault: monitor["skip_default"].(bool),
		})
	}
	return monitors
}

func expandServiceSet(set *schema.Set) []*mackerel.NotificationGroupService {
	services := make([]*mackerel.NotificationGroupService, 0, set.Len())
	for _, service := range set.List() {
		service := service.(map[string]interface{})
		services = append(services, &mackerel.NotificationGroupService{Name: service["name"].(string)})
	}
	return services
}
