package mackerel

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mackerelio/mackerel-client-go"
)

// ErrNotFound indicates that the requested resource does not exist in Mackerel.
// Resource implementations should return this from their Read methods so that the provider can remove the resource from the state instead of failing.
var ErrNotFound = errors.New("not found")

// wrapErrNotFoundIfHttp404 wraps err with ErrNotFound if err represents a 404 response from the Mackerel API.
func wrapErrNotFoundIfHttp404(err error) error {
	var apiErr *mackerel.APIError
	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
		return fmt.Errorf("%w: %w", ErrNotFound, err)
	}
	return err
}
