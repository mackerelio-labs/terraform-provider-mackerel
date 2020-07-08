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
	var metadata mackerel.ServiceMetaData
	if err := json.Unmarshal([]byte(d.Get("metadata_json").(string)), &metadata); err != nil {
		return err
	}

	client := meta.(*mackerel.Client)
	if err := client.PutServiceMetaData(service, namespace, metadata); err != nil {
		return err
	}
	d.SetId(makeServiceMetadataID(service, namespace))

	return resourceMackerelServiceMetadataRead(d, meta)
}

func resourceMackerelServiceMetadataRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	resp, err := client.GetServiceMetaData(d.Get("service").(string), d.Get("namespace").(string))
	if err != nil {
		return err
	}

	metadataJson, err := structure.FlattenJsonToString(resp.ServiceMetaData.(map[string]interface{}))
	if err != nil {
		return err
	}
	d.Set("metadata_json", metadataJson)

	return nil
}

func resourceMackerelServiceMetadataUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceMackerelServiceMetadataCreate(d, meta)
}

func resourceMackerelServiceMetadataDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	return client.DeleteServiceMetaData(d.Get("service").(string), d.Get("namespace").(string))
}

func resourceMackerelServiceMetadataImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), "/metadata/") {
		s := strings.Split(d.Id(), "/metadata/")
		d.Set("service", s[0])
		d.Set("namespace", s[1])
	}

	return []*schema.ResourceData{d}, nil
}

func makeServiceMetadataID(service, namespace string) string {
	return fmt.Sprintf("%s/metadata/%s", service, namespace)
}
