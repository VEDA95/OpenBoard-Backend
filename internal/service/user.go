package service

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/http/validators"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService() (*UserService, error) {
	userRepository, err := repository.NewUserRepository()
	if err != nil {
		return nil, err
	}

	return &UserService{repo: userRepository}, nil
}

func (userService *UserService) GetUsers() ([]*models.User, error) {
	users, err := userService.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (userService *UserService) GetUser(ID string) (*models.User, error) {
	user, err := userService.repo.FindByID(ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userService *UserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := userService.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userService *UserService) CreateUser(data *validators.CreateUserValidator) (*models.User, error) {
}

func (userService *UserService) UpdateUser(data *validators.UpdateUserValidator) (*models.User, error) {
}

func (userService *UserService) DeleteUser(ID string) error {
	return userService.repo.Delete(ID)
}
