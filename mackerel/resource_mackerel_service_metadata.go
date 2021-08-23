package mackerel

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelServiceMetadata() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMackerelServiceMetadataCreate,
		ReadContext:   resourceMackerelServiceMetadataRead,
		UpdateContext: resourceMackerelServiceMetadataUpdate,
		DeleteContext: resourceMackerelServiceMetadataDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMackerelServiceMetadataImport,
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

func resourceMackerelServiceMetadataCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	service := d.Get("service").(string)
	namespace := d.Get("namespace").(string)
	metadata, err := expandServiceMetadata(d.Get("metadata_json").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	client := m.(*mackerel.Client)
	if err := client.PutServiceMetaData(service, namespace, metadata); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strings.Join([]string{service, namespace}, "/"))
	return resourceMackerelServiceMetadataRead(ctx, d, m)
}

func resourceMackerelServiceMetadataRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mackerel.Client)
	resp, err := client.GetServiceMetaData(d.Get("service").(string), d.Get("namespace").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	return flattenServiceMetadata(resp.ServiceMetaData, d)
}

func resourceMackerelServiceMetadataUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMackerelServiceMetadataCreate(ctx, d, m)
}

func resourceMackerelServiceMetadataDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*mackerel.Client)
	if err := client.DeleteServiceMetaData(d.Get("service").(string), d.Get("namespace").(string)); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceMackerelServiceMetadataImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	idParts := strings.SplitN(d.Id(), "/", 2)
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		return nil, fmt.Errorf("the ID must be in the form '<service name>/<namespace>'")
	}
	d.Set("service", idParts[0])
	d.Set("namespace", idParts[1])

	return []*schema.ResourceData{d}, nil
}

func expandServiceMetadata(jsonString string) (mackerel.ServiceMetaData, error) {
	var metadata mackerel.ServiceMetaData
	err := json.Unmarshal([]byte(jsonString), &metadata)
	return metadata, err
}
