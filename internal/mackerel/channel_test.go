package mackerel

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

const (
	testChannelSlackURL   = "https://slack.test/services/xxx/yyy/zzz"
	testChannelWebhookURL = "https://example.test/hook"
)

func Test_ReadChannel_NotFound(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"channels":[]}`))
	}))
	t.Cleanup(server.Close)

	client, err := mackerel.NewClientWithOptions("dummy-api-key", server.URL, false)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = ReadChannel(t.Context(), client, "missing-id")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

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
				EnabledGraphImage: new(false),
				URL:               testChannelSlackURL,
				Events:            new([]string{}),
			},

			model: ChannelModel{
				ID:   types.StringValue("5eKHBxJS5u9"),
				Name: types.StringValue("slack"),
				Slack: []ChannelSlackModel{{
					URL:               types.StringValue(testChannelSlackURL),
					Mentions:          map[string]string{},
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
				EnabledGraphImage: new(true),
				URL:               testChannelSlackURL,
				Events:            new([]string{"alert"}),
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
				Events: new([]string{}),
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
				Emails:  new([]string{"john.doe@example.test"}),
				UserIDs: new([]string{"john"}),
				Events:  new([]string{"alertGroup"}),
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

func Test_Channel_Create_AddToDefaultNotificationGroupFalse(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	model := ChannelModel{
		Name: types.StringValue("slack"),
		Slack: []ChannelSlackModel{{
			URL:               types.StringValue(testChannelSlackURL),
			EnabledGraphImage: types.BoolValue(false),
		}},
	}
	client := &channelCreatorTester{ID: "channel-id"}

	if err := model.createInner(ctx, client); err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	if client.Request.AddToDefaultNotificationGroup == nil {
		t.Fatal("AddToDefaultNotificationGroup is nil")
	}
	if *client.Request.AddToDefaultNotificationGroup {
		t.Fatal("AddToDefaultNotificationGroup should be false")
	}
	if got := model.ID.ValueString(); got != "channel-id" {
		t.Fatalf("ID = %q, want %q", got, "channel-id")
	}
}

func Test_Channel_Update_DoesNotSendAddToDefaultNotificationGroup(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	model := ChannelModel{
		ID:   types.StringValue("channel-id"),
		Name: types.StringValue("slack"),
		Slack: []ChannelSlackModel{{
			URL:               types.StringValue(testChannelSlackURL),
			EnabledGraphImage: types.BoolValue(false),
		}},
	}
	client := &channelUpdaterTester{}

	if err := model.updateInner(ctx, client); err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	if client.ID != "channel-id" {
		t.Fatalf("updated ID = %q, want %q", client.ID, "channel-id")
	}
	if client.Request.AddToDefaultNotificationGroup != nil {
		t.Fatalf("AddToDefaultNotificationGroup should be nil, got %#v", *client.Request.AddToDefaultNotificationGroup)
	}
}

type channelCreatorTester struct {
	ID      string
	Request mackerel.Channel
}

func (ct *channelCreatorTester) CreateChannel(param *mackerel.Channel) (*mackerel.Channel, error) {
	ct.Request = *param
	data := *param
	data.ID = ct.ID
	return &data, nil
}

type channelUpdaterTester struct {
	ID      string
	Request mackerel.Channel
}

func (ut *channelUpdaterTester) UpdateChannelContext(_ context.Context, id string, param *mackerel.Channel) (*mackerel.Channel, error) {
	ut.ID = id
	ut.Request = *param
	data := *param
	data.ID = id
	return &data, nil
}

//go:fix inline
func ptr[T any](x T) *T {
	return new(x)
}
