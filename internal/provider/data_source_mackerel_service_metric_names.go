package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/mackerel"
)

var (
	_ datasource.DataSourceWithConfigure = (*mackerelServiceMetricNamesDataSource)(nil)
)

func NewMackerelServiceMetricNamesDataSource() datasource.DataSource {
	return &mackerelServiceMetricNamesDataSource{}
}

type mackerelServiceMetricNamesDataSource struct {
	Client *mackerel.Client
}

func (*mackerelServiceMetricNamesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_metric_names"
}

func (*mackerelServiceMetricNamesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source allows access to details of a specific Service metric names.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the service.",

				Required:   true,
				Validators: []validator.String{mackerel.ServiceNameValidator()},
			},
			"prefix": schema.StringAttribute{
				Description: "Prefix of the metric names.",

				Optional: true,
			},
			"metric_names": schema.SetAttribute{
				Description: "Set of the service metric names.",

				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

func (d *mackerelServiceMetricNamesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	d.Client = client
}

func (d *mackerelServiceMetricNamesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config mackerel.ServiceMetricNamesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := mackerel.ReadServiceMetricNames(ctx, d.Client, config)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read Service Metric Names from Service: %s", config.Name.ValueString()),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
