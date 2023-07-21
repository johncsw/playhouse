package processor

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"log"
	"os"
	"playhouse-server/repository"
	"playhouse-server/requestbody"
	"playhouse-server/responsebody"
	"playhouse-server/util"
	"time"
)

type UploadChunkProcessor struct {
	WSConn       *websocket.Conn
	VideoID      int
	SessionID    int
	RepoFact     *repository.Factory
	NumsOfChunks int
	URLToStream  string
}

func (p *UploadChunkProcessor) preProcess(resultPipe chan<- responsebody.UploadChunkWSBody) {
	urlToStream, err := p.setUpChunkStorageDirectory()
	if err != nil || urlToStream == "" {
		errMsg := fmt.Sprintf("failed to set up chunk storage directory. videoID=%v", p.VideoID)
		if err != nil {
			errMsg = err.Error()
		}
		resultPipe <- responsebody.UploadChunkWSBody{
			Status:       responsebody.UploadChunkWSBodyStatusFailed,
			ErrorMessage: errMsg,
		}
		return
	}
	p.URLToStream = urlToStream
}

func (p *UploadChunkProcessor) postProcess() error {
	pendingChunks := p.NumsOfChunks

	transaction := p.RepoFact.NewTransaction()
	err := transaction.Execute(func(tx *gorm.DB) error {
		uploadedChunks, getChunkErr := p.RepoFact.NewChunkRepo().GetNumberOfNotUploadedChunks(p.VideoID, tx)
		if getChunkErr != nil {
			return getChunkErr
		}

		pendingChunks -= uploadedChunks

		updateVideoErr := p.RepoFact.NewVideoRepo().UpdatePendingChunks(p.VideoID, int32(pendingChunks), tx)
		if updateVideoErr != nil {
			return updateVideoErr
		}

		return nil
	})

	if err != nil {
		return err
	}

	if pendingChunks == 0 {
		go func() {
			p := TranscodeVideoProcessor{videoID: p.VideoID, repoFact: p.RepoFact}
			p.Process()
		}()
	}

	return nil
}

// The whole processing flow in a nutshell are:
// 1. set up result pipe to consume the result of saving chunks
// 2. consume messages from websocket connection
// 3. save the chunk to repository, raw data to file system and metadata to database
// 4. send the result of saving chunks to result pipe
// 5. send the result of saving chunks to client from result pipe
// 6. operations from 2-6 are done concurrently

func (p *UploadChunkProcessor) Process() {
	resultPipe, quit := p.setUpResultPipe(p.preProcess, p.postProcess, p.NumsOfChunks, time.Now())
	log.Printf("Start saving chunks. videoID=%v sessionID=%v", p.VideoID, p.SessionID)
	go func() {
		log.Printf("Start listening chunk upload. videoID=%v sessionID=%v", p.VideoID, p.SessionID)
		for {
			// the chunk metadata is sent first followed by its raw data, so we need to read the
			// message twice in order to get both of them
			b := &requestbody.UploadChunkWSBody{}
			bodyErr := p.WSConn.ReadJSON(b)
			if bodyErr != nil {
				resultPipe <- responsebody.UploadChunkWSBody{
					Status:       responsebody.UploadChunkWSBodyStatusFailed,
					ErrorMessage: bodyErr.Error(),
				}
				return
			}

			msgType, rawData, rawDataErr := p.WSConn.ReadMessage()
			if rawDataErr != nil {
				resultPipe <- responsebody.UploadChunkWSBody{
					Status:       responsebody.UploadChunkWSBodyStatusFailed,
					ErrorMessage: rawDataErr.Error(),
				}
				return
			}
			if msgType != websocket.BinaryMessage {
				resultPipe <- responsebody.UploadChunkWSBody{
					Status:       responsebody.UploadChunkWSBodyStatusFailed,
					ErrorMessage: "invalid message type",
				}
				return
			}

			b.Content = rawData
			result := p.saveToRepo(b, quit)
			go func() {
				timeout := time.After(5 * time.Second)
				select {
				case r := <-result:
					resultPipe <- r
				case <-timeout:
					resultPipe <- responsebody.UploadChunkWSBody{
						Status:       responsebody.UploadChunkWSBodyStatusFailed,
						ErrorMessage: "Get timeout when saving chunk",
					}
					return
				}
			}()
		}
	}()
}

