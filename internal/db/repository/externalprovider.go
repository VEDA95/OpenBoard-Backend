package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

type ExternalAuthProviderRepository struct {
	db *gorm.DB
}

func NewExternalAuthProviderRepository(db *gorm.DB) *ExternalAuthProviderRepository {
	return &ExternalAuthProviderRepository{db: db}
}

func (repo *ExternalAuthProviderRepository) FindAll(opts ...QueryOption) ([]*models.ExternalAuthProvider, error) {
	providers := make([]*models.ExternalAuthProvider, 0)

	if err := applyOptions(repo.db, opts).Find(&providers).Error; err != nil {
		return nil, err
	}

	return providers, nil
}

func (repo *ExternalAuthProviderRepository) FindByID(ID string, opts ...QueryOption) (*models.ExternalAuthProvider, error) {
	provider := new(models.ExternalAuthProvider)

	if err := applyOptions(repo.db, opts).Where("id = ?", ID).First(provider).Error; err != nil {
		return nil, err
	}

	return provider, nil
}

func (repo *ExternalAuthProviderRepository) FindByName(name string, opts ...QueryOption) (*models.ExternalAuthProvider, error) {
	provider := new(models.ExternalAuthProvider)

	if err := applyOptions(repo.db, opts).Where("name = ?", name).First(provider).Error; err != nil {
		return nil, err
	}

	return provider, nil
}

func (repo *ExternalAuthProviderRepository) FindDefault(opts ...QueryOption) (*models.ExternalAuthProvider, error) {
	provider := new(models.ExternalAuthProvider)

	if err := applyOptions(repo.db, opts).Where("default_login_method = ?", true).First(provider).Error; err != nil {
		return nil, err
	}

	return provider, nil
}

func (repo *ExternalAuthProviderRepository) FindEnabled(opts ...QueryOption) ([]*models.ExternalAuthProvider, error) {
	providers := make([]*models.ExternalAuthProvider, 0)

	// Providers that have valid client_id and client_secret are considered enabled
	if err := applyOptions(repo.db, opts).Where("client_id != '' AND client_secret != ''").Find(&providers).Error; err != nil {
		return nil, err
	}

	return providers, nil
}

func (repo *ExternalAuthProviderRepository) Create(provider *models.ExternalAuthProvider, opts ...QueryOption) error {
	return applyOptions(repo.db, opts).Create(provider).Error
}

func (repo *ExternalAuthProviderRepository) CreateWithPermissions(provider *models.ExternalAuthProvider, permissionIDs []string, reloadOpts ...QueryOption) error {
	return repo.db.Transaction(func(transaction *gorm.DB) error {
		permissions := make([]*models.Permission, 0)

		if len(permissionIDs) > 0 {
			if err := transaction.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
				return err
			}
		}

		if err := transaction.Create(provider).Error; err != nil {
			return err
		}

		if len(permissions) > 0 {
			if err := transaction.Model(provider).Association("Permissions").Append(permissions); err != nil {
				return err
			}
		}

		return applyOptions(transaction, reloadOpts).First(provider, "id = ?", provider.ID).Error
	})
}

func (repo *ExternalAuthProviderRepository) Update(provider *models.ExternalAuthProvider, opts ...QueryOption) error {
	return applyOptions(repo.db, opts).Save(provider).Error
}

func (repo *ExternalAuthProviderRepository) UpdateWithPermissions(provider *models.ExternalAuthProvider, permissionIDs []string, reloadOpts ...QueryOption) error {
	return repo.db.Transaction(func(transaction *gorm.DB) error {
		if err := transaction.Save(provider).Error; err != nil {
			return err
		}

		permissions := make([]*models.Permission, 0)

		if len(permissionIDs) > 0 {
			if err := transaction.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
				return err
			}
		}

		if err := transaction.Model(provider).Association("Permissions").Replace(permissions); err != nil {
			return err
		}

		return applyOptions(transaction, reloadOpts).First(provider, "id = ?", provider.ID).Error
	})
}

func (repo *ExternalAuthProviderRepository) Delete(ID string) error {
	return repo.db.Delete(&models.ExternalAuthProvider{}, "id = ?", ID).Error
}

func (repo *ExternalAuthProviderRepository) Exists(ID string) bool {
	exists := false
	repo.db.Raw("SELECT EXISTS(SELECT 1 FROM external_auth_providers WHERE id = ?) AS found", ID).Find(&exists)
	return exists
}

func (repo *ExternalAuthProviderRepository) ExistsByName(name string) bool {
	exists := false
	repo.db.Raw("SELECT EXISTS(SELECT 1 FROM external_auth_providers WHERE name = ?) AS found", name).Find(&exists)
	return exists
}
