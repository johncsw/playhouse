package repository

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"playhouse-server/model"
	"playhouse-server/util"
	"time"
)

type sessionrepo struct {
	db *gorm.DB
}

func (r *sessionrepo) NewSession() *model.Session {
	now := time.Now().UTC()
	env := util.NewEnv()
	sessionTTLHour := env.SessionTTLHour()
	due := now.Add(time.Hour * time.Duration(sessionTTLHour))

	s := &model.Session{
		IsAvailable: true,
		DueAt:       &due,
		CreatedAt:   &now,
	}
	result := r.db.Create(s)

	if err := result.Error; err != nil {
		panic(util.ResponseErr{
			Code:    http.StatusInternalServerError,
			ErrBody: err,
		})
	}

	return s
}

func (r *sessionrepo) IsSessionAvailable(ID int) bool {
	if ID <= 0 {
		return false
	}

	var s model.Session

	if err := r.db.Find(&s, ID).Error; err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = http.StatusNotFound
		}

		panic(util.ResponseErr{
			Code:    code,
			ErrBody: err,
		})
	}

	isNotDue := s.DueAt.After(time.Now().UTC())

	return isNotDue && s.IsAvailable
}
