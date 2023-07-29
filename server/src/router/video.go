package router

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"playhouse-server/env"
	"playhouse-server/middleware"
	"playhouse-server/repo"
	"playhouse-server/response"
	"strconv"
)

func newVideoRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.AuthHandler)

	r.Group(func(r chi.Router) {
		r.Get("/streaming/{videoID}", getManifestHandler())
		r.Get("/streaming/{videoID}/{m4sFileName}", getStreamingContentHandler())
		r.Get("/all", getAllUploadedVideo())
	})
	return r
}

func getAllUploadedVideo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videoRepo := repo.VideoRepo()
		sessionID := r.Context().Value("sessionID").(int)

		videos, err := videoRepo.GetAllUploadedVideo(sessionID)
		if err != nil {
			panic(response.Error{
				Code:  http.StatusBadRequest,
				Cause: err,
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

		builder := response.Builder{Writer: w}
		builder.Status(http.StatusOK).BuildWithJsonList(result)
	}
}

func getManifestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videoID, convErr := strconv.Atoi(chi.URLParam(r, "videoID"))
		if convErr != nil {
			panic(response.Error{
				Code:  http.StatusBadRequest,
				Cause: convErr,
			})
		}

		videoRepo := repo.VideoRepo()
		URLToStream, isTransCoded, videoErr := videoRepo.IsVideoAvailableToStream(videoID)
		if videoErr != nil {
			panic(response.Error{
				Code:  http.StatusBadRequest,
				Cause: videoErr,
			})
		}

		if URLToStream == "" {
			panic(response.Error{
				Code:  http.StatusNotFound,
				Cause: errors.New("video not found"),
			})
		}

		if !isTransCoded {
			panic(response.Error{
				Code:  http.StatusServiceUnavailable,
				Cause: errors.New("transcoding to the video is not finished yet"),
			})
		}

		manifestPath := fmt.Sprintf("%s/%d-out.mpd", URLToStream, videoID)
		http.ServeFile(w, r, manifestPath)

	}
}

func getStreamingContentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videoID, convErr := strconv.Atoi(chi.URLParam(r, "videoID"))
		if convErr != nil {
			panic(response.Error{
				Code:  http.StatusBadRequest,
				Cause: convErr,
			})
		}

		videoRepo := repo.VideoRepo()
		URLToStream, isTransCoded, videoErr := videoRepo.IsVideoAvailableToStream(videoID)
		if videoErr != nil {
			panic(response.Error{
				Code:  http.StatusBadRequest,
				Cause: videoErr,
			})
		}

		if URLToStream == "" {
			panic(response.Error{
				Code:  http.StatusNotFound,
				Cause: errors.New("video not found"),
			})
		}

		if !isTransCoded {
			panic(response.Error{
				Code:  http.StatusServiceUnavailable,
				Cause: errors.New("transcoding to the video is not finished yet"),
			})
		}

		m4sFileName := chi.URLParam(r, "m4sFileName")
		m4sFilePath := fmt.Sprintf("%s/%s", URLToStream, m4sFileName)
		http.ServeFile(w, r, m4sFilePath)
	}
}
