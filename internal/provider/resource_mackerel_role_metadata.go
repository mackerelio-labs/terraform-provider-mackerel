package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/mackerel"
)

var (
	_ resource.Resource                = (*mackerelRoleMetadataResource)(nil)
	_ resource.ResourceWithConfigure   = (*mackerelRoleMetadataResource)(nil)
	_ resource.ResourceWithImportState = (*mackerelRoleMetadataResource)(nil)
)

func NewMackerelRoleMetadataResource() resource.Resource {
	return &mackerelRoleMetadataResource{}
}

type mackerelRoleMetadataResource struct {
	Client *mackerel.Client
}

func (r *mackerelRoleMetadataResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role_metadata"
}

func (r *mackerelRoleMetadataResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource creates and manages a Role Metadata.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(), // immutable
				},
			},
			"service": schema.StringAttribute{
				Description: "The name of the service.",
				Required:    true,
				Validators: []validator.String{
					mackerel.ServiceNameValidator(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(), // force new
				},
			},
			"role": schema.StringAttribute{
				Description: "The name of the role.",
				Required:    true,
				Validators: []validator.String{
					mackerel.RoleNameValidator(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(), // force new
				},
			},
			"namespace": schema.StringAttribute{
				Description: "The identifier for the metadata.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(), // force new
				},
			},
			"metadata_json": schema.StringAttribute{
				Description: "The arbitrary JSON data for the role.",
				Required:    true,
				CustomType:  jsontypes.NormalizedType{},
			},
		},
	}
}

func (r *mackerelRoleMetadataResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	r.Client = client
}

func (r *mackerelRoleMetadataResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mackerel.RoleMetadataModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Create(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Role Metadata",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelRoleMetadataResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mackerel.RoleMetadataModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Read(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Role Metadata",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelRoleMetadataResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mackerel.RoleMetadataModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Update(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to update Role Metadata",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelRoleMetadataResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mackerel.RoleMetadataModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Delete(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Role Metadata",
			err.Error(),
		)
		return
	}
}

func (r *mackerelRoleMetadataResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	data, err := mackerel.ImportRoleMetadata(req.ID)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("id"),
			"Invalid ID",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
