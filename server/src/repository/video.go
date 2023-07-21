package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math"
	"playhouse-server/model"
	"playhouse-server/requestbody"
	"playhouse-server/util"
	"time"
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

func (r *videorepo) GetPendingUploadVideo(videoID int, sessionID int) (*model.Video, error) {
	var v model.Video
	result := r.db.Model(&model.Video{}).Select("pending_chunks").Where("id = ? AND session_id = ? AND is_deleted = false AND pending_chunks > 0 AND is_transcoded = false", videoID, sessionID).First(&v)
	return &v, result.Error
}

func (r *videorepo) UpdatePendingChunks(videoID int, pendingChunks int32, tx *gorm.DB) error {
	var executor *gorm.DB
	if tx == nil {
		executor = r.db
	} else {
		executor = tx
	}

	now := time.Now().UTC()

	// GORM wouldn't update filed if the value is 0, so we use map to update the field
	updatedFields := map[string]interface{}{"pending_chunks": pendingChunks}
	if pendingChunks == 0 {
		updatedFields["uploaded_at"] = &now
	}
	result := executor.Model(&model.Video{}).Where("id = ?", videoID).Updates(updatedFields)

	err := result.Error
	if err != nil {
		return err
	}

	notUpdatedRight := result.RowsAffected != 1
	if notUpdatedRight {
		return errors.New(fmt.Sprintf("the video is not updated correctly. updated rows: %v", result.RowsAffected))
	}

	return nil
}

func (r *videorepo) GetVideoURLToStream(videoID int) (string, error) {
	var v model.Video
	err := r.db.Model(&model.Video{}).Select("url_to_stream").Where("id = ? ", videoID).First(&v).Error
	return v.URLToStream, err
}

func (r *videorepo) SetVideoURLToStream(videoID int, url string) error {
	// update video url to stream
	result := r.db.Model(&model.Video{}).Where("id = ?", videoID).Update("url_to_stream", url)
	err := result.Error
	if err != nil {
		return err
	}
	updatedRows := result.RowsAffected
	if updatedRows != 1 {
		errMsg := fmt.Sprintf("url stream of the video is not updated correctly. videoID=%v updatedRows=%v", videoID, updatedRows)
		return errors.New(errMsg)
	}

	return nil
}

func (r *videorepo) IsVideoReadyToTranscode(videoID int) (string, error) {
	var v model.Video
	err := r.db.Model(&model.Video{}).Select("url_to_stream").Where("id = ? AND is_deleted = false AND is_transcoded=false AND type='video/mp4' AND length(url_to_stream) > 0 AND pending_chunks = 0", videoID).First(&v).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return v.URLToStream, nil
}
