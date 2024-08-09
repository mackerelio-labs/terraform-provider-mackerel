package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/mackerel"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/validatorutil"
)

type mackerelProvider struct{}

var _ provider.Provider = (*mackerelProvider)(nil)

func New() provider.Provider {
	return &mackerelProvider{}
}

func (m *mackerelProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mackerel"
}

func (m *mackerelProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "Mackerel API Key",
				Optional:    true,
				Sensitive:   true,
			},
			"api_base": schema.StringAttribute{
				Description: "Mackerel API BASE URL",
				Optional:    true,
				Sensitive:   true,
				Validators:  []validator.String{validatorutil.IsURLWithHTTPorHTTPS()},
			},
		},
	}
}

func (m *mackerelProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var schemaConfig mackerel.ClientConfigModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &schemaConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := mackerel.NewClientConfigFromEnv()
	// merge config
	if config.APIKey.IsUnknown() {
		config.APIKey = schemaConfig.APIKey
	}
	if config.APIBase.IsUnknown() {
		config.APIBase = schemaConfig.APIBase
	}

	client, err := config.NewClient()
	if err != nil {
		if errors.Is(err, mackerel.ErrNoAPIKey) {
			resp.Diagnostics.AddError(
				"No API Key",
				err.Error(),
			)
		} else {
			resp.Diagnostics.AddError(
				"Unable to create Mackerel Client",
				err.Error(),
			)
		}
		return
	}

	resp.ResourceData = client
	resp.DataSourceData = client
}

func (m *mackerelProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewMackerelMonitorResource,
		NewMackerelNotificationGroupResource,
		NewMackerelRoleResource,
		NewMackerelRoleMetadataResource,
		NewMackerelServiceResource,
		NewMackerelServiceMetadataResource,
	}
}

func (m *mackerelProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewMackerelMonitorDataSource,
		NewMackerelNotificationGroupDataSource,
		NewMackerelRoleDataSource,
		NewMackerelRoleMetadataDataSource,
		NewMackerelServiceDataSource,
		NewMackerelServiceMetadataDataSource,
		NewMackerelServiceMetricNamesDataSource,
	}
}

func retrieveClient(_ context.Context, providerData any) (client *mackerel.Client, diags diag.Diagnostics) {
	if /* ConfigureProvider RPC is not called */ providerData == nil {
		return
	}

	client, ok := providerData.(*mackerel.Client)
	if !ok {
		diags.AddError(

			"No Mackerel Client is configured.",
			fmt.Sprintf("Expected configured Mackerel client, but got: %T. Please report this issue.", providerData),
		)
		return
	}
	return
}
