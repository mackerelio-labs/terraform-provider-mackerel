package mackerel

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func dataSourceMackerelChannel() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMackerelChannelRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"emails": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"user_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"events": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"slack": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mentions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ok": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"warning": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"critical": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"enabled_graph_image": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"events": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"webhook": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"events": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceMackerelChannelRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Get("id").(string)

	client := meta.(*mackerel.Client)

	channels, err := client.FindChannels()
	if err != nil {
		return err
	}

	var channel *mackerel.Channel
	for _, c := range channels {
		if c.ID == id {
			channel = c
			break
		}
	}
	if channel == nil {
		return fmt.Errorf(`the ID '%s' does not match any channel in mackerel.io`, id)
	}
	d.SetId(channel.ID)
	return flattenChannel(channel, d)
}
