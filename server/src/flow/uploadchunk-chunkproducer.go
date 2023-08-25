package flow

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"playhouse-server/request"
)

type chunkProducer struct {
	uploadChunkSupport *UploadChunkFlowSupport
	producingQueue     chan<- request.UploadChunkWebsocketBody
}

func (p *chunkProducer) start() <-chan error {
	producerErrCh := make(chan error)
	p.produce(producerErrCh)
	return producerErrCh
}

func (p *chunkProducer) produce(producerErrSender chan<- error) {
	go func() {
		defer func() {
			close(producerErrSender)
		}()
		support := p.uploadChunkSupport
		log.Printf("start listening chunk upload. videoID=%v sessionID=%v", support.VideoID, support.SessionID)

		pendingChunks := support.NumsOfChunks
		for pendingChunks > 0 {
			// the chunk metadata is first sent followed by its raw data, so we need to read the
			// message twice in order to get both
			type4Json, jsonBin, jsonErr := support.WebsocketConn.ReadMessage()
			if jsonErr != nil {
				producerErrSender <- jsonErr
				return
			}

			// ignore empty message
			if len(jsonBin) <= 0 {
				continue
			}

			if type4Json != websocket.TextMessage {
				producerErrSender <- fmt.Errorf("invalid message type")
				return
			}

			chunkBody := request.UploadChunkWebsocketBody{}
			if decodeErr := json.Unmarshal(jsonBin, &chunkBody); decodeErr != nil {
				producerErrSender <- decodeErr
				return
			}

			type4ChunkBin, chunkBin, chunkBinErr := support.WebsocketConn.ReadMessage()
			if chunkBinErr != nil {
				producerErrSender <- chunkBinErr
				return
			}

			if type4ChunkBin != websocket.BinaryMessage {
				producerErrSender <- fmt.Errorf("invalid message type")
				return
			}
			chunkBody.Content = chunkBin
			p.producingQueue <- chunkBody
			pendingChunks--
			log.Printf("chunk uploaded. videoID=%v sessionID=%v chunkCode=%v", support.VideoID, support.SessionID, chunkBody.Code)
		}
	}()
}
