package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/mackerel"
)

var (
	_ datasource.DataSourceWithConfigure = (*mackerelServiceMetadataDataSource)(nil)
)

func NewMackerelServiceMetadataDataSource() datasource.DataSource {
	return &mackerelServiceMetadataDataSource{}
}

type mackerelServiceMetadataDataSource struct {
	Client *mackerel.Client
}

func (d *mackerelServiceMetadataDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_metadata"
}

func (d *mackerelServiceMetadataDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source allows access to details of a specific Service Metadata.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"service": schema.StringAttribute{
				Description: "The name of the service.",

				Required:   true,
				Validators: []validator.String{mackerel.ServiceNameValidator()},
			},
			"namespace": schema.StringAttribute{
				Description: "Identifier for the metadata.",

				Required: true,
			},
			"metadata_json": schema.StringAttribute{
				Description: "Arbitrary JSON data for the service.",

				Computed:   true,
				CustomType: jsontypes.NormalizedType{},
			},
		},
	}
}

func (d *mackerelServiceMetadataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	d.Client = client
}

func (d *mackerelServiceMetadataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mackerel.ServiceMetadataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	remoteData, err := mackerel.ReadServiceMetadata(ctx, d.Client, data)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read Service Metadata: %s/%s", data.ServiceName.ValueString(), data.Namespace.ValueString()),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &remoteData)...)
}
