package util

import (
	"net/http"
)

type Wrapper struct{}

var (
	wrapper *Wrapper
)

func NewWrapper() *Wrapper {
	if wrapper == nil {
		wrapper = &Wrapper{}
	}

	return wrapper
}

func (Wrapper) SuccessfulResponse(body []byte, headers map[string]string, w http.ResponseWriter) {
	for h, v := range headers {
		w.Header().Set(h, v)
	}
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(body)
	if err != nil {
		panic(ResponseErr{
			Code:    http.StatusInternalServerError,
			ErrBody: err,
		})
	}
}
