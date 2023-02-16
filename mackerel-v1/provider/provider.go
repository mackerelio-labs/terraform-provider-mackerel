package provider

import (
	"context"
	"os"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

type MackerelProvider struct{}

type MackerelProviderModel struct {
	APIKey  types.String `tfsdk:"api_key"`
	APIBase types.String `tfsdk:"api_base"`
}

var _ provider.Provider = &MackerelProvider{}

func New() provider.Provider {
	return &MackerelProvider{}
}

func (p *MackerelProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Mackerel API Key",
			},
			"api_base": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Mackerel API BASE URL",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`/^https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&\/=]*)$/`),
						"expected api_base to be a valid url",
					),
				},
			},
		},
	}
}

var Client *mackerel.Client

func (p *MackerelProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	apiKey := os.Getenv("MACKEREL_APIKEY")
	if apiKey == "" {
		apiKey = os.Getenv("MACKEREL_API_KEY")
	}
	apiBase := os.Getenv("API_BASE")

	var data MackerelProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if data.APIKey.ValueString() != "" {
		apiKey = data.APIKey.ValueString()
	}
	if data.APIBase.ValueString() != "" {
		apiBase = data.APIBase.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing Mackerel APIKey",
			"hoge",
		)
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
				err.Error(),
			)
		}
	}

	resp.ResourceData = client
	resp.DataSourceData = client
}

func (p *MackerelProvider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {

}

func (p *MackerelProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {

}

func (p *MackerelProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *MackerelProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}
