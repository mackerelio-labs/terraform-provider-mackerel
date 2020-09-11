package mackerel

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/mackerelio/mackerel-client-go"
)

type Config struct {
	APIKey string
}

func (c *Config) Client() (client *mackerel.Client, diags diag.Diagnostics) {
	if c.APIKey == "" {
		return nil, diag.Errorf("no API Key for Mackerel")
	}
	client = mackerel.NewClient(c.APIKey)
	client.HTTPClient.Transport = logging.NewTransport("Mackerel", http.DefaultTransport)
	return client, diags
}
