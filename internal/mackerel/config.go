package mackerel

import (
	"errors"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/mackerelio/mackerel-client-go"
)

type Client = mackerel.Client

type ClientConfigModel struct {
	APIKey  types.String `tfsdk:"api_key"`
	APIBase types.String `tfsdk:"api_base"`
}

var (
	ErrNoAPIKey = errors.New("API Key for Mackerel is not found.")
)

func NewClientConfigFromEnv() ClientConfigModel {
	var data ClientConfigModel

	var apiKey string
	for _, env := range []string{"MACKEREL_APIKEY", "MACKEREL_API_KEY"} {
		if apiKey == "" {
			apiKey = os.Getenv(env)
		}
	}
	if apiKey != "" {
		data.APIKey = types.StringValue(apiKey)
	}

	apiBase := os.Getenv("API_BASE")
	if apiBase != "" {
		data.APIBase = types.StringValue(apiBase)
	}

	return data
}

func (m *ClientConfigModel) NewClient() (*Client, error) {
	apiKey := m.APIKey.ValueString()
	if apiKey == "" {
		return nil, ErrNoAPIKey
	}

	apiBase := m.APIBase.ValueString()

	var client *mackerel.Client
	if apiBase == "" {
		client = mackerel.NewClient(apiKey)
	} else {
		// TODO: use logging transport with tflog (FYI: https://github.com/hashicorp/terraform-plugin-log/issues/91)
		c, err := mackerel.NewClientWithOptions(apiKey, apiBase, false)
		if err != nil {
			return nil, err
		}
		client = c
	}
	client.HTTPClient.Transport = logging.NewSubsystemLoggingHTTPTransport("Mackerel", http.DefaultTransport)
	return client, nil
}
