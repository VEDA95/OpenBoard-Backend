package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

type ListRepository struct {
	db *gorm.DB
}

func NewListRepository(db *gorm.DB) *ListRepository {
	return &ListRepository{db: db}
}

func (listRepo *ListRepository) FindAll(options QueryOptions) ([]*models.List, error) {
	lists := make([]*models.List, 0)
	if err := options.AppendToQuery(listRepo.db).Find(&lists).Error; err != nil {
		return nil, err
	}
	return lists, nil
}

func (listRepo *ListRepository) FindByID(ID string, options QueryOptions) (*models.List, error) {
	list := new(models.List)
	if err := options.AppendToQuery(listRepo.db).Where("id = ?", ID).First(list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (listRepo *ListRepository) FindByBoardID(ID string, options QueryOptions) ([]*models.List, error) {
	lists := make([]*models.List, 0)
	if err := options.AppendToQuery(listRepo.db).Where("board_id = ?", ID).Order("position ASC").Find(&lists).Error; err != nil {
		return nil, err
	}
	return lists, nil
}

func (listRepo *ListRepository) Create(list *models.List, options QueryOptions) error {
	return options.AppendToQuery(listRepo.db).Create(list).Error
}

func (listRepo *ListRepository) Update(list *models.List, options QueryOptions) error {
	return options.AppendToQuery(listRepo.db).Save(list).Error
}

func (listRepo *ListRepository) Delete(ID string) error {
	return listRepo.db.Where("id = ?", ID).Delete(models.List{}).Error
}

func (listRepo *ListRepository) UpdatePositions(boardID string, positions map[string]int) error {
	return listRepo.db.Transaction(func(tx *gorm.DB) error {
		for listID, position := range positions {
			if err := tx.Model(&models.List{}).Where("id = ? AND board_id = ?", listID, boardID).Update("position", position).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (listRepo *ListRepository) Exists(ID string) bool {
	exists := false

	listRepo.db.Raw("SELECT EXISTS(SELECT 1 FROM lists WHERE id = ?) AS found", ID).Find(&exists)

	return exists
}
