package flow

import (
	"fmt"
	"github.com/gorilla/websocket"
	"playhouse-server/repo"
	"playhouse-server/request"
)

type UploadChunkFlowSupport struct {
	WebsocketConn     *websocket.Conn
	VideoID           int
	SessionID         int
	NumsOfChunks      int
	ChunkSavingDirURL string
}

func UploadChunk(support *UploadChunkFlowSupport) <-chan bool {
	flowSuccess := make(chan bool)
	go func() {
		chunkSavingDirURL, prepareErr := prepareChunkSavingDirURL(support)
		if prepareErr != nil {
			report(support, prepareErr)
			return
		}

		support.ChunkSavingDirURL = chunkSavingDirURL

		chunkQueue := make(chan request.UploadChunkWebsocketBody)
		// produce chunks by receiving them from websocket
		producer := chunkProducer{
			producingQueue:     chunkQueue,
			uploadChunkSupport: support,
		}
		producerFinalErrReceiver := producer.start()

		consumer := chunkConsumer{
			consumingQueue:     chunkQueue,
			uploadChunkSupport: support,
		}
		chunkSuccessReceiver, consumerFinalErrReceiver := consumer.start()

		errorHandler := errorHandler{
			support:             support,
			producerErrReceiver: producerFinalErrReceiver,
			consumerErrReceiver: consumerFinalErrReceiver,
			flowSuccessSender:   flowSuccess,
		}
		errorHandler.listen()

		successHandler := successHandler{
			support:              support,
			chunkSuccessReceiver: chunkSuccessReceiver,
			flowSuccessSender:    flowSuccess,
		}
		successHandler.listen()
	}()
	return flowSuccess
}

func prepareChunkSavingDirURL(support *UploadChunkFlowSupport) (string, error) {
	existingURL, dirErr := repo.VideoRepo().GetChunkSavingDirURL(support.VideoID)
	if dirErr != nil {
		return "", dirErr
	}
	if existingURL != "" {
		return existingURL, nil
	}

	newURL, createErr := repo.VideoRepo().CreateChunkSavingDir(support.VideoID)
	if createErr != nil {
		return "", createErr
	}
	if newURL == "" {
		return "", fmt.Errorf("fail to create chunk saving dir. videoID=%v sessionID=%v", support.VideoID, support.SessionID)
	}

	return newURL, nil
}
