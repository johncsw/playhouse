package repo

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math"
	"os"
	"playhouse-server/env"
	"playhouse-server/model"
	"playhouse-server/request"
	"playhouse-server/util"
	"time"
)

type videorepo struct {
}

func (r *videorepo) NewVideo(b *request.UploadRegistrationBody, sessionID int, tx *gorm.DB) (*model.Video, error) {
	var executor *gorm.DB
	if tx == nil {
		executor = db
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
	err := db.Select("size").First(&v, id).Error
	return v.Size, err
}

func (r *videorepo) GetPendingUploadVideo(videoID int, sessionID int) (*model.Video, error) {
	var v model.Video
	result := db.Model(&model.Video{}).Select("pending_chunks").Where("id = ? AND session_id = ? AND is_deleted = false AND pending_chunks > 0 AND is_transcoded = false", videoID, sessionID).First(&v)
	return &v, result.Error
}

func (r *videorepo) UpdatePendingChunks(videoID int, pendingChunks int32, tx *gorm.DB) error {
	var executor *gorm.DB
	if tx == nil {
		executor = db
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
	err := db.Model(&model.Video{}).Select("url_to_stream").Where("id = ? ", videoID).First(&v).Error
	return v.URLToStream, err
}

func (r *videorepo) SetVideoURLToStream(videoID int, url string) error {
	// update video url to stream
	result := db.Model(&model.Video{}).Where("id = ?", videoID).Update("url_to_stream", url)
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
	err := db.Model(&model.Video{}).Select("url_to_stream").Where("id = ? AND is_deleted = false AND is_transcoded=false AND type='video/mp4' AND length(url_to_stream) > 0 AND pending_chunks = 0", videoID).First(&v).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return v.URLToStream, nil
}
func (r *videorepo) IsVideoAvailableToStream(videoID int) (string, bool, error) {
	var v model.Video
	err := db.Model(&model.Video{}).Select("url_to_stream", "is_transcoded").Where("id = ? AND is_deleted = false AND type='video/mp4' AND length(url_to_stream) > 0 AND pending_chunks = 0", videoID).First(&v).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", v.IsTranscoded, nil
	}
	if err != nil {
		return "", v.IsTranscoded, err
	}
	return v.URLToStream, v.IsTranscoded, nil
}

func (r *videorepo) UpdateVideoAsTranscoded(videoID int) error {
	result := db.Model(&model.Video{}).Where("id = ?", videoID).Update("is_transcoded", true)
	err := result.Error
	if err != nil {
		return err
	}
	updatedRows := result.RowsAffected
	if updatedRows != 1 {
		errMsg := fmt.Sprintf("the video is not marked as transcodede correctly. videoID=%v updatedRows=%v", videoID, updatedRows)
		return errors.New(errMsg)
	}

	return nil
}

func (r *videorepo) GetAllUploadedVideo(sessionID int) ([]model.Video, error) {
	var videos []model.Video
	result := db.Model(&model.Video{}).Select("id", "name").Where("session_id = ? AND is_deleted = false AND pending_chunks = 0", sessionID).Find(&videos)
	return videos, result.Error
}

func (r *videorepo) GetChunkSavingDirURL(videoID int) (string, error) {
	urlToStream, urlErr := r.GetVideoURLToStream(videoID)
	if urlErr != nil {
		return "", urlErr
	}
	return urlToStream, nil
}

func (r *videorepo) CreateChunkSavingDir(videoID int) (string, error) {
	urlToStream := fmt.Sprintf("%v/%v", env.CHUNK_STORAGE_PATH(), videoID)
	dirErr := os.Mkdir(urlToStream, 0755)
	if dirErr != nil {
		return "", dirErr
	}
	saveErr := r.SetVideoURLToStream(videoID, urlToStream)
	if saveErr != nil {
		return "", saveErr
	}

	return urlToStream, nil
}
