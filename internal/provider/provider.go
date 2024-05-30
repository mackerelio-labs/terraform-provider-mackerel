package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/mackerelio/mackerel-client-go"
)

type mackerelProvider struct{}
type MackerelProviderModel struct {
	APIKey  types.String `tfsdk:"api_key"`
	APIBase types.String `tfsdk:"api_base"`
}

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
				Validators:  []validator.String{IsURLWithHTTPorHTTPS()},
			},
		},
	}
}

func (m *mackerelProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data MackerelProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := data.APIKey.ValueString()
	if data.APIKey.IsUnknown() || data.APIKey.IsNull() {
		apiKey = os.Getenv("MACKEREL_APIKEY")
		if apiKey == "" {
			apiKey = os.Getenv("MACKEREL_API_KEY")
		}
	}
	if apiKey == "" {
		resp.Diagnostics.AddError(
			"no API Key", "no API Key for Mackerel is found",
		)
	}

	apiBase := data.APIBase.ValueString()
	if data.APIBase.IsUnknown() || data.APIBase.IsNull() {
		apiBase = os.Getenv("API_BASE")
	}

	var client *mackerel.Client
	if apiBase == "" {
		client = mackerel.NewClient(apiKey)
	} else {
		var err error
		client, err = mackerel.NewClientWithOptions(apiKey, apiBase, false)
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to create mackerel client",
				fmt.Sprintf("failed to create mackerel client: %v", err),
			)
		}
	}

	// TODO: use logging transport with tflog (FYI: https://github.com/hashicorp/terraform-plugin-log/issues/91)
	client.HTTPClient.Transport = logging.NewSubsystemLoggingHTTPTransport("Mackerel", http.DefaultTransport)

	resp.ResourceData = client
}

func (m *mackerelProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewMackerelServiceResource,
	}
}

func (m *mackerelProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func retrieveClient(_ context.Context, providerData any) (client *mackerel.Client, diags diag.Diagnostics) {
	if /* ConfigureProvider RPC is not called */ providerData == nil {
		return
	}

	client, ok := providerData.(*mackerel.Client)
	if !ok {
		diags.AddError(

			"Unconfigured Mackerel Client",
			fmt.Sprintf("Expected configured Mackerel client, but got: %T. Please report this issue.", providerData),
		)
		return
	}
	return
}
