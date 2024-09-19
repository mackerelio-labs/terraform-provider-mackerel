package mackerel

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

type DowntimeModel struct {
	ID                   types.String         `tfsdk:"id"`
	Name                 types.String         `tfsdk:"name"`
	Memo                 types.String         `tfsdk:"memo"`
	Start                types.Int64          `tfsdk:"start"`
	Duration             types.Int64          `tfsdk:"duration"`
	Recurrence           []DowntimeRecurrence `tfsdk:"recurrence"` // length <= 1
	ServiceScopes        []string             `tfsdk:"service_scopes"`
	ServiceExcludeScopes []string             `tfsdk:"service_exclude_scopes"`
	RoleScopes           []string             `tfsdk:"role_scopes"`
	RoleExcludeScopes    []string             `tfsdk:"role_exclude_scopes"`
	MonitorScopes        []string             `tfsdk:"monitor_scopes"`
	MonitorExcludeScopes []string             `tfsdk:"monitor_exclude_scopes"`
}
type DowntimeRecurrence struct {
	Type     types.String `tfsdk:"type"`
	Interval types.Int64  `tfsdk:"interval"`
	Weekdays []string     `tfsdk:"weekdays"`
	Until    types.Int64  `tfsdk:"until"`
}

func ReadDowntime(_ context.Context, client *Client, id string) (*DowntimeModel, error) {
	return readDowntime(client, id)
}

type downtimeFinder interface {
	FindDowntimes() ([]*mackerel.Downtime, error)
}

func readDowntime(client downtimeFinder, id string) (*DowntimeModel, error) {
	downtimes, err := client.FindDowntimes()
	if err != nil {
		return nil, err
	}

	downtimeIdx := slices.IndexFunc(downtimes, func(d *mackerel.Downtime) bool {
		return d.ID == id
	})
	if downtimeIdx < 0 {
		return nil, fmt.Errorf("the ID '%s' does not match any downtime in mackerel.io", id)
	}

	return newDowntime(*downtimes[downtimeIdx]), nil
}

func (d *DowntimeModel) Create(_ context.Context, client *Client) error {
	createdDowntime, err := client.CreateDowntime(d.mackerelDowntime())
	if err != nil {
		return err
	}

	d.ID = types.StringValue(createdDowntime.ID)
	return nil
}

func (d *DowntimeModel) Read(_ context.Context, client *Client) error {
	newModel, err := readDowntime(client, d.ID.ValueString())
	if err != nil {
		return err
	}
	d.merge(*newModel)
	return nil
}

func (d *DowntimeModel) Update(_ context.Context, client *Client) error {
	if _, err := client.UpdateDowntime(d.ID.ValueString(), d.mackerelDowntime()); err != nil {
		return err
	}
	return nil
}

func (d *DowntimeModel) Delete(_ context.Context, client *Client) error {
	if _, err := client.DeleteDowntime(d.ID.ValueString()); err != nil {
		return err
	}
	return nil
}

func newDowntime(d mackerel.Downtime) *DowntimeModel {
	model := &DowntimeModel{
		ID:                   types.StringValue(d.ID),
		Name:                 types.StringValue(d.Name),
		Memo:                 types.StringValue(d.Memo),
		Start:                types.Int64Value(d.Start),
		Duration:             types.Int64Value(d.Duration),
		ServiceScopes:        d.ServiceScopes,
		ServiceExcludeScopes: d.ServiceExcludeScopes,
		RoleScopes:           d.RoleScopes,
		RoleExcludeScopes:    d.RoleExcludeScopes,
		MonitorScopes:        d.MonitorScopes,
		MonitorExcludeScopes: d.MonitorExcludeScopes,
	}
	if d.Recurrence != nil {
		recurrence := DowntimeRecurrence{
			Type:     types.StringValue(d.Recurrence.Type.String()),
			Interval: types.Int64Value(d.Recurrence.Interval),
			Until:    types.Int64Value(d.Recurrence.Until),
		}
		recurrence.Weekdays = make([]string, 0, len(d.Recurrence.Weekdays))
		for _, wd := range d.Recurrence.Weekdays {
			recurrence.Weekdays = append(recurrence.Weekdays, wd.String())
		}
		model.Recurrence = []DowntimeRecurrence{recurrence}
	}
	return model
}

