package flow

import (
	"fmt"
	"log"
	"playhouse-server/response"
	"playhouse-server/util"
	"time"
)

type successHandler struct {
	support              *UploadChunkSupport
	chunkSuccessReceiver <-chan response.UploadChunkWebsocketBody
	flowSuccessSender    chan<- bool
}

func (h *successHandler) listen() {
	go func() {
		pendingChunks := h.support.NumsOfChunks
		for pendingChunks > 0 {
			select {
			case resBody := <-h.chunkSuccessReceiver:
				if writeErr := h.support.WebsocketConn.WriteJSON(resBody); writeErr != nil {
					report(h.support, writeErr)
					return
				}
				pendingChunks--
			case <-time.After(5 * time.Second):
				report(h.support, fmt.Errorf("timeout sending success response. videoID=%v sessionID=%v", h.support.VideoID, h.support.SessionID))
				return
			}
		}

		if completeErr := h.support.WebsocketConn.WriteJSON(response.UploadChunkWebsocketBody{
			Status: response.UploadChunkStatusCompleted,
		}); completeErr != nil {
			report(h.support, completeErr)
			return
		}

		if err := h.support.WebsocketConn.Close(); err != nil {
			util.LogError(err, "")
		}

		h.flowSuccessSender <- true
		log.Printf("upload completed. videoID=%v sessionID=%v", h.support.VideoID, h.support.SessionID)
	}()
}
