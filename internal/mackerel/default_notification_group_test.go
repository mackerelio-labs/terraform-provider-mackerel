package mackerel

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

func Test_DefaultNotificationGroup_Read(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		inClient notificationGroupFinderFunc
		wantErr  bool
		wants    DefaultNotificationGroupModel
	}{
		"valid": {
			inClient: func() ([]*mackerel.NotificationGroup, error) {
				return []*mackerel.NotificationGroup{
					{ID: "ng0", Name: "custom"},
					{
						ID:                        "default",
						Type:                      mackerel.NotificationGroupTypeGroupDefault,
						Name:                      "Default",
						NotificationLevel:         mackerel.NotificationLevelAll,
						ChildNotificationGroupIDs: []string{"child-ng"},
						ChildChannelIDs:           []string{"channel"},
						Monitors: []*mackerel.NotificationGroupMonitor{{
							ID:          "monitor",
							SkipDefault: true,
						}},
						Services: []*mackerel.NotificationGroupService{{Name: "service"}},
					},
				}, nil
			},
			wants: DefaultNotificationGroupModel{
				ID:                        types.StringValue("default"),
				NotificationLevel:         types.StringValue("all"),
				ChildNotificationGroupIDs: []types.String{types.StringValue("child-ng")},
				ChildChannelIDs:           []types.String{types.StringValue("channel")},
			},
		},
		"selects_group_default_type_over_name": {
			inClient: func() ([]*mackerel.NotificationGroup, error) {
				return []*mackerel.NotificationGroup{
					{ID: "normal", Type: mackerel.NotificationGroupTypeGroup, Name: "Default"},
					{
						ID:                        "default",
						Type:                      mackerel.NotificationGroupTypeGroupDefault,
						Name:                      "renamed default",
						NotificationLevel:         mackerel.NotificationLevelCritical,
						ChildNotificationGroupIDs: []string{"child-ng"},
						ChildChannelIDs:           []string{"channel"},
					},
				}, nil
			},
			wants: DefaultNotificationGroupModel{
				ID:                        types.StringValue("default"),
				NotificationLevel:         types.StringValue("critical"),
				ChildNotificationGroupIDs: []types.String{types.StringValue("child-ng")},
				ChildChannelIDs:           []types.String{types.StringValue("channel")},
			},
		},
		"missing": {
			inClient: func() ([]*mackerel.NotificationGroup, error) {
				return []*mackerel.NotificationGroup{{ID: "ng0", Type: mackerel.NotificationGroupTypeGroup, Name: "Default"}}, nil
			},
			wantErr: true,
		},
		"duplicate": {
			inClient: func() ([]*mackerel.NotificationGroup, error) {
				return []*mackerel.NotificationGroup{
					{ID: "default0", Type: mackerel.NotificationGroupTypeGroupDefault, Name: "Default"},
					{ID: "default1", Type: mackerel.NotificationGroupTypeGroupDefault, Name: "Default"},
				}, nil
			},
			wantErr: true,
		},
	}

	ctx := context.Background()
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := readDefaultNotificationGroupInner(ctx, tt.inClient)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %+v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if diff := cmp.Diff(got, tt.wants); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func Test_DefaultNotificationGroup_UpdatePreservesUnmanagedFields(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := &defaultNotificationGroupUpdaterTester{
		Groups: []*mackerel.NotificationGroup{{
			ID:                        "default",
			Type:                      mackerel.NotificationGroupTypeGroupDefault,
			Name:                      "Default",
			NotificationLevel:         mackerel.NotificationLevelCritical,
			ChildNotificationGroupIDs: []string{"child-ng"},
			ChildChannelIDs:           []string{"old-channel"},
			Monitors: []*mackerel.NotificationGroupMonitor{{
				ID:          "monitor",
				SkipDefault: true,
			}},
			Services: []*mackerel.NotificationGroupService{{Name: "service"}},
		}},
	}
	model := DefaultNotificationGroupModel{
		NotificationLevel:         types.StringValue("all"),
		ChildNotificationGroupIDs: []types.String{types.StringValue("new-child-ng")},
		ChildChannelIDs:           []types.String{types.StringValue("new-channel")},
	}

	if err := model.updateInner(ctx, client); err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	wantReq := mackerel.NotificationGroup{
		ID:                        "default",
		Name:                      "Default",
		NotificationLevel:         mackerel.NotificationLevelAll,
		ChildNotificationGroupIDs: []string{"new-child-ng"},
		ChildChannelIDs:           []string{"new-channel"},
		Monitors: []*mackerel.NotificationGroupMonitor{{
			ID:          "monitor",
			SkipDefault: true,
		}},
		Services: []*mackerel.NotificationGroupService{{Name: "service"}},
	}
	if client.ID != "default" {
		t.Fatalf("updated ID = %q, want %q", client.ID, "default")
	}
	if diff := cmp.Diff(client.Request, wantReq); diff != "" {
		t.Error(diff)
	}
	if got := model.ID.ValueString(); got != "default" {
		t.Fatalf("model ID = %q, want %q", got, "default")
	}
}

func Test_DefaultNotificationGroup_DeleteDoesNotCallAPI(t *testing.T) {
	t.Parallel()

	var model DefaultNotificationGroupModel
	if err := model.Delete(context.Background(), nil); err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
}

type defaultNotificationGroupUpdaterTester struct {
	Groups  []*mackerel.NotificationGroup
	ID      string
	Request mackerel.NotificationGroup
}

func (ut *defaultNotificationGroupUpdaterTester) FindNotificationGroups() ([]*mackerel.NotificationGroup, error) {
	return ut.Groups, nil
}

func (ut *defaultNotificationGroupUpdaterTester) UpdateNotificationGroup(id string, param *mackerel.NotificationGroup) (*mackerel.NotificationGroup, error) {
	ut.ID = id
	ut.Request = *param
	data := *param
	return &data, nil
}
