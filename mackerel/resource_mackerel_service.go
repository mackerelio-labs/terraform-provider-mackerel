package mackerel

import (
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
	service, err := client.CreateService(&mackerel.CreateServiceParam{
		Name: d.Get("name").(string),
		Memo: d.Get("memo").(string),
	})
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
	for _, service := range services {
		if service.Name == d.Id() {
			d.Set("name", service.Name)
			d.Set("memo", service.Memo)
			break
		}
	}
	return nil
}

func resourceMackerelServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	_, err := client.DeleteService(d.Id())
	return err
}
