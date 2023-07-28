package repo

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"os"
	"playhouse-server/model"
	"playhouse-server/requestbody"
	"playhouse-server/util"
	"time"
)

type chunkrepo struct {
}

func (r *chunkrepo) NewChunks(v *model.Video, tx *gorm.DB) error {
	var executor *gorm.DB
	if tx == nil {
		executor = db
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

func (r *chunkrepo) GetChunkCodeByIsUploaded(videoID int, isUploaded bool) ([]int, error) {
	var chunks []model.Chunk
	err := db.Select("code").Where("video_id = ? AND is_uploaded = ?", videoID, isUploaded).Find(&chunks).Error
	codes := make([]int, len(chunks))
	for i, c := range chunks {
		codes[i] = c.Code
	}
	return codes, err
}

func (r *chunkrepo) SaveUploadedChunk(videoID int, urlToStream string, b *requestbody.UploadChunkWSBody, tx *gorm.DB) error {
	filePath := fmt.Sprintf("%v/%v-%v.bin", urlToStream, videoID, b.Code)
	fileErr := os.WriteFile(filePath, b.Content, 0444) // Read only to everyone
	if fileErr != nil {
		util.LogError(fileErr, "")
		return fileErr
	}

	var executor *gorm.DB
	if tx == nil {
		executor = db
	} else {
		executor = tx
	}

	now := time.Now().UTC()
	result := executor.Model(model.Chunk{}).Where("video_id = ? AND code = ?", videoID, b.Code).Updates(model.Chunk{
		Size:       b.Size,
		IsUploaded: true,
		UploadedAt: &now,
	})

	err := result.Error
	if err != nil {
		return err
	}

	notUpdatedRight := result.RowsAffected != 1
	if notUpdatedRight {
		return errors.New(fmt.Sprintf("the chunk is not updated correctly. updated rows: %v", result.RowsAffected))
	}

	return nil
}

func (r *chunkrepo) GetNumberOfNotUploadedChunks(videoID int, tx *gorm.DB) (int, error) {
	var executor *gorm.DB
	if tx == nil {
		executor = db
	} else {
		executor = tx
	}

	var count int64
	err := executor.Model(&model.Chunk{}).Where("video_id = ? AND is_uploaded = true", videoID).Count(&count).Error
	return int(count), err
}
