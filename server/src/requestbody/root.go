package requestbody

import (
	"encoding/json"
	"io"
	"net/http"
	"playhouse-server/responsebody"
)

type RequestBody interface {
	UploadRegistrationBody | UploadChunkBody
}

type UnCheckedRequestBody interface {
	isValid() bool
}

func ToRequestBody[T RequestBody](b *T, r *http.Request) {
	bodyStreamer := r.Body
	defer func(bodyStreamer io.ReadCloser) {
		err := bodyStreamer.Close()
		if err != nil {
			panic(responsebody.ResponseErr{
				Code:    http.StatusInternalServerError,
				ErrBody: err,
			})
		}
	}(bodyStreamer)

	decoder := json.NewDecoder(bodyStreamer)

	if err := decoder.Decode(b); err != nil {
		panic(responsebody.ResponseErr{
			Code:    http.StatusInternalServerError,
			ErrBody: err,
		})
	}
}

func IsNotValid(b UnCheckedRequestBody) bool {
	return !b.isValid()
}
