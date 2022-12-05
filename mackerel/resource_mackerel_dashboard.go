package mackerel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

var dashboardGraphResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"host": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"host_id": {
						Type:     schema.TypeString,
						Required: true,
					},
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"role": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"role_fullname": {
						Type:     schema.TypeString,
						Required: true,
					},
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"is_stacked": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},
		"service": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"service_name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"expression": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"expression": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
	},
}

var dashboardRangeResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"relative": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"period": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"offset": {
						Type:     schema.TypeInt,
						Required: true,
					},
				},
			},
		},
		"absolute": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"start": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"end": {
						Type:     schema.TypeInt,
						Required: true,
					},
				},
			},
		},
	},
}

var dashboardLayoutResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"x": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"y": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"width": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"height": {
			Type:     schema.TypeInt,
			Required: true,
		},
	},
}

var dashboardMetricResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"host": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"host_id": {
						Type:     schema.TypeString,
						Required: true,
					},
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"service": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"service_name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"expression": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"expression": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
	},
}

func resourceMackerelDashboard() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMackerelDashboardCreate,
		ReadContext:   resourceMackerelDashboardRead,
		UpdateContext: resourceMackerelDashboardUpdate,
		DeleteContext: resourceMackerelDashboardDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"memo": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"url_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"graph": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Required: true,
						},
						"graph": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     dashboardGraphResource,
						},
						"range": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     dashboardRangeResource,
						},
						"layout": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     dashboardLayoutResource,
						},
					},
				},
			},
			"value": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Required: true,
						},
						"metric": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     dashboardMetricResource,
						},
						"fraction_size": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"suffix": {
							Type:     schema.TypeString,
							Required: true,
						},
						"layout": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     dashboardLayoutResource,
						},
					},
				},
			},
			"markdown": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Required: true,
						},
						"markdown": {
							Type:     schema.TypeString,
							Required: true,
						},
						"layout": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     dashboardLayoutResource,
						},
					},
				},
			},
			"alert_status": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Required: true,
						},
						"roll_fullname": {
							Type:     schema.TypeString,
							Required: true,
						},
						"layout": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     dashboardLayoutResource,
						},
					},
				},
			},
		},
	}
}

func resourceMackerelDashboardCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mackerel.Client)
	dashboard, err := client.CreateDashboard(expandDashboard(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dashboard.ID)
	return resourceMackerelDashboardRead(ctx, d, m)
}

func resourceMackerelDashboardRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mackerel.Client)
	dashboard, err := client.FindDashboard(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return flattenDashboard(dashboard, d)
}

func resourceMackerelDashboardUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mackerel.Client)
	_, err := client.UpdateDashboard(d.Id(), expandDashboard(d))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceMackerelDashboardRead(ctx, d, m)
}

func resourceMackerelDashboardDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*mackerel.Client)
	_, err := client.DeleteDashboard(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func expandDashboard(d *schema.ResourceData) *mackerel.Dashboard {
	dashboard := &mackerel.Dashboard{
		Title:   d.Get("title").(string),
		Memo:    d.Get("memo").(string),
		URLPath: d.Get("url_path").(string),
		Widgets: expandDashboardWidgets(d),
	}
	return dashboard
}

func expandDashboardWidgets(d *schema.ResourceData) []mackerel.Widget {
	var widgets []mackerel.Widget
	if _, ok := d.GetOk("graph"); ok {
		if _, ok := d.GetOk("graph.0.graph.0.host"); ok {
			widgets = append(widgets, mackerel.Widget{
				Type:  "graph",
				Title: d.Get("graph.0.title").(string),
				Graph: expandDashboardGraphHost(d),
				Range: expandDashboardRange(d),
				//Layout: expandDashboardLayout(d, "graph"),
			})
		}
		if _, ok := d.GetOk("graph.0.graph.0.role"); ok {
			widgets = append(widgets, mackerel.Widget{
				Type:  "graph",
				Title: d.Get("graph.0.title").(string),
				Graph: expandDashboardGraphRole(d),
				Range: expandDashboardRange(d),
				//Layout: expandDashboardLayout(d, "graph"),
			})
		}
		if _, ok := d.GetOk("graph.0.graph.0.service"); ok {
			widgets = append(widgets, mackerel.Widget{
				Type:  "graph",
				Title: d.Get("graph.0.title").(string),
				Graph: expandDashboardGraphService(d),
				Range: expandDashboardRange(d),
				//Layout: expandDashboardLayout(d, "graph"),
			})
		}
	}
	if _, ok := d.GetOk("markdown"); ok {
		markdowns := d.Get("markdown").([]interface{})
		for _, markdown := range markdowns {
			m := markdown.(map[string]interface{})
			widgets = append(widgets, mackerel.Widget{
				Type:     "markdown",
				Title:    m["title"].(string),
				Markdown: m["markdown"].(string),
				Layout:   expandDashboardLayout(m["layout"].([]interface{})[0].(map[string]interface{})),
			})
		}
	}

	return widgets
}

func expandDashboardGraphHost(d *schema.ResourceData) mackerel.Graph {
	return mackerel.Graph{
		Type:   "host",
		HostID: d.Get("graph.0.graph.0.host.0.host_id").(string),
		Name:   d.Get("graph.0.graph.0.host.0.name").(string),
	}
}

func expandDashboardGraphRole(d *schema.ResourceData) mackerel.Graph {
	return mackerel.Graph{
		Type:         "role",
		RoleFullName: d.Get("graph.0.graph.0.role.0.role_fullname").(string),
		Name:         d.Get("graph.0.graph.0.role.0.name").(string),
		IsStacked:    d.Get("graph.0.graph.0.role.0.is_stacked").(bool),
	}
}

func expandDashboardGraphService(d *schema.ResourceData) mackerel.Graph {
	return mackerel.Graph{
		Type:   "service",
		HostID: d.Get("graph.0.graph.0.service.0.service_name").(string),
		Name:   d.Get("graph.0.graph.0.service.0.name").(string),
	}
}

func expandDashboardRange(d *schema.ResourceData) mackerel.Range {
	if _, ok := d.GetOk("graph.0.range.0.relative"); ok {
		return mackerel.Range{
			Type:   "relative",
			Period: int64(d.Get("graph.0.range.0.relative.0.period").(int)),
			Offset: int64(d.Get("graph.0.range.0.relative.0.offset").(int)),
		}
	}
	if _, ok := d.GetOk("graph.0.range.0.absolute"); ok {
		return mackerel.Range{
			Type:  "absolute",
			Start: int64(d.Get("graph.0.range.0.relative.0.start").(int)),
			End:   int64(d.Get("graph.0.range.0.relative.0.end").(int)),
		}
	}
	return mackerel.Range{}
}

func expandDashboardLayout(layout map[string]interface{}) mackerel.Layout {
	return mackerel.Layout{
		X:      int64(layout["x"].(int)),
		Y:      int64(layout["y"].(int)),
		Width:  int64(layout["width"].(int)),
		Height: int64(layout["height"].(int)),
	}
}
