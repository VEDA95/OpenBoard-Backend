package service

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/http/validators"
	"errors"
	"time"

	"gorm.io/gorm"
)

type UserService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
}

func NewUserService(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (userService *UserService) GetUsers() ([]*models.User, error) {
	users, err := userService.userRepo.FindAll(repository.QueryOptions{
		Preload: []string{"Roles", "Roles.Permissions"},
		Omit:    []string{"Sessions"},
	})
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (userService *UserService) GetUser(ID string) (*models.User, error) {
	user, err := userService.userRepo.FindByID(ID, repository.QueryOptions{
		Preload: []string{"Roles", "Roles.Permissions"},
		Omit:    []string{"Sessions"},
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userService *UserService) GetUserByUsername(username string) (*models.User, error) {
	user, err := userService.userRepo.FindByUsername(username, repository.QueryOptions{
		Preload: []string{"Roles", "Roles.Permissions"},
		Omit:    []string{"Sessions"},
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userService *UserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := userService.userRepo.FindByEmail(email, repository.QueryOptions{
		Preload: []string{"Roles", "Roles.Permissions"},
		Omit:    []string{"Sessions"},
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userService *UserService) CreateUser(data *validators.CreateUserValidator) (*models.User, error) {
	existingUserByUsername, err := userService.userRepo.FindByUsername(data.Username, repository.QueryOptions{
		Select: []string{"id"},
	})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existingUserByUsername != nil {
		return nil, errors.New("username already exists")
	}

	existingUserByEmail, err := userService.userRepo.FindByEmail(data.Email, repository.QueryOptions{
		Select: []string{"id"},
	})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
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
		err := userService.userRepo.CreateWithRoles(
			user,
			*data.Roles,
			repository.QueryOptions{
				Omit: []string{"created_at", "updated_at", "Roles", "Sessions"},
			},
			repository.QueryOptions{
				Omit: []string{"Roles.Permissions"},
			},
		)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	err2 := userService.userRepo.Create(user, repository.QueryOptions{
		Omit: []string{"Sessions"},
	})

	if err2 != nil {
		return nil, err2
	}

	return user, nil
}

func (userService *UserService) UpdateUser(ID string, data *validators.UpdateUserValidator) (*models.User, error) {
	user, err := userService.userRepo.FindByID(ID, repository.QueryOptions{
		Preload: []string{"Roles.Permissions"},
		Omit:    []string{"Sessions"},
	})
	if err != nil {
		return nil, err
	}

	columns := []string{"updated_at"}

	if data.Username != nil && len(*data.Username) > 0 && *data.Username != user.Username {
		user.Username = *data.Username
		columns = append(columns, "username")
	}

	if data.Email != nil && len(*data.Email) > 0 && *data.Email != user.Email {
		user.Email = *data.Email
		columns = append(columns, "email")
	}

	if data.FirstName != nil && data.FirstName != user.FirstName {
		if len(*data.FirstName) == 0 && user.FirstName != nil {
			user.FirstName = nil
		} else {
			user.FirstName = data.FirstName
		}

		columns = append(columns, "first_name")
	}

	if data.LastName != nil && data.LastName != user.LastName {
		if len(*data.LastName) == 0 && user.LastName != nil {
			user.LastName = nil
		} else {
			user.LastName = data.LastName
		}

		columns = append(columns, "last_name")
	}

	now := time.Now()
	user.UpdatedAt = &now

	if data.Roles != nil {
		err := userService.userRepo.UpdateWithRoles(
			user,
			*data.Roles,
			repository.QueryOptions{
				Select: columns,
			},
			repository.QueryOptions{
				Omit: []string{"Roles.Permissions"},
			},
		)
		if err != nil {
			return nil, err
		}

		return user, nil
	}
	err2 := userService.userRepo.Update(user, repository.QueryOptions{
		Select: columns,
	})
	if err2 != nil {
		return nil, err2
	}

	return user, nil
}

func (userService *UserService) DeleteUser(ID string) error {
	return userService.userRepo.Delete(ID)
}
