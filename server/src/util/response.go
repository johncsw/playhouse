package util

import (
	"net/http"
	"playhouse-server/middleware"
)

func ReturnSuccess(body []byte, headers map[string]string, w http.ResponseWriter) {
	for h, v := range headers {
		w.Header().Set(h, v)
	}
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(body)
	if err != nil {
		panic(middleware.ResponseErr{
			Code:    http.StatusInternalServerError,
			ErrBody: err,
		})
	}
}
