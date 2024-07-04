package mackerel

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

func Test_NotificationGroup_Read(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		inID     string
		inClient notificationGroupFinderFunc

		wantErr bool
		wants   NotificationGroupModel
	}{
		"valid": {
			inID: "ng0",
			inClient: func() ([]*mackerel.NotificationGroup, error) {
				return []*mackerel.NotificationGroup{{
					ID:                        "ng0",
					Name:                      "notification group",
					NotificationLevel:         mackerel.NotificationLevelAll,
					ChildNotificationGroupIDs: []string{"cng0"},
					ChildChannelIDs:           []string{"cc0"},
					Monitors: []*mackerel.NotificationGroupMonitor{{
						ID:          "monitor0",
						SkipDefault: true,
					}},
					Services: []*mackerel.NotificationGroupService{{Name: "service"}},
				}}, nil
			},

			wants: NotificationGroupModel{
				ID:                        types.StringValue("ng0"),
				Name:                      types.StringValue("notification group"),
				NotificationLevel:         types.StringValue("all"),
				ChildNotificationGroupIDs: []types.String{types.StringValue("cng0")},
				ChildChannelIDs:           []types.String{types.StringValue("cc0")},
				Monitors: []NotificationTargetMonitorModel{{
					ID:          types.StringValue("monitor0"),
					SkipDefault: types.BoolValue(true),
				}},
				Services: []NotificationTargetServiceModel{{Name: types.StringValue("service")}},
			},
		},
		"missing": {
			inID: "ng1",
			inClient: func() ([]*mackerel.NotificationGroup, error) {
				return []*mackerel.NotificationGroup{{ID: "ng0"}}, nil
			},

			wantErr: true,
		},
	}

	ctx := context.Background()
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ng, err := readNotificationGroupInner(ctx, tt.inClient, tt.inID)
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Errorf("expect error, but got no error")
				} else {
					t.Errorf("unexpected error: %+v", err)
				}
				return
			}

			if diff := cmp.Diff(ng, tt.wants); diff != "" {
				t.Error(diff)
			}
		})
	}
}

type notificationGroupFinderFunc func() ([]*mackerel.NotificationGroup, error)

func (f notificationGroupFinderFunc) FindNotificationGroups() ([]*mackerel.NotificationGroup, error) {
	return f()
}

func Test_NotificationGroup_Create(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in   NotificationGroupModel
		inID string

		wantReq mackerel.NotificationGroup
		wants   NotificationGroupModel
	}{
		"valid": {
			in: NotificationGroupModel{
				ID:                        types.StringUnknown(),
				Name:                      types.StringValue("notification group"),
				NotificationLevel:         types.StringValue("all"),
				ChildNotificationGroupIDs: []types.String{types.StringValue("cng0")},
				ChildChannelIDs:           []types.String{types.StringValue("cc0")},
				Monitors: []NotificationTargetMonitorModel{{
					ID:          types.StringValue("monitor0"),
					SkipDefault: types.BoolValue(true),
				}},
				Services: []NotificationTargetServiceModel{{Name: types.StringValue("service")}},
			},
			inID: "ng0",

			wantReq: mackerel.NotificationGroup{
				ID:                        "",
				Name:                      "notification group",
				NotificationLevel:         mackerel.NotificationLevelAll,
				ChildNotificationGroupIDs: []string{"cng0"},
				ChildChannelIDs:           []string{"cc0"},
				Monitors: []*mackerel.NotificationGroupMonitor{{
					ID:          "monitor0",
					SkipDefault: true,
				}},
				Services: []*mackerel.NotificationGroupService{{Name: "service"}},
			},
			wants: NotificationGroupModel{
				ID:                        types.StringValue("ng0"),
				Name:                      types.StringValue("notification group"),
				NotificationLevel:         types.StringValue("all"),
				ChildNotificationGroupIDs: []types.String{types.StringValue("cng0")},
				ChildChannelIDs:           []types.String{types.StringValue("cc0")},
				Monitors: []NotificationTargetMonitorModel{{
					ID:          types.StringValue("monitor0"),
					SkipDefault: types.BoolValue(true),
				}},
				Services: []NotificationTargetServiceModel{{Name: types.StringValue("service")}},
			},
		},
	}

	ctx := context.Background()
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := tt.in
			client := &notificationGroupCreatorTester{ID: tt.inID}
			if err := m.createInner(ctx, client); err != nil {
				t.Errorf("unexpected error: %+v", err)
				return
			}

			if diff := cmp.Diff(client.Request, tt.wantReq); diff != "" {
				t.Errorf("invalid request:\n%s", diff)
			}
			if diff := cmp.Diff(m, tt.wants); diff != "" {
				t.Errorf("invalid model:\n%s", diff)
			}
		})
	}
}

type notificationGroupCreatorTester struct {
	ID string

	Request mackerel.NotificationGroup
}

func (ct *notificationGroupCreatorTester) CreateNotificationGroup(param *mackerel.NotificationGroup) (*mackerel.NotificationGroup, error) {
	ct.Request = *param
	data := *param
	data.ID = ct.ID
	return &data, nil
}
