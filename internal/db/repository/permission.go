package repository

import (
	"VEDA95/open_board/api/internal/db"
	models "VEDA95/open_board/api/internal/db/model"
	"errors"

	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository() (*PermissionRepository, error) {
	if db.Instance == nil {
		return nil, errors.New("database not initialized")
	}

	return &PermissionRepository{db: db.Instance}, nil
}

func (permissionRepo *PermissionRepository) FindAll() ([]*models.Permission, error) {
	permissions := make([]*models.Permission, 0)

	if err := permissionRepo.db.Find(permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

func (permissionRepo *PermissionRepository) FindByPaths(paths ...string) ([]*models.Permission, error) {
	permissions := make([]*models.Permission, 0)

	if err := permissionRepo.db.Where("path IN ?", paths).Find(permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

func (permissionRepo *PermissionRepository) FindByID(ID string) (*models.Permission, error) {
	permission := new(models.Permission)

	if err := permissionRepo.db.First(permission, ID).Error; err != nil {
		return nil, err
	}

	return permission, nil
}

func (permissionRepo *PermissionRepository) Create(permission *models.Permission) error {
	return permissionRepo.db.Create(permission).Error
}

func (permissionRepo *PermissionRepository) Update(permission *models.Permission) error {
	return permissionRepo.db.Save(permission).Error
}

func (permissionRepo *PermissionRepository) Delete(ID string) error {
	return permissionRepo.db.Delete(&models.Permission{}, ID).Error
}
