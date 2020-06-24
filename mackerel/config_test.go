package mackerel

import (
	"testing"
)

func TestConfigEmptyAPIKey(t *testing.T) {
	config := Config{
		APIKey: "",
	}
	if _, err := config.Client(); err == nil {
		t.Fatalf("expected error, but got nil")
	}
}
