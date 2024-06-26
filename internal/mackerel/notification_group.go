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

// Reads a notification group by `id` or `name`, and returns merged state.
func ReadNotificationGroup(_ context.Context, client *Client, data NotificationGroupModel) (NotificationGroupModel, error) {
	ngs, err := client.FindNotificationGroups()
	if err != nil {
		return NotificationGroupModel{}, err
	}

	var ngIdx int
	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		id := data.ID.ValueString()
		ngIdx = slices.IndexFunc(ngs, func(ng *mackerel.NotificationGroup) bool {
			return ng.ID == id
		})
		if ngIdx < 0 {
			return NotificationGroupModel{}, fmt.Errorf("the ID '%s' does not match any notification group in mackerel.io", id)
		}
	} else if !data.Name.IsNull() && !data.Name.IsUnknown() {
		name := data.Name.ValueString()
		ngIdx = slices.IndexFunc(ngs, func(ng *mackerel.NotificationGroup) bool {
			return ng.Name == name
		})
		if ngIdx < 0 {
			return NotificationGroupModel{}, fmt.Errorf("the name '%s' does not match any notification group in mackerel.io", name)
		}
	} else {
		return NotificationGroupModel{}, fmt.Errorf("missing name or ID")
	}

	return newNotificationGroupModel(*ngs[ngIdx]), nil
}

// API -> Models
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
