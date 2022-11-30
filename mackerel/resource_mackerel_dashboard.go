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
			Required: true,
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
							Type:     schema.TypeSet,
							Required: true,
							Elem:     dashboardGraphResource,
						},
						"range": {
							Type:     schema.TypeSet,
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
	if _, ok := d.GetOk("markdown"); ok {
		widgets = append(widgets, mackerel.Widget{
			Type:     "markdown",
			Title:    d.Get("markdown.0.title").(string),
			Markdown: d.Get("markdown.0.markdown").(string),
			Layout:   expandDashboardLayout(d),
		})
	}

	return widgets
}

func expandDashboardLayout(d *schema.ResourceData) mackerel.Layout {
	return mackerel.Layout{
		X:      int64(d.Get("markdown.0.layout.0.x").(int)),
		Y:      int64(d.Get("markdown.0.layout.0.y").(int)),
		Width:  int64(d.Get("markdown.0.layout.0.width").(int)),
		Height: int64(d.Get("markdown.0.layout.0.height").(int)),
	}
}
