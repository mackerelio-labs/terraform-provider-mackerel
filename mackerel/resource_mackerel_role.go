package mackerel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelRoleCreate,
		Read:   resourceMackerelRoleRead,
		Delete: resourceMackerelRoleDelete,

		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"memo": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceMackerelRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	role, err := client.CreateRole(d.Get("service").(string), &mackerel.CreateRoleParam{
		Name: d.Get("name").(string),
		Memo: d.Get("memo").(string),
	})
	if err != nil {
		return err
	}
	d.SetId(role.Name)
	return resourceMackerelRoleRead(d, meta)
}

func resourceMackerelRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	roles, err := client.FindRoles(d.Get("service").(string))
	if err != nil {
		return err
	}
	for _, role := range roles {
		if role.Name == d.Id() {
			_ = d.Set("name", role.Name)
			_ = d.Set("memo", role.Memo)
			break
		}
	}
	return nil
}

func resourceMackerelRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	_, err := client.DeleteRole(d.Get("service").(string), d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
