package mackerel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

type (
	DashboardModel struct {
		ID        types.String `tfsdk:"id"`
		Title     types.String `tfsdk:"title"`
		Memo      types.String `tfsdk:"memo"`
		URLPath   types.String `tfsdk:"url_path"`
		CreatedAt types.Int64  `tfsdk:"created_at"`
		UpdatedAt types.Int64  `tfsdk:"updated_at"`

		Graph       []DashboardWidgetGraph       `tfsdk:"graph"`
		Value       []DashboardWidgetValue       `tfsdk:"value"`
		Markdown    []DashboardWidgetMarkdown    `tfsdk:"markdown"`
		AlertStatus []DashboardWidgetAlertStatus `tfsdk:"alert_status"`
	}

	DashboardLayout struct {
		X      types.Int64 `tfsdk:"x"`
		Y      types.Int64 `tfsdk:"y"`
		Width  types.Int64 `tfsdk:"width"`
		Height types.Int64 `tfsdk:"height"`
	}
	DashboardWidget struct {
		Title  types.String      `tfsdk:"title"`
		Layout []DashboardLayout `tfsdk:"layout"`
	}

	DashboardWidgetGraph struct {
		DashboardWidget
		Range []DashboardRange `tfsdk:"range"`

		Host       []DashboardGraphHost       `tfsdk:"host"`
		Role       []DashboardGraphRole       `tfsdk:"role"`
		Service    []DashboardGraphService    `tfsdk:"service"`
		Expression []DashboardGraphExpression `tfsdk:"expression"`
		Query      []DashboardGraphQuery      `tfsdk:"query"`
	}
	DashboardRange struct {
		Relative []DashboardRangeRelative `tfsdk:"relative"`
		Absolute []DashboardRangeAbsolute `tfsdk:"absolute"`
	}
	DashboardRangeRelative struct {
		Period types.Int64 `tfsdk:"period"`
		Offset types.Int64 `tfsdk:"offset"`
	}
	DashboardRangeAbsolute struct {
		Start types.Int64 `tfsdk:"start"`
		End   types.Int64 `tfsdk:"end"`
	}
	DashboardGraphHost struct {
		HostID types.String `tfsdk:"host_id"`
		Name   types.String `tfsdk:"name"`
	}
	DashboardGraphRole struct {
		RoleFullname types.String `tfsdk:"role_fullname"`
		Name         types.String `tfsdk:"name"`
		IsStacked    types.Bool   `tfsdk:"is_stacked"`
	}
	DashboardGraphService struct {
		ServiceName types.String `tfsdk:"service_name"`
		Name        types.String `tfsdk:"name"`
	}
	DashboardGraphExpression struct {
		Expression types.String `tfsdk:"expression"`
	}
	DashboardGraphQuery struct {
		Query  types.String `tfsdk:"query"`
		Legend types.String `tfsdk:"legend"`
	}

	DashboardWidgetValue struct {
		DashboardWidget
		Metric       []DashboardMetric `tfsdk:"metric"`
		FractionSize types.Int64       `tfsdk:"fraction_size"`
		Suffix       types.String      `tfsdk:"suffix"`
	}
	DashboardMetric struct {
		Host       []DashboardMetricHost       `tfsdk:"host"`
		Service    []DashboardMetricService    `tfsdk:"service"`
		Expression []DashboardMetricExpression `tfsdk:"expression"`
		Query      []DashboardMetricQuery      `tfsdk:"query"`
	}
	DashboardMetricHost struct {
		HostID types.String `tfsdk:"host_id"`
		Name   types.String `tfsdk:"name"`
	}
	DashboardMetricService struct {
		ServiceName types.String `tfsdk:"service_name"`
		Name        types.String `tfsdk:"name"`
	}
	DashboardMetricExpression struct {
		Expression types.String `tfsdk:"expression"`
	}
	DashboardMetricQuery struct {
		Query  types.String `tfsdk:"query"`
		Legend types.String `tfsdk:"legend"`
	}

	DashboardWidgetMarkdown struct {
		DashboardWidget
		Markdown types.String `tfsdk:"markdown"`
	}

	DashboardWidgetAlertStatus struct {
		DashboardWidget
		RoleFullname types.String `tfsdk:"role_fullname"`
	}
)

