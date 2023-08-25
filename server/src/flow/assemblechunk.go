package flow

import (
	"fmt"
	"io"
	"os"
	"playhouse-server/util"
)

func assembleChunks(support *transcodeVideoFlowSupport) (*transcodeVideoFlowSupport, error) {

	videoPath := fmt.Sprintf("%s/%d-out.mp4", support.dirPath, support.videoID)

	_, err := os.Stat(videoPath)
	if os.IsExist(err) {
		removeErr := os.Remove(videoPath)
		if removeErr != nil {
			return nil, removeErr
		}
	}

	videoFile, copyErr := os.Create(videoPath)
	if copyErr != nil {
		return nil, copyErr
	}

	defer func() { util.LogError(videoFile.Close(), "") }()

	for _, code := range support.chunkCodes {
		chunkPath := fmt.Sprintf("%s/%d-%d.bin", support.dirPath, support.videoID, code)
		chunk, chunkErr := os.Open(chunkPath)
		if chunkErr != nil {
			return nil, chunkErr
		}

		_, copyErr = io.Copy(videoFile, chunk)
		if copyErr != nil {
			return nil, copyErr
		}

		closeErr := chunk.Close()
		if closeErr != nil {
			return nil, closeErr
		}
	}

	newSupport := transcodeVideoFlowSupport{
		videoID:   support.videoID,
		dirPath:   support.dirPath,
		videoPath: videoPath,
	}

	return &newSupport, nil
}
