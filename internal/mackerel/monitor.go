package mackerel

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/typeutil"
	"github.com/mackerelio/mackerel-client-go"
)

type MonitorModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Memo                 types.String `tfsdk:"memo"`
	IsMute               types.Bool   `tfsdk:"is_mute"`
	NotificationInterval types.Int64  `tfsdk:"notification_interval"`
	// #region: one-of
	HostMetricMonitor       []MonitorHostMetric       `tfsdk:"host_metric"`
	ServiceMetricMonitor    []MonitorServiceMetric    `tfsdk:"service_metric"`
	ExpressionMonitor       []MonitorExpression       `tfsdk:"expression"`
	QueryMonitor            []MonitorQuery            `tfsdk:"query"`
	ConnectivityMonitor     []MonitorConnectivity     `tfsdk:"connectivity"`
	ExternalMonitor         []MonitorExternal         `tfsdk:"external"`
	AnomalyDetectionMonitor []MonitorAnomalyDetection `tfsdk:"anomaly_detection"`
	// #endregion
}

type MonitorHostMetric struct {
	MetricName       types.String         `tfsdk:"metric"`
	Operator         types.String         `tfsdk:"operator"`
	Warning          typeutil.FloatString `tfsdk:"warning"`
	Critical         typeutil.FloatString `tfsdk:"critical"`
	Duration         types.Int64          `tfsdk:"duration"`
	MaxCheckAttempts types.Int64          `tfsdk:"max_check_attempts"`
	Scopes           []string             `tfsdk:"scopes"`
	ExcludeScopes    []string             `tfsdk:"exclude_scopes"`
}

type MonitorServiceMetric struct {
	ServiceName             types.String         `tfsdk:"service"`
	MetricName              types.String         `tfsdk:"metric"`
	Operator                types.String         `tfsdk:"operator"`
	Warning                 typeutil.FloatString `tfsdk:"warning"`
	Critical                typeutil.FloatString `tfsdk:"critical"`
	Duration                types.Int64          `tfsdk:"duration"`
	MaxCheckAttempts        types.Int64          `tfsdk:"max_check_attempts"`
	MissingDurationWarning  types.Int64          `tfsdk:"missing_duration_warning"`
	MissingDurationCritical types.Int64          `tfsdk:"missing_duration_critical"`
}

type MonitorExpression struct {
	Expression types.String         `tfsdk:"expression"`
	Operator   types.String         `tfsdk:"operator"`
	Warning    typeutil.FloatString `tfsdk:"warning"`
	Critical   typeutil.FloatString `tfsdk:"critical"`
}

type MonitorQuery struct {
	Query    types.String         `tfsdk:"query"`
	Legend   types.String         `tfsdk:"legend"`
	Operator types.String         `tfsdk:"operator"`
	Warning  typeutil.FloatString `tfsdk:"warning"`
	Critical typeutil.FloatString `tfsdk:"critical"`
}

type MonitorConnectivity struct {
	Scopes            []string     `tfsdk:"scopes"`
	ExcludeScopes     []string     `tfsdk:"exclude_scopes"`
	AlertStatusOnGone types.String `tfsdk:"alert_status_on_gone"`
}

type MonitorExternal struct {
	Method                          types.String      `tfsdk:"method"`
	URL                             types.String      `tfsdk:"url"`
	MaxCheckAttempts                types.Int64       `tfsdk:"max_check_attempts"`
	ServiceName                     types.String      `tfsdk:"service"`
	ResponseTimeCritical            types.Float64     `tfsdk:"response_time_critical"`
	ResponseTimeWarning             types.Float64     `tfsdk:"response_time_warning"`
	ResponseTimeDuration            types.Int64       `tfsdk:"response_time_duration"`
	RequestBody                     types.String      `tfsdk:"request_body"`
	ContainsString                  types.String      `tfsdk:"contains_string"`
	CertificationExpirationCritical types.Int64       `tfsdk:"certification_expiration_critical"`
	CertificationExpirationWarning  types.Int64       `tfsdk:"certification_expiration_warning"`
	SkipCertificateVerification     types.Bool        `tfsdk:"skip_certificate_verification"`
	Headers                         map[string]string `tfsdk:"headers"`
	FollowRedirect                  types.Bool        `tfsdk:"follow_redirect"`
}

