package mackerel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelServiceMetadata() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelServiceMetadataCreate,
		Read:   resourceMackerelServiceMetadataRead,
		Update: resourceMackerelServiceMetadataUpdate,
		Delete: resourceMackerelServiceMetadataDelete,

		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceMackerelServiceMetadataCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	if err := client.PutServiceMetaData(
		d.Get("service").(string),
		d.Get("namespace").(string),
		d.Get("metadata").(mackerel.ServiceMetaData),
	); err != nil {
		return err
	}
	d.SetId(d.Get("namespace").(string))
	return resourceMackerelServiceMetadataRead(d, meta)
}

func resourceMackerelServiceMetadataRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	resp, err := client.GetServiceMetaData(d.Get("service").(string), d.Get("namespace").(string))
	if err != nil {
		return err
	}
	metadata := make(map[string]interface{})
	for k, v := range resp.ServiceMetaData.(map[string]interface{}) {
		metadata[k] = v
	}
	_ = d.Set("metadata", metadata)
	return nil
}

func resourceMackerelServiceMetadataUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceMackerelServiceMetadataCreate(d, meta)
}

func resourceMackerelServiceMetadataDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	err := client.DeleteServiceMetaData(d.Get("service").(string), d.Get("namespace").(string))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
