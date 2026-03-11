package service

import (
	"fmt"

	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/http/validators"
)

type CheckListItemService struct {
	checkListItemRepo *repository.CheckListItemRepository
	cardRepo          *repository.CardRepository
	activityRepo      *repository.ActivityRepository
}

func NewCheckListItemService(
	checkListItemRepo *repository.CheckListItemRepository,
	cardRepo *repository.CardRepository,
	activityRepo *repository.ActivityRepository,
) *CheckListItemService {
	return &CheckListItemService{
		checkListItemRepo: checkListItemRepo,
		cardRepo:          cardRepo,
		activityRepo:      activityRepo,
	}
}

func (service *CheckListItemService) GetCheckListItems() ([]*models.CheckListItem, error) {
	return service.checkListItemRepo.FindAll(repository.WithPreload("Card"))
}

func (service *CheckListItemService) GetCheckListItemByID(ID string) (*models.CheckListItem, error) {
	return service.checkListItemRepo.FindByID(ID, repository.WithPreload("Card"))
}

func (service *CheckListItemService) GetCheckListItemsByCardID(cardID string) ([]*models.CheckListItem, error) {
	items := make([]*models.CheckListItem, 0)
	allItems, err := service.checkListItemRepo.FindAll(repository.WithOmit("Card"))
	if err != nil {
		return nil, err
	}

	for _, item := range allItems {
		if item.CardID == cardID {
			items = append(items, item)
		}
	}

	for i := 0; i < len(items)-1; i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].Position > items[j].Position {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	return items, nil
}

func (service *CheckListItemService) CreateCheckListItem(data *validators.CreateCheckListItemValidator) (*models.CheckListItem, error) {
	existingItems, err := service.GetCheckListItemsByCardID(data.CardID)
	if err != nil {
		return nil, err
	}

	maxPosition := 0
	for _, item := range existingItems {
		if item.Position > maxPosition {
			maxPosition = item.Position
		}
	}

	checkListItem := &models.CheckListItem{
		Name:      data.Name,
		CardID:    data.CardID,
		IsChecked: false,
		Position:  maxPosition + 1,
	}

	if data.Position != nil {
		checkListItem.Position = *data.Position
	}

	if err := service.checkListItemRepo.Create(checkListItem); err != nil {
		return nil, err
	}

	activity := &models.CardActivity{
		Activity: fmt.Sprintf("added checklist item: %s", data.Name),
		UserID:   data.UserID,
		CardID:   data.CardID,
	}
	if err := service.activityRepo.Create(activity); err != nil {
		return nil, err
	}

	return checkListItem, nil
}

func (service *CheckListItemService) UpdateCheckListItem(ID string, userID string, data *validators.UpdateCheckListItemValidator) (*models.CheckListItem, error) {
	checkListItem, err := service.checkListItemRepo.FindByID(ID,
		repository.WithPreload("Card", "Card.List", "Card.List.Board"),
	)
	if err != nil {
		return nil, err
	}

	updated := false
	activityMsg := ""

	if data.Name != nil && *data.Name != checkListItem.Name {
		checkListItem.Name = *data.Name
		updated = true
		activityMsg = fmt.Sprintf("renamed checklist item to: %s", *data.Name)
	}

	if data.IsChecked != nil && *data.IsChecked != checkListItem.IsChecked {
		checkListItem.IsChecked = *data.IsChecked
		updated = true

		if *data.IsChecked {
			activityMsg = fmt.Sprintf("checked: %s", checkListItem.Name)
		} else {
			activityMsg = fmt.Sprintf("unchecked: %s", checkListItem.Name)
		}
	}

	if data.Position != nil {
		checkListItem.Position = *data.Position
		updated = true
	}

	if !updated {
		return checkListItem, nil
	}

	if err := service.checkListItemRepo.Update(checkListItem); err != nil {
		return nil, err
	}

	if activityMsg != "" {
		activity := &models.CardActivity{
			Activity: activityMsg,
			UserID:   userID,
			CardID:   checkListItem.CardID,
		}

		if err := service.activityRepo.Create(activity); err != nil {
			return nil, err
		}
	}

	return checkListItem, nil
}

func (service *CheckListItemService) DeleteCheckListItem(ID string, userID string) error {
	checkListItem, err := service.checkListItemRepo.FindByID(ID,
		repository.WithPreload("Card", "Card.List", "Card.List.Board"),
	)
	if err != nil {
		return err
	}

	if err := service.checkListItemRepo.Delete(ID); err != nil {
		return err
	}

	activity := &models.CardActivity{
		Activity: fmt.Sprintf("removed checklist item: %s", checkListItem.Name),
		UserID:   userID,
		CardID:   checkListItem.CardID,
	}

	if err := service.activityRepo.Create(activity); err != nil {
		return err
	}

	return nil
}

func (service *CheckListItemService) ReorderCheckListItems(ID string, data *validators.ReorderCheckListItemsValidator) error {
	if err := service.checkListItemRepo.UpdatePositions(data.CardID, data.Positions); err != nil {
		return err
	}

	activity := &models.CardActivity{
		Activity: "reordered check list items",
		UserID:   ID,
		CardID:   data.CardID,
	}

	if err := service.activityRepo.Create(activity); err != nil {
		return err
	}

	return nil
}
