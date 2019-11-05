package mackerel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
