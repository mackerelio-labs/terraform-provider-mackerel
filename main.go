package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/mackerelio-labs/terraform-provider-mackerel/mackerel"
)

const (
	providerAddr = "registry.terraform.io/mackerelio-labs/mackerel"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "run as debug-mode")

	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: mackerel.Provider,
	}

	if debug {
		c := context.TODO() // to support cancellation operations such as signal handling in the future.
		if err := plugin.Debug(c, providerAddr, opts); err != nil {
			log.Fatal(err)
		}
		return
	}

	plugin.Serve(opts)
}
