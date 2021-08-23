package mackerel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelAlertGroupSetting() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMackerelAlertGroupSettingCreate,
		ReadContext:   resourceMackerelAlertGroupSettingRead,
		UpdateContext: resourceMackerelAlertGroupSettingUpdate,
		DeleteContext: resourceMackerelAlertGroupSettingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"service_scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"role_scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"monitor_scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"notification_interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceMackerelAlertGroupSettingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mackerel.Client)
	setting, err := client.CreateAlertGroupSetting(expandAlertGroupSetting(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(setting.ID)
	return resourceMackerelAlertGroupSettingRead(ctx, d, m)
}

func resourceMackerelAlertGroupSettingRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mackerel.Client)
	setting, err := client.GetAlertGroupSetting(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return flattenAlertGroupSetting(setting, d)
}

func resourceMackerelAlertGroupSettingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mackerel.Client)
	setting, err := client.UpdateAlertGroupSetting(d.Id(), expandAlertGroupSetting(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(setting.ID)
	return resourceMackerelAlertGroupSettingRead(ctx, d, m)
}

func resourceMackerelAlertGroupSettingDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*mackerel.Client)
	_, err := client.DeleteAlertGroupSetting(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func expandAlertGroupSetting(d *schema.ResourceData) *mackerel.AlertGroupSetting {
	return &mackerel.AlertGroupSetting{
		Name:                 d.Get("name").(string),
		Memo:                 d.Get("memo").(string),
		ServiceScopes:        expandStringListFromSet(d.Get("service_scopes").(*schema.Set)),
		RoleScopes:           expandStringListFromSet(d.Get("role_scopes").(*schema.Set)),
		MonitorScopes:        expandStringListFromSet(d.Get("monitor_scopes").(*schema.Set)),
		NotificationInterval: uint64(d.Get("notification_interval").(int)),
	}
}
