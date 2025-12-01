package service

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
	"errors"
	"fmt"
	"mime"
	"mime/multipart"
	"path/filepath"
)

type FileUploadService struct {
	fileUploadRepo *repository.FileUploadRepository
	userRepo       *repository.UserRepository
	cardRepo       *repository.CardRepository
}

func NewFileUploadService(fileUploadRepo *repository.FileUploadRepository, userRepo *repository.UserRepository, cardRepo *repository.CardRepository) *FileUploadService {
	return &FileUploadService{
		fileUploadRepo: fileUploadRepo,
		userRepo:       userRepo,
		cardRepo:       cardRepo,
	}
}

func (fileUploadService *FileUploadService) GetAll() ([]*models.FileUpload, error) {
	return fileUploadService.fileUploadRepo.FindAll(repository.QueryOptions{
		Preload: []string{"User"},
		Omit:    []string{"UserID"},
	})
}

func (fileUploadService *FileUploadService) GetByID(ID string) (*models.FileUpload, error) {
	if !fileUploadService.fileUploadRepo.Exists(ID) {
		return nil, errors.New("file upload entry does not exist")
	}

	return fileUploadService.fileUploadRepo.FindByID(ID, repository.QueryOptions{
		Preload: []string{"User"},
		Omit:    []string{"UserID"},
	})
}

func (fileUploadService *FileUploadService) UploadFileAsUserThumbnail(user *models.User, file *multipart.FileHeader) (*models.FileUpload, error) {
	fileExtension := filepath.Ext(file.Filename)
	fileData := &models.FileUpload{
		Name:      filepath.Base(file.Filename),
		Extension: fileExtension,
		Type:      mime.TypeByExtension(fileExtension),
		Size:      int(file.Size),
		Path:      fmt.Sprintf("uploads/%s/%s", user.ID, file.Filename),
		UserID:    user.ID,
	}

	if err := fileUploadService.fileUploadRepo.Create(fileData, repository.QueryOptions{Preload: []string{"User"}}); err != nil {
		return nil, err
	}

	user.ThumbnailID = fileData.ID

	if err := fileUploadService.userRepo.Update(user, repository.QueryOptions{Select: []string{"thumbnail_id"}}); err != nil {
		return nil, err
	}

	return fileData, nil
}

func (fileUploadService *FileUploadService) UploadFileAsCardAttachment(ID string, user *models.User, file *multipart.FileHeader) (*models.FileUpload, error) {
	if !fileUploadService.cardRepo.Exists(ID) {
		return nil, errors.New("card does not exist")
	}

	card, err := fileUploadService.cardRepo.FindByID(ID, repository.QueryOptions{
		Preload: []string{"Attachments"},
	})
	if err != nil {
		return nil, err
	}

	fileExtension := filepath.Ext(file.Filename)
	fileData := &models.FileUpload{
		Name:      filepath.Base(file.Filename),
		Extension: fileExtension,
		Type:      mime.TypeByExtension(fileExtension),
		Size:      int(file.Size),
		Path:      fmt.Sprintf("uploads/%s/%s", user.ID, file.Filename),
		UserID:    user.ID,
	}

	if err := fileUploadService.fileUploadRepo.Create(fileData, repository.QueryOptions{Preload: []string{"User"}}); err != nil {
		return nil, err
	}

	card.Attachments = append(card.Attachments, fileData)

	if err := fileUploadService.cardRepo.Update(card, repository.QueryOptions{Preload: []string{"Attachments"}}); err != nil {
		return nil, err
	}

	return fileData, nil
}

func (fileUploadService *FileUploadService) Delete(ID string) error {
	return fileUploadService.fileUploadRepo.Delete(ID)
}
