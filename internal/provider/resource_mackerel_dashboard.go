package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/mackerel"
)

var (
	_ resource.Resource                = (*mackerelDashboardResource)(nil)
	_ resource.ResourceWithConfigure   = (*mackerelDashboardResource)(nil)
	_ resource.ResourceWithImportState = (*mackerelDashboardResource)(nil)
)

func NewMackerelDashboardResource() resource.Resource {
	return &mackerelDashboardResource{}
}

type mackerelDashboardResource struct {
	Client *mackerel.Client
}

func (r *mackerelDashboardResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dashboard"
}

func (r *mackerelDashboardResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemaDashboardResource()
}

func (r *mackerelDashboardResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	r.Client = client
}

func (r *mackerelDashboardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mackerel.DashboardModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Create(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to create a dashboard",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelDashboardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mackerel.DashboardModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Read(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to read a dashboard",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelDashboardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mackerel.DashboardModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Update(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to update a dashboard",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelDashboardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mackerel.DashboardModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Delete(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete a dashboard",
			err.Error(),
		)
		return
	}
}

func (r *mackerelDashboardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

const (
	schemaDashboardIDDesc        = "The ID of the dashboard."
	schemaDashboardTitleDesc     = "The name of the dashboard."
	schemaDashboardMemoDesc      = "The notes regarding the dashboard."
	schemaDashboardURLPathDesc   = "The URL path for the dashboard."
	schemaDashboardCreatedAtDesc = "The time (in epoch seconds) at the dashboard created."
	schemaDashboardUpdatedAtDesc = "The time (in epoch seconds) at the dashboard last updated."
)

func schemaDashboardResource() schema.Schema {
	s := schema.Schema{
		Description: "This resource allows creating and management of dashboards.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: schemaDashboardIDDesc,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Description: schemaDashboardTitleDesc,
				Required:    true,
			},
			"memo": schema.StringAttribute{
				Description: schemaDashboardMemoDesc,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"url_path": schema.StringAttribute{
				Description: schemaDashboardURLPathDesc,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"created_at": schema.Int64Attribute{
				Description: schemaDashboardCreatedAtDesc,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.Int64Attribute{
				Description: schemaDashboardUpdatedAtDesc,
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"graph":        schemaDashboardResource_graph,
			"value":        schemaDashboardResource_value,
			"markdown":     schemaDashboardResource_markdown,
			"alert_status": schemaDashboardResource_alertStatus,
		},
	}
	return s
}

const (
	schemaDashboardWidget_titleDesc  = "The title of the widget."
	schemaDashboardWidget_layoutDesc = "The layout of the widget."
)

var schemaDashboardResource_widgetTitle = schema.StringAttribute{
	Description: schemaDashboardWidget_titleDesc,
	Required:    true,
}
var schemaDashboardResource_widgetLayout = schema.ListNestedBlock{
	Description: schemaDashboardWidget_layoutDesc,
	Validators: []validator.List{
		listvalidator.SizeBetween(1, 1), // Required
	},
	NestedObject: schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"x": schema.Int64Attribute{
				Required: true,
			},
			"y": schema.Int64Attribute{
				Required: true,
			},
			"width": schema.Int64Attribute{
				Required: true,
			},
			"height": schema.Int64Attribute{
				Required: true,
			},
		},
	},
}

const (
	schemaDashboardGraph_rangeDesc               = "The graph display range."
	schemaDashboardGraph_rangeRelativeDesc       = "The relative display range. From (current time + `offset` - `period`) to (current time + `offset`)."
	schemaDashboardGraph_rangeRelativePeriodDesc = "The length of the period (in seconds)."
	schemaDashboardGraph_rangeRelativeOffsetDesc = "The difference from the current time (in seconds)."
	schemaDashboardGraph_rangeAbsoluteDesc       = "The absolute display range. From `start` to `end`."
	schemaDashboardGraph_rangeAbsoluteStartDesc  = "The start time (in epoch seconds)."
	schemaDashboardGraph_rangeAbsoluteEndDesc    = "The end time (in epoch seconds)."
)

var schemaDashboardResource_widgetRange = schema.ListNestedBlock{
	Description: schemaDashboardGraph_rangeDesc,
	Validators: []validator.List{
		listvalidator.SizeAtMost(1), // Single Optional
	},
	NestedObject: schema.NestedBlockObject{
		Validators: []validator.Object{},
		Blocks: map[string]schema.Block{
			"relative": schema.ListNestedBlock{
				Description: schemaDashboardGraph_rangeRelativeDesc,
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.ExactlyOneOf(
						path.MatchRelative(),
						path.MatchRelative().AtParent().AtName("absolute"),
					),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"period": schema.Int64Attribute{
							Description: schemaDashboardGraph_rangeRelativePeriodDesc,
							Required:    true,
						},
						"offset": schema.Int64Attribute{
							Description: schemaDashboardGraph_rangeRelativeOffsetDesc,
							Required:    true,
						},
					},
				},
			},
			"absolute": schema.ListNestedBlock{
				Description: schemaDashboardGraph_rangeAbsoluteDesc,
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"start": schema.Int64Attribute{
							Description: schemaDashboardGraph_rangeAbsoluteStartDesc,
							Required:    true,
						},
						"end": schema.Int64Attribute{
							Description: schemaDashboardGraph_rangeAbsoluteEndDesc,
							Required:    true,
						},
					},
				},
			},
		},
	},
}

