package mackerel

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func dataSourceMackerelServiceMetadata() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMackerelServiceMetadataRead,
		Schema: map[string]*schema.Schema{
			"service": {
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

func dataSourceMackerelServiceMetadataRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	service := d.Get("service").(string)
	namespace := d.Get("namespace").(string)

	client := m.(*mackerel.Client)
	resp, err := client.GetServiceMetaData(service, namespace)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strings.Join([]string{service, namespace}, "/"))
	if err := flattenServiceMetadata(resp.ServiceMetaData, d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}
