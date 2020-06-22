package mackerel

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				ForceNew: true,
			},
			"metadata_json": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
			},
		},
	}
}

func resourceMackerelServiceMetadataCreate(d *schema.ResourceData, meta interface{}) error {
	var metadata mackerel.ServiceMetaData
	if err := json.Unmarshal([]byte(d.Get("metadata_json").(string)), &metadata); err != nil {
		return err
	}
	namespace := d.Get("namespace").(string)

	client := meta.(*mackerel.Client)
	if err := client.PutServiceMetaData(d.Get("service").(string), namespace, metadata); err != nil {
		return err
	}
	d.SetId(namespace)

	return resourceMackerelServiceMetadataRead(d, meta)
}

func resourceMackerelServiceMetadataRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	resp, err := client.GetServiceMetaData(d.Get("service").(string), d.Id())
	if err != nil {
		return err
	}

	metadataJsonBytes, err := json.Marshal(resp.ServiceMetaData)
	if err != nil {
		return err
	}

	metadataJson, err := structure.NormalizeJsonString(string(metadataJsonBytes))
	if err != nil {
		return err
	}

	if err := d.Set("metadata_json", metadataJson); err != nil {
		return err
	}

	return nil
}

func resourceMackerelServiceMetadataUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceMackerelServiceMetadataCreate(d, meta)
}

func resourceMackerelServiceMetadataDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	return client.DeleteServiceMetaData(d.Get("service").(string), d.Id())
}
