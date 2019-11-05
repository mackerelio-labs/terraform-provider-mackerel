package mackerel

import (
	"net/http"
)

func closeResponse(resp *http.Response) {
	if resp != nil {
		_ = resp.Body.Close()
	}
}

func expandStringSlice(s []interface{}) []string {
	vs := make([]string, 0, len(s))
	for _, v := range s {
		vs = append(vs, v.(string))
	}
	return vs
}
