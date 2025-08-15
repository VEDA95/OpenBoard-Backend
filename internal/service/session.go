package service

import "VEDA95/open_board/api/internal/db/repository"

type SessionService struct {
	repo *repository.SessionRepository
}

func NewSessionService() (*SessionService, error) {
	sessionRepository, err := repository.NewSessionRepository()
	if err != nil {
		return nil, err
	}

	return &SessionService{repo: sessionRepository}, nil
}
