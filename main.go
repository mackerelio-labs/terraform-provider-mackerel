package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"

	"github.com/xcezx/terraform-provider-mackerel/mackerel"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: mackerel.Provider})
}