var stringsToMackerelRecurrenceType = map[string]mackerel.DowntimeRecurrenceType{
	"hourly":  mackerel.DowntimeRecurrenceTypeHourly,
	"daily":   mackerel.DowntimeRecurrenceTypeDaily,
	"weekly":  mackerel.DowntimeRecurrenceTypeWeekly,
	"monthly": mackerel.DowntimeRecurrenceTypeMonthly,
	"yearly":  mackerel.DowntimeRecurrenceTypeYearly,
}

var stringToMackerelWeekday = map[string]mackerel.DowntimeWeekday{
	"Sunday":    mackerel.DowntimeWeekday(time.Sunday),
	"Monday":    mackerel.DowntimeWeekday(time.Monday),
	"Tuesday":   mackerel.DowntimeWeekday(time.Tuesday),
	"Wednesday": mackerel.DowntimeWeekday(time.Wednesday),
	"Thursday":  mackerel.DowntimeWeekday(time.Thursday),
	"Friday":    mackerel.DowntimeWeekday(time.Friday),
	"Saturday":  mackerel.DowntimeWeekday(time.Saturday),
}

func (d *DowntimeModel) mackerelDowntime() *mackerel.Downtime {
	mackerelDowntime := &mackerel.Downtime{
		ID:                   d.ID.ValueString(),
		Name:                 d.Name.ValueString(),
		Memo:                 d.Memo.ValueString(),
		Start:                d.Start.ValueInt64(),
		Duration:             d.Duration.ValueInt64(),
		ServiceScopes:        d.ServiceScopes,
		ServiceExcludeScopes: d.ServiceExcludeScopes,
		RoleScopes:           d.RoleScopes,
		RoleExcludeScopes:    d.RoleExcludeScopes,
		MonitorScopes:        d.MonitorScopes,
		MonitorExcludeScopes: d.MonitorExcludeScopes,
	}
	if len(d.Recurrence) == 1 {
		recurrence := d.Recurrence[0]

		recurrenceType, ok := stringsToMackerelRecurrenceType[recurrence.Type.ValueString()]
		if !ok {
			panic(fmt.Errorf("invalid recurrence type: %v", recurrence.Type))
		}

		mackerelWeekdays := make([]mackerel.DowntimeWeekday, 0, len(recurrence.Weekdays))
		for _, weekday := range recurrence.Weekdays {
			mackerelWeekday, ok := stringToMackerelWeekday[weekday]
			if !ok {
				panic(fmt.Errorf("invalid weekday: %s", weekday))
			}
			mackerelWeekdays = append(mackerelWeekdays, mackerelWeekday)
		}

		mackerelRecurrence := &mackerel.DowntimeRecurrence{
			Type:     recurrenceType,
			Interval: recurrence.Interval.ValueInt64(),
			Weekdays: mackerelWeekdays,
			Until:    recurrence.Until.ValueInt64(),
		}
		mackerelDowntime.Recurrence = mackerelRecurrence

	}
	return mackerelDowntime
}

func (d *DowntimeModel) merge(newModel DowntimeModel) {
	// preserve nil
	if len(d.Recurrence) == 1 && len(d.Recurrence[0].Weekdays) == 0 &&
		len(newModel.Recurrence) == 1 && len(newModel.Recurrence[0].Weekdays) == 0 {
		newModel.Recurrence[0].Weekdays = d.Recurrence[0].Weekdays
	}
	if len(d.ServiceScopes) == 0 && len(newModel.ServiceScopes) == 0 {
		newModel.ServiceScopes = d.ServiceScopes
	}
	if len(d.ServiceExcludeScopes) == 0 && len(newModel.ServiceExcludeScopes) == 0 {
		newModel.ServiceExcludeScopes = d.ServiceExcludeScopes
	}
	if len(d.RoleScopes) == 0 && len(newModel.RoleScopes) == 0 {
		newModel.RoleScopes = d.RoleScopes
	}
	if len(d.RoleExcludeScopes) == 0 && len(newModel.RoleExcludeScopes) == 0 {
		newModel.RoleExcludeScopes = d.RoleExcludeScopes
	}
	if len(d.MonitorScopes) == 0 && len(newModel.MonitorScopes) == 0 {
		newModel.MonitorScopes = d.MonitorScopes
	}
	if len(d.MonitorExcludeScopes) == 0 && len(newModel.MonitorExcludeScopes) == 0 {
		newModel.MonitorExcludeScopes = d.MonitorExcludeScopes
	}
	*d = newModel
}
