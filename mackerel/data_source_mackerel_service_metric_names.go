package mackerel

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func dataSourceMackerelServiceMetricNames() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMackerelServiceMetricNamesRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metric_names": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceMackerelServiceMetricNamesRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	prefix := d.Get("prefix").(string)

	client := m.(*mackerel.Client)
	names, err := client.ListServiceMetricNames(name)
	if err != nil {
		return diag.FromErr(err)
	}

	metricNames := make([]string, 0, len(names))
	for _, n := range names {
		if strings.HasPrefix(n, prefix) {
			metricNames = append(metricNames, n)
		}
	}

	d.SetId(name + ":" + prefix)
	return flattenServiceMetricNames(name, metricNames, d)
}
