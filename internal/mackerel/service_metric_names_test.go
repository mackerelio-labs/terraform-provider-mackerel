package mackerel

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type serviceMetricNamesReaderFunc func(string) ([]string, error)

func (f serviceMetricNamesReaderFunc) ListServiceMetricNames(name string) ([]string, error) {
	return f(name)
}

func Test_ServiceMetricNames_read(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		inClient serviceMetricNamesReaderFunc
		inState  ServiceMetricNamesModel

		want ServiceMetricNamesModel
	}{
		"success": {
			inClient: func(name string) ([]string, error) {
				if name != "service0" {
					return nil, fmt.Errorf("no service: %s", name)
				}
				return []string{
					"metric",
					"prefixed_metric",
				}, nil
			},
			inState: ServiceMetricNamesModel{
				Name: types.StringValue("service0"),
			},

			want: ServiceMetricNamesModel{
				ID:   types.StringValue("service0:"),
				Name: types.StringValue("service0"),
				MetricNames: []types.String{
					types.StringValue("metric"),
					types.StringValue("prefixed_metric"),
				},
			},
		},
		"prefix": {
			inClient: func(name string) ([]string, error) {
				if name != "service0" {
					return nil, fmt.Errorf("no service: %s", name)
				}
				return []string{
					"metric",
					"prefixed_metric",
				}, nil
			},
			inState: ServiceMetricNamesModel{
				Name:   types.StringValue("service0"),
				Prefix: types.StringValue("prefixed"),
			},

			want: ServiceMetricNamesModel{
				ID:     types.StringValue("service0:prefixed"),
				Name:   types.StringValue("service0"),
				Prefix: types.StringValue("prefixed"),
				MetricNames: []types.String{
					types.StringValue("prefixed_metric"),
				},
			},
		},
	}

	ctx := context.Background()
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			data, err := readServiceMetricNamesInner(ctx, tt.inClient, tt.inState)
			if err != nil {
				t.Errorf("unexpected error: %+v", err)
				return
			}
			if diff := cmp.Diff(tt.want, data); diff != "" {
				t.Error(diff)
			}
		})
	}
}
