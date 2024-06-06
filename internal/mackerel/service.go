package mackerel

import (
	"context"
	"fmt"
	"regexp"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

func ServiceNameValidator() validator.String {
	return stringvalidator.All(
		stringvalidator.LengthBetween(2, 63),
		stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-_]+$`),
			"Must include only alphabets, numbers, hyphen and underscore, and it can not begin a hyphen or underscore"),
	)
}

type ServiceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Memo types.String `tfsdk:"memo"`
}

func ReadService(_ context.Context, client *Client, name string) (*ServiceModel, error) {
	services, err := client.FindServices()
	if err != nil {
		return nil, err
	}

	serviceIdx := slices.IndexFunc(services, func(s *mackerel.Service) bool {
		return s.Name == name
	})
	if serviceIdx == -1 {
		return nil, fmt.Errorf("the name '%s' does not match any service in mackerel.io", name)
	}

	service := services[serviceIdx]
	return &ServiceModel{
		ID:   types.StringValue(service.Name),
		Name: types.StringValue(service.Name),
		Memo: types.StringValue(service.Memo),
	}, nil
}

func (m *ServiceModel) Set(newData ServiceModel) {
	if !newData.ID.IsUnknown() {
		m.ID = newData.ID
	}
	m.Name = newData.Name
	if newData.Memo.ValueString() != "" || m.Memo.IsUnknown() {
		m.Memo = newData.Memo
	}
}

func (m *ServiceModel) Create(_ context.Context, client *Client) error {
	param := mackerel.CreateServiceParam{
		Name: m.Name.ValueString(),
		Memo: m.Memo.ValueString(),
	}

	service, err := client.CreateService(&param)
	if err != nil {
		return err
	}

	m.ID = types.StringValue(service.Name)
	return nil
}

func (m *ServiceModel) Delete(_ context.Context, client *Client) error {
	if _, err := client.DeleteService(m.ID.ValueString()); err != nil {
		return err
	}
	return nil
}
