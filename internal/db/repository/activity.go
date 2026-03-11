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

func (activityRepo *ActivityRepository) FindAll(opts ...QueryOption) ([]*models.CardActivity, error) {
	activities := make([]*models.CardActivity, 0)
	if err := applyOptions(activityRepo.db, opts).Find(&activities).Error; err != nil {
		return nil, err
	}

	return activities, nil
}

func (activityRepo *ActivityRepository) FindByID(ID string, opts ...QueryOption) (*models.CardActivity, error) {
	activity := new(models.CardActivity)
	if err := applyOptions(activityRepo.db, opts).Where("id = ?", ID).First(activity).Error; err != nil {
		return nil, err
	}

	return activity, nil
}

func (activityRepo *ActivityRepository) Create(activity *models.CardActivity, opts ...QueryOption) error {
	return applyOptions(activityRepo.db, opts).Create(activity).Error
}

func (activityRepo *ActivityRepository) Update(activity *models.CardActivity, opts ...QueryOption) error {
	return applyOptions(activityRepo.db, opts).Save(activity).Error
}

func (activityRepo *ActivityRepository) Delete(ID string) error {
	return activityRepo.db.Where("id = ?", ID).Delete(models.CardActivity{}).Error
}
