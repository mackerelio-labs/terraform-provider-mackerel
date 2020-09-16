package mackerel

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMackerelServiceCreate,
		ReadContext:   resourceMackerelServiceRead,
		DeleteContext: resourceMackerelServiceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
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

func resourceMackerelServiceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mackerel.Client)
	service, err := client.CreateService(expandCreateServiceParam(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(service.Name)
	return resourceMackerelServiceRead(ctx, d, m)
}

func resourceMackerelServiceRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mackerel.Client)
	services, err := client.FindServices()
	if err != nil {
		return diag.FromErr(err)
	}

	var service *mackerel.Service
	for _, s := range services {
		if s.Name == d.Id() {
			service = s
			break
		}
	}
	if service == nil {
		return diag.Errorf("the name '%s' does not match any service in mackerel.io", d.Id())
	}
	return flattenService(service, d)
}

func resourceMackerelServiceDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*mackerel.Client)
	_, err := client.DeleteService(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func expandCreateServiceParam(d *schema.ResourceData) *mackerel.CreateServiceParam {
	return &mackerel.CreateServiceParam{
		Name: d.Get("name").(string),
		Memo: d.Get("memo").(string),
	}
}

// func flattenService(service *mackerel.Service, d *schema.ResourceData) error {
// 	d.Set("name", service.Name)
// 	d.Set("memo", service.Memo)
// 	return nil
// }
