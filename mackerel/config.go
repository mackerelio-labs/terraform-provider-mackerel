package mackerel

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	"github.com/mackerelio/mackerel-client-go"
)

type config struct {
	APIKey string
}

func (c *config) Client() (*mackerel.Client, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("no API Key for Mackerel")
	}
	client := mackerel.NewClient(c.APIKey)
	client.HTTPClient.Transport = logging.NewTransport("Mackerel", http.DefaultTransport)
	return client, nil
}
