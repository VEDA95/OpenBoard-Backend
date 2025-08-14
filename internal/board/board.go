package board

import (
	"VEDA95/open_board/api/internal/auth"
	"VEDA95/open_board/api/internal/db"
	"errors"
	"maps"
	"slices"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/huandu/go-sqlbuilder"
)

var BoardColumns = []string{
	"open_board_board.id AS board_id",
	"open_board_board.date_created AS board_date_created",
	"open_board_board.date_updated AS board_date_updated",
	"open_board_board.name AS board_name",
	"open_board_board.is_public AS board_is_public",
	"open_board_user.id",
	"open_board_user.date_created",
	"open_board_user.date_updated",
	"open_board_user.last_login",
	"open_board_user.username",
	"open_board_user.email",
	"open_board_user.first_name",
	"open_board_user.last_name",
	"open_board_user.hashed_password",
	"open_board_user.enabled",
	"open_board_user.email_verified",
}

var WorkspaceColumns = []string{
	"open_board_workspace.id AS workspace_id",
	"open_board_workspace.date_created AS workspace_date_created",
	"open_board_workspace.date_updated AS workspace_date_updated",
	"open_board_workspace.name AS workspace_name",
	"open_board_workspace.description AS workspace_description",
	"open_board_workspace.is_public AS workspace_is_public",
	"open_board_user.id",
	"open_board_user.date_created",
	"open_board_user.date_updated",
	"open_board_user.last_login",
	"open_board_user.username",
	"open_board_user.email",
	"open_board_user.first_name",
	"open_board_user.last_name",
	"open_board_user.hashed_password",
	"open_board_user.enabled",
	"open_board_user.email_verified",
}

type Workspace struct {
	ID          string                `db:"workspace_id"`
	Name        string                `db:"workspace_name"`
	IsPublic    bool                  `db:"workspace_is_public"`
	DateCreated time.Time             `db:"workspace_date_created"`
	DateUpdated *time.Time            `db:"workspace_date_updated,omitempty"`
	Description *string               `db:"workspace_description,omitempty"`
	Permissions []auth.RolePermission `db:"-"`
	Boards      []Board               `db:"-"`
	User        *auth.User            `db:""`
}

type Board struct {
	ID          string     `db:"board_id"`
	Name        string     `db:"board_name"`
	DateCreated time.Time  `db:"board_date_created"`
	DateUpdated *time.Time `db:"board_date_updated,omitempty"`
	IsPublic    bool       `db:"board_is_public"`
	User        *auth.User `db:""`
}

func WorkspaceQuery() *sqlbuilder.SelectBuilder {
	return sqlbuilder.Select(WorkspaceColumns...).
		From("open_board_workspace").
		Join("open_board_user", "open_board_workspace.user_id = open_board_user.id")
}

func BoardsQuery() *sqlbuilder.SelectBuilder {
	return sqlbuilder.Select(BoardColumns...).
		From("open_board_board").
		Join("open_board_user", "open_board_board.user_id = open_board_user.id")
}

func WorkspacePerimssionsQuery() *sqlbuilder.SelectBuilder {
	return sqlbuilder.Select("open_board_role_permission.id AS permission_id", "open_board_role_permission.path AS permission_path").
		From("open_board_workspace_permissions").
		Join("open_board_workspace", "open_board_workspace_permissions.workspace_id = open_board_workspace.id").
		JoinWithOption(sqlbuilder.LeftJoin, "open_board_role_permission", "open_board_workspace_permissions.permission_id = open_board_role_permission.id")
}

