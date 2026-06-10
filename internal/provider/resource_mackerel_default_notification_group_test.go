package provider_test

import (
	"context"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
)

func Test_MackerelDefaultNotificationGroupResource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwresource.SchemaRequest{}
	resp := &fwresource.SchemaResponse{}
	provider.NewMackerelDefaultNotificationGroupResource().Schema(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %+v", resp.Diagnostics)
	}

	if _, ok := resp.Schema.Blocks["monitor"]; ok {
		t.Fatal("default notification group resource must not expose monitor block")
	}
	if _, ok := resp.Schema.Blocks["service"]; ok {
		t.Fatal("default notification group resource must not expose service block")
	}
	if _, ok := resp.Schema.Attributes["id"]; !ok {
		t.Fatal("default notification group resource must expose id attribute")
	}
	if _, ok := resp.Schema.Attributes["name"]; ok {
		t.Fatal("default notification group resource must not expose name attribute")
	}
	if _, ok := resp.Schema.Attributes["notification_level"]; !ok {
		t.Fatal("default notification group resource must expose notification_level attribute")
	}
	if _, ok := resp.Schema.Attributes["child_notification_group_ids"]; !ok {
		t.Fatal("default notification group resource must expose child_notification_group_ids attribute")
	}
	assertRequiredSetAttribute(t, resp.Schema.Attributes["child_notification_group_ids"])
	assertRequiredSetAttribute(t, resp.Schema.Attributes["child_channel_ids"])

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema validation diagnostics: %+v", diags)
	}
}

func assertRequiredSetAttribute(t *testing.T, attr schema.Attribute) {
	t.Helper()

	setAttr, ok := attr.(schema.SetAttribute)
	if !ok {
		t.Fatalf("attribute type = %T, want schema.SetAttribute", attr)
	}
	if !setAttr.IsRequired() {
		t.Fatal("attribute must be required")
	}
	if setAttr.IsOptional() {
		t.Fatal("attribute must not be optional")
	}
	if setAttr.IsComputed() {
		t.Fatal("attribute must not be computed")
	}
}
