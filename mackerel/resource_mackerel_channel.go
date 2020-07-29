package mackerel

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelChannelCreate,
		Read:   resourceMackerelChannelRead,
		Delete: resourceMackerelChannelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"email": {
				Type:         schema.TypeList,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"email", "slack", "webhook"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"emails": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"user_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"events": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"alert", "alertGroup"}, false),
							},
						},
					},
				},
			},
			"slack": {
				Type:         schema.TypeList,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"email", "slack", "webhook"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"mentions": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ok": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"warning": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"critical": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"enabled_graph_image": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
						},
						"events": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"alert", "alertGroup", "hostStatus", "hostRegister", "hostRetire", "monitor"}, false),
							},
						},
					},
				},
			},
			"webhook": {
				Type:         schema.TypeList,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"email", "slack", "webhook"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"events": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"alert", "alertGroup", "hostStatus", "hostRegister", "hostRetire", "monitor"}, false),
							},
						},
					},
				},
			},
		},
	}
}

func resourceMackerelChannelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	channel, err := client.CreateChannel(expandChannel(d))
	if err != nil {
		return err
	}
	d.SetId(channel.ID)
	return resourceMackerelChannelRead(d, meta)
}

func resourceMackerelChannelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	channels, err := client.FindChannels()
	if err != nil {
		return err
	}
	var channel *mackerel.Channel
	for _, c := range channels {
		if c.ID == d.Id() {
			channel = c
			break
		}
	}
	if channel == nil {
		return fmt.Errorf("the ID '%s' does not match any channel in mackerel.io", d.Id())
	}
	return flattenChannel(channel, d)
}

func resourceMackerelChannelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	_, err := client.DeleteChannel(d.Id())
	return err
}

func expandChannel(d *schema.ResourceData) *mackerel.Channel {
	channel := &mackerel.Channel{
		Name: d.Get("name").(string),
	}
	if _, ok := d.GetOk("email"); ok {
		channel.Type = "email"

		emails := expandStringListFromSet(d.Get("email.0.emails").(*schema.Set))
		channel.Emails = &emails

		userIDs := expandStringListFromSet(d.Get("email.0.user_ids").(*schema.Set))
		channel.UserIDs = &userIDs

		events := expandStringListFromSet(d.Get("email.0.events").(*schema.Set))
		channel.Events = &events
	}
	if _, ok := d.GetOk("slack"); ok {
		channel.Type = "slack"
		channel.URL = d.Get("slack.0.url").(string)
		channel.Mentions = mackerel.Mentions{
			OK:       d.Get("slack.0.mentions.ok").(string),
			Warning:  d.Get("slack.0.mentions.warning").(string),
			Critical: d.Get("slack.0.mentions.critical").(string),
		}

		enabledGraphImage := d.Get("slack.0.enabled_graph_image").(bool)
		channel.EnabledGraphImage = &enabledGraphImage

		events := expandStringListFromSet(d.Get("slack.0.events").(*schema.Set))
		channel.Events = &events
	}
	if _, ok := d.GetOk("webhook"); ok {
		channel.Type = "webhook"
		channel.URL = d.Get("webhook.0.url").(string)

		events := expandStringListFromSet(d.Get("webhook.0.events").(*schema.Set))
		channel.Events = &events
	}
	return channel
}

func flattenChannel(channel *mackerel.Channel, d *schema.ResourceData) error {
	d.Set("name", channel.Name)
	switch channel.Type {
	case "email":
		d.Set("email", []map[string]interface{}{
			{
				"emails":   flattenStringListToSet(*channel.Emails),
				"user_ids": flattenStringListToSet(*channel.UserIDs),
				"events":   flattenStringListToSet(*channel.Events),
			},
		})
	case "slack":
		mentions := make(map[string]string)
		for k, v := range map[string]string{
			"ok":       channel.Mentions.OK,
			"warning":  channel.Mentions.Warning,
			"critical": channel.Mentions.Critical,
		} {
			if v != "" {
				mentions[k] = v
			}
		}
		d.Set("slack", []map[string]interface{}{
			{
				"url":                 channel.URL,
				"mentions":            mentions,
				"enabled_graph_image": channel.EnabledGraphImage,
				"events":              flattenStringListToSet(*channel.Events),
			},
		})
	case "webhook":
		d.Set("webhook", []map[string]interface{}{
			{
				"url":    channel.URL,
				"events": flattenStringListToSet(*channel.Events),
			},
		})
	}
	return nil
}
