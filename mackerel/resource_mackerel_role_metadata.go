package mackerel

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelRoleMetadata() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMackerelRoleMetadataCreate,
		ReadContext:   resourceMackerelRoleMetadataRead,
		UpdateContext: resourceMackerelRoleMetadataUpdate,
		DeleteContext: resourceMackerelRoleMetadataDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMackerelRoleMetadataImport,
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

func resourceMackerelRoleMetadataCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	service := d.Get("service").(string)
	role := d.Get("role").(string)
	namespace := d.Get("namespace").(string)
	metadata, err := expandRoleMetadata(d.Get("metadata_json").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	client := m.(*mackerel.Client)
	if err := client.PutRoleMetaData(service, role, namespace, metadata); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s:%s/%s", service, role, namespace))
	return resourceMackerelRoleMetadataRead(ctx, d, m)
}

func resourceMackerelRoleMetadataRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mackerel.Client)
	resp, err := client.GetRoleMetaData(d.Get("service").(string), d.Get("role").(string), d.Get("namespace").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	return flattenRoleMetadata(resp.RoleMetaData, d)
}

func resourceMackerelRoleMetadataUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMackerelRoleMetadataCreate(ctx, d, m)
}

func resourceMackerelRoleMetadataDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*mackerel.Client)
	if err := client.DeleteRoleMetaData(d.Get("service").(string), d.Get("role").(string), d.Get("namespace").(string)); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceMackerelRoleMetadataImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
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
