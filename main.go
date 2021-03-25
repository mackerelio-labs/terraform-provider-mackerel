package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"

	"github.com/mackerelio-labs/terraform-provider-mackerel/mackerel"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: mackerel.Provider})
}
