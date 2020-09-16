package mackerel

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/mackerelio/mackerel-client-go"
)

func flattenService(service *mackerel.Service, d *schema.ResourceData) (diags diag.Diagnostics) {
	d.Set("name", service.Name)
	d.Set("memo", service.Memo)
	return diags
}

func flattenServiceMetadata(metadata mackerel.ServiceMetaData, d *schema.ResourceData) (diags diag.Diagnostics) {
	metadataJSON, err := structure.FlattenJsonToString(metadata.(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("metadata_json", metadataJSON)
	return diags
}

func flattenRole(role *mackerel.Role, d *schema.ResourceData) (diags diag.Diagnostics) {
	d.Set("name", role.Name)
	d.Set("memo", role.Memo)
	return diags
}

func flattenRoleMetadata(metadata mackerel.RoleMetaData, d *schema.ResourceData) (diags diag.Diagnostics) {
	metadataJSON, err := structure.FlattenJsonToString(metadata.(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("metadata_json", metadataJSON)
	return diags
}

func flattenMonitor(monitor mackerel.Monitor, d *schema.ResourceData) (diags diag.Diagnostics) {
	if v, ok := monitor.(*mackerel.MonitorHostMetric); ok {
		diags = flattenMonitorHostMetric(v, d)
	}
	if v, ok := monitor.(*mackerel.MonitorConnectivity); ok {
		diags = flattenMonitorConnectivity(v, d)
	}
	if v, ok := monitor.(*mackerel.MonitorServiceMetric); ok {
		diags = flattenMonitorServiceMetric(v, d)
	}
	if v, ok := monitor.(*mackerel.MonitorExternalHTTP); ok {
		diags = flattenMonitorExternalHTTP(v, d)
	}
	if v, ok := monitor.(*mackerel.MonitorExpression); ok {
		diags = flattenMonitorExpression(v, d)
	}
	if v, ok := monitor.(*mackerel.MonitorAnomalyDetection); ok {
		diags = flattenMonitorAnomalyDetection(v, d)
	}
	return diags
}

func flattenMonitorHostMetric(monitor *mackerel.MonitorHostMetric, d *schema.ResourceData) (diags diag.Diagnostics) {
	d.Set("name", monitor.Name)
	d.Set("memo", monitor.Memo)
	d.Set("is_mute", monitor.IsMute)
	d.Set("notification_interval", monitor.NotificationInterval)
	normalizedScopes := make([]string, 0, len(monitor.Scopes))
	for _, s := range monitor.Scopes {
		normalizedScopes = append(normalizedScopes, strings.ReplaceAll(s, " ", ""))
	}
	normalizedExcludeScopes := make([]string, 0, len(monitor.ExcludeScopes))
	for _, s := range monitor.ExcludeScopes {
		normalizedExcludeScopes = append(normalizedExcludeScopes, strings.ReplaceAll(s, " ", ""))
	}
	d.Set("host_metric", []map[string]interface{}{
		{
			"metric":             monitor.Metric,
			"operator":           monitor.Operator,
			"warning":            monitor.Warning,
			"critical":           monitor.Critical,
			"duration":           monitor.Duration,
			"max_check_attempts": monitor.MaxCheckAttempts,
			"scopes":             flattenStringListToSet(normalizedScopes),
			"exclude_scopes":     flattenStringListToSet(normalizedExcludeScopes),
		},
	})
	return diags
}

func flattenMonitorConnectivity(monitor *mackerel.MonitorConnectivity, d *schema.ResourceData) (diags diag.Diagnostics) {
	d.Set("name", monitor.Name)
	d.Set("memo", monitor.Memo)
	d.Set("is_mute", monitor.IsMute)
	d.Set("notification_interval", monitor.NotificationInterval)
	normalizedScopes := make([]string, 0, len(monitor.Scopes))
	for _, s := range monitor.Scopes {
		normalizedScopes = append(normalizedScopes, strings.ReplaceAll(s, " ", ""))
	}
	normalizedExcludeScopes := make([]string, 0, len(monitor.ExcludeScopes))
	for _, s := range monitor.ExcludeScopes {
		normalizedExcludeScopes = append(normalizedExcludeScopes, strings.ReplaceAll(s, " ", ""))
	}
	d.Set("connectivity", []map[string]interface{}{
		{
			"scopes":         flattenStringListToSet(normalizedScopes),
			"exclude_scopes": flattenStringListToSet(normalizedExcludeScopes),
		},
	})
	return diags
}

func flattenMonitorServiceMetric(monitor *mackerel.MonitorServiceMetric, d *schema.ResourceData) (diags diag.Diagnostics) {
	d.Set("name", monitor.Name)
	d.Set("memo", monitor.Memo)
	d.Set("is_mute", monitor.IsMute)
	d.Set("notification_interval", monitor.NotificationInterval)
	d.Set("service_metric", []map[string]interface{}{
		{
			"service":                   monitor.Service,
			"metric":                    monitor.Metric,
			"operator":                  monitor.Operator,
			"warning":                   monitor.Warning,
			"critical":                  monitor.Critical,
			"duration":                  monitor.Duration,
			"max_check_attempts":        monitor.MaxCheckAttempts,
			"missing_duration_warning":  monitor.MissingDurationWarning,
			"missing_duration_critical": monitor.MissingDurationCritical,
		},
	})
	return diags
}

func flattenMonitorExternalHTTP(monitor *mackerel.MonitorExternalHTTP, d *schema.ResourceData) (diags diag.Diagnostics) {
	d.Set("name", monitor.Name)
	d.Set("memo", monitor.Memo)
	d.Set("is_mute", monitor.IsMute)
	d.Set("notification_interval", monitor.NotificationInterval)
	headers := make(map[string]interface{}, len(monitor.Headers))
	for _, f := range monitor.Headers {
		headers[f.Name] = f.Value
	}
	external := map[string]interface{}{
		"method":                            monitor.Method,
		"url":                               monitor.URL,
		"max_check_attempts":                monitor.MaxCheckAttempts,
		"service":                           monitor.Service,
		"response_time_critical":            monitor.ResponseTimeCritical,
		"response_time_warning":             monitor.ResponseTimeWarning,
		"response_time_duration":            monitor.ResponseTimeDuration,
		"request_body":                      monitor.RequestBody,
		"contains_string":                   monitor.ContainsString,
		"certification_expiration_critical": monitor.CertificationExpirationCritical,
		"certification_expiration_warning":  monitor.CertificationExpirationWarning,
		"skip_certificate_verification":     monitor.SkipCertificateVerification,
		"headers":                           headers,
	}
	d.Set("external", []map[string]interface{}{external})
	return diags
}

func flattenMonitorExpression(monitor *mackerel.MonitorExpression, d *schema.ResourceData) (diags diag.Diagnostics) {
	d.Set("name", monitor.Name)
	d.Set("memo", monitor.Memo)
	d.Set("is_mute", monitor.IsMute)
	d.Set("notification_interval", monitor.NotificationInterval)
	d.Set("expression", []map[string]interface{}{
		{
			"expression": monitor.Expression,
			"operator":   monitor.Operator,
			"warning":    monitor.Warning,
			"critical":   monitor.Critical,
		},
	})
	return diags
}

func flattenMonitorAnomalyDetection(monitor *mackerel.MonitorAnomalyDetection, d *schema.ResourceData) (diags diag.Diagnostics) {
	d.Set("name", monitor.Name)
	d.Set("memo", monitor.Memo)
	d.Set("is_mute", monitor.IsMute)
	d.Set("notification_interval", monitor.NotificationInterval)
	normalizedScopes := make([]string, 0, len(monitor.Scopes))
	for _, s := range monitor.Scopes {
		normalizedScopes = append(normalizedScopes, strings.ReplaceAll(s, " ", ""))
	}
	d.Set("anomaly_detection", []map[string]interface{}{
		{
			"warning_sensitivity":  monitor.WarningSensitivity,
			"critical_sensitivity": monitor.CriticalSensitivity,
			"training_period_from": monitor.TrainingPeriodFrom,
			"max_check_attempts":   monitor.MaxCheckAttempts,
			"scopes":               flattenStringListToSet(normalizedScopes),
		},
	})
	return diags
}

func flattenDowntime(downtime *mackerel.Downtime, d *schema.ResourceData) (diags diag.Diagnostics) {
	d.Set("name", downtime.Name)
	d.Set("memo", downtime.Memo)
	d.Set("start", downtime.Start)
	d.Set("duration", downtime.Duration)
	if downtime.Recurrence != nil {
		weekdays := make([]string, 0, len(downtime.Recurrence.Weekdays))
		for _, weekday := range downtime.Recurrence.Weekdays {
			weekdays = append(weekdays, weekday.String())
		}
		d.Set("recurrence", []map[string]interface{}{
			{
				"type":     downtime.Recurrence.Type.String(),
				"interval": downtime.Recurrence.Interval,
				"weekdays": flattenStringListToSet(weekdays),
				"until":    downtime.Recurrence.Until,
			},
		})
	}
	d.Set("service_scopes", flattenStringListToSet(downtime.ServiceScopes))
	d.Set("service_exclude_scopes", flattenStringListToSet(downtime.ServiceExcludeScopes))
	d.Set("role_scopes", flattenStringListToSet(downtime.RoleScopes))
	d.Set("role_exclude_scopes", flattenStringListToSet(downtime.RoleExcludeScopes))
	d.Set("monitor_scopes", flattenStringListToSet(downtime.MonitorScopes))
	d.Set("monitor_exclude_scopes", flattenStringListToSet(downtime.MonitorExcludeScopes))
	return diags
}

func flattenChannel(channel *mackerel.Channel, d *schema.ResourceData) (diags diag.Diagnostics) {
	d.Set("name", channel.Name)
	switch channel.Type {
	case "email":
		d.Set("email", []map[string]interface{}{
			{
				"emails":   flattenStringListToSet(*channel.Emails),
				"user_ids": flattenStringListToSet(*channel.UserIDs),
				"events":   flattenStringListToSet(*channel.Events),
			},
		})
	case "slack":
		mentions := make(map[string]string)
		for k, v := range map[string]string{
			"ok":       channel.Mentions.OK,
			"warning":  channel.Mentions.Warning,
			"critical": channel.Mentions.Critical,
		} {
			if v != "" {
				mentions[k] = v
			}
		}
		d.Set("slack", []map[string]interface{}{
			{
				"url":                 channel.URL,
				"mentions":            mentions,
				"enabled_graph_image": channel.EnabledGraphImage,
				"events":              flattenStringListToSet(*channel.Events),
			},
		})
	case "webhook":
		d.Set("webhook", []map[string]interface{}{
			{
				"url":    channel.URL,
				"events": flattenStringListToSet(*channel.Events),
			},
		})
	}
	return diags
}

func flattenNotificationGroup(group *mackerel.NotificationGroup, d *schema.ResourceData) (diags diag.Diagnostics) {
	d.Set("name", group.Name)
	d.Set("notification_level", group.NotificationLevel)
	d.Set("child_notification_group_ids", flattenStringListToSet(group.ChildNotificationGroupIDs))
	d.Set("child_channel_ids", flattenStringListToSet(group.ChildChannelIDs))
	monitors := make([]interface{}, 0, len(group.Monitors))
	for _, monitor := range group.Monitors {
		monitors = append(monitors, map[string]interface{}{
			"id":           monitor.ID,
			"skip_default": monitor.SkipDefault,
		})
	}
	d.Set("monitor", schema.NewSet(schema.HashResource(monitorResource), monitors))
	services := make([]interface{}, 0, len(group.Services))
	for _, service := range group.Services {
		services = append(services, map[string]interface{}{
			"name": service.Name,
		})
	}
	d.Set("service", schema.NewSet(schema.HashResource(serviceResource), services))
	return diags
}

func flattenAlertGroupSetting(setting *mackerel.AlertGroupSetting, d *schema.ResourceData) (diags diag.Diagnostics) {
	d.Set("name", setting.Name)
	d.Set("memo", setting.Memo)
	d.Set("service_scopes", flattenStringListToSet(setting.ServiceScopes))
	normalizedRoleScopes := make([]string, 0, len(setting.RoleScopes))
	for _, r := range setting.RoleScopes {
		normalizedRoleScopes = append(normalizedRoleScopes, strings.ReplaceAll(r, " ", ""))
	}
	d.Set("role_scopes", flattenStringListToSet(normalizedRoleScopes))
	d.Set("monitor_scopes", flattenStringListToSet(setting.MonitorScopes))
	d.Set("notification_interval", setting.NotificationInterval)
	return diags
}
