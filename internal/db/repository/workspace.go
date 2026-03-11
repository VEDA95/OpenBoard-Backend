package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

var WorkspaceFullLoad = []QueryOption{
	WithPreload("User", "Permissions", "Boards"),
	WithOmit("UserID"),
}

type WorkspaceRepository struct {
	db *gorm.DB
}

func NewWorkspaceRepository(db *gorm.DB) *WorkspaceRepository {
	return &WorkspaceRepository{db: db}
}

func (workspaceRepo *WorkspaceRepository) FindAll(opts ...QueryOption) ([]*models.Workspace, error) {
	workspaces := make([]*models.Workspace, 0)
	if err := applyOptions(workspaceRepo.db, opts).Find(&workspaces).Error; err != nil {
		return nil, err
	}

	return workspaces, nil
}

func (workspaceRepo *WorkspaceRepository) FindByID(ID string, opts ...QueryOption) (*models.Workspace, error) {
	workspace := new(models.Workspace)
	if err := applyOptions(workspaceRepo.db, opts).Where("id = ?", ID).First(workspace).Error; err != nil {
		return nil, err
	}

	return workspace, nil
}

func (workspaceRepo *WorkspaceRepository) FindByUserID(ID string, opts ...QueryOption) ([]*models.Workspace, error) {
	workspaces := make([]*models.Workspace, 0)

	if err := applyOptions(workspaceRepo.db, opts).Where("user_id = ?", ID).Find(&workspaces).Error; err != nil {
		return nil, err
	}

	return workspaces, nil
}

func (workspaceRepo *WorkspaceRepository) Create(workspace *models.Workspace, opts ...QueryOption) error {
	return applyOptions(workspaceRepo.db, opts).Create(workspace).Error
}

func (workspaceRepo *WorkspaceRepository) CreateWithPermissions(workspace *models.Workspace, permissionIDs []string, reloadOpts ...QueryOption) error {
	return workspaceRepo.db.Transaction(func(transaction *gorm.DB) error {
		permissions := make([]*models.Permission, 0)
		if err := transaction.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			return err
		}

		if err := transaction.Create(workspace).Error; err != nil {
			return err
		}

		if len(permissions) == 0 {
			return nil
		}

		if err := transaction.Model(workspace).Association("Permissions").Append(permissions); err != nil {
			return err
		}

		return applyOptions(transaction, reloadOpts).First(workspace, "id = ?", workspace.ID).Error
	})
}

func (workspaceRepo *WorkspaceRepository) Update(workspace *models.Workspace, opts ...QueryOption) error {
	return applyOptions(workspaceRepo.db, opts).Save(workspace).Error
}

func (workspaceRepo *WorkspaceRepository) UpdateWithPermissions(workspace *models.Workspace, permissionIDs []string, reloadOpts ...QueryOption) error {
	return workspaceRepo.db.Transaction(func(transaction *gorm.DB) error {
		permissions := make([]*models.Permission, 0)
		if err := transaction.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			return err
		}

		if err := transaction.Save(workspace).Error; err != nil {
			return err
		}

		if err := transaction.Model(workspace).Association("Permissions").Replace(permissions); err != nil {
			return err
		}

		return applyOptions(transaction, reloadOpts).First(workspace, "id = ?", workspace.ID).Error
	})
}

func (workspaceRepo *WorkspaceRepository) Delete(ID string) error {
	return workspaceRepo.db.Where("id = ?", ID).Delete(models.Workspace{}).Error
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