func GetWorkspaces() ([]*Workspace, error) {
	if db.Instance == nil {
		return nil, errors.New("database not initialized")
	}

	workspaceQuery := WorkspaceQuery()
	workspaces := make([]*Workspace, 0)

	if err := db.Instance.Many(workspaceQuery, workspaces); err != nil {
		return nil, err
	}

	workspaceIds := make([]any, len(workspaces))
	userIds := make([]any, len(workspaces))
	users := make([]*auth.User, len(workspaces))

	for index, workspace := range workspaces {
		workspaceIds[index] = workspace.ID
		userIds[index] = workspace.User.Id
		users[index] = workspace.User
	}

	workspacePermissionRows := make([]map[string]any, 0)
	workspaceBoardRows := make([]map[string]any, 0)
	userRows := make([]map[string]any, 0)
	workspacePermissionQuery := sqlbuilder.Select(
		"open_board_workspace.id AS workspace_id",
		"open_board_role_permission.id AS permission_id",
		"open_board_role_permission.path AS permission_path",
	).From("open_board_workspace_permissions")
	workspaceBoardRowsQuery := BoardsQuery()
	usersRolesQuery := sqlbuilder.Select(
		"open_board_user_roles.user_id",
		"open_board_role.id AS role_identifier",
		"open_board_role.name AS role_name",
		"open_board_role_permission.id AS permission_identifier",
		"open_board_role_permission.path AS permission_path",
	).From("open_board_user_roles")

	workspaceBoardRowsQuery.Where(workspaceBoardRowsQuery.In("open_board_board.workspace_id", workspaceIds...))
	usersRolesQuery.
		Join("open_board_role", "open_board_user_roles.role_id = open_board_role.id").
		JoinWithOption(sqlbuilder.LeftJoin, "open_board_role_permissions", "open_board_role.id = open_board_role_permissions.role_id").
		JoinWithOption(sqlbuilder.LeftJoin, "open_board_role_permission", "open_board_role_permissions.permission_id = open_board_role_permission.id").
		Where(usersRolesQuery.In("open_board_user_roles.user_id", userIds...))
	workspacePermissionQuery.Where(workspacePermissionQuery.In("open_board_workspace_permissions.workspace_id", workspaceIds...)).
		Join("open_board_workspace", "open_board_workspace_permissions.workspace_id = open_board_workspace.id").
		JoinWithOption(sqlbuilder.LeftJoin, "open_board_role_permission", "open_board_workspace_permissions.permission_id = open_board_role_permission.id")

	if err := db.Instance.Many(workspaceBoardRowsQuery, workspaceBoardRows); err != nil {
		return nil, err
	}

	if err := db.Instance.Many(workspacePermissionQuery, workspacePermissionRows); err != nil {
		return nil, err
	}

	if err := db.Instance.Many(usersRolesQuery, userRows); err != nil {
		return nil, err
	}

	auth.AppendRolesToUsers(userRows, &users)
	AppendPermissionsToWorkspaces(workspacePermissionRows, workspaces)

	return workspaces, nil
}

func GetWorkspace(ID string) (*Workspace, error) {
	if db.Instance == nil {
		return nil, errors.New("database not initialized")
	}

	workspaceQuery := WorkspaceQuery()
	workspace := new(Workspace)

	workspaceQuery.Where(workspaceQuery.Equal("id", ID))

	if err := db.Instance.One(workspaceQuery, workspace); err != nil {
		return nil, err
	}

	workspaceRows := make([]map[string]any, 0)
	userRows := make([]map[string]any, 0)
	workspacePermissionQuery := WorkspacePerimssionsQuery()
	usersRolesQuery := auth.UsersRolesQuery()

	usersRolesQuery.Where(usersRolesQuery.Equal("open_board_user_roles.user_id", workspace.User.Id))
	workspacePermissionQuery.Where(workspacePermissionQuery.Equal("open_board_workspace_permissions.workspace_id", ID))

	if err := db.Instance.Many(workspacePermissionQuery, workspaceRows); err != nil {
		return nil, err
	}

	if err := db.Instance.Many(usersRolesQuery, userRows); err != nil {
		return nil, err
	}

	auth.AppendRolesToUser(userRows, workspace.User)
	AppendPermissionsToWorkspace(workspaceRows, workspace)

	return workspace, nil
}

