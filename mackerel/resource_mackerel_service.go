package mackerel

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelService() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelServiceCreate,
		Read:   resourceMackerelServiceRead,
		Delete: resourceMackerelServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceMackerelServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	service, err := client.CreateService(expandCreateServiceParam(d))
	if err != nil {
		return err
	}
	d.SetId(service.Name)
	return resourceMackerelServiceRead(d, meta)
}

func resourceMackerelServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	services, err := client.FindServices()
	if err != nil {
		return err
	}

	var service *mackerel.Service
	for _, s := range services {
		if s.Name == d.Id() {
			service = s
			break
		}
	}
	if service == nil {
		return fmt.Errorf("the name '%s' does not match any service in mackerel.io", d.Id())
	}
	return flattenService(service, d)
}

func resourceMackerelServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	_, err := client.DeleteService(d.Id())
	return err
}

func expandCreateServiceParam(d *schema.ResourceData) *mackerel.CreateServiceParam {
	return &mackerel.CreateServiceParam{
		Name: d.Get("name").(string),
		Memo: d.Get("memo").(string),
	}
}

func flattenService(service *mackerel.Service, d *schema.ResourceData) error {
	d.Set("name", service.Name)
	d.Set("memo", service.Memo)
	return nil
}
