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

func (repo *ExternalAuthProviderRepository) FindAll(options QueryOptions) ([]*models.ExternalAuthProvider, error) {
	providers := make([]*models.ExternalAuthProvider, 0)

	if err := options.AppendToQuery(repo.db).Find(&providers).Error; err != nil {
		return nil, err
	}

	return providers, nil
}

func (repo *ExternalAuthProviderRepository) FindByID(ID string, options QueryOptions) (*models.ExternalAuthProvider, error) {
	provider := new(models.ExternalAuthProvider)

	if err := options.AppendToQuery(repo.db).Where("id = ?", ID).First(provider).Error; err != nil {
		return nil, err
	}

	return provider, nil
}

func (repo *ExternalAuthProviderRepository) FindByName(name string, options QueryOptions) (*models.ExternalAuthProvider, error) {
	provider := new(models.ExternalAuthProvider)

	if err := options.AppendToQuery(repo.db).Where("name = ?", name).First(provider).Error; err != nil {
		return nil, err
	}

	return provider, nil
}

func (repo *ExternalAuthProviderRepository) FindDefault(options QueryOptions) (*models.ExternalAuthProvider, error) {
	provider := new(models.ExternalAuthProvider)

	if err := options.AppendToQuery(repo.db).Where("default_login_method = ?", true).First(provider).Error; err != nil {
		return nil, err
	}

	return provider, nil
}

func (repo *ExternalAuthProviderRepository) FindEnabled(options QueryOptions) ([]*models.ExternalAuthProvider, error) {
	providers := make([]*models.ExternalAuthProvider, 0)

	// Providers that have valid client_id and client_secret are considered enabled
	if err := options.AppendToQuery(repo.db).Where("client_id != '' AND client_secret != ''").Find(&providers).Error; err != nil {
		return nil, err
	}

	return providers, nil
}

func (repo *ExternalAuthProviderRepository) Create(provider *models.ExternalAuthProvider, options QueryOptions) error {
	return options.AppendToQuery(repo.db).Create(provider).Error
}

func (repo *ExternalAuthProviderRepository) CreateWithPermissions(provider *models.ExternalAuthProvider, permissionIDs []string, providerOptions QueryOptions, permissionOptions QueryOptions) error {
	return repo.db.Transaction(func(transaction *gorm.DB) error {
		permissions := make([]*models.Permission, 0)

		if len(permissionIDs) > 0 {
			if err := permissionOptions.AppendToQuery(transaction).Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
				return err
			}
		}

		if err := providerOptions.AppendToQuery(transaction).Create(provider).Error; err != nil {
			return err
		}

		if len(permissions) > 0 {
			return transaction.Model(provider).Association("Permissions").Append(permissions)
		}

		return nil
	})
}

func (repo *ExternalAuthProviderRepository) Update(provider *models.ExternalAuthProvider, options QueryOptions) error {
	return options.AppendToQuery(repo.db).Save(provider).Error
}

func (repo *ExternalAuthProviderRepository) UpdateWithPermissions(provider *models.ExternalAuthProvider, permissionIDs []string, providerOptions QueryOptions, permissionOptions QueryOptions) error {
	return repo.db.Transaction(func(transaction *gorm.DB) error {
		if err := providerOptions.AppendToQuery(transaction).Save(provider).Error; err != nil {
			return err
		}

		permissions := make([]*models.Permission, 0)

		if len(permissionIDs) > 0 {
			if err := permissionOptions.AppendToQuery(transaction).Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
				return err
			}
		}

		return transaction.Model(provider).Association("Permissions").Replace(permissions)
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
