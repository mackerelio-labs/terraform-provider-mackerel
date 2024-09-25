package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/mackerel"
)

var (
	_ resource.Resource                = (*mackerelDowntimeResource)(nil)
	_ resource.ResourceWithConfigure   = (*mackerelDowntimeResource)(nil)
	_ resource.ResourceWithImportState = (*mackerelDowntimeResource)(nil)
)

func NewMackerelDowntimeResource() resource.Resource {
	return &mackerelDowntimeResource{}
}

type mackerelDowntimeResource struct {
	Client *mackerel.Client
}

func (r *mackerelDowntimeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_downtime"
}

func (r *mackerelDowntimeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemaDowntimeResource()
}

func (r *mackerelDowntimeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	r.Client = client
}

func (r *mackerelDowntimeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mackerel.DowntimeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Create(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to create a downtime",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelDowntimeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mackerel.DowntimeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Read(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to read a downtime",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelDowntimeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mackerel.DowntimeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Update(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to update a downtime",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelDowntimeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mackerel.DowntimeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Delete(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete a downtime",
			err.Error(),
		)
		return
	}
}

func (r *mackerelDowntimeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

const (
	schemaDowntimeIDDesc              = "The id of the downtime."
	schemaDowntimeNameDesc            = "The name of the downtime."
	schemaDowntimeMemoDesc            = "The notes for the downtime."
	schemaDowntimeStartDesc           = "The starting time (in epoch seconds) of the downtime."
	schemaDowntimeDurationDesc        = "The duration (in minutes) of the downtime."
	schemaDowntimeRecurrenceDesc      = "The configuration for repeating occurrences."
	schemaDowntimeRecurrence_typeDesc = "The recurrence options." +
		"Valid options are `hourly`, `daily`, `weekly`, `monthly` or `yearly`."
	schemaDowntimeRecurrence_intervalDesc = "The recurrence interval."
	schemaDowntimeRecurrence_weekdaysDesc = "The set of the day of the week." +
		"Valid values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` or `Saturday`." +
		"Only available when the `type` is set to `weekly`."
	schemaDowntimeRecurrence_untilDesc     = "The time at which recurrence ends (in epoch seconds)."
	schemaDowntimeServiceScopesDesc        = "The set of target service names."
	schemaDowntimeServiceExcludeScopesDesc = "The set of excluded service names."
	schemaDowntimeRoleScopesDesc           = "The set of target role IDs."
	schemaDowntimeRoleExcludeScopesDesc    = "The set of excluded role IDs."
	schemaDowntimeMonitorScopesDesc        = "The set of target monitor IDs."
	schemaDowntimeMonitorExcludeScopesDesc = "The set of excluded monitor IDs."
)

func schemaDowntimeResource() schema.Schema {
	s := schema.Schema{
		Description: "This resource allows creating and management of downtimes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: schemaDowntimeIDDesc,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(), // immutable
				},
			},
			"name": schema.StringAttribute{
				Description: schemaDowntimeNameDesc,
				Required:    true,
			},
			"memo": schema.StringAttribute{
				Description: schemaDowntimeMemoDesc,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"start": schema.Int64Attribute{
				Description: schemaDowntimeStartDesc,
				Required:    true,
			},
			"duration": schema.Int64Attribute{
				Description: schemaDowntimeDurationDesc,
				Required:    true,
			},
			"service_scopes": schema.SetAttribute{
				Description: schemaDowntimeServiceScopesDesc,
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(mackerel.ServiceNameValidator()),
				},
			},
			"service_exclude_scopes": schema.SetAttribute{
				Description: schemaDowntimeServiceExcludeScopesDesc,
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(mackerel.ServiceNameValidator()),
				},
			},
			"role_scopes": schema.SetAttribute{
				Description: schemaDowntimeRoleScopesDesc,
				ElementType: types.StringType,
				Optional:    true,
			},
			"role_exclude_scopes": schema.SetAttribute{
				Description: schemaDowntimeRoleExcludeScopesDesc,
				ElementType: types.StringType,
				Optional:    true,
			},
			"monitor_scopes": schema.SetAttribute{
				Description: schemaDowntimeMonitorScopesDesc,
				ElementType: types.StringType,
				Optional:    true,
			},
			"monitor_exclude_scopes": schema.SetAttribute{
				Description: schemaDowntimeMonitorExcludeScopesDesc,
				ElementType: types.StringType,
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"recurrence": schema.ListNestedBlock{
				Description: schemaDowntimeRecurrenceDesc,
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: schemaDowntimeRecurrence_typeDesc,
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("hourly", "daily", "weekly", "monthly", "yearly"),
							},
						},
						"interval": schema.Int64Attribute{
							Description: schemaDowntimeRecurrence_intervalDesc,
							Required:    true,
						},
						"weekdays": schema.SetAttribute{
							Description:         schemaDowntimeRecurrence_weekdaysDesc,
							MarkdownDescription: schemaDowntimeRecurrence_weekdaysDesc,
							ElementType:         types.StringType,
							Optional:            true,
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.OneOf(
										"Sunday", "Monday", "Tuesday", "Wednesday",
										"Thursday", "Friday", "Saturday",
									),
								),
							},
						},
						"until": schema.Int64Attribute{
							Description: schemaDowntimeRecurrence_untilDesc,
							Optional:    true,
							Computed:    true,
							Default:     int64default.StaticInt64(0),
						},
					},
				},
			},
		},
	}
	return s
}
