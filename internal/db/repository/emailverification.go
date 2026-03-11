package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

type EmailVerificationRepository struct {
	db *gorm.DB
}

func NewEmailVerificationRepository(db *gorm.DB) *EmailVerificationRepository {
	return &EmailVerificationRepository{db: db}
}

func (repo *EmailVerificationRepository) FindByID(ID string, opts ...QueryOption) (*models.EmailVerificationToken, error) {
	token := new(models.EmailVerificationToken)

	if err := applyOptions(repo.db, opts).Where("id = ?", ID).First(token).Error; err != nil {
		return nil, err
	}

	return token, nil
}

func (repo *EmailVerificationRepository) FindByUserID(userID string, opts ...QueryOption) (*models.EmailVerificationToken, error) {
	token := new(models.EmailVerificationToken)

	if err := applyOptions(repo.db, opts).Where("user_id = ?", userID).First(token).Error; err != nil {
		return nil, err
	}

	return token, nil
}

func (repo *EmailVerificationRepository) FindAllByUserID(userID string, opts ...QueryOption) ([]*models.EmailVerificationToken, error) {
	tokens := make([]*models.EmailVerificationToken, 0)

	if err := applyOptions(repo.db, opts).Where("user_id = ?", userID).Find(&tokens).Error; err != nil {
		return nil, err
	}

	return tokens, nil
}

func (repo *EmailVerificationRepository) Create(token *models.EmailVerificationToken, opts ...QueryOption) error {
	return applyOptions(repo.db, opts).Create(token).Error
}

func (repo *EmailVerificationRepository) Delete(ID string) error {
	return repo.db.Delete(&models.EmailVerificationToken{}, "id = ?", ID).Error
}

func (repo *EmailVerificationRepository) DeleteByUserID(userID string) error {
	return repo.db.Delete(&models.EmailVerificationToken{}, "user_id = ?", userID).Error
}

func (repo *EmailVerificationRepository) DeleteExpired() error {
	return repo.db.Delete(&models.EmailVerificationToken{}, "expires_on < NOW()").Error
}

func (repo *EmailVerificationRepository) Exists(ID string) bool {
	exists := false
	repo.db.Raw("SELECT EXISTS(SELECT 1 FROM email_verification_tokens WHERE id = ?) AS found", ID).Find(&exists)
	return exists
}

func (repo *EmailVerificationRepository) ExistsByUserID(userID string) bool {
	exists := false
	repo.db.Raw("SELECT EXISTS(SELECT 1 FROM email_verification_tokens WHERE user_id = ?) AS found", userID).Find(&exists)
	return exists
}
