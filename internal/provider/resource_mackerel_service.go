package provider

import (
	"context"
	"fmt"
	"regexp"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
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
	client *mackerel.Client
}

type mackerelServiceResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Memo types.String `tfsdk:"memo"`
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
					stringvalidator.LengthBetween(2, 63),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-_]+$`),
						"Must include only alphabets, numbers, hyphen and underscore, and it can not begin a hyphen or underscore"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"memo": schema.StringAttribute{
				Optional:    true,
				Description: "Notes related to this service.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
	r.client = client
}

func (r *mackerelServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mackerelServiceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	param := mackerel.CreateServiceParam{
		Name: data.Name.ValueString(),
		Memo: data.Memo.ValueString(),
	}

	service, err := r.client.CreateService(&param)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Service",
			fmt.Sprintf("An unexpected error occurred while attempting to create the service: %+v", err),
		)
		return
	}

	data.ID = types.StringValue(service.Name)

	resp.Diagnostics.Append(r.read(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mackerelServiceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.read(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelServiceResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Unable to update Service",
		"Mackerel services are cannot be updated. Please report this issue.",
	)
}

func (r *mackerelServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mackerelServiceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := data.ID.ValueString()

	if _, err := r.client.DeleteService(id); err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Service",
			fmt.Sprintf("An unexpected error occurred while attempting to delete the service: %+v", err),
		)
		return
	}
}

func (r *mackerelServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *mackerelServiceResource) read(_ context.Context, data *mackerelServiceResourceModel) (diags diag.Diagnostics) {
	id := data.ID.ValueString()

	services, err := r.client.FindServices()
	if err != nil {
		diags.AddError(
			"Unable to read services",
			fmt.Sprintf("An unexpected error occurred while attempting to fetch the services: %+v", err),
		)
		return
	}

	serviceIdx := slices.IndexFunc(services, func(s *mackerel.Service) bool {
		return s.Name == id
	})
	if serviceIdx == -1 {
		diags.AddError(
			"No Service Found",
			fmt.Sprintf("The name '%s' does not match any service in mackerel.io", id),
		)
		return
	}

	service := services[serviceIdx]
	data.Name = types.StringValue(service.Name)
	data.Memo = types.StringValue(service.Memo)
	return
}
