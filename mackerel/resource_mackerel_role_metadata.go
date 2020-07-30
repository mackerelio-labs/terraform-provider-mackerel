package mackerel

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelRoleMetadata() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelRoleMetadataCreate,
		Read:   resourceMackerelRoleMetadataRead,
		Update: resourceMackerelRoleMetadataUpdate,
		Delete: resourceMackerelRoleMetadataDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMackerelRoleMetadataImport,
		},
		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role": {
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

func resourceMackerelRoleMetadataCreate(d *schema.ResourceData, meta interface{}) error {
	service := d.Get("service").(string)
	role := d.Get("role").(string)
	namespace := d.Get("namespace").(string)
	metadata, err := expandRoleMetadata(d.Get("metadata_json").(string))
	if err != nil {
		return err
	}
	client := meta.(*mackerel.Client)
	if err := client.PutRoleMetaData(service, role, namespace, metadata); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s:%s/%s", service, role, namespace))
	return resourceMackerelRoleMetadataRead(d, meta)
}

func resourceMackerelRoleMetadataRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	resp, err := client.GetRoleMetaData(d.Get("service").(string), d.Get("role").(string), d.Get("namespace").(string))
	if err != nil {
		return err
	}
	return flattenRoleMetadata(resp.RoleMetaData, d)
}

func resourceMackerelRoleMetadataUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceMackerelRoleMetadataCreate(d, meta)
}

func resourceMackerelRoleMetadataDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	return client.DeleteRoleMetaData(d.Get("service").(string), d.Get("role").(string), d.Get("namespace").(string))
}

func resourceMackerelRoleMetadataImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	r := regexp.MustCompile(`^([a-zA-Z0-9-_]+):([a-zA-Z0-9-_]+)/(.*)$`)
	idParts := r.FindStringSubmatch(d.Id())
	if idParts == nil || idParts[1] == "" || idParts[2] == "" || idParts[3] == "" {
		return nil, fmt.Errorf("the ID must be in the form '<service name>:<role name>/<namespace>'")
	}
	d.Set("service", idParts[1])
	d.Set("role", idParts[2])
	d.Set("namespace", idParts[3])

	return []*schema.ResourceData{d}, nil
}

func expandRoleMetadata(jsonString string) (mackerel.RoleMetaData, error) {
	var metadata mackerel.RoleMetaData
	err := json.Unmarshal([]byte(jsonString), &metadata)
	return metadata, err
}

func flattenRoleMetadata(metadata mackerel.RoleMetaData, d *schema.ResourceData) error {
	metadataJSON, err := structure.FlattenJsonToString(metadata.(map[string]interface{}))
	if err != nil {
		return err
	}
	d.Set("metadata_json", metadataJSON)
	return nil
}
