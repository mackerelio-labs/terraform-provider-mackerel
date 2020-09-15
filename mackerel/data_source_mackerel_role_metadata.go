package mackerel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func dataSourceMackerelRoleMetadata() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMackerelRoleMetadataRead,

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
			"metadata_json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMackerelRoleMetadataRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	service := d.Get("service").(string)
	role := d.Get("role").(string)
	namespace := d.Get("namespace").(string)

	client := m.(*mackerel.Client)
	resp, err := client.GetRoleMetaData(service, role, namespace)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s:%s/%s", service, role, namespace))
	if err := flattenRoleMetadata(resp.RoleMetaData, d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}
