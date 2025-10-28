package service

import "VEDA95/open_board/api/internal/db/repository"

type BoardService struct {
	boardRepository *repository.BoardRepository
}
