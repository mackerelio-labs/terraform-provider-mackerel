package mackerel

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

type (
	ChannelModel struct {
		ID      types.String          `tfsdk:"id"`
		Name    types.String          `tfsdk:"name"`
		Email   []ChannelEmailModel   `tfsdk:"email"`
		Slack   []ChannelSlackModel   `tfsdk:"slack"`
		Webhook []ChannelWebhookModel `tfsdk:"webhook"`
	}
	ChannelEmailModel struct {
		Emails  []string `tfsdk:"emails"`
		UserIDs []string `tfsdk:"user_ids"`
		Events  []string `tfsdk:"events"`
	}
	ChannelSlackModel struct {
		URL               types.String      `tfsdk:"url"`
		Mentions          map[string]string `tfsdk:"mentions"`
		EnabledGraphImage types.Bool        `tfsdk:"enabled_graph_image"`
		Events            []string          `tfsdk:"events"`
	}
	ChannelWebhookModel struct {
		URL    types.String `tfsdk:"url"`
		Events []string     `tfsdk:"events"`
	}
)

// Reads a channel by the ID.
// Currently this function is NOT cancelable.
func ReadChannel(_ context.Context, client *Client, id string) (ChannelModel, error) {
	channels, err := client.FindChannels()
	if err != nil {
		return ChannelModel{}, err
	}

	channelIdx := slices.IndexFunc(channels, func(c *mackerel.Channel) bool {
		return c.ID == id
	})
	if channelIdx < 0 {
		return ChannelModel{}, fmt.Errorf("the ID '%s' does not match any channel in mackerel.io", id)
	}

	channel, err := newChannel(*channels[channelIdx])
	if err != nil {
		return ChannelModel{}, err
	}

	return channel, nil
}

// Creates a new channel.
// Currently this function is NOT cancelable.
func (m *ChannelModel) Create(_ context.Context, client *Client) error {
	channelParam := m.mackerelChannel()
	channel, err := client.CreateChannel(&channelParam)
	if err != nil {
		return err
	}

	m.ID = types.StringValue(channel.ID)
	return nil
}

// Reads a channel.
// Currently this function is NOT cancelable.
func (m *ChannelModel) Read(ctx context.Context, client *Client) error {
	newModel, err := ReadChannel(ctx, client, m.ID.ValueString())
	if err != nil {
		return err
	}

	m.merge(newModel)
	return nil
}

// Deletes a channel.
func (m *ChannelModel) Delete(_ context.Context, client *Client) error {
	if _, err := client.DeleteChannel(m.ID.ValueString()); err != nil {
		return err
	}
	return nil
}

func newChannel(mackerelChannel mackerel.Channel) (ChannelModel, error) {
	model := ChannelModel{
		ID:   types.StringValue(mackerelChannel.ID),
		Name: types.StringValue(mackerelChannel.Name),
	}
	switch mackerelChannel.Type {
	case "email":
		model.Email = []ChannelEmailModel{{
			Emails:  *mackerelChannel.Emails,
			UserIDs: *mackerelChannel.UserIDs,
			Events:  *mackerelChannel.Events,
		}}
		return model, nil
	case "slack":
		slackModel := ChannelSlackModel{
			URL:               types.StringValue(mackerelChannel.URL),
			EnabledGraphImage: types.BoolPointerValue(mackerelChannel.EnabledGraphImage),
			Events:            *mackerelChannel.Events,
		}
		if mackerelChannel.Mentions != (mackerel.Mentions{}) {
			mentions := make(map[string]string, 3)
			if mackerelChannel.Mentions.OK != "" {
				mentions["ok"] = mackerelChannel.Mentions.OK
			}
			if mackerelChannel.Mentions.Warning != "" {
				mentions["warning"] = mackerelChannel.Mentions.Warning
			}
			if mackerelChannel.Mentions.Critical != "" {
				mentions["critical"] = mackerelChannel.Mentions.Critical
			}
			slackModel.Mentions = mentions
		}
		model.Slack = []ChannelSlackModel{slackModel}
		return model, nil
	case "webhook":
		model.Webhook = []ChannelWebhookModel{{
			URL:    types.StringValue(mackerelChannel.URL),
			Events: *mackerelChannel.Events,
		}}
		return model, nil
	default:
		return ChannelModel{}, fmt.Errorf("unsupported channel type: %s", mackerelChannel.Type)
	}
}

