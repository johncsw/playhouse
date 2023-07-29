package router

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"playhouse-server/env"
	"playhouse-server/middleware"
	"playhouse-server/repo"
	"playhouse-server/responsebody"
	"strconv"
)

func newVideoRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.AuthHandler)

	r.Group(func(r chi.Router) {
		r.Get("/streaming/{videoID}", GetManifestHandler())
		r.Get("/streaming/{videoID}/{m4sFileName}", GetStreamingContentHanlder())
		r.Get("/all", GetAllUploadedVideo())
	})
	return r
}

func GetAllUploadedVideo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videoRepo := repo.VideoRepo()
		sessionID := r.Context().Value("sessionID").(int)

		videos, err := videoRepo.GetAllUploadedVideo(sessionID)
		if err != nil {
			panic(responsebody.ResponseErr{
				Code:    http.StatusBadRequest,
				ErrBody: err,
			})
		}

		var result []map[string]any
		for _, v := range videos {
			link := fmt.Sprintf("%s/video/%d", env.CLIENT_URL(), v.ID)
			result = append(result, map[string]any{
				"name": v.Name,
				"link": link,
			})
		}

		wrapper := responsebody.Wrapper{Writer: w}
		wrapper.Status(http.StatusOK).JsonListBodyFromMap(result)
	}
}

func GetManifestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videoID, convErr := strconv.Atoi(chi.URLParam(r, "videoID"))
		if convErr != nil {
			panic(responsebody.ResponseErr{
				Code:    http.StatusBadRequest,
				ErrBody: convErr,
			})
		}

		videoRepo := repo.VideoRepo()
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

func GetStreamingContentHanlder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videoID, convErr := strconv.Atoi(chi.URLParam(r, "videoID"))
		if convErr != nil {
			panic(responsebody.ResponseErr{
				Code:    http.StatusBadRequest,
				ErrBody: convErr,
			})
		}

		videoRepo := repo.VideoRepo()
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
