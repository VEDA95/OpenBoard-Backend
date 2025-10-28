package service

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
)

type WorkspaceService struct {
	workspaceRepository *repository.WorkspaceRepository
}

func NewWorkspaceService(workspaceRepository *repository.WorkspaceRepository) *WorkspaceService {
	return &WorkspaceService{
		workspaceRepository: workspaceRepository,
	}
}

func (workspaceService *WorkspaceService) getWorkspaces() ([]*models.Worksapce, error) {
	workspaces, err := workspaceService.workspaceRepository.FindAll(repository.QueryOptions{
		Preload: []string{"User", "Permissions", "Board"},
		Omit:    []string{"UserID"},
	})
	if err != nil {
		return nil, err
	}

	return workspaces, err
}
