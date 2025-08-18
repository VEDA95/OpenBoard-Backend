package repository

import (
	"VEDA95/open_board/api/internal/db"
	models "VEDA95/open_board/api/internal/db/model"
	"errors"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository() (*RoleRepository, error) {
	if db.Instance == nil {
		return nil, errors.New("database not initialized")
	}

	return &RoleRepository{db: db.Instance}, nil
}

func (roleRepo *RoleRepository) FindAll() ([]*models.Role, error) {
	roles := make([]*models.Role, 0)

	if err := roleRepo.db.Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (roleRepo *RoleRepository) FindByID(ID string) (*models.Role, error) {
	role := new(models.Role)

	if err := roleRepo.db.First(role, ID).Error; err != nil {
		return nil, err
	}

	return role, nil
}

func (roleRepo *RoleRepository) FindByIDs(IDs ...string) ([]*models.Role, error) {
	roles := make([]*models.Role, 0)

	if err := roleRepo.db.Where("id IN ?", IDs).Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (roleRepo *RoleRepository) FindByName(name string) (*models.Role, error) {
	role := new(models.Role)

	if err := roleRepo.db.Where("name = ?", name).First(role).Error; err != nil {
		return nil, err
	}

	return role, nil
}

func (roleRepo *RoleRepository) Create(role *models.Role) error {
	return roleRepo.db.Create(role).Error
}

func (roleRepo *RoleRepository) Update(role *models.Role) error {
	return roleRepo.db.Save(role).Error
}

func (roleRepo *RoleRepository) Delete(ID string) error {
	return roleRepo.db.Delete(&models.Role{}, ID).Error
}
