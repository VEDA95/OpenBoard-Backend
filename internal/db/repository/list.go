package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

var ListFullLoad = []QueryOption{
	WithPreload("Board", "Cards"),
}

type ListRepository struct {
	db *gorm.DB
}

func NewListRepository(db *gorm.DB) *ListRepository {
	return &ListRepository{db: db}
}

func (listRepo *ListRepository) FindAll(opts ...QueryOption) ([]*models.List, error) {
	lists := make([]*models.List, 0)
	if err := applyOptions(listRepo.db, opts).Find(&lists).Error; err != nil {
		return nil, err
	}
	return lists, nil
}

func (listRepo *ListRepository) FindByID(ID string, opts ...QueryOption) (*models.List, error) {
	list := new(models.List)
	if err := applyOptions(listRepo.db, opts).Where("id = ?", ID).First(list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (listRepo *ListRepository) FindByBoardID(ID string, opts ...QueryOption) ([]*models.List, error) {
	lists := make([]*models.List, 0)
	if err := applyOptions(listRepo.db, opts).Where("board_id = ?", ID).Order("position ASC").Find(&lists).Error; err != nil {
		return nil, err
	}
	return lists, nil
}

func (listRepo *ListRepository) Create(list *models.List, opts ...QueryOption) error {
	return applyOptions(listRepo.db, opts).Create(list).Error
}

func (listRepo *ListRepository) Update(list *models.List, opts ...QueryOption) error {
	return applyOptions(listRepo.db, opts).Save(list).Error
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
