package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/mackerel"
)

var (
	_ datasource.DataSource              = (*mackerelDowntimeDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*mackerelDowntimeDataSource)(nil)
)

func NewMackerelDowntimeDataSource() datasource.DataSource {
	return &mackerelDowntimeDataSource{}
}

type mackerelDowntimeDataSource struct {
	Client *mackerel.Client
}

func (d *mackerelDowntimeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_downtime"
}

func (d *mackerelDowntimeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemaDowntimeDataSource()
}

func (d *mackerelDowntimeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	d.Client = client
}

func (d *mackerelDowntimeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config mackerel.DowntimeModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := mackerel.ReadDowntime(ctx, d.Client, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read a downtime",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func schemaDowntimeDataSource() schema.Schema {
	s := schema.Schema{
		Description: "This data source allows access to details of a specific downtime setting.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: schemaDowntimeIDDesc,
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: schemaDowntimeNameDesc,
				Computed:    true,
			},
			"memo": schema.StringAttribute{
				Description: schemaDowntimeMemoDesc,
				Computed:    true,
			},
			"start": schema.Int64Attribute{
				Description: schemaDowntimeStartDesc,
				Computed:    true,
			},
			"duration": schema.Int64Attribute{
				Description: schemaDowntimeDurationDesc,
				Computed:    true,
			},
			"recurrence": schema.ListAttribute{
				Description: schemaDowntimeRecurrenceDesc,
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"type":     types.StringType,
						"interval": types.Int64Type,
						"weekdays": types.SetType{
							ElemType: types.StringType,
						},
						"until": types.Int64Type,
					},
				},
			},
			"service_scopes": schema.SetAttribute{
				Description: schemaDowntimeServiceScopesDesc,
				ElementType: types.StringType,
				Computed:    true,
			},
			"service_exclude_scopes": schema.SetAttribute{
				Description: schemaDowntimeServiceExcludeScopesDesc,
				ElementType: types.StringType,
				Computed:    true,
			},
			"role_scopes": schema.SetAttribute{
				Description: schemaDowntimeRoleScopesDesc,
				ElementType: types.StringType,
				Computed:    true,
			},
			"role_exclude_scopes": schema.SetAttribute{
				Description: schemaDowntimeRoleExcludeScopesDesc,
				ElementType: types.StringType,
				Computed:    true,
			},
			"monitor_scopes": schema.SetAttribute{
				Description: schemaDowntimeMonitorScopesDesc,
				ElementType: types.StringType,
				Computed:    true,
			},
			"monitor_exclude_scopes": schema.SetAttribute{
				Description: schemaDowntimeMonitorExcludeScopesDesc,
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
	return s
}
