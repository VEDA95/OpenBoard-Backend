package board

import (
	"VEDA95/open_board/api/internal/auth"
	"time"

	"github.com/huandu/go-sqlbuilder"
)

var boardColumns = []string{
	"open_board_board.id AS board_id",
	"open_board_board.date_created AS board_date_created",
	"open_board_board.date_updated AS board_date_updated",
	"open_board_board.name AS board_name",
	"open_board_board.is_public AS board_is_public",
	"open_board_workspace.id AS workspace_id",
	"open_board_workspace.date_created AS workspace_date_created",
	"open_board_workspace.date_updated AS workspace_date_updated",
	"open_board_workspace.name AS workspace_name",
	"open_board_workspace.description AS workspace_description",
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
	ID          string     `db:"workspace_id"`
	Name        string     `db:"workspace_name"`
	DateCreated time.Time  `db:"workspace_date_created"`
	DateUpdated *time.Time `db:"workspace_date_updated,omitempty"`
	Description *string    `db:"workspace_description,omitempty"`
	User        *auth.User `db:""`
}

type Board struct {
	ID          string     `db:"id"`
	Name        string     `db:"name"`
	DateCreated time.Time  `db:"date_created"`
	DateUpdated *time.Time `db:"date_updated,omitempty"`
	IsPublic    bool       `db:"is_public"`
	Workspace   *Workspace `db:""`
	User        *auth.User `db:""`
}

func BoardsQuery() *sqlbuilder.SelectBuilder {
	return sqlbuilder.Select(boardColumns...).
		From("open_board_board").
		Join("open_board_workspace", "open_board_board.workspace_id = open_board_workspace.id").
		Join("open_board_user", "open_board_board.user_id = open_board_user.id")
}

func GetWorkspaces() []Workspace {}

func GetBoards() []Board {}

func GetWorkspace(ID string) *Workspace {}

func GetBoard(ID string) *Board {}
