package mackerel

import (
	"testing"
)

func TestConfigEmptyAPIKey(t *testing.T) {
	config := Config{
		APIKey: "",
	}
	_, diags := config.Client()
	if !diags.HasError() {
		t.Error("expected diagnostics has an error, but get no error")
	}
}
