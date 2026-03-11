package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

var CardFullLoad = []QueryOption{
	WithPreload("List", "Labels", "Comments", "Comments.User"),
}

var CardDetailLoad = []QueryOption{
	WithPreload("List", "Labels", "Comments", "Comments.User", "Activities", "Attachments"),
}

type CardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) *CardRepository {
	return &CardRepository{db: db}
}

func (cardRepo *CardRepository) FindAll(opts ...QueryOption) ([]*models.Card, error) {
	cards := make([]*models.Card, 0)
	if err := applyOptions(cardRepo.db, opts).Find(&cards).Error; err != nil {
		return nil, err
	}

	return cards, nil
}

func (cardRepo *CardRepository) FindByID(ID string, opts ...QueryOption) (*models.Card, error) {
	card := new(models.Card)
	if err := applyOptions(cardRepo.db, opts).Where("id = ?", ID).First(card).Error; err != nil {
		return nil, err
	}

	return card, nil
}

func (cardRepo *CardRepository) Create(card *models.Card, opts ...QueryOption) error {
	return applyOptions(cardRepo.db, opts).Create(card).Error
}

func (cardRepo *CardRepository) CreateWithAssociations(card *models.Card, labelIDs []string, attachmentIDs []string, reloadOpts ...QueryOption) error {
	if len(labelIDs) == 0 && len(attachmentIDs) == 0 {
		return cardRepo.db.Create(card).Error
	}

	return cardRepo.db.Transaction(func(transaction *gorm.DB) error {
		if err := transaction.Create(card).Error; err != nil {
			return err
		}

		if len(labelIDs) > 0 {
			labels := make([]*models.Label, 0)
			if err := transaction.Where("id IN ?", labelIDs).Find(&labels).Error; err != nil {
				return err
			}

			if err := transaction.Model(card).Association("Labels").Append(labels); err != nil {
				return err
			}
		}

		if len(attachmentIDs) > 0 {
			attachments := make([]*models.FileUpload, 0)
			if err := transaction.Where("id IN ?", attachmentIDs).Find(&attachments).Error; err != nil {
				return err
			}

			if err := transaction.Model(card).Association("Attachments").Append(attachments); err != nil {
				return err
			}
		}

		return applyOptions(transaction, reloadOpts).First(card, "id = ?", card.ID).Error
	})
}

func (cardRepo *CardRepository) Update(card *models.Card, opts ...QueryOption) error {
	return applyOptions(cardRepo.db, opts).Save(card).Error
}

func (cardRepo *CardRepository) UpdateWithAssociations(card *models.Card, labelIDs *[]string, attachmentIDs *[]string, reloadOpts ...QueryOption) error {
	if labelIDs == nil && attachmentIDs == nil {
		return cardRepo.db.Save(card).Error
	}

	return cardRepo.db.Transaction(func(transaction *gorm.DB) error {
		if err := transaction.Save(card).Error; err != nil {
			return err
		}

		if labelIDs != nil {
			labels := make([]*models.Label, 0)
			if err := transaction.Where("id IN ?", *labelIDs).Find(&labels).Error; err != nil {
				return err
			}

			if err := transaction.Model(card).Association("Labels").Replace(labels); err != nil {
				return err
			}
		}

		if attachmentIDs != nil {
			attachments := make([]*models.FileUpload, 0)
			if err := transaction.Where("id IN ?", *attachmentIDs).Find(&attachments).Error; err != nil {
				return err
			}

			if err := transaction.Model(card).Association("Attachments").Replace(attachments); err != nil {
				return err
			}
		}

		return applyOptions(transaction, reloadOpts).First(card, "id = ?", card.ID).Error
	})
}

func (cardRepo *CardRepository) Delete(ID string) error {
	return cardRepo.db.Where("id = ?", ID).Delete(models.Card{}).Error
}

func (cardRepo *CardRepository) Exists(ID string) bool {
	exists := false

	cardRepo.db.Raw("SELECT EXISTS(SELECT 1 FROM cards WHERE id = ?) AS found", ID).Find(&exists)

	return exists
}
