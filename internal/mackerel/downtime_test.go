package mackerel

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

func Test_Downtime_ReadDowntime(t *testing.T) {
	t.Parallel()

	defaultClient := func() ([]*mackerel.Downtime, error) {
		return []*mackerel.Downtime{
			{
				ID:       "5ghjb6vgDFN",
				Name:     "basic",
				Start:    1735707600,
				Duration: 3600,
			},
			{
				ID:       "5ghjbbVABY5",
				Name:     "full",
				Memo:     "This downtime is managed by Terraform.",
				Start:    1735707600,
				Duration: 3600,
				Recurrence: &mackerel.DowntimeRecurrence{
					Type:     mackerel.DowntimeRecurrenceTypeWeekly,
					Interval: 2,
					Weekdays: []mackerel.DowntimeWeekday{
						mackerel.DowntimeWeekday(time.Wednesday),
						mackerel.DowntimeWeekday(time.Thursday),
					},
					Until: 1767193199,
				},
				ServiceScopes:        []string{"include-svc"},
				ServiceExcludeScopes: []string{"exclude-svc"},
				RoleScopes:           []string{"svc: include-role"},
				RoleExcludeScopes:    []string{"svc: exclude-role"},
				MonitorScopes:        []string{"5ghjb7CrJ43"},
				MonitorExcludeScopes: []string{"5ghjbaziveA"},
			},
		}, nil
	}

	cases := map[string]struct {
		inClient downtimeFinderFunc
		inID     string

		wants   DowntimeModel
		wantErr bool
	}{
		"basic": {
			inClient: defaultClient,
			inID:     "5ghjb6vgDFN",

			wants: DowntimeModel{
				ID:                   types.StringValue("5ghjb6vgDFN"),
				Name:                 types.StringValue("basic"),
				Memo:                 types.StringValue(""),
				Start:                types.Int64Value(1735707600),
				Duration:             types.Int64Value(3600),
				ServiceScopes:        []string{},
				ServiceExcludeScopes: []string{},
				RoleScopes:           []string{},
				RoleExcludeScopes:    []string{},
				MonitorScopes:        []string{},
				MonitorExcludeScopes: []string{},
			},
		},
		"no downtime": {
			inClient: defaultClient,
			inID:     "nonexist",

			wantErr: true,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			model, err := readDowntime(tt.inClient, tt.inID)
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error: %+v", err)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(*model, tt.wants); diff != "" {
				t.Error(diff)
			}
		})
	}
}

type downtimeFinderFunc func() ([]*mackerel.Downtime, error)

func (f downtimeFinderFunc) FindDowntimes() ([]*mackerel.Downtime, error) {
	return f()
}

func Test_Downtime_conv(t *testing.T) {
	t.Parallel()

	// api <-> model
	cases := map[string]struct {
		api   mackerel.Downtime
		model DowntimeModel
	}{
		"basic": {
			api: mackerel.Downtime{
				ID:                   "5ghjb6vgDFN",
				Name:                 "basic",
				Start:                1735707600,
				Duration:             3600,
				ServiceScopes:        []string{},
				ServiceExcludeScopes: []string{},
				RoleScopes:           []string{},
				RoleExcludeScopes:    []string{},
				MonitorScopes:        []string{},
				MonitorExcludeScopes: []string{},
			},
			model: DowntimeModel{
				ID:                   types.StringValue("5ghjb6vgDFN"),
				Name:                 types.StringValue("basic"),
				Memo:                 types.StringValue(""),
				Start:                types.Int64Value(1735707600),
				Duration:             types.Int64Value(3600),
				ServiceScopes:        []string{},
				ServiceExcludeScopes: []string{},
				RoleScopes:           []string{},
				RoleExcludeScopes:    []string{},
				MonitorScopes:        []string{},
				MonitorExcludeScopes: []string{},
			},
		},
		"full": {
			api: mackerel.Downtime{
				ID:       "5ghjbbVABY5",
				Name:     "full",
				Memo:     "This downtime is managed by Terraform.",
				Start:    1735707600,
				Duration: 3600,
				Recurrence: &mackerel.DowntimeRecurrence{
					Type:     mackerel.DowntimeRecurrenceTypeWeekly,
					Interval: 2,
					Weekdays: []mackerel.DowntimeWeekday{
						mackerel.DowntimeWeekday(time.Wednesday),
						mackerel.DowntimeWeekday(time.Thursday),
					},
					Until: 1767193199,
				},
				ServiceScopes:        []string{"include-svc"},
				ServiceExcludeScopes: []string{"exclude-svc"},
				RoleScopes:           []string{"svc: include-role"},
				RoleExcludeScopes:    []string{"svc: exclude-role"},
				MonitorScopes:        []string{"5ghjb7CrJ43"},
				MonitorExcludeScopes: []string{"5ghjbaziveA"},
			},
			model: DowntimeModel{
				ID:       types.StringValue("5ghjbbVABY5"),
				Name:     types.StringValue("full"),
				Memo:     types.StringValue("This downtime is managed by Terraform."),
				Start:    types.Int64Value(1735707600),
				Duration: types.Int64Value(3600),
				Recurrence: []DowntimeRecurrence{{
					Type:     types.StringValue("weekly"),
					Interval: types.Int64Value(2),
					Weekdays: []string{"Wednesday", "Thursday"},
					Until:    types.Int64Value(1767193199),
				}},
				ServiceScopes:        []string{"include-svc"},
				ServiceExcludeScopes: []string{"exclude-svc"},
				RoleScopes:           []string{"svc: include-role"},
				RoleExcludeScopes:    []string{"svc: exclude-role"},
				MonitorScopes:        []string{"5ghjb7CrJ43"},
				MonitorExcludeScopes: []string{"5ghjbaziveA"},
			},
		},
	}

	for name, tt := range cases {
		t.Run(name+"/fromAPI", func(t *testing.T) {
			t.Parallel()

			model := newDowntime(tt.api)
			if diff := cmp.Diff(*model, tt.model); diff != "" {
				t.Error(diff)
			}
		})
		t.Run(name+"/toAPI", func(t *testing.T) {
			t.Parallel()

			api := tt.model.mackerelDowntime()
			if diff := cmp.Diff(*api, tt.api); diff != "" {
				t.Error(diff)
			}
		})
	}
}
