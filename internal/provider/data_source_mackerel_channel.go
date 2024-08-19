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
	_ datasource.DataSource              = (*mackerelChannelDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*mackerelChannelDataSource)(nil)
)

func NewMackerelChannelDataSource() datasource.DataSource {
	return &mackerelChannelDataSource{}
}

type mackerelChannelDataSource struct {
	Client *mackerel.Client
}

func (d *mackerelChannelDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_channel"
}

func (d *mackerelChannelDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemaChannelDataSource()
}

func (d *mackerelChannelDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	d.Client = client
}

func (d *mackerelChannelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config mackerel.ChannelModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := mackerel.ReadChannel(ctx, d.Client, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read a channel",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func schemaChannelDataSource() schema.Schema {
	schema := schema.Schema{
		Description: "This data source allows access to details of a specific notification channel.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: schemaChannelIDDesc,
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: schemaChannelNameDesc,
				Computed:    true,
			},
			"email": schema.ListAttribute{
				Description: schemaChannelEmailDesc,
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"emails": types.SetType{
							ElemType: types.StringType,
						},
						"user_ids": types.SetType{
							ElemType: types.StringType,
						},
						"events": types.SetType{
							ElemType: types.StringType,
						},
					},
				},
			},
			"slack": schema.ListAttribute{
				Description: schemaChannelSlackDesc,
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"url": types.StringType,
						"mentions": types.MapType{
							ElemType: types.StringType,
						},
						"enabled_graph_image": types.BoolType,
						"events": types.SetType{
							ElemType: types.StringType,
						},
					},
				},
			},
			"webhook": schema.ListAttribute{
				Description: schemaChannelWebhookDesc,
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"url": types.StringType,
						"events": types.SetType{
							ElemType: types.StringType,
						},
					},
				},
			},
		},
	}
	return schema
}
