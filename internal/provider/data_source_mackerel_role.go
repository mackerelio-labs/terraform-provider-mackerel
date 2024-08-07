package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/mackerel"
)

var (
	_ datasource.DataSource              = (*mackerelRoleDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*mackerelRoleDataSource)(nil)
)

func NewMackerelRoleDataSource() datasource.DataSource {
	return &mackerelRoleDataSource{}
}

type mackerelRoleDataSource struct {
	Client *mackerel.Client
}

func (d *mackerelRoleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (d *mackerelRoleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source allows access to details of a specific Role.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"service": schema.StringAttribute{
				Description: "The name of the service.",
				Required:    true,
				Validators:  []validator.String{mackerel.ServiceNameValidator()},
			},
			"name": schema.StringAttribute{
				Description: "The name of the role.",
				Required:    true,
				Validators:  []validator.String{mackerel.RoleNameValidator()},
			},
			"memo": schema.StringAttribute{
				Description: "Notes related to this role.",
				Computed:    true,
			},
		},
	}
}

func (d *mackerelRoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	d.Client = client
}

func (d *mackerelRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config mackerel.RoleModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := mackerel.ReadRole(ctx, d.Client, config.ServiceName.ValueString(), config.RoleName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Role",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
