package mackerel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MACKEREL_API_KEY", nil),
				Description: "Mackerel API Key",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"mackerel_alert_group_setting": resourceMackerelAlertGroupSetting(),
			"mackerel_channel":             resourceMackerelChannel(),
			"mackerel_downtime":            resourceMackerelDowntime(),
			"mackerel_monitor":             resourceMackerelMonitor(),
			"mackerel_notification_group":  resourceMackerelNotificationGroup(),
			"mackerel_role":                resourceMackerelRole(),
			"mackerel_role_metadata":       resourceMackerelRoleMetadata(),
			"mackerel_service":             resourceMackerelService(),
			"mackerel_service_metadata":    resourceMackerelServiceMetadata(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"mackerel_alert_group_setting": dataSourceMackerelAlertGroupSetting(),
			"mackerel_channel":             dataSourceMackerelChannel(),
			"mackerel_downtime":            dataSourceMackerelDowntime(),
			"mackerel_monitor":             dataSourceMackerelMonitor(),
			"mackerel_notification_group":  dataSourceMackerelNotificationGroup(),
			"mackerel_role":                dataSourceMackerelRole(),
			"mackerel_role_metadata":       dataSourceMackerelRoleMetadata(),
			"mackerel_service":             dataSourceMackerelService(),
			"mackerel_service_metadata":    dataSourceMackerelServiceMetadata(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(data *schema.ResourceData) (interface{}, error) {
	config := Config{
		APIKey: data.Get("api_key").(string),
	}
	return config.Client()
}
