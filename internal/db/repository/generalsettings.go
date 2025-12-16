package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

type GeneralSettingsRepository struct {
	db *gorm.DB
}

func NewGeneralSettingsRepository(db *gorm.DB) *GeneralSettingsRepository {
	return &GeneralSettingsRepository{
		db: db,
	}
}

func (generalSettingsRepo *GeneralSettingsRepository) Find() (*models.GeneralSettings, error) {
	settings := new(models.GeneralSettings)

	if err := generalSettingsRepo.db.First(settings).Error; err != nil {
		return nil, err
	}

	return settings, nil
}

func (generalSettingsRepo *GeneralSettingsRepository) Create() error {
	return generalSettingsRepo.db.Create(&models.GeneralSettings{}).Error
}

func (generalSettingsRepo *GeneralSettingsRepository) Update(settings *models.GeneralSettings) error {
	return generalSettingsRepo.db.Save(settings).Error
}

func (generalSettingsRepo *GeneralSettingsRepository) Delete() error {
	return generalSettingsRepo.db.Delete(models.GeneralSettings{}).Error
}
