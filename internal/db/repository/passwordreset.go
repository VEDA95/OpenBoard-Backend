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

func (passwordResetRepo *PasswordResetRepository) FindAll(options QueryOptions) ([]*models.PasswordResetToken, error) {
	passwordResetTokens := make([]*models.PasswordResetToken, 0)

	if err := options.AppendToQuery(passwordResetRepo.db).Find(passwordResetTokens).Error; err != nil {
		return nil, err
	}

	return passwordResetTokens, nil
}

func (passwordResetRepo *PasswordResetRepository) FindByID(ID string, options QueryOptions) (*models.PasswordResetToken, error) {
	passwordResetToken := new(models.PasswordResetToken)

	if err := options.AppendToQuery(passwordResetRepo.db).First(passwordResetToken, ID).Error; err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}

func (passwordResetRepo *PasswordResetRepository) FindByToken(token string, options QueryOptions) (*models.PasswordResetToken, error) {
	passwordResetToken := new(models.PasswordResetToken)

	if err := options.AppendToQuery(passwordResetRepo.db).Where("token = ?", token).First(passwordResetToken).Error; err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}

func (passwordResetRepo *PasswordResetRepository) Create(passwordResetToken *models.PasswordResetToken, options QueryOptions) error {
	return options.AppendToQuery(passwordResetRepo.db).Create(passwordResetToken).Error
}

func (passwordResetRepo *PasswordResetRepository) Update(passwordResetToken *models.PasswordResetToken, options QueryOptions) error {
	return options.AppendToQuery(passwordResetRepo.db).Save(passwordResetToken).Error
}

func (passwordResetRepo *PasswordResetRepository) Delete(ID string) error {
	return passwordResetRepo.db.Delete(&models.PasswordResetToken{}, ID).Error
}