func ReadDashboard(_ context.Context, client *Client, id string) (DashboardModel, error) {
	d, err := client.FindDashboard(id)
	if err != nil {
		return DashboardModel{}, err
	}
	return newDashboard(*d)
}

func (d *DashboardModel) Create(_ context.Context, client *Client) error {
	param := d.mackerelDashboard()
	md, err := client.CreateDashboard(&param)
	if err != nil {
		return err
	}
	d.ID = types.StringValue(md.ID)
	d.CreatedAt = types.Int64Value(md.CreatedAt)
	d.UpdatedAt = types.Int64Value(md.UpdatedAt)
	return nil
}

func (d *DashboardModel) Read(ctx context.Context, client *Client) error {
	nd, err := ReadDashboard(ctx, client, d.ID.ValueString())
	if err != nil {
		return err
	}
	*d = nd
	return nil
}

func (d *DashboardModel) Update(_ context.Context, client *Client) error {
	param := d.mackerelDashboard()
	md, err := client.UpdateDashboard(d.ID.ValueString(), &param)
	if err != nil {
		return err
	}
	d.UpdatedAt = types.Int64Value(md.UpdatedAt)
	return nil
}

func (d DashboardModel) Delete(_ context.Context, client *Client) error {
	if _, err := client.DeleteDashboard(d.ID.ValueString()); err != nil {
		return err
	}
	return nil
}

const (
	dashboardWidgetTypeGraph       = "graph"
	dashboardWidgetTypeValue       = "value"
	dashboardWidgetTypeMarkdown    = "markdown"
	dashboardWidgetTypeAlertStatus = "alertStatus"

	dashboardGraphTypeHost       = "host"
	dashboardGraphTypeRole       = "role"
	dashboardGraphTypeService    = "service"
	dashboardGraphTypeExpression = "expression"
	dashboardGraphTypeQuery      = "query"

	dashboardRangeTypeRelative = "relative"
	dashboardRangeTypeAbsolute = "absolute"

	dashboardMetricTypeHost       = "host"
	dashboardMetricTypeService    = "service"
	dashboardMetricTypeExpression = "expression"
	dashboardMetricTypeQuery      = "query"
)

func newDashboard(d mackerel.Dashboard) (DashboardModel, error) {
	m := DashboardModel{
		ID:        types.StringValue(d.ID),
		Title:     types.StringValue(d.Title),
		Memo:      types.StringValue(d.Memo),
		URLPath:   types.StringValue(d.URLPath),
		CreatedAt: types.Int64Value(d.CreatedAt),
		UpdatedAt: types.Int64Value(d.UpdatedAt),

		Graph:       []DashboardWidgetGraph{},
		Value:       []DashboardWidgetValue{},
		Markdown:    []DashboardWidgetMarkdown{},
		AlertStatus: []DashboardWidgetAlertStatus{},
	}

	for _, w := range d.Widgets {
		// unsupported features
		if len(w.ReferenceLines) != 0 {
			return m, fmt.Errorf("referenceLines is unsupported.")
		}
		if len(w.FormatRules) != 0 {
			return m, fmt.Errorf("formatFules is unsupported.")
		}
		switch w.Type {
		case dashboardWidgetTypeGraph:
			wg, err := newDashboardWidgetGraph(w)
			if err != nil {
				return m, err
			}
			m.Graph = append(m.Graph, wg)
		case dashboardWidgetTypeValue:
			wv, err := newDashboardWidgetValue(w)
			if err != nil {
				return m, err
			}
			m.Value = append(m.Value, wv)
		case dashboardWidgetTypeMarkdown:
			wm, err := newDashboardWidgetMarkdown(w)
			if err != nil {
				return m, err
			}
			m.Markdown = append(m.Markdown, wm)
		case dashboardWidgetTypeAlertStatus:
			was, err := newDashboardWidgetAlertStatus(w)
			if err != nil {
				return m, err
			}
			m.AlertStatus = append(m.AlertStatus, was)
		default:
			return m, fmt.Errorf("unsupported widget type: %q", w.Type)
		}
	}

	return m, nil
}

