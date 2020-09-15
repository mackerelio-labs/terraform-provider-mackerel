package mackerel

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelChannel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMackerelChannelCreate,
		ReadContext:   resourceMackerelChannelRead,
		DeleteContext: resourceMackerelChannelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourceMackerelChannelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mackerel.Client)
	channel, err := client.CreateChannel(expandChannel(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(channel.ID)
	return resourceMackerelChannelRead(ctx, d, m)
}

func resourceMackerelChannelRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*mackerel.Client)
	channels, err := client.FindChannels()
	if err != nil {
		return diag.FromErr(err)
	}
	var channel *mackerel.Channel
	for _, c := range channels {
		if c.ID == d.Id() {
			channel = c
			break
		}
	}
	if channel == nil {
		return diag.Errorf("the ID '%s' does not match any channel in mackerel.io", d.Id())
	}
	if err := flattenChannel(channel, d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceMackerelChannelDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*mackerel.Client)
	_, err := client.DeleteChannel(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
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
