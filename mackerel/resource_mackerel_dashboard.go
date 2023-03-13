package mackerel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

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
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Required: true,
						},
						"host": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
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
							MaxItems: 1,
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
							MaxItems: 1,
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
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"expression": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"range": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     dashboardRangeResource,
						},
						"layout": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem:     dashboardLayoutResource,
						},
					},
				},
			},
			"value": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Required: true,
						},
						"metric": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem:     dashboardMetricResource,
						},
						"fraction_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"suffix": {
							Type:     schema.TypeString,
							Required: true,
						},
						"layout": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem:     dashboardLayoutResource,
						},
					},
				},
			},
			"markdown": {
				Type:     schema.TypeSet,
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
							MaxItems: 1,
							Elem:     dashboardLayoutResource,
						},
					},
				},
			},
			"alert_status": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role_fullname": {
							Type:     schema.TypeString,
							Required: true,
						},
						"layout": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
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
		graphs := d.Get("graph").(*schema.Set)
		for _, graph := range graphs.List() {
			g := graph.(map[string]interface{})
			var r mackerel.Range
			if v, ok := g["range"].([]interface{}); ok && len(v) > 0 {
				r = expandDashboardRange(v)
			}
			if len(g["host"].([]interface{})) > 0 {
				widgets = append(widgets, mackerel.Widget{
					Type:   "graph",
					Title:  g["title"].(string),
					Graph:  expandDashboardGraphHost(g["host"].([]interface{})),
					Range:  r,
					Layout: expandDashboardLayout(g["layout"].([]interface{})[0].(map[string]interface{})),
				})
			}
			if len(g["role"].([]interface{})) > 0 {
				widgets = append(widgets, mackerel.Widget{
					Type:   "graph",
					Title:  g["title"].(string),
					Graph:  expandDashboardGraphRole(g["role"].([]interface{})),
					Range:  r,
					Layout: expandDashboardLayout(g["layout"].([]interface{})[0].(map[string]interface{})),
				})
			}
			if len(g["service"].([]interface{})) > 0 {
				widgets = append(widgets, mackerel.Widget{
					Type:   "graph",
					Title:  g["title"].(string),
					Graph:  expandDashboardGraphService(g["service"].([]interface{})),
					Range:  r,
					Layout: expandDashboardLayout(g["layout"].([]interface{})[0].(map[string]interface{})),
				})
			}
			if len(g["expression"].([]interface{})) > 0 {
				widgets = append(widgets, mackerel.Widget{
					Type:   "graph",
					Title:  g["title"].(string),
					Graph:  expandDashboardGraphExpression(g["expression"].([]interface{})),
					Range:  r,
					Layout: expandDashboardLayout(g["layout"].([]interface{})[0].(map[string]interface{})),
				})
			}
		}
	}
	if _, ok := d.GetOk("value"); ok {
		values := d.Get("value").(*schema.Set)
		for _, value := range values.List() {
			v := value.(map[string]interface{})
			host := v["metric"].([]interface{})[0].(map[string]interface{})["host"].([]interface{})
			service := v["metric"].([]interface{})[0].(map[string]interface{})["service"].([]interface{})
			expression := v["metric"].([]interface{})[0].(map[string]interface{})["expression"].([]interface{})
			if len(host) > 0 {
				widgets = append(widgets, mackerel.Widget{
					Type:         "value",
					Title:        v["title"].(string),
					Metric:       expandDashboardValueHost(host),
					FractionSize: pointer(int64(v["fraction_size"].(int))),
					Suffix:       v["suffix"].(string),
					Layout:       expandDashboardLayout(v["layout"].([]interface{})[0].(map[string]interface{})),
				})
			}
			if len(service) > 0 {
				widgets = append(widgets, mackerel.Widget{
					Type:         "value",
					Title:        v["title"].(string),
					Metric:       expandDashboardValueService(service),
					FractionSize: pointer(int64(v["fraction_size"].(int))),
					Suffix:       v["suffix"].(string),
					Layout:       expandDashboardLayout(v["layout"].([]interface{})[0].(map[string]interface{})),
				})
			}
			if len(expression) > 0 {
				widgets = append(widgets, mackerel.Widget{
					Type:         "value",
					Title:        v["title"].(string),
					Metric:       expandDashboardValueExpression(expression),
					FractionSize: pointer(int64(v["fraction_size"].(int))),
					Suffix:       v["suffix"].(string),
					Layout:       expandDashboardLayout(v["layout"].([]interface{})[0].(map[string]interface{})),
				})
			}
		}
	}
	if _, ok := d.GetOk("markdown"); ok {
		markdowns := d.Get("markdown").(*schema.Set)
		for _, markdown := range markdowns.List() {
			m := markdown.(map[string]interface{})
			widgets = append(widgets, mackerel.Widget{
				Type:     "markdown",
				Title:    m["title"].(string),
				Markdown: m["markdown"].(string),
				Layout:   expandDashboardLayout(m["layout"].([]interface{})[0].(map[string]interface{})),
			})
		}
	}
	if _, ok := d.GetOk("alert_status"); ok {
		alert_statuses := d.Get("alert_status").(*schema.Set)
		for _, alert_status := range alert_statuses.List() {
			a := alert_status.(map[string]interface{})
			widgets = append(widgets, mackerel.Widget{
				Type:         "alertStatus",
				Title:        a["title"].(string),
				RoleFullName: a["role_fullname"].(string),
				Layout:       expandDashboardLayout(a["layout"].([]interface{})[0].(map[string]interface{})),
			})
		}
	}

	return widgets
}