func GetBoards() ([]Board, error) {
	if db.Instance == nil {
		return nil, errors.New("database not initialized")
	}

	boardQuery := BoardsQuery()
	boards := make([]Board, 0)

	if err := db.Instance.Many(boardQuery, &boards); err != nil {
		return nil, err
	}

	workspaceIds := make([]any, len(boards))
	userIds := make([]any, len(boards))
	workspaceUserIds := make([]any, len(boards))
	workspaces := make([]*Workspace, len(boards))
	users := make([]*auth.User, len(boards))
	workspaceUsers := make([]*auth.User, len(boards))

	for index, board := range boards {
		workspaceIds[index] = board.Workspace.ID
		userIds[index] = board.User.Id
		workspaceUserIds[index] = board.Workspace.User.Id
		workspaces[index] = board.Workspace
		users[index] = board.User
		workspaceUsers[index] = board.Workspace.User
	}

	workspaceRows := make([]map[string]any, 0)
	userRows := make([]map[string]any, 0)
	workspaceUserRows := make([]map[string]any, 0)
	workspacePermissionQuery := sqlbuilder.Select(
		"open_board_workspace.id AS workspace_id",
		"open_board_role_permission.id AS permission_id",
		"open_board_role_permission.path AS permission_path",
	).From("open_board_workspace_permissions")
	usersRolesQuery := sqlbuilder.Select(
		"open_board_user_roles.user_id",
		"open_board_role.id AS role_identifier",
		"open_board_role.name AS role_name",
		"open_board_role_permission.id AS permission_identifier",
		"open_board_role_permission.path AS permission_path",
	).From("open_board_user_roles")
	workspaceUserRolesQuery := sqlbuilder.Select(
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
	workspaceUserRolesQuery.
		Join("open_board_role", "open_board_user_roles.role_id = open_board_role.id").
		JoinWithOption(sqlbuilder.LeftJoin, "open_board_role_permissions", "open_board_role.id = open_board_role_permissions.role_id").
		JoinWithOption(sqlbuilder.LeftJoin, "open_board_role_permission", "open_board_role_permissions.permission_id = open_board_role_permission.id").
		Where(usersRolesQuery.In("open_board_user_roles.user_id", workspaceUserIds...))
	workspacePermissionQuery.Where(workspacePermissionQuery.In("open_board_workspace_permissions.workspace_id", workspaceIds...)).
		Join("open_board_workspace", "open_board_workspace_permissions.workspace_id = open_board_workspace.id").
		JoinWithOption(sqlbuilder.LeftJoin, "open_board_role_permission", "open_board_workspace_permissions.permission_id = open_board_role_permission.id")

	if err := db.Instance.Many(workspacePermissionQuery, workspaceRows); err != nil {
		return nil, err
	}

	if err := db.Instance.Many(usersRolesQuery, userRows); err != nil {
		return nil, err
	}

	if err := db.Instance.Many(workspaceUserRolesQuery, workspaceUserRows); err != nil {
		return nil, err
	}

	auth.AppendRolesToUsers(userRows, &users)
	auth.AppendRolesToUsers(workspaceUserRows, &workspaceUsers)
	AppendPermissionsToWorkspaces(workspaceRows, workspaces)

	return boards, nil
}

func GetBoard(ID string) (*Board, error) {
	if db.Instance == nil {
		return nil, errors.New("database not initialized")
	}

	boardQuery := BoardsQuery()
	board := new(Board)

	boardQuery.Where(boardQuery.Equal("id", ID))

	if err := db.Instance.One(boardQuery, board); err != nil {
		return nil, err
	}

	workspaceRows := make([]map[string]any, 0)
	userRows := make([]map[string]any, 0)
	workspaceUserRows := make([]map[string]any, 0)
	workspacePermissionQuery := sqlbuilder.Select(
		"open_board_workspace.id AS workspace_id",
		"open_board_role_permission.id AS permission_id",
		"open_board_role_permission.path AS permission_path",
	).From("open_board_workspace_permissions")
	usersRolesQuery := sqlbuilder.Select(
		"open_board_user_roles.user_id",
		"open_board_role.id AS role_identifier",
		"open_board_role.name AS role_name",
		"open_board_role_permission.id AS permission_identifier",
		"open_board_role_permission.path AS permission_path",
	).From("open_board_user_roles")
	workspaceUserRolesQuery := sqlbuilder.Select(
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
		Where(usersRolesQuery.Equal("open_board_user_roles.user_id", board.User.Id))
	workspaceUserRolesQuery.
		Join("open_board_role", "open_board_user_roles.role_id = open_board_role.id").
		JoinWithOption(sqlbuilder.LeftJoin, "open_board_role_permissions", "open_board_role.id = open_board_role_permissions.role_id").
		JoinWithOption(sqlbuilder.LeftJoin, "open_board_role_permission", "open_board_role_permissions.permission_id = open_board_role_permission.id").
		Where(usersRolesQuery.Equal("open_board_user_roles.user_id", board.Workspace.User.Id))
	workspacePermissionQuery.Where(workspacePermissionQuery.Equal("open_board_workspace_permissions.workspace_id", board.Workspace.ID)).
		Join("open_board_workspace", "open_board_workspace_permissions.workspace_id = open_board_workspace.id").
		JoinWithOption(sqlbuilder.LeftJoin, "open_board_role_permission", "open_board_workspace_permissions.permission_id = open_board_role_permission.id")

	if err := db.Instance.Many(workspacePermissionQuery, workspaceRows); err != nil {
		return nil, err
	}

	if err := db.Instance.Many(usersRolesQuery, userRows); err != nil {
		return nil, err
	}

	if err := db.Instance.Many(workspaceUserRolesQuery, workspaceUserRows); err != nil {
		return nil, err
	}

	auth.AppendRolesToUser(userRows, board.User)
	auth.AppendRolesToUser(workspaceUserRows, board.Workspace.User)
	AppendPermissionsToWorkspace(workspaceRows, board.Workspace)

	return board, nil
}

func AppendPermissionsToWorkspaces(rows []map[string]any, workspaces []*Workspace) {
	permissionMap := make(map[string]auth.RolePermission)

	for _, row := range rows {
		permissionId := row["permission_id"].(uuid.UUID).String()

		if _, ok := permissionMap[permissionId]; !ok {
			permissionMap[permissionId] = auth.RolePermission{
				Id:   permissionId,
				Path: row["permission_path"].(string),
			}
		}
	}

	for index := range workspaces {
		permissions := make([]auth.RolePermission, 0)

		for _, row := range rows {
			if row["workspace_id"].(uuid.UUID).String() == workspaces[index].ID {
				permissions = append(permissions, permissionMap[row["permission_id"].(uuid.UUID).String()])
			}
		}

		workspaces[index].Permissions = permissions
	}
}

func AppendPermissionsToWorkspace(rows []map[string]any, workspace *Workspace) {
	permissionMap := make(map[string]auth.RolePermission)

	for _, row := range rows {
		permissionId := row["permission_id"].(uuid.UUID).String()

		if _, ok := permissionMap[permissionId]; !ok {
			permissionMap[permissionId] = auth.RolePermission{
				Id:   permissionId,
				Path: row["permission_path"].(string),
			}
		}
	}

	workspace.Permissions = append(workspace.Permissions, slices.Collect(maps.Values(permissionMap))...)
}

func AppendBoardsToWorkspaces(rows []map[string]any, workspaces []*Workspace) {
	boardMap := make(map[string]Board)

	for _, row := range rows {
		boardId := row["board_id"].(uuid.UUID).String()

		if _, ok := boardMap[boardId]; !ok {
			dateCreated, err := time.Parse(time.DateTime, row["board_date_created"].(string))
			if err != nil {
				return
			}

			dateUpdated, err := time.Parse(time.DateTime, row["board_date_updated"].(string))
			if err != nil {
				return
			}

			boardMap[boardId] = Board{
				ID:          boardId,
				Name:        row["board_name"].(string),
				DateCreated: dateCreated,
				DateUpdated: &dateUpdated,
				IsPublic:    row["is_public"].(bool),
				User:        row["user"].(*auth.User),
			}
		}
	}

	for index := range workspaces {
		boards := make([]Board, 0)

		for _, row := range rows {
			if row["workspace_id"].(uuid.UUID).String() == workspaces[index].ID {
				boards = append(boards, boardMap[row["board_id"].(uuid.UUID).String()])
			}
		}

		workspaces[index].Boards = boards
	}
}
