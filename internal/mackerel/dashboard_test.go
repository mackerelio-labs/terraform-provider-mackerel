package mackerel

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

func Test_Dashboard_conv(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		api   mackerel.Dashboard
		model DashboardModel
	}{
		"role graph": {
			api: mackerel.Dashboard{
				ID:      "2c5bLca8d",
				Title:   "role graph dashboard title",
				Memo:    "role graph dashboard",
				URLPath: "role-graph-dashboard",
				Widgets: []mackerel.Widget{{
					Type:  "graph",
					Title: "role graph",
					Graph: mackerel.Graph{
						Type:         "role",
						RoleFullName: "service:role",
						Name:         "loadavg5",
						IsStacked:    true,
					},
					Range: mackerel.Range{
						Type:   "relative",
						Period: 3600,
						Offset: 1800,
					},
					Layout: mackerel.Layout{
						X:      2,
						Y:      12,
						Width:  10,
						Height: 8,
					},
				}},
				CreatedAt: 1439346145003,
				UpdatedAt: 1439346145003,
			},
			model: DashboardModel{
				ID:      types.StringValue("2c5bLca8d"),
				Title:   types.StringValue("role graph dashboard title"),
				Memo:    types.StringValue("role graph dashboard"),
				URLPath: types.StringValue("role-graph-dashboard"),
				Graph: []DashboardWidgetGraph{{
					DashboardWidget: DashboardWidget{
						Title: types.StringValue("role graph"),
						Layout: []DashboardLayout{{
							X:      types.Int64Value(2),
							Y:      types.Int64Value(12),
							Width:  types.Int64Value(10),
							Height: types.Int64Value(8),
						}},
					},
					Role: []DashboardGraphRole{{
						RoleFullname: types.StringValue("service:role"),
						Name:         types.StringValue("loadavg5"),
						IsStacked:    types.BoolValue(true),
					}},
					Range: []DashboardRange{{
						Relative: []DashboardRangeRelative{{
							Period: types.Int64Value(3600),
							Offset: types.Int64Value(1800),
						}},
					}},
				}},
				Value:       []DashboardWidgetValue{},
				Markdown:    []DashboardWidgetMarkdown{},
				AlertStatus: []DashboardWidgetAlertStatus{},
				CreatedAt:   types.Int64Value(1439346145003),
				UpdatedAt:   types.Int64Value(1439346145003),
			},
		},
	}

	for name, tt := range cases {
		t.Run(name+"-toModel", func(t *testing.T) {
			t.Parallel()

			model, err := newDashboard(tt.api)
			if err != nil {
				t.Errorf("unexpected error: %+v", err)
				return
			}

			if diff := cmp.Diff(model, tt.model); diff != "" {
				t.Error(diff)
			}
		})
		t.Run(name+"-toAPI", func(t *testing.T) {
			t.Parallel()

			api := tt.model.mackerelDashboard()
			if diff := cmp.Diff(api, tt.api); diff != "" {
				t.Error(diff)
			}
		})
	}
}
