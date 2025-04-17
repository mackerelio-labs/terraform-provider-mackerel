package mackerel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
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
		"query": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"query": {
						Type:     schema.TypeString,
						Required: true,
					},
					"legend": {
						Type:     schema.TypeString,
						Optional: true,
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
		CustomizeDiff: customdiff.All(
			customdiff.ValidateValue("graph", func(ctx context.Context, value, meta interface{}) error {
				graphs := value.([]interface{})
				for i, g := range graphs {
					graph := g.(map[string]interface{})
					if ranges, ok := graph["range"].([]interface{}); ok && len(ranges) > 0 {
						if ranges[0] == nil {
							return fmt.Errorf("graph[%d].range: exactly one of 'relative' or 'absolute' must be specified", 0)
						}

						r := ranges[0].(map[string]interface{})
						relativeExists := len(r["relative"].([]interface{})) > 0
						absoluteExists := len(r["absolute"].([]interface{})) > 0

						if !relativeExists && !absoluteExists {
							return fmt.Errorf("graph[%d].range: exactly one of 'relative' or 'absolute' must be specified", i)
						}

						if relativeExists && absoluteExists {
							return fmt.Errorf("graph[%d].range: cannot specify both 'relative' and 'absolute'", i)
						}
					}
				}
				return nil
			}),
		),
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
						"query": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"query": {
										Type:     schema.TypeString,
										Required: true,
									},
									"legend": {
										Type:     schema.TypeString,
										Optional: true,
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
				Type:     schema.TypeList,
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
							MaxItems: 1,
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
		graphs := d.Get("graph").([]interface{})
		for _, graph := range graphs {
			g := graph.(map[string]interface{})
			var r mackerel.Range
			if v, ok := g["range"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
				r = expandDashboardRange(v)
			}
			title := g["title"].(string)
			layout := expandDashboardLayout(g["layout"].([]interface{})[0].(map[string]interface{}))

			if host := g["host"].([]interface{}); len(host) > 0 {
				widgets = append(widgets, mackerel.Widget{
					Type:   "graph",
					Title:  title,
					Graph:  expandDashboardGraphHost(host),
					Range:  r,
					Layout: layout,
				})
			}
			if role := g["role"].([]interface{}); len(role) > 0 {
				widgets = append(widgets, mackerel.Widget{
					Type:   "graph",
					Title:  title,
					Graph:  expandDashboardGraphRole(role),
					Range:  r,
					Layout: layout,
				})
			}
			if service := g["service"].([]interface{}); len(service) > 0 {
				widgets = append(widgets, mackerel.Widget{
					Type:   "graph",
					Title:  title,
					Graph:  expandDashboardGraphService(service),
					Range:  r,
					Layout: layout,
				})
			}
			if expression := g["expression"].([]interface{}); len(expression) > 0 {
				widgets = append(widgets, mackerel.Widget{
					Type:   "graph",
					Title:  title,
					Graph:  expandDashboardGraphExpression(expression),
					Range:  r,
					Layout: layout,
				})
			}
			if query := g["query"].([]interface{}); len(query) > 0 {
				widgets = append(widgets, mackerel.Widget{
					Type:   "graph",
					Title:  title,
					Graph:  expandDashboardGraphQuery(query),
					Range:  r,
					Layout: layout,
				})
			}
		}
	}
	if _, ok := d.GetOk("value"); ok {
		values := d.Get("value").([]interface{})
		for _, value := range values {
			v := value.(map[string]interface{})
			title := v["title"].(string)
			metric := v["metric"].([]interface{})[0].(map[string]interface{})
			var fractionSize int64
			if fs, ok := v["fraction_size"].(int); ok {
				fractionSize = int64(fs)
			}
			suffix := v["suffix"].(string)
			layout := expandDashboardLayout(v["layout"].([]interface{})[0].(map[string]interface{}))

			if host := metric["host"].([]interface{}); len(host) > 0 {
				widgets = append(widgets, mackerel.Widget{
					Type:         "value",
					Title:        title,
					Metric:       expandDashboardValueHost(host),
					FractionSize: &fractionSize,
					Suffix:       suffix,
					Layout:       layout,
				})
			}
			if service := metric["service"].([]interface{}); len(service) > 0 {
				widgets = append(widgets, mackerel.Widget{
					Type:         "value",
					Title:        title,
					Metric:       expandDashboardValueService(service),
					FractionSize: &fractionSize,
					Suffix:       suffix,
					Layout:       layout,
				})
			}
			if expression := metric["expression"].([]interface{}); len(expression) > 0 {
				widgets = append(widgets, mackerel.Widget{
					Type:         "value",
					Title:        title,
					Metric:       expandDashboardValueExpression(expression),
					FractionSize: &fractionSize,
					Suffix:       suffix,
					Layout:       layout,
				})
			}
			if query := metric["query"].([]interface{}); len(query) > 0 {
				widgets = append(widgets, mackerel.Widget{
					Type:         "value",
					Title:        title,
					Metric:       expandDashboardValueQuery(query),
					FractionSize: &fractionSize,
					Suffix:       suffix,
					Layout:       layout,
				})
			}
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
	if _, ok := d.GetOk("alert_status"); ok {
		alert_statuses := d.Get("alert_status").([]interface{})
		for _, alert_status := range alert_statuses {
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

func expandDashboardGraphQuery(query []interface{}) mackerel.Graph {
	q := query[0].(map[string]interface{})
	g := mackerel.Graph{
		Type:  "query",
		Query: q["query"].(string),
	}
	if legend, ok := q["legend"].(string); ok {
		g.Legend = legend
	}
	return g
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

func expandDashboardValueQuery(query []interface{}) mackerel.Metric {
	q := query[0].(map[string]interface{})
	m := mackerel.Metric{
		Type:  "query",
		Query: q["query"].(string),
	}
	if legend, ok := q["legend"]; ok {
		m.Legend = legend.(string)
	}
	return m
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
