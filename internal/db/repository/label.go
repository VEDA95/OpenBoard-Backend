package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

type LabelRepository struct {
	db *gorm.DB
}

func NewLabelRepository(db *gorm.DB) *LabelRepository {
	return &LabelRepository{db: db}
}

func (labelRepo *LabelRepository) FindAll(opts ...QueryOption) ([]*models.Label, error) {
	labels := make([]*models.Label, 0)
	if err := applyOptions(labelRepo.db, opts).Find(&labels).Error; err != nil {
		return nil, err
	}

	return labels, nil
}

func (labelRepo *LabelRepository) FindBoardID(ID string, opts ...QueryOption) ([]*models.Label, error) {
	labels := make([]*models.Label, 0)
	if err := applyOptions(labelRepo.db, opts).Where("board_id = ?", ID).Find(&labels).Error; err != nil {
		return nil, err
	}

	return labels, nil
}

func (labelRepo *LabelRepository) FindByID(ID string, opts ...QueryOption) (*models.Label, error) {
	label := new(models.Label)
	if err := applyOptions(labelRepo.db, opts).Where("id = ?", ID).First(label).Error; err != nil {
		return nil, err
	}

	return label, nil
}

func (labelRepo *LabelRepository) Create(label *models.Label, opts ...QueryOption) error {
	return applyOptions(labelRepo.db, opts).Create(label).Error
}

func (labelRepo *LabelRepository) Update(label *models.Label, opts ...QueryOption) error {
	return applyOptions(labelRepo.db, opts).Save(label).Error
}

func (labelRepo *LabelRepository) Delete(ID string) error {
	return labelRepo.db.Where("id = ?", ID).Delete(models.Label{}).Error
}

func (labelRepo *LabelRepository) Exists(ID string) bool {
	exists := false

	labelRepo.db.Raw("SELECT EXISTS(SELECT 1 FROM lists WHERE id = ?) AS found", ID).Find(&exists)

	return exists
}