func (d DashboardModel) mackerelDashboard() mackerel.Dashboard {
	md := mackerel.Dashboard{
		ID:        d.ID.ValueString(),
		Title:     d.Title.ValueString(),
		Memo:      d.Memo.ValueString(),
		URLPath:   d.URLPath.ValueString(),
		CreatedAt: d.CreatedAt.ValueInt64(),
		UpdatedAt: d.UpdatedAt.ValueInt64(),
		Widgets:   make([]mackerel.Widget, 0, len(d.Graph)+len(d.Value)+len(d.Markdown)+len(d.AlertStatus)),
	}
	for _, g := range d.Graph {
		md.Widgets = append(md.Widgets, g.mackerelWidget())
	}
	for _, v := range d.Value {
		md.Widgets = append(md.Widgets, v.mackerelWidget())
	}
	for _, m := range d.Markdown {
		md.Widgets = append(md.Widgets, m.mackerelWidget())
	}
	for _, a := range d.AlertStatus {
		md.Widgets = append(md.Widgets, a.mackerelWidget())
	}
	return md
}

func newDashboardWidget(w mackerel.Widget) DashboardWidget {
	return DashboardWidget{
		Title: types.StringValue(w.Title),
		Layout: []DashboardLayout{{
			X:      types.Int64Value(w.Layout.X),
			Y:      types.Int64Value(w.Layout.Y),
			Width:  types.Int64Value(w.Layout.Width),
			Height: types.Int64Value(w.Layout.Height),
		}},
	}
}

func (w DashboardWidget) mackerelWidget() mackerel.Widget {
	return mackerel.Widget{
		Title: w.Title.ValueString(),
		Layout: mackerel.Layout{
			X:      w.Layout[0].X.ValueInt64(),
			Y:      w.Layout[0].Y.ValueInt64(),
			Width:  w.Layout[0].Width.ValueInt64(),
			Height: w.Layout[0].Height.ValueInt64(),
		},
	}
}

func newDashboardWidgetGraph(w mackerel.Widget) (DashboardWidgetGraph, error) {
	if w.Type != dashboardWidgetTypeGraph {
		return DashboardWidgetGraph{}, fmt.Errorf("expect graph widget, but got: %s", w.Type)
	}

	g := DashboardWidgetGraph{
		DashboardWidget: newDashboardWidget(w),
	}

	switch w.Range.Type {
	case dashboardRangeTypeAbsolute:
		g.Range = []DashboardRange{{
			Absolute: []DashboardRangeAbsolute{{
				Start: types.Int64Value(w.Range.Start),
				End:   types.Int64Value(w.Range.End),
			}},
		}}
	case dashboardRangeTypeRelative:
		g.Range = []DashboardRange{{
			Relative: []DashboardRangeRelative{{
				Period: types.Int64Value(w.Range.Period),
				Offset: types.Int64Value(w.Range.Offset),
			}},
		}}
	default:
		return g, fmt.Errorf("unsupported range type: %s", w.Range.Type)
	}

	switch w.Graph.Type {
	case dashboardGraphTypeHost:
		g.Host = []DashboardGraphHost{{
			HostID: types.StringValue(w.Graph.HostID),
			Name:   types.StringValue(w.Graph.Name),
		}}
	case dashboardGraphTypeRole:
		g.Role = []DashboardGraphRole{{
			RoleFullname: types.StringValue(w.Graph.RoleFullName),
			Name:         types.StringValue(w.Graph.Name),
			IsStacked:    types.BoolValue(w.Graph.IsStacked),
		}}
	case dashboardGraphTypeService:
		g.Service = []DashboardGraphService{{
			ServiceName: types.StringValue(w.Graph.ServiceName),
			Name:        types.StringValue(w.Graph.Name),
		}}
	case dashboardGraphTypeExpression:
		g.Expression = []DashboardGraphExpression{{
			Expression: types.StringValue(w.Graph.Expression),
		}}
	case dashboardGraphTypeQuery:
		g.Query = []DashboardGraphQuery{{
			Query:  types.StringValue(w.Graph.Query),
			Legend: types.StringValue(w.Graph.Legend),
		}}
	default:
		return g, fmt.Errorf("unsupported graph type: %s", w.Graph.Type)
	}

	return g, nil
}