type MonitorAnomalyDetection struct {
	WarningSensitivity  types.String `tfsdk:"warning_sensitivity"`
	CriticalSensitivity types.String `tfsdk:"critical_sensitivity"`
	MaxCheckAttempts    types.Int64  `tfsdk:"max_check_attempts"`
	TrainingPeriodFrom  types.Int64  `tfsdk:"training_period_from"`
	Scopes              []string     `tfsdk:"scopes"`
}

// Reads the monitor by the id.
// Currently non-cancelable.
func ReadMonitor(_ context.Context, client *Client, id string) (MonitorModel, error) {
	m, err := client.GetMonitor(id)
	if err != nil {
		return MonitorModel{}, err
	}
	return newMonitor(m)
}

func (m *MonitorModel) Create(ctx context.Context, client *Client) error {
	monitor, err := client.CreateMonitor(m.mackerelMonitor())
	if err != nil {
		return err
	}
	m.ID = types.StringValue(monitor.MonitorID())
	return nil
}

func (m *MonitorModel) Read(ctx context.Context, client *Client) error {
	remote, err := ReadMonitor(ctx, client, m.ID.ValueString())
	if err != nil {
		return err
	}

	*m = remote

	return nil
}

func (m MonitorModel) Update(ctx context.Context, client *Client) error {
	if _, err := client.UpdateMonitor(m.ID.ValueString(), m.mackerelMonitor()); err != nil {
		return err
	}
	return nil
}

func (m MonitorModel) Delete(ctx context.Context, client *Client) error {
	if _, err := client.DeleteMonitor(m.ID.ValueString()); err != nil {
		return err
	}
	return nil
}

