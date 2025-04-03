package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/mackerel"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/typeutil"
)

var (
	_ datasource.DataSource              = (*mackerelMonitorDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*mackerelMonitorDataSource)(nil)
)

type mackerelMonitorDataSource struct {
	Client *mackerel.Client
}

func NewMackerelMonitorDataSource() datasource.DataSource {
	return &mackerelMonitorDataSource{}
}

func (d *mackerelMonitorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (d *mackerelMonitorDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: schemaMonitorIDDesc,
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: schemaMonitorNameDesc,
				Computed:    true,
			},
			"memo": schema.StringAttribute{
				Description: schemaMonitorMemoDesc,
				Computed:    true,
			},
			"is_mute": schema.BoolAttribute{
				Description: schemaMonitorIsMuteDesc,
				Computed:    true,
			},
			"notification_interval": schema.Int64Attribute{
				Description: schemaMonitorNotificationIntervalDesc,
				Computed:    true,
			},
			"host_metric": schema.ListAttribute{
				Description: schemaMonitorHostMetricDesc,
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"metric":             types.StringType,
						"operator":           types.StringType,
						"warning":            typeutil.FloatStringType{},
						"critical":           typeutil.FloatStringType{},
						"duration":           types.Int64Type,
						"max_check_attempts": types.Int64Type,
						"scopes": types.SetType{
							ElemType: types.StringType,
						},
						"exclude_scopes": types.SetType{
							ElemType: types.StringType,
						},
					},
				},
			},
			"service_metric": schema.ListAttribute{
				Description: schemaMonitorServiceMetricDesc,
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"service":                   types.StringType,
						"metric":                    types.StringType,
						"operator":                  types.StringType,
						"warning":                   typeutil.FloatStringType{},
						"critical":                  typeutil.FloatStringType{},
						"duration":                  types.Int64Type,
						"max_check_attempts":        types.Int64Type,
						"missing_duration_warning":  types.Int64Type,
						"missing_duration_critical": types.Int64Type,
					},
				},
			},
			"expression": schema.ListAttribute{
				Description: schemaMonitorExpressionDesc,
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"expression": types.StringType,
						"operator":   types.StringType,
						"warning":    typeutil.FloatStringType{},
						"critical":   typeutil.FloatStringType{},
					},
				},
			},
			"query": schema.ListAttribute{
				Description: schemaMonitorQueryDesc,
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"query":    types.StringType,
						"legend":   types.StringType,
						"operator": types.StringType,
						"warning":  typeutil.FloatStringType{},
						"critical": typeutil.FloatStringType{},
					},
				},
			},
			"connectivity": schema.ListAttribute{
				Description: schemaMonitorConnectivityDesc,
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"alert_status_on_gone": types.StringType,
						"scopes": types.SetType{
							ElemType: types.StringType,
						},
						"exclude_scopes": types.SetType{
							ElemType: types.StringType,
						},
					},
				},
			},
			"external": schema.ListAttribute{
				Description: schemaMonitorExternalDesc,
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"max_check_attempts": types.Int64Type,
						"url":                types.StringType,
						"method":             types.StringType,
						"request_body":       types.StringType,
						"headers": types.MapType{
							ElemType: types.StringType,
						},
						"service":                types.StringType,
						"response_time_critical": types.Float64Type,
						"response_time_warning":  types.Float64Type,
						"response_time_duration": types.Int64Type,

						"contains_string":      types.StringType,
						"follow_redirect":      types.BoolType,
						"expected_status_code": types.Int64Type,

						"skip_certificate_verification":     types.BoolType,
						"certification_expiration_critical": types.Int64Type,
						"certification_expiration_warning":  types.Int64Type,
					},
				},
			},
			"anomaly_detection": schema.ListAttribute{
				Description: schemaMonitorAnomalyDetectionDesc,
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"warning_sensitivity":  types.StringType,
						"critical_sensitivity": types.StringType,
						"max_check_attempts":   types.Int64Type,
						"training_period_from": types.Int64Type,
						"scopes": types.SetType{
							ElemType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *mackerelMonitorDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, diags := retrieveClient(ctx, req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	d.Client = client
}

func (d *mackerelMonitorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mackerel.MonitorModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := data.Read(ctx, d.Client); err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Monitor",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
