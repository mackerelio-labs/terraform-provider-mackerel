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
	_ datasource.DataSource              = (*mackerelDashboardDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*mackerelDashboardDataSource)(nil)
)

type mackerelDashboardDataSource struct {
	Client *mackerel.Client
}

func NewMackerelDashboardDataSource() datasource.DataSource {
	return &mackerelDashboardDataSource{}
}

func (d *mackerelDashboardDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dashboard"
}

func (d *mackerelDashboardDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemaDashboardDataSource()
}

func (d *mackerelDashboardDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	d.Client = client
}

func (d *mackerelDashboardDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config mackerel.DashboardModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := mackerel.ReadDashboard(ctx, d.Client, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read a dashboard",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func schemaDashboardDataSource() schema.Schema {
	layoutType := types.ListType{
		ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"x":      types.Int64Type,
				"y":      types.Int64Type,
				"width":  types.Int64Type,
				"height": types.Int64Type,
			},
		},
	}
	return schema.Schema{
		Description: "This data source allows access to details of a specific dashboard.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: schemaDashboardIDDesc,
				Required:    true,
			},
			"title": schema.StringAttribute{
				Description: schemaDashboardTitleDesc,
				Computed:    true,
			},
			"memo": schema.StringAttribute{
				Description: schemaDashboardMemoDesc,
				Computed:    true,
			},
			"url_path": schema.StringAttribute{
				Description: schemaDashboardURLPathDesc,
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: schemaDashboardCreatedAtDesc,
				Computed:    true,
			},
			"updated_at": schema.Int64Attribute{
				Description: schemaDashboardUpdatedAtDesc,
				Computed:    true,
			},
			"graph": schema.ListAttribute{
				Description: schemaDashboardGraphDesc,
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"title":  types.StringType,
						"layout": layoutType,
						"range": types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"relative": types.ListType{
										ElemType: types.ObjectType{
											AttrTypes: map[string]attr.Type{
												"period": types.Int64Type,
												"offset": types.Int64Type,
											},
										},
									},
									"absolute": types.ListType{
										ElemType: types.ObjectType{
											AttrTypes: map[string]attr.Type{
												"start": types.Int64Type,
												"end":   types.Int64Type,
											},
										},
									},
								},
							},
						},

						"host": types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"host_id": types.StringType,
									"name":    types.StringType,
								},
							},
						},
						"role": types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"role_fullname": types.StringType,
									"name":          types.StringType,
									"is_stacked":    types.BoolType,
								},
							},
						},
						"service": types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"service_name": types.StringType,
									"name":         types.StringType,
								},
							},
						},
						"expression": types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"expression": types.StringType,
								},
							},
						},
						"query": types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"query":  types.StringType,
									"legend": types.StringType,
								},
							},
						},
					},
				},
			},
			"value": schema.ListAttribute{
				Description: schemaDashboardValueDesc,
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"title":         types.StringType,
						"layout":        layoutType,
						"fraction_size": types.Int64Type,
						"suffix":        types.StringType,
						"metric": types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"host": types.ListType{
										ElemType: types.ObjectType{
											AttrTypes: map[string]attr.Type{
												"host_id": types.StringType,
												"name":    types.StringType,
											},
										},
									},
									"service": types.ListType{
										ElemType: types.ObjectType{
											AttrTypes: map[string]attr.Type{
												"service_name": types.StringType,
												"name":         types.StringType,
											},
										},
									},
									"expression": types.ListType{
										ElemType: types.ObjectType{
											AttrTypes: map[string]attr.Type{
												"expression": types.StringType,
											},
										},
									},
									"query": types.ListType{
										ElemType: types.ObjectType{
											AttrTypes: map[string]attr.Type{
												"query":  types.StringType,
												"legend": types.StringType,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"markdown": schema.ListAttribute{
				Description: schemaDashboardMarkdownDesc,
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"title":    types.StringType,
						"layout":   layoutType,
						"markdown": types.StringType,
					},
				},
			},
			"alert_status": schema.ListAttribute{
				Description: schemaDashboardAlertStatusDesc,
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"title":         types.StringType,
						"layout":        layoutType,
						"role_fullname": types.StringType,
					},
				},
			},
		},
	}
}
