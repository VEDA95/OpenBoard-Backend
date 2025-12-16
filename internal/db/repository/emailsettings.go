package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

type EmailSettingsRepository struct {
	db *gorm.DB
}

func NewEmailSettingsRepository(db *gorm.DB) *EmailSettingsRepository {
	return &EmailSettingsRepository{
		db: db,
	}
}

func (generalSettingsRepo *EmailSettingsRepository) Find() (*models.EmailSettings, error) {
	settings := new(models.EmailSettings)

	if err := generalSettingsRepo.db.First(settings).Error; err != nil {
		return nil, err
	}

	return settings, nil
}

func (generalSettingsRepo *EmailSettingsRepository) Create() error {
	return generalSettingsRepo.db.Create(&models.EmailSettings{}).Error
}

func (generalSettingsRepo *EmailSettingsRepository) Update(settings *models.EmailSettings) error {
	return generalSettingsRepo.db.Save(settings).Error
}

func (generalSettingsRepo *EmailSettingsRepository) Delete() error {
	return generalSettingsRepo.db.Delete(models.EmailSettings{}).Error
}
