package mackerel

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/mackerelio/mackerel-client-go"
)

const (
	channelTypeEmail   = "email"
	channelTypeSlack   = "slack"
	channelTypeWebhook = "webhook"
)

const (
	channelEventAlert        = "alert"
	channelEventAlertGroup   = "alertGroup"
	channelEventHostStatus   = "hostStatus"
	channelEventHostRegister = "hostRegister"
	channelEventHostRetire   = "hostRetire"
	channelEventMonitor      = "monitor"
)

type emailChannel struct {
	ID      string   `json:"id"`
	Type    string   `json:"type"`
	Name    string   `json:"name"`
	Emails  []string `json:"emails"`
	UserIds []string `json:"userIds"`
	Events  []string `json:"events"`
}

type slackChannel struct {
	ID                string    `json:"id"`
	Type              string    `json:"type"`
	Name              string    `json:"name"`
	URL               string    `json:"url"`
	Mentions          *mentions `json:"mentions"`
	EnabledGraphImage bool      `json:"enabledGraphImage"`
	Events            []string  `json:"events"`
}
type mentions struct {
	Ok       string `json:"ok"`
	Warning  string `json:"warning"`
	Critical string `json:"critical"`
}

type webhookChannel struct {
	ID     string   `json:"id"`
	Type   string   `json:"type"`
	Name   string   `json:"name"`
	URL    string   `json:"url"`
	Events []string `json:"events"`
}

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
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"slack", "webhook"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"emails": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"user_ids": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"events": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									channelEventAlert,
									channelEventAlertGroup}, false),
							},
						},
					},
				},
			},
			"slack": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"email", "webhook"},
				MaxItems:      1,
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
							Default:  true,
						},
						"events": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									channelEventAlert,
									channelEventAlertGroup,
									channelEventHostStatus,
									channelEventHostRegister,
									channelEventHostRetire,
									channelEventMonitor,
								}, false),
							},
						},
					},
				},
			},
			"webhook": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"email", "slack"},
				MaxItems:      1,
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
									channelEventAlert,
									channelEventAlertGroup,
									channelEventHostStatus,
									channelEventHostRegister,
									channelEventHostRetire,
									channelEventMonitor,
								}, false),
							},
						},
					},
				},
			},
		},
	}
}

func buildChannelStruct(d *schema.ResourceData) (interface{}, error) {
	if d.Get("email.#").(int) > 0 {
		emailChannel := &emailChannel{
			Type:    channelTypeEmail,
			Name:    d.Get("name").(string),
			Emails:  []string{},
			UserIds: []string{},
			Events:  []string{},
		}
		if email, ok := d.GetOk("email.0"); ok {
			attr := email.(map[string]interface{})
			emailChannel.Emails = expandStringSlice(attr["emails"].([]interface{}))
			emailChannel.UserIds = expandStringSlice(attr["user_ids"].([]interface{}))
			emailChannel.Events = expandStringSlice(attr["events"].([]interface{}))
		}
		return emailChannel, nil
	}

	if slack, ok := d.GetOk("slack.0"); ok {
		attr := slack.(map[string]interface{})
		return &slackChannel{
			Type:              channelTypeSlack,
			Name:              d.Get("name").(string),
			URL:               attr["url"].(string),
			Mentions:          &mentions{},
			EnabledGraphImage: attr["enabled_graph_image"].(bool),
			Events:            expandStringSlice(attr["events"].([]interface{})),
		}, nil
	}

	if webhook, ok := d.GetOk("webhook.0"); ok {
		attr := webhook.(map[string]interface{})
		return &webhookChannel{
			Type:   channelTypeWebhook,
			Name:   d.Get("name").(string),
			URL:    attr["url"].(string),
			Events: expandStringSlice(attr["events"].([]interface{})),
		}, nil
	}

	return nil, errors.New("unknown channel type")
}

func resourceMackerelChannelCreate(d *schema.ResourceData, meta interface{}) error {
	ch, err := buildChannelStruct(d)
	if err != nil {
		return err
	}

	client := meta.(*mackerel.Client)
	resp, err := client.PostJSON("/api/v0/channels", ch)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	var data struct {
		ID string `json:"id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}
	d.SetId(data.ID)
	return resourceMackerelChannelRead(d, meta)
}

func resourceMackerelChannelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	newURL := client.BaseURL
	newURL.Path = "/api/v0/channels"
	req, err := http.NewRequest("GET", newURL.String(), nil)
	if err != nil {
		return err
	}
	resp, err := client.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return err
	}

	var data struct {
		Channels []json.RawMessage `json:"channels"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}

	for _, rawCh := range data.Channels {
		var chData struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		}
		if err := json.Unmarshal(rawCh, &chData); err != nil {
			return err
		}
		if chData.ID == d.Id() {
			switch chData.Type {
			case channelTypeEmail:
				ch := &emailChannel{}
				if err := json.Unmarshal(rawCh, ch); err != nil {
					return err
				}
				_ = d.Set("email", ch)
			case channelTypeSlack:
				ch := &slackChannel{}
				if err := json.Unmarshal(rawCh, ch); err != nil {
					return err
				}
				_ = d.Set("slack", ch)
			case channelTypeWebhook:
				ch := &webhookChannel{}
				if err := json.Unmarshal(rawCh, ch); err != nil {
					return err
				}
				_ = d.Set("webhook", ch)
			}
			break
		}
	}

	return nil
}

func resourceMackerelChannelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	newURL := client.BaseURL
	newURL.Path = fmt.Sprintf("/api/v0/channels/%s", d.Id())
	req, err := http.NewRequest("DELETE", newURL.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
