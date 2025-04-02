package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/mackerel"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/planmodifierutil"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/typeutil"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/validatorutil"
)

var (
	_ resource.Resource                     = (*mackerelMonitorResource)(nil)
	_ resource.ResourceWithConfigure        = (*mackerelMonitorResource)(nil)
	_ resource.ResourceWithConfigValidators = (*mackerelMonitorResource)(nil)
	_ resource.ResourceWithImportState      = (*mackerelMonitorResource)(nil)
)

func NewMackerelMonitorResource() resource.Resource {
	return &mackerelMonitorResource{}
}

type mackerelMonitorResource struct {
	Client *mackerel.Client
}

func (r *mackerelMonitorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (r *mackerelMonitorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema, _ = schemaMonitorResource()
}

func (r *mackerelMonitorResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	_, validators := schemaMonitorResource()
	return validators
}

func (r *mackerelMonitorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	r.Client = client
}

func (r *mackerelMonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mackerel.MonitorModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Create(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Monitor",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelMonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mackerel.MonitorModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Read(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Monitor",
			err.Error(),
		)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelMonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mackerel.MonitorModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Update(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to update Monitor",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mackerelMonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mackerel.MonitorModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Delete(ctx, r.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Monitor",
			err.Error(),
		)
		return
	}
}

func (r *mackerelMonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema

const (
	schemaMonitorIDDesc                   = "The ID of the monitor."
	schemaMonitorNameDesc                 = "The name of the monitor."
	schemaMonitorMemoDesc                 = "The notes for the monitoring configuration."
	schemaMonitorIsMuteDesc               = "Whether monitoring is muted or not."
	schemaMonitorNotificationIntervalDesc = "The time interval (in minutes) for re-sending notifications." +
		"If this field is empty, notifications will not be re-sent."
)

func schemaMonitorResource() (schema.Schema, []resource.ConfigValidator) {
	return schema.Schema{
			Description: "This resource allows creating and management of monitors.",
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Description: schemaMonitorIDDesc,
					Computed:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"name": schema.StringAttribute{
					Description: schemaMonitorNameDesc,
					Required:    true,
				},
				"memo": schema.StringAttribute{
					Description: schemaMonitorMemoDesc,
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
				},
				"is_mute": schema.BoolAttribute{
					Description: schemaMonitorIsMuteDesc,
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
				},
				"notification_interval": schema.Int64Attribute{
					Description: schemaMonitorNotificationIntervalDesc,
					Optional:    true,
					Computed:    true,
					Default:     int64default.StaticInt64(0), // TODO(schema upgrade): handle null
					Validators: []validator.Int64{
						int64validator.AtLeast(0),
					},
				},
			},
			Blocks: map[string]schema.Block{
				"host_metric":       schemaMonitorResourceHostMetricBlock(),
				"service_metric":    schemaMonitorResourceServiceMetricBlock(),
				"expression":        schemaMonitorResourceExpressionBlock(),
				"query":             schemaMonitorResourceQueryBlock(),
				"connectivity":      schemaMonitorResourceConnectivityBlock(),
				"external":          schemaMonitorResourceExternalBlock(),
				"anomaly_detection": schemaMonitorResourceAnomalyDetectionBlock(),
			},
		}, []resource.ConfigValidator{
			resourcevalidator.ExactlyOneOf(
				path.MatchRoot("host_metric"),
				path.MatchRoot("service_metric"),
				path.MatchRoot("expression"),
				path.MatchRoot("query"),
				path.MatchRoot("connectivity"),
				path.MatchRoot("external"),
				path.MatchRoot("anomaly_detection"),
			),
		}
}

const (
	schemaMonitorOperatorDesc = "The operator that determines the conditions that state whether the designated variable is greater (>) or less than (<)." +
		`The observed value is on the left of ">" or "<" and the designated value is on the right`
	schemaMonitorWarningDesc          = "The threshold that generates a warning alert."
	schemaMonitorCriticalDesc         = "The threshold that generates a critical alert."
	schemaMonitorDurationDesc         = "The average value of the designated interval (in minutes) will be monitored."
	schemaMonitorMaxCheckAttepmtsDesc = "The number of consecutive warning/critical instances before an alert is made."
	schemaMonitorScopesDesc           = "The set of service names and/or role ids to be monitored."
	schemaMonitorExcludeScopesDesc    = "The set of service names and/or role ids to except from monitoring."
)

