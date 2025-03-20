package auth

import (
	"VEDA95/open_board/api/internal/db"
	"errors"
	"github.com/gofrs/uuid/v5"
	"github.com/huandu/go-sqlbuilder"
)

var RolePermissionsQueryColumns = []string{
	"open_board_role.id AS role_identifier",
	"open_board_role.name AS role_name",
	"open_board_role_permission.id AS permission_identifier",
	"open_board_role_permission.path AS permission_path",
}

type Role struct {
	Id          string            `json:"id" db:"role_identifier"`
	Name        string            `json:"name" db:"role_name"`
	Permissions []*RolePermission `json:"permissions" db:"permissions"`
}

type Permission struct {
	Id   string `json:"id" db:"id"`
	Path string `json:"path" db:"path"`
}

type RolePermission struct {
	Id   string `json:"id" db:"permission_identifier"`
	Path string `json:"path" db:"permission_path"`
}

func GetRoles() ([]*Role, error) {
	if db.Instance == nil {
		return nil, errors.New("db is not initialized")
	}

	rows := make([]map[string]interface{}, 0)
	rolesQuery := RolesQuery()

	if err := db.Instance.Many(rolesQuery, &rows); err != nil {
		return nil, err
	}

	output := make([]*Role, 0)

	AppendPermissionsToRoles(rows, &output)

	return output, nil
}

func GetRole(id string) (*Role, error) {
	if db.Instance == nil {
		return nil, errors.New("db is not initialized")
	}

	rows := make([]map[string]interface{}, 0)
	roleQuery := RolesQuery()
	roleQuery.Where(roleQuery.Equal("open_board_role.id", id))

	if err := db.Instance.Many(roleQuery, &rows); err != nil {
		return nil, err
	}

	var output Role

	AppendPermissionsToRole(rows, &output)

	return &output, nil
}

func GetPermissions() ([]Permission, error) {
	if db.Instance == nil {
		return nil, errors.New("db is not initialized")
	}

	output := make([]Permission, 0)
	permissionQuery := sqlbuilder.Select("*").From("open_board_role_permission")

	if err := db.Instance.Many(permissionQuery, &output); err != nil {
		return nil, err
	}

	return output, nil
}

func GetPermission(id string) (*Permission, error) {
	if db.Instance == nil {
		return nil, errors.New("db is not initialized")
	}

	var output Permission
	permissionQuery := sqlbuilder.Select("*").From("open_board_role_permission")
	permissionQuery.Where(permissionQuery.Equal("open_board_role_permission.id", id))

	if err := db.Instance.One(permissionQuery, &output); err != nil {
		return nil, err
	}

	return &output, nil
}

func AppendPermissionsToRole(rows []map[string]interface{}, role *Role) {
	for _, row := range rows {
		if len(role.Id) == 0 {
			role.Id = row["role_identifier"].(uuid.UUID).String()
			role.Name = row["role_name"].(string)
		}

		if row["permission_identifier"] != nil {
			role.Permissions = append(role.Permissions, &RolePermission{
				Id:   row["permission_identifier"].(uuid.UUID).String(),
				Path: row["permission_path"].(string),
			})
		}
	}
}

func AppendPermissionsToRoles(rows []map[string]interface{}, roles *[]*Role) {
	roleMap := make(map[string]*Role)

	for _, row := range rows {
		roleId := row["role_identifier"].(uuid.UUID).String()

		if _, ok := roleMap[roleId]; !ok {
			roleMap[roleId] = &Role{
				Id:   roleId,
				Name: row["role_name"].(string),
			}
		}

		if row["permission_identifier"] != nil {
			roleMap[roleId].Permissions = append(roleMap[roleId].Permissions, &RolePermission{
				Id:   row["permission_identifier"].(uuid.UUID).String(),
				Path: row["permission_path"].(string),
			})
		}
	}

	for _, role := range roleMap {
		*roles = append(*roles, role)
	}
}

func RolesQuery() *sqlbuilder.SelectBuilder {
	return sqlbuilder.Select(RolePermissionsQueryColumns...).
		From("open_board_role").
		JoinWithOption(sqlbuilder.LeftJoin, "open_board_role_permissions", "open_board_role.id = open_board_role_permissions.role_id").
		JoinWithOption(sqlbuilder.LeftJoin, "open_board_role_permission", "open_board_role_permissions.permission_id = open_board_role_permission.id")
}
