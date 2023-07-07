package model

import "time"

type Chunk struct {
	Code       int
	Size       int
	IsUploaded bool
	Content    []byte
	CreatedAt  *time.Time
	UploadedAt *time.Time
	VideoID    int
	SessionID  int
}

func (Chunk) TableName() string {
	return "chunk"
}
