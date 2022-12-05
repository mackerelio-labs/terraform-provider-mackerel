package mackerel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

var dashboardGraphDataResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"host": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"host_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"role": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"role_fullname": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"is_stacked": {
						Type:     schema.TypeBool,
						Computed: true,
					},
				},
			},
		},
		"service": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"service_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"expression": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"expression": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	},
}

var dashboardRangeDataResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"relative": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"period": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"offset": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"absolute": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"start": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"end": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
	},
}

var dashboardLayoutDataResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"x": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"y": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"width": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"height": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	},
}

var dashboardMetricDataResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"host": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"host_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"service": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"service_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"expression": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"expression": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	},
}

func dataSourceMackerelDashboard() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMackerelDashboardRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"memo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"graph": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"graph": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dashboardGraphDataResource,
						},
						"range": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dashboardRangeDataResource,
						},
						"layout": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dashboardLayoutDataResource,
						},
					},
				},
			},
			"value": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metric": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     dashboardMetricDataResource,
						},
						"fraction_size": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"suffix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"layout": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     dashboardLayoutDataResource,
						},
					},
				},
			},
			"markdown": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"markdown": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"layout": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dashboardLayoutDataResource,
						},
					},
				},
			},
			"alert_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"roll_fullname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"layout": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     dashboardLayoutDataResource,
						},
					},
				},
			},
		},
	}
}

func dataSourceMackerelDashboardRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Get("id").(string)

	client := m.(*mackerel.Client)

	dashboards, err := client.FindDashboards()

	if err != nil {
		return diag.FromErr(err)
	}
	var dashboard *mackerel.Dashboard
	for _, a := range dashboards {
		if a.ID == id {
			dashboard = a
			break
		}
	}
	if dashboard == nil {
		return diag.Errorf(`the id '%s' does not match any dashboard in mackerel.io`, id)
	}

	dashboardWithWidgets, err := client.FindDashboard(id)
	dashboard.Widgets = dashboardWithWidgets.Widgets

	d.SetId(dashboard.ID)

	return flattenDashboard(dashboard, d)
}
