package service

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/http/validators"
)

type RoleService struct {
	repo *repository.RoleRepository
}

func NewRoleService() (*RoleService, error) {
	roleRepo, err := repository.NewRoleRepository()
	if err != nil {
		return nil, err
	}

	return &RoleService{repo: roleRepo}, nil
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

func (roleService *RoleService) CreateRole(data *validators.CreateRoleValidator) (*models.Role, error) {
	role := &models.Role{Name: data.Name}

	if data.Permissions != nil && len(*data.Permissions) > 0 {
		if err := roleService.repo.CreateWithPermissions(role, *data.Permissions...); err != nil {
			return nil, err
		}

		return role, nil
	}

	if err := roleService.repo.Create(role); err != nil {
		return nil, err
	}

	return role, nil
}

func (roleService *RoleService) UpdateRole(ID string, data *validators.UpdateRoleValidator) (*models.Role, error) {
	role, err := roleService.repo.FindByID(ID)
	if err != nil {
		return nil, err
	}

	if data.Name != nil && *data.Name != role.Name {
		role.Name = *data.Name
	}

	if data.Permissions != nil {
		if err := roleService.repo.UpdateWithPermissions(role, *data.Permissions...); err != nil {
			return nil, err
		}

		return role, nil
	}

	if err := roleService.repo.Update(role); err != nil {
		return nil, err
	}

	return role, nil
}

func (roleService *RoleService) DeleteRole(ID string) error {
	return roleService.repo.Delete(ID)
}
