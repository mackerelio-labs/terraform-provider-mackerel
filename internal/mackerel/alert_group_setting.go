package mackerel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

type AlertGroupSettingModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Memo                 types.String `tfsdk:"memo"`
	ServiceScopes        []string     `tfsdk:"service_scopes"`
	RoleScopes           []string     `tfsdk:"role_scopes"`
	MonitorScopes        []string     `tfsdk:"monitor_scopes"`
	NotificationInterval types.Int64  `tfsdk:"notification_interval"`
}

func ReadAlertGroupSetting(_ context.Context, client *Client, id string) (AlertGroupSettingModel, error) {
	mag, err := client.GetAlertGroupSetting(id)
	if err != nil {
		return AlertGroupSettingModel{}, err
	}
	return newAlertGroupSetting(*mag), nil
}

func (ag *AlertGroupSettingModel) Create(_ context.Context, client *Client) error {
	param := ag.mackerelAlertGroupSetting()
	mag, err := client.CreateAlertGroupSetting(&param)
	if err != nil {
		return err
	}

	ag.ID = types.StringValue(mag.ID)
	return nil
}

func (ag *AlertGroupSettingModel) Read(ctx context.Context, client *Client) error {
	newAg, err := ReadAlertGroupSetting(ctx, client, ag.ID.ValueString())
	if err != nil {
		return err
	}

	ag.merge(newAg)
	return nil
}

func (ag AlertGroupSettingModel) Update(_ context.Context, client *Client) error {
	param := ag.mackerelAlertGroupSetting()
	if _, err := client.UpdateAlertGroupSetting(ag.ID.ValueString(), &param); err != nil {
		return err
	}
	return nil
}

func (ag AlertGroupSettingModel) Delete(_ context.Context, client *Client) error {
	if _, err := client.DeleteAlertGroupSetting(ag.ID.ValueString()); err != nil {
		return err
	}
	return nil
}

func newAlertGroupSetting(ag mackerel.AlertGroupSetting) AlertGroupSettingModel {
	return AlertGroupSettingModel{
		ID:                   types.StringValue(ag.ID),
		Name:                 types.StringValue(ag.Name),
		Memo:                 types.StringValue(ag.Memo),
		ServiceScopes:        ag.ServiceScopes,
		RoleScopes:           normalizeScopes(ag.RoleScopes),
		MonitorScopes:        ag.MonitorScopes,
		NotificationInterval: types.Int64Value(int64(ag.NotificationInterval)),
	}
}

func (ag AlertGroupSettingModel) mackerelAlertGroupSetting() mackerel.AlertGroupSetting {
	return mackerel.AlertGroupSetting{
		ID:                   ag.ID.ValueString(),
		Name:                 ag.Name.ValueString(),
		Memo:                 ag.Memo.ValueString(),
		ServiceScopes:        ag.ServiceScopes,
		RoleScopes:           ag.RoleScopes,
		MonitorScopes:        ag.MonitorScopes,
		NotificationInterval: uint64(ag.NotificationInterval.ValueInt64()),
	}
}

func (ag *AlertGroupSettingModel) merge(newAg AlertGroupSettingModel) {
	// Distinct null and [] by preserving old state
	if len(ag.ServiceScopes) == 0 && len(newAg.ServiceScopes) == 0 {
		newAg.ServiceScopes = ag.ServiceScopes
	}
	if len(ag.RoleScopes) == 0 && len(newAg.RoleScopes) == 0 {
		newAg.RoleScopes = ag.RoleScopes
	}
	if len(ag.MonitorScopes) == 0 && len(newAg.MonitorScopes) == 0 {
		newAg.MonitorScopes = ag.MonitorScopes
	}
	*ag = newAg
}
