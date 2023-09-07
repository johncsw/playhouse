package repo

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"playhouse-server/env"
	"playhouse-server/model"
	"playhouse-server/response"
	"time"
)

type sessionrepo struct {
	db *gorm.DB
}

func (r *sessionrepo) NewSession(usrID int) (*model.Session, error) {
	now := time.Now().UTC()
	sessionTTLHour := env.SESSION_TTL_HOUR()
	due := now.Add(time.Hour * time.Duration(sessionTTLHour))

	s := &model.Session{
		IsAvailable: true,
		DueAt:       &due,
		CreatedAt:   &now,
		UserID:      usrID,
	}
	err := db.Create(s).Error

	return s, err
}

func (r *sessionrepo) IsSessionAvailable(ID int) bool {
	if ID <= 0 {
		return false
	}

	var s model.Session

	if err := db.First(&s, ID).Error; err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = http.StatusForbidden
		}

		panic(response.Error{
			Code:  code,
			Cause: errors.New("session not found"),
		})
	}

	isDue := time.Now().UTC().After(*s.DueAt)
	if isDue {
		panic(response.Error{
			Code:  http.StatusForbidden,
			Cause: errors.New("session expired"),
		})
	}

	return s.IsAvailable
}

func (r *sessionrepo) DeleteAll(tx *gorm.DB) error {
	var executor *gorm.DB
	if tx == nil {
		executor = db
	} else {
		executor = tx
	}

	return executor.Where("id > 0").Delete(&model.Session{}).Error
}
