package mackerel

import (
	"errors"
	"net/http"
	"testing"

	"github.com/mackerelio/mackerel-client-go"
)

func Test_wrapErrNotFoundIfHttp404(t *testing.T) {
	t.Parallel()

	t.Run("404 APIError is wrapped with ErrNotFound", func(t *testing.T) {
		t.Parallel()

		apiErr := &mackerel.APIError{StatusCode: http.StatusNotFound, Message: "not found"}
		err := wrapErrNotFoundIfHttp404(apiErr)
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("expected ErrNotFound, got: %v", err)
		}
		if !errors.Is(err, apiErr) {
			t.Errorf("expected the original error to be preserved, got: %v", err)
		}
	})

	t.Run("non-404 APIError is not wrapped", func(t *testing.T) {
		t.Parallel()

		apiErr := &mackerel.APIError{StatusCode: http.StatusInternalServerError, Message: "server error"}
		if err := wrapErrNotFoundIfHttp404(apiErr); errors.Is(err, ErrNotFound) {
			t.Errorf("expected the error not to be ErrNotFound, got: %v", err)
		}
	})

	t.Run("non-APIError is not wrapped", func(t *testing.T) {
		t.Parallel()

		plainErr := errors.New("some error")
		if err := wrapErrNotFoundIfHttp404(plainErr); errors.Is(err, ErrNotFound) {
			t.Errorf("expected the error not to be ErrNotFound, got: %v", err)
		}
	})
}