func schemaMonitorResourceOperatorAttr() schema.StringAttribute {
	return schema.StringAttribute{
		Description: schemaMonitorOperatorDesc,
		Required:    true,
		Validators: []validator.String{
			stringvalidator.OneOf(">", "<"),
		},
	}
}

// TODO(schema upgrade): handle null and use float64 directly
func schemaMonitorResourceWarningAttr() schema.StringAttribute {
	return schema.StringAttribute{
		Description: schemaMonitorWarningDesc,
		Optional:    true,
		Computed:    true,
		CustomType:  typeutil.FloatStringType{},
		Default:     stringdefault.StaticString(""),
	}
}

func schemaMonitorResourceCriticalAttr() schema.StringAttribute {
	return schema.StringAttribute{
		Description: schemaMonitorCriticalDesc,
		Optional:    true,
		Computed:    true,
		CustomType:  typeutil.FloatStringType{},
		Default:     stringdefault.StaticString(""),
	}
}

func schemaMonitorResourceThresholdValidator() validator.Object {
	return objectvalidator.AtLeastOneOf(path.MatchRoot("warning"), path.MatchRoot("critical"))
}

func schemaMonitorResourceDurationAttr() schema.Int64Attribute {
	return schema.Int64Attribute{
		Description: schemaMonitorDurationDesc,
		Required:    true,
		Validators: []validator.Int64{
			int64validator.Between(1, 10),
		},
	}
}

func schemaMonitorMaxCheckAttemptsAttr(defaultAttempts uint) schema.Int64Attribute {
	return schema.Int64Attribute{
		Description: schemaMonitorMaxCheckAttepmtsDesc,
		Optional:    true,
		Computed:    true,
		Default:     int64default.StaticInt64(int64(defaultAttempts)),
		Validators: []validator.Int64{
			int64validator.Between(1, 10),
		},
	}
}

func schemaMonitorResourceScopesAttr() schema.SetAttribute {
	return schema.SetAttribute{
		ElementType:   types.StringType,
		Description:   schemaMonitorScopesDesc,
		Optional:      true,
		Computed:      true,
		PlanModifiers: []planmodifier.Set{planmodifierutil.NilRelaxedSet()},
	}
}

func schemaMonitorResourceExcludeScopesAttr() schema.SetAttribute {
	return schema.SetAttribute{
		ElementType:   types.StringType,
		Description:   schemaMonitorExcludeScopesDesc,
		Optional:      true,
		Computed:      true,
		PlanModifiers: []planmodifier.Set{planmodifierutil.NilRelaxedSet()},
	}
}

const (
	schemaMonitorHostMetricDesc        = "The settings for the host metric monitoring."
	schemaMonitorHostMetric_MetricDesc = "The name of the host metric to be monitored."
)

func schemaMonitorResourceHostMetricBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: schemaMonitorHostMetricDesc,
		Validators:  []validator.List{listvalidator.SizeAtMost(1)},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"metric": schema.StringAttribute{
					Description: schemaMonitorHostMetric_MetricDesc,
					Required:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"operator":           schemaMonitorResourceOperatorAttr(),
				"warning":            schemaMonitorResourceWarningAttr(),
				"critical":           schemaMonitorResourceCriticalAttr(),
				"duration":           schemaMonitorResourceDurationAttr(),
				"max_check_attempts": schemaMonitorMaxCheckAttemptsAttr(1),
				"scopes":             schemaMonitorResourceScopesAttr(),
				"exclude_scopes":     schemaMonitorResourceExcludeScopesAttr(),
			},
			Validators: []validator.Object{
				schemaMonitorResourceThresholdValidator(),
			},
		},
	}
}

