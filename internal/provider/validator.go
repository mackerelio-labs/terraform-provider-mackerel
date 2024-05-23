package provider

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type urlSchemeValidator struct {
	validSchemes []string
}

var _ validator.String = (*urlSchemeValidator)(nil)

func IsURLWithHTTPorHTTPS() validator.String {
	return &urlSchemeValidator{validSchemes: []string{"http", "https"}}
}

func (uv *urlSchemeValidator) Description(context.Context) string {
	return fmt.Sprintf("url string with scheme of: %q", strings.Join(uv.validSchemes, ","))
}

func (uv *urlSchemeValidator) MarkdownDescription(context.Context) string {
	return fmt.Sprintf("url string with scheme of: `%q`", strings.Join(uv.validSchemes, "`,`"))
}

func (uv *urlSchemeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	v := req.ConfigValue.ValueString()
	if v == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Empty String",
			"expected url string, but got empty string",
		)
		return
	}

	u, err := url.Parse(v)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid URL",
			fmt.Sprintf("expected url string: %+v", err),
		)
		return
	}

	if u.Host == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"No Host",
			"expected to have a host",
		)
	}

	isSchemeValid := slices.Index(uv.validSchemes, u.Scheme) > 0
	if !isSchemeValid {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Scheme",
			fmt.Sprintf("expected to have a url with scheme of: %q", strings.Join(uv.validSchemes, ",")),
		)
	}
}
