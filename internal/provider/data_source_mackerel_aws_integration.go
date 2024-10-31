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
	_ datasource.DataSource              = (*mackerelAWSIntegrationDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*mackerelAWSIntegrationDataSource)(nil)
)

type mackerelAWSIntegrationDataSource struct {
	Client *mackerel.Client
}

func NewMackerelAWSIntegrationDataSource() datasource.DataSource {
	return &mackerelAWSIntegrationDataSource{}
}

func (d *mackerelAWSIntegrationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aws_integration"
}

func (d *mackerelAWSIntegrationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemaAWSIntegrationDataSource()
}

func (d *mackerelAWSIntegrationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	d.Client = client
}

func (d *mackerelAWSIntegrationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config mackerel.AWSIntegrationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := mackerel.ReadAWSIntegration(ctx, d.Client, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read an AWS integration",
			err.Error(),
		)
		return
	}
	dataSourceModel := mackerel.AWSIntegrationDataSourceModel(*data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &dataSourceModel)...)
}

func schemaAWSIntegrationDataSource() schema.Schema {
	attrs := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: schemaAWSIntegrationIDDesc,
			Required:    true,
		},
		"name": schema.StringAttribute{
			Description: schemaAWSIntegrationNameDesc,
			Computed:    true,
		},
		"memo": schema.StringAttribute{
			Description: schemaAWSIntegrationMemoDesc,
			Computed:    true,
		},
		"key": schema.StringAttribute{
			Description: schemaAWSIntegrationKeyDesc,
			Computed:    true,
		},
		"role_arn": schema.StringAttribute{
			Description: schemaAWSIntegrationRoleARNDesc,
			Computed:    true,
		},
		"external_id": schema.StringAttribute{
			Description: schemaAWSIntegrationExternalIDDesc,
			Computed:    true,
		},
		"region": schema.StringAttribute{
			Description: schemaAWSIntegrationRegionDesc,
			Computed:    true,
		},
		"included_tags": schema.StringAttribute{
			Description: schemaAWSIntegrationIncludedTagsDesc,
			Computed:    true,
		},
		"excluded_tags": schema.StringAttribute{
			Description: schemaAWSIntegrationExcludedTagsDesc,
			Computed:    true,
		},
	}

	for name, spec := range awsIntegrationServices {
		if spec.supportsAutoRetire {
			attrs[name] = schema.SetAttribute{
				Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"enable": types.BoolType,
						"role":   types.StringType,
						"excluded_metrics": types.ListType{
							ElemType: types.StringType,
						},
						"retire_automatically": types.BoolType,
					},
				},
			}
		} else {
			attrs[name] = schema.SetAttribute{
				Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"enable": types.BoolType,
						"role":   types.StringType,
						"excluded_metrics": types.ListType{
							ElemType: types.StringType,
						},
					},
				},
			}
		}
	}

	s := schema.Schema{
		Attributes: attrs,
	}
	return s
}
