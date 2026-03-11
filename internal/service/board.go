package service

import (
	"errors"

	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/http/validators"

	models "VEDA95/open_board/api/internal/db/model"
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
	boards, err := boardService.boardRepository.FindAll(repository.WithPreload("User", "Permissions", "Workspace", "Lists"), repository.WithOmit("UserID", "WorkspaceID"))
	if err != nil {
		return nil, err
	}

	return boards, nil
}

func (boardService *BoardService) GetBoardByID(ID string) (*models.Board, error) {
	board, err := boardService.boardRepository.FindByID(ID, repository.WithPreload("User", "Permissions", "Workspace", "Lists"), repository.WithOmit("UserID", "WorkspaceID"))
	if err != nil {
		return nil, err
	}

	return board, nil
}

func (boardService *BoardService) GetBoardsByUserID(ID string) ([]*models.Board, error) {
	boards, err := boardService.boardRepository.FindByUserID(ID, repository.WithPreload("User", "Permissions", "Workspace", "Lists"), repository.WithOmit("UserID", "WorkspaceID"))
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
		if board.IsPublic || board.User.ID == user.ID || board.Workspace.User.ID == user.ID {
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

	if !user.IsSuperuser() && (board.User.ID != user.ID || board.Workspace.User.ID != user.ID || !user.IsAuthorizedPartial(permissionPaths...)) {
		return nil, errors.New("unauthorized")
	}

	return board, nil
}

func (boardService *BoardService) CreateBoard(user *models.User, data *validators.CreateBoardValidator) (*models.Board, error) {
	canManage := user.IsSuperuser() || user.IsAuthorized("boards:manage_all")

	if data.UserID != nil && !canManage {
		return nil, errors.New("unauthorized: no permission for creating boards for other users")
	}

	if data.UserID == nil || (data.UserID == nil && canManage) {
		data.UserID = &user.ID
	}

	if boardService.boardRepository.ExistsForUserByName(*data.UserID, data.Name) {
		return nil, errors.New("board already exists")
	}

	board := &models.Board{
		Name:        data.Name,
		WorkspaceID: data.WorkspaceID,
		UserID:      *data.UserID,
		IsPublic:    data.IsPublic,
	}

	if data.PermissionIDs != nil && len(*data.PermissionIDs) > 0 {
		err := boardService.boardRepository.CreateWithPermissions(
			board,
			*data.PermissionIDs,
			repository.WithPreload("User", "Permissions", "Workspace", "Lists"),
		)
		if err != nil {
			return nil, err
		}

		return board, nil
	}

	if err := boardService.boardRepository.Create(board, repository.WithPreload("User", "Permissions", "Workspace", "Lists")); err != nil {
		return nil, err
	}

	return board, nil
}

func (boardService *BoardService) UpdateBoard(ID string, user *models.User, data *validators.UpdateBoardValidator) (*models.Board, error) {
	if !boardService.boardRepository.Exists(ID) {
		return nil, errors.New("board does not exist")
	}

	board, err := boardService.boardRepository.FindByID(ID, repository.WithPreload("User", "Permissions", "Workspace", "Lists"), repository.WithOmit("UserID", "WorkspaceID"))
	if err != nil {
		return nil, err
	}

	insufficientPermission := (board.UserID != user.ID) &&
		(board.Workspace == nil || board.Workspace.UserID != user.ID) &&
		(!user.IsSuperuser() && !user.IsAuthorized("boards:manage_all"))

	if insufficientPermission {
		if len(board.Permissions) > 0 {
			permissionPaths := make([]string, len(board.Permissions))

			for index, permission := range board.Permissions {
				permissionPaths[index] = permission.Path
			}

			if !user.IsAuthorizedPartial(permissionPaths...) {
				return nil, errors.New("unauthorized: no permission to update this board")
			}

		}

		return nil, errors.New("unauthorized: no permission to update this board")
	}

	if data.Name != nil && *data.Name != board.Name {
		board.Name = *data.Name
	}

	if data.IsPublic != nil && *data.IsPublic != board.IsPublic {
		if insufficientPermission {
			return nil, errors.New("unauthorized: only admins and owners can switch this board to public")
		}

		board.IsPublic = *data.IsPublic
	}

	if data.UserID != nil && *data.UserID != board.UserID {
		if insufficientPermission {
			return nil, errors.New("unauthorized: only admins, managers, and owners can transfer ownership of this board to another user")
		}

		board.UserID = *data.UserID
	}

	if data.WorkspaceID != nil && *data.WorkspaceID != board.WorkspaceID {
		if insufficientPermission {
			return nil, errors.New("unauthorized: only admins, managers, and owners can transfer this board to another workspace")
		}

		board.WorkspaceID = *data.WorkspaceID
	}

	if data.PermissionIDs != nil && len(*data.PermissionIDs) > 0 {
		if insufficientPermission {
			return nil, errors.New("unauthorized: only admins, managers, and owners can update permissions for this board")
		}

		err := boardService.boardRepository.UpdateWithPermissions(
			board,
			*data.PermissionIDs,
			repository.WithPreload("User", "Permissions", "Workspace", "Lists"),
		)
		if err != nil {
			return nil, err
		}

		return board, nil
	}

	err2 := boardService.boardRepository.Update(board, repository.WithPreload("User", "Permissions", "Workspace", "Lists"))
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
