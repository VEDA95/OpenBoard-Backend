package service

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/http/validators"
)

type WorkspaceService struct {
	workspaceRepository *repository.WorkspaceRepository
}

func NewWorkspaceService(workspaceRepository *repository.WorkspaceRepository) *WorkspaceService {
	return &WorkspaceService{
		workspaceRepository: workspaceRepository,
	}
}

func (workspaceService *WorkspaceService) GetWorkspaces() ([]*models.Worksapce, error) {
	workspaces, err := workspaceService.workspaceRepository.FindAll(repository.QueryOptions{
		Preload: []string{"User", "Permissions", "Board"},
		Omit:    []string{"UserID"},
	})
	if err != nil {
		return nil, err
	}

	return workspaces, err
}

func (workspaceService *WorkspaceService) GetWorkspaceByID(id string) (*models.Worksapce, error) {
	workspace, err := workspaceService.workspaceRepository.FindByID(id, repository.QueryOptions{
		Preload: []string{"User", "Permissions", "Board"},
		Omit:    []string{"UserID"},
	})
	if err != nil {
		return nil, err
	}

	return workspace, nil
}

func (workspaceService *WorkspaceService) CreateWorkspace(data *validators.CreateBoardValidator) (*models.Worksapce, error) {
}
