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

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-_]{1,62}$`),
					"must include only alphabets, numbers, hyphen and underscore; in addition, it must be within 2 to 63 characters, and it can not begin a hyphen or underscore"),
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
			_ = d.Set("name", service.Name)
			_ = d.Set("memo", service.Memo)
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
