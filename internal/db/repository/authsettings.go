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

	rpID := os.Getenv("WEBAUTHN_RP_ID")
	if len(rpID) == 0 {
		rpID = "localhost"
	}

	rpDisplayName := os.Getenv("WEBAUTHN_RP_DISPLAY_NAME")
	if len(rpDisplayName) == 0 {
		rpDisplayName = "Open Board"
	}

	rpOrigins := os.Getenv("WEBAUTHN_RP_ORIGINS")
	if len(rpOrigins) == 0 {
		rpOrigins = frontendURL
	}

	return generalSettingsRepo.db.Create(&models.AuthSettings{
		CORSDomain:            frontendURL,
		WebAuthnRPID:          rpID,
		WebAuthnRPDisplayName: rpDisplayName,
		WebAuthnRPOrigins:     rpOrigins,
	}).Error
}

func (generalSettingsRepo *AuthSettingsRepository) Update(settings *models.AuthSettings) error {
	return generalSettingsRepo.db.Save(settings).Error
}

func (generalSettingsRepo *AuthSettingsRepository) Delete() error {
	return generalSettingsRepo.db.Delete(models.AuthSettings{}).Error
}
