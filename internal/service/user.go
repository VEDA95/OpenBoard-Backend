package service

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/http/validators"
	"errors"
	"time"
)

type UserService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
}

func NewUserService() (*UserService, error) {
	userRepo, err := repository.NewUserRepository()
	if err != nil {
		return nil, err
	}

	roleRepo, err := repository.NewRoleRepository()
	if err != nil {
		return nil, err
	}

	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}, nil
}

func (userService *UserService) GetUsers() ([]*models.User, error) {
	users, err := userService.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (userService *UserService) GetUser(ID string) (*models.User, error) {
	user, err := userService.userRepo.FindByID(ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userService *UserService) GetUserByUsername(username string) (*models.User, error) {
	user, err := userService.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userService *UserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := userService.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userService *UserService) CreateUser(data *validators.CreateUserValidator) (*models.User, error) {
	existingUserByUsername, err := userService.userRepo.FindByUsername(data.Username)
	if err != nil {
		return nil, err
	}

	if existingUserByUsername != nil {
		return nil, errors.New("username already exists")
	}

	existingUserByEmail, err := userService.userRepo.FindByEmail(data.Email)
	if err != nil {
		return nil, err
	}

	if existingUserByEmail != nil {
		return nil, errors.New("user with the provided email already exists")
	}

	user := new(models.User)
	user.Username = data.Username
	user.Email = data.Email

	if err := user.HashPassword(data.Password); err != nil {
		return nil, err
	}

	if data.FirstName == nil {
		user.FirstName = data.FirstName
	}

	if data.LastName == nil {
		user.LastName = data.LastName
	}

	if data.Roles != nil && len(*data.Roles) > 0 {
		if err := userService.userRepo.CreateWithRoles(user, *data.Roles...); err != nil {
			return nil, err
		}

		return user, nil
	}

	if err := userService.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (userService *UserService) UpdateUser(ID string, data *validators.UpdateUserValidator) (*models.User, error) {
	user, err := userService.userRepo.FindByID(ID)
	if err != nil {
		return nil, err
	}

	if data.Username != nil && len(*data.Username) > 0 && *data.Username != user.Username {
		user.Username = *data.Username
	}

	if data.Email != nil && len(*data.Email) > 0 && *data.Email != user.Email {
		user.Email = *data.Email
	}

	if data.FirstName != nil && data.FirstName != user.FirstName {
		if len(*data.FirstName) == 0 && user.FirstName != nil {
			user.FirstName = nil
		} else {
			user.FirstName = data.FirstName
		}
	}

	if data.LastName != nil && data.LastName != user.LastName {
		if len(*data.LastName) == 0 && user.LastName != nil {
			user.LastName = nil
		} else {
			user.LastName = data.LastName
		}
	}

	now := time.Now()
	user.UpdatedAt = &now

	if data.Roles != nil {
		if err := userService.userRepo.UpdateWithRoles(user, *data.Roles...); err != nil {
			return nil, err
		}

		return user, nil
	}

	if err := userService.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (userService *UserService) DeleteUser(ID string) error {
	return userService.userRepo.Delete(ID)
}
