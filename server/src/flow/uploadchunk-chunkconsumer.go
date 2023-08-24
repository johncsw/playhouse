package flow

import (
	"fmt"
	"log"
	"playhouse-server/repo"
	"playhouse-server/request"
	"playhouse-server/response"
	"time"
)

type chunkConsumer struct {
	uploadChunkSupport *UploadChunkSupport
	consumingQueue     <-chan request.UploadChunkWebsocketBody
}

func (c *chunkConsumer) start() (<-chan response.UploadChunkWebsocketBody, <-chan error) {
	consumerErrCh := make(chan error)
	quit := make(chan error)
	successCh := make(chan response.UploadChunkWebsocketBody)
	support := c.uploadChunkSupport
	go func() {
		for i := 0; i < support.NumsOfChunks; i++ {
			select {
			case quitErr := <-quit:
				consumerErrCh <- quitErr
				return
			case chunkBody := <-c.consumingQueue:
				c.consume(chunkBody, quit, successCh)
			case <-time.After(5 * time.Second):
				consumerErrCh <- fmt.Errorf("timeout consuming chunks. videoID=%v sessionID=%v", support.VideoID, support.SessionID)
				return
			}
		}
	}()
	return successCh, consumerErrCh
}

func (c *chunkConsumer) consume(chunkBody request.UploadChunkWebsocketBody, quit chan<- error, successCh chan<- response.UploadChunkWebsocketBody) {
	go func() {
		support := c.uploadChunkSupport
		log.Printf("saving chunks. code=%v videoID=%v sessionID=%v\n", chunkBody.Code, support.VideoID, support.SessionID)
		saveErr := repo.ChunkRepo().SaveUploadedChunk(support.VideoID, support.ChunkSavingDirURL, &chunkBody, nil)
		if saveErr != nil {
			quit <- saveErr
			return
		}

		successCh <- response.UploadChunkWebsocketBody{
			Status: response.UploadChunkStatusSuccess,
			Code:   chunkBody.Code,
			Size:   chunkBody.Size,
		}
	}()
}
