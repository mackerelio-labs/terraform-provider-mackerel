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

var serviceNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-_]+$`)

func ServiceNameValidator() validator.String {
	return stringvalidator.All(
		stringvalidator.LengthBetween(2, 63),
		stringvalidator.RegexMatches(serviceNameRegex,
			"Must include only alphabets, numbers, hyphen and underscore, and it can not begin a hyphen or underscore"),
	)
}

type ServiceModel = ServiceModelV0

type ServiceModelV0 struct {
	ID   types.String `tfsdk:"id"`
	Name string       `tfsdk:"name"`
	Memo types.String `tfsdk:"memo"`
}

func ImportService(_ context.Context, id string) (ServiceModel, error) {
	return ServiceModelV0{
		ID:   types.StringValue(id),
		Name: id,
	}, nil
}

func ReadService(ctx context.Context, client *Client, name string) (ServiceModel, error) {
	return readServiceInner(ctx, client, name)
}

type serviceFinder interface {
	FindServices() ([]*mackerel.Service, error)
}

func readServiceInner(_ context.Context, client serviceFinder, name string) (ServiceModel, error) {
	services, err := client.FindServices()
	if err != nil {
		return ServiceModel{}, err
	}

	serviceIdx := slices.IndexFunc(services, func(s *mackerel.Service) bool {
		return s.Name == name
	})
	if serviceIdx == -1 {
		return ServiceModel{}, fmt.Errorf("the name '%s' does not match any service in mackerel.io", name)
	}

	service := services[serviceIdx]
	return ServiceModelV0{
		ID:   types.StringValue(service.Name),
		Name: service.Name,
		Memo: types.StringValue(service.Memo),
	}, nil
}

func (m *ServiceModel) Create(_ context.Context, client *Client) error {
	param := mackerel.CreateServiceParam{
		Name: m.Name,
		Memo: m.Memo.ValueString(),
	}

	service, err := client.CreateService(&param)
	if err != nil {
		return err
	}

	m.ID = types.StringValue(service.Name)
	return nil
}

func (m *ServiceModel) Read(ctx context.Context, client *Client) error {
	var name string
	if !m.ID.IsUnknown() {
		name = m.ID.ValueString()
	} else {
		name = m.Name
	}
	remoteData, err := readServiceInner(ctx, client, name)
	if err != nil {
		return err
	}
	*m = remoteData
	return nil
}

func (m ServiceModel) Delete(_ context.Context, client *Client) error {
	if _, err := client.DeleteService(m.ID.ValueString()); err != nil {
		return err
	}
	return nil
}
