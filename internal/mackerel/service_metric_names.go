package mackerel

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ServiceMetricNamesModel struct {
	ID          types.String   `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	MetricNames []types.String `tfsdk:"metric_names"`
	Prefix      types.String   `tfsdk:"prefix"`
}

func ReadServiceMetricNames(ctx context.Context, client *Client, state ServiceMetricNamesModel) (ServiceMetricNamesModel, error) {
	return readServiceMetricNamesInner(ctx, client, state)
}

type serviceMetricNamesReader interface {
	ListServiceMetricNames(string) ([]string, error)
}

func readServiceMetricNamesInner(_ context.Context, client serviceMetricNamesReader, state ServiceMetricNamesModel) (ServiceMetricNamesModel, error) {
	name := state.Name.ValueString()
	prefix := state.Prefix.ValueString()

	data := state
	data.ID = types.StringValue(name + ":" + prefix)

	names, err := client.ListServiceMetricNames(name)
	if err != nil {
		return data, err
	}

	data.MetricNames = make([]types.String, 0, len(names))
	for _, name := range names {
		if strings.HasPrefix(name, prefix) {
			data.MetricNames = append(data.MetricNames, types.StringValue(name))
		}
	}
	return data, nil
}
