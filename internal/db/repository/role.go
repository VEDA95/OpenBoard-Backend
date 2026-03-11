package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

var RoleFullLoad = []QueryOption{
	WithPreload("Permissions"),
}

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (roleRepo *RoleRepository) FindAll(opts ...QueryOption) ([]*models.Role, error) {
	roles := make([]*models.Role, 0)

	if err := applyOptions(roleRepo.db, opts).Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (roleRepo *RoleRepository) FindByID(ID string, opts ...QueryOption) (*models.Role, error) {
	role := new(models.Role)

	if err := applyOptions(roleRepo.db, opts).Where("id = ?", ID).First(role).Error; err != nil {
		return nil, err
	}

	return role, nil
}

func (roleRepo *RoleRepository) FindByIDs(IDs []string, opts ...QueryOption) ([]*models.Role, error) {
	roles := make([]*models.Role, 0)

	if err := applyOptions(roleRepo.db, opts).Where("id IN ?", IDs).Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (roleRepo *RoleRepository) FindByName(name string, opts ...QueryOption) (*models.Role, error) {
	role := new(models.Role)

	if err := applyOptions(roleRepo.db, opts).Where("name = ?", name).First(role).Error; err != nil {
		return nil, err
	}

	return role, nil
}

func (roleRepo *RoleRepository) Create(role *models.Role, opts ...QueryOption) error {
	return applyOptions(roleRepo.db, opts).Create(role).Error
}

func (roleRepo *RoleRepository) CreateWithPermissions(role *models.Role, permissionIDs []string, reloadOpts ...QueryOption) error {
	return roleRepo.db.Transaction(func(transaction *gorm.DB) error {
		permissions := make([]*models.Permission, 0)

		if err := transaction.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			return err
		}

		if err := transaction.Create(role).Error; err != nil {
			return err
		}

		if err := transaction.Model(role).Association("Permissions").Append(permissions); err != nil {
			return err
		}

		return applyOptions(transaction, reloadOpts).First(role, "id = ?", role.ID).Error
	})
}

func (roleRepo *RoleRepository) Update(role *models.Role, opts ...QueryOption) error {
	return applyOptions(roleRepo.db, opts).Save(role).Error
}

func (roleRepo *RoleRepository) UpdateWithPermissions(role *models.Role, permissionsIDs []string, reloadOpts ...QueryOption) error {
	return roleRepo.db.Transaction(func(transaction *gorm.DB) error {
		if err := transaction.Save(role).Error; err != nil {
			return err
		}

		if len(permissionsIDs) == 0 {
			return nil
		}

		permissions := make([]*models.Permission, 0)
		if err := transaction.Where("id IN ?", permissionsIDs).Find(&permissions).Error; err != nil {
			return err
		}

		if err := transaction.Model(role).Association("Permissions").Replace(permissions); err != nil {
			return err
		}

		return applyOptions(transaction, reloadOpts).First(role, "id = ?", role.ID).Error
	})
}

func (roleRepo *RoleRepository) Delete(ID string) error {
	return roleRepo.db.Delete(&models.Role{}, ID).Error
}

func (roleRepo *RoleRepository) Exists(ID string) bool {
	exists := false

	roleRepo.db.Raw("SELECT EXISTS(SELECT 1 FROM roles WHERE id = ?) AS found", ID).Find(&exists)

	return exists
}

func (roleRepo *RoleRepository) ExistsByName(name string) bool {
	exists := false

	roleRepo.db.Raw("SELECT EXISTS(SELECT 1 FROM roles WHERE name = ?) AS found", name).Find(&exists)

	return exists
}
