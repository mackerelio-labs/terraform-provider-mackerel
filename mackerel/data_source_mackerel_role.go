package mackerel

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func dataSourceMackerelRole() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMackerelRoleRead,

		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Required: true,
			},
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

func dataSourceMackerelRoleRead(d *schema.ResourceData, meta interface{}) error {
	service := d.Get("service").(string)
	name := d.Get("name").(string)

	client := meta.(*mackerel.Client)
	roles, err := client.FindRoles(service)
	if err != nil {
		return err
	}

	var role *mackerel.Role
	for _, r := range roles {
		if r.Name == name {
			role = r
			break
		}
	}
	if role == nil {
		return fmt.Errorf("the name '%s' does not match any role in mackerel.io", name)
	}
	d.SetId(fmt.Sprintf("%s:%s", service, role.Name))
	return flattenRole(role, d)
}
