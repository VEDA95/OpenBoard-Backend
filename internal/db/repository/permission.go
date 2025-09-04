package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (permissionRepo *PermissionRepository) FindAll(options QueryOptions) ([]*models.Permission, error) {
	permissions := make([]*models.Permission, 0)

	if err := options.AppendToQuery(permissionRepo.db).Find(permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

func (permissionRepo *PermissionRepository) FindByPaths(paths []string, options QueryOptions) ([]*models.Permission, error) {
	permissions := make([]*models.Permission, 0)

	if err := options.AppendToQuery(permissionRepo.db).Where("path IN ?", paths).Find(permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

func (permissionRepo *PermissionRepository) FindByID(ID string, options QueryOptions) (*models.Permission, error) {
	permission := new(models.Permission)

	if err := options.AppendToQuery(permissionRepo.db).Where("id = ?", ID).First(permission).Error; err != nil {
		return nil, err
	}

	return permission, nil
}

func (permissionRepo *PermissionRepository) Create(permission *models.Permission, options QueryOptions) error {
	return options.AppendToQuery(permissionRepo.db).Create(permission).Error
}

func (permissionRepo *PermissionRepository) Update(permission *models.Permission, options QueryOptions) error {
	return options.AppendToQuery(permissionRepo.db).Save(permission).Error
}

func (permissionRepo *PermissionRepository) Delete(ID string) error {
	return permissionRepo.db.Delete(&models.Permission{}, ID).Error
}
