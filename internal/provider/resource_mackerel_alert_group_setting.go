package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
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
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/planmodifierutil"
)

var (
	_ resource.Resource                = (*mackerelAlertGroupSettingResource)(nil)
	_ resource.ResourceWithConfigure   = (*mackerelAlertGroupSettingResource)(nil)
	_ resource.ResourceWithImportState = (*mackerelAlertGroupSettingResource)(nil)
)

func NewMackerelAlertGroupSettingResource() resource.Resource {
	return &mackerelAlertGroupSettingResource{}
}

type mackerelAlertGroupSettingResource struct {
	Client *mackerel.Client
}

func (r *mackerelAlertGroupSettingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_group_setting"
}

func (r *mackerelAlertGroupSettingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemaAlertGroupSettingResource
}

func (r *mackerelAlertGroupSettingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	r.Client = client
}

func (r *mackerelAlertGroupSettingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mackerel.AlertGroupSettingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Create(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to create an alert group setting.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelAlertGroupSettingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mackerel.AlertGroupSettingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Read(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to read an alert group setting.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelAlertGroupSettingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mackerel.AlertGroupSettingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Update(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to update an alert group setting.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelAlertGroupSettingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mackerel.AlertGroupSettingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Delete(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete an alert group setting.",
			err.Error(),
		)
		return
	}
}

func (r *mackerelAlertGroupSettingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

const (
	schemaAlertGroupSettingIDDesc                   = "The ID of the alert group setting."
	schemaAlertGroupSettingNameDesc                 = "The name of the alert group setting."
	schemaAlertGroupSettingMemoDesc                 = "The notes related to the alert group setting."
	schemaAlertGroupSettingServiceScopesDesc        = "The set of the target service names."
	schemaAlertGroupSettingRoleScopesDesc           = "The set of the target role IDs."
	schemaAlertGroupSettingMonitorScopesDesc        = "The set of the target monitor IDs."
	schemaAlertGroupSettingNotificationIntervalDesc = "The time interval (in minutes) for resending notifications."
)

var schemaAlertGroupSettingResource = schema.Schema{
	Description: "This resource allows creating and managemd of alert group settings.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: schemaAlertGroupSettingIDDesc,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Description: schemaAlertGroupSettingNameDesc,
			Required:    true,
		},
		"memo": schema.StringAttribute{
			Description: schemaAlertGroupSettingMemoDesc,
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString(""),
		},
		"service_scopes": schema.SetAttribute{
			Description: schemaAlertGroupSettingServiceScopesDesc,
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			Validators: []validator.Set{
				setvalidator.ValueStringsAre(mackerel.ServiceNameValidator()),
			},
			PlanModifiers: []planmodifier.Set{planmodifierutil.NilRelaxedSet()},
		},
		"role_scopes": schema.SetAttribute{
			Description:   schemaAlertGroupSettingRoleScopesDesc,
			ElementType:   types.StringType,
			Optional:      true,
			Computed:      true,
			PlanModifiers: []planmodifier.Set{planmodifierutil.NilRelaxedSet()},
		},
		"monitor_scopes": schema.SetAttribute{
			Description:   schemaAlertGroupSettingMonitorScopesDesc,
			ElementType:   types.StringType,
			Optional:      true,
			Computed:      true,
			PlanModifiers: []planmodifier.Set{planmodifierutil.NilRelaxedSet()},
		},
		"notification_interval": schema.Int64Attribute{
			Description: schemaAlertGroupSettingNotificationIntervalDesc,
			Optional:    true,
			Computed:    true,
			Default:     int64default.StaticInt64(0),
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
	},
}
