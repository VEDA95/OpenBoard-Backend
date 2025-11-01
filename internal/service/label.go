package service

import (
	"errors"

	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/http/validators"
)

type LabelService struct {
	labelRepo *repository.LabelRepository
	boardRepo *repository.BoardRepository
}

func NewLabelService(
	labelRepo *repository.LabelRepository,
	boardRepo *repository.BoardRepository,
) *LabelService {
	return &LabelService{
		labelRepo: labelRepo,
		boardRepo: boardRepo,
	}
}

func (labelService *LabelService) GetLabels() ([]*models.Label, error) {
	return labelService.labelRepo.FindAll(repository.QueryOptions{
		Preload: []string{"Board"},
	})
}

func (labelService *LabelService) GetLabelByID(ID string) (*models.Label, error) {
	return labelService.labelRepo.FindByID(ID, repository.QueryOptions{
		Preload: []string{"Board"},
	})
}

func (labelService *LabelService) GetLabelsByBoardID(boardID string) ([]*models.Label, error) {
	return labelService.labelRepo.FindBoardID(boardID, repository.QueryOptions{})
}

func (labelService *LabelService) CreateLabel(data *validators.CreateLabelValidator) (*models.Label, error) {
	label := &models.Label{
		Name:    data.Name,
		Color:   data.Color,
		BoardID: data.BoardID,
	}

	if err := labelService.labelRepo.Create(label, repository.QueryOptions{}); err != nil {
		return nil, err
	}

	return label, nil
}

func (labelService *LabelService) UpdateLabel(ID string, data *validators.UpdateLabelValidator) (*models.Label, error) {
	label, err := labelService.labelRepo.FindByID(ID, repository.QueryOptions{})
	if err != nil {
		return nil, err
	}

	if data.Name != nil {
		label.Name = *data.Name
	}

	if data.Color != nil {
		label.Color = *data.Color
	}

	if err := labelService.labelRepo.Update(label, repository.QueryOptions{}); err != nil {
		return nil, err
	}

	return label, nil
}

func (labelService *LabelService) DeleteLabel(ID string) error {
	if !labelService.labelRepo.Exists(ID) {
		return errors.New("label does not exist")
	}

	if err := labelService.labelRepo.Delete(ID); err != nil {
		return err
	}

	return nil
}
