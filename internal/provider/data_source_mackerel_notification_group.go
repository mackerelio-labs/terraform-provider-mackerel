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
	_ datasource.DataSource              = (*mackerelNotificationGroupDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*mackerelNotificationGroupDataSource)(nil)
)

func NewMackerelNotificationGroupDataSource() datasource.DataSource {
	return &mackerelNotificationGroupDataSource{}
}

type mackerelNotificationGroupDataSource struct {
	Client *mackerel.Client
}

func (d *mackerelNotificationGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_group"
}

func (d *mackerelNotificationGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"notification_level": schema.StringAttribute{
				Computed: true,
			},
			"child_notification_group_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"child_channel_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"monitor": schema.SetAttribute{
				Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"id":           types.StringType,
						"skip_default": types.BoolType,
					},
				},
			},
			"service": schema.SetAttribute{
				Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name": types.StringType,
					},
				},
			},
		},
	}
}

func (d *mackerelNotificationGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	d.Client = client
}

func (d *mackerelNotificationGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config mackerel.NotificationGroupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := mackerel.ReadNotificationGroup(ctx, d.Client, config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Notification Group.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
