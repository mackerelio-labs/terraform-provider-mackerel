package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"

	"github.com/mackerelio-labs/terraform-provider-mackerel/mackerel"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: mackerel.Provider})
}
