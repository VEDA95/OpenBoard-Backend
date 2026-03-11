package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

type MultiAuthMethodRepository struct {
	db *gorm.DB
}

func NewMultiAuthMethodRepository(db *gorm.DB) *MultiAuthMethodRepository {
	return &MultiAuthMethodRepository{db: db}
}

func (repo *MultiAuthMethodRepository) FindAll(opts ...QueryOption) ([]*models.MultiAuthMethod, error) {
	methods := make([]*models.MultiAuthMethod, 0)

	if err := applyOptions(repo.db, opts).Find(&methods).Error; err != nil {
		return nil, err
	}

	return methods, nil
}

func (repo *MultiAuthMethodRepository) FindByID(ID string, opts ...QueryOption) (*models.MultiAuthMethod, error) {
	method := new(models.MultiAuthMethod)

	if err := applyOptions(repo.db, opts).Where("id = ?", ID).First(method).Error; err != nil {
		return nil, err
	}

	return method, nil
}

func (repo *MultiAuthMethodRepository) FindByUserID(userID string, opts ...QueryOption) ([]*models.MultiAuthMethod, error) {
	methods := make([]*models.MultiAuthMethod, 0)

	if err := applyOptions(repo.db, opts).Where("user_id = ?", userID).Find(&methods).Error; err != nil {
		return nil, err
	}

	return methods, nil
}

func (repo *MultiAuthMethodRepository) FindByUserIDAndType(userID string, methodType string, opts ...QueryOption) (*models.MultiAuthMethod, error) {
	method := new(models.MultiAuthMethod)

	if err := applyOptions(repo.db, opts).Where("user_id = ? AND type = ?", userID, methodType).First(method).Error; err != nil {
		return nil, err
	}

	return method, nil
}

func (repo *MultiAuthMethodRepository) FindByUserIDAndName(userID string, name string, opts ...QueryOption) (*models.MultiAuthMethod, error) {
	method := new(models.MultiAuthMethod)

	if err := applyOptions(repo.db, opts).Where("user_id = ? AND name = ?", userID, name).First(method).Error; err != nil {
		return nil, err
	}

	return method, nil
}

func (repo *MultiAuthMethodRepository) Create(method *models.MultiAuthMethod, opts ...QueryOption) error {
	return applyOptions(repo.db, opts).Create(method).Error
}

func (repo *MultiAuthMethodRepository) Update(method *models.MultiAuthMethod, opts ...QueryOption) error {
	return applyOptions(repo.db, opts).Save(method).Error
}

func (repo *MultiAuthMethodRepository) Delete(ID string) error {
	return repo.db.Delete(&models.MultiAuthMethod{}, "id = ?", ID).Error
}

func (repo *MultiAuthMethodRepository) DeleteByUserID(userID string) error {
	return repo.db.Delete(&models.MultiAuthMethod{}, "user_id = ?", userID).Error
}

func (repo *MultiAuthMethodRepository) DeleteByUserIDAndType(userID string, methodType string) error {
	return repo.db.Delete(&models.MultiAuthMethod{}, "user_id = ? AND type = ?", userID, methodType).Error
}

func (repo *MultiAuthMethodRepository) Exists(ID string) bool {
	exists := false
	repo.db.Raw("SELECT EXISTS(SELECT 1 FROM multi_auth_methods WHERE id = ?) AS found", ID).Find(&exists)
	return exists
}

func (repo *MultiAuthMethodRepository) ExistsByUserIDAndType(userID string, methodType string) bool {
	exists := false
	repo.db.Raw("SELECT EXISTS(SELECT 1 FROM multi_auth_methods WHERE user_id = ? AND type = ?) AS found", userID, methodType).Find(&exists)
	return exists
}

func (repo *MultiAuthMethodRepository) ExistsByUserIDAndName(userID string, name string) bool {
	exists := false
	repo.db.Raw("SELECT EXISTS(SELECT 1 FROM multi_auth_methods WHERE user_id = ? AND name = ?) AS found", userID, name).Find(&exists)
	return exists
}

func (repo *MultiAuthMethodRepository) CountByUserID(userID string) int64 {
	var count int64
	repo.db.Model(&models.MultiAuthMethod{}).Where("user_id = ?", userID).Count(&count)
	return count
}
