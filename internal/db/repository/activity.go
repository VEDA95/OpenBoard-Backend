package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

type ActivityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (activityRepo *ActivityRepository) FindAll(options QueryOptions) ([]*models.CardActivity, error) {
	activities := make([]*models.CardActivity, 0)
	if err := options.AppendToQuery(activityRepo.db).Find(&activities).Error; err != nil {
		return nil, err
	}

	return activities, nil
}

func (activityRepo *ActivityRepository) FindByID(ID string, options QueryOptions) (*models.CardActivity, error) {
	activity := new(models.CardActivity)
	if err := options.AppendToQuery(activityRepo.db).Where("id = ?", ID).First(activity).Error; err != nil {
		return nil, err
	}

	return activity, nil
}

func (activityRepo *ActivityRepository) Create(activity *models.CardActivity, options QueryOptions) error {
	return options.AppendToQuery(activityRepo.db).Create(activity).Error
}

func (activityRepo *ActivityRepository) Update(activity *models.CardActivity, options QueryOptions) error {
	return options.AppendToQuery(activityRepo.db).Save(activity).Error
}

func (activityRepo *ActivityRepository) Delete(ID string) error {
	return activityRepo.db.Where("id = ?", ID).Delete(models.CardActivity{}).Error
}
