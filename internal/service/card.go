package service

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/http/validators"
)

type CardService struct {
	cardRepo *repository.CardRepository
}

func NewCardService(cardRepo *repository.CardRepository) *CardService {
	return &CardService{cardRepo: cardRepo}
}

func (cardService *CardService) GetCards() ([]*models.Card, error) {
	return cardService.cardRepo.FindAll(repository.QueryOptions{
		Preload: []string{"List", "Labels", "Comments", "Comments.User"},
	})
}

func (cardService *CardService) GetCardByID(ID string) (*models.Card, error) {
	return cardService.cardRepo.FindByID(ID, repository.QueryOptions{
		Preload: []string{"List", "Labels", "Comments", "Comments.User", "Activities", "Attachments"},
	})
}

func (cardService *CardService) CreateCard(data *validators.CreateCardValidator) (*models.Card, error) {
	cards, err := cardService.cardRepo.FindAll(repository.QueryOptions{
		Select: []string{"position"},
		Omit:   []string{"Comments", "Labels", "Activities", "Attachments"},
	})
	if err != nil {
		return nil, err
	}

	maxPosition := 0
	for _, card := range cards {
		if card.ListID == data.ListID && card.Position > maxPosition {
			maxPosition = card.Position
		}
	}

	card := &models.Card{
		Name:        data.Name,
		Description: data.Description,
		ListID:      data.ListID,
		Position:    maxPosition + 1,
		IsActive:    true,
		Color:       data.Color,
	}

	if data.DueDate != nil {
		card.DueDate = data.DueDate
	}

	var labelIDs []string
	var attachmentIDs []string

	if data.LabelIDs != nil {
		labelIDs = *data.LabelIDs
	}

	if data.AttachmentIDs != nil {
		attachmentIDs = *data.AttachmentIDs
	}

	if err := cardService.cardRepo.CreateWithAssociations(
		card,
		labelIDs,
		attachmentIDs,
		repository.QueryOptions{},
		repository.QueryOptions{},
		repository.QueryOptions{},
	); err != nil {
		return nil, err
	}

	return card, nil
}

func (cardService *CardService) UpdateCard(ID string, data *validators.UpdateCardValidator) (*models.Card, error) {
	card, err := cardService.cardRepo.FindByID(ID, repository.QueryOptions{})
	if err != nil {
		return nil, err
	}

	if data.Name != nil && *data.Name != card.Name {
		card.Name = *data.Name
	}

	if data.Description != nil && data.Description != card.Description {
		card.Description = data.Description
	}

	if data.Color != nil && data.Color != card.Color {
		card.Color = data.Color
	}

	if data.DueDate != nil && data.DueDate != card.DueDate {
		card.DueDate = data.DueDate
	}

	if data.IsActive != nil && *data.IsActive != card.IsActive {
		card.IsActive = *data.IsActive
	}

	if err := cardService.cardRepo.UpdateWithAssociations(
		card,
		data.LabelIDs,
		data.AttachmentIDs,
		repository.QueryOptions{},
		repository.QueryOptions{},
		repository.QueryOptions{},
	); err != nil {
		return nil, err
	}

	return card, nil
}

func (cardService *CardService) MoveCard(ID string, data *validators.MoveCardValidator) (*models.Card, error) {
	card, err := cardService.cardRepo.FindByID(ID, repository.QueryOptions{})
	if err != nil {
		return nil, err
	}

	if data.ListID != nil && *data.ListID != card.ListID {
		card.ListID = *data.ListID
	}

	if data.Position != nil && *data.Position != card.Position {
		card.Position = *data.Position
	}

	if err := cardService.cardRepo.Update(card, repository.QueryOptions{}); err != nil {
		return nil, err
	}

	return card, nil
}

func (cardService *CardService) DeleteCard(ID string) error {
	return cardService.cardRepo.Delete(ID)
}
