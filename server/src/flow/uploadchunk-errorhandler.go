package flow

import (
	"github.com/gorilla/websocket"
	"log"
	"playhouse-server/response"
	"playhouse-server/util"
)

type errorHandler struct {
	support             *UploadChunkFlowSupport
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
			log.Printf("error occurred in upload producer. videoID=%v sessionID=%v err=%v", support.VideoID, support.SessionID, producerErr)
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

func report(support *UploadChunkFlowSupport, err error) {
	util.LogError(err, "")
	writeErr := support.WebsocketConn.WriteJSON(response.UploadChunkWebsocketBody{
		Status:       response.UploadChunkStatusFailed,
		ErrorMessage: err.Error(),
	})
	if writeErr != nil {
		util.LogError(writeErr, "")
	}

	msgErr := support.WebsocketConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if msgErr != nil {
		util.LogError(msgErr, "")
		return
	}

	if closeErr := support.WebsocketConn.Close(); closeErr != nil {
		util.LogError(closeErr, "")
	}
}
