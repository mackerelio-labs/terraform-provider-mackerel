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

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"email": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"emails": {
							Type:         schema.TypeList,
							Optional:     true,
							ForceNew:     true,
							Elem:         &schema.Schema{Type: schema.TypeString},
							AtLeastOneOf: []string{"email.0.emails", "email.0.user_ids", "email.0.events"},
						},
						"user_ids": {
							Type:         schema.TypeList,
							Optional:     true,
							ForceNew:     true,
							Elem:         &schema.Schema{Type: schema.TypeString},
							AtLeastOneOf: []string{"email.0.emails", "email.0.user_ids", "email.0.events"},
						},
						"events": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"alert", "alertGroup"}, false),
							},
							AtLeastOneOf: []string{"email.0.emails", "email.0.user_ids", "email.0.events"},
						},
					},
				},
			},
			"slack": {
				Type:         schema.TypeList,
				Optional:     true,
				ForceNew:     true,
				MaxItems:     1,
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
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"enabled_graph_image": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
						},
						"events": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"alert",
									"alertGroup",
									"hostStatus",
									"hostRegister",
									"hostRetire",
									"monitor",
								}, false),
							},
						},
					},
				},
			},
			"webhook": {
				Type:         schema.TypeList,
				Optional:     true,
				ForceNew:     true,
				MaxItems:     1,
				ExactlyOneOf: []string{"email", "slack", "webhook"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"events": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"alert",
									"alertGroup",
									"hostStatus",
									"hostRegister",
									"hostRetire",
									"monitor",
								}, false),
							},
						},
					},
				},
			},
		},
	}
}

func resourceMackerelChannelCreate(d *schema.ResourceData, meta interface{}) error {
	var param *mackerel.Channel

	if email, ok := d.GetOk("email.0"); ok {
		param = expandEmailChannel(d.Get("name").(string), email.(map[string]interface{}))
	}
	if slack, ok := d.GetOk("slack.0"); ok {
		param = expandSlackChannel(d.Get("name").(string), slack.(map[string]interface{}))
	}
	if webhook, ok := d.GetOk("webhook.0"); ok {
		param = expandWebhookChannel(d.Get("name").(string), webhook.(map[string]interface{}))
	}

	client := meta.(*mackerel.Client)
	channel, err := client.CreateChannel(param)
	if err != nil {
		return err
	}
	d.SetId(channel.ID)

	return resourceMackerelChannelRead(d, meta)
}

func expandEmailChannel(name string, email map[string]interface{}) *mackerel.Channel {
	emails := expandStringList(email["emails"].([]interface{}))
	userIDs := expandStringList(email["user_ids"].([]interface{}))
	events := expandStringList(email["events"].([]interface{}))

	return &mackerel.Channel{
		Name:    name,
		Type:    "email",
		Emails:  &emails,
		UserIDs: &userIDs,
		Events:  &events,
	}
}

func expandSlackChannel(name string, slack map[string]interface{}) *mackerel.Channel {
	var mentions mackerel.Mentions
	if v, exists := slack["mentions"]; exists {
		mentionsMap := v.(map[string]interface{})
		if ok, exists := mentionsMap["ok"]; exists {
			mentions.OK = ok.(string)
		}
		if warning, exists := mentionsMap["warning"]; exists {
			mentions.Warning = warning.(string)
		}
		if critical, exists := mentionsMap["critical"]; exists {
			mentions.Critical = critical.(string)
		}
	}
	events := expandStringList(slack["events"].([]interface{}))
	enabledGraphImage := slack["enabled_graph_image"].(bool)

	return &mackerel.Channel{
		Name:              name,
		Type:              "slack",
		URL:               slack["url"].(string),
		Mentions:          mentions,
		EnabledGraphImage: &enabledGraphImage,
		Events:            &events,
	}
}

func expandWebhookChannel(name string, webhook map[string]interface{}) *mackerel.Channel {
	events := expandStringList(webhook["events"].([]interface{}))
	return &mackerel.Channel{
		Name:   name,
		Type:   "webhook",
		URL:    webhook["url"].(string),
		Events: &events,
	}
}

func expandStringList(s []interface{}) []string {
	vs := make([]string, 0, len(s))
	for _, v := range s {
		vs = append(vs, v.(string))
	}
	return vs
}

func resourceMackerelChannelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	channels, err := client.FindChannels()
	if err != nil {
		return err
	}
	for _, channel := range channels {
		if channel.ID == d.Id() {
			_ = d.Set("type", channel.Type)
			_ = d.Set("name", channel.Name)
			_ = d.Set("events", channel.Events)
			switch channel.Type {
			case "email":
				_ = d.Set("emails", channel.Emails)
				_ = d.Set("user_ids", channel.UserIDs)
			case "slack":
				_ = d.Set("url", channel.URL)
				_ = d.Set("mentions", channel.Mentions)
				_ = d.Set("enabled_graph_image", channel.EnabledGraphImage)
			case "webhook":
				_ = d.Set("url", channel.URL)
			}
		}
	}
	return nil
}

func resourceMackerelChannelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	_, err := client.DeleteChannel(d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
