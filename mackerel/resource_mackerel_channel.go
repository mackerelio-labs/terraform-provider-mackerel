package mackerel

import (
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
	param := &mackerel.Channel{
		Name: d.Get("name").(string),
	}

	if d.Get("email.#") == 1 {
		param.Type = "email"

		emails := expandStringList(d.Get("email.0.emails").(*schema.Set).List())
		param.Emails = &emails

		userIDs := expandStringList(d.Get("email.0.user_ids").(*schema.Set).List())
		param.UserIDs = &userIDs

		events := expandStringList(d.Get("email.0.events").(*schema.Set).List())
		param.Events = &events
	}
	if d.Get("slack.#") == 1 {
		param.Type = "slack"

		param.URL = d.Get("slack.0.url").(string)

		param.Mentions = mackerel.Mentions{
			OK:       d.Get("slack.0.mentions.ok").(string),
			Warning:  d.Get("slack.0.mentions.warning").(string),
			Critical: d.Get("slack.0.mentions.critical").(string),
		}

		enabledGraphImage := d.Get("slack.0.enabled_graph_image").(bool)
		param.EnabledGraphImage = &enabledGraphImage

		events := expandStringList(d.Get("slack.0.events").(*schema.Set).List())
		param.Events = &events
	}
	if d.Get("webhook.#") == 1 {
		param.Type = "webhook"

		param.URL = d.Get("webhook.0.url").(string)

		events := expandStringList(d.Get("webhook.0.events").(*schema.Set).List())
		param.Events = &events
	}

	client := meta.(*mackerel.Client)
	channel, err := client.CreateChannel(param)
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

	for _, channel := range channels {
		if channel.ID == d.Id() {
			d.Set("name", channel.Name)

			switch channel.Type {
			case "email":
				d.Set("email", []map[string]interface{}{
					{
						"emails":   flattenStringSet(*channel.Emails),
						"user_ids": flattenStringSet(*channel.UserIDs),
						"events":   flattenStringSet(*channel.Events),
					},
				})
			case "slack":
				d.Set("slack", []map[string]interface{}{
					{
						"url": channel.URL,
						"mentions": map[string]interface{}{
							"ok":       channel.Mentions.OK,
							"warning":  channel.Mentions.Warning,
							"critical": channel.Mentions.Critical,
						},
						"enabled_graph_image": *channel.EnabledGraphImage,
						"events":              flattenStringSet(*channel.Events),
					},
				})
			case "webhook":
				d.Set("webhook", []map[string]interface{}{
					{
						"url":    channel.URL,
						"events": flattenStringSet(*channel.Events),
					},
				})
			}
			break
		}
	}
	return nil
}

func resourceMackerelChannelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	_, err := client.DeleteChannel(d.Id())
	return err
}