func (m ChannelModel) mackerelChannel() mackerel.Channel {
	channel := mackerel.Channel{
		ID:     m.ID.ValueString(),
		Name:   m.Name.ValueString(),
		Events: &[]string{},
	}
	if len(m.Email) > 0 {
		emailModel := m.Email[0]
		channel.Type = "email"
		if len(emailModel.Emails) > 0 {
			channel.Emails = &emailModel.Emails
		} else {
			channel.Emails = &[]string{}
		}
		if len(emailModel.UserIDs) > 0 {
			channel.UserIDs = &emailModel.UserIDs
		} else {
			channel.UserIDs = &[]string{}
		}
		if len(emailModel.Events) > 0 {
			channel.Events = &emailModel.Events
		}
	} else if len(m.Slack) > 0 {
		slackModel := m.Slack[0]
		channel.Type = "slack"
		channel.URL = slackModel.URL.ValueString()
		if slackModel.Mentions != nil {
			if okMention, ok := slackModel.Mentions["ok"]; ok {
				channel.Mentions.OK = okMention
			}
			if warnMention, ok := slackModel.Mentions["warning"]; ok {
				channel.Mentions.Warning = warnMention
			}
			if critMention, ok := slackModel.Mentions["critical"]; ok {
				channel.Mentions.Critical = critMention
			}
		}
		channel.EnabledGraphImage = slackModel.EnabledGraphImage.ValueBoolPointer()
		if len(slackModel.Events) > 0 {
			channel.Events = &slackModel.Events
		}
	} else if len(m.Webhook) > 0 {
		webhookModel := m.Webhook[0]
		channel.Type = "webhook"
		channel.URL = webhookModel.URL.ValueString()
		if len(webhookModel.Events) > 0 {
			channel.Events = &webhookModel.Events
		}
	}
	return channel
}

func (m *ChannelModel) merge(newModel ChannelModel) {
	if len(m.Email) > 0 && len(newModel.Email) > 0 {
		oldEmail := m.Email[0]
		newEmail := &newModel.Email[0]

		// Distinct between null and [] in email.emails
		// If both side of emails are empty, preserve old one.
		if len(oldEmail.Emails) == 0 && len(newEmail.Emails) == 0 {
			newEmail.Emails = oldEmail.Emails
		}

		// Distinct between null and [] in email.user_ids
		if len(oldEmail.UserIDs) == 0 && len(newEmail.UserIDs) == 0 {
			newEmail.UserIDs = oldEmail.UserIDs
		}

		// Distinct between null and [] in email.events
		if len(oldEmail.Events) == 0 && len(newEmail.Events) == 0 {
			newEmail.Events = oldEmail.Events
		}
	}
	if len(m.Slack) > 0 && len(newModel.Slack) > 0 {
		oldSlack := m.Slack[0]
		newSlack := &newModel.Slack[0]

		// Distinct between null and {} in slack.mentions
		if len(oldSlack.Mentions) == 0 && len(newSlack.Mentions) == 0 {
			newSlack.Mentions = oldSlack.Mentions
		}

		// Distinct between null and [] in slack.events
		if len(oldSlack.Events) == 0 && len(newSlack.Events) == 0 {
			newSlack.Events = oldSlack.Events
		}
	}
	if len(m.Webhook) > 0 && len(newModel.Webhook) > 0 {
		// Distinct between null and [] in webhook.events
		if len(m.Webhook[0].Events) == 0 && len(newModel.Webhook[0].Events) == 0 {
			newModel.Webhook[0].Events = m.Webhook[0].Events
		}
	}

	*m = newModel
}