func newMonitor(mackerelMonitor mackerel.Monitor) (MonitorModel, error) {
	model := MonitorModel{
		ID:   types.StringValue(mackerelMonitor.MonitorID()),
		Name: types.StringValue(mackerelMonitor.MonitorName()),
	}
	switch m := mackerelMonitor.(type) {
	case *mackerel.MonitorHostMetric:
		model.Memo = types.StringValue(m.Memo)
		model.IsMute = types.BoolValue(m.IsMute)
		model.NotificationInterval = types.Int64Value(int64(m.NotificationInterval))

		hm := MonitorHostMetric{
			MetricName:       types.StringValue(m.Metric),
			Operator:         types.StringValue(m.Operator),
			Warning:          newFloatStringV0FromFloatPointer(m.Warning),
			Critical:         newFloatStringV0FromFloatPointer(m.Critical),
			Duration:         types.Int64Value(int64(m.Duration)),
			MaxCheckAttempts: types.Int64Value(int64(m.MaxCheckAttempts)),
		}
		hm.Scopes = normalizeScopes(m.Scopes)
		hm.ExcludeScopes = normalizeScopes(m.ExcludeScopes)

		model.HostMetricMonitor = []MonitorHostMetric{hm}
	case *mackerel.MonitorServiceMetric:
		model.Memo = types.StringValue(m.Memo)                                       // "" as default
		model.IsMute = types.BoolValue(m.IsMute)                                     // false as default
		model.NotificationInterval = types.Int64Value(int64(m.NotificationInterval)) // 0 as default (FIXME: 0 is invalid)

		sm := MonitorServiceMetric{
			ServiceName:             types.StringValue(m.Service),
			MetricName:              types.StringValue(m.Metric),
			Operator:                types.StringValue(m.Operator),
			Warning:                 newFloatStringV0FromFloatPointer(m.Warning),
			Critical:                newFloatStringV0FromFloatPointer(m.Critical),
			Duration:                types.Int64Value(int64(m.Duration)),                // 0 as default
			MissingDurationWarning:  types.Int64Value(int64(m.MissingDurationWarning)),  // 0 as default (FIXME: but 0 is invalid)
			MissingDurationCritical: types.Int64Value(int64(m.MissingDurationCritical)), // 0 as default
		}
		if m.MaxCheckAttempts > 0 {
			sm.MaxCheckAttempts = types.Int64Value(int64(m.MaxCheckAttempts))
		}
		model.ServiceMetricMonitor = []MonitorServiceMetric{sm}
	case *mackerel.MonitorExpression:
		model.Memo = types.StringValue(m.Memo)
		model.IsMute = types.BoolValue(m.IsMute)
		model.NotificationInterval = types.Int64Value(int64(m.NotificationInterval))

		em := MonitorExpression{
			Expression: types.StringValue(m.Expression),
			Operator:   types.StringValue(m.Operator),
			Warning:    newFloatStringV0FromFloatPointer(m.Warning),
			Critical:   newFloatStringV0FromFloatPointer(m.Critical),
		}
		model.ExpressionMonitor = []MonitorExpression{em}
	case *mackerel.MonitorQuery:
		model.Memo = types.StringValue(m.Memo)
		model.IsMute = types.BoolValue(m.IsMute)
		model.NotificationInterval = types.Int64Value(int64(m.NotificationInterval))

		qm := MonitorQuery{
			Query:    types.StringValue(m.Query),
			Legend:   types.StringValue(m.Legend),
			Operator: types.StringValue(m.Operator),
			Warning:  newFloatStringV0FromFloatPointer(m.Warning),
			Critical: newFloatStringV0FromFloatPointer(m.Critical),
		}
		model.QueryMonitor = []MonitorQuery{qm}
	case *mackerel.MonitorConnectivity:
		model.Memo = types.StringValue(m.Memo)
		model.IsMute = types.BoolValue(m.IsMute)
		model.NotificationInterval = types.Int64Value(int64(m.NotificationInterval))

		cm := MonitorConnectivity{
			AlertStatusOnGone: types.StringValue(m.AlertStatusOnGone),
		}
		cm.Scopes = normalizeScopes(m.Scopes)
		cm.ExcludeScopes = normalizeScopes(m.ExcludeScopes)
		model.ConnectivityMonitor = []MonitorConnectivity{cm}
	case *mackerel.MonitorExternalHTTP:
		model.Memo = types.StringValue(m.Memo)
		model.IsMute = types.BoolValue(m.IsMute)
		model.NotificationInterval = types.Int64Value(int64(m.NotificationInterval))

		ehm := MonitorExternal{
			Method:                      types.StringValue(m.Method),
			URL:                         types.StringValue(m.URL),
			MaxCheckAttempts:            types.Int64Value(int64(m.MaxCheckAttempts)),
			ServiceName:                 types.StringValue(m.Service),                   // "" as default
			RequestBody:                 types.StringValue(m.RequestBody),               // "" as default
			ContainsString:              types.StringValue(m.ContainsString),            // "" as default
			SkipCertificateVerification: types.BoolValue(m.SkipCertificateVerification), // false as default
			FollowRedirect:              types.BoolValue(m.FollowRedirect),              // false as default
		}
		if m.ResponseTimeCritical != nil {
			ehm.ResponseTimeCritical = types.Float64Value(*m.ResponseTimeCritical)
		} else /* default */ {
			ehm.ResponseTimeCritical = types.Float64Value(0.0)
		}
		if m.ResponseTimeWarning != nil {
			ehm.ResponseTimeWarning = types.Float64Value(*m.ResponseTimeWarning)
		} else /* default */ {
			ehm.ResponseTimeWarning = types.Float64Value(0.0)
		}
		if m.ResponseTimeDuration != nil {
			ehm.ResponseTimeDuration = types.Int64Value(int64(*m.ResponseTimeDuration))
		} else /* default */ {
			ehm.ResponseTimeDuration = types.Int64Value(0)
		}
		if m.CertificationExpirationCritical != nil {
			ehm.CertificationExpirationCritical = types.Int64Value(int64(*m.CertificationExpirationCritical))
		} else /* default */ {
			ehm.CertificationExpirationCritical = types.Int64Value(0)
		}
		if m.CertificationExpirationWarning != nil {
			ehm.CertificationExpirationWarning = types.Int64Value(int64(*m.CertificationExpirationWarning))
		} else /* default */ {
			ehm.CertificationExpirationWarning = types.Int64Value(0)
		}

		if m.Headers != nil {
			headers := make(map[string]string, len(m.Headers))
			for _, h := range m.Headers {
				headers[h.Name] = h.Value
			}
			ehm.Headers = headers
		}

		model.ExternalMonitor = []MonitorExternal{ehm}
	case *mackerel.MonitorAnomalyDetection:
		model.Memo = types.StringValue(m.Memo)
		model.IsMute = types.BoolValue(m.IsMute)
		model.NotificationInterval = types.Int64Value(int64(m.NotificationInterval))

		adm := MonitorAnomalyDetection{
			MaxCheckAttempts:    types.Int64Value(int64(m.MaxCheckAttempts)),
			WarningSensitivity:  types.StringValue(m.WarningSensitivity),       // "" as default
			CriticalSensitivity: types.StringValue(m.CriticalSensitivity),      // "" as default
			TrainingPeriodFrom:  types.Int64Value(int64(m.TrainingPeriodFrom)), // 0 as default
			Scopes:              normalizeScopes(m.Scopes),
		}
		model.AnomalyDetectionMonitor = []MonitorAnomalyDetection{adm}
	default:
		return model, fmt.Errorf("unimplemented type: %s", mackerelMonitor.MonitorType())
	}

	return model, nil
}
func (m MonitorModel) mackerelMonitor() mackerel.Monitor {
	if len(m.HostMetricMonitor) > 0 {
		hm := m.HostMetricMonitor[0]
		return &mackerel.MonitorHostMetric{
			Type:                 "host",
			ID:                   m.ID.ValueString(),
			Name:                 m.Name.ValueString(),
			Memo:                 m.Memo.ValueString(),
			IsMute:               m.IsMute.ValueBool(),
			NotificationInterval: uint64(m.NotificationInterval.ValueInt64()),

			Metric:           hm.MetricName.ValueString(),
			Operator:         hm.Operator.ValueString(),
			Warning:          hm.Warning.ValueFloat64Pointer(),
			Critical:         hm.Critical.ValueFloat64Pointer(),
			Duration:         uint64(hm.Duration.ValueInt64()),
			MaxCheckAttempts: uint64(hm.MaxCheckAttempts.ValueInt64()),

			Scopes:        hm.Scopes,
			ExcludeScopes: hm.ExcludeScopes,
		}
	}
	if len(m.ServiceMetricMonitor) > 0 {
		sm := m.ServiceMetricMonitor[0]
		return &mackerel.MonitorServiceMetric{
			Type:                 "service",
			ID:                   m.ID.ValueString(),
			Name:                 m.Name.ValueString(),
			Memo:                 m.Memo.ValueString(),
			IsMute:               m.IsMute.ValueBool(),
			NotificationInterval: uint64(m.NotificationInterval.ValueInt64()),

			Service:                 sm.ServiceName.ValueString(),
			Metric:                  sm.MetricName.ValueString(),
			Operator:                sm.Operator.ValueString(),
			Warning:                 sm.Warning.ValueFloat64Pointer(),
			Critical:                sm.Critical.ValueFloat64Pointer(),
			Duration:                uint64(sm.Duration.ValueInt64()),
			MaxCheckAttempts:        uint64(sm.MaxCheckAttempts.ValueInt64()),
			MissingDurationWarning:  uint64(sm.MissingDurationWarning.ValueInt64()),
			MissingDurationCritical: uint64(sm.MissingDurationCritical.ValueInt64()),
		}
	}
	if len(m.ExpressionMonitor) > 0 {
		em := m.ExpressionMonitor[0]
		return &mackerel.MonitorExpression{
			Type:                 "expression",
			ID:                   m.ID.ValueString(),
			Name:                 m.Name.ValueString(),
			Memo:                 m.Memo.ValueString(),
			IsMute:               m.IsMute.ValueBool(),
			NotificationInterval: uint64(m.NotificationInterval.ValueInt64()),

			Expression: em.Expression.ValueString(),
			Operator:   em.Operator.ValueString(),
			Warning:    em.Warning.ValueFloat64Pointer(),
			Critical:   em.Critical.ValueFloat64Pointer(),
		}
	}
	if len(m.QueryMonitor) > 0 {
		qm := m.QueryMonitor[0]
		return &mackerel.MonitorQuery{
			Type:                 "query",
			ID:                   m.ID.ValueString(),
			Name:                 m.Name.ValueString(),
			Memo:                 m.Memo.ValueString(),
			IsMute:               m.IsMute.ValueBool(),
			NotificationInterval: uint64(m.NotificationInterval.ValueInt64()),

			Query:    qm.Query.ValueString(),
			Legend:   qm.Legend.ValueString(),
			Operator: qm.Operator.ValueString(),
			Warning:  qm.Warning.ValueFloat64Pointer(),
			Critical: qm.Critical.ValueFloat64Pointer(),
		}
	}
	if len(m.ConnectivityMonitor) > 0 {
		cm := m.ConnectivityMonitor[0]
		return &mackerel.MonitorConnectivity{
			Type:                 "connectivity",
			ID:                   m.ID.ValueString(),
			Name:                 m.Name.ValueString(),
			Memo:                 m.Memo.ValueString(),
			IsMute:               m.IsMute.ValueBool(),
			NotificationInterval: uint64(m.NotificationInterval.ValueInt64()),

			AlertStatusOnGone: cm.AlertStatusOnGone.ValueString(),
			Scopes:            cm.Scopes,
			ExcludeScopes:     cm.ExcludeScopes,
		}
	}
	if len(m.ExternalMonitor) > 0 {
		ehm := m.ExternalMonitor[0]
		mon := mackerel.MonitorExternalHTTP{
			Type:                 "external",
			ID:                   m.ID.ValueString(),
			Name:                 m.Name.ValueString(),
			Memo:                 m.Memo.ValueString(),
			IsMute:               m.IsMute.ValueBool(),
			NotificationInterval: uint64(m.NotificationInterval.ValueInt64()),

			Method:                      ehm.Method.ValueString(),
			URL:                         ehm.URL.ValueString(),
			MaxCheckAttempts:            uint64(ehm.MaxCheckAttempts.ValueInt64()),
			Service:                     ehm.ServiceName.ValueString(),
			RequestBody:                 ehm.RequestBody.ValueString(),
			ContainsString:              ehm.ContainsString.ValueString(),
			SkipCertificateVerification: ehm.SkipCertificateVerification.ValueBool(),
			FollowRedirect:              ehm.FollowRedirect.ValueBool(),
		}
		// Zero values are default, and need to be omitted.
		if resTimeCrit := ehm.ResponseTimeCritical.ValueFloat64(); resTimeCrit != 0.0 {
			mon.ResponseTimeCritical = &resTimeCrit
		}
		if resTimeWarn := ehm.ResponseTimeWarning.ValueFloat64(); resTimeWarn != 0.0 {
			mon.ResponseTimeWarning = &resTimeWarn
		}
		if resTimeDur := ehm.ResponseTimeDuration.ValueInt64(); resTimeDur > 0 {
			resTimeDurU64 := uint64(resTimeDur)
			mon.ResponseTimeDuration = &resTimeDurU64
		}
		if certExpCrit := ehm.CertificationExpirationCritical.ValueInt64(); certExpCrit > 0 {
			certExpCritU64 := uint64(certExpCrit)
			mon.CertificationExpirationCritical = &certExpCritU64
		}
		if certExpWarn := ehm.CertificationExpirationWarning.ValueInt64(); certExpWarn > 0 {
			certExpWarnU64 := uint64(certExpWarn)
			mon.CertificationExpirationWarning = &certExpWarnU64
		}

		// Headers
		fields := make([]mackerel.HeaderField, 0, len(ehm.Headers))
		for name, value := range ehm.Headers {
			fields = append(fields, mackerel.HeaderField{Name: name, Value: value})
		}
		slices.SortStableFunc(fields, func(a, b mackerel.HeaderField) int {
			return strings.Compare(a.Name, b.Name)
		})
		mon.Headers = fields

		return &mon
	}
	if len(m.AnomalyDetectionMonitor) > 0 {
		adm := m.AnomalyDetectionMonitor[0]
		return &mackerel.MonitorAnomalyDetection{
			Type:                 "anomalyDetection",
			ID:                   m.ID.ValueString(),
			Name:                 m.Name.ValueString(),
			Memo:                 m.Memo.ValueString(),
			IsMute:               m.IsMute.ValueBool(),
			NotificationInterval: uint64(m.NotificationInterval.ValueInt64()),

			WarningSensitivity:  adm.WarningSensitivity.ValueString(),
			CriticalSensitivity: adm.CriticalSensitivity.ValueString(),
			MaxCheckAttempts:    uint64(adm.MaxCheckAttempts.ValueInt64()),
			TrainingPeriodFrom:  uint64(adm.TrainingPeriodFrom.ValueInt64()),
			Scopes:              adm.Scopes,
		}
	}

	panic("unimplemented type")
}

func normalizeScopes(scopes []string) []string {
	normalizedScopes := make([]string, 0, len(scopes))
	for _, s := range scopes {
		// API returns `<service>: <role>`
		normalizedScopes = append(normalizedScopes, strings.Replace(s, ": ", ":", 1))
	}
	return normalizedScopes
}

func newFloatStringV0FromFloatPointer(value *float64) typeutil.FloatString {
	if value == nil {
		return typeutil.NewFloatStringValue("")
	}
	return typeutil.NewFloatStringValue(strconv.FormatFloat(*value, 'f', -1, 64))
}
