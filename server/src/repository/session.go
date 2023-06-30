package repository

import (
	"gorm.io/gorm"
	"net/http"
	"playhouse-server/middleware"
	"playhouse-server/model"
	"playhouse-server/util"
	"time"
)

type SessionRepo struct {
	DB *gorm.DB
}

func (sr *SessionRepo) NewSession() *model.Session {
	now := time.Now().UTC()
	sessionTTLHour := util.EnvGetSessionTTLHour()
	due := now.Add(time.Hour * time.Duration(sessionTTLHour))

	s := &model.Session{
		IsAvailable: true,
		DueAt:       &due,
		CreatedAt:   &now,
	}
	result := sr.DB.Create(s)

	if err := result.Error; err != nil {
		panic(middleware.ResponseErr{
			Code:    http.StatusInternalServerError,
			ErrBody: err,
		})
	}

	return s
}
