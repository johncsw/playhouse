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
		r.Get("/streaming/{videoID}", GetManifestHandler(repoFact))
		r.Get("/streaming/{videoID}/{m4sFileName}", GetStreamingContentHanlder(repoFact))
	})
	return r
}

func GetManifestHandler(repoFact *repository.Factory) http.HandlerFunc {
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

func GetStreamingContentHanlder(repoFact *repository.Factory) http.HandlerFunc {
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

		m4sFileName := chi.URLParam(r, "m4sFileName")
		m4sFilePath := fmt.Sprintf("%s/%s", URLToStream, m4sFileName)
		http.ServeFile(w, r, m4sFilePath)
	}
}
