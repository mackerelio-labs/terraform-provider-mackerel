package mackerel

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

func Test_AlertGroupSetting_conv(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		api   mackerel.AlertGroupSetting
		model AlertGroupSettingModel
	}{
		"basic": {
			api: mackerel.AlertGroupSetting{
				ID:            "5fCTLAQFbhy",
				Name:          "tf-alert-group-basic",
				ServiceScopes: []string{},
				RoleScopes:    []string{},
				MonitorScopes: []string{},
			},
			model: AlertGroupSettingModel{
				ID:                   types.StringValue("5fCTLAQFbhy"),
				Name:                 types.StringValue("tf-alert-group-basic"),
				Memo:                 types.StringValue(""),
				ServiceScopes:        []string{},
				RoleScopes:           []string{},
				MonitorScopes:        []string{},
				NotificationInterval: types.Int64Value(0),
			},
		},
		"full": {
			api: mackerel.AlertGroupSetting{
				ID:                   "5fCTLFmatvs",
				Name:                 "tf-alert-group-full",
				Memo:                 "This alert group setting is managed by Terraform.",
				ServiceScopes:        []string{"tf-svc"},
				RoleScopes:           []string{"tf-svc:tf-role"},
				MonitorScopes:        []string{"5fCTLBchbQG"},
				NotificationInterval: 60,
			},
			model: AlertGroupSettingModel{
				ID:                   types.StringValue("5fCTLFmatvs"),
				Name:                 types.StringValue("tf-alert-group-full"),
				Memo:                 types.StringValue("This alert group setting is managed by Terraform."),
				ServiceScopes:        []string{"tf-svc"},
				RoleScopes:           []string{"tf-svc:tf-role"},
				MonitorScopes:        []string{"5fCTLBchbQG"},
				NotificationInterval: types.Int64Value(60),
			},
		},
	}

	for name, tt := range cases {
		t.Run(name+"-toModel", func(t *testing.T) {
			t.Parallel()

			model := newAlertGroupSetting(tt.api)
			if diff := cmp.Diff(model, tt.model); diff != "" {
				t.Error(diff)
			}
		})
		t.Run(name+"-toAPI", func(t *testing.T) {
			t.Parallel()

			api := tt.model.mackerelAlertGroupSetting()
			if diff := cmp.Diff(api, tt.api); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func Test_AlertGroupSetting_merge(t *testing.T) {
	t.Parallel()

	// lhs <- rhs = wants
	cases := map[string]struct {
		lhs   AlertGroupSettingModel
		rhs   AlertGroupSettingModel
		wants AlertGroupSettingModel
	}{
		"nil preserving": {
			lhs: AlertGroupSettingModel{
				Name: types.StringValue("before"),
			},
			rhs: AlertGroupSettingModel{
				Name:          types.StringValue("after"),
				ServiceScopes: []string{},
				RoleScopes:    []string{},
				MonitorScopes: []string{},
			},
			wants: AlertGroupSettingModel{
				Name: types.StringValue("after"),
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := tt.lhs
			m.merge(tt.rhs)
			if diff := cmp.Diff(m, tt.wants); diff != "" {
				t.Error(diff)
			}

		})
	}
}
