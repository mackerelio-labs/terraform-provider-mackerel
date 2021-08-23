package mackerel

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelDowntime() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMackerelDowntimeCreate,
		ReadContext:   resourceMackerelDowntimeRead,
		UpdateContext: resourceMackerelDowntimeUpdate,
		DeleteContext: resourceMackerelDowntimeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"memo": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"recurrence": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"hourly", "daily", "weekly", "monthly", "yearly"}, false),
						},
						"interval": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"weekdays": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, false),
							},
						},
						"until": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
					},
				},
			},
			"service_scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"service_exclude_scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"role_scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"role_exclude_scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"monitor_scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"monitor_exclude_scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceMackerelDowntimeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mackerel.Client)
	dt, err := client.CreateDowntime(expandDowntime(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dt.ID)
	return resourceMackerelDowntimeRead(ctx, d, m)
}

func resourceMackerelDowntimeRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mackerel.Client)
	downtimes, err := client.FindDowntimes()
	if err != nil {
		return diag.FromErr(err)
	}
	var downtime *mackerel.Downtime
	for _, dt := range downtimes {
		if dt.ID == d.Id() {
			downtime = dt
			break
		}
	}
	if downtime == nil {
		return diag.Errorf("the ID '%s' does not match any downtime in mackerel.io", d.Id())
	}
	return flattenDowntime(downtime, d)
}

func resourceMackerelDowntimeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mackerel.Client)
	_, err := client.UpdateDowntime(d.Id(), expandDowntime(d))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceMackerelDowntimeRead(ctx, d, m)
}

func resourceMackerelDowntimeDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*mackerel.Client)
	_, err := client.DeleteDowntime(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

var stringToRecurrenceType = map[string]mackerel.DowntimeRecurrenceType{
	"hourly":  mackerel.DowntimeRecurrenceTypeHourly,
	"daily":   mackerel.DowntimeRecurrenceTypeDaily,
	"weekly":  mackerel.DowntimeRecurrenceTypeWeekly,
	"monthly": mackerel.DowntimeRecurrenceTypeMonthly,
	"yearly":  mackerel.DowntimeRecurrenceTypeYearly,
}

var stringToWeekday = map[string]mackerel.DowntimeWeekday{
	"Sunday":    mackerel.DowntimeWeekday(time.Sunday),
	"Monday":    mackerel.DowntimeWeekday(time.Monday),
	"Tuesday":   mackerel.DowntimeWeekday(time.Tuesday),
	"Wednesday": mackerel.DowntimeWeekday(time.Wednesday),
	"Thursday":  mackerel.DowntimeWeekday(time.Thursday),
	"Friday":    mackerel.DowntimeWeekday(time.Friday),
	"Saturday":  mackerel.DowntimeWeekday(time.Saturday),
}

func expandDowntime(d *schema.ResourceData) *mackerel.Downtime {
	downtime := &mackerel.Downtime{
		Name:                 d.Get("name").(string),
		Memo:                 d.Get("memo").(string),
		Start:                int64(d.Get("start").(int)),
		Duration:             int64(d.Get("duration").(int)),
		Recurrence:           nil,
		ServiceScopes:        expandStringListFromSet(d.Get("service_scopes").(*schema.Set)),
		ServiceExcludeScopes: expandStringListFromSet(d.Get("service_exclude_scopes").(*schema.Set)),
		RoleScopes:           expandStringListFromSet(d.Get("role_scopes").(*schema.Set)),
		RoleExcludeScopes:    expandStringListFromSet(d.Get("role_exclude_scopes").(*schema.Set)),
		MonitorScopes:        expandStringListFromSet(d.Get("monitor_scopes").(*schema.Set)),
		MonitorExcludeScopes: expandStringListFromSet(d.Get("monitor_exclude_scopes").(*schema.Set)),
	}
	if _, ok := d.GetOk("recurrence"); ok {
		var recurrence mackerel.DowntimeRecurrence
		if v, ok := d.GetOk("recurrence.0.type"); ok {
			if rType, ok := stringToRecurrenceType[v.(string)]; ok {
				recurrence.Type = rType
			}
		}
		if v, ok := d.GetOk("recurrence.0.interval"); ok {
			recurrence.Interval = int64(v.(int))
		}
		if v, ok := d.GetOk("recurrence.0.weekdays"); ok {
			set := v.(*schema.Set)
			weekdays := make([]mackerel.DowntimeWeekday, 0, set.Len())
			for _, weekday := range set.List() {
				if rWeekday, ok := stringToWeekday[weekday.(string)]; ok {
					weekdays = append(weekdays, rWeekday)
				}
			}
			recurrence.Weekdays = weekdays
		}
		if v, ok := d.GetOk("recurrence.0.until"); ok {
			recurrence.Until = int64(v.(int))
		}
		downtime.Recurrence = &recurrence
	}
	return downtime
}
