package mackerel

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Test_ClientConfig_apiKeyCompat(t *testing.T) {
	cases := map[string]struct {
		MACKEREL_APIKEY  string
		MACKEREL_API_KEY string
		want             types.String
	}{
		"no key": {
			want: types.StringNull(),
		},
		"MACKEREL_API_KEY": {
			MACKEREL_API_KEY: "api_key",

			want: types.StringValue("api_key"),
		},
		"MACKEREL_APIKEY": {
			MACKEREL_APIKEY: "apikey",

			want: types.StringValue("apikey"),
		},
		"both": {
			MACKEREL_APIKEY:  "apikey",
			MACKEREL_API_KEY: "api_key",

			want: types.StringValue("apikey"),
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("MACKEREL_APIKEY", tt.MACKEREL_APIKEY)
			t.Setenv("MACKEREL_API_KEY", tt.MACKEREL_API_KEY)

			c := NewClientConfigFromEnv()
			if c.APIKey != tt.want {
				t.Errorf("expected to be '%s', but got '%s'.", tt.want, c.APIKey)
			}
		})
	}
}

func Test_ClientConfig_noApiKey(t *testing.T) {
	t.Parallel()

	var config ClientConfigModel
	if _, err := config.NewClient(); !errors.Is(err, ErrNoAPIKey) {
		t.Errorf("expected to ErrNoAPIKey, but got: %v", err)
	}
}
