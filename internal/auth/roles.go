package auth

import (
	"VEDA95/open_board/api/internal/db"
	"errors"
	"github.com/huandu/go-sqlbuilder"
	"time"
)

var RolePermissionsQueryColumns = []string{
	"open_board_role.id AS role_identifier",
	"open_board_role.name AS role_name",
	"open_board_role.date_created AS role_date_created",
	"open_board_role_permission.id AS permission_identifier",
	"open_board_role_permission.path AS permission_path",
	"open_board_role_permission.date_created AS permission_date_created",
}

type Role struct {
	Id          string            `json:"id" db:"role_identifier"`
	DateCreated time.Time         `json:"date_created" db:"role_date_created"`
	Name        string            `json:"name" db:"role_name"`
	Permissions []*RolePermission `json:"permissions" db:"permissions"`
}

type Permission struct {
	Id          string    `json:"id" db:"id"`
	DateCreated time.Time `json:"date_created" db:"date_created"`
	Path        string    `json:"path" db:"path"`
}

type RolePermission struct {
	Id          string    `json:"id" db:"permission_identifier"`
	DateCreated time.Time `json:"date_created" db:"permission_date_created"`
	Path        string    `json:"path" db:"permission_path"`
}

func GetRoles() ([]*Role, error) {
	if db.Instance == nil {
		return nil, errors.New("db is not initialized")
	}

	rows := make([]map[string]interface{}, 0)
	rolesQuery := sqlbuilder.Select(RolePermissionsQueryColumns...).From("open_board_role_permissions")
	rolesQuery.
		Join("open_board_role", "open_board_role_permissions.role_id = open_board_role.id").
		Join("open_board_role_permission", "open_board_role_permissions.permission_id = open_board_role_permission.id")

	if err := db.Instance.Many(rolesQuery, &rows); err != nil {
		return nil, err
	}

	roleMap := make(map[string]*Role)
	output := make([]*Role, 0)

	for _, row := range rows {
		roleId := row["role_identifier"].(string)

		if _, ok := roleMap[roleId]; !ok {
			roleMap[roleId] = &Role{
				Id:          roleId,
				DateCreated: row["role_date_created"].(time.Time),
				Name:        row["role_name"].(string),
			}
		}

		roleMap[roleId].Permissions = append(roleMap[roleId].Permissions, &RolePermission{
			Id:          row["permission_identifier"].(string),
			DateCreated: row["permission_date_created"].(time.Time),
			Path:        row["permission_path"].(string),
		})
	}

	for _, role := range roleMap {
		output = append(output, role)
	}

	return output, nil
}

func GetRole(id string) (*Role, error) {
	if db.Instance == nil {
		return nil, errors.New("db is not initialized")
	}

	rows := make([]map[string]interface{}, 0)
	roleQuery := sqlbuilder.Select(RolePermissionsQueryColumns...).From("open_board_role_permissions")
	roleQuery.
		Join("open_board_role", "open_board_role_permissions.role_id = open_board_role.id").
		Join("open_board_role_permission", "open_board_role_permissions.permission_id = open_board_role_permission.id").
		Where(roleQuery.Equal("open_board_role_permissions.role_id", id))

	if err := db.Instance.Many(roleQuery, &rows); err != nil {
		return nil, err
	}

	var output Role

	for _, row := range rows {
		if len(output.Id) == 0 {
			output = Role{
				Id:          row["role_identifier"].(string),
				DateCreated: row["role_date_created"].(time.Time),
				Name:        row["role_name"].(string),
			}
		}

		output.Permissions = append(output.Permissions, &RolePermission{
			Id:          row["permission_identifier"].(string),
			DateCreated: row["permission_date_created"].(time.Time),
			Path:        row["permission_path"].(string),
		})
	}

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
