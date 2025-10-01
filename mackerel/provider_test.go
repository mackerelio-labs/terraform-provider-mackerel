package mackerel

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProvider *schema.Provider
var testAccProtoV5ProviderFactories map[string]func() (tfprotov5.ProviderServer, error)

func init() {
	testAccProvider = Provider()

	providerServer := protoV5ProviderServer(testAccProvider)
	testAccProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
		"mackerel": func() (tfprotov5.ProviderServer, error) {
			return providerServer, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_apiKeyCompat(t *testing.T) {
	config := make(map[string]interface{})
	c := terraform.NewResourceConfigRaw(config)
	ctx := context.Background()
	t.Run("MACKEREL_API_KEY", func(t *testing.T) {
		testSetenv(t, "MACKEREL_API_KEY", "apikey1")
		testSetenv(t, "MACKEREL_APIKEY", "")
		if err := Provider().Configure(ctx, c); err != nil {
			t.Errorf("Configure: %v", err)
		}
	})
	t.Run("MACKEREL_APIKEY", func(t *testing.T) {
		testSetenv(t, "MACKEREL_API_KEY", "")
		testSetenv(t, "MACKEREL_APIKEY", "apikey1")
		if err := Provider().Configure(ctx, c); err != nil {
			t.Errorf("Configure: %v", err)
		}
	})
}

func testSetenv(t testing.TB, name, val string) {
	t.Helper()

	setenv := func(name, val string) error {
		if val == "" {
			return os.Unsetenv(name)
		} else {
			return os.Setenv(name, val)
		}
	}

	s := os.Getenv(name)
	t.Cleanup(func() {
		if err := setenv(name, s); err != nil {
			t.Fatalf("Setenv(%q, %q): %v", name, val, err)
		}
	})
	if err := setenv(name, val); err != nil {
		t.Fatalf("Setenv(%q, %q): %v", name, val, err)
	}
}

func TestProviderImpl(t *testing.T) {
	var _ = Provider()
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("MACKEREL_API_KEY") == "" {
		t.Fatal("MACKEREL_API_KEY must be set for acceptance tests")
	}
}
