package mackerel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func dataSourceMackerelService() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMackerelServiceRead,

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

func dataSourceMackerelServiceRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	name := d.Get("name").(string)

	client := m.(*mackerel.Client)
	services, err := client.FindServices()
	if err != nil {
		return diag.FromErr(err)
	}

	var service *mackerel.Service
	for _, s := range services {
		if s.Name == name {
			service = s
			break
		}
	}
	if service == nil {
		return diag.Errorf("the name '%s' does not match any service in mackerel.io", name)
	}
	d.SetId(service.Name)
	if err := flattenService(service, d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}
