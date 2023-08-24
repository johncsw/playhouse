package flow

import (
	"log"
	"playhouse-server/response"
	"playhouse-server/util"
)

type errorHandler struct {
	support             *UploadChunkSupport
	consumerErrReceiver <-chan error
	producerErrReceiver <-chan error
	flowSuccessSender   chan<- bool
}

func (h *errorHandler) listen() {
	go func() {
		support := h.support
		// producer finishes first
		producerErr := <-h.producerErrReceiver
		if producerErr != nil {
			log.Printf("Error occurred in upload producer. videoID=%v sessionID=%v err=%v", support.VideoID, support.SessionID, producerErr)
			report(support, producerErr)
			h.flowSuccessSender <- false
			return
		}

		consumerErr := <-h.consumerErrReceiver
		if consumerErr != nil {
			log.Printf("Error occurred in upload consumer. videoID=%v sessionID=%v err=%v", support.VideoID, support.SessionID, consumerErr)
			report(support, consumerErr)
			h.flowSuccessSender <- false
		}
	}()
}

func report(support *UploadChunkSupport, err error) {
	writeErr := support.WebsocketConn.WriteJSON(response.UploadChunkWebsocketBody{
		Status:       response.UploadChunkStatusFailed,
		ErrorMessage: err.Error(),
	})
	if writeErr != nil {
		util.LogError(writeErr, "")
	}
	if closeErr := support.WebsocketConn.Close(); closeErr != nil {
		util.LogError(closeErr, "")
	}
}
