package service

import (
	"errors"

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
		Preload: []string{"User", "Permissions", "Boards"},
		Omit:    []string{"UserID"},
	})
	if err != nil {
		return nil, err
	}

	return workspaces, err
}

func (workspaceService *WorkspaceService) GetWorkspaceByID(ID string) (*models.Worksapce, error) {
	workspace, err := workspaceService.workspaceRepository.FindByID(ID, repository.QueryOptions{
		Preload: []string{"User", "Permissions", "Boards"},
		Omit:    []string{"UserID"},
	})
	if err != nil {
		return nil, err
	}

	return workspace, nil
}

func (workspaceService *WorkspaceService) GetUserAccessibleWorkspaceByID(ID string, user *models.User) (*models.Worksapce, error) {
	workspace, err := workspaceService.workspaceRepository.FindByID(ID, repository.QueryOptions{
		Preload: []string{"User", "Permissions", "Boards"},
		Omit:    []string{"UserID"},
	})
	if err != nil {
		return nil, err
	}

	if !user.IsSuperuser() && (!workspace.IsPublic || workspace.User.ID != user.ID) {
		permissionPaths := make([]string, len(workspace.Permissions))

		for index, permission := range workspace.Permissions {
			permissionPaths[index] = permission.Path
		}

		if !user.IsAuthorizedPartial(permissionPaths...) {
			return nil, errors.New("unauthorized")
		}
	}

	return workspace, nil
}

func (workspaceService *WorkspaceService) GetWorkspaceByUserID(ID string) ([]*models.Worksapce, error) {
	workspaces, err := workspaceService.workspaceRepository.FindByUserID(ID, repository.QueryOptions{
		Preload: []string{"User", "Permissions", "Boards"},
		Omit:    []string{"UserID"},
	})
	if err != nil {
		return nil, err
	}

	return workspaces, nil
}

func (workspaceService *WorkspaceService) GetUserAccessibleWorkspaces(user *models.User) ([]*models.Worksapce, error) {
	allWorkspaces, err := workspaceService.GetWorkspaces()
	if err != nil {
		return nil, err
	}

	if user.IsSuperuser() {
		return allWorkspaces, nil
	}

	workspaces := make([]*models.Worksapce, 0)

	for _, workspace := range allWorkspaces {
		if workspace.IsPublic || workspace.User.ID == user.ID {
			workspaces = append(workspaces, workspace)
			continue
		}

		permissionPaths := make([]string, len(workspace.Permissions))

		for index, permission := range workspace.Permissions {
			permissionPaths[index] = permission.Path
		}

		if user.IsAuthorizedPartial(permissionPaths...) {
			workspaces = append(workspaces, workspace)
			continue
		}
	}

	return workspaces, nil
}

func (workspaceService *WorkspaceService) CreateWorkspace(data *validators.CreateWorkspaceValidator) (*models.Worksapce, error) {
	if workspaceService.workspaceRepository.ExistsForUserByName(data.UserID, data.Name) {
		return nil, errors.New("workspace already exists")
	}

	workspace := &models.Worksapce{
		UserID:      data.UserID,
		Name:        data.Name,
		Description: data.Description,
	}

	if data.PermissionIDs != nil && len(*data.PermissionIDs) > 0 {
		err := workspaceService.workspaceRepository.CreateWithPermissions(
			workspace,
			*data.PermissionIDs,
			repository.QueryOptions{
				Preload: []string{"User", "Permissions", "Boards"},
				Omit:    []string{"UserID"},
			},
			repository.QueryOptions{},
		)
		if err != nil {
			return nil, err
		}

		return workspace, nil
	}

	err := workspaceService.workspaceRepository.Create(
		workspace,
		repository.QueryOptions{
			Preload: []string{"User", "Permissions", "Boards"},
			Omit:    []string{"UserID"},
		},
	)
	if err != nil {
		return nil, err
	}

	return workspace, nil
}

func (workspaceService *WorkspaceService) UpdateWorkspace(ID string, data *validators.UpdateWorkspaceValidator) (*models.Worksapce, error) {
	if !workspaceService.workspaceRepository.Exists(ID) {
		return nil, errors.New("workspace does not exist")
	}

	workspace, err := workspaceService.workspaceRepository.FindByID(ID, repository.QueryOptions{})
	if err != nil {
		return nil, err
	}

	if data.Name != nil && *data.Name != workspace.Name {
		workspace.Name = *data.Name
	}

	if data.Description != nil && *data.Description != *workspace.Description {
		workspace.Name = *data.Description
	}

	if data.UserID != nil && *data.UserID != workspace.UserID {
		workspace.UserID = *data.UserID
	}

	if data.PermissionIDs != nil {
		err := workspaceService.workspaceRepository.UpdateWithPermissions(
			workspace,
			*data.PermissionIDs,
			repository.QueryOptions{
				Preload: []string{"User", "Permissions", "Boards"},
				Omit:    []string{"UserID"},
			},
			repository.QueryOptions{},
		)
		if err != nil {
			return nil, err
		}

		return workspace, nil
	}

	err2 := workspaceService.workspaceRepository.Update(
		workspace,
		repository.QueryOptions{
			Preload: []string{"User", "Permissions", "Boards"},
			Omit:    []string{"UserID"},
		},
	)

	if err2 != nil {
		return nil, err2
	}

	return workspace, nil
}

func (workspaceService *WorkspaceService) DeleteWorkspace(ID string) error {
	if !workspaceService.workspaceRepository.Exists(ID) {
		return errors.New("workspace does not exist")
	}

	return workspaceService.workspaceRepository.Delete(ID)
}
