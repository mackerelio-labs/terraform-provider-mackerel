package mackerel

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func dataSourceMackerelService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMackerelServiceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"memo": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMackerelServiceRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	client := meta.(*mackerel.Client)
	services, err := client.FindServices()
	if err != nil {
		return err
	}

	var service *mackerel.Service
	for _, s := range services {
		if s.Name == name {
			service = s
			break
		}
	}
	if service == nil {
		return fmt.Errorf("the name '%s' does not match any service in mackerel.io", name)
	}
	d.SetId(service.Name)
	return flattenService(service, d)
}