func pointer(x int64) *int64 {
	return &x
}

func expandDashboardGraphHost(host []interface{}) mackerel.Graph {
	return mackerel.Graph{
		Type:   "host",
		HostID: host[0].(map[string]interface{})["host_id"].(string),
		Name:   host[0].(map[string]interface{})["name"].(string),
	}
}

func expandDashboardGraphRole(role []interface{}) mackerel.Graph {
	return mackerel.Graph{
		Type:         "role",
		RoleFullName: role[0].(map[string]interface{})["role_fullname"].(string),
		Name:         role[0].(map[string]interface{})["name"].(string),
		IsStacked:    role[0].(map[string]interface{})["is_stacked"].(bool),
	}
}

func expandDashboardGraphService(service []interface{}) mackerel.Graph {
	return mackerel.Graph{
		Type:        "service",
		ServiceName: service[0].(map[string]interface{})["service_name"].(string),
		Name:        service[0].(map[string]interface{})["name"].(string),
	}
}

func expandDashboardGraphExpression(expression []interface{}) mackerel.Graph {
	return mackerel.Graph{
		Type:       "expression",
		Expression: expression[0].(map[string]interface{})["expression"].(string),
	}
}

func expandDashboardValueHost(host []interface{}) mackerel.Metric {
	return mackerel.Metric{
		Type:   "host",
		HostID: host[0].(map[string]interface{})["host_id"].(string),
		Name:   host[0].(map[string]interface{})["name"].(string),
	}
}

func expandDashboardValueService(service []interface{}) mackerel.Metric {
	return mackerel.Metric{
		Type:        "service",
		ServiceName: service[0].(map[string]interface{})["service_name"].(string),
		Name:        service[0].(map[string]interface{})["name"].(string),
	}
}

func expandDashboardValueExpression(expression []interface{}) mackerel.Metric {
	return mackerel.Metric{
		Type:       "expression",
		Expression: expression[0].(map[string]interface{})["expression"].(string),
	}
}

func expandDashboardRange(r []interface{}) mackerel.Range {
	if len(r[0].(map[string]interface{})["relative"].([]interface{})) > 0 {
		return mackerel.Range{
			Type:   "relative",
			Period: int64(r[0].(map[string]interface{})["relative"].([]interface{})[0].(map[string]interface{})["period"].(int)),
			Offset: int64(r[0].(map[string]interface{})["relative"].([]interface{})[0].(map[string]interface{})["offset"].(int)),
		}
	}
	if len(r[0].(map[string]interface{})["absolute"].([]interface{})) > 0 {
		return mackerel.Range{
			Type:  "absolute",
			Start: int64(r[0].(map[string]interface{})["absolute"].([]interface{})[0].(map[string]interface{})["start"].(int)),
			End:   int64(r[0].(map[string]interface{})["absolute"].([]interface{})[0].(map[string]interface{})["end"].(int)),
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
