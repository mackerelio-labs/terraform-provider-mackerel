package mackerel

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func dataSourceMackerelRoleMetadata() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMackerelRoleMetadataRead,

		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role": {
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

func dataSourceMackerelRoleMetadataRead(d *schema.ResourceData, meta interface{}) error {
	service := d.Get("service").(string)
	role := d.Get("role").(string)
	namespace := d.Get("namespace").(string)

	client := meta.(*mackerel.Client)
	resp, err := client.GetRoleMetaData(service, role, namespace)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s:%s/%s", service, role, namespace))
	return flattenRoleMetadata(resp.RoleMetaData, d)
}
