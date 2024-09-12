package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/mackerel"
)

var (
	_ datasource.DataSource              = (*mackerelAlertGroupSettingDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*mackerelAlertGroupSettingDataSource)(nil)
)

func NewMackerelAlertGroupSettingDataSource() datasource.DataSource {
	return &mackerelAlertGroupSettingDataSource{}
}

type mackerelAlertGroupSettingDataSource struct {
	Client *mackerel.Client
}

func (d *mackerelAlertGroupSettingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_group_setting"
}

func (d *mackerelAlertGroupSettingDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemaAlertGroupSettingDataSource
}

func (d *mackerelAlertGroupSettingDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	d.Client = client
}

func (d *mackerelAlertGroupSettingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config mackerel.AlertGroupSettingModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := mackerel.ReadAlertGroupSetting(ctx, d.Client, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read an alert group setting.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

var schemaAlertGroupSettingDataSource = schema.Schema{
	Description: "This data source allows access to details of a specific alert group setting.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: schemaAlertGroupSettingIDDesc,
			Required:    true,
		},
		"name": schema.StringAttribute{
			Description: schemaAlertGroupSettingNameDesc,
			Computed:    true,
		},
		"memo": schema.StringAttribute{
			Description: schemaAlertGroupSettingMemoDesc,
			Computed:    true,
		},
		"service_scopes": schema.SetAttribute{
			ElementType: types.StringType,
			Description: schemaAlertGroupSettingServiceScopesDesc,
			Computed:    true,
		},
		"role_scopes": schema.SetAttribute{
			ElementType: types.StringType,
			Description: schemaAlertGroupSettingRoleScopesDesc,
			Computed:    true,
		},
		"monitor_scopes": schema.SetAttribute{
			ElementType: types.StringType,
			Description: schemaAlertGroupSettingMonitorScopesDesc,
			Computed:    true,
		},
		"notification_interval": schema.Int64Attribute{
			Description: schemaAlertGroupSettingNotificationIntervalDesc,
			Computed:    true,
		},
	},
}
