package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/mackerel"
)

var (
	_ resource.Resource                = (*mackerelServiceResource)(nil)
	_ resource.ResourceWithConfigure   = (*mackerelServiceResource)(nil)
	_ resource.ResourceWithImportState = (*mackerelServiceResource)(nil)
)

func NewMackerelServiceResource() resource.Resource {
	return &mackerelServiceResource{}
}

type mackerelServiceResource struct {
	Client *mackerel.Client
}

func (r *mackerelServiceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (r *mackerelServiceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The `mackerel_service` resource allows creating and management of Service.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of service.",
				Validators: []validator.String{
					mackerel.ServiceNameValidator(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"memo": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Notes related to this service.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: stringdefault.StaticString(""),
			},
			"roles": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Set of roles in the service. This is a computed field and will be populated after the service is created.",
			},
		},
	}
}

func (r *mackerelServiceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.Client = client
}

func (r *mackerelServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mackerel.ServiceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Create(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Service",
			err.Error(),
		)
		return
	}

	// Read back the full service data to get the roles field
	if err := data.Read(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Service after creation",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mackerel.ServiceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Read(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Service",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelServiceResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Unable to update Service",
		"Mackerel services cannot be updated in-place. Please report this issue.",
	)
}

func (r *mackerelServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mackerel.ServiceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Delete(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Service",
			err.Error(),
		)
		return
	}
}

func (r *mackerelServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	data, err := mackerel.ImportService(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to import Service",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
