package mackerel

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Test_IPRanges_read(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		status  int
		body    string
		want    func(url string) IPRangesModel
		wantErr bool
	}{
		"success (sorted)": {
			status: http.StatusOK,
			body: `{
  "incoming": ["192.0.2.2/32", "192.0.2.10/32", "192.0.2.1/32"],
  "incoming_ipv6": ["2001:db8::/56"],
  "outgoing": ["198.51.100.1/32"],
  "outgoing_ipv6": ["2001:db8:1::/56"]
}`,
			want: func(url string) IPRangesModel {
				return IPRangesModel{
					ID: types.StringValue(url),
					Incoming: []types.String{
						types.StringValue("192.0.2.1/32"),
						types.StringValue("192.0.2.10/32"),
						types.StringValue("192.0.2.2/32"),
					},
					IncomingIPv6: []types.String{types.StringValue("2001:db8::/56")},
					Outgoing:     []types.String{types.StringValue("198.51.100.1/32")},
					OutgoingIPv6: []types.String{types.StringValue("2001:db8:1::/56")},
				}
			},
		},
		"missing fields": {
			status: http.StatusOK,
			body:   `{"incoming": ["192.0.2.1/32"]}`,
			want: func(url string) IPRangesModel {
				return IPRangesModel{
					ID:           types.StringValue(url),
					Incoming:     []types.String{types.StringValue("192.0.2.1/32")},
					IncomingIPv6: []types.String{},
					Outgoing:     []types.String{},
					OutgoingIPv6: []types.String{},
				}
			},
		},
		"server error": {
			status:  http.StatusInternalServerError,
			body:    `internal server error`,
			wantErr: true,
		},
		"invalid json": {
			status:  http.StatusOK,
			body:    `{`,
			wantErr: true,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got := r.Header.Get("User-Agent"); got != "test-agent" {
					t.Errorf("unexpected User-Agent: %q", got)
				}
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer ts.Close()

			got, err := readIPRangesInner(context.Background(), ts.Client(), ts.URL, "test-agent")
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if diff := cmp.Diff(tt.want(ts.URL), got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
