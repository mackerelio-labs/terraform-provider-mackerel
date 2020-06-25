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
	var metadata mackerel.RoleMetaData
	if err := json.Unmarshal([]byte(d.Get("metadata_json").(string)), &metadata); err != nil {
		return err
	}

	client := meta.(*mackerel.Client)
	if err := client.PutRoleMetaData(service, role, namespace, metadata); err != nil {
		return err
	}
	d.SetId(makeRoleMetadataID(service, role, namespace))
	return resourceMackerelRoleMetadataRead(d, meta)
}

func resourceMackerelRoleMetadataRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	resp, err := client.GetRoleMetaData(d.Get("service").(string), d.Get("role").(string), d.Get("namespace").(string))
	if err != nil {
		return err
	}

	metadataJSONBytes, err := json.Marshal(resp.RoleMetaData)
	if err != nil {
		return err
	}

	metadataJSON, err := structure.NormalizeJsonString(string(metadataJSONBytes))
	if err != nil {
		return err
	}

	if err := d.Set("metadata_json", metadataJSON); err != nil {
		return err
	}

	return nil
}

func resourceMackerelRoleMetadataUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceMackerelRoleMetadataCreate(d, meta)
}

func resourceMackerelRoleMetadataDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	return client.DeleteRoleMetaData(d.Get("service").(string), d.Get("role").(string), d.Get("namespace").(string))
}

func resourceMackerelRoleMetadataImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	r := regexp.MustCompile(`^([a-zA-Z0-9-_]+)/roles/([a-zA-Z0-9-_]+)/metadata/(.*)$`)
	if v := r.FindStringSubmatch(d.Id()); v != nil {
		if err := d.Set("service", v[1]); err != nil {
			return nil, err
		}
		if err := d.Set("role", v[2]); err != nil {
			return nil, err
		}
		if err := d.Set("namespace", v[3]); err != nil {
			return nil, err
		}
	}
	return []*schema.ResourceData{d}, nil
}

func makeRoleMetadataID(service, role, namespace string) string {
	return fmt.Sprintf("%s/roles/%s/metadata/%s", service, role, namespace)
}
