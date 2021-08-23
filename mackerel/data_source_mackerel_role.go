package mackerel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func dataSourceMackerelRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMackerelRoleRead,

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

func dataSourceMackerelRoleRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	service := d.Get("service").(string)
	name := d.Get("name").(string)

	client := m.(*mackerel.Client)
	roles, err := client.FindRoles(service)
	if err != nil {
		return diag.FromErr(err)
	}

	var role *mackerel.Role
	for _, r := range roles {
		if r.Name == name {
			role = r
			break
		}
	}
	if role == nil {
		return diag.Errorf("the name '%s' does not match any role in mackerel.io", name)
	}
	d.SetId(fmt.Sprintf("%s:%s", service, role.Name))
	return flattenRole(role, d)
}
