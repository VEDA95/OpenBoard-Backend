package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (roleRepo *RoleRepository) FindAll(options QueryOptions) ([]*models.Role, error) {
	roles := make([]*models.Role, 0)

	if err := options.AppendToQuery(roleRepo.db).Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (roleRepo *RoleRepository) FindByID(ID string, options QueryOptions) (*models.Role, error) {
	role := new(models.Role)

	if err := options.AppendToQuery(roleRepo.db).Where("id = ?", ID).First(role).Error; err != nil {
		return nil, err
	}

	return role, nil
}

func (roleRepo *RoleRepository) FindByIDs(IDs []string, options QueryOptions) ([]*models.Role, error) {
	roles := make([]*models.Role, 0)

	if err := options.AppendToQuery(roleRepo.db).Where("id IN ?", IDs).Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (roleRepo *RoleRepository) FindByName(name string, options QueryOptions) (*models.Role, error) {
	role := new(models.Role)

	if err := options.AppendToQuery(roleRepo.db).Where("name = ?", name).First(role).Error; err != nil {
		return nil, err
	}

	return role, nil
}

func (roleRepo *RoleRepository) Create(role *models.Role, options QueryOptions) error {
	return options.AppendToQuery(roleRepo.db).Create(role).Error
}

func (roleRepo *RoleRepository) CreateWithPermissions(role *models.Role, permissionIDs []string, roleOptions QueryOptions, permissionOptions QueryOptions) error {
	return roleRepo.db.Transaction(func(transaction *gorm.DB) error {
		permissions := make([]*models.Permission, 0)

		if err := permissionOptions.AppendToQuery(transaction).Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			return err
		}

		if err := roleOptions.AppendToQuery(transaction).Create(role).Error; err != nil {
			return err
		}

		return transaction.Model(role).Association("Permissions").Append(permissions)
	})
}

func (roleRepo *RoleRepository) Update(role *models.Role, options QueryOptions) error {
	return options.AppendToQuery(roleRepo.db).Save(role).Error
}

func (roleRepo *RoleRepository) UpdateWithPermissions(role *models.Role, permissionsIDs []string, roleOptions QueryOptions, permissionOptions QueryOptions) error {
	return roleRepo.db.Transaction(func(transaction *gorm.DB) error {
		if err := roleOptions.AppendToQuery(transaction).Save(role); err != nil {
			return nil
		}

		if len(permissionsIDs) == 0 {
			return nil
		}

		permissions := make([]*models.Permission, 0)
		if err := permissionOptions.AppendToQuery(transaction).Where("id IN ?", permissionsIDs).Find(&permissions).Error; err != nil {
			return err
		}

		return transaction.Model(role).Association("Permissions").Replace(permissions)
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
