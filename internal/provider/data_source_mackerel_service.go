package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/mackerel"
)

var (
	_ datasource.DataSource              = (*mackerelServiceDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*mackerelServiceDataSource)(nil)
)

func NewMackerelServiceDataSource() datasource.DataSource {
	return &mackerelServiceDataSource{}
}

type mackerelServiceDataSource struct {
	Client *mackerel.Client
}

func (d *mackerelServiceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (d *mackerelServiceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source allows access to details of a specific Service.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the service.",
				Validators:  []validator.String{mackerel.ServiceNameValidator()},
			},
			"memo": schema.StringAttribute{
				Computed:    true,
				Description: "Notes related to this service.",
			},
		},
	}
}

func (d *mackerelServiceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.Client = client
}

func (d *mackerelServiceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mackerel.ServiceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	newData, err := mackerel.ReadService(ctx, d.Client, name)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read Service: %s", name),
			err.Error(),
		)
		return
	}

	data.Set(*newData)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
