package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

var BoardFullLoad = []QueryOption{
	WithPreload("User", "Permissions", "Workspace", "Lists"),
	WithOmit("UserID", "WorkspaceID"),
}

type BoardRepository struct {
	db *gorm.DB
}

func NewBoardRepository(db *gorm.DB) *BoardRepository {
	return &BoardRepository{db: db}
}

func (boardRepo *BoardRepository) FindAll(opts ...QueryOption) ([]*models.Board, error) {
	boards := make([]*models.Board, 0)
	if err := applyOptions(boardRepo.db, opts).Find(&boards).Error; err != nil {
		return nil, err
	}

	return boards, nil
}

func (boardRepo *BoardRepository) FindByID(ID string, opts ...QueryOption) (*models.Board, error) {
	board := new(models.Board)
	if err := applyOptions(boardRepo.db, opts).Where("id = ?", ID).First(board).Error; err != nil {
		return nil, err
	}

	return board, nil
}

func (boardRepo *BoardRepository) FindByUserID(ID string, opts ...QueryOption) ([]*models.Board, error) {
	boards := make([]*models.Board, 0)

	if err := applyOptions(boardRepo.db, opts).Where("user_id = ?", ID).Find(&boards).Error; err != nil {
		return nil, err
	}

	return boards, nil
}

func (boardRepo *BoardRepository) Create(board *models.Board, opts ...QueryOption) error {
	return boardRepo.db.Transaction(func(transaction *gorm.DB) error {
		if err := applyOptions(transaction, opts).Create(board).Error; err != nil {
			return err
		}

		lists := []*models.List{
			{Name: "To Do", Position: 1, BoardID: board.ID},
			{Name: "In Progress", Position: 2, BoardID: board.ID},
			{Name: "Under Review", Position: 3, BoardID: board.ID},
			{Name: "Complete", Position: 4, BoardID: board.ID},
		}

		return transaction.Model(board).Association("Lists").Append(lists)
	})
}

func (boardRepo *BoardRepository) CreateWithPermissions(board *models.Board, permissionIDs []string, reloadOpts ...QueryOption) error {
	return boardRepo.db.Transaction(func(transaction *gorm.DB) error {
		permissions := make([]*models.Permission, 0)
		if err := transaction.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			return err
		}

		if err := transaction.Create(board).Error; err != nil {
			return err
		}

		lists := []*models.List{
			{Name: "To Do", Position: 1, BoardID: board.ID},
			{Name: "In Progress", Position: 2, BoardID: board.ID},
			{Name: "Under Review", Position: 3, BoardID: board.ID},
			{Name: "Complete", Position: 4, BoardID: board.ID},
		}

		if err := transaction.Model(board).Association("Lists").Append(lists); err != nil {
			return err
		}

		if len(permissions) > 0 {
			if err := transaction.Model(board).Association("Permissions").Append(permissions); err != nil {
				return err
			}
		}

		return applyOptions(transaction, reloadOpts).First(board, "id = ?", board.ID).Error
	})
}

func (boardRepo *BoardRepository) Update(board *models.Board, opts ...QueryOption) error {
	return applyOptions(boardRepo.db, opts).Save(board).Error
}

func (boardRepo *BoardRepository) UpdateWithPermissions(board *models.Board, permissionIDs []string, reloadOpts ...QueryOption) error {
	return boardRepo.db.Transaction(func(transaction *gorm.DB) error {
		permissions := make([]*models.Permission, 0)
		if err := transaction.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			return err
		}

		if err := transaction.Save(board).Error; err != nil {
			return err
		}

		if err := transaction.Model(board).Association("Permissions").Replace(permissions); err != nil {
			return err
		}

		return applyOptions(transaction, reloadOpts).First(board, "id = ?", board.ID).Error
	})
}

func (boardRepo *BoardRepository) Delete(ID string) error {
	return boardRepo.db.Where("id = ?", ID).Delete(models.Board{}).Error
}

func (boardRepo *BoardRepository) Exists(ID string) bool {
	exists := false

	boardRepo.db.Raw("SELECT EXISTS(SELECT 1 FROM boards WHERE id = ?) AS found", ID).Find(&exists)

	return exists
}

func (boardRepo *BoardRepository) ExistsForUserByName(ID string, name string) bool {
	exists := false

	boardRepo.db.Raw("SELECT EXISTS(SELECT 1 FROM boards WHERE id = ? AND name = ?) AS found", ID, name).Find(&exists)

	return exists
}
