package mackerel

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelRoleCreate,
		Read:   resourceMackerelRoleRead,
		Delete: resourceMackerelRoleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMackerelRoleImport,
		},
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
				ValidateFunc: validation.All(
					validation.StringLenBetween(2, 63),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-_]+$`),
						"must include only alphabets, numbers, hyphen and underscore, and it can not begin a hyphen or underscore"),
				),
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
	service := d.Get("service").(string)
	client := meta.(*mackerel.Client)
	role, err := client.CreateRole(service, expandCreateRoleParam(d))
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s:%s", service, role.Name))
	return resourceMackerelRoleRead(d, meta)
}

func resourceMackerelRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	roles, err := client.FindRoles(d.Get("service").(string))
	if err != nil {
		return err
	}
	var role *mackerel.Role
	for _, r := range roles {
		if r.Name == d.Get("name").(string) {
			role = r
			break
		}
	}
	if role == nil {
		return fmt.Errorf("the name '%s' does not any match role in mackerel.io", d.Get("name").(string))
	}
	return flattenRole(role, d)
}

func resourceMackerelRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	_, err := client.DeleteRole(d.Get("service").(string), d.Get("name").(string))
	return err
}

func resourceMackerelRoleImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	idParts := strings.SplitN(d.Id(), ":", 2)
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		return nil, fmt.Errorf("the ID must be in the form '<service name>:<role name>'")
	}
	d.Set("service", idParts[0])
	d.Set("name", idParts[1])
	return []*schema.ResourceData{d}, nil
}

func expandCreateRoleParam(d *schema.ResourceData) *mackerel.CreateRoleParam {
	return &mackerel.CreateRoleParam{
		Name: d.Get("name").(string),
		Memo: d.Get("memo").(string),
	}
}

func flattenRole(role *mackerel.Role, d *schema.ResourceData) error {
	d.Set("name", role.Name)
	d.Set("memo", role.Memo)
	return nil
}
