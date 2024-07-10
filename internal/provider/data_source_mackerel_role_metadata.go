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
	_ datasource.DataSourceWithConfigure = (*mackerelRoleMetadataDataSource)(nil)
)

func NewMackerelRoleMetadataDataSource() datasource.DataSource {
	return &mackerelRoleMetadataDataSource{}
}

type mackerelRoleMetadataDataSource struct {
	Client *mackerel.Client
}

func (d *mackerelRoleMetadataDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role_metadata"
}

func (d *mackerelRoleMetadataDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source accesses to details of a specific Role Metadata.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"service": schema.StringAttribute{
				Description: "The name of the service.",
				Required:    true,
				Validators:  []validator.String{mackerel.ServiceNameValidator()},
			},
			"role": schema.StringAttribute{
				Description: "The name of the role.",
				Required:    true,
				Validators:  []validator.String{mackerel.RoleNameValidator()},
			},
			"namespace": schema.StringAttribute{
				Description: "The identifier for the metadata.",
				Required:    true,
			},
			"metadata_json": schema.StringAttribute{
				Description: "The arbitrary JSON data for the role.",
				Computed:    true,
				CustomType:  jsontypes.NormalizedType{},
			},
		},
	}
}

func (d *mackerelRoleMetadataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	d.Client = client
}

func (d *mackerelRoleMetadataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config mackerel.RoleMetadataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := config.ServiceName.ValueString()
	roleName := config.RoleName.ValueString()
	namespace := config.Namespace.ValueString()
	data, err := mackerel.ReadRoleMetadata(ctx, d.Client, serviceName, roleName, namespace)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read Role Metadata: service=%s role=%s namespace=%s", serviceName, roleName, namespace),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
