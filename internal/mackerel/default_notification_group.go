package mackerel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

type DefaultNotificationGroupModel struct {
	ID                        types.String   `tfsdk:"id"`
	NotificationLevel         types.String   `tfsdk:"notification_level"`
	ChildNotificationGroupIDs []types.String `tfsdk:"child_notification_group_ids"`
	ChildChannelIDs           []types.String `tfsdk:"child_channel_ids"`
}

func readDefaultNotificationGroupInner(ctx context.Context, client notificationGroupFinder) (DefaultNotificationGroupModel, error) {
	ng, err := findDefaultNotificationGroup(ctx, client)
	if err != nil {
		return DefaultNotificationGroupModel{}, err
	}
	return newDefaultNotificationGroupModel(*ng), nil
}

func findDefaultNotificationGroup(_ context.Context, client notificationGroupFinder) (*mackerel.NotificationGroup, error) {
	ngs, err := client.FindNotificationGroups()
	if err != nil {
		return nil, err
	}

	var defaultNG *mackerel.NotificationGroup
	for _, ng := range ngs {
		if ng.Type != mackerel.NotificationGroupTypeGroupDefault {
			continue
		}
		if defaultNG != nil {
			return nil, fmt.Errorf("multiple default notification groups found")
		}
		defaultNG = ng
	}
	if defaultNG == nil {
		return nil, fmt.Errorf("default notification group is not found")
	}
	return defaultNG, nil
}

func (m *DefaultNotificationGroupModel) Create(ctx context.Context, client *Client) error {
	return m.updateInner(ctx, client)
}

func (m *DefaultNotificationGroupModel) Read(ctx context.Context, client *Client) error {
	data, err := readDefaultNotificationGroupInner(ctx, client)
	if err != nil {
		return err
	}
	*m = data
	return nil
}

type defaultNotificationGroupUpdater interface {
	notificationGroupFinder
	UpdateNotificationGroup(string, *mackerel.NotificationGroup) (*mackerel.NotificationGroup, error)
}

func (m *DefaultNotificationGroupModel) Update(ctx context.Context, client *Client) error {
	return m.updateInner(ctx, client)
}

func (m *DefaultNotificationGroupModel) updateInner(ctx context.Context, client defaultNotificationGroupUpdater) error {
	ng, err := findDefaultNotificationGroup(ctx, client)
	if err != nil {
		return err
	}

	param := *ng
	param.Type = ""
	param.NotificationLevel = mackerel.NotificationLevel(m.NotificationLevel.ValueString())
	param.ChildNotificationGroupIDs = make([]string, 0, len(m.ChildNotificationGroupIDs))
	for _, id := range m.ChildNotificationGroupIDs {
		param.ChildNotificationGroupIDs = append(param.ChildNotificationGroupIDs, id.ValueString())
	}
	param.ChildChannelIDs = make([]string, 0, len(m.ChildChannelIDs))
	for _, id := range m.ChildChannelIDs {
		param.ChildChannelIDs = append(param.ChildChannelIDs, id.ValueString())
	}

	updated, err := client.UpdateNotificationGroup(ng.ID, &param)
	if err != nil {
		return err
	}
	*m = newDefaultNotificationGroupModel(*updated)
	return nil
}

func (m *DefaultNotificationGroupModel) Delete(context.Context, *Client) error {
	return nil
}

func newDefaultNotificationGroupModel(ng mackerel.NotificationGroup) (data DefaultNotificationGroupModel) {
	data.ID = types.StringValue(ng.ID)
	data.NotificationLevel = types.StringValue(string(ng.NotificationLevel))

	data.ChildNotificationGroupIDs = make([]types.String, 0, len(ng.ChildNotificationGroupIDs))
	for _, id := range ng.ChildNotificationGroupIDs {
		data.ChildNotificationGroupIDs = append(data.ChildNotificationGroupIDs, types.StringValue(id))
	}

	data.ChildChannelIDs = make([]types.String, 0, len(ng.ChildChannelIDs))
	for _, id := range ng.ChildChannelIDs {
		data.ChildChannelIDs = append(data.ChildChannelIDs, types.StringValue(id))
	}

	return
}
