package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/mackerel"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/planmodifierutil"
)

var (
	_ resource.Resource                = (*mackerelNotificationGroupResource)(nil)
	_ resource.ResourceWithConfigure   = (*mackerelNotificationGroupResource)(nil)
	_ resource.ResourceWithImportState = (*mackerelNotificationGroupResource)(nil)
)

func NewMackerelNotificationGroupResource() resource.Resource {
	return &mackerelNotificationGroupResource{}
}

type mackerelNotificationGroupResource struct {
	Client *mackerel.Client
}

func (r *mackerelNotificationGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_group"
}

func (r *mackerelNotificationGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource allows creating and management of notification groups",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the notitication group",

				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the notification group",

				Required: true,
			},
			"notification_level": schema.StringAttribute{
				MarkdownDescription: "The level of notitication (`all` or `critical`)",

				Optional: true,
				Computed: true,
				Validators: []validator.String{
					mackerel.NotificationLevelValidator(),
				},
				Default: stringdefault.StaticString("all"),
			},
			"child_notification_group_ids": schema.SetAttribute{
				Description: "A set of notification group IDs",

				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					planmodifierutil.NilRelaxedSet(),
				},
			},
			"child_channel_ids": schema.SetAttribute{
				Description: "A set of notification channel IDs",

				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					planmodifierutil.NilRelaxedSet(),
				},
			},
		},
		// TODO: migrate to nested attributes
		Blocks: map[string]schema.Block{
			"monitor": schema.SetNestedBlock{
				Description: "Configuration block(s) with monitor rules",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The monitor rule ID",

							Required: true,
						},
						"skip_default": schema.BoolAttribute{
							Description: "If true, send notifications to this notification group only.",

							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
						},
					},
				},
			},
			"service": schema.SetNestedBlock{
				Description: "Configuration block(s) with services",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the service",

							Required: true,
							Validators: []validator.String{
								mackerel.ServiceNameValidator(),
							},
						},
					},
				},
			},
		},
	}
}

func (r *mackerelNotificationGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	r.Client = client
}

func (r *mackerelNotificationGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mackerel.NotificationGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Create(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Notification Group",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelNotificationGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mackerel.NotificationGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Read(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable read Notification Group",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelNotificationGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mackerel.NotificationGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Update(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to update Notification Group",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelNotificationGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mackerel.NotificationGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Delete(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Notification Group",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelNotificationGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