func (p *UploadChunkProcessor) setUpResultPipe(
	preProcess func(chan<- responsebody.UploadChunkWSBody), postProcess func() error,
	iterations int, start time.Time) (chan<- responsebody.UploadChunkWSBody, <-chan bool) {
	resultPipe := make(chan responsebody.UploadChunkWSBody)
	quit := make(chan bool)
	go func() {
		defer func() {
			util.LogError(p.WSConn.Close(), "")
			log.Printf("Finished saving chunks. videoID=%v sessionID=%v elapsed=%v", p.VideoID, p.SessionID, time.Since(start))
			close(quit)
		}()

		for iterations > 0 {
			r := <-resultPipe
			if err := p.WSConn.WriteJSON(r); err != nil || r.Status == responsebody.UploadChunkWSBodyStatusFailed {
				if err != nil {
					util.LogError(err, "")
				}

				if r.Status == responsebody.UploadChunkWSBodyStatusFailed {
					jsonData, err := json.Marshal(r)
					if err != nil {
						util.LogError(err, "")
					}
					jsonStr := string(jsonData)
					log.Printf("Fail to upload chunk. result=%v", jsonStr)
				}

				return
			}
			iterations--
		}

		var writeErr error
		postProcessErr := postProcess()
		if postProcessErr != nil {
			writeErr = p.WSConn.WriteJSON(responsebody.UploadChunkWSBody{
				Status:       responsebody.UploadChunkWSBodyStatusFailed,
				ErrorMessage: postProcessErr.Error(),
			})
			if writeErr != nil {
				util.LogError(writeErr, "")
			}
			return
		}

		if writeErr = p.WSConn.WriteJSON(responsebody.UploadChunkWSBody{
			Status: responsebody.UploadChunkWSBodyStatusCompleted,
		}); writeErr != nil {
			util.LogError(writeErr, "")
		}

	}()
	preProcess(resultPipe)
	return resultPipe, quit
}

func (p *UploadChunkProcessor) saveToRepo(b *requestbody.UploadChunkWSBody, quit <-chan bool) <-chan responsebody.UploadChunkWSBody {
	result := make(chan responsebody.UploadChunkWSBody)
	go func() {
		select {
		case <-quit:
			return
		default:
			chunkRepo := p.RepoFact.NewChunkRepo()
			err := chunkRepo.SaveUploadedChunk(p.VideoID, p.URLToStream, b, nil)
			if err != nil {
				result <- responsebody.UploadChunkWSBody{
					Status:       responsebody.UploadChunkWSBodyStatusFailed,
					Code:         b.Code,
					Size:         b.Size,
					ErrorMessage: err.Error(),
				}
				return
			}
			result <- responsebody.UploadChunkWSBody{
				Status: responsebody.UploadChunkWSBodyStatusSuccess,
				Code:   b.Code,
				Size:   b.Size,
			}
		}
	}()

	return result
}

func (p *UploadChunkProcessor) setUpChunkStorageDirectory() (string, error) {
	videoRepo := p.RepoFact.NewVideoRepo()
	urlToStream, err := videoRepo.GetVideoURLToStream(p.VideoID)
	if err != nil {
		return "", err
	}
	// has not been set yet, but a record of a video does exist
	if urlToStream == "" {
		urlToStream = fmt.Sprintf("%v/%v", util.NewEnv().CHUNK_STORAGE_PATH(), p.VideoID)
		dirErr := os.Mkdir(urlToStream, 0755)
		if dirErr != nil {
			return "", dirErr
		}
		saveErr := videoRepo.SetVideoURLToStream(p.VideoID, urlToStream)
		if saveErr != nil {
			return "", saveErr
		}
		return urlToStream, nil
	}

	return urlToStream, nil
}
