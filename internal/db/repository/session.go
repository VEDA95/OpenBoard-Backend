package repository

import (
	models "VEDA95/open_board/api/internal/db/model"

	"gorm.io/gorm"
)

var SessionFullLoad = []QueryOption{
	WithPreload("User", "User.Roles", "User.Roles.Permissions"),
	WithOmit("User.Sessions"),
}

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (sessionRepo *SessionRepository) FindAll(opts ...QueryOption) ([]*models.Session, error) {
	sessions := make([]*models.Session, 0)

	if err := applyOptions(sessionRepo.db, opts).Find(&sessions).Error; err != nil {
		return nil, err
	}

	return sessions, nil
}

func (sessionRepo *SessionRepository) FindByUser(ID string, opts ...QueryOption) ([]*models.Session, error) {
	sessions := make([]*models.Session, 0)

	if err := applyOptions(sessionRepo.db, opts).Where("user_id = ?", ID).Find(&sessions).Error; err != nil {
		return nil, err
	}

	return sessions, nil
}

func (sessionRepo *SessionRepository) FindByID(ID string, opts ...QueryOption) (*models.Session, error) {
	session := new(models.Session)

	if err := applyOptions(sessionRepo.db, opts).Where("id = ?", ID).First(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (sessionRepo *SessionRepository) FindByAccessToken(token string, opts ...QueryOption) (*models.Session, error) {
	session := new(models.Session)

	if err := applyOptions(sessionRepo.db, opts).Where("access_token = ?", token).First(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (sessionRepo *SessionRepository) FindByRefreshToken(token string, opts ...QueryOption) (*models.Session, error) {
	session := new(models.Session)

	if err := applyOptions(sessionRepo.db, opts).Where("refresh_token = ?", token).First(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (sessionRepo *SessionRepository) Create(session *models.Session, opts ...QueryOption) error {
	return applyOptions(sessionRepo.db, opts).Create(session).Error
}

func (sessionRepo *SessionRepository) Update(session *models.Session, opts ...QueryOption) error {
	return applyOptions(sessionRepo.db, opts).Save(session).Error
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
