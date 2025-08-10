package board

import "github.com/huandu/go-sqlbuilder"

var boardColumns = []string{}

func BoardsQuery() *sqlbuilder.SelectBuilder {
	return sqlbuilder.Select(boardColumns...).From("open_board_board")
}
