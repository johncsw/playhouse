package flow

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"playhouse-server/repo"
	"playhouse-server/util"
)

func UpdatePendingChunks(videoID int, pendingChunks int) int {
	newPendingChunks, err := updatePendingChunks(videoID, pendingChunks)
	if err != nil {
		util.LogError(err, "")
	}
	return newPendingChunks
}

func updatePendingChunks(videoID int, pendingChunks int) (int, error) {
	var newPendingChunks int
	err := repo.NewTransaction(func(tx *gorm.DB) error {
		uploadedChunks, countErr := repo.ChunkRepo().GetNumberUploadedChunks(videoID, tx)
		if countErr != nil {
			return countErr
		}

		newPendingChunks = pendingChunks - uploadedChunks
		if newPendingChunks < 0 {
			return errors.New(fmt.Sprintf("new pending chunks is less than 0. videoID=%d pendingChunks=%d uploadedChunks=%d", videoID, pendingChunks, uploadedChunks))
		}

		updateErr := repo.VideoRepo().UpdatePendingChunks(videoID, int32(newPendingChunks), tx)
		if updateErr != nil {
			return updateErr
		}

		return nil
	})

	return newPendingChunks, err
}
