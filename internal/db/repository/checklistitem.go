package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

type CheckListItemRepository struct {
	db *gorm.DB
}

func NewCheckListItemRepository(db *gorm.DB) *CheckListItemRepository {
	return &CheckListItemRepository{db: db}
}

func (checkListItemRepo *CheckListItemRepository) FindAll(options QueryOptions) ([]*models.CheckListItem, error) {
	checkListItems := make([]*models.CheckListItem, 0)
	if err := options.AppendToQuery(checkListItemRepo.db).Find(&checkListItems).Error; err != nil {
		return nil, err
	}

	return checkListItems, nil
}

func (checkListItemRepo *CheckListItemRepository) FindByID(ID string, options QueryOptions) (*models.CheckListItem, error) {
	checkListItem := new(models.CheckListItem)
	if err := options.AppendToQuery(checkListItemRepo.db).Where("id = ?", ID).First(checkListItem).Error; err != nil {
		return nil, err
	}

	return checkListItem, nil
}

func (checkListItemRepo *CheckListItemRepository) Create(checkListItem *models.CheckListItem, options QueryOptions) error {
	return options.AppendToQuery(checkListItemRepo.db).Create(checkListItem).Error
}

func (checkListItemRepo *CheckListItemRepository) Update(checkListItem *models.CheckListItem, options QueryOptions) error {
	return options.AppendToQuery(checkListItemRepo.db).Save(checkListItem).Error
}

func (checkListItemRepo *CheckListItemRepository) Delete(ID string) error {
	return checkListItemRepo.db.Where("id = ?", ID).Delete(models.CheckListItem{}).Error
}
