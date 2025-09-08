package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (commentRepo *CommentRepository) FindAll(options QueryOptions) ([]*models.Comment, error) {
	comments := make([]*models.Comment, 0)
	if err := options.AppendToQuery(commentRepo.db).Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (commentRepo *CommentRepository) FindByUserID(ID string, options QueryOptions) ([]*models.Comment, error) {
	comments := make([]*models.Comment, 0)
	if err := options.AppendToQuery(commentRepo.db).Where("user_id = ?", ID).Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (commentRepo *CommentRepository) FindByCardID(ID string, options QueryOptions) ([]*models.Comment, error) {
	comments := make([]*models.Comment, 0)
	if err := options.AppendToQuery(commentRepo.db).Where("card_id = ?", ID).Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (commentRepo *CommentRepository) FindByID(ID string, options QueryOptions) (*models.Comment, error) {
	comment := new(models.Comment)
	if err := options.AppendToQuery(commentRepo.db).Where("id = ?", ID).First(comment).Error; err != nil {
		return nil, err
	}

	return comment, nil
}

func (commentRepo *CommentRepository) Create(comment *models.Comment, options QueryOptions) error {
	return options.AppendToQuery(commentRepo.db).Create(comment).Error
}

func (commentRepo *CommentRepository) Update(comment *models.Comment, options QueryOptions) error {
	return options.AppendToQuery(commentRepo.db).Save(comment).Error
}

func (commentRepo *CommentRepository) Delete(ID string) error {
	return commentRepo.db.Where("id = ?", ID).Delete(models.Comment{}).Error
}
