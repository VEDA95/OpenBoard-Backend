package auth

import (
	"VEDA95/open_board/api/internal/db"
	"errors"
	"github.com/huandu/go-sqlbuilder"
	"time"
)

type User struct {
	Id             string     `json:"id" db:"id"`
	DateCreated    time.Time  `json:"date_created" db:"date_created"`
	DateUpdated    *time.Time `json:"date_updated" db:"date_updated,omitempty"`
	LastLogin      *time.Time `json:"last_login" db:"last_login,omitempty"`
	Username       string     `json:"username" db:"username"`
	Email          string     `json:"email_address" db:"email"`
	FirstName      *string    `json:"first_name" db:"first_name,omitempty"`
	LastName       *string    `json:"last_name" db:"last_name,omitempty"`
	Enabled        bool       `json:"enabled" db:"enabled"`
	EmailVerified  bool       `json:"email_verified" db:"email_verified"`
	HashedPassword string     `json:"-" db:"hashed_password"`
	Roles          []*Role    `json:"roles" db:"-"`
}

var UserQueryColumns = []string{
	"id",
	"date_created",
	"date_updated",
	"last_login",
	"username",
	"email",
	"first_name",
	"last_name",
	"hashed_password",
	"enabled",
	"email_verified",
}

func GetUsers() ([]User, error) {
	if db.Instance == nil {
		return nil, errors.New("database not initialized")
	}

	output := make([]User, 0)
	usersQuery := sqlbuilder.Select(UserQueryColumns...).From("open_board_user")

	if err := db.Instance.Many(usersQuery, &output); err != nil {
		return nil, err
	}

	userIds := make([]string, len(output))

	for _, user := range output {
		userIds = append(userIds, user.Id)
	}

	rows := make([]map[string]interface{}, 0)
	usersRolesQuery := sqlbuilder.Select(
		"open_board_user_roles.user_id",
		"open_board_role.id AS role_identifier",
		"open_board_role.name AS role_name",
		"open_board_role_permission.id AS permission_identifier",
		"open_board_role_permission.path AS permission_path",
	).From("open_board_user_roles")
	usersRolesQuery.
		Join("open_board_role", "open_board_user_roles.role_id = open_board_role.id").
		JoinWithOption(sqlbuilder.LeftJoin, "open_board_role_permissions", "open_board_role.id = open_board_role_permissions.role_id").
		JoinWithOption(sqlbuilder.LeftJoin, "open_board_role_permission", "open_board_role_permissions.permission_id = open_board_role_permission.id").
		Where(usersRolesQuery.In("open_board_user_roles.user_id", userIds))

	if err := db.Instance.Many(usersRolesQuery, &rows); err != nil {
		return nil, err
	}

	roleMap := make(map[string]*Role)

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

	for _, user := range output {
		roles := make([]*Role, 0)

		for _, row := range rows {
			if row["user_id"].(string) == user.Id {
				roles = append(roles, roleMap[row["role_identifier"].(string)])
			}
		}

		user.Roles = roles
	}

	return output, nil
}

func GetUser(id string) (*User, error) {
	if db.Instance == nil {
		return nil, errors.New("database not initialized")
	}

	var output User
	userQuery := sqlbuilder.Select(UserQueryColumns...).From("open_board_user")
	userQuery.Where(userQuery.Equal("id", id))

	if err := db.Instance.One(userQuery, &output); err != nil {
		return nil, err
	}

	rows := make([]map[string]interface{}, 0)
	usersRolesQuery := sqlbuilder.Select(
		"open_board_role.id AS role_identifier",
		"open_board_role.name AS role_name",
		"open_board_role_permission.id AS permission_identifier",
		"open_board_role_permission.path AS permission_path",
	).From("open_board_user_roles")
	usersRolesQuery.
		Join("open_board_role", "open_board_user_roles.role_id = open_board_role.id").
		JoinWithOption(sqlbuilder.LeftJoin, "open_board_role_permissions", "open_board_role.id = open_board_role_permissions.role_id").
		JoinWithOption(sqlbuilder.LeftJoin, "open_board_role_permission", "open_board_role_permissions.permission_id = open_board_role_permission.id").
		Where(usersRolesQuery.Equal("open_board_user_roles.user_id", output.Id))

	if err := db.Instance.Many(usersRolesQuery, &rows); err != nil {
		return nil, err
	}

	roleMap := make(map[string]*Role)

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
		output.Roles = append(output.Roles, role)
	}

	return &output, nil
}
