package mackerel

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func dataSourceMackerelServiceMetadata() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMackerelServiceMetadataRead,
		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metadata_json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMackerelServiceMetadataRead(d *schema.ResourceData, meta interface{}) error {
	service := d.Get("service").(string)
	namespace := d.Get("namespace").(string)

	client := meta.(*mackerel.Client)
	resp, err := client.GetServiceMetaData(service, namespace)
	if err != nil {
		return err
	}
	d.SetId(strings.Join([]string{service, namespace}, "/"))
	return flattenServiceMetadata(resp.ServiceMetaData, d)
}