func (g DashboardWidgetGraph) mackerelWidget() mackerel.Widget {
	w := g.DashboardWidget.mackerelWidget()
	w.Type = dashboardWidgetTypeGraph

	if len(g.Range) != 1 {
		panic(fmt.Sprintf("expect range length to be 1, but got: %d", len(g.Range)))
	}
	r := g.Range[0]
	if len(g.Range[0].Absolute) == 1 {
		w.Range = mackerel.Range{
			Type:  dashboardRangeTypeAbsolute,
			Start: r.Absolute[0].Start.ValueInt64(),
			End:   r.Absolute[0].End.ValueInt64(),
		}
	} else if len(g.Range[0].Relative) == 1 {
		w.Range = mackerel.Range{
			Type:   dashboardRangeTypeRelative,
			Period: r.Relative[0].Period.ValueInt64(),
			Offset: r.Relative[0].Offset.ValueInt64(),
		}
	} else {
		panic(fmt.Sprintf("invalid range: %+v", g.Range[0]))
	}

	if len(g.Host) == 1 {
		w.Graph = mackerel.Graph{
			Type:   dashboardGraphTypeHost,
			HostID: g.Host[0].HostID.ValueString(),
			Name:   g.Host[0].Name.ValueString(),
		}
	} else if len(g.Role) == 1 {
		w.Graph = mackerel.Graph{
			Type:         dashboardGraphTypeRole,
			RoleFullName: g.Role[0].RoleFullname.ValueString(),
			Name:         g.Role[0].Name.ValueString(),
			IsStacked:    g.Role[0].IsStacked.ValueBool(),
		}
	} else if len(g.Service) == 1 {
		w.Graph = mackerel.Graph{
			Type:        dashboardGraphTypeService,
			ServiceName: g.Service[0].ServiceName.ValueString(),
			Name:        g.Service[0].Name.ValueString(),
		}
	} else if len(g.Expression) == 1 {
		w.Graph = mackerel.Graph{
			Type:       dashboardGraphTypeExpression,
			Expression: g.Expression[0].Expression.ValueString(),
		}
	} else if len(g.Query) == 1 {
		w.Graph = mackerel.Graph{
			Type:   dashboardGraphTypeQuery,
			Query:  g.Query[0].Query.ValueString(),
			Legend: g.Query[0].Legend.ValueString(),
		}
	} else {
		panic(fmt.Sprintf("invalid graph: %+v", g))
	}

	return w
}

