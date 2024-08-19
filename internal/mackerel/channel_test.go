package mackerel

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

const (
	testChannelSlackURL   = "https://slack.test/services/xxx/yyy/zzz"
	testChannelWebhookURL = "https://example.test/hook"
)

func Test_Channel_conv(t *testing.T) {
	t.Parallel()

	// api <-> model
	cases := map[string]struct {
		api   mackerel.Channel
		model ChannelModel
	}{
		"slack": {
			api: mackerel.Channel{
				ID:                "5eKHBxJS5u9",
				Name:              "slack",
				Type:              "slack",
				Mentions:          mackerel.Mentions{},
				EnabledGraphImage: ptr(false),
				URL:               testChannelSlackURL,
				Events:            ptr([]string{}),
			},

			model: ChannelModel{
				ID:   types.StringValue("5eKHBxJS5u9"),
				Name: types.StringValue("slack"),
				Slack: []ChannelSlackModel{{
					URL:               types.StringValue(testChannelSlackURL),
					EnabledGraphImage: types.BoolValue(false),
					Events:            []string{},
				}},
			},
		},
		"slack-full": {
			api: mackerel.Channel{
				ID:   "5eKHBxJS5u9",
				Name: "slack-full",
				Type: "slack",
				Mentions: mackerel.Mentions{
					OK:       "OK!!!",
					Warning:  "WARNING!!!",
					Critical: "CRITICAL!!!",
				},
				EnabledGraphImage: ptr(true),
				URL:               testChannelSlackURL,
				Events:            ptr([]string{"alert"}),
			},

			model: ChannelModel{
				ID:   types.StringValue("5eKHBxJS5u9"),
				Name: types.StringValue("slack-full"),
				Slack: []ChannelSlackModel{{
					URL: types.StringValue(testChannelSlackURL),
					Mentions: map[string]string{
						"ok":       "OK!!!",
						"warning":  "WARNING!!!",
						"critical": "CRITICAL!!!",
					},
					EnabledGraphImage: types.BoolValue(true),
					Events:            []string{"alert"},
				}},
			},
		},
		"webhook": {
			api: mackerel.Channel{
				ID:     "5eKHBxRHcMJ",
				Name:   "webhook",
				Type:   "webhook",
				URL:    testChannelWebhookURL,
				Events: ptr([]string{}),
			},
			model: ChannelModel{
				ID:   types.StringValue("5eKHBxRHcMJ"),
				Name: types.StringValue("webhook"),
				Webhook: []ChannelWebhookModel{{
					URL:    types.StringValue(testChannelWebhookURL),
					Events: []string{},
				}},
			},
		},
		"email": {
			api: mackerel.Channel{
				ID:      "5eKHBzgCmAd",
				Name:    "email",
				Type:    "email",
				Emails:  ptr([]string{"john.doe@example.test"}),
				UserIDs: ptr([]string{"john"}),
				Events:  ptr([]string{"alertGroup"}),
			},
			model: ChannelModel{
				ID:   types.StringValue("5eKHBzgCmAd"),
				Name: types.StringValue("email"),
				Email: []ChannelEmailModel{{
					Emails:  []string{"john.doe@example.test"},
					UserIDs: []string{"john"},
					Events:  []string{"alertGroup"},
				}},
			},
		},
	}

	for name, tt := range cases {
		t.Run(name+"-toModel", func(t *testing.T) {
			t.Parallel()

			m, err := newChannel(tt.api)
			if err != nil {
				return
			}

			if diff := cmp.Diff(m, tt.model); diff != "" {
				t.Error(diff)
			}
		})
		t.Run(name+"-toAPI", func(t *testing.T) {
			t.Parallel()

			c := tt.model.mackerelChannel()
			if diff := cmp.Diff(c, tt.api); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func Test_Channel_merge(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in    ChannelModel
		inNew ChannelModel
		wants ChannelModel
	}{
		"slack": {
			in: ChannelModel{
				ID:   types.StringValue("5eKHBxJS5u9"),
				Name: types.StringValue("slack"),
				Slack: []ChannelSlackModel{{
					URL:               types.StringValue(testChannelSlackURL),
					EnabledGraphImage: types.BoolValue(false),
					Events:            nil,
				}},
			},
			inNew: ChannelModel{
				ID:   types.StringValue("5eKHBxJS5u9"),
				Name: types.StringValue("slack"),
				Slack: []ChannelSlackModel{{
					URL:               types.StringValue(testChannelSlackURL),
					EnabledGraphImage: types.BoolValue(true), // changed
					Events:            []string{},
				}},
			},
			wants: ChannelModel{
				ID:   types.StringValue("5eKHBxJS5u9"),
				Name: types.StringValue("slack"),
				Slack: []ChannelSlackModel{{
					URL:               types.StringValue(testChannelSlackURL),
					EnabledGraphImage: types.BoolValue(true),
					Events:            nil, // preserved
				}},
			},
		},
		"webhook": {
			in: ChannelModel{
				ID:   types.StringValue("5eKHBxRHcMJ"),
				Name: types.StringValue("webhook"),
				Webhook: []ChannelWebhookModel{{
					URL:    types.StringValue(testChannelWebhookURL),
					Events: nil,
				}},
			},
			inNew: ChannelModel{
				ID:   types.StringValue("5eKHBxRHcMJ"),
				Name: types.StringValue("webhook-changed"),
				Webhook: []ChannelWebhookModel{{
					URL:    types.StringValue(testChannelWebhookURL),
					Events: []string{},
				}},
			},
			wants: ChannelModel{
				ID:   types.StringValue("5eKHBxRHcMJ"),
				Name: types.StringValue("webhook-changed"),
				Webhook: []ChannelWebhookModel{{
					URL:    types.StringValue(testChannelWebhookURL),
					Events: nil,
				}},
			},
		},
		"email": {
			in: ChannelModel{
				ID:   types.StringValue("5eKHBzgCmAd"),
				Name: types.StringValue("email"),
				Email: []ChannelEmailModel{{
					Emails:  nil,
					UserIDs: nil,
					Events:  nil,
				}},
			},
			inNew: ChannelModel{
				ID:   types.StringValue("5eKHBzgCmAd"),
				Name: types.StringValue("email-changed"),
				Email: []ChannelEmailModel{{
					Emails:  []string{},
					UserIDs: []string{},
					Events:  []string{},
				}},
			},
			wants: ChannelModel{
				ID:   types.StringValue("5eKHBzgCmAd"),
				Name: types.StringValue("email-changed"),
				Email: []ChannelEmailModel{{
					Emails:  nil,
					UserIDs: nil,
					Events:  nil,
				}},
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := tt.in
			m.merge(tt.inNew)
			if diff := cmp.Diff(m, tt.wants); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func ptr[T any](x T) *T {
	return &x
}
