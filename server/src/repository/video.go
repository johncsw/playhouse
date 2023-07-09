package repository

import (
	"gorm.io/gorm"
	"math"
	"playhouse-server/model"
	"playhouse-server/requestbody"
	"playhouse-server/util"
)

type videorepo struct {
	db *gorm.DB
}

func (r *videorepo) NewVideo(b *requestbody.UploadRegistrationBody, sessionID int, tx *gorm.DB) (*model.Video, error) {
	var executor *gorm.DB
	if tx == nil {
		executor = r.db
	} else {
		executor = tx
	}
	v := &model.Video{
		Size:      b.VideoSize,
		Name:      b.VideoName,
		Type:      b.VideoType,
		SessionID: sessionID,
	}

	maxChunkSize := util.MaxChunkSize(v.Size)
	chunkCounts := int32(math.Ceil(float64(v.Size) / float64(maxChunkSize)))
	v.PendingChunks = chunkCounts

	err := executor.Create(v).Error

	return v, err
}

func (r *videorepo) GetVideoSize(id int) (int, error) {
	var v model.Video
	err := r.db.Select("size").First(&v, id).Error
	return v.Size, err
}
