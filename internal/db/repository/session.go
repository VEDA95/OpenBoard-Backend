package repository

import (
	"VEDA95/open_board/api/internal/db"
	models "VEDA95/open_board/api/internal/db/model"
	"errors"

	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository() (*SessionRepository, error) {
	if db.Instance == nil {
		return nil, errors.New("database not initialized")
	}

	return &SessionRepository{db: db.Instance}, nil
}

func (sessionRepo *SessionRepository) FindAll() ([]*models.Session, error) {
	sessions := make([]*models.Session, 0)

	if err := sessionRepo.db.Find(&sessions).Error; err != nil {
		return nil, err
	}

	return sessions, nil
}

func (sessionRepo *SessionRepository) FindByUser(ID string) ([]*models.Session, error) {
	sessions := make([]*models.Session, 0)

	if err := sessionRepo.db.Where("user_id = ?", ID).Find(&sessions).Error; err != nil {
		return nil, err
	}

	return sessions, nil
}

func (sessionRepo *SessionRepository) FindByID(ID string) (*models.Session, error) {
	session := new(models.Session)

	if err := sessionRepo.db.First(session, ID).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (sessionRepo *SessionRepository) FindByAccessToken(token string) (*models.Session, error) {
	session := new(models.Session)

	if err := sessionRepo.db.Where("access_token = ?", token).First(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (sessionRepo *SessionRepository) FindByRefreshToken(token string) (*models.Session, error) {
	session := new(models.Session)

	if err := sessionRepo.db.Where("refresh_token = ?", token).First(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (sessionRepo *SessionRepository) Create(session *models.Session) error {
	return sessionRepo.db.Create(session).Error
}

func (sessionRepo *SessionRepository) Update(session *models.Session) error {
	return sessionRepo.db.Save(session).Error
}

func (sessionRepo *SessionRepository) Delete(ID string) error {
	return sessionRepo.db.Delete(&models.Session{}, ID).Error
}

func (sessionRepo *SessionRepository) DeleteByUserID(ID string) error {
	return sessionRepo.db.Where("user_id = ?", ID).Delete(&models.Session{}).Error
}
