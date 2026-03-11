package service

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/http/validators"
)

type RoleService struct {
	repo *repository.RoleRepository
}

func NewRoleService(roleRepo *repository.RoleRepository) *RoleService {
	return &RoleService{repo: roleRepo}
}

func (roleService *RoleService) GetRoles() ([]*models.Role, error) {
	roles, err := roleService.repo.FindAll(repository.RoleFullLoad...)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (roleService *RoleService) GetRole(ID string) (*models.Role, error) {
	role, err := roleService.repo.FindByID(ID, repository.RoleFullLoad...)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (roleService *RoleService) CreateRole(data *validators.CreateRoleValidator) (*models.Role, error) {
	role := &models.Role{Name: data.Name}

	if data.Permissions != nil && len(*data.Permissions) > 0 {
		err := roleService.repo.CreateWithPermissions(
			role,
			*data.Permissions,
			repository.RoleFullLoad...,
		)
		if err != nil {
			return nil, err
		}

		return role, nil
	}

	err := roleService.repo.Create(role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (roleService *RoleService) UpdateRole(ID string, data *validators.UpdateRoleValidator) (*models.Role, error) {
	role, err := roleService.repo.FindByID(ID, repository.WithOmit("Permissions"))
	if err != nil {
		return nil, err
	}

	if data.Name != nil && *data.Name != role.Name {
		role.Name = *data.Name
	}

	if data.Permissions != nil {
		err := roleService.repo.UpdateWithPermissions(
			role,
			*data.Permissions,
			repository.RoleFullLoad...,
		)
		if err != nil {
			return nil, err
		}

		return role, nil
	}

	err2 := roleService.repo.Update(role, repository.WithOmit("Permissions"))
	if err2 != nil {
		return nil, err2
	}

	return role, nil
}

func (roleService *RoleService) DeleteRole(ID string) error {
	return roleService.repo.Delete(ID)
}
