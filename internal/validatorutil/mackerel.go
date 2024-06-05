package validatorutil

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func MackerelServiceName() validator.String {
	return stringvalidator.All(
		stringvalidator.LengthBetween(2, 63),
		stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-_]+$`),
			"Must include only alphabets, numbers, hyphen and underscore, and it can not begin a hyphen or underscore"),
	)
}
