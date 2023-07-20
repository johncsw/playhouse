package processor

import (
	"fmt"
	"io"
	"os"
	"playhouse-server/repository"
	"playhouse-server/util"
)

type TranscodeVideoProcessor struct {
	repoFact *repository.Factory
	videoID  int
}
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

func (p *TranscodeVideoProcessor) Process() {
	// 1. get and check metadata of the video and its chunks
	//		- input: videoID
	//		- output: videoID & chunkCodes
	// 2. open a .mp4 video file and write chunks to it
	//		- input: videoID & chunkCodes
	//		- output: path of the .mp4 file
	// 3. call FFmpeg to transcode the video and do error handling
	//		- intput: path of the .mp4 file
	//		- output: result of transcoding
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
