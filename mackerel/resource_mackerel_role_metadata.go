package mackerel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelRoleMetadata() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelRoleMetadataCreate,
		Read:   resourceMackerelRoleMetadataRead,
		Update: resourceMackerelRoleMetadataUpdate,
		Delete: resourceMackerelRoleMetadataDelete,

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
			},
			"metadata": {
				Type:     schema.TypeMap,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceMackerelRoleMetadataCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	if err := client.PutRoleMetaData(
		d.Get("service").(string),
		d.Get("role").(string),
		d.Get("namespace").(string),
		d.Get("metadata").(mackerel.RoleMetaData),
	); err != nil {
		return err
	}
	d.SetId(d.Get("namespace").(string))
	return resourceMackerelRoleMetadataRead(d, meta)
}

func resourceMackerelRoleMetadataRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	resp, err := client.GetRoleMetaData(
		d.Get("service").(string),
		d.Get("role").(string),
		d.Get("namespace").(string),
	)
	if err != nil {
		return err
	}
	metadata := make(map[string]interface{})
	for k, v := range resp.RoleMetaData.(map[string]interface{}) {
		metadata[k] = v
	}
	_ = d.Set("metadata", metadata)
	return nil
}

func resourceMackerelRoleMetadataUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceMackerelRoleMetadataCreate(d, meta)
}

func resourceMackerelRoleMetadataDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	if err := client.DeleteRoleMetaData(
		d.Get("service").(string),
		d.Get("role").(string),
		d.Get("namespace").(string),
	); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
