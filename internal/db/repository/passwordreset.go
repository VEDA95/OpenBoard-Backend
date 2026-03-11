package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

type PasswordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

func (passwordResetRepo *PasswordResetRepository) FindAll(opts ...QueryOption) ([]*models.PasswordResetToken, error) {
	passwordResetTokens := make([]*models.PasswordResetToken, 0)

	if err := applyOptions(passwordResetRepo.db, opts).Find(passwordResetTokens).Error; err != nil {
		return nil, err
	}

	return passwordResetTokens, nil
}

func (passwordResetRepo *PasswordResetRepository) FindByID(ID string, opts ...QueryOption) (*models.PasswordResetToken, error) {
	passwordResetToken := new(models.PasswordResetToken)

	if err := applyOptions(passwordResetRepo.db, opts).Where("id = ?", ID).First(passwordResetToken).Error; err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}

func (passwordResetRepo *PasswordResetRepository) FindByToken(token string, opts ...QueryOption) (*models.PasswordResetToken, error) {
	passwordResetToken := new(models.PasswordResetToken)

	if err := applyOptions(passwordResetRepo.db, opts).Where("token = ?", token).First(passwordResetToken).Error; err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}

func (passwordResetRepo *PasswordResetRepository) Create(passwordResetToken *models.PasswordResetToken, opts ...QueryOption) error {
	return applyOptions(passwordResetRepo.db, opts).Create(passwordResetToken).Error
}

func (passwordResetRepo *PasswordResetRepository) Update(passwordResetToken *models.PasswordResetToken, opts ...QueryOption) error {
	return applyOptions(passwordResetRepo.db, opts).Save(passwordResetToken).Error
}

func (passwordResetRepo *PasswordResetRepository) Delete(ID string) error {
	return passwordResetRepo.db.Delete(&models.PasswordResetToken{}, ID).Error
}
