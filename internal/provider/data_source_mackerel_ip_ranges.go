package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/mackerel"
)

var (
	_ datasource.DataSourceWithConfigure = (*mackerelIPRangesDataSource)(nil)
)

func NewMackerelIPRangesDataSource() datasource.DataSource {
	return &mackerelIPRangesDataSource{}
}

type mackerelIPRangesDataSource struct {
	Client *mackerel.Client
}

func (*mackerelIPRangesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ip_ranges"
}

func (*mackerelIPRangesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source allows access to the IP address ranges used by Mackerel.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"incoming": schema.ListAttribute{
				Description: "IPv4 address ranges that Mackerel receives.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"incoming_ipv6": schema.ListAttribute{
				Description: "IPv6 address ranges that Mackerel receives.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"outgoing": schema.ListAttribute{
				Description: "IPv4 address ranges that Mackerel sends from (e.g. alert notifications).",
				ElementType: types.StringType,
				Computed:    true,
			},
			"outgoing_ipv6": schema.ListAttribute{
				Description: "IPv6 address ranges that Mackerel sends from (e.g. alert notifications).",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

func (d *mackerelIPRangesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	d.Client = client
}

func (d *mackerelIPRangesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	data, err := mackerel.ReadIPRanges(ctx, d.Client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Mackerel IP ranges",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
