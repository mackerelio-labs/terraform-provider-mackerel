package mackerel

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMackerelRoleCreate,
		ReadContext:   resourceMackerelRoleRead,
		DeleteContext: resourceMackerelRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMackerelRoleImport,
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

func resourceMackerelRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	service := d.Get("service").(string)
	client := m.(*mackerel.Client)
	role, err := client.CreateRole(service, expandCreateRoleParam(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s:%s", service, role.Name))
	return resourceMackerelRoleRead(ctx, d, m)
}

func resourceMackerelRoleRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*mackerel.Client)
	roles, err := client.FindRoles(d.Get("service").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	var role *mackerel.Role
	for _, r := range roles {
		if r.Name == d.Get("name").(string) {
			role = r
			break
		}
	}
	if role == nil {
		return diag.Errorf("the name '%s' does not match any role in mackerel.io", d.Get("name").(string))
	}
	if err := flattenRole(role, d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceMackerelRoleDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*mackerel.Client)
	_, err := client.DeleteRole(d.Get("service").(string), d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceMackerelRoleImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
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
