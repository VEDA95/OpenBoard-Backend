package repository

import (
	"VEDA95/open_board/api/internal/db"
	models "VEDA95/open_board/api/internal/db/model"
	"errors"

	"gorm.io/gorm"
)

type PasswordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository() (*PasswordResetRepository, error) {
	if db.Instance == nil {
		return nil, errors.New("database not initialized")
	}

	return &PasswordResetRepository{db: db.Instance}, nil
}

func (passwordResetRepo *PasswordResetRepository) FindAll() ([]*models.PasswordResetToken, error) {
	passwordResetTokens := make([]*models.PasswordResetToken, 0)

	if err := passwordResetRepo.db.Find(passwordResetTokens).Error; err != nil {
		return nil, err
	}

	return passwordResetTokens, nil
}

func (passwordResetRepo *PasswordResetRepository) FindByID(ID string) (*models.PasswordResetToken, error) {
	passwordResetToken := new(models.PasswordResetToken)

	if err := passwordResetRepo.db.First(passwordResetToken, ID).Error; err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}

func (passwordResetRepo *PasswordResetRepository) FindByToken(token string) (*models.PasswordResetToken, error) {
	passwordResetToken := new(models.PasswordResetToken)

	if err := passwordResetRepo.db.Where("token = ?", token).First(passwordResetToken).Error; err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}

func (passwordResetRepo *PasswordResetRepository) Create(passwordResetToken *models.PasswordResetToken) error {
	return passwordResetRepo.db.Create(passwordResetToken).Error
}

func (passwordResetRepo *PasswordResetRepository) Update(passwordResetToken *models.PasswordResetToken) error {
	return passwordResetRepo.db.Save(passwordResetToken).Error
}

func (passwordResetRepo *PasswordResetRepository) Delete(ID string) error {
	return passwordResetRepo.db.Delete(&models.PasswordResetToken{}, ID).Error
}
