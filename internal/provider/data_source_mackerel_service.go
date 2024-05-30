package provider

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/mackerelio/mackerel-client-go"
)

var (
	_ datasource.DataSource              = (*mackerelServiceDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*mackerelServiceDataSource)(nil)
)

func NewMackerelServiceDataSource() datasource.DataSource {
	return &mackerelServiceDataSource{}
}

type mackerelServiceDataSource struct {
	client *mackerel.Client
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
	d.client = client
}

func (d *mackerelServiceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mackerelServiceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	services, err := d.client.FindServices()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read services",
			fmt.Sprintf("An unexpected error occurred while attempting to fetch the services: %v", err),
		)
		return
	}

	name := data.Name.ValueString()
	serviceIdx := slices.IndexFunc(services, func(s *mackerel.Service) bool {
		return s.Name == name
	})
	if serviceIdx == -1 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("No Service Found: %s", name),
			// FIXME: for backwards compatibility
			fmt.Sprintf("the name '%s' does not match any service in mackerel.io", name),
			// fmt.Sprintf("The name '%s' does not match any service in mackerel.io", name),
		)
		return
	}

	data.SetService(services[serviceIdx])
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
