package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (sessionRepo *SessionRepository) FindAll(options QueryOptions) ([]*models.Session, error) {
	sessions := make([]*models.Session, 0)

	if err := options.AppendToQuery(sessionRepo.db).Find(&sessions).Error; err != nil {
		return nil, err
	}

	return sessions, nil
}

func (sessionRepo *SessionRepository) FindByUser(ID string, options QueryOptions) ([]*models.Session, error) {
	sessions := make([]*models.Session, 0)

	if err := options.AppendToQuery(sessionRepo.db).Where("user_id = ?", ID).Find(&sessions).Error; err != nil {
		return nil, err
	}

	return sessions, nil
}

func (sessionRepo *SessionRepository) FindByID(ID string, options QueryOptions) (*models.Session, error) {
	session := new(models.Session)

	if err := options.AppendToQuery(sessionRepo.db).Where("id = ?", ID).First(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (sessionRepo *SessionRepository) FindByAccessToken(token string, options QueryOptions) (*models.Session, error) {
	session := new(models.Session)

	if err := options.AppendToQuery(sessionRepo.db).Where("access_token = ?", token).First(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (sessionRepo *SessionRepository) FindByRefreshToken(token string, options QueryOptions) (*models.Session, error) {
	session := new(models.Session)

	if err := options.AppendToQuery(sessionRepo.db).Where("refresh_token = ?", token).First(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (sessionRepo *SessionRepository) Create(session *models.Session, options QueryOptions) error {
	return options.AppendToQuery(sessionRepo.db).Create(session).Error
}

func (sessionRepo *SessionRepository) Update(session *models.Session, options QueryOptions) error {
	return options.AppendToQuery(sessionRepo.db).Save(session).Error
}

func (sessionRepo *SessionRepository) Delete(ID string) error {
	return sessionRepo.db.Where("id = ?", ID).Delete(&models.Session{}).Error
}

func (sessionRepo *SessionRepository) DeleteByUserID(ID string) error {
	return sessionRepo.db.Where("user_id = ?", ID).Delete(&models.Session{}).Error
}

func (sessionRepo *SessionRepository) Exists(ID string) bool {
	exists := false

	sessionRepo.db.Raw("SELECT EXISTS(SELECT 1 FROM sessions WHERE id = ?) AS found", ID).Find(&exists)

	return exists
}

func (sessionRepo *SessionRepository) ExistsByAcessToken(token string) bool {
	exists := false

	sessionRepo.db.Raw("SELECT EXISTS(SELECT 1 FROM sessions WHERE access_token = ?) AS found", token).Find(&exists)

	return exists
}

func (sessionRepo *SessionRepository) ExistsByRefreshToken(token string) bool {
	exists := false

	sessionRepo.db.Raw("SELECT EXISTS(SELECT 1 FROM sessions WHERE refresh_token = ?) AS found", token).Find(&exists)

	return exists
}
