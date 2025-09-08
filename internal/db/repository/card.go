package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

type CardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) *CardRepository {
	return &CardRepository{db: db}
}

func (cardRepo *CardRepository) FindAll(options QueryOptions) ([]*models.Card, error) {
	cards := make([]*models.Card, 0)
	if err := options.AppendToQuery(cardRepo.db).Find(&cards).Error; err != nil {
		return nil, err
	}

	return cards, nil
}

func (cardRepo *CardRepository) FindByID(ID string, options QueryOptions) (*models.Card, error) {
	card := new(models.Card)
	if err := options.AppendToQuery(cardRepo.db).Where("id = ?", ID).First(card).Error; err != nil {
		return nil, err
	}

	return card, nil
}

func (cardRepo *CardRepository) Create(card *models.Card, options QueryOptions) error {
	return options.AppendToQuery(cardRepo.db).Create(card).Error
}

func (cardRepo *CardRepository) CreateWithAssociations(card *models.Card, labelIDs []string, attachmentIDs []string, options QueryOptions, labelOptions QueryOptions, attachmentOptions QueryOptions) error {
	if len(labelIDs) == 0 && len(attachmentIDs) == 0 {
		return options.AppendToQuery(cardRepo.db).Create(card).Error
	}

	return cardRepo.db.Transaction(func(transaction *gorm.DB) error {
		if err := options.AppendToQuery(transaction).Create(card).Error; err != nil {
			return err
		}

		if len(labelIDs) > 0 {
			labels := make([]*models.Label, 0)
			if err := labelOptions.AppendToQuery(transaction).Where("id = ?", labelIDs).Find(&labels).Error; err != nil {
				return err
			}

			if err := transaction.Model(card).Association("Labels").Append(labels); err != nil {
				return err
			}
		}

		if len(attachmentIDs) > 0 {
			attachments := make([]*models.FileUpload, 0)
			if err := attachmentOptions.AppendToQuery(transaction).Where("id IN ?", attachmentIDs).Find(&attachments).Error; err != nil {
				return err
			}

			if err := transaction.Model(card).Association("Attachments").Append(attachments); err != nil {
				return err
			}
		}

		return nil
	})
}

func (cardRepo *CardRepository) Update(card *models.Card, options QueryOptions) error {
	return options.AppendToQuery(cardRepo.db).Save(card).Error
}

func (cardRepo *CardRepository) UpdateWithAssociations(card *models.Card, labelIDs *[]string, attachmentIDs *[]string, options *QueryOptions, labelOptions *QueryOptions, attachmentOptions *QueryOptions) error {
	if labelIDs == nil && attachmentIDs == nil {
		return options.AppendToQuery(cardRepo.db).Save(card).Error
	}

	return cardRepo.db.Transaction(func(transaction *gorm.DB) error {
		if err := options.AppendToQuery(transaction).Save(card).Error; err != nil {
			return err
		}

		if labelIDs != nil {
			labels := make([]*models.Label, 0)
			if err := labelOptions.AppendToQuery(transaction).Where("id IN ?", labelIDs).First(&labels).Error; err != nil {
				return err
			}

			if err := transaction.Model(card).Association("Labels").Replace(labels); err != nil {
				return err
			}
		}

		if attachmentIDs != nil {
			attachments := make([]*models.FileUpload, 0)
			if err := attachmentOptions.AppendToQuery(cardRepo.db).Where("id IN ?", attachmentIDs).Find(&attachments).Error; err != nil {
				return err
			}

			if err := transaction.Model(card).Association("Attachments").Replace(attachments); err != nil {
				return err
			}
		}

		return nil
	})
}

func (cardRepo *CardRepository) Delete(ID string) error {
	return cardRepo.db.Where("id = ?", ID).Delete(models.Card{}).Error
}
