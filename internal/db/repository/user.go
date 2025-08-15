package repository

import (
	"VEDA95/open_board/api/internal/db"
	models "VEDA95/open_board/api/internal/db/model"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() (*UserRepository, error) {
	if db.Instance == nil {
		return nil, errors.New("database instance not initialized")
	}

	return &UserRepository{db: db.Instance}, nil
}

func (userRepo *UserRepository) FindAll() ([]*models.User, error) {
	users := make([]*models.User, 0)

	if err := userRepo.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (userRepo *UserRepository) FindByID(ID string) (*models.User, error) {
	user := new(models.User)

	if err := userRepo.db.First(user, ID).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (userRepo *UserRepository) FindByEmail(email string) (*models.User, error) {
	user := new(models.User)

	if err := userRepo.db.First(user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (userRepo *UserRepository) Create(user *models.User) error {
	return userRepo.db.Create(user).Error
}

func (userRepo *UserRepository) Update(user *models.User) error {
	return userRepo.db.Save(user).Error
}

func (userRepo *UserRepository) Delete(ID string) error {
	return userRepo.db.Delete(&models.User{}, ID).Error
}
