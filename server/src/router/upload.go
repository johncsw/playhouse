package router

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/http"
	"playhouse-server/env"
	"playhouse-server/middleware"
	"playhouse-server/model"
	"playhouse-server/processor"
	"playhouse-server/repo"
	"playhouse-server/request"
	"playhouse-server/response"
	"playhouse-server/util"
	"strconv"
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
		r.Post("/register", UploadRegistrationHandler())
		r.Get("/chunk-code", GetChunkCodeHandler())
		r.Get("/chunks", ChunkUploadHandler(webSocketUpgrader))
	})
	return r
}

func UploadRegistrationHandler() http.HandlerFunc {

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

func GetChunkCodeHandler() http.HandlerFunc {
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

func ChunkUploadHandler(webSocketUpgrader *websocket.Upgrader) http.HandlerFunc {
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

		conn, socketErr := webSocketUpgrader.Upgrade(w, r, nil)
		if socketErr != nil {
			panic(response.Error{
				Code:  http.StatusInternalServerError,
				Cause: socketErr,
			})
		}

		p := processor.UploadChunkProcessor{
			WSConn:       conn,
			VideoID:      videoID,
			NumsOfChunks: int(v.PendingChunks),
			SessionID:    sessionID,
		}
		p.Process()
	}
}
