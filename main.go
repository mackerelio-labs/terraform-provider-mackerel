package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"

	"github.com/mackerelio-labs/terraform-provider-mackerel/mackerel"
)

const (
	providerAddr = "registry.terraform.io/mackerelio-labs/mackerel"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "run as debug-mode")

	flag.Parse()

	ctx := context.Background()

	providers := []func() tfprotov5.ProviderServer{
		mackerel.Provider().GRPCProvider,
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf5server.ServeOpt
	if debug {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	err = tf5server.Serve(
		providerAddr,
		muxServer.ProviderServer,
		serveOpts...,
	)
	if err != nil {
		log.Fatal(err)
	}
}
