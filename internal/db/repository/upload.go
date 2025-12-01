package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

type FileUploadRepository struct {
	db *gorm.DB
}

func NewFileUploadRepository(db *gorm.DB) *FileUploadRepository {
	return &FileUploadRepository{db: db}
}

func (fileUploadRepo *FileUploadRepository) FindAll(options QueryOptions) ([]*models.FileUpload, error) {
	fileUploads := make([]*models.FileUpload, 0)
	if err := options.AppendToQuery(fileUploadRepo.db).Find(fileUploads).Error; err != nil {
		return nil, err
	}

	return fileUploads, nil
}

func (fileUploadRepo *FileUploadRepository) FindByID(ID string, options QueryOptions) (*models.FileUpload, error) {
	fileUpload := new(models.FileUpload)
	if err := options.AppendToQuery(fileUploadRepo.db).Where("id = ?", ID).First(fileUpload).Error; err != nil {
		return nil, err
	}

	return fileUpload, nil
}

func (fileUploadRepo *FileUploadRepository) FindByUserID(ID string, options QueryOptions) (*models.FileUpload, error) {
	fileUpload := new(models.FileUpload)
	if err := options.AppendToQuery(fileUploadRepo.db).Where("user_id = ?", ID).First(fileUpload).Error; err != nil {
		return nil, err
	}

	return fileUpload, nil
}

func (fileUploadRepo *FileUploadRepository) Create(fileUpload *models.FileUpload, options QueryOptions) error {
	return options.AppendToQuery(fileUploadRepo.db).Create(fileUpload).Error
}

func (fileUploadRepo *FileUploadRepository) Update(fileUpload *models.FileUpload, options QueryOptions) error {
	return options.AppendToQuery(fileUploadRepo.db).Save(fileUpload).Error
}

func (fileUploadRepo *FileUploadRepository) Delete(ID string) error {
	return fileUploadRepo.db.Where("id = ?", ID).Delete(models.FileUpload{}).Error
}

func (fileUploadRepo *FileUploadRepository) Exists(ID string) bool {
	exists := false

	fileUploadRepo.db.Raw("SELECT EXISTS(SELECT 1 FROM file_uploads WHERE id = ?) AS found", ID).Find(&exists)

	return exists
}