const (
	schemaMonitorServiceMetricDesc                         = "The settings for monitoring service metrics."
	schemaMonitorServiceMetric_ServiceDesc                 = "The name of the service to be monitored."
	schemaMonitorServiceMetric_MetricDesc                  = "The name of the service metric to be monitored."
	schemaMonitorServiceMetric_MissingDurationWarningDesc  = "The threshold (in minutes) to generate a warning alert for the interruption monitoring."
	schemaMonitorServiceMetric_MissingDurationCriticalDesc = "The threshold (in minutes) to generate a critical alert for the interruption monitoring."
)

func schemaMonitorResourceServiceMetricBlock() schema.Block {
	missingDurationValidator := int64validator.All(
		int64validator.Between(10, 7*24*60),
		validatorutil.IntDivisibleBy(10),
	)
	return schema.ListNestedBlock{
		Description: schemaMonitorServiceMetricDesc,
		Validators:  []validator.List{listvalidator.SizeAtMost(1)},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"service": schema.StringAttribute{
					Description: schemaMonitorServiceMetric_ServiceDesc,
					Required:    true,
					Validators: []validator.String{
						mackerel.ServiceNameValidator(),
					},
				},
				"metric": schema.StringAttribute{
					Description: schemaMonitorServiceMetric_MetricDesc,
					Required:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"operator":           schemaMonitorResourceOperatorAttr(),
				"warning":            schemaMonitorResourceWarningAttr(),
				"critical":           schemaMonitorResourceCriticalAttr(),
				"duration":           schemaMonitorResourceDurationAttr(),
				"max_check_attempts": schemaMonitorMaxCheckAttemptsAttr(1),
				"missing_duration_warning": schema.Int64Attribute{
					Description: schemaMonitorServiceMetric_MissingDurationWarningDesc,
					Optional:    true,
					Computed:    true,
					Default:     int64default.StaticInt64(0), // TODO(schema upgrade): handle null
					Validators:  []validator.Int64{missingDurationValidator},
				},
				"missing_duration_critical": schema.Int64Attribute{
					Description: schemaMonitorServiceMetric_MissingDurationCriticalDesc,
					Optional:    true,
					Computed:    true,
					Default:     int64default.StaticInt64(0), // TODO(schema upgrade): handle null
					Validators:  []validator.Int64{missingDurationValidator},
				},
			},
			Validators: []validator.Object{
				schemaMonitorResourceThresholdValidator(),
			},
		},
	}
}

const (
	schemaMonitorExpressionDesc            = "The settings for the expression monitoring."
	schemaMonitorExpression_ExpressionDesc = "The expression of the monitoring target. Only valid for graph sequences that become one line."
)

func schemaMonitorResourceExpressionBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: schemaMonitorExpressionDesc,
		Validators:  []validator.List{listvalidator.SizeAtMost(1)},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"expression": schema.StringAttribute{
					Description: schemaMonitorExpression_ExpressionDesc,
					Required:    true,
				},
				"operator": schemaMonitorResourceOperatorAttr(),
				"warning":  schemaMonitorResourceWarningAttr(),
				"critical": schemaMonitorResourceCriticalAttr(),
			},
			Validators: []validator.Object{
				schemaMonitorResourceThresholdValidator(),
			},
		},
	}
}

const (
	schemaMonitorQueryDesc        = "The settings for the query monitoring."
	schemaMonitorQuery_QueryDesc  = "The PromQL-style query of the monitoring target(s)."
	schemaMonitorQuery_LegendDesc = "The graph legend for the alerts."
)

func schemaMonitorResourceQueryBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: schemaMonitorQueryDesc,
		Validators:  []validator.List{listvalidator.SizeAtMost(1)},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"query": schema.StringAttribute{
					Description: schemaMonitorQuery_QueryDesc,
					Required:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(), // force new
					},
				},
				"legend": schema.StringAttribute{
					Description: schemaMonitorQuery_LegendDesc,
					Required:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(), // force new
					},
				},
				"operator": schemaMonitorResourceOperatorAttr(),
				"warning":  schemaMonitorResourceWarningAttr(),
				"critical": schemaMonitorResourceCriticalAttr(),
			},
			Validators: []validator.Object{schemaMonitorResourceThresholdValidator()},
		},
	}
}

