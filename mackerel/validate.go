package mackerel

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// https://github.com/DataDog/terraform-provider-datadog/blob/master/datadog/internal/validators/validators.go#L16
// ValidateFloatString makes sure a string can be parsed into a float
func ValidateFloatString(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringMatch(regexp.MustCompile(`\d*(\.\d*)?`), "value must be a float")(v, k)
}
