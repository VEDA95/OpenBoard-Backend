package service

import (
	"errors"

	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/http/validators"
)

type CommentService struct {
	commentRepo *repository.CommentRepository
	cardRepo    *repository.CardRepository
}

func NewCommentService(commentRepo *repository.CommentRepository, cardRepo *repository.CardRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		cardRepo:    cardRepo,
	}
}

func (commentService *CommentService) GetComments() ([]*models.Comment, error) {
	return commentService.commentRepo.FindAll(repository.CommentFullLoad...)
}

func (commentService *CommentService) GetCommentByID(ID string) (*models.Comment, error) {
	return commentService.commentRepo.FindByID(ID, repository.CommentFullLoad...)
}

func (commentService *CommentService) GetCommentsByCardID(cardID string) ([]*models.Comment, error) {
	return commentService.commentRepo.FindByCardID(cardID,
		repository.WithPreload("User"),
		repository.WithOmit("Card"),
	)
}

func (commentService *CommentService) GetCommentsByUserID(userID string) ([]*models.Comment, error) {
	return commentService.commentRepo.FindByUserID(userID,
		repository.WithPreload("Card"),
		repository.WithOmit("User"),
	)
}

func (commentService *CommentService) CreateComment(data *validators.CreateCommentValidator, userID string) (*models.Comment, error) {
	card, err := commentService.cardRepo.FindByID(data.CardID, repository.IDOnly...)
	if err != nil {
		return nil, errors.New("card not found")
	}

	comment := &models.Comment{
		Comment: data.Comment,
		CardID:  card.ID,
		UserID:  userID,
	}

	if err := commentService.commentRepo.Create(comment,
		repository.WithOmit("User", "Card"),
	); err != nil {
		return nil, err
	}

	// Load the user data for the response
	createdComment, err := commentService.commentRepo.FindByID(comment.ID,
		repository.WithPreload("User"),
		repository.WithOmit("Card"),
	)
	if err != nil {
		return nil, err
	}

	return createdComment, nil
}

func (commentService *CommentService) UpdateComment(ID string, data *validators.UpdateCommentValidator, userID string) (*models.Comment, error) {
	comment, err := commentService.commentRepo.FindByID(ID,
		repository.WithSelect("id", "user_id", "comment"),
	)
	if err != nil {
		return nil, err
	}

	// Check if user owns this comment
	if comment.UserID != userID {
		return nil, errors.New("unauthorized to update this comment")
	}

	if data.Comment != nil {
		comment.Comment = *data.Comment
	}

	if err := commentService.commentRepo.Update(comment,
		repository.WithSelect("comment", "updated_at"),
	); err != nil {
		return nil, err
	}

	// Reload with user data
	updatedComment, err := commentService.commentRepo.FindByID(ID,
		repository.WithPreload("User"),
		repository.WithOmit("Card"),
	)
	if err != nil {
		return nil, err
	}

	return updatedComment, nil
}

func (commentService *CommentService) DeleteComment(ID string, userID string) error {
	comment, err := commentService.commentRepo.FindByID(ID,
		repository.WithSelect("id", "user_id"),
	)
	if err != nil {
		return err
	}

	// Check if user owns this comment or is admin (you might want to add admin check)
	if comment.UserID != userID {
		return errors.New("unauthorized to delete this comment")
	}

	return commentService.commentRepo.Delete(ID)
}
