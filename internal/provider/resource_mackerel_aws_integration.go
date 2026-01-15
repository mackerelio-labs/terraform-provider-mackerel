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
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/planmodifierutil"
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

func (*mackerelAWSIntegrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aws_integration"
}

func (*mackerelAWSIntegrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (*mackerelAWSIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

const (
	schemaAWSIntegrationIDDesc           = "The ID of the AWS integration."
	schemaAWSIntegrationNameDesc         = "The name of the AWS integration."
	schemaAWSIntegrationMemoDesc         = "The notes related to the AWS integration."
	schemaAWSIntegrationKeyDesc          = "The AWS access key ID for the integration."
	schemaAWSIntegrationSecretKeyDesc    = "The AWS secret access key for the integration."
	schemaAWSIntegrationRoleARNDesc      = "The AWS IAM role resource name (ARN) for the integration."
	schemaAWSIntegrationExternalIDDesc   = "The AWS IAM role external ID used during integration."
	schemaAWSIntegrationRegionDesc       = "The AWS region in which the integration will be enabled."
	schemaAWSIntegrationIncludedTagsDesc = "The comma separated list of tags to be included in the integration."
	schemaAWSIntegrationExcludedTagsDesc = "The comma separated list of tags to be excluded in the integration."

	schemaAWSIntegrationServiceDesc                    = "The settings of each AWS service."
	schemaAWSIntegrationServiceEnableDesc              = "Whether integration settings are enabled. Default is `true`."
	schemaAWSIntegrationServiceRoleDesc                = "The ID of the role to be assigned to the service."
	schemaAWSIntegrationServiceExcludedMetricsDesc     = "The list of the metrics to be excluded."
	schemaAWSIntegrationServiceRetireAutomaticallyDesc = "Whether automatic retirement is enabled."
)

var awsIntegrationServices = map[string]struct {
	supportsAutoRetire bool
}{
	"ec2":         {supportsAutoRetire: true},
	"elb":         {},
	"alb":         {supportsAutoRetire: true},
	"nlb":         {supportsAutoRetire: true},
	"rds":         {supportsAutoRetire: true},
	"redshift":    {},
	"elasticache": {supportsAutoRetire: true},
	"sqs":         {},
	"lambda":      {supportsAutoRetire: true},
	"dynamodb":    {},
	"cloudfront":  {},
	"api_gateway": {},
	"kinesis":     {},
	"s3":          {},
	"es":          {},
	"ecs_cluster": {supportsAutoRetire: true},
	"ses":         {},
	"states":      {},
	"efs":         {},
	"firehose":    {},
	"batch":       {},
	"waf":         {},
	"billing":     {},
	"route53":     {},
	"connect":     {},
	"docdb":       {},
	"codebuild":   {},
}

func schemaAWSIntegrationResource() schema.Schema {
	serviceSchema := schema.SetNestedBlock{
		Description: schemaAWSIntegrationServiceDesc,
		Validators: []validator.Set{
			setvalidator.SizeAtMost(1),
		},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"enable": schema.BoolAttribute{
					Description: schemaAWSIntegrationServiceEnableDesc,
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(true),
				},
				"role": schema.StringAttribute{
					Description: schemaAWSIntegrationServiceRoleDesc,
					Optional:    true,
				},
				"excluded_metrics": schema.ListAttribute{
					Description: schemaAWSIntegrationServiceExcludedMetricsDesc,
					ElementType: types.StringType,
					Optional:    true,
					Computed:    true,
					PlanModifiers: []planmodifier.List{
						planmodifierutil.NilRelaxedList(),
					},
				},
			},
		},
	}
	serviceSchemaWithRetireAutomatically := schema.SetNestedBlock{
		Description: schemaAWSIntegrationServiceDesc,
		Validators: []validator.Set{
			setvalidator.SizeAtMost(1),
		},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"enable": schema.BoolAttribute{
					Description: schemaAWSIntegrationServiceEnableDesc,
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(true),
				},
				"role": schema.StringAttribute{
					Description: schemaAWSIntegrationServiceRoleDesc,
					Optional:    true,
				},
				"excluded_metrics": schema.ListAttribute{
					Description: schemaAWSIntegrationServiceExcludedMetricsDesc,
					ElementType: types.StringType,
					Optional:    true,
					Computed:    true,
					PlanModifiers: []planmodifier.List{
						planmodifierutil.NilRelaxedList(),
					},
				},
				"retire_automatically": schema.BoolAttribute{
					Description: schemaAWSIntegrationServiceRetireAutomaticallyDesc,
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
				},
			},
		},
	}

	services := make(map[string]schema.Block, len(awsIntegrationServices))
	for name, spec := range awsIntegrationServices {
		var block schema.Block
		if spec.supportsAutoRetire {
			block = serviceSchemaWithRetireAutomatically
		} else {
			block = serviceSchema
		}
		services[name] = block
	}

	schema := schema.Schema{
		Description: "This resource allows creating and management of the AWS integration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: schemaAWSIntegrationIDDesc,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(), // immutable
				},
			},
			"name": schema.StringAttribute{
				Description: schemaAWSIntegrationNameDesc,
				Required:    true,
			},
			"memo": schema.StringAttribute{
				Description: schemaAWSIntegrationMemoDesc,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"key": schema.StringAttribute{
				Description: schemaAWSIntegrationKeyDesc,
				Optional:    true,
				Sensitive:   true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Validators: []validator.String{
					// With Access Key, secret access key is need too.
					stringvalidator.AlsoRequires(path.MatchRoot("secret_key")),
				},
			},
			"secret_key": schema.StringAttribute{
				Description: schemaAWSIntegrationSecretKeyDesc,
				Optional:    true,
				Sensitive:   true,
				Validators: []validator.String{
					// Secret access key cannot be set alone
					stringvalidator.AlsoRequires(path.MatchRoot("key")),
				},
			},
			"role_arn": schema.StringAttribute{
				Description: schemaAWSIntegrationRoleARNDesc,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"external_id": schema.StringAttribute{
				Description: schemaAWSIntegrationExternalIDDesc,
				Optional:    true,
				Sensitive:   true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Validators: []validator.String{
					// External ID cannot be set alone
					stringvalidator.AlsoRequires(path.MatchRoot("role_arn")),
				},
			},
			"region": schema.StringAttribute{
				Description: schemaAWSIntegrationRegionDesc,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"included_tags": schema.StringAttribute{
				Description: schemaAWSIntegrationIncludedTagsDesc,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"excluded_tags": schema.StringAttribute{
				Description: schemaAWSIntegrationExcludedTagsDesc,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
		},
		Blocks: services,
	}
	return schema
}
