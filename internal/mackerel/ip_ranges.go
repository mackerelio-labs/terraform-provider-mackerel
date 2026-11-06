package mackerel

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IPRangesURL is the endpoint of the Mackerel Metadata API that provides IP address ranges.
// See https://mackerel.io/docs/entry/spec/metadata
const IPRangesURL = "https://mackerel.io/meta/ip-ranges.json"

type IPRangesModel struct {
	ID           types.String   `tfsdk:"id"`
	Incoming     []types.String `tfsdk:"incoming"`
	IncomingIPv6 []types.String `tfsdk:"incoming_ipv6"`
	Outgoing     []types.String `tfsdk:"outgoing"`
	OutgoingIPv6 []types.String `tfsdk:"outgoing_ipv6"`
}

type ipRangesResponse struct {
	Incoming     []string `json:"incoming"`
	IncomingIPv6 []string `json:"incoming_ipv6"`
	Outgoing     []string `json:"outgoing"`
	OutgoingIPv6 []string `json:"outgoing_ipv6"`
}

// ReadIPRanges fetches the IP address ranges used by Mackerel.
// The Metadata API does not require an API key, so the client is used only for
// its HTTP transport and User-Agent when available.
func ReadIPRanges(ctx context.Context, client *Client) (IPRangesModel, error) {
	httpClient := http.DefaultClient
	userAgent := ""
	if client != nil {
		if client.HTTPClient != nil {
			httpClient = client.HTTPClient
		}
		userAgent = client.UserAgent
	}
	return readIPRangesInner(ctx, httpClient, IPRangesURL, userAgent)
}

func readIPRangesInner(ctx context.Context, httpClient *http.Client, url, userAgent string) (IPRangesModel, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return IPRangesModel{}, err
	}
	req.Header.Set("Accept", "application/json")
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return IPRangesModel{}, err
	}
	defer resp.Body.Close() // nolint

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return IPRangesModel{}, fmt.Errorf("unexpected status code from %s: %d %s", url, resp.StatusCode, string(body))
	}

	var r ipRangesResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return IPRangesModel{}, fmt.Errorf("failed to decode response from %s: %w", url, err)
	}

	return IPRangesModel{
		ID:           types.StringValue(url),
		Incoming:     toSortedStringValues(r.Incoming),
		IncomingIPv6: toSortedStringValues(r.IncomingIPv6),
		Outgoing:     toSortedStringValues(r.Outgoing),
		OutgoingIPv6: toSortedStringValues(r.OutgoingIPv6),
	}, nil
}

// toSortedStringValues returns the values sorted so that the result is stable
// regardless of the order returned by the API.
func toSortedStringValues(ss []string) []types.String {
	sorted := slices.Sorted(slices.Values(ss))
	vs := make([]types.String, 0, len(sorted))
	for _, s := range sorted {
		vs = append(vs, types.StringValue(s))
	}
	return vs
}
