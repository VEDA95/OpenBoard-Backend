package service

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/http/validators"
)

type ListService struct {
	listRepo *repository.ListRepository
}

func NewListService(listRepo *repository.ListRepository) *ListService {
	return &ListService{listRepo: listRepo}
}

func (listService *ListService) GetLists() ([]*models.List, error) {
	return listService.listRepo.FindAll(repository.ListFullLoad...)
}

func (listService *ListService) GetListByID(ID string) (*models.List, error) {
	return listService.listRepo.FindByID(ID, repository.ListFullLoad...)
}

func (listService *ListService) GetListsByBoardID(boardID string) ([]*models.List, error) {
	return listService.listRepo.FindByBoardID(boardID, repository.WithPreload("Cards"))
}

func (listService *ListService) CreateList(data *validators.CreateListValidator) (*models.List, error) {
	lists, err := listService.listRepo.FindByBoardID(data.BoardID)
	maxPosition := 0

	if err != nil {
		return nil, err
	}

	for _, list := range lists {
		if list.Position > maxPosition {
			maxPosition = list.Position
		}
	}

	list := &models.List{
		Name:     data.Name,
		BoardID:  data.BoardID,
		Position: maxPosition + 1,
		Color:    data.Color,
	}

	if err := listService.listRepo.Create(list); err != nil {
		return nil, err
	}

	return list, nil
}

func (listService *ListService) UpdateList(ID string, data *validators.UpdateListValidator) (*models.List, error) {
	list, err := listService.listRepo.FindByID(ID)
	if err != nil {
		return nil, err
	}

	if data.Name != nil {
		list.Name = *data.Name
	}

	if data.Color != nil {
		list.Color = data.Color
	}

	if data.Position != nil {
		list.Position = *data.Position
	}

	if err := listService.listRepo.Update(list); err != nil {
		return nil, err
	}

	return list, nil
}

func (listService *ListService) DeleteList(ID string) error {
	return listService.listRepo.Delete(ID)
}

func (listService *ListService) ReorderLists(boardID string, positions map[string]int) error {
	return listService.listRepo.UpdatePositions(boardID, positions)
}
