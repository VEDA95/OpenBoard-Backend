package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

type BoardRepository struct {
	db *gorm.DB
}

func NewBoardRepository(db *gorm.DB) *BoardRepository {
	return &BoardRepository{db: db}
}

func (boardRepo *BoardRepository) FindAll(options QueryOptions) ([]*models.Board, error) {
	boards := make([]*models.Board, 0)
	if err := options.AppendToQuery(boardRepo.db).Find(&boards).Error; err != nil {
		return nil, err
	}

	return boards, nil
}

func (boardRepo *BoardRepository) FindByID(ID string, options QueryOptions) (*models.Board, error) {
	board := new(models.Board)
	if err := options.AppendToQuery(boardRepo.db).Where("id = ?", ID).First(board).Error; err != nil {
		return nil, err
	}

	return board, nil
}

func (boardRepo *BoardRepository) Create(board *models.Board, options QueryOptions) error {
	return boardRepo.db.Transaction(func(transaction *gorm.DB) error {
		if err := options.AppendToQuery(transaction).Create(board).Error; err != nil {
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

func (boardRepo *BoardRepository) CreateWithPermissions(board *models.Board, permissionIDs []string, options QueryOptions, permissionOptions QueryOptions) error {
	return boardRepo.db.Transaction(func(transaction *gorm.DB) error {
		permissions := make([]*models.Permission, 0)
		if err := permissionOptions.AppendToQuery(transaction).Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			return err
		}

		if err := options.AppendToQuery(transaction).Create(board).Error; err != nil {
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

		if len(permissions) == 0 {
			return nil
		}

		return transaction.Model(board).Association("Permissions").Append(permissions)
	})
}

func (boardRepo *BoardRepository) Update(board *models.Board, options QueryOptions) error {
	return options.AppendToQuery(boardRepo.db).Save(board).Error
}

func (boardRepo *BoardRepository) UpdateWithPermissions(board *models.Board, permissionIDs []string, options QueryOptions, permissionOptions QueryOptions) error {
	return boardRepo.db.Transaction(func(transaction *gorm.DB) error {
		permissions := make([]*models.Permission, 0)
		if err := permissionOptions.AppendToQuery(transaction).Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			return err
		}

		if err := options.AppendToQuery(transaction).Create(board).Error; err != nil {
			return err
		}

		return transaction.Model(board).Association("Permissions").Replace(permissions)
	})
}

func (boardRepo *BoardRepository) Delete(ID string) error {
	return boardRepo.db.Where("id = ?", ID).Delete(models.Board{}).Error
}
