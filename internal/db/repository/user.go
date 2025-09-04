package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (userRepo *UserRepository) FindAll(options QueryOptions) ([]*models.User, error) {
	users := make([]*models.User, 0)

	if err := options.AppendToQuery(userRepo.db).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (userRepo *UserRepository) FindByID(ID string, options QueryOptions) (*models.User, error) {
	user := new(models.User)

	if err := options.AppendToQuery(userRepo.db).Where("id = ?", ID).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (userRepo *UserRepository) FindByUsername(username string, options QueryOptions) (*models.User, error) {
	user := new(models.User)

	if err := options.AppendToQuery(userRepo.db).Where("username = ?", username).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (userRepo *UserRepository) FindByEmail(email string, options QueryOptions) (*models.User, error) {
	user := new(models.User)

	if err := options.AppendToQuery(userRepo.db).Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (userRepo *UserRepository) Create(user *models.User, options QueryOptions) error {
	return options.AppendToQuery(userRepo.db).Create(user).Error
}

func (userRepo *UserRepository) CreateWithRoles(user *models.User, roleIDs []string, userOptions QueryOptions, roleOptions QueryOptions) error {
	return userRepo.db.Transaction(func(transaction *gorm.DB) error {
		roles := make([]*models.Role, 0)

		if err := roleOptions.AppendToQuery(transaction).Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
			return err
		}

		if err := userOptions.AppendToQuery(transaction).Create(user).Error; err != nil {
			return err
		}

		return transaction.Model(user).Association("Roles").Append(roles)
	})
}

func (userRepo *UserRepository) Update(user *models.User, options QueryOptions) error {
	return options.AppendToQuery(userRepo.db).Save(user).Error
}

func (userRepo *UserRepository) UpdateWithRoles(user *models.User, roleIDs []string, userOptions QueryOptions, roleOptions QueryOptions) error {
	return userRepo.db.Transaction(func(transaction *gorm.DB) error {
		roles := make([]*models.Role, 0)

		if err := userOptions.AppendToQuery(transaction).Save(user).Error; err != nil {
			return err
		}

		if len(roleIDs) == 0 {
			return nil
		}

		if err := roleOptions.AppendToQuery(transaction).Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
			return err
		}

		return transaction.Model(user).Association("Roles").Replace(roles)
	})
}

func (userRepo *UserRepository) Delete(ID string) error {
	return userRepo.db.Delete(&models.User{}, ID).Error
}