const (
	schemaDashboardGraphDesc      = "The graph widget."
	schemaDashboardGraph_nameDesc = "The name of the graph."

	schemaDashboardGraph_hostDesc        = "The host graph."
	schemaDashboardGraph_host_hostIDDesc = "The ID of the host."

	schemaDashboardGraph_roleDesc              = "The role graph."
	schemaDashboardGraph_role_roleFullNameDesc = "The service name or role ID."
	schemaDashboardGraph_role_isStackedDesc    = "Whether the graph is stacked or line graph."

	schemaDashboardGraph_serviceDesc             = "The service graph."
	schemaDashboardGraph_service_serviceNameDesc = "The name of the service."

	schemaDashboardGraph_expressionDesc            = "The expression graph."
	schemaDashboardGraph_expression_expressionDesc = "The expression representing the graph."

	schemaDashboardGraph_queryDesc        = "The query graph."
	schemaDashboardGraph_query_queryDesc  = "The PromQL-style query."
	schemaDashboardGraph_query_legendDesc = "The query legend."
)

var schemaDashboardResource_graph = schema.ListNestedBlock{
	NestedObject: schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"title": schemaDashboardResource_widgetTitle,
		},
		Blocks: map[string]schema.Block{
			"range":  schemaDashboardResource_widgetRange,
			"layout": schemaDashboardResource_widgetLayout,

			"host": schema.ListNestedBlock{
				Description: schemaDashboardGraph_hostDesc,
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.ExactlyOneOf(
						path.MatchRelative(),
						path.MatchRelative().AtParent().AtName("role"),
						path.MatchRelative().AtParent().AtName("service"),
						path.MatchRelative().AtParent().AtName("service"),
						path.MatchRelative().AtParent().AtName("expression"),
						path.MatchRelative().AtParent().AtName("query"),
					),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"host_id": schema.StringAttribute{
							Description: schemaDashboardGraph_host_hostIDDesc,
							Required:    true,
						},
						"name": schema.StringAttribute{
							Description: schemaDashboardGraph_nameDesc,
							Required:    true,
						},
					},
				},
			},
			"role": schema.ListNestedBlock{
				Description: schemaDashboardGraph_roleDesc,
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"role_fullname": schema.StringAttribute{
							Description: schemaDashboardGraph_role_roleFullNameDesc,
							Required:    true,
						},
						"name": schema.StringAttribute{
							Description: schemaDashboardGraph_nameDesc,
							Required:    true,
						},
						"is_stacked": schema.BoolAttribute{
							Description: schemaDashboardGraph_role_isStackedDesc,
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
					},
				},
			},
			"service": schema.ListNestedBlock{
				Description: schemaDashboardGraph_serviceDesc,
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"service_name": schema.StringAttribute{
							Description: schemaDashboardGraph_service_serviceNameDesc,
							Required:    true,
						},
						"name": schema.StringAttribute{
							Description: schemaDashboardGraph_nameDesc,
							Required:    true,
						},
					},
				},
			},
			"expression": schema.ListNestedBlock{
				Description: schemaDashboardGraph_expressionDesc,
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"expression": schema.StringAttribute{
							Description: schemaDashboardGraph_expression_expressionDesc,
							Required:    true,
						},
					},
				},
			},
			"query": schema.ListNestedBlock{
				Description: schemaDashboardGraph_queryDesc,
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"query": schema.StringAttribute{
							Description: schemaDashboardGraph_query_queryDesc,
							Required:    true,
						},
						"legend": schema.StringAttribute{
							Description: schemaDashboardGraph_query_legendDesc,
							Required:    true,
						},
					},
				},
			},
		},
	},
}

