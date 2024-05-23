package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"

	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
	"github.com/mackerelio-labs/terraform-provider-mackerel/mackerel"
)

const (
	providerAddr = "registry.terraform.io/mackerelio-labs/mackerel"
)

func main() {
	ctx := context.Background()

	var debug bool
	flag.BoolVar(&debug, "debug", false, "run as debug-mode")

	flag.Parse()

	providers := []func() tfprotov5.ProviderServer{
		mackerel.Provider().GRPCProvider,
		providerserver.NewProtocol5(provider.New()),
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf5server.ServeOpt
	if debug {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	if err := tf5server.Serve(
		providerAddr,
		muxServer.ProviderServer,
		serveOpts...,
	); err != nil {
		log.Fatal(err)
	}
}
