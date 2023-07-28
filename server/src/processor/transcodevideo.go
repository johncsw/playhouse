package processor

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"playhouse-server/env"
	"playhouse-server/repository"
	"playhouse-server/util"
	"strings"
	"time"
)

type TranscodeVideoProcessor struct {
	repoFact *repository.Factory
	videoID  int
}

const transcodingCommand = "ffmpeg -i %s %s"

type checkVideoAndChunkOutput struct {
	passCheck  bool
	videoPath  string
	err        error
	chunkCodes []int
}

type assembleChunksOutput struct {
	videoPath string
	err       error
}

type transcodeVideoOutput struct {
	err error
}

func (p *TranscodeVideoProcessor) Process() {
	start := time.Now()
	log.Printf("Start transcoding videoID=%d\n", p.videoID)
	checkOutput := p.checkVideoAndChunks()
	didNotPass := !checkOutput.passCheck
	if didNotPass {
		if checkOutput.err != nil {
			util.LogError(checkOutput.err, "")
		}
		return
	}

	assembleOutput := p.assembleChunks(&checkOutput)
	if assembleOutput.err != nil {
		util.LogError(assembleOutput.err, "")
		return
	}

	transcodeOutput := p.transcodeVideo(&assembleOutput)
	if transcodeOutput.err != nil {
		util.LogError(transcodeOutput.err, "")
		return
	}

	markErr := p.markVideoAsTranscoded()
	if markErr != nil {
		util.LogError(markErr, "")
	}

	log.Printf("Finished transcoding video. videoID=%d elapsed=%v\n ", p.videoID, time.Since(start))
}

func (p *TranscodeVideoProcessor) checkVideoAndChunks() checkVideoAndChunkOutput {
	output := checkVideoAndChunkOutput{passCheck: false}

	videoRepo := p.repoFact.NewVideoRepo()
	URLTostream, videoErr := videoRepo.IsVideoReadyToTranscode(p.videoID)

	if videoErr != nil || URLTostream == "" {
		if videoErr != nil {
			output.err = videoErr
		}
		return output
	}

	chunkRepo := p.repoFact.NewChunkRepo()
	chunkCodes, chunkErr := chunkRepo.GetChunkCodeByIsUploaded(p.videoID, true)
	if chunkErr != nil {
		output.err = chunkErr
		return output
	}

	output.passCheck = true
	output.videoPath = URLTostream
	output.chunkCodes = chunkCodes
	return output
}

func (p *TranscodeVideoProcessor) assembleChunks(input *checkVideoAndChunkOutput) assembleChunksOutput {

	output := assembleChunksOutput{}
	outputPath := fmt.Sprintf("%s/%d-out.mp4", input.videoPath, p.videoID)

	_, err := os.Stat(outputPath)
	if os.IsExist(err) {
		removeErr := os.Remove(outputPath)
		if removeErr != nil {
			output.err = removeErr
			return output
		}
	}

	outputVideo, createErr := os.Create(outputPath)
	if createErr != nil {
		output.err = createErr
		return output
	}

	defer func() { util.LogError(outputVideo.Close(), "") }()

	for _, code := range input.chunkCodes {
		chunkPath := fmt.Sprintf("%s/%d-%d.bin", input.videoPath, p.videoID, code)
		chunk, chunkErr := os.Open(chunkPath)
		if chunkErr != nil {
			output.err = createErr
			return output
		}

		_, createErr = io.Copy(outputVideo, chunk)
		if createErr != nil {
			output.err = createErr
			return output
		}

		closeErr := chunk.Close()
		if closeErr != nil {
			output.err = closeErr
			return output
		}
	}

	output.videoPath = outputPath

	return output
}

func (p *TranscodeVideoProcessor) transcodeVideo(input *assembleChunksOutput) transcodeVideoOutput {
	output := transcodeVideoOutput{}

	outputPath := strings.Replace(input.videoPath, ".mp4", ".mpd", 1)
	cmdStr := fmt.Sprintf(transcodingCommand, input.videoPath, outputPath)
	cmd := exec.Command(env.SHELL_PATH(), "-c", cmdStr)

	stderr, _ := cmd.StderrPipe()
	if cmdErr := cmd.Start(); cmdErr != nil {
		output.err = cmdErr
		return output
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		log.Printf("transcoding video %d: %s\n", p.videoID, scanner.Text())
	}

	return output
}

func (p *TranscodeVideoProcessor) markVideoAsTranscoded() error {
	return p.repoFact.NewVideoRepo().UpdateVideoAsTranscoded(p.videoID)
}
