package mackerel

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

type RoleMetadataModel struct {
	ID           types.String         `tfsdk:"id"`
	ServiceName  types.String         `tfsdk:"service"`
	RoleName     types.String         `tfsdk:"role"`
	Namespace    types.String         `tfsdk:"namespace"`
	MetadataJSON jsontypes.Normalized `tfsdk:"metadata_json"`
}

func roleMetadataID(serviceName, roleName, namespace string) string {
	return fmt.Sprintf("%s:%s/%s", serviceName, roleName, namespace)
}

func parseRoleMetadataID(id string) (serviceName, roleName, namespace string, err error) {
	sn, rest, foundColon := strings.Cut(id, ":")
	rn, ns, foundSlash := strings.Cut(rest, "/")
	if !foundColon || !foundSlash {
		return "", "", "", fmt.Errorf("The ID is expected to have `<service>:<role>/<namespace>` format, but got: '%s'.", id)
	}
	return sn, rn, ns, nil
}

func ReadRoleMetadata(ctx context.Context, client *Client, serviceName, roleName, namespace string) (RoleMetadataModel, error) {
	return readRoleMetadata(client, serviceName, roleName, namespace)
}

type roleMetadataReader interface {
	GetRoleMetaData(serviceName, roleName, namespace string) (*mackerel.RoleMetaDataResp, error)
}

func readRoleMetadata(client roleMetadataReader, serviceName, roleName, namespace string) (RoleMetadataModel, error) {
	metadataResp, err := client.GetRoleMetaData(serviceName, roleName, namespace)
	if err != nil {
		return RoleMetadataModel{}, err
	}

	metadataJSON, err := json.Marshal(metadataResp.RoleMetaData)
	if err != nil {
		return RoleMetadataModel{}, fmt.Errorf("failed to marshal result: %w", err)
	}

	id := roleMetadataID(serviceName, roleName, namespace)
	return RoleMetadataModel{
		ID:           types.StringValue(id),
		ServiceName:  types.StringValue(serviceName),
		RoleName:     types.StringValue(roleName),
		Namespace:    types.StringValue(namespace),
		MetadataJSON: jsontypes.NewNormalizedValue(string(metadataJSON)),
	}, nil
}

func ImportRoleMetadata(id string) (RoleMetadataModel, error) {
	serviceName, roleName, namespace, err := parseRoleMetadataID(id)
	if err != nil {
		return RoleMetadataModel{}, err
	}
	return RoleMetadataModel{
		ID:          types.StringValue(id),
		ServiceName: types.StringValue(serviceName),
		RoleName:    types.StringValue(roleName),
		Namespace:   types.StringValue(namespace),
	}, nil
}

func (m *RoleMetadataModel) Create(ctx context.Context, client *Client) error {
	return m.create(client)
}

func (m *RoleMetadataModel) create(client roleMetadataUpdator) error {
	if err := m.update(client); err != nil {
		return err
	}

	m.ID = types.StringValue(
		roleMetadataID(m.ServiceName.ValueString(), m.RoleName.ValueString(), m.Namespace.ValueString()),
	)

	return nil
}

func (m *RoleMetadataModel) Read(ctx context.Context, client *Client) error {
	data, err := readRoleMetadata(
		client,
		m.ServiceName.ValueString(),
		m.RoleName.ValueString(),
		m.Namespace.ValueString(),
	)
	if err != nil {
		return err
	}

	m.ID = data.ID // computed
	m.MetadataJSON = data.MetadataJSON
	return nil
}

func (m RoleMetadataModel) Update(ctx context.Context, client *Client) error {
	return m.update(client)
}

type roleMetadataUpdator interface {
	PutRoleMetaData(serviceName, roleName, namespace string, metadata mackerel.RoleMetaData) error
}

func (m *RoleMetadataModel) update(client roleMetadataUpdator) error {
	var metadata mackerel.RoleMetaData
	if err := json.Unmarshal([]byte(m.MetadataJSON.ValueString()), &metadata); err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %w", err)
	}
	if err := client.PutRoleMetaData(
		m.ServiceName.ValueString(),
		m.RoleName.ValueString(),
		m.Namespace.ValueString(),
		metadata,
	); err != nil {
		return err
	}
	return nil
}

func (m RoleMetadataModel) Delete(_ context.Context, client *Client) error {
	if err := client.DeleteRoleMetaData(
		m.ServiceName.ValueString(),
		m.RoleName.ValueString(),
		m.Namespace.ValueString(),
	); err != nil {
		return err
	}
	return nil
}
