package router

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/http"
	"playhouse-server/auth"
	"playhouse-server/env"
	"playhouse-server/middleware"
	"playhouse-server/model"
	"playhouse-server/processor"
	"playhouse-server/repository"
	"playhouse-server/requestbody"
	"playhouse-server/responsebody"
	"playhouse-server/util"
	"strconv"
)

func newUploadRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.AuthHandler)

	repoFact := repository.NewFactory()
	authenticator := auth.NewSessionAuthenticator()
	webSocketUpgrader := &websocket.Upgrader{
		ReadBufferSize:  2 * util.MB,
		WriteBufferSize: 1 * util.KB,
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get("Origin") == env.CORS_ALLOWED_WEBSITE()
		},
	}

	r.Group(func(r chi.Router) {
		r.Post("/register", UploadRegistrationHandler(repoFact, authenticator))
		r.Get("/chunk-code", GetChunkCodeHandler(repoFact))
		r.Get("/chunks", ChunkUploadHandler(repoFact, authenticator, webSocketUpgrader))
	})
	return r
}

func UploadRegistrationHandler(
	repoFact *repository.Factory, authenticator *auth.SessionAuthenticator) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		b := &requestbody.UploadRegistrationBody{}
		requestbody.ToRequestBody(b, r)
		if requestbody.IsNotValid(b) {
			panic(responsebody.ResponseErr{
				Code:    http.StatusBadRequest,
				ErrBody: errors.New("not a valid request body"),
			})
		}

		var newVideo model.Video
		sessionID := authenticator.GetSessionId(r)

		transaction := repoFact.NewTransaction()
		err := transaction.Execute(func(tx *gorm.DB) error {
			videoRepo := repoFact.NewVideoRepo()

			v, verr := videoRepo.NewVideo(b, sessionID, tx)
			if verr != nil {
				return verr
			}

			chunkRepo := repoFact.NewChunkRepo()
			if cerr := chunkRepo.NewChunks(v, tx); cerr != nil {
				return cerr
			}

			newVideo = *v
			return nil
		})

		if err != nil {
			panic(responsebody.ResponseErr{
				Code:    http.StatusInternalServerError,
				ErrBody: err,
			})
		}

		wrapper := responsebody.Wrapper{Writer: w}
		videoID := strconv.Itoa(newVideo.ID)
		wrapper.Status(http.StatusCreated).JsonBodyFromMap(map[string]any{
			"videoID": videoID,
		})
	}
}

func GetChunkCodeHandler(repoFact *repository.Factory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videoID, convErr := strconv.Atoi(r.URL.Query().Get("video-id"))
		if convErr != nil {
			panic(responsebody.ResponseErr{
				Code:    http.StatusInternalServerError,
				ErrBody: convErr,
			})
		}

		videoRepo := repoFact.NewVideoRepo()
		videoSize, videoErr := videoRepo.GetVideoSize(videoID)
		if videoErr != nil {
			errCode := http.StatusInternalServerError
			err := videoErr
			if errors.Is(videoErr, gorm.ErrRecordNotFound) {
				errCode = http.StatusNotFound
				err = errors.New("video not found")
			}
			panic(responsebody.ResponseErr{
				Code:    errCode,
				ErrBody: err,
			})
		}
		maxChunkSize := util.MaxChunkSize(videoSize)

		chunkRepo := repoFact.NewChunkRepo()
		codes, dbErr := chunkRepo.GetChunkCodeByIsUploaded(videoID, false)
		if dbErr != nil {
			panic(responsebody.ResponseErr{
				Code:    http.StatusInternalServerError,
				ErrBody: dbErr,
			})
		}

		wrapper := responsebody.Wrapper{Writer: w}
		wrapper.Status(http.StatusOK).JsonBodyFromMap(
			map[string]any{
				"maxChunkSize": maxChunkSize,
				"chunkCodes":   codes,
			})
	}
}

func ChunkUploadHandler(repoFact *repository.Factory, authenticator *auth.SessionAuthenticator, webSocketUpgrader *websocket.Upgrader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videoID, convErr := strconv.Atoi(r.URL.Query().Get("video-id"))
		if convErr != nil {
			panic(responsebody.ResponseErr{
				Code:    http.StatusInternalServerError,
				ErrBody: convErr,
			})
		}
		sessionID := authenticator.GetSessionId(r)

		videoRepo := repoFact.NewVideoRepo()
		v, videoErr := videoRepo.GetPendingUploadVideo(videoID, sessionID)
		if videoErr != nil {
			errCode := http.StatusInternalServerError
			errMsg := videoErr.Error()
			if errors.Is(videoErr, gorm.ErrRecordNotFound) {
				errCode = http.StatusBadRequest
				errMsg = "not a valid video for upload"
			}
			panic(responsebody.ResponseErr{
				Code:    errCode,
				ErrBody: errors.New(errMsg),
			})
		}

		conn, socketErr := webSocketUpgrader.Upgrade(w, r, nil)
		if socketErr != nil {
			panic(responsebody.ResponseErr{
				Code:    http.StatusInternalServerError,
				ErrBody: socketErr,
			})
		}

		p := processor.UploadChunkProcessor{
			WSConn:       conn,
			VideoID:      videoID,
			RepoFact:     repoFact,
			NumsOfChunks: int(v.PendingChunks),
			SessionID:    sessionID,
		}
		p.Process()
	}
}
