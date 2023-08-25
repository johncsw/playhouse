package flow

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"playhouse-server/env"
	"playhouse-server/repo"
	"playhouse-server/util"
	"strings"
)

type transcodeVideoFlowSupport struct {
	videoID    int
	dirPath    string
	videoPath  string
	chunkCodes []int
}

func TranscodeVideo(videoID int, sessionID int) {
	err := transcodeVideo(videoID, sessionID)
	if err != nil {
		util.LogError(err, "")
	}
	log.Printf("transcoding video %d complete. sessionID=%d\n", videoID, sessionID)
}

func transcodeVideo(videoID int, sessionID int) error {
	dirPath, checkErr := repo.VideoRepo().IsVideoReadyToTranscode(videoID)
	if checkErr != nil || dirPath == "" {
		if checkErr != nil {
			return checkErr
		}
		return errors.New(fmt.Sprintf("video directory path is empty. videoID=%d sessionID=%v", videoID, sessionID))
	}

	chunkCodes, codeErr := repo.ChunkRepo().GetChunkCodeByIsUploaded(videoID, true)
	if codeErr != nil {
		return codeErr
	}
	if len(chunkCodes) == 0 {
		return errors.New(fmt.Sprintf("no chunks to transcode. videoID=%d sessionID=%v\n", videoID, sessionID))
	}

	support := &transcodeVideoFlowSupport{
		videoID:    videoID,
		dirPath:    dirPath,
		chunkCodes: chunkCodes,
	}

	assemSupport, assemErr := assembleChunks(support)
	if assemErr != nil {
		return assemErr
	}
	if assemSupport == nil {
		return errors.New(fmt.Sprintf("missing support for transcoding. videoID=%d sessionID=%v", videoID, sessionID))
	}

	transcodeErr := transcode(assemSupport)
	if transcodeErr != nil {
		return transcodeErr
	}

	return nil
}

const transcodingCommand = "ffmpeg -i %s %s"

func transcode(support *transcodeVideoFlowSupport) error {

	manifestPath := strings.Replace(support.videoPath, ".mp4", ".mpd", 1)
	cmdStr := fmt.Sprintf(transcodingCommand, support.videoPath, manifestPath)
	cmd := exec.Command(env.SHELL_PATH(), "-c", cmdStr)

	stderr, _ := cmd.StderrPipe()
	if cmdErr := cmd.Start(); cmdErr != nil {
		return cmdErr
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		log.Printf("transcoding video %d: %s\n", support.videoID, scanner.Text())
	}

	return nil
}
