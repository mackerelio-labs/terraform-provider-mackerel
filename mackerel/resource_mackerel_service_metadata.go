package mackerel

import (
	"encoding/json"
	"fmt"
	"strings"

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
		Importer: &schema.ResourceImporter{
			State: resourceMackerelServiceMetadataImport,
		},
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
	service := d.Get("service").(string)
	namespace := d.Get("namespace").(string)
	metadata, err := expandServiceMetadata(d)
	if err != nil {
		return err
	}
	client := meta.(*mackerel.Client)
	if err := client.PutServiceMetaData(service, namespace, metadata); err != nil {
		return err
	}
	d.SetId(strings.Join([]string{service, namespace}, "/"))
	return resourceMackerelServiceMetadataRead(d, meta)
}

func resourceMackerelServiceMetadataRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	resp, err := client.GetServiceMetaData(d.Get("service").(string), d.Get("namespace").(string))
	if err != nil {
		return err
	}
	return flattenServiceMetadata(resp.ServiceMetaData, d)
}

func resourceMackerelServiceMetadataUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceMackerelServiceMetadataCreate(d, meta)
}

func resourceMackerelServiceMetadataDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	return client.DeleteServiceMetaData(d.Get("service").(string), d.Get("namespace").(string))
}

func resourceMackerelServiceMetadataImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	idParts := strings.SplitN(d.Id(), "/", 2)
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		return nil, fmt.Errorf("the ID must be in the form '<service name>:<namespace>'")
	}
	d.Set("service", idParts[0])
	d.Set("namespace", idParts[1])

	return []*schema.ResourceData{d}, nil
}

func expandServiceMetadata(d *schema.ResourceData) (mackerel.ServiceMetaData, error) {
	var metadata mackerel.ServiceMetaData
	if err := json.Unmarshal([]byte(d.Get("metadata_json").(string)), &metadata); err != nil {
		return nil, err
	}
	return metadata, nil
}

func flattenServiceMetadata(metadata mackerel.ServiceMetaData, d *schema.ResourceData) error {
	metadataJSON, err := structure.FlattenJsonToString(metadata.(map[string]interface{}))
	if err != nil {
		return err
	}
	d.Set("metadata_json", metadataJSON)
	return nil
}
