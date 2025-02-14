package mackerel

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

type NotificationGroupModel struct {
	ID                        types.String                     `tfsdk:"id"`
	Name                      types.String                     `tfsdk:"name"`
	NotificationLevel         types.String                     `tfsdk:"notification_level"`
	ChildNotificationGroupIDs []types.String                   `tfsdk:"child_notification_group_ids"`
	ChildChannelIDs           []types.String                   `tfsdk:"child_channel_ids"`
	Monitors                  []NotificationTargetMonitorModel `tfsdk:"monitor"`
	Services                  []NotificationTargetServiceModel `tfsdk:"service"`
}
type NotificationTargetMonitorModel struct {
	ID          types.String `tfsdk:"id"`
	SkipDefault types.Bool   `tfsdk:"skip_default"`
}
type NotificationTargetServiceModel struct {
	Name types.String `tfsdk:"name"`
}

func NotificationLevelValidator() validator.String {
	return stringvalidator.OneOf(
		string(mackerel.NotificationLevelAll),
		string(mackerel.NotificationLevelCritical),
	)
}

// Reads a notification group by `id`
func ReadNotificationGroup(ctx context.Context, client *Client, id string) (NotificationGroupModel, error) {
	return readNotificationGroupInner(ctx, client, id)
}

type notificationGroupFinder interface {
	FindNotificationGroups() ([]*mackerel.NotificationGroup, error)
}

func readNotificationGroupInner(_ context.Context, client notificationGroupFinder, id string) (NotificationGroupModel, error) {
	ngs, err := client.FindNotificationGroups()
	if err != nil {
		return NotificationGroupModel{}, err
	}

	ngIdx := slices.IndexFunc(ngs, func(ng *mackerel.NotificationGroup) bool {
		return ng.ID == id
	})
	if ngIdx < 0 {
		return NotificationGroupModel{}, fmt.Errorf("the ID '%s' does not match any notification group in mackerel.io", id)
	}

	return newNotificationGroupModel(*ngs[ngIdx]), nil
}

// Creates a notification group
func (m *NotificationGroupModel) Create(ctx context.Context, client *Client) error {
	return m.createInner(ctx, client)
}

type notificationGroupCreator interface {
	CreateNotificationGroup(*mackerel.NotificationGroup) (*mackerel.NotificationGroup, error)
}

func (m *NotificationGroupModel) createInner(_ context.Context, client notificationGroupCreator) error {
	param := m.mackerelNotificationGroup()
	ng, err := client.CreateNotificationGroup(&param)
	if err != nil {
		return err
	}

	m.ID = types.StringValue(ng.ID)
	return nil
}

// Reads the notification group
func (m *NotificationGroupModel) Read(ctx context.Context, client *Client) error {
	data, err := ReadNotificationGroup(ctx, client, m.ID.ValueString())
	if err != nil {
		return err
	}

	// merge states
	if m.ID.ValueString() != data.ID.ValueString() {
		return fmt.Errorf("ID cannot be updated")
	}
	*m = data

	return nil
}

// Updates the notification group
func (m *NotificationGroupModel) Update(_ context.Context, client *Client) error {
	param := m.mackerelNotificationGroup()
	if _, err := client.UpdateNotificationGroup(m.ID.ValueString(), &param); err != nil {
		return err
	}
	return nil
}

// Deletes the notification group
func (m *NotificationGroupModel) Delete(_ context.Context, client *Client) error {
	if _, err := client.DeleteNotificationGroup(m.ID.ValueString()); err != nil {
		return err
	}
	return nil
}

// API -> Model
func newNotificationGroupModel(ng mackerel.NotificationGroup) (data NotificationGroupModel) {
	data.ID = types.StringValue(ng.ID)
	data.Name = types.StringValue(ng.Name)
	data.NotificationLevel = types.StringValue(string(ng.NotificationLevel))

	data.ChildNotificationGroupIDs = make([]types.String, 0, len(ng.ChildNotificationGroupIDs))
	for _, id := range ng.ChildNotificationGroupIDs {
		data.ChildNotificationGroupIDs = append(data.ChildNotificationGroupIDs, types.StringValue(id))
	}

	data.ChildChannelIDs = make([]types.String, 0, len(ng.ChildChannelIDs))
	for _, id := range ng.ChildChannelIDs {
		data.ChildChannelIDs = append(data.ChildChannelIDs, types.StringValue(id))
	}

	data.Monitors = make([]NotificationTargetMonitorModel, 0, len(ng.Monitors))
	for _, monitor := range ng.Monitors {
		data.Monitors = append(data.Monitors, NotificationTargetMonitorModel{
			ID:          types.StringValue(monitor.ID),
			SkipDefault: types.BoolValue(monitor.SkipDefault),
		})
	}

	data.Services = make([]NotificationTargetServiceModel, 0, len(ng.Services))
	for _, service := range ng.Services {
		data.Services = append(data.Services, NotificationTargetServiceModel{
			Name: types.StringValue(service.Name),
		})
	}

	return
}

// Model -> API
func (m NotificationGroupModel) mackerelNotificationGroup() (ng mackerel.NotificationGroup) {
	ng.ID = m.ID.ValueString()
	ng.Name = m.Name.ValueString()
	ng.NotificationLevel = mackerel.NotificationLevel(m.NotificationLevel.ValueString())

	ng.ChildNotificationGroupIDs = make([]string, 0, len(m.ChildNotificationGroupIDs))
	for _, id := range m.ChildNotificationGroupIDs {
		ng.ChildNotificationGroupIDs = append(ng.ChildNotificationGroupIDs, id.ValueString())
	}

	ng.ChildChannelIDs = make([]string, 0, len(m.ChildChannelIDs))
	for _, id := range m.ChildChannelIDs {
		ng.ChildChannelIDs = append(ng.ChildChannelIDs, id.ValueString())
	}

	ng.Monitors = make([]*mackerel.NotificationGroupMonitor, 0, len(m.Monitors))
	for _, monitor := range m.Monitors {
		ng.Monitors = append(ng.Monitors, &mackerel.NotificationGroupMonitor{
			ID:          monitor.ID.ValueString(),
			SkipDefault: monitor.SkipDefault.ValueBool(),
		})
	}

	ng.Services = make([]*mackerel.NotificationGroupService, 0, len(m.Services))
	for _, service := range m.Services {
		ng.Services = append(ng.Services, &mackerel.NotificationGroupService{
			Name: service.Name.ValueString(),
		})
	}

	return
}
