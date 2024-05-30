package mackerel

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	mackerelfwprovider "github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
)

// Provider returns a *schema.Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"MACKEREL_APIKEY",
					"MACKEREL_API_KEY",
				}, nil),
				Description: "Mackerel API Key",
				Sensitive:   true,
			},
			"api_base": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("API_BASE", nil),
				Description:  "Mackerel API BASE URL",
				Sensitive:    true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"mackerel_alert_group_setting": resourceMackerelAlertGroupSetting(),
			"mackerel_aws_integration":     resourceMackerelAWSIntegration(),
			"mackerel_channel":             resourceMackerelChannel(),
			"mackerel_downtime":            resourceMackerelDowntime(),
			"mackerel_dashboard":           resourceMackerelDashboard(),
			"mackerel_monitor":             resourceMackerelMonitor(),
			"mackerel_notification_group":  resourceMackerelNotificationGroup(),
			"mackerel_role":                resourceMackerelRole(),
			"mackerel_role_metadata":       resourceMackerelRoleMetadata(),
			"mackerel_service":             resourceMackerelService(),
			"mackerel_service_metadata":    resourceMackerelServiceMetadata(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"mackerel_alert_group_setting":  dataSourceMackerelAlertGroupSetting(),
			"mackerel_aws_integration":      dataSourceMackerelAWSIntegration(),
			"mackerel_channel":              dataSourceMackerelChannel(),
			"mackerel_dashboard":            dataSourceMackerelDashboard(),
			"mackerel_downtime":             dataSourceMackerelDowntime(),
			"mackerel_monitor":              dataSourceMackerelMonitor(),
			"mackerel_notification_group":   dataSourceMackerelNotificationGroup(),
			"mackerel_role":                 dataSourceMackerelRole(),
			"mackerel_role_metadata":        dataSourceMackerelRoleMetadata(),
			"mackerel_service":              dataSourceMackerelService(),
			"mackerel_service_metadata":     dataSourceMackerelServiceMetadata(),
			"mackerel_service_metric_names": dataSourceMackerelServiceMetricNames(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func ProtoV5ProviderServer() tfprotov5.ProviderServer {
	return protoV5ProviderServer(Provider())
}

func protoV5ProviderServer(provider *schema.Provider) tfprotov5.ProviderServer {
	fwFlag := os.Getenv("MACKEREL_EXPERIMENTAL_TFFRAMEWORK")
	if fwFlag == "1" || fwFlag == "true" {
		log.Printf("[INFO] mackerel: use terraform-plugin-framework based implementation")

		// Resources
		delete(provider.ResourcesMap, "mackerel_service")

		// Data Sources
		delete(provider.DataSourcesMap, "mackerel_service")

		mux, err := tf5muxserver.NewMuxServer(
			context.Background(),
			providerserver.NewProtocol5(mackerelfwprovider.New()),
			provider.GRPCProvider,
		)
		if err != nil {
			panic(err)
		}
		return mux.ProviderServer()
	}

	return provider.GRPCProvider()
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		APIKey:  d.Get("api_key").(string),
		APIBase: d.Get("api_base").(string),
	}
	return config.Client()
}
