package router

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"playhouse-server/auth"
	"playhouse-server/middleware"
	"playhouse-server/model"
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

	r.Group(func(r chi.Router) {
		r.Post("/register", UploadRegistrationHandler(repoFact, authenticator))
		r.Get("/chunk-code", GetChunkCodeHandler(repoFact, authenticator))
	})
	return r
}

func GetChunkCodeHandler(repoFact *repository.Factory, authenticator *auth.SessionAuthenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := authenticator.GetSessionId(r)
		videoID, convErr := strconv.Atoi(r.URL.Query().Get("video-id"))
		if convErr != nil {
			panic(util.ResponseErr{
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
			panic(util.ResponseErr{
				Code:    errCode,
				ErrBody: err,
			})
		}
		maxChunkSize := util.MaxChunkSize(videoSize)

		chunkRepo := repoFact.NewChunkRepo()
		codes, dbErr := chunkRepo.GetUnUploadedChunkCode(videoID, sessionID)
		if dbErr != nil {
			panic(util.ResponseErr{
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

func UploadRegistrationHandler(
	repoFact *repository.Factory, authenticator *auth.SessionAuthenticator) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		b := &requestbody.UploadRegistrationBody{}
		requestbody.ToRequestBody(b, r)
		if requestbody.IsNotValid(b) {
			panic(util.ResponseErr{
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
			panic(util.ResponseErr{
				Code:    http.StatusInternalServerError,
				ErrBody: err,
			})
		}

		wrapper := responsebody.Wrapper{Writer: w}
		wrapper.Status(http.StatusOK).JsonBodyFromMap(map[string]any{
			"videoID": newVideo.ID,
		})
	}
}
