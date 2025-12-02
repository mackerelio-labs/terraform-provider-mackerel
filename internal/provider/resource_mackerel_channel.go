package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/mackerel"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/planmodifierutil"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/validatorutil"
)

var (
	_ resource.Resource                     = (*mackerelChannelResource)(nil)
	_ resource.ResourceWithConfigValidators = (*mackerelChannelResource)(nil)
	_ resource.ResourceWithConfigure        = (*mackerelChannelResource)(nil)
	_ resource.ResourceWithImportState      = (*mackerelChannelResource)(nil)
)

func NewMackerelChannelResource() resource.Resource {
	return &mackerelChannelResource{}
}

type mackerelChannelResource struct {
	Client *mackerel.Client
}

func (r *mackerelChannelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_channel"
}

func (r *mackerelChannelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema, _ = schemaChannelResource()
}

func (r *mackerelChannelResource) ConfigValidators(context.Context) []resource.ConfigValidator {
	_, validators := schemaChannelResource()
	return validators
}

func (r *mackerelChannelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	r.Client = client
}

func (r *mackerelChannelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mackerel.ChannelModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Create(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to create a channel",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelChannelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mackerel.ChannelModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Read(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to read a channel",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelChannelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mackerel.ChannelModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := data.Update(ctx, r.Client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update a channel",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelChannelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mackerel.ChannelModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Delete(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete a channel",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelChannelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

const (
	schemaChannelIDDesc     = "The ID of the notification channel."
	schemaChannelNameDesc   = "The name of the notification channel."
	schemaChannelEventsDesc = "The set of notification event types to be received."

	schemaChannelEmailDesc         = "The settings for an email notification channel."
	schemaChannelEmail_EmailsDesc  = "The set of email addresses specified to receive notifications."
	schemaChannelEmail_UserIDsDesc = "The set of user IDs specified to receive notifications."

	schemaChannelSlackDesc          = "The settings for a slack notification channel."
	schemaChannelSlack_URLDesc      = "The incoming webhook URL for Slack."
	schemaChannelSlack_MentionsDesc = "The map of the condition (ok, warning, critical)" +
		"and the text accompanying the alert notification."
	schemaChannelSlack_EnabledGraphImageDesc = "Whether or not the corresponding graph is posted to Slack."

	schemaChannelWebhookDesc     = "The settings for a webhook notification channel."
	schemaChannelWebhook_URLDesc = "The URL that will receive HTTP request."
)

// requiresReplaceIfChannelTypeChanges returns a plan modifier that requires replace
// when the channel type changes
func requiresReplaceIfChannelTypeChanges() planmodifier.List {
	return listplanmodifier.RequiresReplaceIf(
		func(ctx context.Context, req planmodifier.ListRequest, res *listplanmodifier.RequiresReplaceIfFuncResponse) {
			// Only require replace if list size changes (0->1 or 1->0)
			// This indicates a channel type change
			oldSize := len(req.StateValue.Elements())
			newSize := len(req.PlanValue.Elements())
			if (oldSize == 0 && newSize > 0) || (oldSize > 0 && newSize == 0) {
				res.RequiresReplace = true
			}
		},
		"Channel type cannot be changed in-place",
		"Channel type cannot be changed in-place",
	)
}

func schemaChannelResource() (schema.Schema, []resource.ConfigValidator) {
	eventsAttr := schema.SetAttribute{
		Description: schemaChannelEventsDesc,
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		Validators: []validator.Set{setvalidator.ValueStringsAre(
			stringvalidator.OneOf(
				"alert",
				"alertGroup",
				"hostStatus",
				"hostRegister",
				"hostRetire",
				"monitor",
			),
		)},
		PlanModifiers: []planmodifier.Set{
			planmodifierutil.NilRelaxedSet(),
		},
	}
	schema := schema.Schema{
		Description: "This resource allows creating and management of channel, which manages either email, slack or webhook.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: schemaChannelIDDesc,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(), // immutable
				},
			},
			"name": schema.StringAttribute{
				Description: schemaChannelNameDesc,
				Required:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"email": schema.ListNestedBlock{
				Description: schemaChannelEmailDesc,
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				PlanModifiers: []planmodifier.List{
					requiresReplaceIfChannelTypeChanges(),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"emails": schema.SetAttribute{
							ElementType: types.StringType,
							Description: schemaChannelEmail_EmailsDesc,
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.Set{
								planmodifierutil.NilRelaxedSet(),
							},
						},
						"user_ids": schema.SetAttribute{
							ElementType: types.StringType,
							Description: schemaChannelEmail_UserIDsDesc,
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.Set{
								planmodifierutil.NilRelaxedSet(),
							},
						},
						"events": eventsAttr,
					},
				},
			},
			"slack": schema.ListNestedBlock{
				Description: schemaChannelSlackDesc,
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				PlanModifiers: []planmodifier.List{
					requiresReplaceIfChannelTypeChanges(),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"url": schema.StringAttribute{
							Description: schemaChannelSlack_URLDesc,
							Required:    true,
							Validators: []validator.String{
								validatorutil.IsURLWithHTTPorHTTPS(),
							},
						},
						// FIXME(needs schema upgrade): use nested attribute
						"mentions": schema.MapAttribute{
							ElementType: types.StringType,
							Description: schemaChannelSlack_MentionsDesc,
							Optional:    true,
							Computed:    true,
							Validators: []validator.Map{
								mapvalidator.KeysAre(stringvalidator.OneOf("ok", "warning", "critical")),
								mapvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
							},
							PlanModifiers: []planmodifier.Map{
								planmodifierutil.NilRelaxedMap(),
							},
						},
						"enabled_graph_image": schema.BoolAttribute{
							Description: schemaChannelSlack_EnabledGraphImageDesc,
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"events": eventsAttr,
					},
				},
			},
			"webhook": schema.ListNestedBlock{
				Description: schemaChannelWebhookDesc,
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				PlanModifiers: []planmodifier.List{
					requiresReplaceIfChannelTypeChanges(),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"url": schema.StringAttribute{
							Description: schemaChannelWebhook_URLDesc,
							Required:    true,
							Validators: []validator.String{
								validatorutil.IsURLWithHTTPorHTTPS(),
							},
						},
						"events": eventsAttr,
					},
				},
			},
		},
	}
	validators := []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("email"),
			path.MatchRoot("slack"),
			path.MatchRoot("webhook"),
		),
	}
	return schema, validators
}
