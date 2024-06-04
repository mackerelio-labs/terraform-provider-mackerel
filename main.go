package main

import (
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"

	"github.com/mackerelio-labs/terraform-provider-mackerel/mackerel"
)

const (
	providerAddr = "registry.terraform.io/mackerelio-labs/mackerel"
)

func main() {
	// No timestamp to logs
	// FYI: https://developer.hashicorp.com/terraform/plugin/log/writing#duplicate-timestamp-and-incorrect-level-messages
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	var debug bool
	flag.BoolVar(&debug, "debug", false, "run as debug-mode")

	flag.Parse()

	var serveOpts []tf5server.ServeOpt
	if debug {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	if err := tf5server.Serve(
		providerAddr,
		mackerel.ProtoV5ProviderServer,
		serveOpts...,
	); err != nil {
		log.Printf("[ERROR] failed to start server: %v", err)
		panic(err)
	}
}
