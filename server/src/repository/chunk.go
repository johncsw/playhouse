package repository

import (
	"gorm.io/gorm"
	"playhouse-server/model"
)

type chunkrepo struct {
	db *gorm.DB
}

func (r *chunkrepo) NewChunks(v *model.Video, tx *gorm.DB) error {
	var executor *gorm.DB
	if tx == nil {
		executor = r.db
	} else {
		executor = tx
	}

	for i := int32(0); i < v.PendingChunks; i++ {
		chunkCode := int(i)
		err := executor.Create(&model.Chunk{
			VideoID:   v.ID,
			SessionID: v.SessionID,
			Code:      chunkCode,
		}).Error

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *chunkrepo) GetUnUploadedChunkCode(videoID int, sessionID int) ([]int, error) {
	var chunks []model.Chunk
	err := r.db.Select("code").Where("video_id = ? AND session_id = ? AND is_uploaded = false", videoID, sessionID).Find(&chunks).Error
	codes := make([]int, len(chunks))
	for i, c := range chunks {
		codes[i] = c.Code
	}
	return codes, err
}
