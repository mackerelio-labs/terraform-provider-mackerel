package mackerel

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

type ServiceMetadataModel struct {
	ID           types.String         `tfsdk:"id"`
	ServiceName  types.String         `tfsdk:"service"`
	Namespace    types.String         `tfsdk:"namespace"`
	MetadataJSON jsontypes.Normalized `tfsdk:"metadata_json"`
}

func serviceMetadataID(serviceName, namespace string) string {
	return strings.Join([]string{serviceName, namespace}, "/")
}

func parseServiceMetadataID(id string) (serviceName, namespace string, err error) {
	first, last, ok := strings.Cut(id, "/")
	if !ok {
		return "", "", fmt.Errorf("The ID is expected to have `<service_name>/<namespace>` format, but got: '%s'.", id)
	}
	return first, last, nil
}

func ReadServiceMetadata(ctx context.Context, client *Client, data ServiceMetadataModel) (ServiceMetadataModel, error) {
	return readServiceMetadataInner(ctx, client, data)
}

type serviceMetadataGetter interface {
	GetServiceMetaData(string, string) (*mackerel.ServiceMetaDataResp, error)
}

func readServiceMetadataInner(_ context.Context, client serviceMetadataGetter, data ServiceMetadataModel) (ServiceMetadataModel, error) {
	serviceName, namespace, err := data.getID()
	if err != nil {
		return ServiceMetadataModel{}, err
	}

	metadataResp, err := client.GetServiceMetaData(serviceName, namespace)
	if err != nil {
		return ServiceMetadataModel{}, err
	}

	data.ID = types.StringValue(serviceMetadataID(serviceName, namespace))
	data.ServiceName = types.StringValue(serviceName)
	data.Namespace = types.StringValue(namespace)

	if metadataResp.ServiceMetaData == nil {
		if /* expected not to be deleted */ !data.MetadataJSON.IsNull() {
			data.MetadataJSON = jsontypes.NewNormalizedValue("")
		}
		return data, nil
	}

	metadataJSON, err := json.Marshal(metadataResp.ServiceMetaData)
	if err != nil {
		return ServiceMetadataModel{}, fmt.Errorf("failed to marshal result: %w", err)
	}

	data.MetadataJSON = jsontypes.NewNormalizedValue(string(metadataJSON))
	return data, nil
}

func (m *ServiceMetadataModel) Validate(base path.Path) (diags diag.Diagnostics) {
	if m.ID.IsNull() || m.ID.IsUnknown() {
		return
	}
	id := m.ID.ValueString()
	idPath := base.AtName("id")

	serviceName, namespace, err := parseServiceMetadataID(id)
	if err != nil {
		diags.AddAttributeError(
			idPath,
			"Invalid ID",
			err.Error(),
		)
		return
	}

	if !m.ServiceName.IsNull() && !m.ServiceName.IsUnknown() && m.ServiceName.ValueString() != serviceName {
		diags.AddAttributeError(
			idPath,
			"Invalid ID",
			fmt.Sprintf("ID is expected to start with '%s/', but got: '%s'", m.ServiceName.ValueString(), id),
		)
	}
	if !m.Namespace.IsNull() && !m.Namespace.IsUnknown() && m.Namespace.ValueString() != namespace {
		diags.AddAttributeError(
			idPath,
			"Invalid ID",
			fmt.Sprintf("ID is expected to end with '/%s', but got: '%s'", m.Namespace.ValueString(), id),
		)
	}

	return
}

func (m *ServiceMetadataModel) CreateOrUpdateMetadata(_ context.Context, client *Client) error {
	serviceName, namespace, err := m.getID()
	if err != nil {
		return err
	}

	var metadata mackerel.ServiceMetaData

	if err := json.Unmarshal(
		[]byte(m.MetadataJSON.ValueString()), &metadata,
	); err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	if err := client.PutServiceMetaData(serviceName, namespace, metadata); err != nil {
		return err
	}

	m.ID = types.StringValue(serviceMetadataID(serviceName, namespace))
	return nil
}

func (m *ServiceMetadataModel) Delete(_ context.Context, client *Client) error {
	serviceName, namespace, err := m.getID()
	if err != nil {
		return err
	}

	if err := client.DeleteServiceMetaData(serviceName, namespace); err != nil {
		return err
	}

	return nil
}

func (m *ServiceMetadataModel) getID() (serviceName, namespace string, err error) {
	if !m.ID.IsNull() && !m.ID.IsUnknown() {
		return parseServiceMetadataID(m.ID.ValueString())
	}
	return m.ServiceName.ValueString(), m.Namespace.ValueString(), nil
}