const (
	schemaDashboardValueDesc              = "The value widget."
	schemaDashboardValue_metricDesc       = "The metric configuration."
	schemaDashboardValue_fractionSizeDesc = "The decimal places displayed on the widget (0-16)."
	schemaDashboardValue_suffixDesc       = "The units to be displayed after the value."
	schemaDashboardValue_nameDesc         = "The name of the metric."

	schemaDashboardValue_hostDesc        = "The host metrics."
	schemaDashboardValue_host_hostIDDesc = "The ID of the host."

	schemaDashboardValue_serviceDesc             = "The service metrics."
	schemaDashboardValue_service_serviceNameDesc = "The name of the service."

	schemaDashboardValue_expressionDesc            = "The expression."
	schemaDashboardValue_expression_expressionDesc = "The expression representing metrics."

	schemaDashboardValue_queryDesc        = "The query."
	schemaDashboardValue_query_queryDesc  = "The PromQL-style query representing metrics."
	schemaDashboardValue_query_legendDesc = "The query legend."
)

var schemaDashboardResource_value = schema.ListNestedBlock{
	Description: schemaDashboardValueDesc,
	NestedObject: schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"title": schemaDashboardResource_widgetTitle,
			"fraction_size": schema.Int64Attribute{
				Description: schemaDashboardValue_fractionSizeDesc,
				Optional:    true,
			},
			"suffix": schema.StringAttribute{
				Description: schemaDashboardValue_suffixDesc,
				Required:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"layout": schemaDashboardResource_widgetLayout,
			"metric": schema.ListNestedBlock{
				Description: schemaDashboardValue_metricDesc,
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"host": schema.ListNestedBlock{
							Description: schemaDashboardValue_hostDesc,
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
								listvalidator.ExactlyOneOf(
									path.MatchRelative(),
									path.MatchRelative().AtParent().AtName("service"),
									path.MatchRelative().AtParent().AtName("expression"),
									path.MatchRelative().AtParent().AtName("query"),
								),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"host_id": schema.StringAttribute{
										Description: schemaDashboardValue_host_hostIDDesc,
										Required:    true,
									},
									"name": schema.StringAttribute{
										Description: schemaDashboardValue_nameDesc,
										Required:    true,
									},
								},
							},
						},
						"service": schema.ListNestedBlock{
							Description: schemaDashboardValue_serviceDesc,
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"service_name": schema.StringAttribute{
										Description: schemaDashboardValue_service_serviceNameDesc,
										Required:    true,
									},
									"name": schema.StringAttribute{
										Description: schemaDashboardValue_nameDesc,
										Required:    true,
									},
								},
							},
						},
						"expression": schema.ListNestedBlock{
							Description: schemaDashboardValue_expressionDesc,
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"expression": schema.StringAttribute{
										Description: schemaDashboardValue_expression_expressionDesc,
										Required:    true,
									},
								},
							},
						},
						"query": schema.ListNestedBlock{
							Description: schemaDashboardValue_queryDesc,
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"query": schema.StringAttribute{
										Description: schemaDashboardValue_query_queryDesc,
										Required:    true,
									},
									"legend": schema.StringAttribute{
										Description: schemaDashboardValue_query_legendDesc,
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

const (
	schemaDashboardMarkdownDesc          = "The markdown widget."
	schemaDashboardMarkdown_markdownDesc = "The markdown formatted string."
)

var schemaDashboardResource_markdown = schema.ListNestedBlock{
	Description: schemaDashboardMarkdownDesc,
	NestedObject: schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"title": schemaDashboardResource_widgetTitle,
			"markdown": schema.StringAttribute{
				Description: schemaDashboardMarkdown_markdownDesc,
				Required:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"layout": schemaDashboardResource_widgetLayout,
		},
	},
}

const (
	schemaDashboardAlertStatusDesc              = "The alert status widget."
	schemaDashboardAlertStatus_roleFullNameDesc = "The service name or role ID."
)

var schemaDashboardResource_alertStatus = schema.ListNestedBlock{
	Description: schemaDashboardAlertStatusDesc,
	NestedObject: schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"title": schemaDashboardResource_widgetTitle,
			"role_fullname": schema.StringAttribute{
				Description: schemaDashboardAlertStatus_roleFullNameDesc,
				Required:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"layout": schemaDashboardResource_widgetLayout,
		},
	},
}
