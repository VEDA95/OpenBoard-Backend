package service

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
)

type RoleService struct {
	repo *repository.RoleRepository
}

func (roleService *RoleService) GetRoles() ([]*models.Role, error) {
	roles, err := roleService.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (roleService *RoleService) GetRole(ID string) (*models.Role, error) {
	role, err := roleService.repo.FindByID(ID)
	if err != nil {
		return nil, err
	}

	return role, nil
}
