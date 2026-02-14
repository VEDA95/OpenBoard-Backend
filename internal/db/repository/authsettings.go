package repository

import (
	models "VEDA95/open_board/api/internal/db/model"
	"os"

	"gorm.io/gorm"
)

type AuthSettingsRepository struct {
	db *gorm.DB
}

func NewAuthSettingsRepository(db *gorm.DB) *AuthSettingsRepository {
	return &AuthSettingsRepository{
		db: db,
	}
}

func (generalSettingsRepo *AuthSettingsRepository) Find() (*models.AuthSettings, error) {
	settings := new(models.AuthSettings)

	if err := generalSettingsRepo.db.First(settings).Error; err != nil {
		return nil, err
	}

	return settings, nil
}

func (generalSettingsRepo *AuthSettingsRepository) Create() error {
	frontendURL := os.Getenv("FRONTEND_URL")

	if len(frontendURL) == 0 {
		frontendURL = "http://localhost:3000"
	}

	return generalSettingsRepo.db.Create(&models.AuthSettings{CORSDomain: frontendURL}).Error
}

func (generalSettingsRepo *AuthSettingsRepository) Update(settings *models.AuthSettings) error {
	return generalSettingsRepo.db.Save(settings).Error
}

func (generalSettingsRepo *AuthSettingsRepository) Delete() error {
	return generalSettingsRepo.db.Delete(models.AuthSettings{}).Error
}
