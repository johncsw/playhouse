package router

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"playhouse-server/middleware"
	"playhouse-server/repository"
	"playhouse-server/responsebody"
	"strconv"
)

func newVideoRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.AuthHandler)
	repoFact := repository.NewFactory()
	r.Group(func(r chi.Router) {
		r.Get("/streaming/{videoID}", VideoStreamingHandler(repoFact))
	})
	return r
}

func VideoStreamingHandler(repoFact *repository.Factory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videoID, convErr := strconv.Atoi(chi.URLParam(r, "videoID"))
		if convErr != nil {
			panic(responsebody.ResponseErr{
				Code:    http.StatusBadRequest,
				ErrBody: convErr,
			})
		}

		videoRepo := repoFact.NewVideoRepo()
		URLToStream, isTransCoded, videoErr := videoRepo.IsVideoAvailableToStream(videoID)
		if videoErr != nil {
			panic(responsebody.ResponseErr{
				Code:    http.StatusBadRequest,
				ErrBody: videoErr,
			})
		}

		if URLToStream == "" {
			panic(responsebody.ResponseErr{
				Code:    http.StatusNotFound,
				ErrBody: errors.New("video not found"),
			})
		}

		if !isTransCoded {
			panic(responsebody.ResponseErr{
				Code:    http.StatusServiceUnavailable,
				ErrBody: errors.New("transcoding to the video is not finished yet"),
			})
		}

		manifestPath := fmt.Sprintf("%s/%d-out.mpd", URLToStream, videoID)
		http.ServeFile(w, r, manifestPath)
	}
}
