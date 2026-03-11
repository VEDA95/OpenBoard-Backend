package service

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/http/validators"
)

type PermissionService struct {
	repo *repository.PermissionRepository
}

func NewPermissionService(permissionRepo *repository.PermissionRepository) *PermissionService {
	return &PermissionService{repo: permissionRepo}
}

func (permissionService *PermissionService) GetPermissions() ([]*models.Permission, error) {
	permissions, err := permissionService.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (permissionService *PermissionService) GetPermission(ID string) (*models.Permission, error) {
	permission, err := permissionService.repo.FindByID(ID)
	if err != nil {
		return nil, err
	}

	return permission, nil
}

func (permissionService *PermissionService) GetPermissionsByPaths(paths ...string) ([]*models.Permission, error) {
	permissions, err := permissionService.repo.FindByPaths(paths)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (permissionService *PermissionService) CreatePermission(data *validators.CreatePermissionValidator) (*models.Permission, error) {
	permission := &models.Permission{Path: data.Path}
	if err := permissionService.repo.Create(permission); err != nil {
		return nil, err
	}

	return permission, nil
}

func (permissionService *PermissionService) UpdatePermission(ID string, data *validators.UpdatePermissionValidator) (*models.Permission, error) {
	permission, err := permissionService.repo.FindByID(ID,
		repository.WithSelect("id", "path"),
	)
	if err != nil {
		return nil, err
	}

	if data.Path != nil && *data.Path != permission.Path {
		permission.Path = *data.Path
	}

	if err := permissionService.repo.Update(permission); err != nil {
		return nil, err
	}

	return permission, nil
}

func (permissionService *PermissionService) DeletePermission(ID string) error {
	return permissionService.repo.Delete(ID)
}
