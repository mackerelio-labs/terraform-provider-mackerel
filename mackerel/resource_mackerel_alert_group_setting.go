package mackerel

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelAlertGroupSetting() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelAlertGroupSettingCreate,
		Read:   resourceMackerelAlertGroupSettingRead,
		Update: resourceMackerelAlertGroupSettingUpdate,
		Delete: resourceMackerelAlertGroupSettingDelete,
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

func resourceMackerelAlertGroupSettingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	setting, err := client.CreateAlertGroupSetting(expandAlertGroupSetting(d))
	if err != nil {
		return err
	}
	d.SetId(setting.ID)
	return resourceMackerelAlertGroupSettingRead(d, meta)
}

func resourceMackerelAlertGroupSettingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	setting, err := client.GetAlertGroupSetting(d.Id())
	if err != nil {
		return err
	}
	return flattenAlertGroupSetting(setting, d)
}

func resourceMackerelAlertGroupSettingUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	setting, err := client.UpdateAlertGroupSetting(d.Id(), expandAlertGroupSetting(d))
	if err != nil {
		return err
	}
	d.SetId(setting.ID)
	return resourceMackerelAlertGroupSettingRead(d, meta)
}

func resourceMackerelAlertGroupSettingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	_, err := client.DeleteAlertGroupSetting(d.Id())
	return err
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

func flattenAlertGroupSetting(setting *mackerel.AlertGroupSetting, d *schema.ResourceData) error {
	d.Set("name", setting.Name)
	d.Set("memo", setting.Memo)
	d.Set("service_scopes", flattenStringListToSet(setting.ServiceScopes))
	normalizedRoleScopes := make([]string, 0, len(setting.RoleScopes))
	for _, r := range setting.RoleScopes {
		normalizedRoleScopes = append(normalizedRoleScopes, strings.ReplaceAll(r, " ", ""))
	}
	d.Set("role_scopes", flattenStringListToSet(normalizedRoleScopes))
	d.Set("monitor_scopes", flattenStringListToSet(setting.MonitorScopes))
	d.Set("notification_interval", setting.NotificationInterval)
	return nil
}
