package service

import (
	"errors"

	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/http/validators"
)

type BoardService struct {
	boardRepository *repository.BoardRepository
}

func NewBoardService(boardRepository *repository.BoardRepository) *BoardService {
	return &BoardService{
		boardRepository: boardRepository,
	}
}

func (boardService *BoardService) GetBoards() ([]*models.Board, error) {
	boards, err := boardService.boardRepository.FindAll(repository.QueryOptions{
		Preload: []string{"User", "Permissions", "Workspace", "Lists"},
		Omit:    []string{"UserID", "WorkspaceID"},
	})
	if err != nil {
		return nil, err
	}

	return boards, nil
}

func (boardService *BoardService) GetBoardByID(ID string) (*models.Board, error) {
	board, err := boardService.boardRepository.FindByID(ID, repository.QueryOptions{
		Preload: []string{"User", "Permissions", "Workspace", "Lists"},
		Omit:    []string{"UserID", "WorkspaceID"},
	})
	if err != nil {
		return nil, err
	}

	return board, nil
}

func (boardService *BoardService) GetBoardsByUserID(ID string) ([]*models.Board, error) {
	boards, err := boardService.boardRepository.FindByUserID(ID, repository.QueryOptions{
		Preload: []string{"User", "Permissions", "Workspace", "Lists"},
		Omit:    []string{"UserID", "WorkspaceID"},
	})
	if err != nil {
		return nil, err
	}

	return boards, nil
}

func (boardService *BoardService) GetUserAccessibleBoards(user *models.User) ([]*models.Board, error) {
	allBoards, err := boardService.GetBoards()
	if err != nil {
		return nil, err
	}

	if user.IsSuperuser() {
		return allBoards, nil
	}

	boards := make([]*models.Board, 0)

	for _, board := range allBoards {
		if board.IsPublic || board.User.ID == user.ID {
			boards = append(boards, board)
			continue
		}

		permissionPaths := make([]string, len(board.Permissions))

		for index, permission := range board.Permissions {
			permissionPaths[index] = permission.Path
		}

		if board.User.IsAuthorizedPartial(permissionPaths...) {
			boards = append(boards, board)
			continue
		}
	}

	return boards, nil
}

func (boardService *BoardService) GetUserAccessibleBoardByID(ID string, user *models.User) (*models.Board, error) {
	board, err := boardService.GetBoardByID(ID)
	if err != nil {
		return nil, err
	}

	permissionPaths := make([]string, len(board.Permissions))

	for index, permission := range board.Permissions {
		permissionPaths[index] = permission.Path
	}

	if !user.IsSuperuser() && (board.User.ID != user.ID || !user.IsAuthorizedPartial(permissionPaths...)) {
		return nil, errors.New("unauthorized")
	}

	return board, nil
}

func (boardService *BoardService) CreateBoard(data *validators.CreateBoardValidator) (*models.Board, error) {
	if boardService.boardRepository.ExistsForUserByName(data.UserID, data.Name) {
		return nil, errors.New("board already exists")
	}

	board := &models.Board{
		Name:        data.Name,
		WorksapceID: data.WorkspaceID,
		UserID:      data.UserID,
		IsPublic:    *data.IsPublic,
	}

	if data.PermissionIDs != nil && len(*data.PermissionIDs) > 0 {
		err := boardService.boardRepository.CreateWithPermissions(
			board,
			*data.PermissionIDs,
			repository.QueryOptions{
				Preload: []string{"User", "Permissions", "Workspace", "Lists"},
				Omit:    []string{"UserID", "WorkspaceID"},
			},
			repository.QueryOptions{},
		)
		if err != nil {
			return nil, err
		}

		return board, nil
	}

	err := boardService.boardRepository.Create(board, repository.QueryOptions{
		Preload: []string{"User", "Permissions", "Workspace", "Lists"},
		Omit:    []string{"UserID", "WorkspaceID"},
	})
	if err != nil {
		return nil, err
	}

	return board, nil
}

func (boardService *BoardService) UpdateBoard(ID string, data *validators.UpdateBoardValidator) (*models.Board, error) {
	if !boardService.boardRepository.Exists(ID) {
		return nil, errors.New("board does not exist")
	}

	board, err := boardService.boardRepository.FindByID(ID, repository.QueryOptions{
		Preload: []string{"User", "Permissions", "Workspace", "Lists"},
		Omit:    []string{"UserID", "WorkspaceID"},
	})
	if err != nil {
		return nil, err
	}

	if data.Name != nil && *data.Name != board.Name {
		board.Name = *data.Name
	}

	if data.IsPublic != nil && *data.IsPublic != board.IsPublic {
		board.IsPublic = *data.IsPublic
	}

	if data.UserID != nil && *data.UserID != board.UserID {
		board.UserID = *data.UserID
	}

	if data.WorkspaceID != nil && *data.WorkspaceID != board.WorksapceID {
		board.WorksapceID = *data.WorkspaceID
	}

	if data.PermissionIDs != nil && len(*data.PermissionIDs) > 0 {
		err := boardService.boardRepository.UpdateWithPermissions(
			board,
			*data.PermissionIDs,
			repository.QueryOptions{
				Preload: []string{"User", "Permissions", "Workspace", "Lists"},
				Omit:    []string{"UserID", "WorkspaceID"},
			},
			repository.QueryOptions{},
		)
		if err != nil {
			return nil, err
		}

		return board, nil
	}

	err2 := boardService.boardRepository.Update(board, repository.QueryOptions{
		Preload: []string{"User", "Permissions", "Workspace", "Lists"},
		Omit:    []string{"UserID", "WorkspaceID"},
	})

	if err2 != nil {
		return nil, err2
	}

	return board, nil
}

func (boardService *BoardService) DeleteBoard(ID string) error {
	if !boardService.boardRepository.Exists(ID) {
		return errors.New("board does not exist")
	}

	return boardService.boardRepository.Delete(ID)
}
