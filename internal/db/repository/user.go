package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

var UserFullLoad = []QueryOption{
	WithPreload("Roles", "Roles.Permissions"),
	WithOmit("Sessions"),
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (userRepo *UserRepository) FindAll(opts ...QueryOption) ([]*models.User, error) {
	users := make([]*models.User, 0)

	if err := applyOptions(userRepo.db, opts).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (userRepo *UserRepository) FindByID(ID string, opts ...QueryOption) (*models.User, error) {
	user := new(models.User)

	if err := applyOptions(userRepo.db, opts).Where("id = ?", ID).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (userRepo *UserRepository) FindByUsername(username string, opts ...QueryOption) (*models.User, error) {
	user := new(models.User)

	if err := applyOptions(userRepo.db, opts).Where("username = ?", username).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (userRepo *UserRepository) FindByEmail(email string, opts ...QueryOption) (*models.User, error) {
	user := new(models.User)

	if err := applyOptions(userRepo.db, opts).Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (userRepo *UserRepository) Create(user *models.User, opts ...QueryOption) error {
	return applyOptions(userRepo.db, opts).Create(user).Error
}

func (userRepo *UserRepository) CreateWithRoles(user *models.User, roleIDs []string, writeOpts []QueryOption, roleOpts ...QueryOption) error {
	return userRepo.db.Transaction(func(transaction *gorm.DB) error {
		roles := make([]*models.Role, 0)

		if err := applyOptions(transaction, roleOpts).Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
			return err
		}

		if err := applyOptions(transaction, writeOpts).Create(user).Error; err != nil {
			return err
		}

		return transaction.Model(user).Association("Roles").Append(roles)
	})
}

func (userRepo *UserRepository) Update(user *models.User, opts ...QueryOption) error {
	return applyOptions(userRepo.db, opts).Save(user).Error
}

func (userRepo *UserRepository) UpdateWithRoles(user *models.User, roleIDs []string, writeOpts []QueryOption, roleOpts ...QueryOption) error {
	return userRepo.db.Transaction(func(transaction *gorm.DB) error {
		roles := make([]*models.Role, 0)

		if err := applyOptions(transaction, writeOpts).Save(user).Error; err != nil {
			return err
		}

		if len(roleIDs) == 0 {
			return nil
		}

		if err := applyOptions(transaction, roleOpts).Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
			return err
		}

		return transaction.Model(user).Association("Roles").Replace(roles)
	})
}

func (userRepo *UserRepository) Delete(ID string) error {
	return userRepo.db.Delete(&models.User{}, ID).Error
}

func (userRepo *UserRepository) Exists(ID string) bool {
	exists := false

	userRepo.db.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?) AS found", ID).Find(&exists)

	return exists
}

func (userRepo *UserRepository) ExistsByUsernameOrEmail(username string, email string) bool {
	exists := false

	userRepo.db.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE username = ? OR email = ?) AS found", username, email).Find(&exists)

	return exists
}

func (userRepo *UserRepository) ExistsByUsername(username string) bool {
	exists := false

	userRepo.db.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?) AS found", username).Find(&exists)

	return exists
}