func newDashboardWidgetValue(w mackerel.Widget) (DashboardWidgetValue, error) {
	if w.Type != dashboardWidgetTypeValue {
		return DashboardWidgetValue{}, fmt.Errorf("expect value widget, but got: %s", w.Type)
	}
	v := DashboardWidgetValue{
		DashboardWidget: newDashboardWidget(w),
		Metric:          []DashboardMetric{{}},
		FractionSize:    types.Int64PointerValue(w.FractionSize),
		Suffix:          types.StringValue(w.Suffix),
	}
	switch w.Metric.Type {
	case dashboardMetricTypeHost:
		v.Metric[0].Host = []DashboardMetricHost{{
			HostID: types.StringValue(w.Metric.HostID),
			Name:   types.StringValue(w.Metric.Name),
		}}
	case dashboardMetricTypeService:
		v.Metric[0].Service = []DashboardMetricService{{
			ServiceName: types.StringValue(w.Metric.ServiceName),
			Name:        types.StringValue(w.Metric.Name),
		}}
	case dashboardMetricTypeExpression:
		v.Metric[0].Expression = []DashboardMetricExpression{{
			Expression: types.StringValue(w.Metric.Expression),
		}}
	case dashboardMetricTypeQuery:
		v.Metric[0].Query = []DashboardMetricQuery{{
			Query:  types.StringValue(w.Metric.Query),
			Legend: types.StringValue(w.Metric.Legend),
		}}
	default:
		return v, fmt.Errorf("unsupported metric type: %s", w.Metric.Type)
	}
	return v, nil
}

func (v DashboardWidgetValue) mackerelWidget() mackerel.Widget {
	w := v.DashboardWidget.mackerelWidget()
	w.Type = dashboardWidgetTypeValue
	w.FractionSize = v.FractionSize.ValueInt64Pointer()
	w.Suffix = v.Suffix.ValueString()

	if len(v.Metric) != 1 {
		panic(fmt.Sprintf("invalid metric length: %d", len(v.Metric)))
	}
	m := v.Metric[0]
	if len(m.Host) == 1 {
		w.Metric = mackerel.Metric{
			Type:   dashboardMetricTypeHost,
			HostID: m.Host[0].HostID.ValueString(),
			Name:   m.Host[0].Name.ValueString(),
		}
	} else if len(m.Service) == 1 {
		w.Metric = mackerel.Metric{
			Type:        dashboardMetricTypeService,
			ServiceName: m.Service[0].ServiceName.ValueString(),
			Name:        m.Service[0].Name.ValueString(),
		}
	} else if len(m.Expression) == 1 {
		w.Metric = mackerel.Metric{
			Type:       dashboardMetricTypeExpression,
			Expression: m.Expression[0].Expression.ValueString(),
		}
	} else if len(m.Query) == 1 {
		w.Metric = mackerel.Metric{
			Type:   dashboardMetricTypeQuery,
			Query:  m.Query[0].Query.ValueString(),
			Legend: m.Query[0].Legend.ValueString(),
		}
	} else {
		panic(fmt.Sprintf("invalid metric: %+v", m))
	}

	return w
}

func newDashboardWidgetMarkdown(w mackerel.Widget) (DashboardWidgetMarkdown, error) {
	if w.Type != dashboardWidgetTypeMarkdown {
		return DashboardWidgetMarkdown{}, fmt.Errorf("expect markdown widget, but got: %s", w.Type)
	}
	return DashboardWidgetMarkdown{
		DashboardWidget: newDashboardWidget(w),
		Markdown:        types.StringValue(w.Markdown),
	}, nil
}

func (m DashboardWidgetMarkdown) mackerelWidget() mackerel.Widget {
	w := m.DashboardWidget.mackerelWidget()
	w.Type = dashboardWidgetTypeMarkdown
	w.Markdown = m.Markdown.ValueString()
	return w
}

func newDashboardWidgetAlertStatus(w mackerel.Widget) (DashboardWidgetAlertStatus, error) {
	if w.Type != dashboardWidgetTypeAlertStatus {
		return DashboardWidgetAlertStatus{}, fmt.Errorf("expect alet status widget, but got: %s", w.Type)
	}
	return DashboardWidgetAlertStatus{
		DashboardWidget: newDashboardWidget(w),
		RoleFullname:    types.StringValue(w.RoleFullName),
	}, nil
}

func (a DashboardWidgetAlertStatus) mackerelWidget() mackerel.Widget {
	w := a.DashboardWidget.mackerelWidget()
	w.Type = dashboardWidgetTypeAlertStatus
	w.RoleFullName = a.RoleFullname.ValueString()
	return w
}
