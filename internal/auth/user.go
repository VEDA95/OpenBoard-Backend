package auth

import (
	"VEDA95/open_board/api/internal/db"
	"errors"
	"github.com/gofrs/uuid/v5"
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

func GetUsers() ([]*User, error) {
	if db.Instance == nil {
		return nil, errors.New("database not initialized")
	}

	output := make([]*User, 0)
	usersQuery := UsersQuery()

	if err := db.Instance.Many(usersQuery, &output); err != nil {
		return nil, err
	}

	userIds := make([]interface{}, len(output))

	for index, user := range output {
		userIds[index] = user.Id
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
		Where(usersRolesQuery.In("open_board_user_roles.user_id", userIds...))

	if err := db.Instance.Many(usersRolesQuery, &rows); err != nil {
		return nil, err
	}

	AppendRolesToUsers(rows, &output)

	return output, nil
}

func GetUser(id string) (*User, error) {
	if db.Instance == nil {
		return nil, errors.New("database not initialized")
	}

	var output User
	userQuery := UsersQuery()
	userQuery.Where(userQuery.Equal("id", id))

	if err := db.Instance.One(userQuery, &output); err != nil {
		return nil, err
	}

	rows := make([]map[string]interface{}, 0)
	usersRolesQuery := UsersRolesQuery()
	usersRolesQuery.Where(usersRolesQuery.Equal("open_board_user_roles.user_id", output.Id))

	if err := db.Instance.Many(usersRolesQuery, &rows); err != nil {
		return nil, err
	}

	AppendRolesToUser(rows, &output)

	return &output, nil
}

func UsersQuery() *sqlbuilder.SelectBuilder {
	return sqlbuilder.Select(UserQueryColumns...).From("open_board_user")
}

func UsersRolesQuery() *sqlbuilder.SelectBuilder {
	return sqlbuilder.Select(
		"open_board_role.id AS role_identifier",
		"open_board_role.name AS role_name",
		"open_board_role_permission.id AS permission_identifier",
		"open_board_role_permission.path AS permission_path",
	).From("open_board_user_roles").
		Join("open_board_role", "open_board_user_roles.role_id = open_board_role.id").
		JoinWithOption(sqlbuilder.LeftJoin, "open_board_role_permissions", "open_board_role.id = open_board_role_permissions.role_id").
		JoinWithOption(sqlbuilder.LeftJoin, "open_board_role_permission", "open_board_role_permissions.permission_id = open_board_role_permission.id")
}

func AppendRolesToUsers(rows []map[string]interface{}, users *[]*User) {
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

	for index := range *users {
		roles := make([]*Role, 0)

		for _, row := range rows {
			if row["user_id"].(uuid.UUID).String() == (*users)[index].Id {
				roles = append(roles, roleMap[row["role_identifier"].(uuid.UUID).String()])
			}
		}

		(*users)[index].Roles = roles
	}
}

func AppendRolesToUser(rows []map[string]interface{}, user *User) {
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
		user.Roles = append(user.Roles, role)
	}
}
