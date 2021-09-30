package mackerel

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/mackerelio/mackerel-client-go"
)

type Config struct {
	APIKey  string
	APIBase string
}

func (c *Config) Client() (client *mackerel.Client, diags diag.Diagnostics) {
	var err error

	if c.APIKey == "" {
		return nil, diag.Errorf("no API Key for Mackerel")
	}

	if c.APIBase == "" {
		client = mackerel.NewClient(c.APIKey)
	} else {
		client, err = mackerel.NewClientWithOptions(c.APIKey, c.APIBase, false)
		if err != nil {
			return nil, diag.Errorf("failed to create mackerel client: %s", err)
		}
	}

	client.HTTPClient.Transport = logging.NewTransport("Mackerel", http.DefaultTransport)
	return client, diags
}
