package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

var CommentFullLoad = []QueryOption{
	WithPreload("User", "Card"),
}

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (commentRepo *CommentRepository) FindAll(opts ...QueryOption) ([]*models.Comment, error) {
	comments := make([]*models.Comment, 0)
	if err := applyOptions(commentRepo.db, opts).Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (commentRepo *CommentRepository) FindByUserID(ID string, opts ...QueryOption) ([]*models.Comment, error) {
	comments := make([]*models.Comment, 0)
	if err := applyOptions(commentRepo.db, opts).Where("user_id = ?", ID).Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (commentRepo *CommentRepository) FindByCardID(ID string, opts ...QueryOption) ([]*models.Comment, error) {
	comments := make([]*models.Comment, 0)
	if err := applyOptions(commentRepo.db, opts).Where("card_id = ?", ID).Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (commentRepo *CommentRepository) FindByID(ID string, opts ...QueryOption) (*models.Comment, error) {
	comment := new(models.Comment)
	if err := applyOptions(commentRepo.db, opts).Where("id = ?", ID).First(comment).Error; err != nil {
		return nil, err
	}

	return comment, nil
}

func (commentRepo *CommentRepository) Create(comment *models.Comment, opts ...QueryOption) error {
	return applyOptions(commentRepo.db, opts).Create(comment).Error
}

func (commentRepo *CommentRepository) Update(comment *models.Comment, opts ...QueryOption) error {
	return applyOptions(commentRepo.db, opts).Save(comment).Error
}

func (commentRepo *CommentRepository) Delete(ID string) error {
	return commentRepo.db.Where("id = ?", ID).Delete(models.Comment{}).Error
}

func (commentRepo *CommentRepository) Exists(ID string) bool {
	exists := false

	commentRepo.db.Raw("SELECT EXISTS(SELECT 1 FROM comments WHERE id = ?) AS found", ID).Find(&exists)

	return exists
}
