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
	role, err := client.CreateRole(service, &mackerel.CreateRoleParam{
		Name: d.Get("name").(string),
		Memo: d.Get("memo").(string),
	})
	if err != nil {
		return err
	}
	d.SetId(makeRoleID(service, role.Name))
	return resourceMackerelRoleRead(d, meta)
}

func resourceMackerelRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	roles, err := client.FindRoles(d.Get("service").(string))
	if err != nil {
		return err
	}
	for _, role := range roles {
		if role.Name == d.Get("name").(string) {
			d.Set("memo", role.Memo)
			break
		}
	}
	return nil
}

func resourceMackerelRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	_, err := client.DeleteRole(d.Get("service").(string), d.Get("name").(string))
	return err
}

func resourceMackerelRoleImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), "/roles/") {
		s := strings.Split(d.Id(), "/roles/")
		d.Set("service", s[0])
		d.Set("name", s[1])
	}

	return []*schema.ResourceData{d}, nil
}

func makeRoleID(service, name string) string {
	return fmt.Sprintf("%s/roles/%s", service, name)
}
