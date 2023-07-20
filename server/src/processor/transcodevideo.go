package processor

import (
	"playhouse-server/repository"
	"playhouse-server/util"
)

type TranscodeVideoProcessor struct {
	repoFact *repository.Factory
	videoID  int
}
type checkVideoAndChunkOutput struct {
	passCheck bool
	videoPath string
	err       error
	chunkCode []int
}

func (p *TranscodeVideoProcessor) Process() {
	// 1. get and check metadata of the video and its chunks
	//		- input: videoID
	//		- output: videoID & chunkCode
	// 2. open a .mp4 video file and write chunks to it
	//		- input: videoID & chunkCode
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
	output.chunkCode = chunkCodes
	return output
}