const (
	schemaMonitorConnectivityDesc                   = "The settings for the host connectivity monitoring"
	schemaMonitorConnectivity_AlertStatusOnGoneDesc = "The status of an alert generated by this monitor."
)

func schemaMonitorResourceConnectivityBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: schemaMonitorConnectivityDesc,
		Validators:  []validator.List{listvalidator.SizeAtMost(1)},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"alert_status_on_gone": schema.StringAttribute{
					Description: schemaMonitorConnectivity_AlertStatusOnGoneDesc,
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("CRITICAL"),
					Validators: []validator.String{
						stringvalidator.OneOf("CRITICAL", "WARNING"),
					},
				},
				"scopes":         schemaMonitorResourceScopesAttr(),
				"exclude_scopes": schemaMonitorResourceExcludeScopesAttr(),
			},
		},
	}
}

const (
	schemaMonitorExternalDesc = "The settings for the external HTTP monitoring."

	schemaMonitorExternal_URLDesc         = "The URL to be monitored."
	schemaMonitorExternal_MethodDesc      = "The request method."
	schemaMonitorExternal_RequestBodyDesc = "The request body."
	schemaMonitorExternal_HeadersDesc     = "The request headers."

	schemaMonitorExternal_ServiceDesc = "The service name. " +
		"When response time is monitored, it will be graphed as a service metric of this."
	schemaMonitorExternal_ResponseTimeCriticalDesc = "The threshold (in milliseconds) of the response time for critical alerts."
	schemaMonitorExternal_ResponseTimeWarningDesc  = "The threshold (in milliseconds) of the response time for warning alerts."
	schemaMonitorExternal_ResponseTimeDurationDesc = "The duration (in minutes) for calcurating the average response time to be monitored."

	schemaMonitorExternal_ContainsStringDesc = "The string which should be contained in the response body."
	schemaMonitorExternal_FollowRedirectDesc = "Whether or not to track the redirection destination of the response."

	schemaMonitorExternal_SkipCertificateVerificationDesc     = "Whether or not to skip the verification of certificate."
	schemaMonitorExternal_CertificationExpirationCriticalDesc = "The threshold (in days) of the certification expiration date for critical alerts."
	schemaMonitorExternal_CertificationExpirationWarningDesc  = "The threshold (in days) of the certification expiration date for warning alerts."
	schemaMonitorExternal_ExpectedStatusCodeDesc              = "Specify the status code that is judged as OK. If not specified, 2xx or 3xx will be judged as OK."
)

func schemaMonitorResourceExternalBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: schemaMonitorExternalDesc,
		Validators:  []validator.List{listvalidator.SizeAtMost(1)},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"max_check_attempts": schemaMonitorMaxCheckAttemptsAttr(1),
				"url": schema.StringAttribute{
					Description: schemaMonitorExternal_URLDesc,
					Required:    true,
					Validators: []validator.String{
						validatorutil.IsURLWithHTTPorHTTPS(),
					},
				},
				"method": schema.StringAttribute{
					Description: schemaMonitorExternal_MethodDesc,
					Required:    true,
					Validators: []validator.String{
						stringvalidator.OneOf("GET", "POST", "PUT", "DELETE"),
					},
				},
				"request_body": schema.StringAttribute{
					Description: schemaMonitorExternal_RequestBodyDesc,
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
				},
				"headers": schema.MapAttribute{
					ElementType:   types.StringType,
					Description:   schemaMonitorExternal_HeadersDesc,
					Sensitive:     true,
					Optional:      true,
					Computed:      true,
					PlanModifiers: []planmodifier.Map{planmodifierutil.NilRelaxedMap()},
				},

				"service": schema.StringAttribute{
					Description: schemaMonitorExternal_ServiceDesc,
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					Validators: []validator.String{
						mackerel.ServiceNameValidator(),
					},
				},
				"response_time_critical": schema.Float64Attribute{
					Description: schemaMonitorExternal_ResponseTimeCriticalDesc,
					Optional:    true,
					Computed:    true,
					Default:     float64default.StaticFloat64(0.0),
					Validators: []validator.Float64{
						float64validator.AlsoRequires(path.MatchRelative().AtParent().AtName("service")),
					},
				},
				"response_time_warning": schema.Float64Attribute{
					Description: schemaMonitorExternal_ResponseTimeWarningDesc,
					Optional:    true,
					Computed:    true,
					Default:     float64default.StaticFloat64(0.0),
					Validators: []validator.Float64{
						float64validator.AlsoRequires(path.MatchRelative().AtParent().AtName("service")),
					},
				},
				"response_time_duration": schema.Int64Attribute{
					Description: schemaMonitorExternal_ResponseTimeDurationDesc,
					Optional:    true,
					Computed:    true,
					Default:     int64default.StaticInt64(0),
					Validators: []validator.Int64{
						int64validator.AlsoRequires(path.MatchRelative().AtParent().AtName("service")),
					},
				},

				"contains_string": schema.StringAttribute{
					Description: schemaMonitorExternal_ContainsStringDesc,
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
				},
				"follow_redirect": schema.BoolAttribute{
					Description: schemaMonitorExternal_FollowRedirectDesc,
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
				},

				"skip_certificate_verification": schema.BoolAttribute{
					Description: schemaMonitorExternal_SkipCertificateVerificationDesc,
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
				},
				"certification_expiration_critical": schema.Int64Attribute{
					Description: schemaMonitorExternal_CertificationExpirationCriticalDesc,
					Optional:    true,
					Computed:    true,
					Default:     int64default.StaticInt64(0),
				},
				"certification_expiration_warning": schema.Int64Attribute{
					Description: schemaMonitorExternal_CertificationExpirationWarningDesc,
					Optional:    true,
					Computed:    true,
					Default:     int64default.StaticInt64(0),
				},
				"expected_status_code": schema.Int64Attribute{
					Description: schemaMonitorExternal_ExpectedStatusCodeDesc,
					Optional:    true,
				},
			},
		},
	}
}

const (
	schemaMonitorAnomalyDetectionDesc                     = "The settings for monitoring with anomaly detection."
	schemaMonitorAnomalyDetection_WarningSensitivityDesc  = "The sensitivity that generates warning alerts."
	schemaMonitorAnomalyDetection_CriticalSensitivityDesc = "The sensitivity that generates critical alerts."
	schemaMonitorAnomalyDetection_TrainingPeriodFromDesc  = "The specified training period."
)

func schemaMonitorResourceAnomalyDetectionBlock() schema.Block {
	sensitivityValidator := stringvalidator.OneOf("insensitive", "normal", "sensitive")
	return schema.ListNestedBlock{
		Description: schemaMonitorAnomalyDetectionDesc,
		Validators:  []validator.List{listvalidator.SizeAtMost(1)},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"warning_sensitivity": schema.StringAttribute{
					Description: schemaMonitorAnomalyDetection_WarningSensitivityDesc,
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					Validators: []validator.String{
						sensitivityValidator,
					},
				},
				"critical_sensitivity": schema.StringAttribute{
					Description: schemaMonitorAnomalyDetection_CriticalSensitivityDesc,
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					Validators: []validator.String{
						sensitivityValidator,
					},
				},
				"max_check_attempts": schemaMonitorMaxCheckAttemptsAttr(3),
				"training_period_from": schema.Int64Attribute{
					Description: schemaMonitorAnomalyDetection_TrainingPeriodFromDesc,
					Optional:    true,
					Computed:    true,
					Default:     int64default.StaticInt64(0),
				},
				"scopes": schema.SetAttribute{
					ElementType: types.StringType,
					Description: schemaMonitorScopesDesc,
					Required:    true,
				},
			},
			Validators: []validator.Object{
				objectvalidator.AtLeastOneOf(path.MatchRoot("warning_sensitivity"), path.MatchRoot("critical_sensitivity")),
			},
		},
	}
}
