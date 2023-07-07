package model

import "time"

type Video struct {
	ID            int `gorm:"primaryKey"`
	Name          string
	Type          string
	Size          int
	URLToStream   string
	PendingChunks int32
	IsDeleted     bool
	IsTranscoded  bool
	CreatedAt     *time.Time
	UploadedAt    *time.Time
	SessionID     int
}

func (Video) TableName() string {
	return "video"
}
