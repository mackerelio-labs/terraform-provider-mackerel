package mackerel

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

type RoleModel struct {
	ID          types.String `tfsdk:"id"`
	ServiceName types.String `tfsdk:"service"`
	RoleName    types.String `tfsdk:"name"`
	Memo        types.String `tfsdk:"memo"`
}

func roleID(serviceName, roleName string) string {
	return fmt.Sprintf("%s:%s", serviceName, roleName)
}

func parseRoleID(id string) (serviceName, roleName string, err error) {
	serviceName, roleName, ok := strings.Cut(id, ":")
	if !ok {
		return "", "", fmt.Errorf("The ID is expected to have `<service>:<name>` format, but got: '%s'.", id)
	}
	return serviceName, roleName, nil
}

var roleNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-_]+$`)

func RoleNameValidator() validator.String {
	return stringvalidator.All(
		stringvalidator.LengthBetween(2, 63),
		stringvalidator.RegexMatches(
			roleNameRegex,
			"it can only contain letters, numbers, hyphens, and underscores and cannot start with a hyphen or underscore",
		),
	)
}

func ReadRole(ctx context.Context, client *Client, serviceName, roleName string) (RoleModel, error) {
	return readRoleInner(ctx, client, serviceName, roleName)
}

type roleFinder interface {
	FindRoles(string) ([]*mackerel.Role, error)
}

func readRoleInner(_ context.Context, client roleFinder, serviceName, roleName string) (RoleModel, error) {
	roles, err := client.FindRoles(serviceName)
	if err != nil {
		return RoleModel{}, err
	}

	roleIdx := slices.IndexFunc(roles, func(r *mackerel.Role) bool {
		return r.Name == roleName
	})
	if roleIdx < 0 {
		return RoleModel{}, fmt.Errorf("the name '%s' does not match any role in mackerel.io", roleName)
	}

	role := roles[roleIdx]
	return RoleModel{
		ID:          types.StringValue(roleID(serviceName, roleName)),
		ServiceName: types.StringValue(serviceName),
		RoleName:    types.StringValue(roleName),
		Memo:        types.StringValue(role.Memo),
	}, nil
}

func (m *RoleModel) Create(_ context.Context, client *Client) error {
	serviceName := m.ServiceName.ValueString()
	if _, err := client.CreateRole(serviceName, &mackerel.CreateRoleParam{
		Name: m.RoleName.ValueString(),
		Memo: m.Memo.ValueString(),
	}); err != nil {
		return err
	}

	m.ID = types.StringValue(roleID(serviceName, m.RoleName.ValueString()))

	return nil
}

func (m *RoleModel) Read(ctx context.Context, client *Client) error {
	return m.readInner(ctx, client)
}
func (m *RoleModel) readInner(ctx context.Context, client roleFinder) error {
	// In ImportState, attributes other than `id` are unset.
	var serviceName, roleName string
	if !m.ID.IsNull() && !m.ID.IsUnknown() {
		s, r, err := parseRoleID(m.ID.ValueString())
		if err != nil {
			return err
		}
		serviceName, roleName = s, r
	} else {
		serviceName = m.ServiceName.ValueString()
		roleName = m.RoleName.ValueString()
	}

	r, err := readRoleInner(ctx, client, serviceName, roleName)
	if err != nil {
		return err
	}

	m.ID = r.ID                   // computed
	m.ServiceName = r.ServiceName // required
	m.RoleName = r.RoleName       // required

	// optional
	if /* preserve null */ !m.Memo.IsNull() || r.Memo.ValueString() != "" {
		m.Memo = r.Memo
	}

	return nil
}

func (m *RoleModel) Delete(_ context.Context, client *Client) error {
	if _, err := client.DeleteRole(
		m.ServiceName.ValueString(),
		m.RoleName.ValueString(),
	); err != nil {
		return err
	}
	return nil
}
