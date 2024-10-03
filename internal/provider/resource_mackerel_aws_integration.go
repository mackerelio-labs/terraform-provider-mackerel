package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
)

var (
	_ resource.Resource                = (*mackerelAWSIntegrationResource)(nil)
	_ resource.ResourceWithConfigure   = (*mackerelAWSIntegrationResource)(nil)
	_ resource.ResourceWithImportState = (*mackerelAWSIntegrationResource)(nil)
)

func NewMackerelAWSIntegrationResource() resource.Resource {
	return &mackerelAWSIntegrationResource{}
}

type mackerelAWSIntegrationResource struct {
	Client *mackerel.Client
}

func (_ *mackerelAWSIntegrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aws_integration"
}

func (_ *mackerelAWSIntegrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemaAWSIntegrationResource()
}

func (r *mackerelAWSIntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	r.Client = client
}

func (r *mackerelAWSIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mackerel.AWSIntegrationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Create(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to create aws integration settings",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelAWSIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mackerel.AWSIntegrationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Read(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to read aws integration settings",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelAWSIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mackerel.AWSIntegrationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Update(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to update aws integration settings",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelAWSIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mackerel.AWSIntegrationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Delete(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete aws integration settings",
			err.Error(),
		)
		return
	}
}

func (_ *mackerelAWSIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func schemaAWSIntegrationResource() schema.Schema {
	serviceSchema := schema.SetNestedBlock{
		Validators: []validator.Set{
			setvalidator.SizeAtMost(1),
		},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"enable": schema.BoolAttribute{
					Optional: true,
					Computed: true,
					Default:  booldefault.StaticBool(true),
				},
				"role": schema.StringAttribute{
					Optional: true,
				},
				"excluded_metrics": schema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
				},
			},
		},
	}
	serviceSchemaWithRetireAutomatically := schema.SetNestedBlock{
		Validators: []validator.Set{
			setvalidator.SizeAtMost(1),
		},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"enable": schema.BoolAttribute{
					Optional: true,
					Computed: true,
					Default:  booldefault.StaticBool(true),
				},
				"role": schema.StringAttribute{
					Optional: true,
				},
				"excluded_metrics": schema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
				},
				"retire_automatically": schema.BoolAttribute{
					Optional: true,
					Computed: true,
					Default:  booldefault.StaticBool(false),
				},
			},
		},
	}

	services := map[string]schema.Block{
		"ec2":         serviceSchemaWithRetireAutomatically,
		"elb":         serviceSchema,
		"alb":         serviceSchema,
		"nlb":         serviceSchema,
		"rds":         serviceSchemaWithRetireAutomatically,
		"redshift":    serviceSchema,
		"elasticache": serviceSchemaWithRetireAutomatically,
		"sqs":         serviceSchema,
		"lambda":      serviceSchema,
		"dynamodb":    serviceSchema,
		"cloudfront":  serviceSchema,
		"api_gateway": serviceSchema,
		"kinesis":     serviceSchema,
		"s3":          serviceSchema,
		"es":          serviceSchema,
		"ecs_cluster": serviceSchema,
		"ses":         serviceSchema,
		"states":      serviceSchema,
		"efs":         serviceSchema,
		"firehose":    serviceSchema,
		"batch":       serviceSchema,
		"waf":         serviceSchema,
		"billing":     serviceSchema,
		"route53":     serviceSchema,
		"connect":     serviceSchema,
		"docdb":       serviceSchema,
		"codebuild":   serviceSchema,
	}

	schema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(), // immutable
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"memo": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Computed:  true,
				Default:   stringdefault.StaticString(""),
				Validators: []validator.String{
					// With Access Key, secret access key is need too.
					stringvalidator.AlsoRequires(path.MatchRoot("secret_key")),
				},
			},
			"secret_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Computed:  true,
				Default:   stringdefault.StaticString(""),
				Validators: []validator.String{
					// Secret access key cannot be set alone
					stringvalidator.AlsoRequires(path.MatchRoot("key")),
				},
			},
			"role_arn": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"external_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					// External ID cannot be set alone
					stringvalidator.AlsoRequires(path.MatchRoot("role_arn")),
				},
			},
			"region": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"included_tags": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"excluded_tags": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
		},
		Blocks: services,
	}
	return schema
}
