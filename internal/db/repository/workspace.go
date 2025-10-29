package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

type WorkspaceRepository struct {
	db *gorm.DB
}

func NewWorkspaceRepository(db *gorm.DB) *WorkspaceRepository {
	return &WorkspaceRepository{db: db}
}

func (workspaceRepo *WorkspaceRepository) FindAll(options QueryOptions) ([]*models.Worksapce, error) {
	workspaces := make([]*models.Worksapce, 0)
	if err := options.AppendToQuery(workspaceRepo.db).Find(&workspaces).Error; err != nil {
		return nil, err
	}

	return workspaces, nil
}

func (workspaceRepo *WorkspaceRepository) FindByID(ID string, options QueryOptions) (*models.Worksapce, error) {
	workspace := new(models.Worksapce)
	if err := options.AppendToQuery(workspaceRepo.db).Where("id = ?", ID).First(workspace).Error; err != nil {
		return nil, err
	}

	return workspace, nil
}

func (workspaceRepo *WorkspaceRepository) FindByUserID(ID string, options QueryOptions) ([]*models.Worksapce, error) {
	workspaces := make([]*models.Worksapce, 0)

	if err := options.AppendToQuery(workspaceRepo.db).Where("user_id = ?", ID).Find(&workspaces).Error; err != nil {
		return nil, err
	}

	return workspaces, nil
}

func (workspaceRepo *WorkspaceRepository) Create(workspace *models.Worksapce, options QueryOptions) error {
	return options.AppendToQuery(workspaceRepo.db).Create(workspace).Error
}

func (workspaceRepo *WorkspaceRepository) CreateWithPermissions(workspace *models.Worksapce, permissionIDs []string, options QueryOptions, permissionOptions QueryOptions) error {
	return workspaceRepo.db.Transaction(func(transaction *gorm.DB) error {
		permissions := make([]*models.Permission, 0)
		if err := permissionOptions.AppendToQuery(transaction).Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			return err
		}

		if err := options.AppendToQuery(transaction).Create(workspace).Error; err != nil {
			return err
		}

		if len(permissions) == 0 {
			return nil
		}

		return transaction.Model(workspace).Association("Permissions").Append(permissions)
	})
}

func (workspaceRepo *WorkspaceRepository) Update(workspace *models.Worksapce, options QueryOptions) error {
	return options.AppendToQuery(workspaceRepo.db).Save(workspace).Error
}

func (workspaceRepo *WorkspaceRepository) UpdateWithPermissions(workspace *models.Worksapce, permissionIDs []string, options QueryOptions, permissionOptions QueryOptions) error {
	return workspaceRepo.db.Transaction(func(transaction *gorm.DB) error {
		permissions := make([]*models.Permission, 0)
		if err := permissionOptions.AppendToQuery(transaction).Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			return err
		}

		if err := options.AppendToQuery(transaction).Save(workspace).Error; err != nil {
			return err
		}

		return transaction.Model(workspace).Association("Permissions").Replace(permissions)
	})
}

func (workspaceRepo *WorkspaceRepository) Delete(ID string) error {
	return workspaceRepo.db.Where("id = ?", ID).Delete(models.Worksapce{}).Error
}

func (workspaceRepo *WorkspaceRepository) Exists(ID string) bool {
	exists := false

	workspaceRepo.db.Raw("SELECT EXISTS(SELECT 1 FROM workspaces WHERE id = ?) AS found", ID).Find(&exists)

	return exists
}

func (workspaceRepo *WorkspaceRepository) ExistsForUserByName(ID string, name string) bool {
	exists := false

	workspaceRepo.db.Raw("SELECT EXISTS(SELECT 1 FROM workspaces WHERE user_id = ? AND name = ?) AS found", ID, name).Find(&exists)

	return exists
}
