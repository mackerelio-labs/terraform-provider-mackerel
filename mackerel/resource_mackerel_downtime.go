package mackerel

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelDowntime() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelDowntimeCreate,
		Read:   resourceMackerelDowntimeRead,
		Update: resourceMackerelDowntimeUpdate,
		Delete: resourceMackerelDowntimeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceMackerelDowntimeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	dt, err := client.CreateDowntime(buildDowntimeStruct(d))
	if err != nil {
		return err
	}
	d.SetId(dt.ID)
	return resourceMackerelDowntimeRead(d, meta)
}

func resourceMackerelDowntimeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	downtimes, err := client.FindDowntimes()
	if err != nil {
		return err
	}
	for _, downtime := range downtimes {
		if downtime.ID == d.Id() {
			d.Set("name", downtime.Name)
			d.Set("memo", downtime.Memo)
			d.Set("start", downtime.Start)
			d.Set("duration", downtime.Duration)
			if downtime.Recurrence != nil {
				weekdays := make([]string, 0, len(downtime.Recurrence.Weekdays))
				for _, weekday := range downtime.Recurrence.Weekdays {
					weekdays = append(weekdays, weekday.String())
				}
				d.Set("recurrence", []map[string]interface{}{
					{
						"type":     downtime.Recurrence.Type.String(),
						"interval": downtime.Recurrence.Interval,
						"weekdays": flattenStringListToSet(weekdays),
						"until":    downtime.Recurrence.Until,
					},
				})
			}
			d.Set("service_scopes", flattenStringListToSet(downtime.ServiceScopes))
			d.Set("service_exclude_scopes", flattenStringListToSet(downtime.ServiceExcludeScopes))
			d.Set("role_scopes", flattenStringListToSet(downtime.RoleScopes))
			d.Set("role_exclude_scopes", flattenStringListToSet(downtime.RoleExcludeScopes))
			d.Set("monitor_scopes", flattenStringListToSet(downtime.MonitorScopes))
			d.Set("monitor_exclude_scopes", flattenStringListToSet(downtime.MonitorExcludeScopes))
			break
		}
	}
	return nil
}

func resourceMackerelDowntimeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	_, err := client.UpdateDowntime(d.Id(), buildDowntimeStruct(d))
	if err != nil {
		return err
	}

	return resourceMackerelDowntimeRead(d, meta)
}

func resourceMackerelDowntimeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)
	_, err := client.DeleteDowntime(d.Id())
	return err
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

func buildDowntimeStruct(d *schema.ResourceData) *mackerel.Downtime {
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
