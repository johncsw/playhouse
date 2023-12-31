package router

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"log"
	"net/http"
	"playhouse-server/env"
	"playhouse-server/flow"
	"playhouse-server/middleware"
	"playhouse-server/model"
	"playhouse-server/repo"
	"playhouse-server/request"
	"playhouse-server/response"
	"playhouse-server/util"
	"strconv"
	"sync/atomic"
)

func newUploadRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.AuthHandler)

	webSocketUpgrader := &websocket.Upgrader{
		ReadBufferSize:  2 * util.MB,
		WriteBufferSize: 1 * util.KB,
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get("Origin") == env.CORS_ALLOWED_WEBSITE()
		},
	}

	r.Group(func(r chi.Router) {
		r.Post("/register", uploadRegistrationHandler())
		r.Get("/chunk-code", getChunkCodeHandler())
		r.Get("/chunks", chunkUploadHandler(webSocketUpgrader))
	})
	return r
}

func uploadRegistrationHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		b := &request.UploadRegistrationBody{}
		request.ToRequestBody(b, r)
		if request.IsNotValid(b) {
			panic(response.Error{
				Code:  http.StatusBadRequest,
				Cause: errors.New("not a valid request body"),
			})
		}

		var newVideo model.Video
		sessionID := r.Context().Value("sessionID").(int)
		err := repo.NewTransaction(func(tx *gorm.DB) error {
			videoRepo := repo.VideoRepo()

			v, verr := videoRepo.NewVideo(b, sessionID, tx)
			if verr != nil {
				return verr
			}

			chunkRepo := repo.ChunkRepo()
			if cerr := chunkRepo.NewChunks(v, tx); cerr != nil {
				return cerr
			}

			newVideo = *v
			return nil
		})

		if err != nil {
			panic(response.Error{
				Code:  http.StatusInternalServerError,
				Cause: err,
			})
		}

		builder := response.Builder{Writer: w}
		videoID := strconv.Itoa(newVideo.ID)
		builder.Status(http.StatusCreated).BuildWithJson(map[string]any{
			"videoID": videoID,
		})
	}
}

func getChunkCodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videoID, convErr := strconv.Atoi(r.URL.Query().Get("video-id"))
		if convErr != nil {
			panic(response.Error{
				Code:  http.StatusInternalServerError,
				Cause: convErr,
			})
		}

		videoRepo := repo.VideoRepo()
		videoSize, videoErr := videoRepo.GetVideoSize(videoID)
		if videoErr != nil {
			errCode := http.StatusInternalServerError
			err := videoErr
			if errors.Is(videoErr, gorm.ErrRecordNotFound) {
				errCode = http.StatusNotFound
				err = errors.New("video not found")
			}
			panic(response.Error{
				Code:  errCode,
				Cause: err,
			})
		}
		maxChunkSize := util.MaxChunkSize(videoSize)

		chunkRepo := repo.ChunkRepo()
		codes, dbErr := chunkRepo.GetChunkCodeByIsUploaded(videoID, false)
		if dbErr != nil {
			panic(response.Error{
				Code:  http.StatusInternalServerError,
				Cause: dbErr,
			})
		}

		builder := response.Builder{Writer: w}
		builder.Status(http.StatusOK).BuildWithJson(
			map[string]any{
				"maxChunkSize": maxChunkSize,
				"chunkCodes":   codes,
			})
	}
}

func chunkUploadHandler(webSocketUpgrader *websocket.Upgrader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videoID, convErr := strconv.Atoi(r.URL.Query().Get("video-id"))
		if convErr != nil {
			panic(response.Error{
				Code:  http.StatusInternalServerError,
				Cause: convErr,
			})
		}

		sessionID := r.Context().Value("sessionID").(int)
		videoRepo := repo.VideoRepo()
		v, videoErr := videoRepo.GetPendingUploadVideo(videoID, sessionID)
		if videoErr != nil {
			errCode := http.StatusInternalServerError
			errMsg := videoErr.Error()
			if errors.Is(videoErr, gorm.ErrRecordNotFound) {
				errCode = http.StatusBadRequest
				errMsg = "not a valid video for upload"
			}
			panic(response.Error{
				Code:  errCode,
				Cause: errors.New(errMsg),
			})
		}

		currentUploads := atomic.LoadInt32(&flow.AllowedUploads)
		if currentUploads >= 2 {
			panic(response.Error{
				Code:  http.StatusServiceUnavailable,
				Cause: errors.New(fmt.Sprintf("fail to upload video. only two uploads are allowed at a time. videoID=%v sessionID=%v", videoID, sessionID)),
			})
		}

		conn, socketErr := webSocketUpgrader.Upgrade(w, r, nil)
		if socketErr != nil {
			panic(response.Error{
				Code:  http.StatusInternalServerError,
				Cause: socketErr,
			})
		}

		atomic.AddInt32(&flow.AllowedUploads, 1)
		go func() {
			defer atomic.AddInt32(&flow.AllowedUploads, -1)
			success := <-flow.UploadChunk(&flow.UploadChunkFlowSupport{
				WebsocketConn: conn,
				VideoID:       videoID,
				VideoSize:     v.Size,
				NumsOfChunks:  int(v.PendingChunks),
				SessionID:     sessionID,
			})
			newPendingChunks := flow.UpdatePendingChunks(videoID, int(v.PendingChunks))
			if success && newPendingChunks == 0 {
				flow.TranscodeVideo(videoID, sessionID)
				markErr := repo.VideoRepo().UpdateVideoAsTranscodeComplete(videoID)
				if markErr != nil {
					log.Printf("failed to mark video as transcode complete. videoID=%d sessionID=%d\n", videoID, sessionID)
				}
			}
		}()
	}
}
